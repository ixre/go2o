/**
 * Copyright 2015 @ S1N1 Team.
 * name : factor
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package promotion

import (
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/domain/interface/sale"
	"time"
)

func FactoryPromotion(rep promotion.IPromotionRep, saleRep sale.ISaleRep,
	v *promotion.ValuePromotion, dv interface{}) promotion.IPromotion {
	p := newPromotion(rep, saleRep, v)

	switch p.Type() {
	case promotion.TypeFlagCashBack:
		return createCashBackPromotion(p, v)
	case promotion.TypeFlagCoupon:
		return createCouponPromotion(p, v)
	}

	//todo: other promotion
	return p
}

// 创建
func createCashBackPromotion(p *Promotion, v interface{}) promotion.IPromotion {
	var pv *promotion.ValueCashBack

	if v != nil {
		pv, _ = v.(*promotion.ValueCashBack)
	}

	if pv == nil {
		pv = p._promRep.GetValueCashBack(p.GetAggregateRootId())
		if pv == nil {
			pv = &promotion.ValueCashBack{
				Id: p.GetAggregateRootId(),
			}
		}
	}

	return &CashBackPromotion{
		Promotion:      p,
		_cashBackValue: pv,
	}
}

func createCouponPromotion(p *Promotion, v interface{}) promotion.IPromotion {
	var pv *promotion.ValueCoupon

	if v != nil {
		pv, _ = v.(*promotion.ValueCoupon)
	}

	if pv == nil {
		pv = p._promRep.GetValueCoupon(p.GetAggregateRootId())
		if pv == nil {
			pv = &promotion.ValueCoupon{
				Id:         p.GetAggregateRootId(),
				CreateTime: time.Now().Unix(),
			}
		}
	}

	pv.Amount = pv.TotalAmount

	return newCoupon(p, pv, p.promRep, p.memberRep)
}

func DeletePromotion(p promotion.IPromotion) error {
	var err error
	var pi = p.(*Promotion)
	if p.Type() == promotion.TypeFlagCashBack {
		v := p.(*CashBackPromotion)
		err = v._promRep.DeleteValueCashBack(v.GetAggregateRootId())
	}
	if err == nil {
		err = pi._promRep.DeletePromotion(p.GetAggregateRootId())
	}
	return err
}
