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
	"go2o/core/domain/interface/sale/product"
	"go2o/core/domain/interface/valueobject"
)

var _ sale.ISale = new(saleImpl)

type saleImpl struct {
	mchId        int32
	saleRepo     sale.ISaleRepo
	labelRepo    sale.ISaleLabelRepo
	cateRepo     sale.ICategoryRepo
	goodsRepo    goods.IGoodsRepo
	valRepo      valueobject.IValueRepo
	expressRepo  express.IExpressRepo
	promRepo     promotion.IPromotionRepo
	cateManager  sale.ICategoryManager
	labelManager sale.ILabelManager
	itemManager  sale.IItemManager
	itemRepo     product.IProductRepo
	goodsManager sale.IGoodsManager
}

func NewSale(mchId int32, saleRepo sale.ISaleRepo, valRepo valueobject.IValueRepo,
	cateRepo sale.ICategoryRepo, itemRepo product.IProductRepo, goodsRepo goods.IGoodsRepo,
	tagRepo sale.ISaleLabelRepo, expressRepo express.IExpressRepo,
	promRepo promotion.IPromotionRepo) sale.ISale {
	return (&saleImpl{
		mchId:       mchId,
		cateRepo:    cateRepo,
		saleRepo:    saleRepo,
		labelRepo:   tagRepo,
		itemRepo:    itemRepo,
		goodsRepo:   goodsRepo,
		expressRepo: expressRepo,
		promRepo:    promRepo,
		valRepo:     valRepo,
	}).init()
}

func (s *saleImpl) init() sale.ISale {
	return s
}

// 分类服务
func (s *saleImpl) CategoryManager() sale.ICategoryManager {
	if s.cateManager == nil {
		s.cateManager = NewCategoryManager(
			s.GetAggregateRootId(), s.cateRepo, s.valRepo)
	}
	return s.cateManager
}

// 标签管理器
func (s *saleImpl) LabelManager() sale.ILabelManager {
	if s.labelManager == nil {
		s.labelManager = NewLabelManager(
			s.GetAggregateRootId(), s.labelRepo, s.valRepo)
	}
	return s.labelManager
}

// 货品服务
func (s *saleImpl) ItemManager() sale.IItemManager {
	if s.itemManager == nil {
		s.itemManager = NewItemManager(
			s.GetAggregateRootId(), s, s.itemRepo,
			s.expressRepo, s.valRepo)
	}
	return s.itemManager
}

// 商品服务
func (s *saleImpl) GoodsManager() sale.IGoodsManager {
	if s.goodsManager == nil {
		s.goodsManager = NewGoodsManager(
			s.GetAggregateRootId(), s, s.valRepo)
	}
	return s.goodsManager
}

func (s *saleImpl) GetAggregateRootId() int32 {
	return s.mchId
}
