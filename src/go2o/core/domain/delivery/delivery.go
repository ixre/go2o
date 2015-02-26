/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-14 15:37
 * description :
 * history :
 */
package delivery

import (
	"go2o/core/domain/interface/delivery"
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
func (this *Delivery) GetArea(addr string)[]*delivery.AreaValue{

    return nil
}

//　获取覆盖区域
func (this *Delivery)  GetCoverageArea(id int) delivery.ICoverageArea{
    val := this._rep.GetCoverageArea(this._id,id)
    return newCoverageArea(val,this._rep)
}


// 查看单个所在的区域
func (this *Delivery) FindSingleCoverageArea(lng, lat float32) delivery.ICoverageArea{
    _ = this._rep.GetAllConverageAreas(this._id)
    return nil
}

// 查找所有所在的区域
func (this *Delivery) FindCoverageAreas(lng, lat float32) []delivery.ICoverageArea{
    return nil
}