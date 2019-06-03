/**
 * Copyright 2015 @ to2.net.
 * name : promotion
 * author : jarryliu
 * date : 2016-06-24 17:52
 * description :
 * history :
 */
package dto

type (
	//member_rep LINE:562
	SimpleCoupon struct {
		// 优惠券编号
		Id int `db:"id"`

		Num int `db:"num"`

		// 优惠券标题
		Title string `db:"title"`

		//优惠码
		Code string `db:"code"`

		//优惠金额
		Fee int `db:"fee"`

		//订单折扣(不打折为100)
		Discount int `db:"discount"`

		//是否使用
		IsUsed int `db:"is_used"`

		//结束日期
		OverTime int64 `db:"over_time"`
	}
)
