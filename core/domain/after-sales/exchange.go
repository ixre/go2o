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
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/tmp"
	"time"
)

var _ afterSales.IExchangeOrder = new(exchangeOrderImpl)
var _ afterSales.IAfterSalesOrder = new(exchangeOrderImpl)

// 换货单
//todo: 是否也需要限制退货数量
type exchangeOrderImpl struct {
	*afterSalesOrderImpl
	_excValue *afterSales.ExchangeOrder
}

func newExchangeOrderImpl(v *afterSalesOrderImpl) *exchangeOrderImpl {
	if v._value.Type != afterSales.TypeExchange {
		panic(errors.New("售后单类型不是换货单"))
	}
	return &exchangeOrderImpl{
		afterSalesOrderImpl: v,
	}
}

func (e *exchangeOrderImpl) getValue() *afterSales.ExchangeOrder {
	if e._excValue == nil {
		if e.GetDomainId() <= 0 {
			panic(errors.New("换货单还未提交"))
		}
		v := &afterSales.ExchangeOrder{}
		if tmp.Db().GetOrm().Get(e.GetDomainId(), v) == nil {
			e._excValue = v
		}
		panic(errors.New("换货单不存在"))
	}
	return e._excValue
}

// 获取售后单数据
func (e *exchangeOrderImpl) Value() afterSales.AfterSalesOrder {
	v := e.afterSalesOrderImpl.Value()
	v2 := e.getValue()
	v.Data = *v2
	return v
}

// 提交售后申请
func (e *exchangeOrderImpl) Submit() error {
	err := e.afterSalesOrderImpl.Submit()
	// 提交换货单
	if err == nil {
		e._excValue = &afterSales.ExchangeOrder{
			Id:          e.afterSalesOrderImpl.GetDomainId(),
			IsShipped:   0,
			ShipSpName:  "",
			ShipSpOrder: "",
			IsReceived:  1,
			ShipTime:    0,
			ReceiveTime: 0,
		}
		_, err = orm.Save(tmp.Db().GetOrm(), e._excValue, 0)
	}
	return err
}

// 保存换货单
func (e *exchangeOrderImpl) saveExchangeOrder(v *afterSales.ExchangeOrder) error {
	_, err := orm.Save(tmp.Db().GetOrm(), v, v.Id)
	return err
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
	return e.saveExchangeOrder(v)
}
