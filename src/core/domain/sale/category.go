/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:53
 * description :
 * history :
 */

package sale

import (
	"go2o/src/core/domain/interface/sale"
)

var _ sale.ICategory = new(Category)

type Category struct {
	value *sale.ValueCategory
	rep   sale.ISaleRep
}

func newCategory(saleRep sale.ISaleRep, v *sale.ValueCategory) sale.ICategory {
	return &Category{
		value: v,
		rep:   saleRep,
	}
}

func (this *Category) GetDomainId() int {
	return this.value.Id
}

func (this *Category) GetValue() sale.ValueCategory {
	return *this.value
}

func (this *Category) SetValue(v *sale.ValueCategory) error {
	val := this.value
	if val.Id == v.Id {
		val.Description = v.Description
		val.Enabled = v.Enabled
		val.Name = v.Name
		val.OrderIndex = v.OrderIndex
		val.ParentId = v.ParentId
	}
	return nil
}

func (this *Category) Save() (int, error) {
	return this.rep.SaveCategory(this.value)
}
