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
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/valueobject"
)

var _ sale.ISale = new(SaleImpl)

const MAX_CACHE_SIZE int = 5000

type SaleImpl struct {
	_mchId        int
	_saleRep      sale.ISaleRep
	_labelRep     sale.ISaleLabelRep
	_cateRep      sale.ICategoryRep
	_goodsRep     goods.IGoodsRep
	_valRep       valueobject.IValueRep
	_promRep      promotion.IPromotionRep
	_proCache     map[int]sale.IItem
	_cateManager  sale.ICategoryManager
	_labelManager sale.ILabelManager
	_itemManager  sale.IItemManager
	_goodsManager sale.IGoodsManager
}

func NewSale(merchantId int, saleRep sale.ISaleRep, valRep valueobject.IValueRep,
	cateRep sale.ICategoryRep, goodsRep goods.IGoodsRep, tagRep sale.ISaleLabelRep,
	promRep promotion.IPromotionRep) sale.ISale {
	return (&SaleImpl{
		_mchId:    merchantId,
		_cateRep:  cateRep,
		_saleRep:  saleRep,
		_labelRep: tagRep,
		_goodsRep: goodsRep,
		_promRep:  promRep,
		_valRep:   valRep,
	}).init()
}

func (this *SaleImpl) init() sale.ISale {
	this._proCache = make(map[int]sale.IItem)
	return this
}

// 分类服务
func (this *SaleImpl) CategoryManager() sale.ICategoryManager {
	if this._cateManager == nil {
		this._cateManager = NewCategoryManager(
			this.GetAggregateRootId(), this._cateRep, this._valRep)
	}
	return this._cateManager
}

// 标签管理器
func (this *SaleImpl) LabelManager() sale.ILabelManager {
	if this._labelManager == nil {
		this._labelManager = NewLabelManager(
			this.GetAggregateRootId(), this._labelRep, this._valRep)
	}
	return this._labelManager
}

// 货品服务
func (this *SaleImpl) ItemManager() sale.IItemManager {
	if this._itemManager == nil {
		this._itemManager = NewItemManager(
			this.GetAggregateRootId(), this, this._valRep)
	}
	return this._itemManager
}

// 商品服务
func (this *SaleImpl) GoodsManager() sale.IGoodsManager {
	if this._goodsManager == nil {
		this._goodsManager = NewGoodsManager(
			this.GetAggregateRootId(), this, this._valRep)
	}
	return this._goodsManager
}

func (this *SaleImpl) clearCache(goodsId int) {
	delete(this._proCache, goodsId)
}

func (this *SaleImpl) chkCache() {
	if len(this._proCache) >= MAX_CACHE_SIZE {
		this._proCache = make(map[int]sale.IItem)
	}
}

func (this *SaleImpl) GetAggregateRootId() int {
	return this._mchId
}
