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
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/valueobject"
)

var _ sale.ISaleLabel = new(saleLabelImpl)

type saleLabelImpl struct {
	_rep        sale.ISaleTagRep
	_merchantId int
	_value      *sale.SaleLabel
}

func NewSaleLabel(mchId int, value *sale.SaleLabel,
	rep sale.ISaleTagRep) sale.ISaleLabel {
	return &saleLabelImpl{
		_rep:        rep,
		_merchantId: mchId,
		_value:      value,
	}
}

func (this *saleLabelImpl) GetDomainId() int {
	if this._value != nil {
		return this._value.Id
	}
	return 0
}

func (this *saleLabelImpl) GetValue() *sale.SaleLabel {
	return this._value
}

// 是否为系统内置
func (this *saleLabelImpl) System() bool {
	return this._value.MerchantId == 0
}

// 设置值
func (this *saleLabelImpl) SetValue(v *sale.SaleLabel) error {
	if v != nil {
		// 如果为系统内置，不能修改名称
		if !this.System() {
			this._value.Enabled = v.Enabled
			this._value.TagCode = v.TagCode
		}
		this._value.TagName = v.TagName
		this._value.LabelImage = v.LabelImage
		if len(v.TagCode) == 0 {
			this._value.TagCode = v.TagCode
		}
	}
	return nil
}

func (this *saleLabelImpl) Save() (int, error) {
	this._value.MerchantId = this._merchantId
	return this._rep.SaveSaleTag(this._merchantId, this._value)
}

// 获取标签下的商品
func (this *saleLabelImpl) GetValueGoods(sortBy string, begin, end int) []*valueobject.Goods {
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
func (this *saleLabelImpl) GetPagedValueGoods(sortBy string, begin, end int) (int, []*valueobject.Goods) {
	if begin < 0 || begin > end {
		begin = 0
	}
	if end <= 0 {
		end = 5
	}
	return this._rep.GetPagedValueGoodsBySaleTag(this._merchantId,
		this.GetDomainId(), sortBy, begin, end)
}
