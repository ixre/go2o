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
	value  *shipment.ShipmentOrder
	rep    shipment.IShipmentRep
	expRep express.IExpressRep
	expSp  *express.ExpressProvider
}

func NewShipmentOrder(v *shipment.ShipmentOrder, rep shipment.IShipmentRep,
	expRep express.IExpressRep) shipment.IShipmentOrder {
	return &shipmentOrderImpl{
		value:  v,
		rep:    rep,
		expRep: expRep,
	}
}

// 获取聚合根编号
func (s *shipmentOrderImpl) GetAggregateRootId() int32 {
	return s.value.Id
}

// 获取值
func (s *shipmentOrderImpl) Value() shipment.ShipmentOrder {
	return *s.value
}

func (s *shipmentOrderImpl) getExpressProvider(spId int32) *express.ExpressProvider {
	if s.expSp == nil {
		s.expSp = s.expRep.GetExpressProvider(spId)
	}
	return s.expSp
}

// 发货
func (s *shipmentOrderImpl) Ship(spId int32, spOrderNo string) error {
	if e := s.getExpressProvider(spId); e == nil || e.Enabled != 1 {
		return express.ErrNotSupportProvider
	}
	if s.value.SpId != spId {
		s.expSp = nil
		s.value.SpId = spId
	}
	s.value.SpOrderNo = spOrderNo
	s.value.Stat = shipment.StatShipped
	s.value.ShipTime = time.Now().Unix()
	return s.save()
}

// 保存
func (s *shipmentOrderImpl) save() error {
	s.value.UpdateTime = time.Now().Unix()
	if s.GetAggregateRootId() > 0 {
		_, err := s.rep.SaveShipmentOrder(s.value)
		return err
	}
	id, err := s.rep.SaveShipmentOrder(s.value)
	if err == nil {
		s.value.Id = id
		items := s.value.Items
		if items != nil && len(items) > 0 {
			for _, v := range items {
				v.OrderId = id
				v.Id, err = s.rep.SaveShipmentItem(v)
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}

// 发货完成
func (s *shipmentOrderImpl) Completed() error {
	s.value.Stat = shipment.StatCompleted
	return s.save()
}

// 更新快递记录
func (s *shipmentOrderImpl) UpdateLog() error {
	panic(errors.New("not implement!"))
}
