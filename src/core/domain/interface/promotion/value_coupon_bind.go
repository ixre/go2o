/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:19
 * description :
 * history :
 */

package promotion

//优惠券绑定
type ValueCouponBind struct {
	Id int `db:"id" auto:"yes" pk:"yes"`

	//会员编号
	MemberId int `db:"member_id"`

	//优惠券编号
	CouponId int `db:"coupon_id"`

	//绑定时间
	BindTime int64 `db:"bind_time"`

	//是否使用
	IsUsed int `db:"is_used"`

	//使用时间
	UseTime int64 `db:"use_time"`
}
