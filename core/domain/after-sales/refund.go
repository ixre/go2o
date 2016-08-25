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
	"go2o/core/domain/tmp"
)

var _ afterSales.IRefundOrder = new(refundOrderImpl)

type refundOrderImpl struct {
	*afterSalesOrderImpl
	_refValue  *afterSales.RefundOrder
	_memberRep member.IMemberRep
}

func newRefundOrder(v *afterSalesOrderImpl, memberRep member.IMemberRep) *refundOrderImpl {
	if v._value.Type != afterSales.TypeRefund {
		panic(errors.New("售后单类型不是退款单"))
	}
	return &refundOrderImpl{
		afterSalesOrderImpl: v,
		_memberRep:          memberRep,
	}
}

func (r *refundOrderImpl) getValue() *afterSales.RefundOrder {
	if r._refValue == nil {
		if r.GetDomainId() <= 0 {
			panic(errors.New("退款单还未提交"))
		}
		v := &afterSales.RefundOrder{}
		if tmp.Db().GetOrm().Get(r.GetDomainId(), v) != nil {
			panic(errors.New("退款单不存在"))
		}
		r._refValue = v
	}
	return r._refValue
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
	_, err := orm.Save(tmp.Db().GetOrm(), r._refValue, r.GetDomainId())
	return err
}

// 设置要退回货物信息
func (r *refundOrderImpl) SetItem(snapshotId int, quantity int) error {
	o := r.GetOrder()
	for _, v := range o.Items() {
		if v.SnapshotId == snapshotId {
			// 判断是否超过数量,减去已退货数量
			if v.Quantity-v.ReturnQuantity < quantity {
				return afterSales.ErrOutOfQuantity
			}
			// 设置退回商品
			r._value.SnapshotId = snapshotId
			r._value.Quantity = quantity
			return nil
		}
	}
	return afterSales.ErrNoSuchOrderItem
}

// 提交退款申请
func (r *refundOrderImpl) Submit() (int, error) {
	o := r.GetOrder()
	if o.GetValue().State >= order.StatShipped {
		return 0, afterSales.ErrRefundAfterShipped
	}
	id, err := r.afterSalesOrderImpl.Submit()
	// 提交退款单
	if err == nil {
		// 锁定退货数量
		err = o.Return(r._value.SnapshotId, r._value.Quantity)
		if err == nil {
			// 生成退款单
			err = r.submitRefundOrder()
		}
	}
	return id, err
}

// 提交退款单
func (r *refundOrderImpl) submitRefundOrder() (err error) {
	r._refValue = &afterSales.RefundOrder{
		Id:       r.afterSalesOrderImpl.GetDomainId(),
		IsRefund: 0,
	}
	// 计算退款金额
	o := r.GetOrder()
	for _, v := range o.Items() {
		if v.SnapshotId == r._value.SnapshotId {
			price := v.FinalAmount / float32(v.Quantity) // 计算单价
			r._refValue.Amount = price * float32(r._value.Quantity)
			break
		}
	}
	if r._refValue.Amount <= 0 {
		return afterSales.ErrOrderAmount
	}
	_, err = orm.Save(tmp.Db().GetOrm(), r._refValue, 0)
	return err
}

// 取消申请
func (r *refundOrderImpl) Cancel() error {
	err := r.afterSalesOrderImpl.Cancel()
	if err == nil {
		// 撤销退货数量
		err = r.GetOrder().RevertReturn(r._value.SnapshotId, r._value.Quantity)
	}
	return err
}

// 退回申请
func (r *refundOrderImpl) Reject(remark string) error {
	err := r.afterSalesOrderImpl.Reject(remark)
	if err == nil {
		// 撤销退货数量
		err = r.GetOrder().RevertReturn(r._value.SnapshotId, r._value.Quantity)
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
	mm := r._memberRep.GetMember(r._value.BuyerId)
	if mm == nil {
		return member.ErrNoSuchMember
	}
	acc := mm.GetAccount()
	return acc.ChargeForBalance(member.ChargeByRefund, "订单退款",
		o.OrderNo, amount, member.DefaultRelateUser)
}
