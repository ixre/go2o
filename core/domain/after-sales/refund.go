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
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/tmp"
	"time"
)

var _ afterSales.IRefundOrder = new(refundOrderImpl)

type refundOrderImpl struct {
	_rep      afterSales.IAfterSalesRep
	_value    *afterSales.RefundOrder
	_order    order.ISubOrder
	_orderRep order.IOrderRep
}

func NewRefundOrder(v *afterSales.RefundOrder, rep afterSales.IAfterSalesRep,
	orderRep order.IOrderRep) *refundOrderImpl {
	return &refundOrderImpl{
		_rep:      rep,
		_value:    v,
		_orderRep: orderRep,
	}
}

// 获取领域对象编号
func (r *refundOrderImpl) GetDomainId() int {
	return r._value.Id
}

// 获取值
func (r *refundOrderImpl) GetValue() afterSales.RefundOrder {
	return *r._value
}

// 获取订单
func (r *refundOrderImpl) Order() order.ISubOrder {
	if r._order == nil {
		r._order = r._orderRep.Manager().GetSubOrder(r._value.OrderId)
	}
	return r._order
}

func (r *refundOrderImpl) save() error {
	r._value.UpdateTime = time.Now().Unix()
	id, err := orm.Save(tmp.Db().GetOrm(), r._value, r.GetDomainId())
	r._value.Id = id
	return err
}

// 提交退款申请
func (r *refundOrderImpl) Submit() error {
	if r.GetDomainId() > 0 {
		return nil
	}
	err := r.save()
	if err == nil {
		err = r.Order().SubmitRefund(r._value.Reason)
	}
	return err
}

// 取消申请退款
func (r *refundOrderImpl) Cancel() error {
	r._value.State = afterSales.RefundStatCancelled
	err := r.save()
	if err == nil {
		err = r.Order().CancelRefund()
	}
	return err
}

// 拒绝退款
func (r *refundOrderImpl) Decline(remark string) error {
	r._value.State = afterSales.RefundStatVendorDecline
	r._value.VendorRemark = remark
	err := r.save()
	if err == nil {
		err = r.Order().Decline(remark)
	}
	return err
}

// 同意退款
func (r *refundOrderImpl) Agree() error {
	r._value.State = afterSales.RefundStatAwaittingConfirm
	return r.save()
}

// 确认退款
func (r *refundOrderImpl) Confirm() error {
	r._value.State = afterSales.RefundStatCompleted
	err := r.save()
	if err == nil {
		err = r.Order().Refund()
	}
	return err
}

// 申请调解
func (r *refundOrderImpl) RequestIntercede() error {
	r._value.State = afterSales.RefundStatIntercede
	return r.save()
}

// 调解后直接操作退款或退回
func (r *refundOrderImpl) IntercedeHandle(pass bool, remark string) error {
	r._value.Remark = remark
	if pass {
		return r.Confirm()
	}
	return nil
}
