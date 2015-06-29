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
	// 获取促销
	GetPromotion(id int) IPromotion

	// 获取促销
	CreatePromotion(v *ValuePromotion, dv interface{}) IPromotion

	// 获取促销
	GetValuePromotion(id int) *ValuePromotion

	// 保存促销
	SaveValuePromotion(*ValuePromotion) (int, error)

	// 删除促销
	DeletePromotion(id int) error

	// 保存返现促销
	SaveValueCashBack(v *ValueCashBack, create bool) (int, error)

	// 获取返现促销
	GetValueCashBack(int) *ValueCashBack

	// 删除现金返现促销
	DeleteValueCashBack(id int) error

	// 获取商品的促销编号
	GetGoodsPromotionId(goodsId int, promFlag int) int

	/** =======  优惠券 ==========**/

	// 获取优惠券
	GetValueCoupon(id int) *ValueCoupon

	// 保存优惠券值
	SaveValueCoupon(v *ValueCoupon,isCreate bool) (id int, err error)

	GetCouponTake(couponId, takeId int) *ValueCouponTake

	SaveCouponTake(*ValueCouponTake) error

	GetCouponTakes(couponId int) []ValueCouponTake

	GetCouponBind(couponId, bindId int) *ValueCouponBind

	SaveCouponBind(*ValueCouponBind) error

	GetCouponBinds(couponId int) []ValueCouponBind

	// 根据优惠券代码获取优惠券
	GetValueCouponByCode(partnerId int, couponCode string) *ValueCoupon

	// 根据代码获取优惠券
	GetCouponByCode(partnerId int, code string) IPromotion

	// 获取会员的优惠券绑定
	GetCouponBindByMemberId(couponId, memberId int) (*ValueCouponBind, error)

	// 获取会员的优惠券占用
	GetCouponTakeByMemberId(couponId, memberId int) (*ValueCouponTake, error)
}
