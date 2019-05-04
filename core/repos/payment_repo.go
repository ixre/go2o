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
	"bytes"
	"database/sql"
	"fmt"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"go2o/core"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/valueobject"
	payImpl "go2o/core/domain/payment"
	"go2o/core/variable"
	"log"
	"strings"
	"time"
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
	if p.Connector.GetOrm().GetBy(e, "order_id= $1", orderId) == nil {
		return p.CreatePaymentOrder(e)
	}
	return nil
}

// 根据订单号获取支付单
func (p *paymentRepoImpl) GetPaymentOrderByOrderNo(orderType int, orderNo string) payment.IPaymentOrder {
	e := &payment.Order{}
	if p.Connector.GetOrm().GetBy(e, "out_order_no= $1 AND order_type= $2",
		orderNo, orderType) == nil {
		return p.CreatePaymentOrder(e)
	}
	return nil
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
		p.ExecScalar("SELECT id FROM pay_order where trade_no= $1", &id, paymentNo)
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
	p.Connector.ExecScalar("SELECT id FROM pay_order WHERE trade_no= $1 AND id<> $2", &i, tradeNo, id)
	return i == 0
}

func (p *paymentRepoImpl) GetTradeChannelItems(tradeNo string) []*payment.TradeMethodData {
	list := make([]*payment.TradeMethodData, 0)
	err := p.GetOrm().Select(&list, "trade_no= $1 LIMIT $2", tradeNo, 10)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PayTradeChan")
	}
	return list
}

func (p *paymentRepoImpl) SavePaymentTradeChan(tradeNo string, tradeChan *payment.TradeMethodData) (int, error) {
	tradeChan.TradeNo = tradeNo
	id, err := orm.Save(p.GetOrm(), tradeChan, tradeChan.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PayTradeChan")
	}
	return id, err
}

func (p *paymentRepoImpl) GetMergePayOrders(mergeTradeNo string) []payment.IPaymentOrder {
	var tradeNo = ""
	var tradeNoArr []string
	// 查询支付单号
	p.Connector.Query("SELECT order_trade_no FROM pay_merge_order WHERE merge_trade_no= $1 LIMIT $2",
		func(rows *sql.Rows) {
			for rows.Next() {
				rows.Scan(&tradeNo)
				tradeNoArr = append(tradeNoArr, "'"+tradeNo+"'")
			}
		}, mergeTradeNo, 10)

	var arr = make([]payment.IPaymentOrder, 0)
	// 查询支付单
	if l := len(tradeNoArr); l > 0 {
		list := make([]*payment.Order, 0)
		p.Connector.GetOrm().Select(&list, "trade_no IN ("+strings.Join(tradeNoArr, ",")+
			") LIMIT $1", len(tradeNoArr))
		for _, v := range list {
			arr = append(arr, p.CreatePaymentOrder(v))
		}
	}
	return arr
}

func (p *paymentRepoImpl) ResetMergePaymentOrders(tradeNos []string) error {
	buf := bytes.NewBuffer([]byte("("))
	for i, v := range tradeNos {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("'")
		buf.WriteString(v)
		buf.WriteString("'")
	}
	buf.WriteString(")")
	_, err := p.Connector.GetOrm().Delete(&payment.MergeOrder{},
		"order_trade_no in "+buf.String())
	return err
}

func (p *paymentRepoImpl) SaveMergePaymentOrders(mergeTradeNo string, tradeNos []string) error {
	unix := time.Now().Unix()
	orm := p.Connector.GetOrm()
	for _, v := range tradeNos {
		order := &payment.MergeOrder{
			MergeTradeNo: mergeTradeNo,
			OrderTradeNo: v,
			SubmitTime:   unix,
		}
		if _, _, err := orm.Save(nil, order); err != nil {
			return err
		}
	}
	return nil
}
