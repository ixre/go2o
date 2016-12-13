/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 11:09
 * description :
 * history :
 */

package repository

import (
	"github.com/jsix/gof/db"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/valueobject"
	saleImpl "go2o/core/domain/sale"
)

var _ sale.ISaleRepo = new(saleRepo)

type saleRepo struct {
	db.Connector
	_cache       map[int32]sale.ISale
	_tagRepo     item.ISaleLabelRepo
	_promRepo    promotion.IPromotionRepo
	_itemRepo    product.IProductRepo
	_goodsRepo   item.IGoodsRepo
	_cateRepo    product.ICategoryRepo
	_expressRepo express.IExpressRepo
	_valRepo     valueobject.IValueRepo
}

func NewSaleRepo(c db.Connector, cateRepo product.ICategoryRepo,
	valRepo valueobject.IValueRepo, saleLabelRepo item.ISaleLabelRepo,
	itemRepo product.IProductRepo, expressRepo express.IExpressRepo,
	goodsRepo item.IGoodsRepo, promRepo promotion.IPromotionRepo) sale.ISaleRepo {
	return (&saleRepo{
		Connector:    c,
		_tagRepo:     saleLabelRepo,
		_promRepo:    promRepo,
		_itemRepo:    itemRepo,
		_goodsRepo:   goodsRepo,
		_cateRepo:    cateRepo,
		_expressRepo: expressRepo,
		_valRepo:     valRepo,
	}).init()
}

func (s *saleRepo) init() sale.ISaleRepo {
	s._cache = make(map[int32]sale.ISale)
	return s
}

func (s *saleRepo) GetSale(mchId int32) sale.ISale {
	v, ok := s._cache[mchId]
	if !ok {
		v = saleImpl.NewSale(mchId, s, s._valRepo, s._cateRepo,
			s._itemRepo, s._goodsRepo, s._tagRepo, s._expressRepo,
			s._promRepo)
		s._cache[mchId] = v
	}
	return v
}
