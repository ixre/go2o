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
	"database/sql"
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
	"log"
)

var _ payment.IPaymentRepo = new(paymentRepoImpl)

type paymentRepoImpl struct {
	db.Connector
	Storage storage.Interface
	*payImpl.RepoBase
	memberRepo member.IMemberRepo
	valueRepo  valueobject.IValueRepo
	orderRepo  order.IOrderRepo
}

func NewPaymentRepo(sto storage.Interface, conn db.Connector, mmRepo member.IMemberRepo,
	orderRepo order.IOrderRepo, valRepo valueobject.IValueRepo) payment.IPaymentRepo {
	return &paymentRepoImpl{
		Storage:    sto,
		Connector:  conn,
		memberRepo: mmRepo,
		valueRepo:  valRepo,
		orderRepo:  orderRepo,
	}
}

// 根据订单号获取支付单
func (p *paymentRepoImpl) GetPaymentBySalesOrderId(orderId int64) payment.IPaymentOrder {
	e := &payment.Order{}
	if p.Connector.GetOrm().GetBy(e, "order_id=?", orderId) == nil {
		return p.CreatePaymentOrder(e)
	}
	return nil
}

// 根据订单号获取支付单
func (p *paymentRepoImpl) GetPaymentOrderByOrderNo(orderType int, orderNo string) payment.IPaymentOrder {
	e := &payment.Order{}
	if p.Connector.GetOrm().GetBy(e, "out_order_no=? AND order_type=?",
		orderNo, orderType) == nil {
		return p.CreatePaymentOrder(e)
	}
	return nil
}

func (p *paymentRepoImpl) GetMergePayOrders(mergeTradeNo string) []payment.IPaymentOrder {
	var list []*payment.Order
	p.Connector.GetOrm().Select(&list, "merge_trade_no=? AND state=? LIMIT 10",
		mergeTradeNo, payment.StateAwaitingPayment)
	var arr = make([]payment.IPaymentOrder, len(list))
	for i, v := range list {
		arr[i] = p.CreatePaymentOrder(v)
	}
	return arr
}

func (p *paymentRepoImpl) getPaymentOrderCk(id int) string {
	return fmt.Sprintf("go2o:repo:pay:order:%d", id)
}
func (p *paymentRepoImpl) getPaymentOrderCkByNo(orderNO string) string {
	return fmt.Sprintf("go2o:repo:pay:order:%s", orderNO)
}

// 根据编号获取支付单
func (p *paymentRepoImpl) GetPaymentOrderById(id int) payment.IPaymentOrder {
	if id <= 0 {
		return nil
	}
	e := &payment.Order{}
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
func (p *paymentRepoImpl) GetPaymentOrder(paymentNo string) payment.IPaymentOrder {
	k := p.getPaymentOrderCkByNo(paymentNo)
	id, err := p.Storage.GetInt(k)
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
func (p *paymentRepoImpl) CreatePaymentOrder(
	o *payment.Order) payment.IPaymentOrder {
	return p.RepoBase.CreatePaymentOrder(o, p,
		p.memberRepo, p.orderRepo.Manager(), p.valueRepo)
}

// 保存支付单
func (p *paymentRepoImpl) SavePaymentOrder(v *payment.Order) (int, error) {
	stat := v.State
	if v.ID > 0 {
		stat = p.GetPaymentOrderById(v.ID).Get().State
	}
	id, err := orm.Save(p.GetOrm(), v, v.ID)
	if err == nil {
		v.ID = id
		// 缓存订单
		p.Storage.SetExpire(p.getPaymentOrderCk(id), *v, DefaultCacheSeconds)
		// 缓存订单号与订单的关系
		p.Storage.SetExpire(p.getPaymentOrderCkByNo(v.TradeNo), v.ID, DefaultCacheSeconds*10)
		// 已经更改过状态,且为已成功,则推送到队列中
		if stat != v.State && v.State == payment.StateFinished {
			p.notifyPaymentFinish(v.ID)
		}
	}
	return id, err
}

// 通知支付单完成
func (p *paymentRepoImpl) notifyPaymentFinish(paymentOrderId int) error {
	rc := core.GetRedisConn()
	defer rc.Close()
	_, err := rc.Do("RPUSH", variable.KvPaymentOrderFinishQueue, paymentOrderId)
	//log.Println("--  推送支付单成功", paymentOrderId,err)
	return err
}

// 检查交易单号是否匹配
func (p *paymentRepoImpl) CheckTradeNoMatch(tradeNo string, id int) bool {
	i := 0
	p.Connector.ExecScalar("SELECT id FROM pay_order WHERE trade_no=? AND id<>?", &i, tradeNo, id)
	return i == 0
}

func (p *paymentRepoImpl) GetTradeChannelItems(tradeNo string) []*payment.TradeChan {
	var list []*payment.TradeChan
	err := p.GetOrm().Select(&list, "trade_no=?", tradeNo)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PayTradeChan")
	}
	return list
}

func (p *paymentRepoImpl) SavePaymentTradeChan(tradeNo string, tradeChan *payment.TradeChan) (int, error) {
	tradeChan.TradeNo = tradeNo
	id, err := orm.Save(p.GetOrm(), tradeChan, tradeChan.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PayTradeChan")
	}
	return id, err
}
