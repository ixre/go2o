/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:25
 * description :
 * history :
 */

package promotion

type ValueCouponTake struct {
	Id int `db:"id" auto:"yes" pk:"yes"`

	//会员编号
	MemberId int `db:"member_id"`

	//优惠券编号
	CouponId int `db:"coupon_id"`

	//占用时间
	TakeTime int64 `db:"take_time"`

	//释放时间,超过该时间，优惠券释放
	ExtraTime int64 `db:"extra_time"`

	//是否应用到订单
	IsApply int `db:"is_apply"`

	//更新时间
	ApplyTime int64 `db:"apply_time"`
}
