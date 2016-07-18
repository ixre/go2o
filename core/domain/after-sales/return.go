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
	"go2o/core/domain/tmp"
)

var _ afterSales.IAfterSalesOrder = new(returnOrderImpl)
var _ afterSales.IReturnOrder = new(returnOrderImpl)

type returnOrderImpl struct {
	*afterSalesOrderImpl
	_returnValue *afterSales.ReturnOrder
	_memberRep   member.IMemberRep
}

func newReturnOrderImpl(v *afterSalesOrderImpl,
	memberRep member.IMemberRep) *returnOrderImpl {
	if v._value.Type != afterSales.TypeReturn {
		panic(errors.New("售后单类型不是退货单"))
	}
	return &returnOrderImpl{
		afterSalesOrderImpl: v,
		_memberRep:          memberRep,
	}
}

func (r *returnOrderImpl) getValue() *afterSales.ReturnOrder {
	if r._returnValue == nil {
		if r.GetDomainId() <= 0 {
			panic(errors.New("退货单还未提交"))
		}
		v := &afterSales.ReturnOrder{}
		if tmp.Db().GetOrm().Get(r.GetDomainId(), v) == nil {
			r._returnValue = v
		}
		panic(errors.New("退货单不存在"))
	}
	return r._returnValue
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
	_, err := orm.Save(tmp.Db().GetOrm(), r._returnValue, r.GetDomainId())
	return err
}

// 提交售后申请
func (r *returnOrderImpl) Submit() error {
	err := r.afterSalesOrderImpl.Submit()
	// 提交退货单
	if err != nil {
		return err
	}
	r._returnValue = &afterSales.ReturnOrder{
		Id:       r.afterSalesOrderImpl.GetDomainId(),
		IsRefund: 0,
	}
	o := r.GetOrder()
	for _, v := range o.Items() {
		if v.SnapshotId == r._value.SnapshotId {
			r._returnValue.Amount = v.FinalAmount * float32(r._value.Quantity)
			break
		}
	}
	if r._returnValue.Amount <= 0 {
		panic(errors.New("退货单货款为零"))
	}
	_, err = orm.Save(tmp.Db().GetOrm(), r._returnValue, 0)
	return err
}

// 完成退货
func (r *returnOrderImpl) Confirm() error {
	err := r.afterSalesOrderImpl.Confirm()
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
	mm := r._memberRep.GetMember(r._value.BuyerId)
	if mm == nil {
		return member.ErrNoSuchMember
	}
	acc := mm.GetAccount()
	return acc.ChargeBalance(member.TypeBalanceOrderRefund, "订单退款",
		o.OrderNo, amount)
}
