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
	"go2o/core/domain/interface/after-sales"
)

var _ afterSales.IAfterSalesRep

type afterSalesRep struct {
	db.Connector
}

func NewAfterSalesRep(conn db.Connector) afterSales.IAfterSalesRep {
	return &afterSalesRep{
		Connector: conn,
	}
}

// 创建退款单
func (a *afterSalesRep) CreateRefundOrder(v *afterSales.RefundOrder) afterSales.IRefundOrder {
	return nil
}
