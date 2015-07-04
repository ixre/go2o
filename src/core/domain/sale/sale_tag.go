/**
 * Copyright 2015 @ S1N1 Team.
 * name : sale_tag
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

import (
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/domain/interface/valueobject"
)

var _ sale.ISaleTag = new(SaleTag)

type SaleTag struct {
	_rep       sale.ISaleTagRep
	_partnerId int
	_value     *sale.ValueSaleTag
}

func NewSaleTag(partnerId int, value *sale.ValueSaleTag, rep sale.ISaleTagRep) sale.ISaleTag {
	return &SaleTag{
		_rep:       rep,
		_partnerId: partnerId,
		_value:     value,
	}
}

func (this *SaleTag) GetDomainId() int {
	if this._value != nil {
		return this._value.Id
	}
	return 0
}

func (this *SaleTag) GetValue() *sale.ValueSaleTag {
	return this._value
}

func (this *SaleTag) SetValue(v *sale.ValueSaleTag) error {
	if v != nil {
		this._value.Enabled = v.Enabled
		this._value.GoodsImage = v.GoodsImage
		if len(v.TagCode) == 0 {
			this._value.TagCode = v.TagCode
		}
		this._value.TagName = v.TagName
	}
	return nil
}

func (this *SaleTag) Save() (int, error) {
	this._value.PartnerId = this._partnerId
	return this._rep.SaveSaleTag(this._partnerId, this._value)
}

// 获取标签下的商品
func (this *SaleTag) GetValueGoods(begin, end int) []*valueobject.Goods {
	if begin < 0 || begin > end {
		begin = 0
	}
	if end <= 0 {
		end = 5
	}
	return this._rep.GetValueGoodsBySaleTag(this._partnerId, this._value.Id, begin, end)
}

// 获取标签下的分页商品
func (this *SaleTag) GetPagedValueGoods(begin, end int) (int, []*valueobject.Goods) {
	if begin < 0 || begin > end {
		begin = 0
	}
	if end <= 0 {
		end = 5
	}
	return this._rep.GetPagedValueGoodsBySaleTag(this._partnerId, this._value.Id, begin, end)
}
