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
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/domain/interface/valueobject"
)

var _ sale.ISale = new(saleImpl)

const MAX_CACHE_SIZE int = 5000

type saleImpl struct {
	_mchId        int
	_saleRep      sale.ISaleRep
	_labelRep     sale.ISaleLabelRep
	_cateRep      sale.ICategoryRep
	_goodsRep     goods.IGoodsRep
	_valRep       valueobject.IValueRep
	_expressRep   express.IExpressRep
	_promRep      promotion.IPromotionRep
	_proCache     map[int]sale.IItem
	_cateManager  sale.ICategoryManager
	_labelManager sale.ILabelManager
	_itemManager  sale.IItemManager
	_itemRep      item.IItemRep
	_goodsManager sale.IGoodsManager
}

func NewSale(merchantId int, saleRep sale.ISaleRep, valRep valueobject.IValueRep,
	cateRep sale.ICategoryRep, itemRep item.IItemRep, goodsRep goods.IGoodsRep,
	tagRep sale.ISaleLabelRep, expressRep express.IExpressRep,
	promRep promotion.IPromotionRep) sale.ISale {
	return (&saleImpl{
		_mchId:      merchantId,
		_cateRep:    cateRep,
		_saleRep:    saleRep,
		_labelRep:   tagRep,
		_itemRep:    itemRep,
		_goodsRep:   goodsRep,
		_expressRep: expressRep,
		_promRep:    promRep,
		_valRep:     valRep,
	}).init()
}

func (s *saleImpl) init() sale.ISale {
	s._proCache = make(map[int]sale.IItem)
	return s
}

// 分类服务
func (s *saleImpl) CategoryManager() sale.ICategoryManager {
	if s._cateManager == nil {
		s._cateManager = NewCategoryManager(
			s.GetAggregateRootId(), s._cateRep, s._valRep)
	}
	return s._cateManager
}

// 标签管理器
func (s *saleImpl) LabelManager() sale.ILabelManager {
	if s._labelManager == nil {
		s._labelManager = NewLabelManager(
			s.GetAggregateRootId(), s._labelRep, s._valRep)
	}
	return s._labelManager
}

// 货品服务
func (s *saleImpl) ItemManager() sale.IItemManager {
	if s._itemManager == nil {
		s._itemManager = NewItemManager(
			s.GetAggregateRootId(), s, s._itemRep,
			s._expressRep, s._valRep)
	}
	return s._itemManager
}

// 商品服务
func (s *saleImpl) GoodsManager() sale.IGoodsManager {
	if s._goodsManager == nil {
		s._goodsManager = NewGoodsManager(
			s.GetAggregateRootId(), s, s._valRep)
	}
	return s._goodsManager
}

func (s *saleImpl) clearCache(goodsId int) {
	delete(s._proCache, goodsId)
}

func (s *saleImpl) chkCache() {
	if len(s._proCache) >= MAX_CACHE_SIZE {
		s._proCache = make(map[int]sale.IItem)
	}
}

func (s *saleImpl) GetAggregateRootId() int {
	return s._mchId
}
