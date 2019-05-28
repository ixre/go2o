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

import "go2o/core/infrastructure/domain"

var (
	ErrCanNotApplied *domain.DomainError = domain.NewError(
		"name_exists", "无法应用此优惠")

	ErrExistsSamePromotionFlag *domain.DomainError = domain.NewError(
		"exists_same_promotion_flag", "已存在相同的促销")

	ErrNoSuchPromotion *domain.DomainError = domain.NewError(
		"no_such_promotion", "促销不存在")

	ErrNoDetailsPromotion *domain.DomainError = domain.NewError(
		"no_details_promotion", "促销信息不完整")
)

// 促销聚合根
type IPromotion interface {
	// 获取聚合根编号
	GetAggregateRootId() int32

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
	Save() (int32, error)

	// 获取优惠券
	//GetCoupon(id int32) ICouponPromotion

	// 创建优惠券
	//CreateCoupon(val *ValueCoupon) ICouponPromotion
}

type PromotionInfo struct {
	// 促销编号
	Id int32 `db:"id" pk:"yes" auto:"yes"`

	// 商户编号
	MerchantId int32 `db:"mch_id"`

	// 促销简称
	ShortName string `db:"short_name"`

	// 促销描述
	Description string `db:"description"`

	// 类型位值
	TypeFlag int `db:"type_flag"`

	// 商品编号(为0则应用订单)
	GoodsId int64 `db:"goods_id"`

	// 是否启用
	Enabled int `db:"enabled"`

	// 修改时间
	UpdateTime int64 `db:"update_time"`
}
