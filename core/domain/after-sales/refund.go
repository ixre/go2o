/**
 * Copyright 2015 @ z3q.net.
 * name : refund
 * author : jarryliu
 * date : 2016-07-17 11:43
 * description :
 * history :
 */
package afterSales

import (
	"errors"
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/tmp"
	"math"
)

var _ afterSales.IRefundOrder = new(refundOrderImpl)

type refundOrderImpl struct {
	*afterSalesOrderImpl
	refValue    *afterSales.RefundOrder
	memberRepo  member.IMemberRepo
	paymentRepo payment.IPaymentRepo
}

func newRefundOrder(v *afterSalesOrderImpl, memberRepo member.IMemberRepo,
	paymentRepo payment.IPaymentRepo) *refundOrderImpl {
	if v.value.Type != afterSales.TypeRefund {
		panic(errors.New("售后单类型不是退款单"))
	}
	return &refundOrderImpl{
		afterSalesOrderImpl: v,
		memberRepo:          memberRepo,
		paymentRepo:         paymentRepo,
	}
}

func (r *refundOrderImpl) getValue() *afterSales.RefundOrder {
	if r.refValue == nil {
		if r.GetDomainId() <= 0 {
			panic(errors.New("退款单还未提交"))
		}
		v := &afterSales.RefundOrder{}
		if tmp.Db().GetOrm().Get(r.GetDomainId(), v) != nil {
			panic(errors.New("退款单不存在"))
		}
		r.refValue = v
	}
	return r.refValue
}

// 获取售后单数据
func (r *refundOrderImpl) Value() afterSales.AfterSalesOrder {
	v := r.afterSalesOrderImpl.Value()
	v2 := r.getValue()
	v.Data = *v2
	return v
}

// 保存
func (r *refundOrderImpl) saveRefundOrder() error {
	_, err := orm.Save(tmp.Db().GetOrm(), r.refValue, int(r.GetDomainId()))
	return err
}

// 设置要退回货物信息
func (r *refundOrderImpl) SetItem(snapshotId int64, quantity int32) error {
	o := r.GetOrder()
	for _, v := range o.Items() {
		if v.SnapshotId == snapshotId {
			// 判断是否超过数量,减去已退货数量
			if v.Quantity-v.ReturnQuantity < quantity {
				return afterSales.ErrOutOfQuantity
			}
			// 设置退回商品
			r.value.SnapshotId = snapshotId
			r.value.Quantity = quantity
			return nil
		}
	}
	return afterSales.ErrNoSuchOrderItem
}

// 提交退款申请
func (r *refundOrderImpl) Submit() (int32, error) {
	o := r.GetOrder()
	if o.GetValue().State >= order.StatShipped {
		return 0, afterSales.ErrRefundAfterShipped
	}
	id, err := r.afterSalesOrderImpl.Submit()
	// 提交退款单
	if err == nil {
		// 锁定退货数量
		err = o.Return(r.value.SnapshotId, r.value.Quantity)
		if err == nil {
			// 生成退款单
			err = r.submitRefundOrder()
		}
	}
	return id, err
}

// 提交退款单
func (r *refundOrderImpl) submitRefundOrder() (err error) {
	r.refValue = &afterSales.RefundOrder{
		Id:       r.afterSalesOrderImpl.GetDomainId(),
		IsRefund: 0,
	}
	// 计算退款金额
	o := r.GetOrder()
	for _, v := range o.Items() {
		if v.SnapshotId == r.value.SnapshotId {
			price := v.FinalAmount / float32(v.Quantity) // 计算单价
			r.refValue.Amount = price * float32(r.value.Quantity)
			break
		}
	}
	if r.refValue.Amount <= 0 || math.IsNaN(float64(r.refValue.Amount)) {
		return afterSales.ErrOrderAmount
	}
	_, err = orm.Save(tmp.Db().GetOrm(), r.refValue, 0)
	return err
}

// 取消申请
func (r *refundOrderImpl) Cancel() error {
	err := r.afterSalesOrderImpl.Cancel()
	if err == nil {
		// 撤销退货数量
		err = r.GetOrder().RevertReturn(r.value.SnapshotId, r.value.Quantity)
	}
	return err
}

// 退回申请
func (r *refundOrderImpl) Reject(remark string) error {
	err := r.afterSalesOrderImpl.Reject(remark)
	if err == nil {
		// 撤销退货数量
		err = r.GetOrder().RevertReturn(r.value.SnapshotId, r.value.Quantity)
	}
	return err
}

// 完成退款
func (r *refundOrderImpl) Process() error {
	err := r.afterSalesOrderImpl.Process()
	if err == nil {
		err = r.handleReturn()
	}
	return err
}

// 处理退款
func (r *refundOrderImpl) handleReturn() error {
	v := r.getValue()
	if v.IsRefund == 1 {
		return nil
	}
	v.IsRefund = 1
	err := r.saveRefundOrder()
	if err == nil {
		err = r.backAmount(v.Amount)
	}
	return err
}

// 退款
func (r *refundOrderImpl) backAmount(amount float32) error {
	o := r.GetOrder().GetValue()
	mm := r.memberRepo.GetMember(r.value.BuyerId)
	if mm == nil {
		return member.ErrNoSuchMember
	}
	acc := mm.GetAccount()
	//支付单与父订单关联。多个子订单合并付款
	po := r.paymentRepo.GetPaymentBySalesOrderId(o.OrderId)
	//如果支付单已清理数据，则全部退回到余额
	if po == nil {
		return acc.Refund(member.AccountBalance,
			member.KindBalanceRefund, "订单退款",
			o.OrderNo, amount, member.DefaultRelateUser)
	}
	//原路退回
	pv := po.GetValue()
	if pv.BalanceDiscount > 0 {
		if err := acc.Refund(member.AccountBalance,
			member.KindBalanceRefund,
			"订单退款", o.OrderNo, pv.BalanceDiscount,
			member.DefaultRelateUser); err == nil {
			amount -= pv.BalanceDiscount
		}
	}
	//退积分
	if pv.IntegralDiscount > 0 {
		//todo : 退换积分,暂时积分抵扣的不退款
	}
	//多退少补
	if pv.FinalFee > amount {
		amount = pv.FinalFee
	}
	//退到钱包账户
	if pv.PaymentSign == payment.SignWalletAccount {
		return acc.Refund(member.AccountWallet,
			member.KindWalletPaymentRefund,
			"订单退款", o.OrderNo, amount,
			member.DefaultRelateUser)
	}
	//原路退回，暂时不实现。直接退到钱包账户
	return acc.Refund(member.AccountWallet,
		member.KindWalletPaymentRefund,
		"订单退款", o.OrderNo, amount,
		member.DefaultRelateUser)
}
