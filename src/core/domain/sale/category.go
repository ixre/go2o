/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:53
 * description :
 * history :
 */

package sale

import (
	"bytes"
	"go2o/src/core/domain/interface/sale"
	"strconv"
	"strings"
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
		if len(v.Url) == 0 || (strings.HasPrefix(val.Url, "/c-") && v.ParentId != val.ParentId) {
			val.Url = this.getAutomaticUrl(val.PartnerId, v.Id)
		}
	}
	return nil
}

func (this *Category) Save() (int, error) {
	id, err := this.rep.SaveCategory(this.value)
	if err == nil && this.GetDomainId() == 0 {
		this.value.Id = id
		if len(this.value.Url) == 0 {
			this.value.Url = this.getAutomaticUrl(this.value.PartnerId, id)
			return this.Save()
		}
	}
	return id, err
}

func (this *Category) getAutomaticUrl(partnerId, id int) string {
	var relCategories []*sale.ValueCategory = this.rep.GetRelationCategories(partnerId, id)
	var buf *bytes.Buffer = bytes.NewBufferString("/c")
	var l int = len(relCategories)
	for i := l; i > 0; i-- {
		buf.WriteString(strconv.Itoa(relCategories[i-1].Id))
	}
	buf.WriteString(".htm")
	return buf.String()
}
