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
	"fmt"
	"go2o/core/domain/interface/member"
	"go2o/core/dto"
	"time"
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
func (this *giftCardManagerImpl) PagedAvailableCoupon(start, end int) (
	total int, rows []*dto.ValueCoupon) {
	// 未使用,且未过有效期的
	unix := time.Now().Unix()
	return this._rep.GetMemberPagedCoupon(this._memberId, start, end,
		fmt.Sprintf("over_time > %d AND is_used = 0", unix))
}

// 所有的优惠券
func (this *giftCardManagerImpl) PagedAllCoupon(start, end int) (
	total int, rows []*dto.ValueCoupon) {
	return this._rep.GetMemberPagedCoupon(this._memberId, start, end, "1=1")
}

// 过期的优惠券
func (this *giftCardManagerImpl) PagedExpiresCoupon(start, end int) (
	total int, rows []*dto.ValueCoupon) {
	//未使用且已超过有效期
	unix := time.Now().Unix()
	return this._rep.GetMemberPagedCoupon(this._memberId, start, end,
		fmt.Sprintf("over_time < %d AND is_used = 0", unix))
}
