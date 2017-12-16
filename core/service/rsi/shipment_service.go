/**
 * Copyright 2015 @ z3q.net.
 * name : express_service.go
 * author : jarryliu
 * date : 2016-07-05 18:57
 * description :
 * history :
 */
package rsi

import (
	"github.com/jsix/gof/log"
	"go2o/core/domain/interface/delivery"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/shipment"
	"go2o/core/module"
	"go2o/gen-code/thrift/define"
)

var _ define.ShipmentService = new(shipmentServiceImpl)

type shipmentServiceImpl struct {
	_repo         shipment.IShipmentRepo
	_deliveryRepo delivery.IDeliveryRepo
	_expressRepo  express.IExpressRepo
}

// 获取快递服务
func NewShipmentService(rep shipment.IShipmentRepo,
	deliveryRepo delivery.IDeliveryRepo, expressRepo express.IExpressRepo) *shipmentServiceImpl {
	return &shipmentServiceImpl{
		_repo:         rep,
		_deliveryRepo: deliveryRepo,
		_expressRepo:  expressRepo,
	}
}

// 创建一个配送覆盖的区域
func (s *shipmentServiceImpl) CreateCoverageArea(c *delivery.CoverageValue) (int32, error) {
	return s._deliveryRepo.SaveCoverageArea(c)
}

// 获取订单的发货单信息
func (s *shipmentServiceImpl) GetShipOrderOfOrder(orderId int64, sub bool) *shipment.ShipmentOrder {
	arr := s._repo.GetShipOrders(orderId, sub)
	if arr != nil && len(arr) > 0 {
		v := arr[0].Value()
		return &v
	}
	return nil
}

func (s *shipmentServiceImpl) GetLogisticFlowTrack(shipperCode string,
	logisticCode string, invert bool) (r *define.SShipOrderTrack, err error) {
	em := module.Get(module.M_EXPRESS).(*module.ExpressModule)
	flow, err := em.GetLogisticFlowTrack(shipperCode, logisticCode, invert)
	if err == nil {
		return s.logisticFlowTrackDto(flow), nil
	}
	return &define.SShipOrderTrack{
		Code:    1,
		Message: err.Error(),
	}, nil
}
func (s *shipmentServiceImpl) logisticFlowTrackDto(o *shipment.ShipOrderTrack) *define.SShipOrderTrack {
	if o == nil {
		return &define.SShipOrderTrack{
			Code:    1,
			Message: "无法获取物流信息",
		}
	}
	r := &define.SShipOrderTrack{
		LogisticCode: o.LogisticCode,
		ShipperName:  o.ShipperName,
		ShipperCode:  o.ShipperCode,
		ShipState:    o.ShipState,
		UpdateTime:   o.UpdateTime,
		Flows:        make([]*define.SShipFlow, 0),
	}
	for _, v := range o.Flows {
		r.Flows = append(r.Flows, &define.SShipFlow{
			Subject:    v.Subject,
			CreateTime: v.CreateTime,
			Remark:     v.Remark,
		})
	}
	return r
}

// 获取发货单的物流追踪信息,
// - shipOrderId:发货单编号
func (s *shipmentServiceImpl) ShipOrderLogisticTrack(shipOrderId int64, invert bool) (r *define.SShipOrderTrack, err error) {
	so := s._repo.GetShipmentOrder(shipOrderId)
	if so != nil {
		sp := s._expressRepo.GetExpressProvider(so.Value().SpId)
		if sp == nil {
			log.Println("[ Go2o][ Service][ Warning]: no such express provider id ", so.Value().SpId)
		} else {
			//spOrder = "462681586678"
			//sp.ApiCode = "ZTO"
			spOrder := so.Value().SpOrder
			r, err := s.GetLogisticFlowTrack(sp.ApiCode, spOrder, invert)
			r.ShipperName = sp.Name
			return r, err
		}
	}
	return nil, nil
}
