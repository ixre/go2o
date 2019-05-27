/**
 * Copyright 2015 @ z3q.net.
 * name : after_sales_repo.go
 * author : jarryliu
 * date : 2016-07-17 08:36
 * description :
 * history :
 */
package repos

import (
	"github.com/ixre/gof/db"
	asImpl "go2o/core/domain/after-sales"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
)

var _ afterSales.IAfterSalesRepo = new(afterSalesRepo)

type afterSalesRepo struct {
	db.Connector
	orderRepo   order.IOrderRepo
	memberRepo  member.IMemberRepo
	paymentRepo payment.IPaymentRepo
}

func NewAfterSalesRepo(conn db.Connector, orderRepo order.IOrderRepo,
	memberRepo member.IMemberRepo, paymentRepo payment.IPaymentRepo) afterSales.IAfterSalesRepo {
	return &afterSalesRepo{
		Connector:   conn,
		orderRepo:   orderRepo,
		memberRepo:  memberRepo,
		paymentRepo: paymentRepo,
	}

}

// 创建售后单
func (a *afterSalesRepo) CreateAfterSalesOrder(v *afterSales.AfterSalesOrder) afterSales.IAfterSalesOrder {
	return asImpl.NewAfterSalesOrder(v, a, a.orderRepo, a.memberRepo, a.paymentRepo)
}

// 获取售后单
func (a *afterSalesRepo) GetAfterSalesOrder(id int32) afterSales.IAfterSalesOrder {
	v := &afterSales.AfterSalesOrder{}
	if a.GetOrm().Get(id, v) == nil {
		return a.CreateAfterSalesOrder(v)
	}
	return nil
}

// 获取订单的售后单
func (a *afterSalesRepo) GetAllOfSaleOrder(orderId int64) []afterSales.IAfterSalesOrder {
	list := []*afterSales.AfterSalesOrder{}
	orders := []afterSales.IAfterSalesOrder{}
	if a.GetOrm().Select(&list, "order_id= $1", orderId) == nil {
		for _, v := range list {
			orders = append(orders, a.CreateAfterSalesOrder(v))
		}
	}
	return orders
}
