/**
 * Copyright 2015 @ z3q.net.
 * name : return
 * author : jarryliu
 * date : 2016-07-17 17:29
 * description :
 * history :
 */
package afterSales

import (
	"errors"
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/tmp"
	"math"
)

var _ afterSales.IAfterSalesOrder = new(returnOrderImpl)
var _ afterSales.IReturnOrder = new(returnOrderImpl)

type returnOrderImpl struct {
	*afterSalesOrderImpl
	refValue    *afterSales.ReturnOrder
	memberRepo  member.IMemberRepo
	paymentRepo payment.IPaymentRepo
}

func newReturnOrderImpl(v *afterSalesOrderImpl, memberRepo member.IMemberRepo,
	paymentRepo payment.IPaymentRepo) *returnOrderImpl {
	if v.value.Type != afterSales.TypeReturn {
		panic(errors.New("售后单类型不是退货单"))
	}
	return &returnOrderImpl{
		afterSalesOrderImpl: v,
		memberRepo:          memberRepo,
		paymentRepo:         paymentRepo,
	}
}

func (r *returnOrderImpl) getValue() *afterSales.ReturnOrder {
	if r.refValue == nil {
		if r.GetDomainId() <= 0 {
			panic(errors.New("退货单还未提交"))
		}
		v := &afterSales.ReturnOrder{}
		if tmp.Db().GetOrm().Get(r.GetDomainId(), v) != nil {
			panic(errors.New("退货单不存在"))
		}
		r.refValue = v
	}
	return r.refValue
}

// 获取售后单数据
func (r *returnOrderImpl) Value() afterSales.AfterSalesOrder {
	v := r.afterSalesOrderImpl.Value()
	v2 := r.getValue()
	v.Data = *v2
	return v
}

// 保存
func (r *returnOrderImpl) saveReturnOrder() error {
	_, err := orm.Save(tmp.Db().GetOrm(), r.refValue, int(r.GetDomainId()))
	return err
}

// 设置要退回货物信息
func (r *returnOrderImpl) SetItem(snapshotId int64, quantity int32) error {
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

// 提交售后申请
func (r *returnOrderImpl) Submit() (int32, error) {
	//o := r.GetOrder()
	//if o.GetValue().State == order.StatCompleted {
	//    return 0, afterSales.ErrReturnAfterReceived
	//}
	id, err := r.afterSalesOrderImpl.Submit()
	// 提交退货单
	if err == nil {
		// 锁定退货数量
		err = r.GetOrder().Return(r.value.SnapshotId, r.value.Quantity)
		if err == nil {
			// 生成退货单
			err = r.submitReturnOrder()
		}
	}
	return id, err
}

// 提交退货单
func (r *returnOrderImpl) submitReturnOrder() (err error) {
	r.refValue = &afterSales.ReturnOrder{
		Id:       r.afterSalesOrderImpl.GetDomainId(),
		IsRefund: 0,
	}
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
func (r *returnOrderImpl) Cancel() error {
	err := r.afterSalesOrderImpl.Cancel()
	if err == nil {
		// 撤销退货数量
		err = r.GetOrder().RevertReturn(r.value.SnapshotId, r.value.Quantity)
	}
	return err
}

// 退回申请
func (r *returnOrderImpl) Reject(remark string) error {
	err := r.afterSalesOrderImpl.Reject(remark)
	if err == nil {
		// 撤销退货数量
		err = r.GetOrder().RevertReturn(r.value.SnapshotId, r.value.Quantity)
	}
	return err
}

// 完成退货
func (r *returnOrderImpl) Process() error {
	err := r.afterSalesOrderImpl.Process()
	if err == nil {
		err = r.handleReturn()
	}
	return err
}

// 处理退货
func (r *returnOrderImpl) handleReturn() error {

	//todo:??添加库存,或计入残次品

	v := r.getValue()
	if v.IsRefund == 1 {
		return nil
	}
	v.IsRefund = 1
	err := r.saveReturnOrder()
	if err == nil {
		err = r.backAmount(v.Amount)
	}
	return err
}

// 退款
func (r *returnOrderImpl) backAmount(amount float32) error {
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
			o.OrderNo, amount,
			member.DefaultRelateUser)
	}
	//原路退回
	pv := po.GetValue()
	if pv.BalanceDiscount > 0 {
		if err := acc.Refund(member.AccountBalance, member.ChargeByRefund, "订单退款",
			o.OrderNo, pv.BalanceDiscount, member.DefaultRelateUser); err == nil {
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
