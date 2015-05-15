/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-12-08 11:44
 * description :
 * history :
 */

package sale

import (
	"go2o/src/core/domain/interface/sale"
	"time"
)

var _ sale.ISale = new(Sale)

const MAX_CACHE_SIZE int = 1000

type Sale struct {
	partnerId  int
	saleRep    sale.ISaleRep
	proCache   map[int]sale.IGoods
	categories []sale.ICategory
}

func NewSale(partnerId int, saleRep sale.ISaleRep) sale.ISale {
	return (&Sale{
		partnerId: partnerId,
		saleRep:   saleRep,
	}).init()
}

func (this *Sale) init() sale.ISale {
	this.proCache = make(map[int]sale.IGoods)
	return this
}

func (this *Sale) clearCache(goodsId int) {
	delete(this.proCache, goodsId)
}

func (this *Sale) chkCache() {
	if len(this.proCache) >= MAX_CACHE_SIZE {
		this.proCache = make(map[int]sale.IGoods)
	}
}

func (this *Sale) GetAggregateRootId() int {
	return this.partnerId
}

func (this *Sale) CreateGoods(v *sale.ValueGoods) sale.IGoods {
	if v.CreateTime == 0 {
		v.CreateTime = time.Now().Unix()
	}

	if v.UpdateTime == 0 {
		v.UpdateTime = v.CreateTime
	}

	//todo: 判断category

	return newGoods(this, v, this.saleRep)
}

// 根据产品编号获取产品
func (this *Sale) GetGoods(goodsId int) sale.IGoods {
	p, ok := this.proCache[goodsId]
	if !ok {
		this.chkCache()
		pv := this.saleRep.GetValueGoods(this.GetAggregateRootId(), goodsId)

		if pv != nil {
			p = this.CreateGoods(pv)
			this.proCache[goodsId] = p
		}
	}
	return p
}

// 删除商品
func (this *Sale) DeleteGoods(goodsId int) error {
	err := this.saleRep.DeleteGoods(this.GetAggregateRootId(), goodsId)
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

	return newCategory(this.saleRep, v)
}

// 获取分类
func (this *Sale) GetCategory(id int) sale.ICategory {
	v := this.saleRep.GetCategory(this.GetAggregateRootId(), id)
	if v != nil {
		return this.CreateCategory(v)
	}
	return nil
}

// 获取所有分类
func (this *Sale) GetCategories() []sale.ICategory {
	//if this.categories == nil {
	list := this.saleRep.GetCategories(this.GetAggregateRootId())
	this.categories = make([]sale.ICategory, len(list))
	for i, v := range list {
		this.categories[i] = this.CreateCategory(v)
	}
	//}
	return this.categories
}

// 删除分类
func (this *Sale) DeleteCategory(id int) error {
	//todo: 删除应放到这里来处理
	return this.saleRep.DeleteCategory(this.GetAggregateRootId(), id)
}

// 获取指定的商品快照
func (this *Sale) GetGoodsSnapshot(id int) *sale.GoodsSnapshot {
	return this.saleRep.GetGoodsSnapshot(id)
}

// 根据Key获取商品快照
func (this *Sale) GetGoodsSnapshotByKey(key string) *sale.GoodsSnapshot {
	return this.saleRep.GetGoodsSnapshotByKey(key)
}
