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
	"go2o/gen-code/thrift/define"
)
var _ define.ShipmentService = new(shipmentServiceImpl)
type shipmentServiceImpl struct {
	_rep          shipment.IShipmentRepo
	_deliveryRepo delivery.IDeliveryRepo
	_shipRepo   shipment.IShipmentRepo
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
	panic("implement me")
}