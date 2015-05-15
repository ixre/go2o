/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-14 15:37
 * description :
 * history :
 */
package delivery

import (
	"go2o/src/core/domain/interface/delivery"
	"go2o/src/core/infrastructure/domain"
)

var _ delivery.IDelivery = new(Delivery)

type Delivery struct {
	_id  int
	_rep delivery.IDeliveryRep
}

func NewDelivery(id int, dlvRep delivery.IDeliveryRep) delivery.IDelivery {
	return &Delivery{
		_id:  id,
		_rep: dlvRep,
	}
}

// 返回聚合编号
func (this *Delivery) GetAggregateRootId() int {
	return this._id
}

// 等同于GetAggregateRootId()
func (this *Delivery) GetPartnerId() int {
	return this._id
}

// 根据地址获取地区(可能会有重复的区名)
func (this *Delivery) GetArea(addr string) ([]*delivery.AreaValue, error) {
	name, err := domain.GetAreaName(addr)
	if err != nil {
		return nil, err
	}
	return this._rep.GetAreaByArea(name), nil
}

//　获取覆盖区域
func (this *Delivery) GetCoverageArea(id int) delivery.ICoverageArea {
	val := this._rep.GetCoverageArea(this._id, id)
	return newCoverageArea(val, this._rep)
}

// 获取最近的配送区域
func (this *Delivery) GetNearestCoverage(lng, lat float64) delivery.ICoverageArea {
	var distance int
	var nearest delivery.ICoverageArea = nil
	areas := this.FindCoverageAreas(lng, lat)

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
func (this *Delivery) FindSingleCoverageArea(lng, lat float64) delivery.ICoverageArea {
	var covers []*delivery.CoverageValue = this._rep.GetAllCoverageAreas(this._id)
	if covers != nil {
		return newCoverageArea(covers[0], this._rep)
	}
	return nil
}

// 查找所有所在的区域
func (this *Delivery) FindCoverageAreas(lng, lat float64) []delivery.ICoverageArea {
	var covers []*delivery.CoverageValue = this._rep.GetAllCoverageAreas(this._id)
	if covers != nil {
		var list []delivery.ICoverageArea = make([]delivery.ICoverageArea, len(covers))
		for i, v := range covers {
			list[i] = newCoverageArea(v, this._rep)
		}
		return list
	}
	return nil
}

// 获取配送信息
func (this *Delivery) GetDeliveryInfo(coverageId int) (shopId, deliverUsrId int, err error) {
	v := this._rep.GetDeliveryBind(this.GetAggregateRootId(), coverageId)
	if v != nil {
		return v.ShopId, v.DeliverUsrId, nil
	}
	return -1, -1, nil
}
