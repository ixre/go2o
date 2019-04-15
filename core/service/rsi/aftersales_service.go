/**
 * Copyright 2015 @ z3q.net.
 * name : aftersales_service.go
 * author : jarryliu
 * date : 2016-07-18 17:16
 * description :
 * history :
 */
package rsi

import (
	"github.com/ixre/gof/db"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/interface/order"
	"go2o/core/dto"
	"go2o/core/infrastructure/format"
	"go2o/core/query"
)

type afterSalesService struct {
	_orderRepo order.IOrderRepo
	_rep       afterSales.IAfterSalesRepo
	_query     *query.AfterSalesQuery
	db.Connector
}

func NewAfterSalesService(rep afterSales.IAfterSalesRepo,
	q *query.AfterSalesQuery, orderRepo order.IOrderRepo) *afterSalesService {
	return &afterSalesService{
		_rep:       rep,
		_orderRepo: orderRepo,
		_query:     q,
	}
}

// 提交售后单
func (a *afterSalesService) SubmitAfterSalesOrder(orderId int64, asType int,
	snapshotId int64, quantity int32, reason string, imgUrl string) (int32, error) {
	ro := a._rep.CreateAfterSalesOrder(&afterSales.AfterSalesOrder{
		// 订单编号
		OrderId: orderId,
		// 类型，退货、换货、维修
		Type: asType,
		// 售后原因
		Reason: reason,
		// 上传截图
		ImageUrl: imgUrl,
	})
	err := ro.SetItem(snapshotId, quantity)
	if err == nil {
		return ro.Submit()
	}
	return 0, err
}

// 获取订单的所有售后单
func (a *afterSalesService) GetAllAfterSalesOrderOfSaleOrder(orderId int64) []afterSales.AfterSalesOrder {
	list := a._rep.GetAllOfSaleOrder(orderId)
	arr := make([]afterSales.AfterSalesOrder, len(list))
	for i, v := range list {
		arr[i] = v.Value()
		arr[i].StateText = afterSales.Stat(arr[i].State).String()
	}
	return arr
}

// 获取会员的分页售后单
func (a *afterSalesService) QueryPagerAfterSalesOrderOfMember(memberId int64, begin,
	size int, where string) (int, []*dto.PagedMemberAfterSalesOrder) {
	return a._query.QueryPagerAfterSalesOrderOfMember(memberId, begin, size, where)
}

// 获取商户的分页售后单
func (a *afterSalesService) QueryPagerAfterSalesOrderOfVendor(vendorId int32, begin,
	size int, where string) (int, []*dto.PagedVendorAfterSalesOrder) {
	return a._query.QueryPagerAfterSalesOrderOfVendor(vendorId, begin, size, where)
}

// 获取售后单
func (a *afterSalesService) GetAfterSaleOrder(id int32) *afterSales.AfterSalesOrder {
	as := a._rep.GetAfterSalesOrder(id)
	if as != nil {
		v := as.Value()
		v.StateText = afterSales.Stat(v.State).String()
		v.ReturnSpImage = format.GetResUrl(v.ReturnSpImage)
		return &v
	}
	return nil
}

// 同意售后
func (a *afterSalesService) AgreeAfterSales(id int32, remark string) error {
	as := a._rep.GetAfterSalesOrder(id)
	return as.Agree()
}

// 拒绝售后
func (a *afterSalesService) DeclineAfterSales(id int32, reason string) error {
	as := a._rep.GetAfterSalesOrder(id)
	return as.Decline(reason)
}

// 申请调解
func (a *afterSalesService) RequestIntercede(id int32) error {
	as := a._rep.GetAfterSalesOrder(id)
	return as.RequestIntercede()
}

// 系统确认
func (a *afterSalesService) ConfirmAfterSales(id int32) error {
	as := a._rep.GetAfterSalesOrder(id)
	return as.Confirm()
}

// 系统退回
func (a *afterSalesService) RejectAfterSales(id int32, remark string) error {
	as := a._rep.GetAfterSalesOrder(id)
	if as == nil {
		return afterSales.ErrNoSuchOrder
	}

	return as.Reject(remark)
}

// 处理退款/退货完成,一般是系统自动调用
func (a *afterSalesService) ProcessAfterSalesOrder(id int32) error {
	as := a._rep.GetAfterSalesOrder(id)
	if as == nil {
		return afterSales.ErrNoSuchOrder
	}
	v := as.Value()
	switch v.Type {
	case afterSales.TypeRefund:
		return as.Process()
	case afterSales.TypeReturn:
		return as.Process()
	}
	return afterSales.ErrAutoProcess
}

// 售后收货
func (a *afterSalesService) ReceiveReturnShipment(id int32) error {
	as := a._rep.GetAfterSalesOrder(id)
	err := as.ReturnReceive()
	if err == nil {
		if as.Value().State != afterSales.TypeExchange {
			err = as.Process()
		}
	}
	return err
}

// 换货发货
func (a *afterSalesService) ExchangeShipment(id int32, spName string, spOrder string) error {
	ex := a._rep.GetAfterSalesOrder(id).(afterSales.IExchangeOrder)
	return ex.ExchangeShip(spName, spOrder)
}

// 换货收货
func (a *afterSalesService) ReceiveExchange(id int32) error {
	ex := a._rep.GetAfterSalesOrder(id).(afterSales.IExchangeOrder)
	return ex.ExchangeReceive()
}
