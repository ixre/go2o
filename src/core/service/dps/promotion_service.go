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
	"go2o/src/core/domain/interface/promotion"
)

type promotionService struct {
	_rep promotion.IOldPromotionRep
	_newRep promotion.IPromotionRep
}

func NewPromotionService(r promotion.IOldPromotionRep,rep promotion.IPromotionRep) *promotionService {
	return &promotionService{
		_rep: r,
		_newRep:rep,
	}
}

// 获取促销
func (this *promotionService) GetPromotion(id int)(*promotion.ValuePromotion,interface{}){
	var prom promotion.IPromotion = this._newRep.GetPromotion(id)
	if prom != nil{
		return prom.GetValue(),prom.GetRelationValue()
	}
	return nil,nil
}

// 保存促销
func (this *promotionService) SavePromotion(v *promotion.ValuePromotion)(int,error){
	var prom promotion.IPromotion
	if v.Id > 0{
		prom = this._newRep.GetPromotion(v.Id)
		prom.SetValue(v)
	}else{
		prom = this._newRep.CreatePromotion(v)
	}
	return prom.Save()
}


/**************   Coupon ************/
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
