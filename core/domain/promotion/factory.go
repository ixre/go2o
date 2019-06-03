/**
 * Copyright 2015 @ to2.net.
 * name : factor
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package promotion

import (
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/promotion"
	"time"
)

func FactoryPromotion(rep promotion.IPromotionRepo, goodsRepo item.IGoodsItemRepo, memRepo member.IMemberRepo,
	v *promotion.PromotionInfo) promotion.IPromotion {
	p := newPromotion(rep, goodsRepo, memRepo, v)

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
func createCashBackPromotion(p *promotionImpl) promotion.IPromotion {
	var pv *promotion.ValueCashBack

	if p.GetAggregateRootId() > 0 {
		pv = p.promRepo.GetValueCashBack(p.GetAggregateRootId())
	}
	if pv == nil {
		pv = &promotion.ValueCashBack{
			Id: p.GetAggregateRootId(),
		}
	}

	return &CashBackPromotion{
		promotionImpl: p,
		cashBackValue: pv,
	}
}

func createCouponPromotion(p *promotionImpl) promotion.IPromotion {
	var pv *promotion.ValueCoupon

	if p.GetAggregateRootId() > 0 {
		pv = p.promRepo.GetValueCoupon(p.GetAggregateRootId())
	}
	if pv == nil {
		pv = &promotion.ValueCoupon{
			Id:         p.GetAggregateRootId(),
			CreateTime: time.Now().Unix(),
		}

	}

	if p.GetAggregateRootId() <= 0 {
		pv.Amount = pv.TotalAmount
	}

	return newCoupon(p, pv, p.promRepo, p.memberRepo)
}

func DeletePromotion(p promotion.IPromotion) error {
	var err error
	var rep promotion.IPromotionRepo = nil
	if p.Type() == promotion.TypeFlagCashBack {
		v := p.(*CashBackPromotion)
		rep = v.promRepo
		err = rep.DeleteValueCashBack(v.GetAggregateRootId())
	} else if p.Type() == promotion.TypeFlagCoupon {
		v := p.(*Coupon)
		rep = v.promRepo
		err = rep.DeleteValueCoupon(v.GetAggregateRootId())
	}
	if err == nil && rep != nil {
		err = rep.DeletePromotion(p.GetAggregateRootId())
	}
	return err
}
