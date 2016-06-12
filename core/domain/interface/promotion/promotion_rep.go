/**
 * Copyright 2014 @ z3q.net.
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
	CreatePromotion(v *PromotionInfo) IPromotion

	// 获取促销
	GetValuePromotion(id int) *PromotionInfo

	// 保存促销
	SaveValuePromotion(*PromotionInfo) (int, error)

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

	// 获取商品可用的促销
	GetPromotionOfGoods(goodsId int) []*PromotionInfo

	// 获取商户订单可用的促销
	GetPromotionOfMerchantOrder(merchantId int) []*PromotionInfo

	/** =======  优惠券 ==========**/

	// 获取优惠券
	GetValueCoupon(id int) *ValueCoupon

	// 保存优惠券值
	SaveValueCoupon(v *ValueCoupon, isCreate bool) (id int, err error)

	// 删除优惠券
	DeleteValueCoupon(id int) error

	GetCouponTake(couponId, takeId int) *ValueCouponTake

	SaveCouponTake(*ValueCouponTake) error

	GetCouponTakes(couponId int) []ValueCouponTake

	GetCouponBind(couponId, bindId int) *ValueCouponBind

	SaveCouponBind(*ValueCouponBind) error

	GetCouponBinds(couponId int) []ValueCouponBind

	// 根据优惠券代码获取优惠券
	GetValueCouponByCode(merchantId int, couponCode string) *ValueCoupon

	// 根据代码获取优惠券
	GetCouponByCode(merchantId int, code string) IPromotion

	// 获取会员的优惠券绑定
	GetCouponBindByMemberId(couponId, memberId int) (*ValueCouponBind, error)

	// 获取会员的优惠券占用
	GetCouponTakeByMemberId(couponId, memberId int) (*ValueCouponTake, error)
}
