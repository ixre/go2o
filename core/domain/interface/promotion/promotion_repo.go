/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 14:11
 * description :
 * history :
 */

package promotion

type IPromotionRepo interface {
	// 获取促销
	GetPromotion(id int32) IPromotion

	// 获取促销
	CreatePromotion(v *PromotionInfo) IPromotion

	// 获取促销
	GetValuePromotion(id int32) *PromotionInfo

	// 保存促销
	SaveValuePromotion(*PromotionInfo) (int32, error)

	// 删除促销
	DeletePromotion(id int32) error

	// 保存返现促销
	SaveValueCashBack(v *ValueCashBack, create bool) (int32, error)

	// 获取返现促销
	GetValueCashBack(id int32) *ValueCashBack

	// 删除现金返现促销
	DeleteValueCashBack(id int32) error

	// 获取商品的促销编号
	GetGoodsPromotionId(goodsId int64, promFlag int) int

	// 获取商品可用的促销
	GetPromotionOfGoods(goodsId int64) []*PromotionInfo

	// 获取商户订单可用的促销
	GetPromotionOfMerchantOrder(mchId int32) []*PromotionInfo

	/** =======  优惠券 ==========**/

	// 获取优惠券
	GetValueCoupon(id int32) *ValueCoupon

	// 保存优惠券值
	SaveValueCoupon(v *ValueCoupon, isCreate bool) (id int32, err error)

	// 删除优惠券
	DeleteValueCoupon(id int32) error

	GetCouponTake(couponId, takeId int32) *ValueCouponTake

	SaveCouponTake(*ValueCouponTake) error

	GetCouponTakes(couponId int32) []ValueCouponTake

	GetCouponBind(couponId, bindId int32) *ValueCouponBind

	SaveCouponBind(*ValueCouponBind) error

	GetCouponBinds(couponId int32) []ValueCouponBind

	// 根据优惠券代码获取优惠券
	GetValueCouponByCode(mchId int32, couponCode string) *ValueCoupon

	// 根据代码获取优惠券
	GetCouponByCode(mchId int32, code string) IPromotion

	// 获取会员的优惠券绑定
	GetCouponBindByMemberId(couponId int32, memberId int64) (*ValueCouponBind, error)

	// 获取会员的优惠券占用
	GetCouponTakeByMemberId(couponId int32, memberId int64) (*ValueCouponTake, error)
}
