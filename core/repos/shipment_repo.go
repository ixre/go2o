/**
 * Copyright 2015 @ z3q.net.
 * name : shipment_rep
 * author : jarryliu
 * date : 2016-07-15 10:28
 * description :
 * history :
 */
package repos

import (
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/shipment"
	shipImpl "go2o/core/domain/shipment"
)

var _ shipment.IShipmentRepo = new(shipmentRepo)

type shipmentRepo struct {
	_expRepo express.IExpressRepo
	db.Connector
}

func NewShipmentRepo(conn db.Connector, expRepo express.IExpressRepo) *shipmentRepo {
	return &shipmentRepo{
		Connector: conn,
		_expRepo:  expRepo,
	}
}

// 创建发货单
func (s *shipmentRepo) CreateShipmentOrder(o *shipment.ShipmentOrder) shipment.IShipmentOrder {
	return shipImpl.NewShipmentOrder(o, s, s._expRepo)
}

func (s *shipmentRepo) getShipOrderById(id int64) *shipment.ShipmentOrder {
	e := &shipment.ShipmentOrder{}
	if s.GetOrm().Get(id, e) == nil {
		return e
	}
	return nil
}

// 获取发货单
func (s *shipmentRepo) GetShipmentOrder(id int64) shipment.IShipmentOrder {
	if e := s.getShipOrderById(id); e != nil {
		return s.CreateShipmentOrder(e)
	}
	return nil
}

// 获取订单对应的发货单
func (s *shipmentRepo) GetShipOrders(orderId int64, sub bool) []shipment.IShipmentOrder {
	var list []*shipment.ShipmentOrder
	if sub {
		s.GetOrm().Select(&list, "sub_orderid= $1", orderId)
	} else {
		s.GetOrm().Select(&list, "order_id= $1", orderId)
	}
	orders := make([]shipment.IShipmentOrder, len(list))
	for i, v := range list {
		orders[i] = s.CreateShipmentOrder(v)
	}
	return orders
}

// 保存发货单
func (s *shipmentRepo) SaveShipmentOrder(o *shipment.ShipmentOrder) (int, error) {
	return orm.Save(s.GetOrm(), o, int(o.ID))
}

// 保存发货商品项
func (s *shipmentRepo) SaveShipmentItem(v *shipment.Item) (int, error) {
	return orm.Save(s.GetOrm(), v, int(v.ID))
}

// 删除发货单
func (s *shipmentRepo) DeleteShipmentOrder(id int64) error {
	return s.GetOrm().DeleteByPk(&shipment.ShipmentOrder{}, id)
}
