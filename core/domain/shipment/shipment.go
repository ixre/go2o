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
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/shipment"
	"time"
)

var _ shipment.IShipmentOrder = new(shipmentOrderImpl)

type shipmentOrderImpl struct {
	value   *shipment.ShipmentOrder
	rep     shipment.IShipmentRepo
	expRepo express.IExpressRepo
	expSp   *express.ExpressProvider
}

func NewShipmentOrder(v *shipment.ShipmentOrder, rep shipment.IShipmentRepo,
	expRepo express.IExpressRepo) shipment.IShipmentOrder {
	return &shipmentOrderImpl{
		value:   v,
		rep:     rep,
		expRepo: expRepo,
	}
}

// 获取聚合根编号
func (s *shipmentOrderImpl) GetAggregateRootId() int64 {
	return s.value.ID
}

// 获取值
func (s *shipmentOrderImpl) Value() shipment.ShipmentOrder {
	return *s.value
}

func (s *shipmentOrderImpl) getExpressProvider(spId int32) *express.ExpressProvider {
	if s.expSp == nil {
		s.expSp = s.expRepo.GetExpressProvider(spId)
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
	s.value.SpOrder = spOrderNo
	s.value.State = shipment.StatShipped
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
	return s.submit()
}

func (s *shipmentOrderImpl) submit() error {
	if s.GetAggregateRootId() > 0 {
		panic("shipment order created!")
	}
	id, err := util.I64Err(s.rep.SaveShipmentOrder(s.value))
	if err == nil {
		s.value.ID = id
		items := s.value.Items
		if items != nil && len(items) > 0 {
			for _, v := range items {
				v.ShipOrder = id
				v.ID, err = util.I64Err(s.rep.SaveShipmentItem(v))
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
	s.value.State = shipment.StatCompleted
	return s.save()
}

// 更新快递记录
func (s *shipmentOrderImpl) UpdateLog() error {
	panic(errors.New("not implement!"))
}

// 智能选择门店
func (t *shipmentOrderImpl) SmartChoiceShop(address string) (shop.IShop, error) {
	panic("not implement")
	/*

	   //todo: 应只选择线下实体店
	   //todo: AggregateRootId
	   dly := t.deliveryRepo.GetDelivery(-1)

	   lng, lat, err := lbs.GetLocation(address)
	   if err != nil {
	       return nil, errors.New("无法识别的地址：" + address)
	   }
	   var cov delivery.ICoverageArea = dly.GetNearestCoverage(lng, lat)
	   if cov == nil {
	       return nil, delivery.ErrNotCoveragedArea
	   }
	   shopId, _, err := dly.GetDeliveryInfo(cov.GetDomainId())
	   return t.mch.ShopManager().GetShop(shopId), err
	*/

}
