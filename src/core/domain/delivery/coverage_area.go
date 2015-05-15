/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-14 15:42
 * description :
 * history :
 */
package delivery

import (
	"errors"
	"go2o/src/core/domain/interface/delivery"
	"go2o/src/core/infrastructure/lbs"
)

var _ delivery.ICoverageArea = new(CoverageArea)

type CoverageArea struct {
	value *delivery.CoverageValue
	rep   delivery.IDeliveryRep
}

func newCoverageArea(v *delivery.CoverageValue, rep delivery.IDeliveryRep) delivery.ICoverageArea {
	return &CoverageArea{
		value: v,
		rep:   rep,
	}
}

// 是否可以配送
// 返回是否可以配送，以及距离(米)
func (this *CoverageArea) CanDeliver(lng, lat float64) (bool, int) {
	distance := lbs.GetLocDistance(
		this.value.Lng, this.value.Lat, lng, lat)
	i := int(distance)
	return i <= this.value.Radius*1000, i
}

// 是否可以配送
// 返回是否可以配送，以及距离(米)
func (this *CoverageArea) CanDeliverTo(address string) (bool, int) {
	lng, lat, err := lbs.GetLocation(address)
	if err != nil {
		return false, -1
	}
	return this.CanDeliver(lng, lat)
}

func (this *CoverageArea) GetDomainId() int {
	return this.value.Id
}

func (this *CoverageArea) GetValue() delivery.CoverageValue {
	return *this.value
}

func (this *CoverageArea) SetValue(v *delivery.CoverageValue) error {
	if v.Id == this.value.Id && v.Id > 0 {
		this.value = v
		return nil
	}
	return errors.New("no such value")
}

func (this *CoverageArea) Save() (int, error) {
	return this.rep.SaveCoverageArea(this.value)
}
