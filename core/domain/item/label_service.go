/**
 * Copyright 2015 @ to2.net.
 * name : sale_tag
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package item

import (
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/valueobject"
)

var _ item.ISaleLabel = new(saleLabelImpl)

type saleLabelImpl struct {
	rep   item.ISaleLabelRepo
	mchId int32
	value *item.Label
}

func NewSaleLabel(mchId int32, value *item.Label,
	rep item.ISaleLabelRepo) item.ISaleLabel {
	return &saleLabelImpl{
		rep:   rep,
		mchId: mchId,
		value: value,
	}
}

func (l *saleLabelImpl) GetDomainId() int32 {
	if l.value != nil {
		return l.value.Id
	}
	return 0
}

func (l *saleLabelImpl) GetValue() *item.Label {
	return l.value
}

// 是否为系统内置
func (l *saleLabelImpl) System() bool {
	return l.value.MerchantId == 0
}

// 设置值
func (l *saleLabelImpl) SetValue(v *item.Label) error {
	if v != nil {
		// 如果为系统内置，不能修改名称
		if !l.System() {
			l.value.Enabled = v.Enabled
			l.value.TagCode = v.TagCode
		}
		l.value.TagName = v.TagName
		l.value.LabelImage = v.LabelImage
		if len(v.TagCode) == 0 {
			l.value.TagCode = v.TagCode
		}
	}
	return nil
}

func (l *saleLabelImpl) Save() (int32, error) {
	l.value.MerchantId = l.mchId
	return l.rep.SaveSaleLabel(l.mchId, l.value)
}

// 获取标签下的商品
func (l *saleLabelImpl) GetValueGoods(sortBy string,
	begin, end int) []*valueobject.Goods {
	if begin < 0 || begin > end {
		begin = 0
	}
	if end <= 0 {
		end = 5
	}
	return l.rep.GetValueGoodsBySaleLabel(l.mchId,
		l.value.Id, sortBy, begin, end)
}

// 获取标签下的分页商品
func (l *saleLabelImpl) GetPagedValueGoods(sortBy string,
	begin, end int) (int, []*valueobject.Goods) {
	if begin < 0 || begin > end {
		begin = 0
	}
	if end <= 0 {
		end = 5
	}
	return l.rep.GetPagedValueGoodsBySaleLabel(l.mchId,
		l.GetDomainId(), sortBy, begin, end)
}

var _ item.ILabelService = new(labelManagerImpl)

type labelManagerImpl struct {
	_rep     item.ISaleLabelRepo
	_valRepo valueobject.IValueRepo
	_mchId   int32
}

func NewLabelManager(mchId int32, rep item.ISaleLabelRepo,
	valRepo valueobject.IValueRepo) item.ILabelService {
	c := &labelManagerImpl{
		_rep:     rep,
		_mchId:   mchId,
		_valRepo: valRepo,
	}
	return c.init()
}

func (l *labelManagerImpl) init() item.ILabelService {
	//mchConf := l._valRepo.GetPlatformConf()
	//if !mchConf.MchGoodsCategory && l._mchId > 0 {

	//todo: mch sale label
	l._mchId = 0
	//}
	return l
}

// 初始化销售标签
func (l *labelManagerImpl) InitSaleLabels() error {
	if len(l.GetAllSaleLabels()) != 0 {
		return nil
	}

	arr := []item.Label{
		{
			TagName: "新品上架",
			TagCode: "new-goods",
		},
		{
			TagName: "热销商品",
			TagCode: "hot-sales",
		},
		{
			TagName: "特色商品",
			TagCode: "special-goods",
		},
		{
			TagName: "优惠促销",
			TagCode: "prom-sales",
		},
		{
			TagName: "尾品清仓",
			TagCode: "clean-goods",
		},
	}

	var err error
	for _, v := range arr {
		v.Enabled = 1
		v.MerchantId = l._mchId
		_, err = l.CreateSaleLabel(&v).Save()
	}

	return err
}

// 获取所有的销售标签
func (l *labelManagerImpl) GetAllSaleLabels() []item.ISaleLabel {
	arr := l._rep.GetAllValueSaleLabels(l._mchId)
	var tags = make([]item.ISaleLabel, len(arr))

	for i, v := range arr {
		tags[i] = l.CreateSaleLabel(v)
	}
	return tags
}

// 获取销售标签
func (l *labelManagerImpl) GetSaleLabel(id int32) item.ISaleLabel {
	return l._rep.GetSaleLabel(l._mchId, id)
}

// 根据Code获取销售标签
func (l *labelManagerImpl) GetSaleLabelByCode(code string) item.ISaleLabel {
	v := l._rep.GetSaleLabelByCode(l._mchId, code)
	return l.CreateSaleLabel(v)
}

// 创建销售标签
func (l *labelManagerImpl) CreateSaleLabel(v *item.Label) item.ISaleLabel {
	if v == nil {
		return nil
	}
	v.MerchantId = l._mchId
	return l._rep.CreateSaleLabel(v)
}

// 删除销售标签
func (l *labelManagerImpl) DeleteSaleLabel(id int32) error {
	v := l.GetSaleLabel(id)
	if v != nil {
		if v.System() {
			return item.ErrInternalDisallow
		}
		return l._rep.DeleteSaleLabel(l._mchId, id)
	}
	return nil
}
