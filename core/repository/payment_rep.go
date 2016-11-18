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
	"fmt"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"github.com/jsix/gof/storage"
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
	Storage storage.Interface
	*payImpl.PaymentRepBase
	_memberRep member.IMemberRep
	_valRep    valueobject.IValueRep
	_orderRep  order.IOrderRep
}

func NewPaymentRep(sto storage.Interface, conn db.Connector, mmRep member.IMemberRep,
	orderRep order.IOrderRep, valRep valueobject.IValueRep) payment.IPaymentRep {
	return &paymentRep{
		Storage:    sto,
		Connector:  conn,
		_memberRep: mmRep,
		_valRep:    valRep,
		_orderRep:  orderRep,
	}
}

// 根据订单号获取支付单
func (p *paymentRep) GetPaymentBySalesOrderId(orderId int64) payment.IPaymentOrder {
	e := &payment.PaymentOrder{}
	if p.Connector.GetOrm().GetBy(e, "order_id=?", orderId) == nil {
		return p.CreatePaymentOrder(e)
	}
	return nil
}

func (p *paymentRep) getPaymentOrderCk(id int64) string {
	return fmt.Sprintf("go2o:rep:pay:order:%d", id)
}
func (p *paymentRep) getPaymentOrderCkByNo(orderNO string) string {
	return fmt.Sprintf("go2o:rep:pay:order:%s", orderNO)
}

// 根据编号获取支付单
func (p *paymentRep) GetPaymentOrder(id int64) payment.IPaymentOrder {
	if id <= 0 {
		return nil
	}
	e := &payment.PaymentOrder{}
	k := p.getPaymentOrderCk(id)
	if err := p.Storage.Get(k, &e); err != nil {
		if p.Connector.GetOrm().Get(id, e) != nil {
			return nil
		}
		p.Storage.SetExpire(k, *e, DefaultCacheSeconds)
	}
	return p.CreatePaymentOrder(e)
}

// 根据支付单号获取支付单
func (p *paymentRep) GetPaymentOrderByNo(paymentNo string) payment.IPaymentOrder {
	k := p.getPaymentOrderCkByNo(paymentNo)
	id, err := p.Storage.GetInt64(k)
	if err != nil {
		p.ExecScalar("SELECT id FROM pay_order where trade_no=?", &id, paymentNo)
		if id == 0 {
			return nil
		}
		p.Storage.SetExpire(k, id, DefaultCacheSeconds*10)
	}
	return p.GetPaymentOrder(id)
}

// 创建支付单
func (p *paymentRep) CreatePaymentOrder(
	o *payment.PaymentOrder) payment.IPaymentOrder {
	return p.PaymentRepBase.CreatePaymentOrder(o, p,
		p._memberRep, p._orderRep.Manager(), p._valRep)
}

// 保存支付单
func (p *paymentRep) SavePaymentOrder(v *payment.PaymentOrder) (int64, error) {
	stat := v.State
	if v.Id > 0 {
		stat = p.GetPaymentOrder(v.Id).GetValue().State
	}
	id, err := orm.Save(p.GetOrm(), v, v.Id)
	if err == nil {
		v.Id = id
		// 缓存订单
		p.Storage.SetExpire(p.getPaymentOrderCk(id), *v, DefaultCacheSeconds)
		// 缓存订单号与订单的关系
		p.Storage.SetExpire(p.getPaymentOrderCkByNo(v.TradeNo), v.Id, DefaultCacheSeconds*10)
		// 已经更改过状态,且为已成功,则推送到队列中
		if stat != v.State && v.State == payment.StateFinishPayment {
			p.notifyPaymentFinish(v.Id)
		}
	}
	return id, err
}

// 通知支付单完成
func (p *paymentRep) notifyPaymentFinish(paymentOrderId int64) error {
	rc := core.GetRedisConn()
	defer rc.Close()
	_, err := rc.Do("RPUSH", variable.KvPaymentOrderFinishQueue, paymentOrderId)
	//log.Println("--  推送支付单成功", paymentOrderId)
	return err
}
