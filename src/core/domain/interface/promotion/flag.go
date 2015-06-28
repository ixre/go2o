/**
 * Copyright 2015 @ S1N1 Team.
 * name : flag
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package promotion

const (
	// 优惠券
	TypeFlagCoupon = 1 << 0

	// 返现
	TypeFlagCashBack = 1 << 1

	//todo: other promotion type
	//TypeFlagCashBack = 1 << 2
)

const (
	// 应用订单
	ApplyForOrder = 1

	// 应用商品
	ApplyForGoods = 2
)
