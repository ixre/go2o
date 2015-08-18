/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:18
 * description :
 * history :
 */

package dps

import (
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/promotion"
	promImpl "go2o/src/core/domain/promotion"
)

type promotionService struct {
	_rep promotion.IPromotionRep
}

func NewPromotionService(rep promotion.IPromotionRep) *promotionService {
	return &promotionService{
		_rep: rep,
	}
}

// 获取促销
func (this *promotionService) GetPromotion(id int) (*promotion.ValuePromotion, interface{}) {
	var prom promotion.IPromotion = this._rep.GetPromotion(id)
	if prom != nil {
		return prom.GetValue(), prom.GetRelationValue()
	}
	return nil, nil
}

// 保存促销
func (this *promotionService) SavePromotion(v *promotion.ValuePromotion) (int, error) {
	var prom promotion.IPromotion
	if v.Id > 0 {
		prom = this._rep.GetPromotion(v.Id)
		err := prom.SetValue(v)
		if err != nil {
			return v.Id, err
		}
	} else {
		prom = this._rep.CreatePromotion(v)
	}
	return prom.Save()
}

// 删除促销
func (this *promotionService) DelPromotion(partnerId int, promId int) error {
	prom := this._rep.GetPromotion(promId)
	if prom == nil {
		return promotion.ErrNoSuchPromotion
	}
	if prom.GetValue().PartnerId != partnerId {
		return partner.ErrPartnerNotMatch
	}

	return promImpl.DeletePromotion(prom)
}

func (this *promotionService) SaveCashBackPromotion(partnerId int, v *promotion.ValuePromotion,
	v1 *promotion.ValueCashBack) (int, error) {
	var prom promotion.IPromotion
	var err error
	if v.Id > 0 {
		prom = this._rep.GetPromotion(v.Id)
		if prom.GetValue().PartnerId != partnerId {
			return -1, partner.ErrPartnerNotMatch
		}
	} else {
		prom = this._rep.CreatePromotion(v)
	}

	if err = prom.SetValue(v); err == nil {
		cb := prom.(promotion.ICashBackPromotion)
		err = cb.SetDetailsValue(v1)
	}

	if err != nil {
		return v.Id, err
	}

	return prom.Save()
}

/**************   Coupon ************/
func (this *promotionService) SaveCoupon(partnerId int, v *promotion.ValuePromotion, v1 *promotion.ValueCoupon) (int, error) {
	var prom promotion.IPromotion
	var err error
	if v.Id > 0 {
		prom = this._rep.GetPromotion(v.Id)
		if prom.GetValue().PartnerId != partnerId {
			return -1, partner.ErrPartnerNotMatch
		}
	} else {
		prom = this._rep.CreatePromotion(v)
	}

	if err = prom.SetValue(v); err == nil {
		cb := prom.(promotion.ICouponPromotion)
		err = cb.SetDetailsValue(v1)
	}

	if err != nil {
		return v.Id, err
	}

	return prom.Save()
}

func (this *promotionService) BindCoupons(partnerId int, id int, members []string) error {
	coupon := this._rep.GetPromotion(id).(promotion.ICouponPromotion)
	return coupon.Binds(members)
}
