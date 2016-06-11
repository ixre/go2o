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
	_rep   sale.ISaleLabelRep
	_mchId int
	_value *sale.Label
}

func NewSaleLabel(mchId int, value *sale.Label,
	rep sale.ISaleLabelRep) sale.ISaleLabel {
	return &saleLabelImpl{
		_rep:        rep,
		_mchId: mchId,
		_value:      value,
	}
}

func (this *saleLabelImpl) GetDomainId() int {
	if this._value != nil {
		return this._value.Id
	}
	return 0
}

func (this *saleLabelImpl) GetValue() *sale.Label {
	return this._value
}

// 是否为系统内置
func (this *saleLabelImpl) System() bool {
	return this._value.MerchantId == 0
}

// 设置值
func (this *saleLabelImpl) SetValue(v *sale.Label) error {
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
	this._value.MerchantId = this._mchId
	return this._rep.SaveSaleLabel(this._mchId, this._value)
}

// 获取标签下的商品
func (this *saleLabelImpl) GetValueGoods(sortBy string, begin, end int) []*valueobject.Goods {
	if begin < 0 || begin > end {
		begin = 0
	}
	if end <= 0 {
		end = 5
	}
	return this._rep.GetValueGoodsBySaleLabel(this._mchId,
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
	return this._rep.GetPagedValueGoodsBySaleLabel(this._mchId,
		this.GetDomainId(), sortBy, begin, end)
}

var _ sale.ILabelManager = new(labelManagerImpl)

type labelManagerImpl struct {
	_rep    sale.ISaleLabelRep
	_valRep valueobject.IValueRep
	_mchId  int
}

func NewLabelManager(mchId int, rep sale.ISaleLabelRep,
	valRep valueobject.IValueRep) sale.ILabelManager {
	c := &labelManagerImpl{
		_rep:    rep,
		_mchId:  mchId,
		_valRep: valRep,
	}
	return c.init()
}

func (this *labelManagerImpl) init() sale.ILabelManager {
	//mchConf := this._valRep.GetPlatformConf()
	//if !mchConf.MchGoodsCategory && this._mchId > 0 {

	//todo: mch sale label
	this._mchId = 0
	//}
	return this
}

// 初始化销售标签
func (this *labelManagerImpl) InitSaleLabels() error {
	if len(this.GetAllSaleLabels()) != 0 {
		return nil
	}

	arr := []sale.Label{
		sale.Label{
			TagName: "新品上架",
			TagCode: "new-goods",
		},
		sale.Label{
			TagName: "热销商品",
			TagCode: "hot-sales",
		},
		sale.Label{
			TagName: "特色商品",
			TagCode: "special-goods",
		},
		sale.Label{
			TagName: "优惠促销",
			TagCode: "prom-sales",
		},
		sale.Label{
			TagName: "尾品清仓",
			TagCode: "clean-goods",
		},
	}

	var err error
	for _, v := range arr {
		v.Enabled = 1
		v.MerchantId = this._mchId
		_, err = this.CreateSaleLabel(&v).Save()
	}

	return err
}

// 获取所有的销售标签
func (this *labelManagerImpl) GetAllSaleLabels() []sale.ISaleLabel {
	arr := this._rep.GetAllValueSaleLabels(this._mchId)
	var tags = make([]sale.ISaleLabel, len(arr))

	for i, v := range arr {
		tags[i] = this.CreateSaleLabel(v)
	}
	return tags
}

// 获取销售标签
func (this *labelManagerImpl) GetSaleLabel(id int) sale.ISaleLabel {
	return this._rep.GetSaleLabel(this._mchId, id)
}

// 根据Code获取销售标签
func (this *labelManagerImpl) GetSaleLabelByCode(code string) sale.ISaleLabel {
	v := this._rep.GetSaleLabelByCode(this._mchId, code)
	return this.CreateSaleLabel(v)
}

// 创建销售标签
func (this *labelManagerImpl) CreateSaleLabel(v *sale.Label) sale.ISaleLabel {
	if v == nil {
		return nil
	}
	v.MerchantId = this._mchId
	return this._rep.CreateSaleLabel(v)
}

// 删除销售标签
func (this *labelManagerImpl) DeleteSaleLabel(id int) error {
	v := this.GetSaleLabel(id)
	if v != nil {
		if v.System() {
			return sale.ErrInternalDisallow
		}
		return this._rep.DeleteSaleLabel(this._mchId, id)
	}
	return nil
}
