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
	"go2o/core/dto"
	"go2o/core/infrastructure/format"
	"go2o/core/query"
)

type afterSalesService struct {
	_orderRep order.IOrderRep
	_rep      afterSales.IAfterSalesRep
	_query    *query.AfterSalesQuery
}

func NewAfterSalesService(rep afterSales.IAfterSalesRep,
	q *query.AfterSalesQuery, orderRep order.IOrderRep) *afterSalesService {
	return &afterSalesService{
		_rep:      rep,
		_orderRep: orderRep,
		_query:    q,
	}
}

// 提交售后单
func (a *afterSalesService) SubmitAfterSalesOrder(orderId int, asType int,
	snapshotId int, quantity int, reason string, img string) (int, error) {
	ro := a._rep.CreateAfterSalesOrder(&afterSales.AfterSalesOrder{
		// 订单编号
		OrderId: orderId,
		// 类型，退货、换货、维修
		Type: asType,
		// 售后原因
		Reason:        reason,
		ReturnSpImage: img,
	})
	err := ro.SetItem(snapshotId,quantity)
	if err == nil {
		return ro.Submit()
	}
	return 0,err
}

// 获取订单的所有售后单
func (a *afterSalesService) GetAllAfterSalesOrderOfSaleOrder(orderId int) []afterSales.AfterSalesOrder {
	list := a._rep.GetAllOfSaleOrder(orderId)
	arr := make([]afterSales.AfterSalesOrder, len(list))
	for i, v := range list {
		arr[i] = v.Value()
		arr[i].StateText = afterSales.Stat(arr[i].State).String()
	}
	return arr
}

// 获取会员的分页售后单
func (a *afterSalesService) QueryPagerAfterSalesOrderOfMember(memberId, begin,
	size int, where string) (int, []*dto.PagedMemberAfterSalesOrder) {
	return a._query.QueryPagerAfterSalesOrderOfMember(memberId, begin, size, where)
}

// 获取售后单
func (a *afterSalesService) GetAfterSaleOrder(id int) *afterSales.AfterSalesOrder {
	as := a._rep.GetAfterSalesOrder(id)
	if as != nil {
		v := as.Value()
		v.StateText = afterSales.Stat(v.State).String()
		v.ReturnSpImage = format.GetResUrl(v.ReturnSpImage)
		return &v
	}
	return nil
}
