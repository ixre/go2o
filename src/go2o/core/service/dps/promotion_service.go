/**
 * Copyright 2014 @ Ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-03 23:18
 * description :
 * history :
 */

package dps

import (
	"errors"
	"go2o/core/domain/interface/promotion"
)

type promotionService struct {
	_rep promotion.IPromotionRep
}

func NewPromotionService(r promotion.IPromotionRep) *promotionService {
	return &promotionService{
		_rep: r,
	}
}

func (this *promotionService) GetCoupon(partnerId int, id int) promotion.ICoupon {
	_prom := this._rep.GetPromotion(partnerId)
	return _prom.GetCoupon(id)
}

func (this *promotionService) SaveCoupon(partnerId int, e *promotion.ValueCoupon) (int, error) {
	prom := this._rep.GetPromotion(partnerId)
	var coupon promotion.ICoupon
	if e.Id > 0 {
		coupon = prom.GetCoupon(e.Id)
		if coupon == nil {
			return 0, errors.New("优惠券不存在")
		}
		err := coupon.SetValue(e)
		if err != nil {
			return 0, err
		}
	} else {
		coupon = prom.CreateCoupon(e)
	}
	return coupon.Save()
}

func (this *promotionService) BindCoupons(partnerId int, id int, members []string) error {
	coupon := this.GetCoupon(partnerId, id)
	return coupon.Binds(members)
}
