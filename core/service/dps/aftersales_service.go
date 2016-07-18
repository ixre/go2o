/**
 * Copyright 2015 @ z3q.net.
 * name : aftersales_service.go
 * author : jarryliu
 * date : 2016-07-18 17:16
 * description :
 * history :
 */
package dps

import (
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/interface/order"
)

type afterSalesService struct {
	_orderRep order.IOrderRep
	_rep      afterSales.IAfterSalesRep
}

func NewAfterSalesService(rep afterSales.IAfterSalesRep, orderRep order.IOrderRep) *afterSalesService {
	return &afterSalesService{
		_rep:      rep,
		_orderRep: orderRep,
	}
}

// 提交售后单
func (a *afterSalesService) SubmitAfterSalesOrder(orderId int, asType int,
	snapshotId int, quantity int, reason string, img string) error {
	ro := a._rep.CreateAfterSalesOrder(&afterSales.AfterSalesOrder{
		// 订单编号
		OrderId: orderId,
		// 类型，退货、换货、维修
		Type: asType,
		// 售后原因
		Reason:        reason,
		ReturnSpImage: img,
	})
	return ro.Submit()
}
