/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-10 10:46
 * description :
 * history :
 */

package shopping

import (
	"go2o/src/core/domain/interface/promotion"
)

// 应用到订单的优惠券
type OrderCoupon struct {
	OrderId      int     `db:"order_id"`
	CouponId     int     `db:"coupon_id"`
	CouponCode   string  `db:"coupon_code"`
	Fee          float32 `db:"coupon_fee"`
	Describe     string  `db:"coupon_describe"`
	SendIntegral int     `db:"send_integral"`
}

func (this *OrderCoupon) Clone(coupon promotion.ICouponPromotion,
	orderId int, orderFee float32) *OrderCoupon {
	v := coupon.GetValue()
	this.CouponCode = v.Code
	this.CouponId = v.Id
	this.OrderId = orderId
	this.Fee = coupon.GetCouponFee(orderFee)
	this.Describe = coupon.GetDescribe()
	this.SendIntegral = v.Integral
	return this
}
