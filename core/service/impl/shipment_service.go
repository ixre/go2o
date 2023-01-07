/**
 * Copyright 2015 @ 56x.net.
 * name : 2.express_service.go
 * author : jarryliu
 * date : 2016-07-05 18:57
 * description :
 * history :
 */
package impl

import (
	"context"
	"fmt"

	"github.com/ixre/go2o/core/domain/interface/delivery"
	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/domain/interface/shipment"
	"github.com/ixre/go2o/core/module"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/log"
)

var _ proto.ShipmentServiceServer = new(shipmentServiceImpl)

type shipmentServiceImpl struct {
	repo         shipment.IShipmentRepo
	deliveryRepo delivery.IDeliveryRepo
	expressRepo  express.IExpressRepo
	serviceUtil
	proto.UnimplementedShipmentServiceServer
}

// 获取快递服务
func NewShipmentService(rep shipment.IShipmentRepo,
	deliveryRepo delivery.IDeliveryRepo, expressRepo express.IExpressRepo) *shipmentServiceImpl {
	return &shipmentServiceImpl{
		repo:         rep,
		deliveryRepo: deliveryRepo,
		expressRepo:  expressRepo,
	}
}

// 创建一个配送覆盖的区域
func (s *shipmentServiceImpl) CreateCoverageArea_(_ context.Context, r *proto.SCoverageValue) (*proto.Result, error) {
	v := s.parseCoverArea(r)
	_, err := s.deliveryRepo.SaveCoverageArea(v)
	return s.error(err), nil
}

// GetShipOrderOfOrder 获取订单的发货单信息
func (s *shipmentServiceImpl) GetShipOrderOfOrder(_ context.Context, r *proto.OrderId) (*proto.ShipmentOrderListResponse, error) {
	list := s.repo.GetShipOrders(r.Value, true)
	arr := make([]*proto.SShipmentOrder, len(list))
	if list != nil {
		for i, v := range list {
			dst := s.parseShipOrderDto(v.Value())
			items := v.Items()
			dst.Items = make([]*proto.SShipmentItem, len(items))
			for i, v := range items {
				dst.Items[i] = s.parseShipItemDto(v)
			}
			arr[i] = dst
		}
	}
	return &proto.ShipmentOrderListResponse{Value: arr}, nil
}

func (s *shipmentServiceImpl) GetLogisticFlowTrack(_ context.Context, r *proto.LogisticFlowTrackRequest) (*proto.SShipOrderTrack, error) {
	em := module.Get(module.EXPRESS).(*module.ExpressModule)
	flow, err := em.GetLogisticFlowTrack(r.ShipperCode, r.LogisticCode, r.Invert)
	if err == nil {
		return s.logisticFlowTrackDto(flow), nil
	}
	return &proto.SShipOrderTrack{
		Code:    1,
		Message: err.Error(),
	}, nil
}
func (s *shipmentServiceImpl) logisticFlowTrackDto(o *shipment.ShipOrderTrack) *proto.SShipOrderTrack {
	if o == nil {
		return &proto.SShipOrderTrack{
			Code:    1,
			Message: "无法获取物流信息",
		}
	}
	r := &proto.SShipOrderTrack{
		LogisticCode: o.LogisticCode,
		ShipperName:  o.ShipperName,
		ShipperCode:  o.ShipperCode,
		ShipState:    o.ShipState,
		UpdateTime:   o.UpdateTime,
		Flows:        make([]*proto.SShipFlow, 0),
	}
	for _, v := range o.Flows {
		r.Flows = append(r.Flows, &proto.SShipFlow{
			Subject:    v.Subject,
			CreateTime: v.CreateTime,
			Remark:     v.Remark,
		})
	}
	return r
}

// 获取发货单的物流追踪信息,
// - shipOrderId:发货单编号
func (s *shipmentServiceImpl) ShipOrderLogisticTrack(ctx context.Context, rq *proto.OrderLogisticTrackRequest) (*proto.SShipOrderTrack, error) {
	so := s.repo.GetShipmentOrder(rq.ShipOrderId)
	if so != nil {
		sp := s.expressRepo.GetExpressProvider(so.Value().SpId)
		if sp == nil {
			log.Println("[ Go2o][ Service][ Warning]: no such express provider id ", so.Value().SpId)
		} else {
			//spOrder = "462681586678"
			//sp.ApiCode = "ZTO"
			spOrder := so.Value().SpOrder
			r, err := s.GetLogisticFlowTrack(ctx, &proto.LogisticFlowTrackRequest{
				ShipperCode:  sp.ApiCode,
				LogisticCode: spOrder,
				Invert:       rq.Invert,
			})
			r.ShipperName = sp.Name
			return r, err
		}
	}
	return nil, fmt.Errorf("no such ship order")
}

func (s *shipmentServiceImpl) parseCoverArea(r *proto.SCoverageValue) *delivery.CoverageValue {
	return &delivery.CoverageValue{
		Id:      int32(r.Id),
		Name:    r.Name,
		Lng:     r.Lng,
		Lat:     r.Lat,
		Radius:  int(r.Radius),
		Address: r.Address,
		AreaId:  int32(r.AreaId),
	}
}

func (s *shipmentServiceImpl) parseShipOrderDto(v shipment.ShipmentOrder) *proto.SShipmentOrder {
	return &proto.SShipmentOrder{
		Id:          v.ID,
		OrderId:     v.OrderId,
		SubOrderId:  v.SubOrderId,
		ExpressSpId: int64(v.SpId),
		ShipOrderNo: v.SpOrder,
		ShipmentLog: v.ShipmentLog,
		Amount:      v.Amount,
		FinalAmount: v.FinalAmount,
		ShipTime:    v.ShipTime,
		State:       int32(v.State),
		UpdateTime:  v.UpdateTime,
		Items:       nil,
	}
}

func (s *shipmentServiceImpl) parseShipItemDto(v *shipment.ShipmentItem) *proto.SShipmentItem {
	return &proto.SShipmentItem{
		Id:          v.ID,
		SnapshotId:  v.SnapshotId,
		Quantity:    v.Quantity,
		Amount:      v.Amount,
		FinalAmount: v.FinalAmount,
	}
}
