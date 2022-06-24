/**
 * Copyright 2015 @ 56x.net.
 * name : after_sales_repo.go
 * author : jarryliu
 * date : 2016-07-17 08:36
 * description :
 * history :
 */
package repos

import (
	asImpl "github.com/ixre/go2o/core/domain/after-sales"
	afterSales "github.com/ixre/go2o/core/domain/interface/after-sales"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

var _ afterSales.IAfterSalesRepo = new(afterSalesRepo)

type afterSalesRepo struct {
	db.Connector
	orderRepo   order.IOrderRepo
	memberRepo  member.IMemberRepo
	paymentRepo payment.IPaymentRepo
	o           orm.Orm
}

// GetFreeOrderNo implements afterSales.IAfterSalesRepo
func (a *afterSalesRepo) GetFreeOrderNo(orderId int64) string {
	//todo:实际应用需要预先生成订单号
	var orderNo string
	for {
		orderNo = domain.NewTradeNo(1, int(orderId))
		var rec int
		if a.ExecScalar(`SELECT id FROM sale_after_order WHERE order_no= $1 LIMIT 1`,
			&rec, orderNo); rec == 0 {
			break
		}
	}
	return orderNo
}

func NewAfterSalesRepo(o orm.Orm, orderRepo order.IOrderRepo,
	memberRepo member.IMemberRepo, paymentRepo payment.IPaymentRepo) afterSales.IAfterSalesRepo {
	return &afterSalesRepo{
		Connector:   o.Connector(),
		o:           o,
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
	if a.o.Get(id, v) == nil {
		return a.CreateAfterSalesOrder(v)
	}
	return nil
}

// 获取订单的售后单
func (a *afterSalesRepo) GetAllOfSaleOrder(orderId int64) []afterSales.IAfterSalesOrder {
	list := []*afterSales.AfterSalesOrder{}
	orders := []afterSales.IAfterSalesOrder{}
	if a.o.Select(&list, "order_id= $1", orderId) == nil {
		for _, v := range list {
			orders = append(orders, a.CreateAfterSalesOrder(v))
		}
	}
	return orders
}
