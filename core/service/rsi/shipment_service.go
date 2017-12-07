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
	"go2o/core/domain/interface/delivery"
	"go2o/core/domain/interface/shipment"
	"go2o/core/module"
	"go2o/gen-code/thrift/define"
)

var _ define.ShipmentService = new(shipmentServiceImpl)

type shipmentServiceImpl struct {
	_rep          shipment.IShipmentRepo
	_deliveryRepo delivery.IDeliveryRepo
	_shipRepo     shipment.IShipmentRepo
}

// 获取快递服务
func NewShipmentService(rep shipment.IShipmentRepo,
	deliveryRepo delivery.IDeliveryRepo) *shipmentServiceImpl {
	return &shipmentServiceImpl{
		_rep:          rep,
		_deliveryRepo: deliveryRepo,
	}
}

// 创建一个配送覆盖的区域
func (s *shipmentServiceImpl) CreateCoverageArea(c *delivery.CoverageValue) (int32, error) {
	return s._deliveryRepo.SaveCoverageArea(c)
}

// 获取订单的发货单信息
func (s *shipmentServiceImpl) GetShipOrderOfOrder(orderId int64, sub bool) *shipment.ShipmentOrder {
	arr := s._rep.GetShipOrders(orderId, sub)
	if arr != nil && len(arr) > 0 {
		v := arr[0].Value()
		return &v
	}
	return nil
}

func (s *shipmentServiceImpl) GetLogisticFlowTrace(shipperCode string, logisticCode string) (r *define.SShipOrderTrace, err error) {
	em := module.Get(module.M_EXPRESS).(*module.ExpressModule)
	flow, err := em.GetLogisticFlowTrace(shipperCode, logisticCode)
	if err == nil {
		return s.logisticFlowTraceDto(flow), nil
	}
	return &define.SShipOrderTrace{
		Code:    1,
		Message: err.Error(),
	}, nil
}
func (s *shipmentServiceImpl) logisticFlowTraceDto(o *shipment.ShipOrderTrace) *define.SShipOrderTrace {
	if o == nil {
		return &define.SShipOrderTrace{
			Code:    1,
			Message: "无法获取物流信息",
		}
	}
	r := &define.SShipOrderTrace{
		LogisticCode: o.LogisticCode,
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
