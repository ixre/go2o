/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:18
 * description :
 * history :
 */

package rsi

import (
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/promotion"
	promImpl "go2o/core/domain/promotion"
)

type promotionService struct {
	_rep promotion.IPromotionRepo
}

func NewPromotionService(rep promotion.IPromotionRepo) *promotionService {
	return &promotionService{
		_rep: rep,
	}
}

// 获取促销
func (p *promotionService) GetPromotion(id int32) (*promotion.PromotionInfo, interface{}) {
	var prom promotion.IPromotion = p._rep.GetPromotion(id)
	if prom != nil {
		return prom.GetValue(), prom.GetRelationValue()
	}
	return nil, nil
}

// 保存促销
func (p *promotionService) SavePromotion(v *promotion.PromotionInfo) (int32, error) {
	var prom promotion.IPromotion
	if v.Id > 0 {
		prom = p._rep.GetPromotion(v.Id)
		err := prom.SetValue(v)
		if err != nil {
			return v.Id, err
		}
	} else {
		prom = p._rep.CreatePromotion(v)
	}
	return prom.Save()
}

// 删除促销
func (p *promotionService) DelPromotion(mchId int32, promId int32) error {
	prom := p._rep.GetPromotion(promId)
	if prom == nil {
		return promotion.ErrNoSuchPromotion
	}
	if prom.GetValue().MerchantId != mchId {
		return merchant.ErrMerchantNotMatch
	}

	return promImpl.DeletePromotion(prom)
}

func (p *promotionService) SaveCashBackPromotion(mchId int32,
	v *promotion.PromotionInfo, v1 *promotion.ValueCashBack) (int32, error) {
	var prom promotion.IPromotion
	var err error
	if v.Id > 0 {
		prom = p._rep.GetPromotion(v.Id)
		if prom.GetValue().MerchantId != mchId {
			return -1, merchant.ErrMerchantNotMatch
		}
	} else {
		prom = p._rep.CreatePromotion(v)
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
func (p *promotionService) SaveCoupon(mchId int32, v *promotion.PromotionInfo,
	v1 *promotion.ValueCoupon) (int32, error) {
	var prom promotion.IPromotion
	var err error
	if v.Id > 0 {
		prom = p._rep.GetPromotion(v.Id)
		if prom.GetValue().MerchantId != mchId {
			return -1, merchant.ErrMerchantNotMatch
		}
	} else {
		prom = p._rep.CreatePromotion(v)
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

func (p *promotionService) BindCoupons(mchId int32, id int32, members []string) error {
	coupon := p._rep.GetPromotion(id).(promotion.ICouponPromotion)
	return coupon.Binds(members)
}
