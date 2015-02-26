/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-14 15:42
 * description :
 * history :
 */
package delivery

import (
	"go2o/core/domain/interface/delivery"
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

func (this *CoverageArea) GetDomainId() int {
	return this.value.Id
}

func (this *CoverageArea) GetValue() CoverageValue {
	return *this.value
}

func (this *CoverageArea) SetValue(v *CoverageValue) error {
	if v.Id == this.value.Id && v.Id > 0 {
		this.value = v
		return nil
	}
	return errors.New("no such value")
}

func (this *CoverageArea) Save() (int, error) {
	return this.rep.SaveCoverageArea(this.value)
}
