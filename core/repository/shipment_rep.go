/**
 * Copyright 2015 @ z3q.net.
 * name : shipment_rep
 * author : jarryliu
 * date : 2016-07-15 10:28
 * description :
 * history :
 */
package repository

import (
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/shipment"
	shipImpl "go2o/core/domain/shipment"
)

var _ shipment.IShipmentRep = new(shipmentRep)

type shipmentRep struct {
	_expRep express.IExpressRep
	db.Connector
}

func NewShipmentRep(conn db.Connector, expRep express.IExpressRep) *shipmentRep {
	return &shipmentRep{
		Connector: conn,
		_expRep:   expRep,
	}
}

// 创建发货单
func (s *shipmentRep) CreateShipmentOrder(o *shipment.ShipmentOrder) shipment.IShipmentOrder {
	return shipImpl.NewShipmentOrder(o, s, s._expRep)
}

func (s *shipmentRep) GetShipOrderById(id int) *shipment.ShipmentOrder {
	e := &shipment.ShipmentOrder{}
	if s.GetOrm().Get(id, &e) == nil {
		return e
	}
	return nil
}

// 获取发货单
func (s *shipmentRep) GetShipmentOrder(id int) shipment.IShipmentOrder {
	if e := s.GetShipOrderById(id); e != nil {
		return s.CreateShipmentOrder(e)
	}
	return nil
}

// 获取订单对应的发货单
func (s *shipmentRep) GetOrders(orderId int) []shipment.IShipmentOrder {
	list := []*shipment.ShipmentOrder{}
	s.GetOrm().Select(&list, "order_id=?", orderId)
	orders := make([]shipment.IShipmentOrder, len(list))
	for i, v := range list {
		orders[i] = s.CreateShipmentOrder(v)
	}
	return orders
}

// 保存发货单
func (s *shipmentRep) SaveShipmentOrder(o *shipment.ShipmentOrder) (int, error) {
	return orm.Save(s.GetOrm(), o, o.Id)
}

// 删除发货单
func (s *shipmentRep) DeleteShipmentOrder(id int) error {
	return s.GetOrm().DeleteByPk(&shipment.ShipmentOrder{}, id)
}
