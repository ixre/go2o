/**
 * Copyright 2015 @ z3q.net.
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
	_merchantId int
	_value     *sale.ValueSaleTag
}

func NewSaleTag(merchantId int, value *sale.ValueSaleTag,
	rep sale.ISaleTagRep) sale.ISaleTag {
	return &SaleTag{
		_rep:       rep,
		_merchantId: merchantId,
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

// 是否为系统内置
func (this *SaleTag) System() bool {
	return this._value.IsInternal == 1
}

// 设置值
func (this *SaleTag) SetValue(v *sale.ValueSaleTag) error {
	if v != nil {
		// 如果为系统内置，不能修改名称
		if !this.System() {
			this._value.Enabled = v.Enabled
			this._value.TagCode = v.TagCode
		}
		this._value.TagName = v.TagName
		this._value.GoodsImage = v.GoodsImage
		if len(v.TagCode) == 0 {
			this._value.TagCode = v.TagCode
		}
	}
	return nil
}

func (this *SaleTag) Save() (int, error) {
	this._value.MerchantId = this._merchantId
	return this._rep.SaveSaleTag(this._merchantId, this._value)
}

// 获取标签下的商品
func (this *SaleTag) GetValueGoods(sortBy string, begin, end int) []*valueobject.Goods {
	if begin < 0 || begin > end {
		begin = 0
	}
	if end <= 0 {
		end = 5
	}
	return this._rep.GetValueGoodsBySaleTag(this._merchantId,
		this._value.Id, sortBy, begin, end)
}

// 获取标签下的分页商品
func (this *SaleTag) GetPagedValueGoods(sortBy string, begin, end int) (int, []*valueobject.Goods) {
	if begin < 0 || begin > end {
		begin = 0
	}
	if end <= 0 {
		end = 5
	}
	return this._rep.GetPagedValueGoodsBySaleTag(this._merchantId,
		this.GetDomainId(), sortBy, begin, end)
}
