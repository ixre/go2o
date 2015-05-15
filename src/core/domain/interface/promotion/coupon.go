/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-03 21:04
 * description :
 * history :
 */

package promotion

import (
	"go2o/src/core/domain/interface/member"
)

type ICoupon interface {
	GetDomainId() int

	GetValue() ValueCoupon

	// 设置值
	SetValue(*ValueCoupon) error

	// 获取绑定
	GetBinds() []ValueCouponBind

	// 获取占用
	GetTakes() []ValueCouponTake

	// 保存
	Save() (id int, err error)

	// 获取优惠券描述
	GetDescribe() string

	// 获取优惠的金额
	GetCouponFee(orderFee float32) float32

	// 是否可用,传递会员及订单金额
	// error返回不适用的详细信息
	CanUse(member.IMember, float32) (bool, error)

	// 是否允许占用
	CanTake() bool
	// 获取占用
	GetTake(memberId int) (*ValueCouponTake, error)
	//占用
	Take(memberId int) error
	// 应用到订单
	ApplyTake(couponTakeId int) error

	// 绑定
	Bind(memberId int) error
	//获取绑定
	GetBind(memberId int) (*ValueCouponBind, error)
	//绑定
	Binds(memberIds []string) error
	//使用优惠券
	UseCoupon(couponBindId int) error
}
