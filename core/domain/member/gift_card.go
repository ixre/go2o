/**
 * Copyright 2015 @ to2.net.
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
	memberId int64
	rep      member.IMemberRepo
}

func newGiftCardManagerImpl(memberId int64, rep member.IMemberRepo) member.IGiftCardManager {
	return &giftCardManagerImpl{
		memberId: memberId,
		rep:      rep,
	}
}

// 可用的优惠券分页数据
func (g *giftCardManagerImpl) PagedAvailableCoupon(start, end int) (
	total int, rows []*dto.SimpleCoupon) {
	// 未使用,且未过有效期的
	unix := time.Now().Unix()
	return g.rep.GetMemberPagedCoupon(g.memberId, start, end,
		fmt.Sprintf("over_time > %d AND is_used = 0", unix))
}

// 所有的优惠券
func (g *giftCardManagerImpl) PagedAllCoupon(start, end int) (
	total int, rows []*dto.SimpleCoupon) {
	return g.rep.GetMemberPagedCoupon(g.memberId, start, end, "1=1")
}

// 过期的优惠券
func (g *giftCardManagerImpl) PagedExpiresCoupon(start, end int) (
	total int, rows []*dto.SimpleCoupon) {
	//未使用且已超过有效期
	unix := time.Now().Unix()
	return g.rep.GetMemberPagedCoupon(g.memberId, start, end,
		fmt.Sprintf("over_time < %d AND is_used = 0", unix))
}
