/**
 * Copyright 2015 @ z3q.net.
 * name : payment_repo.go
 * author : jarryliu
 * date : 2016-07-03 12:52
 * description :
 * history :
 */
package repos

import (
	"fmt"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"github.com/jsix/gof/storage"
	"github.com/jsix/gof/util"
	"go2o/core"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/valueobject"
	payImpl "go2o/core/domain/payment"
	"go2o/core/variable"
)

var _ payment.IPaymentRepo = new(paymentRepo)

type paymentRepo struct {
	db.Connector
	Storage storage.Interface
	*payImpl.PaymentRepBase
	_memberRepo member.IMemberRepo
	_valRepo    valueobject.IValueRepo
	_orderRepo  order.IOrderRepo
}

func NewPaymentRepo(sto storage.Interface, conn db.Connector, mmRepo member.IMemberRepo,
	orderRepo order.IOrderRepo, valRepo valueobject.IValueRepo) payment.IPaymentRepo {
	return &paymentRepo{
		Storage:     sto,
		Connector:   conn,
		_memberRepo: mmRepo,
		_valRepo:    valRepo,
		_orderRepo:  orderRepo,
	}
}

// 根据订单号获取支付单
func (p *paymentRepo) GetPaymentBySalesOrderId(orderId int64) payment.IPaymentOrder {
	e := &payment.PaymentOrder{}
	if p.Connector.GetOrm().GetBy(e, "order_id=?", orderId) == nil {
		return p.CreatePaymentOrder(e)
	}
	return nil
}

func (p *paymentRepo) getPaymentOrderCk(id int32) string {
	return fmt.Sprintf("go2o:repo:pay:order:%d", id)
}
func (p *paymentRepo) getPaymentOrderCkByNo(orderNO string) string {
	return fmt.Sprintf("go2o:repo:pay:order:%s", orderNO)
}

// 根据编号获取支付单
func (p *paymentRepo) GetPaymentOrderById(id int32) payment.IPaymentOrder {
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
func (p *paymentRepo) GetPaymentOrder(paymentNo string) payment.IPaymentOrder {
	k := p.getPaymentOrderCkByNo(paymentNo)
	id, err := util.I32Err(p.Storage.GetInt(k))
	if err != nil {
		p.ExecScalar("SELECT id FROM pay_order where trade_no=?", &id, paymentNo)
		if id == 0 {
			return nil
		}
		p.Storage.SetExpire(k, id, DefaultCacheSeconds*10)
	}
	return p.GetPaymentOrderById(id)
}

// 创建支付单
func (p *paymentRepo) CreatePaymentOrder(
	o *payment.PaymentOrder) payment.IPaymentOrder {
	return p.PaymentRepBase.CreatePaymentOrder(o, p,
		p._memberRepo, p._orderRepo.Manager(), p._valRepo)
}

// 保存支付单
func (p *paymentRepo) SavePaymentOrder(v *payment.PaymentOrder) (int32, error) {
	stat := v.State
	if v.Id > 0 {
		stat = p.GetPaymentOrderById(v.Id).GetValue().State
	}
	id, err := orm.I32(orm.Save(p.GetOrm(), v, int(v.Id)))
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
func (p *paymentRepo) notifyPaymentFinish(paymentOrderId int32) error {
	rc := core.GetRedisConn()
	defer rc.Close()
	_, err := rc.Do("RPUSH", variable.KvPaymentOrderFinishQueue, paymentOrderId)
	//log.Println("--  推送支付单成功", paymentOrderId,err)
	return err
}

// 检查交易单号是否匹配
func (p *paymentRepo) CheckTradeNoMatch(tradeNo string, id int32) bool {
	i := 0
	p.Connector.ExecScalar("SELECT id FROM pay_order WHERE trade_no=? AND id<>?", &i, tradeNo, id)
	return i == 0
}
