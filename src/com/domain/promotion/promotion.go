/**
 * Copyright 2014 @ Ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-04 13:48
 * description :
 * history :
 */

package promotion

import (
	"com/domain/interface/member"
	"com/domain/interface/promotion"
)

type Promotion struct {
	promRep   promotion.IPromotionRep
	memberRep member.IMemberRep
	partnerId int
}

func NewPromotion(partnerId int, promRep promotion.IPromotionRep,
	memberRep member.IMemberRep) promotion.IPromotion {
	return &Promotion{
		partnerId: partnerId,
		promRep:   promRep,
		memberRep: memberRep,
	}
}

func (this *Promotion) GetAggregateRootId() int {
	return this.partnerId
}

func (this *Promotion) GetCoupon(id int) promotion.ICoupon {
	var val *promotion.ValueCoupon = this.promRep.GetCoupon(id)
	return this.CreateCoupon(val)
}

func (this *Promotion) CreateCoupon(val *promotion.ValueCoupon) promotion.ICoupon {
	return newCoupon(val, this.promRep, this.memberRep)
}
