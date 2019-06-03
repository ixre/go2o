/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 21:04
 * description :
 * history :
 */

package promotion

import (
	"go2o/core/domain/interface/member"
)

// 优惠券促销
type ICouponPromotion interface {
	GetDomainId() int32

	// 获取促销内容
	GetDetailsValue() ValueCoupon

	// 设置促销内容
	SetDetailsValue(*ValueCoupon) error

	// 获取绑定
	GetBinds() []ValueCouponBind

	// 获取占用
	GetTakes() []ValueCouponTake

	// 保存
	Save() (id int32, err error)

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
	GetTake(memberId int64) (*ValueCouponTake, error)

	//占用
	Take(memberId int64) error

	// 应用到订单
	ApplyTake(couponTakeId int32) error

	// 绑定
	Bind(memberId int64) error

	//获取绑定
	GetBind(memberId int64) (*ValueCouponBind, error)

	//绑定
	Binds(memberIds []string) error

	//使用优惠券
	UseCoupon(couponBindId int32) error
}

type (
	ValueCoupon struct {
		// 优惠券编号
		Id int32 `db:"id" pk:"yes"`

		//优惠码
		Code string `db:"code"`

		// 优惠码可用数量
		Amount int `db:"amount"`

		// 优惠码数量
		TotalAmount int `db:"total_amount"`

		//优惠金额
		Fee int `db:"fee"`

		//赠送积分
		Integral int `db:"integral"`

		//订单折扣(不打折为100)
		Discount int `db:"discount"`

		//等级限制
		MinLevel int `db:"min_level"`

		//订单金额限制
		MinFee int `db:"min_fee"`

		BeginTime int64 `db:"begin_time"`
		OverTime  int64 `db:"over_time"`

		//是否需要绑定。反之可以直接使用
		NeedBind int `db:"need_bind"`

		CreateTime int64 `db:"create_time"`
	}

	//优惠券绑定
	ValueCouponBind struct {
		Id int32 `db:"id" auto:"yes" pk:"yes"`

		//会员编号
		MemberId int64 `db:"member_id"`

		//优惠券编号
		CouponId int32 `db:"coupon_id"`

		//绑定时间
		BindTime int64 `db:"bind_time"`

		//是否使用
		IsUsed int `db:"is_used"`

		//使用时间
		UseTime int64 `db:"use_time"`
	}

	ValueCouponTake struct {
		Id int32 `db:"id" auto:"yes" pk:"yes"`

		//会员编号
		MemberId int64 `db:"member_id"`

		//优惠券编号
		CouponId int32 `db:"coupon_id"`

		//占用时间
		TakeTime int64 `db:"take_time"`

		//释放时间,超过该时间，优惠券释放
		ExtraTime int64 `db:"extra_time"`

		//是否应用到订单
		IsApply int `db:"is_apply"`

		//更新时间
		ApplyTime int64 `db:"apply_time"`
	}
)
