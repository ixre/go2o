/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-03 14:11
 * description :
 * history :
 */

package promotion

type IPromotionRep interface {
	GetPromotion(partnerId int) IPromotion
	GetCoupon(id int) *ValueCoupon
	SaveCoupon(ValueCoupon) (id int, err error)

	GetCouponTake(couponId, takeId int) *ValueCouponTake

	SaveCouponTake(*ValueCouponTake) error

	GetCouponTakes(couponId int) []ValueCouponTake

	GetCouponBind(couponId, bindId int) *ValueCouponBind

	SaveCouponBind(*ValueCouponBind) error

	GetCouponBinds(couponId int) []ValueCouponBind

	// 根据优惠券编号获取优惠券
	GetCouponByCode(partnerId int, couponCode string) (ICoupon, error)

	// 获取会员的优惠券绑定
	GetCouponBindByMemberId(couponId, memberId int) (*ValueCouponBind, error)

	// 获取会员的优惠券占用
	GetCouponTakeByMemberId(couponId, memberId int) (*ValueCouponTake, error)
}
