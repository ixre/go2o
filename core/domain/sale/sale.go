/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 11:44
 * description :
 * history :
 */

package sale

import (
	"errors"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/valueobject"
	"time"
)

var _ sale.ISale = new(Sale)

const MAX_CACHE_SIZE int = 1000

type Sale struct {
	_merchantId   int
	_saleRep      sale.ISaleRep
	_saleLabelRep sale.ISaleLabelRep
	_cateRep      sale.ICategoryRep
	_goodsRep     sale.IGoodsRep
	_valRep       valueobject.IValueRep
	_promRep      promotion.IPromotionRep
	_proCache     map[int]sale.IItem
	_cateManager  sale.ICategoryManager
}

func NewSale(merchantId int, saleRep sale.ISaleRep, valRep valueobject.IValueRep,
	cateRep sale.ICategoryRep, goodsRep sale.IGoodsRep, tagRep sale.ISaleLabelRep,
	promRep promotion.IPromotionRep) sale.ISale {
	return (&Sale{
		_merchantId:   merchantId,
		_cateRep:      cateRep,
		_saleRep:      saleRep,
		_saleLabelRep: tagRep,
		_goodsRep:     goodsRep,
		_promRep:      promRep,
		_valRep:       valRep,
	}).init()
}

func (this *Sale) init() sale.ISale {
	this._proCache = make(map[int]sale.IItem)
	return this
}

// 分类服务
func (this *Sale) CategoryManager() sale.ICategoryManager {
	if this._cateManager == nil {
		this._cateManager = NewCategoryManager(
			this.GetAggregateRootId(), this._cateRep, this._valRep)
	}
	return this._cateManager
}

func (this *Sale) clearCache(goodsId int) {
	delete(this._proCache, goodsId)
}

func (this *Sale) chkCache() {
	if len(this._proCache) >= MAX_CACHE_SIZE {
		this._proCache = make(map[int]sale.IItem)
	}
}

func (this *Sale) GetAggregateRootId() int {
	return this._merchantId
}

func (this *Sale) CreateItem(v *sale.ValueItem) sale.IItem {
	if v.CreateTime == 0 {
		v.CreateTime = time.Now().Unix()
	}
	if v.UpdateTime == 0 {
		v.UpdateTime = v.CreateTime
	} //todo: 判断category
	return newItem(this, v, this._saleRep, this._saleLabelRep, this._goodsRep, this._promRep)
}

// 创建商品
func (this *Sale) CreateGoods(s *sale.ValueGoods) sale.IGoods {
	return NewSaleGoods(this, nil, s, this._saleRep, this._goodsRep, this._promRep)
}

// 删除货品
func (this *Sale) DeleteItem(id int) error {
	var err error
	num := this._saleRep.GetItemSaleNum(this.GetAggregateRootId(), id)

	if num == 0 {
		err = this._saleRep.DeleteItem(this.GetAggregateRootId(), id)
		if err != nil {
			this.clearCache(id)
		}
	} else {
		err = sale.ErrCanNotDeleteItem
	}
	return err
}

// 根据产品编号获取产品
func (this *Sale) GetItem(itemId int) sale.IItem {
	pv := this._saleRep.GetValueItem(this.GetAggregateRootId(), itemId)
	if pv != nil {
		return this.CreateItem(pv)
	}
	return nil
}

// 创建商品
func (this *Sale) CreateGoodsByItem(item sale.IItem, v *sale.ValueGoods) sale.IGoods {
	return NewSaleGoods(this, item, v, this._saleRep, this._goodsRep, this._promRep)
}

// 根据产品编号获取商品
func (this *Sale) GetGoods(goodsId int) sale.IGoods {
	var v *sale.ValueGoods = this._goodsRep.GetValueGoodsById(goodsId)
	if v != nil {
		pv := this._saleRep.GetValueItem(this.GetAggregateRootId(), v.ItemId)
		if pv != nil {
			return this.CreateGoodsByItem(this.CreateItem(pv), v)
		}
	}
	return nil
}

// 根据产品SKU获取商品
func (this *Sale) GetGoodsBySku(itemId, sku int) sale.IGoods {
	var v *sale.ValueGoods = this._goodsRep.GetValueGoodsBySku(itemId, sku)
	if v != nil {
		pv := this._saleRep.GetValueItem(this.GetAggregateRootId(), v.ItemId)
		if pv != nil {
			return this.CreateGoodsByItem(this.CreateItem(pv), v)
		}
	}
	return nil
}

// 删除商品
func (this *Sale) DeleteGoods(goodsId int) error {
	goods := this.GetGoods(goodsId)
	if goods.GetValue().SaleNum > 0 {
		return sale.ErrNoSuchSnapshot
	}

	//todo: delete goods
	err := this._saleRep.DeleteItem(this.GetAggregateRootId(), goodsId)
	if err != nil {
		this.clearCache(goodsId)
	}
	return err
}

// 初始化销售标签
func (this *Sale) InitSaleLabels() error {
	if len(this.GetAllSaleLabels()) != 0 {
		return errors.New("已经存在数据，无法初始化!")
	}

	arr := []sale.SaleLabel{
		sale.SaleLabel{
			TagName: "新品上架",
			TagCode: "new-goods",
		},
		sale.SaleLabel{
			TagName: "热销商品",
			TagCode: "hot-sales",
		},
		sale.SaleLabel{
			TagName: "特色商品",
			TagCode: "special-goods",
		},
		sale.SaleLabel{
			TagName: "优惠促销",
			TagCode: "prom-sales",
		},
		sale.SaleLabel{
			TagName: "尾品清仓",
			TagCode: "clean-goods",
		},
	}

	var err error
	for _, v := range arr {
		v.Enabled = 1
		v.MerchantId = this._merchantId
		_, err = this.CreateSaleLabel(&v).Save()
	}

	return err
}

// 获取所有的销售标签
func (this *Sale) GetAllSaleLabels() []sale.ISaleLabel {
	arr := this._saleLabelRep.GetAllValueSaleLabels(this._merchantId)
	var tags = make([]sale.ISaleLabel, len(arr))

	for i, v := range arr {
		tags[i] = this.CreateSaleLabel(v)
	}
	return tags
}

// 获取销售标签
func (this *Sale) GetSaleLabel(id int) sale.ISaleLabel {
	return this._saleLabelRep.GetSaleLabel(this._merchantId, id)
}

// 根据Code获取销售标签
func (this *Sale) GetSaleLabelByCode(code string) sale.ISaleLabel {
	v := this._saleLabelRep.GetSaleLabelByCode(this._merchantId, code)
	return this.CreateSaleLabel(v)
}

// 创建销售标签
func (this *Sale) CreateSaleLabel(v *sale.SaleLabel) sale.ISaleLabel {
	if v == nil {
		return nil
	}
	v.MerchantId = this.GetAggregateRootId()
	return this._saleLabelRep.CreateSaleLabel(v)
}

// 删除销售标签
func (this *Sale) DeleteSaleLabel(id int) error {
	v := this.GetSaleLabel(id)
	if v != nil {
		if v.System() {
			return sale.ErrInternalDisallow
		}
		return this._saleLabelRep.DeleteSaleLabel(this._merchantId, id)
	}
	return nil
}

// 获取指定的商品快照
func (this *Sale) GetGoodsSnapshot(id int) *sale.GoodsSnapshot {
	return this._saleRep.GetGoodsSnapshot(id)
}

// 根据Key获取商品快照
func (this *Sale) GetGoodsSnapshotByKey(key string) *sale.GoodsSnapshot {
	return this._saleRep.GetGoodsSnapshotByKey(key)
}

// 获取指定数量已上架的商品
func (this *Sale) GetOnShelvesGoods(start, end int,
	sortBy string) []*valueobject.Goods {
	return this._goodsRep.GetOnShelvesGoods(this.GetAggregateRootId(),
		start, end, sortBy)
}
