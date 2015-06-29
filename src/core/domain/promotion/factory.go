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
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/domain/interface/sale"
	"time"
)

func FactoryPromotion(rep promotion.IPromotionRep, saleRep sale.ISaleRep, memRep member.IMemberRep,
	v *promotion.ValuePromotion) promotion.IPromotion {
	p := newPromotion(rep, saleRep, memRep, v)

	switch p.Type() {
	case promotion.TypeFlagCashBack:
		return createCashBackPromotion(p)
	case promotion.TypeFlagCoupon:
		return createCouponPromotion(p)
	}
	//todo: other promotion
	return p
}

// 创建
func createCashBackPromotion(p *Promotion) promotion.IPromotion {
	var pv *promotion.ValueCashBack

	if p.GetAggregateRootId() > 0 {
		pv = p._promRep.GetValueCashBack(p.GetAggregateRootId())
	}
	if pv == nil {
		pv = &promotion.ValueCashBack{
			Id: p.GetAggregateRootId(),
		}
	}

	return &CashBackPromotion{
		Promotion:      p,
		_cashBackValue: pv,
	}
}

func createCouponPromotion(p *Promotion) promotion.IPromotion {
	var pv *promotion.ValueCoupon

	if p.GetAggregateRootId() > 0 {
		pv = p._promRep.GetValueCoupon(p.GetAggregateRootId())
	}
	if pv == nil {
		pv = &promotion.ValueCoupon{
			Id:         p.GetAggregateRootId(),
			CreateTime: time.Now().Unix(),
		}

	}

	pv.Amount = pv.TotalAmount

	return newCoupon(p, pv, p._promRep, p._memberRep)
}

func DeletePromotion(p promotion.IPromotion) error {
	var err error
	var rep promotion.IPromotionRep = nil
	if p.Type() == promotion.TypeFlagCashBack {
		v := p.(*CashBackPromotion)
		rep = v._promRep
		err = rep.DeleteValueCashBack(v.GetAggregateRootId())
	} else if p.Type() == promotion.TypeFlagCoupon {
		v := p.(*Coupon)
		rep = v._promRep
		err = rep.DeleteValueCoupon(v.GetAggregateRootId())
	}
	if err == nil && rep != nil {
		err = rep.DeletePromotion(p.GetAggregateRootId())
	}
	return err
}
