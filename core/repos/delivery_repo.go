/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2015-02-16 10:28
 * description :
 * history :
 */
package repos

import (
	deliverImpl "github.com/ixre/go2o/core/domain/delivery"
	"github.com/ixre/go2o/core/domain/interface/delivery"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

var _ delivery.IDeliveryRepo = new(deliveryRepo)

type deliveryRepo struct {
	db.Connector
	o orm.Orm
}

func NewDeliverRepo(o orm.Orm) delivery.IDeliveryRepo {
	return &deliveryRepo{
		Connector: o.Connector(),
		o:         o,
	}
}

// 获取配送
func (this *deliveryRepo) GetDelivery(id int32) delivery.IDeliveryAggregateRoot {
	return deliverImpl.NewDelivery(id, this)
}

// 根据区名获取区域
func (this *deliveryRepo) GetAreaByArea(name string) []*delivery.AreaValue {
	arr := make([]*delivery.AreaValue, 0)
	if err := this.o.Select(&arr, "name LIKE $1",
		"%"+name+"%"); err == nil {
		return arr
	}
	return nil
}

// 保存覆盖区域
func (this *deliveryRepo) SaveCoverageArea(v *delivery.CoverageValue) (int32, error) {
	return orm.I32(orm.Save(this.o, v, int(v.Id)))
}

// 获取覆盖区域
func (this *deliveryRepo) GetCoverageArea(areaId, id int32) *delivery.CoverageValue {
	e := new(delivery.CoverageValue)
	err := this.o.GetBy(e, "id= $1 AND area_id= $2", id, areaId)
	if err != nil {
		return nil
	}
	return e
}

// 获取所有的覆盖区域
func (this *deliveryRepo) GetAllCoverageAreas(areaId int32) []*delivery.CoverageValue {
	e := make([]*delivery.CoverageValue, 0)
	err := this.o.Select(&e, "area_id= $1", areaId)
	if err != nil {
		return nil
	}
	return e
}

// 获取配送绑定
func (this *deliveryRepo) GetDeliveryBind(mchId, coverageId int32) *delivery.MerchantDeliverBind {
	e := new(delivery.MerchantDeliverBind)
	err := this.o.GetBy(e, "merchant_id= $1 AND coverage_id= $2", mchId, coverageId)
	if err != nil {
		return nil
	}
	return e
}
