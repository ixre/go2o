/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-08 11:44
 * description :
 * history :
 */

package sale

import (
	"errors"
	"go2o/src/core/domain/interface/sale"
	"time"
)

var _ sale.ISale = new(Sale)

const MAX_CACHE_SIZE int = 1000

type Sale struct {
	_partnerId  int
	_saleRep    sale.ISaleRep
	_saleTagRep sale.ISaleTagRep
	_proCache   map[int]sale.IGoods
	_categories []sale.ICategory
}

func NewSale(partnerId int, saleRep sale.ISaleRep, tagRep sale.ISaleTagRep) sale.ISale {
	return (&Sale{
		_partnerId:  partnerId,
		_saleRep:    saleRep,
		_saleTagRep: tagRep,
	}).init()
}

func (this *Sale) init() sale.ISale {
	this._proCache = make(map[int]sale.IGoods)
	return this
}

func (this *Sale) clearCache(goodsId int) {
	delete(this._proCache, goodsId)
}

func (this *Sale) chkCache() {
	if len(this._proCache) >= MAX_CACHE_SIZE {
		this._proCache = make(map[int]sale.IGoods)
	}
}

func (this *Sale) GetAggregateRootId() int {
	return this._partnerId
}

func (this *Sale) CreateGoods(v *sale.ValueGoods) sale.IGoods {
	if v.CreateTime == 0 {
		v.CreateTime = time.Now().Unix()
	}

	if v.UpdateTime == 0 {
		v.UpdateTime = v.CreateTime
	}

	//todo: 判断category

	return newGoods(this, v, this._saleRep, this._saleTagRep)
}

// 根据产品编号获取产品
func (this *Sale) GetGoods(goodsId int) sale.IGoods {
	p, ok := this._proCache[goodsId]
	if !ok {
		this.chkCache()
		pv := this._saleRep.GetValueGoods(this.GetAggregateRootId(), goodsId)

		if pv != nil {
			p = this.CreateGoods(pv)
			this._proCache[goodsId] = p
		}
	}
	return p
}

// 删除商品
func (this *Sale) DeleteGoods(goodsId int) error {
	err := this._saleRep.DeleteGoods(this.GetAggregateRootId(), goodsId)
	if err != nil {
		this.clearCache(goodsId)
	}
	return err
}

// 创建分类
func (this *Sale) CreateCategory(v *sale.ValueCategory) sale.ICategory {
	if v.CreateTime == 0 {
		v.CreateTime = time.Now().Unix()
	}
	v.PartnerId = this.GetAggregateRootId()

	return newCategory(this._saleRep, v)
}

// 获取分类
func (this *Sale) GetCategory(id int) sale.ICategory {
	v := this._saleRep.GetCategory(this.GetAggregateRootId(), id)
	if v != nil {
		return this.CreateCategory(v)
	}
	return nil
}

// 获取所有分类
func (this *Sale) GetCategories() []sale.ICategory {
	//if this.categories == nil {
	list := this._saleRep.GetCategories(this.GetAggregateRootId())
	this._categories = make([]sale.ICategory, len(list))
	for i, v := range list {
		this._categories[i] = this.CreateCategory(v)
	}
	//}
	return this._categories
}

// 删除分类
func (this *Sale) DeleteCategory(id int) error {
	//todo: 删除应放到这里来处理
	return this._saleRep.DeleteCategory(this.GetAggregateRootId(), id)
}

// 初始化销售标签
func (this *Sale) InitSaleTags() error {
	if len(this.GetAllSaleTags()) != 0 {
		return errors.New("已经存在数据，无法初始化!")
	}

	arr := []sale.ValueSaleTag{
		sale.ValueSaleTag{
			TagName: "新品上架",
			TagCode: "new-goods",
		},
		sale.ValueSaleTag{
			TagName: "热销商品",
			TagCode: "hot-sales",
		},
		sale.ValueSaleTag{
			TagName: "特色商品",
			TagCode: "special-goods",
		},
		sale.ValueSaleTag{
			TagName: "优惠促销",
			TagCode: "prom-sales",
		},
		sale.ValueSaleTag{
			TagName: "尾品清仓",
			TagCode: "clean-goods",
		},
	}

	var err error
	for _, v := range arr {
		v.Enabled = 1
		v.PartnerId = this._partnerId
		_, err = this.CreateSaleTag(&v).Save()
	}

	return err
}

// 获取所有的销售标签
func (this *Sale) GetAllSaleTags() []sale.ISaleTag {
	arr := this._saleTagRep.GetAllValueSaleTags(this._partnerId)
	var tags = make([]sale.ISaleTag, len(arr))

	for i, v := range arr {
		tags[i] = this.CreateSaleTag(v)
	}
	return tags
}

// 获取销售标签
func (this *Sale) GetSaleTag(id int) sale.ISaleTag {
	return this._saleTagRep.GetSaleTag(this._partnerId, id)
}

// 根据Code获取销售标签
func (this *Sale) GetSaleTagByCode(code string) sale.ISaleTag {
	v := this._saleTagRep.GetSaleTagByCode(this._partnerId, code)
	return this.CreateSaleTag(v)
}

// 创建销售标签
func (this *Sale) CreateSaleTag(v *sale.ValueSaleTag) sale.ISaleTag {
	if v == nil {
		return nil
	}
	return this._saleTagRep.CreateSaleTag(v)
}

// 删除销售标签
func (this *Sale) DeleteSaleTag(id int) error {
	return this._saleTagRep.DeleteSaleTag(this._partnerId, id)
}

// 获取指定的商品快照
func (this *Sale) GetGoodsSnapshot(id int) *sale.GoodsSnapshot {
	return this._saleRep.GetGoodsSnapshot(id)
}

// 根据Key获取商品快照
func (this *Sale) GetGoodsSnapshotByKey(key string) *sale.GoodsSnapshot {
	return this._saleRep.GetGoodsSnapshotByKey(key)
}
