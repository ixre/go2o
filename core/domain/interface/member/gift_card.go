/**
 * Copyright 2015 @ to2.net.
 * name : gift.go
 * author : jarryliu
 * date : 2016-06-24 16:50
 * description :
 * history :
 */
package member

import (
	"go2o/core/dto"
)

/** 礼品/卡/券  **/

type (
	IGiftCardManager interface {
		//todo: ???
		// 领用优惠券
		// TakeCoupon()

		// 可用的优惠券分页数据
		PagedAvailableCoupon(start, end int) (total int, rows []*dto.SimpleCoupon)

		// 所有的优惠券
		PagedAllCoupon(start, end int) (total int, rows []*dto.SimpleCoupon)

		// 过期的优惠券
		PagedExpiresCoupon(start, end int) (total int, rows []*dto.SimpleCoupon)
	}
)
