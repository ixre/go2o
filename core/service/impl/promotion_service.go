/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:18
 * description :
 * history :
 */

package impl

import (
	"fmt"

	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/promotion"
	promImpl "github.com/ixre/go2o/core/domain/promotion"
)

type PromotionService struct {
	_repo promotion.IPromotionRepo
}

func NewPromotionService(rep promotion.IPromotionRepo) *PromotionService {
	return &PromotionService{
		_repo: rep,
	}
}

// 获取促销
func (p *PromotionService) GetPromotion(id int32) (*promotion.PromotionInfo, interface{}) {
	var prom = p._repo.GetPromotion(id)
	if prom != nil {
		return prom.GetValue(), prom.GetRelationValue()
	}
	return nil, fmt.Errorf("no such promotion")
}

// 保存促销
func (p *PromotionService) SavePromotion(v *promotion.PromotionInfo) (int32, error) {
	var prom promotion.IPromotion
	if v.Id > 0 {
		prom = p._repo.GetPromotion(v.Id)
		err := prom.SetValue(v)
		if err != nil {
			return v.Id, err
		}
	} else {
		prom = p._repo.CreatePromotion(v)
	}
	return prom.Save()
}

// 删除促销
func (p *PromotionService) DelPromotion(mchId int64, promId int32) error {
	prom := p._repo.GetPromotion(promId)
	if prom == nil {
		return promotion.ErrNoSuchPromotion
	}
	if int64(prom.GetValue().MerchantId) != mchId {
		return merchant.ErrMerchantNotMatch
	}

	return promImpl.DeletePromotion(prom)
}

func (p *PromotionService) SaveCashBackPromotion(mchId int64,
	v *promotion.PromotionInfo, v1 *promotion.ValueCashBack) (int32, error) {
	var prom promotion.IPromotion
	var err error
	if v.Id > 0 {
		prom = p._repo.GetPromotion(v.Id)
		if int64(prom.GetValue().MerchantId) != mchId {
			return -1, merchant.ErrMerchantNotMatch
		}
	} else {
		prom = p._repo.CreatePromotion(v)
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
func (p *PromotionService) SaveCoupon(mchId int64, v *promotion.PromotionInfo,
	v1 *promotion.ValueCoupon) (int32, error) {
	var prom promotion.IPromotion
	var err error
	if v.Id > 0 {
		prom = p._repo.GetPromotion(v.Id)
		if int64(prom.GetValue().MerchantId) != mchId {
			return -1, merchant.ErrMerchantNotMatch
		}
	} else {
		prom = p._repo.CreatePromotion(v)
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

func (p *PromotionService) BindCoupons(mchId int64, id int32, members []string) error {
	coupon := p._repo.GetPromotion(id).(promotion.ICouponPromotion)
	return coupon.Binds(members)
}
