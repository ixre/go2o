/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2015-02-16 10:28
 * description :
 * history :
 */
package repository

import (
	"github.com/atnet/gof/db"
	deliverImpl "go2o/src/core/domain/delivery"
	"go2o/src/core/domain/interface/delivery"
)

var _ delivery.IDeliveryRep = new(DeliveryRep)

type DeliveryRep struct {
	db.Connector
}

func NewDeliverRep(c db.Connector) delivery.IDeliveryRep {
	return &DeliveryRep{
		Connector: c,
	}
}

// 获取配送
func (this *DeliveryRep) GetDelivery(id int) delivery.IDelivery {
	return deliverImpl.NewDelivery(id, this)
}

// 根据区名获取区域
func (this *DeliveryRep) GetAreaByArea(name string) []*delivery.AreaValue {
	arr := make([]*delivery.AreaValue, 0)
	if err := this.Connector.GetOrm().Select(arr, "name LIKE '%?%'", name); err == nil {
		return arr
	}
	return nil
}

// 保存覆盖区域
func (this *DeliveryRep) SaveCoverageArea(v *delivery.CoverageValue) (int, error) {
	orm := this.Connector.GetOrm()
	if v.Id <= 0 {
		_, _, err := orm.Save(nil, v)
		if err == nil {
			this.Connector.ExecScalar(`SELECT MAX(id) FROM dlv_coverage`, &v.Id)
		}
		return v.Id, err
	} else {
		_, _, err := orm.Save(nil, v)
		return v.Id, err
	}
}

// 获取覆盖区域
func (this *DeliveryRep) GetCoverageArea(areaId, id int) *delivery.CoverageValue {
	e := new(delivery.CoverageValue)
	err := this.Connector.GetOrm().GetBy(e, "id=? AND area_id=?", id, areaId)
	if err != nil {
		return nil
	}
	return e
}

// 获取所有的覆盖区域
func (this *DeliveryRep) GetAllCoverageAreas(areaId int) []*delivery.CoverageValue {
	e := make([]*delivery.CoverageValue, 0)
	err := this.Connector.GetOrm().Select(e, "area_id=?", areaId)
	if err != nil {
		return nil
	}
	return e
}

// 获取配送绑定
func (this *DeliveryRep) GetDeliveryBind(partnerId, coverageId int) *delivery.PartnerDeliverBind {
	e := new(delivery.PartnerDeliverBind)
	err := this.Connector.GetOrm().GetBy(e, "partner_id=? AND coverage_id=?", partnerId, coverageId)
	if err != nil {
		return nil
	}
	return e
}
