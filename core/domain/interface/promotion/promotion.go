/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-04 23:46
 * description :
 * history :
 */

//整单折扣 (自动满１)
//满立减
//满就送 (送汤)
//**优惠券

package promotion

// 促销聚合根
type IPromotion interface {
	// 获取聚合根编号
	GetAggregateRootId() int

	// 获取值
	GetValue() *PromotionInfo

	// 获取相关的值
	GetRelationValue() interface{}

	// 设置值
	SetValue(*PromotionInfo) error

	// 应用类型
	ApplyFor() int

	// 促销类型
	Type() int

	// 获取类型名称
	TypeName() string

	// 保存
	Save() (int, error)

	// 获取优惠券
	//GetCoupon(id int) ICouponPromotion

	// 创建优惠券
	//CreateCoupon(val *ValueCoupon) ICouponPromotion
}
