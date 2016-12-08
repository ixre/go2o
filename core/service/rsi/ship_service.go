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
)

type shipmentService struct {
	_rep         shipment.IShipmentRepo
	_deliveryRepo delivery.IDeliveryRepo
}

// 获取快递服务
func NewShipmentService(rep shipment.IShipmentRepo,
	deliveryRepo delivery.IDeliveryRepo) *shipmentService {
	return &shipmentService{
		_rep:         rep,
		_deliveryRepo: deliveryRepo,
	}
}

// 创建一个配送覆盖的区域
func (s *shipmentService) CreateCoverageArea(c *delivery.CoverageValue) (int32, error) {
	return s._deliveryRepo.SaveCoverageArea(c)
}

// 获取订单的发货单信息
func (s *shipmentService) GetShipOrderOfOrder(orderId int32) *shipment.ShipmentOrder {
	arr := s._rep.GetOrders(orderId)
	if arr != nil && len(arr) > 0 {
		v := arr[0].Value()
		return &v
	}
	return nil
}
