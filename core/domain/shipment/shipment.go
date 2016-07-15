/**
 * Copyright 2015 @ z3q.net.
 * name : shipment
 * author : jarryliu
 * date : 2016-07-15 10:08
 * description :
 * history :
 */
package shipment

import (
	"errors"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/shipment"
	"time"
)

var _ shipment.IShipmentOrder = new(shipmentOrderImpl)

type shipmentOrderImpl struct {
	_value  *shipment.ShipmentOrder
	_rep    shipment.IShipmentRep
	_expRep express.IExpressRep
	_expSp  *express.ExpressProvider
}

func NewShipmentOrder(v *shipment.ShipmentOrder, rep shipment.IShipmentRep,
	expRep express.IExpressRep) shipment.IShipmentOrder {
	return &shipmentOrderImpl{
		_value:  v,
		_rep:    rep,
		_expRep: expRep,
	}
}

// 获取聚合根编号
func (s *shipmentOrderImpl) GetAggregateRootId() int {
	return s._value.Id
}

// 获取值
func (s *shipmentOrderImpl) Value() shipment.ShipmentOrder {
	return *s._value
}

func (s *shipmentOrderImpl) getExpressProvider(spId int) *express.ExpressProvider {
	if s._expSp == nil {
		s._expSp = s._expRep.GetExpressProvider(spId)
	}
	return s._expSp
}

// 发货
func (s *shipmentOrderImpl) Ship(spId int, spOrderNo string) error {
	if e := s.getExpressProvider(spId); e == nil || e.Enabled != 1 {
		return express.ErrNotSupportProvider
	}
	if s._value.SpId != spId {
		s._expSp = nil
		s._value.SpId = spId
	}
	s._value.SpOrderNo = spOrderNo
	s._value.ShipTime = time.Now().Unix()
	return s.save()
}

// 保存
func (s *shipmentOrderImpl) save() error {
	s._value.UpdateTime = time.Now().Unix()
	id, err := s._rep.SaveShipmentOrder(s._value)
	if s.GetAggregateRootId() <= 0 {
		s._value.Id = id
	}
	return err
}

// 发货完成
func (s *shipmentOrderImpl) Shipped() error {
	s._value.IsShipped = 1
	return s.save()
}

// 更新快递记录
func (s *shipmentOrderImpl) UpdateLog() error {
	panic(errors.New("not implement!"))
}
