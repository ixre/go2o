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
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/domain/interface/valueobject"
	saleImpl "go2o/core/domain/sale"
)

var _ sale.ISaleRep = new(saleRep)

type saleRep struct {
	db.Connector
	_cache      map[int]sale.ISale
	_tagRep     sale.ISaleLabelRep
	_promRep    promotion.IPromotionRep
	_itemRep    item.IItemRep
	_goodsRep   goods.IGoodsRep
	_cateRep    sale.ICategoryRep
	_expressRep express.IExpressRep
	_valRep     valueobject.IValueRep
}

func NewSaleRep(c db.Connector, cateRep sale.ICategoryRep,
	valRep valueobject.IValueRep, saleLabelRep sale.ISaleLabelRep,
	itemRep item.IItemRep, expressRep express.IExpressRep,
	goodsRep goods.IGoodsRep, promRep promotion.IPromotionRep) sale.ISaleRep {
	return (&saleRep{
		Connector:   c,
		_tagRep:     saleLabelRep,
		_promRep:    promRep,
		_itemRep:    itemRep,
		_goodsRep:   goodsRep,
		_cateRep:    cateRep,
		_expressRep: expressRep,
		_valRep:     valRep,
	}).init()
}

func (s *saleRep) init() sale.ISaleRep {
	s._cache = make(map[int]sale.ISale)
	return s
}

func (s *saleRep) GetSale(mchId int) sale.ISale {
	v, ok := s._cache[mchId]
	if !ok {
		v = saleImpl.NewSale(mchId, s, s._valRep, s._cateRep,
			s._itemRep, s._goodsRep, s._tagRep, s._expressRep,
			s._promRep)
		s._cache[mchId] = v
	}
	return v
}
