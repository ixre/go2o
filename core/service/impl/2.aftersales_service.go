/**
 * Copyright 2015 @ 56x.net.
 * name : aftersales_service.go
 * author : jarryliu
 * date : 2016-07-18 17:16
 * description :
 * history :
 */
package impl

import (
	"context"
	"errors"

	afterSales "github.com/ixre/go2o/core/domain/interface/aftersales"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/db"
)

var _ proto.AfterSalesServiceServer = new(afterSalesService)

type afterSalesService struct {
	_orderRepo order.IOrderRepo
	_rep       afterSales.IAfterSalesRepo
	_query     *query.AfterSalesQuery
	db.Connector
	serviceUtil
	proto.UnimplementedAfterSalesServiceServer
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
func (a *afterSalesService) SubmitAfterSalesOrder(_ context.Context, r *proto.SubmitAfterSalesOrderRequest) (*proto.SubmitAfterSalesOrderResponse, error) {
	af := &afterSales.AfterSalesOrder{
		// 订单编号
		OrderId: r.OrderId,
		// 类型，退货、换货、维修
		Type: int(r.AfterSalesType),
		// 售后原因
		Reason: r.Reason,
	}
	if len(r.Images) > 0 {
		// 上传截图
		af.ImageUrl = r.Images[0]
	}
	if len(af.Reason) < 6 {
		return &proto.SubmitAfterSalesOrderResponse{
			ErrCode: 1,
			ErrMsg:  "申请原因不能为空",
		}, nil
	}
	ro := a._rep.CreateAfterSalesOrder(af)
	err := ro.SetItem(r.ItemSnapshotId, int32(r.Quantity))
	if err == nil {
		_, err = ro.Submit()
	}
	ret := &proto.SubmitAfterSalesOrderResponse{}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	} else {
		ret.AfterSalesOrderNo = ro.Value().OrderNo
	}
	return ret, nil
}

// 获取订单的所有售后单
func (a *afterSalesService) GetAllAfterSalesOrderOfSaleOrder(_ context.Context, r *proto.OriginOrderIdRequest) (*proto.AfterSalesOrderListResponse, error) {
	list := a._rep.GetAllOfSaleOrder(r.OrderId)
	arr := make([]*proto.SAfterSalesOrder, len(list))
	for i, v := range list {
		arr[i] = a.parseAfterSalesDto(v.Value())
		arr[i].StatusText = afterSales.Stat(arr[i].Status).String()
	}
	return &proto.AfterSalesOrderListResponse{
		Value: arr,
	}, nil
}

// 获取会员的分页售后单
func (a *afterSalesService) QueryPagerAfterSalesOrderOfMember(_ context.Context, r *proto.PagingBuyerOrdersRequest) (*proto.PagingBuyerAfterSalesOrderListResponse, error) {
	total, list := a._query.QueryPagerAfterSalesOrderOfMember(
		r.BuyerId, int(r.Params.Begin), int(r.Params.End), r.Params.Where)
	arr := make([]*proto.SPagingBuyerAfterSalesOrder, len(list))
	for i, v := range list {
		arr[i] = &proto.SPagingBuyerAfterSalesOrder{
			Id:             int64(v.Id),
			OrderNo:        v.OrderNo,
			VendorId:       int64(v.VendorId),
			SellerName:     v.VendorName,
			AfterSalesType: int32(v.Type),
			SnapshotId:     int64(v.SnapshotId),
			Quantity:       int32(v.Quantity),
			SkuId:          int64(v.SkuId),
			ItemTitle:      v.GoodsTitle,
			ItemImage:      v.GoodsImage,
			Status:         int32(v.Status),
			CreateTime:     v.CreateTime,
		}
		arr[i].StatusText = afterSales.Stat(arr[i].Status).String()
	}
	return &proto.PagingBuyerAfterSalesOrderListResponse{
		Total: int64(total),
		Data:  arr,
	}, nil
}

// 获取商户的分页售后单
func (a *afterSalesService) QueryPagerAfterSalesOrderOfVendor(_ context.Context, r *proto.PagingSellerOrdersRequest) (*proto.PagingSellerAfterSalesOrderListResponse, error) {
	total, list := a._query.QueryPagerAfterSalesOrderOfVendor(
		r.SellerId, int(r.Params.Begin), int(r.Params.End), r.Params.Where)
	arr := make([]*proto.SPagingSellerAfterSalesOrder, len(list))
	for i, v := range list {
		arr[i] = &proto.SPagingSellerAfterSalesOrder{
			Id:             int64(v.Id),
			OrderNo:        v.OrderNo,
			BuyerId:        int64(v.BuyerId),
			BuyerName:      v.BuyerName,
			AfterSalesType: int32(v.Type),
			SnapshotId:     int64(v.SnapshotId),
			Quantity:       int32(v.Quantity),
			SkuId:          int64(v.SkuId),
			ItemTitle:      v.GoodsTitle,
			ItemImage:      v.GoodsImage,
			Status:         int32(v.Status),
			CreateTime:     v.CreateTime,
			UpdateTime:     v.UpdateTime,
		}
		arr[i].StatusText = afterSales.Stat(arr[i].Status).String()
	}
	return &proto.PagingSellerAfterSalesOrderListResponse{
		Total: int64(total),
		Data:  arr,
	}, nil
}

// 获取售后单
func (a *afterSalesService) GetAfterSaleOrder(_ context.Context, req *proto.AfterSalesOrderNo) (*proto.SAfterSalesOrder, error) {
	as := a._rep.GetAfterSalesOrder(req.OrderNo)
	if as != nil {
		v := as.Value()
		v.StatusText = afterSales.Stat(v.Status).String()
		v.ShipmentImage = format.GetResUrl(v.ShipmentImage)
		return a.parseAfterSalesDto(v), nil
	}
	return nil, order.ErrNoSuchOrder
}

// 同意售后
func (a *afterSalesService) AgreeAfterSales(_ context.Context, req *proto.AfterSalesProcessRequest) (*proto.Result, error) {
	as := a._rep.GetAfterSalesOrder(req.OrderNo)
	err := as.Agree()
	return a.error(err), nil
}

// 拒绝售后
func (a *afterSalesService) DeclineAfterSales(_ context.Context, req *proto.AfterSalesProcessRequest) (*proto.Result, error) {
	as := a._rep.GetAfterSalesOrder(req.OrderNo)
	err := as.Decline(req.Remark)
	return a.error(err), nil
}

// 申请调解
func (a *afterSalesService) RequestIntercede(_ context.Context, req *proto.AfterSalesProcessRequest) (*proto.Result, error) {
	as := a._rep.GetAfterSalesOrder(req.OrderNo)
	err := as.RequestIntercede()
	return a.error(err), nil
}

// 系统确认
func (a *afterSalesService) ConfirmAfterSales(_ context.Context, req *proto.AfterSalesOrderNo) (*proto.Result, error) {
	as := a._rep.GetAfterSalesOrder(req.OrderNo)
	err := as.Confirm()
	return a.error(err), nil
}

// 系统退回
func (a *afterSalesService) RejectAfterSales(_ context.Context, req *proto.AfterSalesProcessRequest) (*proto.Result, error) {
	var err error
	as := a._rep.GetAfterSalesOrder(req.OrderNo)
	if as == nil {
		err = afterSales.ErrNoSuchOrder
	} else {
		err = as.Reject(req.Remark)
	}
	return a.error(err), nil
}

// 处理退款/退货完成,一般是系统自动调用
func (a *afterSalesService) ProcessAfterSalesOrder(_ context.Context, req *proto.AfterSalesProcessRequest) (*proto.Result, error) {
	var err error
	as := a._rep.GetAfterSalesOrder(req.OrderNo)
	if as == nil {
		err = afterSales.ErrNoSuchOrder
	} else {
		switch as.Value().Type {
		case afterSales.TypeRefund:
			err = as.Process()
		case afterSales.TypeReturn:
			err = as.Process()
		default:
			err = afterSales.ErrAutoProcess
		}
	}
	return a.error(err), nil
}

// 换货发货
func (a *afterSalesService) ReturnShipment(_ context.Context, req *proto.ReturnShipmentRequest) (*proto.Result, error) {
	ex := a._rep.GetAfterSalesOrder(req.OrderNo).(afterSales.IReturnAfterSalesOrder)
	return a.error(errors.New("该售后订单不支持当前操作")), nil
	err := ex.ReturnShipment(req.ShipmentExpress, req.ShipmentOrder, "")
	return a.error(err), nil
}

// 换货收货
func (a *afterSalesService) ReceiveItem(_ context.Context, req *proto.AfterSalesOrderNo) (*proto.Result, error) {
	ex := a._rep.GetAfterSalesOrder(req.OrderNo).(afterSales.IReturnAfterSalesOrder)
	if ex == nil {
		return a.error(errors.New("该售后订单不支持当前操作")), nil
	}
	err := ex.ReturnReceive()
	return a.error(err), nil
}

func (a *afterSalesService) parseAfterSalesDto(v afterSales.AfterSalesOrder) *proto.SAfterSalesOrder {
	return &proto.SAfterSalesOrder{
		Id:              int64(v.Id),
		OrderNo:         v.OrderNo,
		RelateOrderId:   v.OrderId,
		VendorId:        v.VendorId,
		BuyerId:         v.BuyerId,
		AfterSalesType:  int32(v.Type),
		SnapshotId:      v.SnapshotId,
		Quantity:        v.Quantity,
		Reason:          v.Reason,
		ImageUrl:        v.ImageUrl,
		PersonName:      v.PersonName,
		PersonPhone:     v.PersonPhone,
		ShipmentExpress: v.ShipmentExpress,
		ShipmentOrderNo: v.ShipmentOrderNo,
		ShipmentImage:   v.ShipmentImage,
		Remark:          v.Remark,
		VendorRemark:    v.VendorRemark,
		Status:          int32(v.Status),
		CreateTime:      v.CreateTime,
		UpdateTime:      v.UpdateTime,
		StatusText:      v.StatusText,
	}
}
