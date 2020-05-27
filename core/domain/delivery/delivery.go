/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 15:37
 * description :
 * history :
 */
package delivery

import (
	"go2o/core/domain/interface/delivery"
	"go2o/core/infrastructure/domain"
)

var _ delivery.IDelivery = new(Delivery)

type Delivery struct {
	id  int32
	rep delivery.IDeliveryRepo
}

func NewDelivery(id int32, dlvRepo delivery.IDeliveryRepo) delivery.IDelivery {
	return &Delivery{
		id:  id,
		rep: dlvRepo,
	}
}

// 返回聚合编号
func (d *Delivery) GetAggregateRootId() int32 {
	return d.id
}

// 等同于GetAggregateRootId()
func (d *Delivery) GetMerchantId() int32 {
	return d.id
}

// 根据地址获取地区(可能会有重复的区名)
func (d *Delivery) GetArea(addr string) ([]*delivery.AreaValue, error) {
	name, err := domain.GetAreaName(addr)
	if err != nil {
		return nil, err
	}
	return d.rep.GetAreaByArea(name), nil
}

//　获取覆盖区域
func (d *Delivery) GetCoverageArea(id int32) delivery.ICoverageArea {
	val := d.rep.GetCoverageArea(d.id, id)
	return newCoverageArea(val, d.rep)
}

// 获取最近的配送区域
func (d *Delivery) GetNearestCoverage(lng, lat float64) delivery.ICoverageArea {
	var distance int
	var nearest delivery.ICoverageArea = nil
	areas := d.FindCoverageAreas(lng, lat)

	// 获取最近的门店
	for _, v := range areas {
		if b, d := v.CanDeliver(lng, lat); b {
			return v
		} else {
			if d < distance || distance == 0 {
				d = distance
				nearest = v
			}
		}
	}
	return nearest
}

// 查看单个所在的区域
func (d *Delivery) FindSingleCoverageArea(lng, lat float64) delivery.ICoverageArea {
	var covers []*delivery.CoverageValue = d.rep.GetAllCoverageAreas(d.id)
	if covers != nil {
		return newCoverageArea(covers[0], d.rep)
	}
	return nil
}

// 查找所有所在的区域
func (d *Delivery) FindCoverageAreas(lng, lat float64) []delivery.ICoverageArea {
	var covers []*delivery.CoverageValue = d.rep.GetAllCoverageAreas(d.id)
	if covers != nil {
		var list []delivery.ICoverageArea = make([]delivery.ICoverageArea, len(covers))
		for i, v := range covers {
			list[i] = newCoverageArea(v, d.rep)
		}
		return list
	}
	return nil
}

// 获取配送信息
func (d *Delivery) GetDeliveryInfo(coverageId int32) (shopId, deliverUsrId int32, err error) {
	v := d.rep.GetDeliveryBind(d.GetAggregateRootId(), coverageId)
	if v != nil {
		return v.ShopId, v.DeliverUsrId, nil
	}
	return -1, -1, nil
}
