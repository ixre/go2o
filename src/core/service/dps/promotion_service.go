/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:18
 * description :
 * history :
 */

package dps

import (
	"errors"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/promotion"
	promImpl "go2o/src/core/domain/promotion"
)

type promotionService struct {
	_rep    promotion.IOldPromotionRep
	_newRep promotion.IPromotionRep
}

func NewPromotionService(r promotion.IOldPromotionRep, rep promotion.IPromotionRep) *promotionService {
	return &promotionService{
		_rep:    r,
		_newRep: rep,
	}
}

// 获取促销
func (this *promotionService) GetPromotion(id int) (*promotion.ValuePromotion, interface{}) {
	var prom promotion.IPromotion = this._newRep.GetPromotion(id)
	if prom != nil {
		return prom.GetValue(), prom.GetRelationValue()
	}
	return nil, nil
}

// 保存促销
func (this *promotionService) SavePromotion(v *promotion.ValuePromotion) (int, error) {
	var prom promotion.IPromotion
	if v.Id > 0 {
		prom = this._newRep.GetPromotion(v.Id)
		prom.SetValue(v)
	} else {
		prom = this._newRep.CreatePromotion(v)
	}
	return prom.Save()
}

// 删除促销
func (this *promotionService) DelPromotion(partnerId int,promId int)error{
	prom := this._newRep.GetPromotion(promId)
	if prom == nil{
		return promotion.ErrNoSuchPromotion
	}
	if prom.GetValue().PartnerId != partnerId{
		return partner.ErrNotMatch
	}

	return promImpl.DeletePromotion(prom)
}

func (this *promotionService) SaveCashBackPromotion(partnerId int, v *promotion.ValuePromotion,
	v1 *promotion.ValueCashBack) (int, error) {
	var prom promotion.IPromotion
	var err error
	if v.Id > 0 {
		prom = this._newRep.GetPromotion(v.Id)
		if prom.GetValue().PartnerId != partnerId {
			return -1, partner.ErrNotMatch
		}
	} else {
		prom = this._newRep.CreatePromotion(v)
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
func (this *promotionService) GetCoupon(partnerId int, id int) promotion.ICouponPromotion {
	_prom := this._rep.GetPromotion(partnerId)
	return _prom.GetCoupon(id)
}

func (this *promotionService) SaveCoupon(partnerId int, e *promotion.ValueCoupon) (int, error) {
	prom := this._rep.GetPromotion(partnerId)
	var coupon promotion.ICouponPromotion
	if e.Id > 0 {
		coupon = prom.GetCoupon(e.Id)
		if coupon == nil {
			return 0, errors.New("优惠券不存在")
		}
		err := coupon.SetDetailsValue(e)
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
