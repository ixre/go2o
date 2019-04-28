/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2015-02-16 10:28
 * description :
 * history :
 */
package repos

import (
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	deliverImpl "go2o/core/domain/delivery"
	"go2o/core/domain/interface/delivery"
)

var _ delivery.IDeliveryRepo = new(deliveryRepo)

type deliveryRepo struct {
	db.Connector
}

func NewDeliverRepo(c db.Connector) delivery.IDeliveryRepo {
	return &deliveryRepo{
		Connector: c,
	}
}

// 获取配送
func (this *deliveryRepo) GetDelivery(id int32) delivery.IDelivery {
	return deliverImpl.NewDelivery(id, this)
}

// 根据区名获取区域
func (this *deliveryRepo) GetAreaByArea(name string) []*delivery.AreaValue {
	arr := make([]*delivery.AreaValue, 0)
	if err := this.Connector.GetOrm().Select(&arr, "name LIKE '%?%'", name); err == nil {
		return arr
	}
	return nil
}

// 保存覆盖区域
func (this *deliveryRepo) SaveCoverageArea(v *delivery.CoverageValue) (int32, error) {
	return orm.I32(orm.Save(this.GetOrm(), v, int(v.Id)))
}

// 获取覆盖区域
func (this *deliveryRepo) GetCoverageArea(areaId, id int32) *delivery.CoverageValue {
	e := new(delivery.CoverageValue)
	err := this.Connector.GetOrm().GetBy(e, "id=? AND area_id=?", id, areaId)
	if err != nil {
		return nil
	}
	return e
}

// 获取所有的覆盖区域
func (this *deliveryRepo) GetAllCoverageAreas(areaId int32) []*delivery.CoverageValue {
	e := make([]*delivery.CoverageValue, 0)
	err := this.Connector.GetOrm().Select(&e, "area_id=?", areaId)
	if err != nil {
		return nil
	}
	return e
}

// 获取配送绑定
func (this *deliveryRepo) GetDeliveryBind(mchId, coverageId int32) *delivery.MerchantDeliverBind {
	e := new(delivery.MerchantDeliverBind)
	err := this.Connector.GetOrm().GetBy(e, "merchant_id=? AND coverage_id=?", mchId, coverageId)
	if err != nil {
		return nil
	}
	return e
}
