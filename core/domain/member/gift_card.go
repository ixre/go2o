/**
 * Copyright 2015 @ z3q.net.
 * name : gift_card
 * author : jarryliu
 * date : 2016-06-24 17:04
 * description :
 * history :
 */
package member

import (
	"go2o/core/domain/interface/member"
	"go2o/core/dto"
)

var _ member.IGiftCardManager = new(giftCardManagerImpl)

type giftCardManagerImpl struct {
	_memberId int
	_rep      member.IMemberRep
}

func newGiftCardManagerImpl(memberId int, rep member.IMemberRep) member.IGiftCardManager {
	return &giftCardManagerImpl{
		_memberId: memberId,
		_rep:      rep,
	}
}

// 可用的优惠券分页数据
func (this *giftCardManagerImpl) PagedAvailableCoupon(start, end int) (total int, rows []*dto.ValueCoupon) {
	return this._rep.GetMemberPagedCoupon(this._memberId, start, end, "")
}

// 已使用的优惠券
func (this *giftCardManagerImpl) PagedAllCoupon(start, end int) (total int, rows []*dto.ValueCoupon) {
	return this._rep.GetMemberPagedCoupon(this._memberId, start, end, "")
}

// 过期的优惠券
func (this *giftCardManagerImpl) PagedExpiresCoupon(start, end int) (total int, rows []*dto.ValueCoupon) {
	return this._rep.GetMemberPagedCoupon(this._memberId, start, end, "")
}
