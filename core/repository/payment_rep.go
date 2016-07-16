/**
 * Copyright 2015 @ z3q.net.
 * name : payment_rep.go
 * author : jarryliu
 * date : 2016-07-03 12:52
 * description :
 * history :
 */
package repository

import (
	"github.com/jsix/gof/db"
	"go2o/core"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/valueobject"
	payImpl "go2o/core/domain/payment"
	"go2o/core/variable"
)

var _ payment.IPaymentRep = new(paymentRep)

type paymentRep struct {
	db.Connector
	*payImpl.PaymentRepBase
	_memberRep member.IMemberRep
	_valRep    valueobject.IValueRep
	_orderRep  order.IOrderRep
}

func NewPaymentRep(conn db.Connector, mmRep member.IMemberRep,
	orderRep order.IOrderRep, valRep valueobject.IValueRep) payment.IPaymentRep {
	return &paymentRep{
		Connector:  conn,
		_memberRep: mmRep,
		_valRep:    valRep,
		_orderRep:  orderRep,
	}
}

// 根据编号获取支付单
func (this *paymentRep) GetPaymentOrder(
	id int) payment.IPaymentOrder {
	e := &payment.PaymentOrderBean{}
	if this.Connector.GetOrm().Get(id, e) == nil {
		return this.CreatePaymentOrder(e)
	}
	return nil
}

// 根据支付单号获取支付单
func (this *paymentRep) GetPaymentOrderByNo(
	paymentNo string) payment.IPaymentOrder {
	e := &payment.PaymentOrderBean{}
	if this.Connector.GetOrm().GetBy(e, "trade_no=?", paymentNo) == nil {
		return this.CreatePaymentOrder(e)
	}
	return nil
}

// 根据订单号获取支付单
func (p *paymentRep) GetPaymentBySalesOrderId(orderId int) payment.IPaymentOrder {
	e := &payment.PaymentOrderBean{}
	if p.Connector.GetOrm().GetBy(e, "order_id=?", orderId) == nil {
		return p.CreatePaymentOrder(e)
	}
	return nil
}

// 创建支付单
func (this *paymentRep) CreatePaymentOrder(
	p *payment.PaymentOrderBean) payment.IPaymentOrder {
	return this.PaymentRepBase.CreatePaymentOrder(p, this,
		this._memberRep, this._orderRep.Manager(), this._valRep)
}

// 保存支付单
func (this *paymentRep) SavePaymentOrder(
	v *payment.PaymentOrderBean) (id int, err error) {
	orm := this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		var id64 int64
		_, id64, err = orm.Save(nil, v)
		v.Id = int(id64)
	}
	return v.Id, err
}

// 通知支付单完成
func (this *paymentRep) NotifyPaymentFinish(paymentOrderId int) error {
	rc := core.GetRedisConn()
	defer rc.Close()
	_, err := rc.Do("RPUSH", variable.KvPaymentOrderFinishQueue, paymentOrderId)
	return err
}
