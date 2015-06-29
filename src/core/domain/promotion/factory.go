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
	v *promotion.ValuePromotion) promotion.IPromotion {
	p := newPromotion(rep, saleRep, v)

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
	pv := p._promRep.GetValueCashBack(p.GetAggregateRootId())
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
	//	val.CreateTime = time.Now().Unix()
	//	val.Amount = val.TotalAmount
	//	return newCoupon(val, this.promRep, this.memberRep)

	pv := p.promRep.GetValueCoupon(p.GetAggregateRootId())
	if pv == nil {
		pv = &promotion.ValueCoupon{
			Id: p.GetAggregateRootId(),
			Fee : pv.TotalAmount,
			CreateTime:time.Now().Unix(),
		}
	}
	return newCoupon(p,pv,p.promRep,p.memberRep)
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
