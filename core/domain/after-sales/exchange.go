/**
 * Copyright 2015 @ z3q.net.
 * name : exchange
 * author : jarryliu
 * date : 2016-07-18 09:55
 * description :
 * history :
 */
package afterSales

import (
	"errors"
	"github.com/ixre/gof/db/orm"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/tmp"
	"time"
)

var _ afterSales.IExchangeOrder = new(exchangeOrderImpl)
var _ afterSales.IAfterSalesOrder = new(exchangeOrderImpl)

// 换货单
//todo: 是否也需要限制退货数量
type exchangeOrderImpl struct {
	*afterSalesOrderImpl
	refValue *afterSales.ExchangeOrder
}

func newExchangeOrderImpl(v *afterSalesOrderImpl) *exchangeOrderImpl {
	if v.value.Type != afterSales.TypeExchange {
		panic(errors.New("售后单类型不是换货单"))
	}
	return &exchangeOrderImpl{
		afterSalesOrderImpl: v,
	}
}

func (e *exchangeOrderImpl) getValue() *afterSales.ExchangeOrder {
	if e.refValue == nil {
		if e.GetDomainId() <= 0 {
			panic(errors.New("换货单还未提交"))
		}
		v := &afterSales.ExchangeOrder{}
		if tmp.Db().GetOrm().Get(e.GetDomainId(), v) != nil {
			panic(errors.New("换货单不存在"))
		}
		e.refValue = v
	}
	return e.refValue
}

// 获取售后单数据
func (e *exchangeOrderImpl) Value() afterSales.AfterSalesOrder {
	v := e.afterSalesOrderImpl.Value()
	v2 := e.getValue()
	v.Data = *v2
	// 自动收货
	if v2.IsShipped == 1 && v2.ReceiveTime < time.Now().Unix() {
		v2.IsReceived = 1
		if err := e.Process(); err == nil {
			e.saveExchangeOrder(v2)
		}
	}
	return v
}

// 提交售后申请
func (e *exchangeOrderImpl) Submit() (int32, error) {
	o := e.GetOrder()
	if o.GetValue().State != order.StatCompleted {
		return 0, afterSales.ErrExchangeNotReceiveItems
	}
	id, err := e.afterSalesOrderImpl.Submit()
	// 提交换货单
	if err == nil {
		e.refValue = &afterSales.ExchangeOrder{
			Id:          e.afterSalesOrderImpl.GetDomainId(),
			IsShipped:   0,
			ShipSpName:  "",
			ShipSpOrder: "",
			IsReceived:  1,
			ShipTime:    0,
			ReceiveTime: 0,
		}
		_, err = orm.Save(tmp.Db().GetOrm(), e.refValue, 0)
	}
	return id, err
}

// 保存换货单
func (e *exchangeOrderImpl) saveExchangeOrder(v *afterSales.ExchangeOrder) error {
	_, err := orm.Save(tmp.Db().GetOrm(), v, int(v.Id))
	return err
}

// 处理完成
func (e *exchangeOrderImpl) Process() error {
	v := e.getValue()
	if v.IsShipped == 0 {
		return afterSales.ErrExchangeOrderNoShipping
	}
	if v.IsReceived == 0 {
		return afterSales.ErrNotReceive
	}
	return e.afterSalesOrderImpl.Process()
}

// 将换货的商品重新发货
func (e *exchangeOrderImpl) ExchangeShip(spName string, spOrder string) error {
	if e.afterSalesOrderImpl.GetDomainId() <= 0 {
		panic(errors.New("换货单尚未提交"))
	}
	v := e.getValue()
	v.ShipSpName = spName
	v.ShipSpOrder = spOrder
	v.ShipTime = time.Now().Unix()
	v.ReceiveTime = time.Now().Add(time.Hour * 24 * 7).Unix() //7天后自动收货
	v.IsShipped = 1
	return e.saveExchangeOrder(v)
}

// 消费者延长收货时间
func (e *exchangeOrderImpl) LongReceive() error {
	v := e.getValue()
	v.ReceiveTime += 1440
	return e.saveExchangeOrder(v)
}

// 接收换货
func (e *exchangeOrderImpl) ExchangeReceive() error {
	v := e.getValue()
	v.IsReceived = 1
	v.ReceiveTime = time.Now().Unix()
	err := e.saveExchangeOrder(v)
	if err == nil {
		err = e.Process()
	}
	return err
}
