/**
 * Copyright 2014 @ Ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-03 23:18
 * description :
 * history :
 */

package dproxy

import (
	"com/domain/interface/promotion"
)

type promotionService struct {
	promRep promotion.IPromotionRep
}

func (this *promotionService) GetCoupon(partnerId int, id int) promotion.ICoupon {
	_prom := this.promRep.GetPromotion(partnerId)
	return _prom.GetCoupon(id)
}

func (this *promotionService) SaveCoupon(partnerId int, e *promotion.ValueCoupon) (int, error) {
	_prom := this.promRep.GetPromotion(partnerId)
	return _prom.CreateCoupon(e).Save()
}

func (this *promotionService) BindCoupons(partnerId int, id int, members []string) error {
	coupon := this.GetCoupon(partnerId, id)
	return coupon.Binds(members)
}
