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
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
)

var _ afterSales.IAfterSalesRep = new(afterSalesRep)

type afterSalesRep struct {
	db.Connector
	_orderRep  order.IOrderRep
	_memberRep member.IMemberRep
}

func NewAfterSalesRep(conn db.Connector, orderRep order.IOrderRep,
	memberRep member.IMemberRep) afterSales.IAfterSalesRep {
	return &afterSalesRep{
		Connector:  conn,
		_orderRep:  orderRep,
		_memberRep: memberRep,
	}

}

// 创建售后单
func (a *afterSalesRep) CreateAfterSalesOrder(v *afterSales.AfterSalesOrder) afterSales.IAfterSalesOrder {
	return asImpl.NewAfterSalesOrder(v, a, a._orderRep, a._memberRep)
}

// 获取售后单
func (a *afterSalesRep) GetAfterSalesOrder(id int) afterSales.IAfterSalesOrder {
	v := &afterSales.AfterSalesOrder{}
	if a.GetOrm().Get(id, v) == nil {
		return a.CreateAfterSalesOrder(v)
	}
	return nil
}

// 获取订单的售后单
func (a *afterSalesRep) GetAllOfSaleOrder(orderId int) []afterSales.IAfterSalesOrder {
	list := []*afterSales.AfterSalesOrder{}
	orders := []afterSales.IAfterSalesOrder{}
	if a.GetOrm().Select(&list, "order_id=?", orderId) == nil {
		for _, v := range list {
			orders = append(orders, a.CreateAfterSalesOrder(v))
		}
	}
	return orders
}
