/**
 * Copyright 2015 @ z3q.net.
 * name : after_sales_rep.go
 * author : jarryliu
 * date : 2016-07-17 08:36
 * description :
 * history :
 */
package repository

import (
	"github.com/jsix/gof/db"
	asImpl "go2o/core/domain/after-sales"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/interface/order"
)

var _ afterSales.IAfterSalesRep = new(afterSalesRep)

type afterSalesRep struct {
	db.Connector
	_orderRep order.IOrderRep
}

func NewAfterSalesRep(conn db.Connector, orderRep order.IOrderRep) afterSales.IAfterSalesRep {
	return &afterSalesRep{
		Connector: conn,
		_orderRep: orderRep,
	}

}

// 创建退款单
func (a *afterSalesRep) CreateRefundOrder(v *afterSales.RefundOrder) afterSales.IRefundOrder {
	return asImpl.NewRefundOrder(v, a, a._orderRep)
}

// 获取退款单
func (a *afterSalesRep) GetRefundOrder(id int) afterSales.IRefundOrder {
	e := &afterSales.RefundOrder{}
	if a.GetOrm().Get(id, e) == nil {
		return a.CreateRefundOrder(e)
	}
	return nil
}

// 获取订单的退款单
func (a *afterSalesRep) GetRefundOrders(orderId int) []afterSales.IRefundOrder {
	list := []*afterSales.RefundOrder{}
	orders := []afterSales.IRefundOrder{}
	if a.GetOrm().Select(&list, "order_id=?", orderId) == nil {
		for _, v := range list {
			orders = append(orders, a.CreateRefundOrder(v))
		}
	}
	return orders
}
