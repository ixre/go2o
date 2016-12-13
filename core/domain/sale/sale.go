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
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/valueobject"
	itemImpl "go2o/core/domain/item"
	productImpl "go2o/core/domain/product"
)

var _ sale.ISale = new(saleImpl)

type saleImpl struct {
	mchId        int32
	saleRepo     sale.ISaleRepo
	labelRepo    item.ISaleLabelRepo
	cateRepo     product.ICategoryRepo
	goodsRepo    item.IGoodsRepo
	valRepo      valueobject.IValueRepo
	expressRepo  express.IExpressRepo
	promRepo     promotion.IPromotionRepo
	cateManager  product.IGlobCatService
	labelManager item.ILabelManager
	productRepo  product.IProductRepo
	goodsManager item.IGoodsManager
}

func NewSale(mchId int32, saleRepo sale.ISaleRepo, valRepo valueobject.IValueRepo,
	cateRepo product.ICategoryRepo, itemRepo product.IProductRepo, goodsRepo item.IGoodsRepo,
	tagRepo item.ISaleLabelRepo, expressRepo express.IExpressRepo,
	promRepo promotion.IPromotionRepo) sale.ISale {
	return (&saleImpl{
		mchId:       mchId,
		cateRepo:    cateRepo,
		saleRepo:    saleRepo,
		labelRepo:   tagRepo,
		productRepo: itemRepo,
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
func (s *saleImpl) CategoryManager() product.IGlobCatService {
	if s.cateManager == nil {
		s.cateManager = productImpl.NewCategoryManager(
			s.GetAggregateRootId(), s.cateRepo, s.valRepo)
	}
	return s.cateManager
}

// 标签管理器
func (s *saleImpl) LabelManager() item.ILabelManager {
	if s.labelManager == nil {
		s.labelManager = itemImpl.NewLabelManager(
			s.GetAggregateRootId(), s.labelRepo, s.valRepo)
	}
	return s.labelManager
}

// 商品服务
func (s *saleImpl) GoodsManager() item.IGoodsManager {
	if s.goodsManager == nil {
		s.goodsManager = itemImpl.NewGoodsManager(
			s.GetAggregateRootId(),
			s.goodsRepo, s.productRepo, s.promRepo, s.valRepo)

	}
	return s.goodsManager
}

func (s *saleImpl) GetAggregateRootId() int32 {
	return s.mchId
}
