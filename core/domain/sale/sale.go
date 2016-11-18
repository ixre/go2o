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

type saleImpl struct {
	mchId        int64
	saleRep      sale.ISaleRep
	labelRep     sale.ISaleLabelRep
	cateRep      sale.ICategoryRep
	goodsRep     goods.IGoodsRep
	valRep       valueobject.IValueRep
	expressRep   express.IExpressRep
	promRep      promotion.IPromotionRep
	cateManager  sale.ICategoryManager
	labelManager sale.ILabelManager
	itemManager  sale.IItemManager
	itemRep      item.IItemRep
	goodsManager sale.IGoodsManager
}

func NewSale(mchId int64, saleRep sale.ISaleRep, valRep valueobject.IValueRep,
	cateRep sale.ICategoryRep, itemRep item.IItemRep, goodsRep goods.IGoodsRep,
	tagRep sale.ISaleLabelRep, expressRep express.IExpressRep,
	promRep promotion.IPromotionRep) sale.ISale {
	return (&saleImpl{
		mchId:      mchId,
		cateRep:    cateRep,
		saleRep:    saleRep,
		labelRep:   tagRep,
		itemRep:    itemRep,
		goodsRep:   goodsRep,
		expressRep: expressRep,
		promRep:    promRep,
		valRep:     valRep,
	}).init()
}

func (s *saleImpl) init() sale.ISale {
	return s
}

// 分类服务
func (s *saleImpl) CategoryManager() sale.ICategoryManager {
	if s.cateManager == nil {
		s.cateManager = NewCategoryManager(
			s.GetAggregateRootId(), s.cateRep, s.valRep)
	}
	return s.cateManager
}

// 标签管理器
func (s *saleImpl) LabelManager() sale.ILabelManager {
	if s.labelManager == nil {
		s.labelManager = NewLabelManager(
			s.GetAggregateRootId(), s.labelRep, s.valRep)
	}
	return s.labelManager
}

// 货品服务
func (s *saleImpl) ItemManager() sale.IItemManager {
	if s.itemManager == nil {
		s.itemManager = NewItemManager(
			s.GetAggregateRootId(), s, s.itemRep,
			s.expressRep, s.valRep)
	}
	return s.itemManager
}

// 商品服务
func (s *saleImpl) GoodsManager() sale.IGoodsManager {
	if s.goodsManager == nil {
		s.goodsManager = NewGoodsManager(
			s.GetAggregateRootId(), s, s.valRep)
	}
	return s.goodsManager
}

func (s *saleImpl) GetAggregateRootId() int64 {
	return s.mchId
}
