/**
 * Copyright 2015 @ 56x.net.
 * name : gift.go
 * author : jarryliu
 * date : 2016-06-24 16:50
 * description :
 * history :
 */
package member

/** 礼品/卡/券  **/

type (
	IGiftCardManager interface {
		//todo: ???
		// 领用优惠券
		// TakeCoupon()

		// 可用的优惠券分页数据
		PagedAvailableCoupon(start, end int) (total int, rows []*SimpleCouponQueryObject)

		// 所有的优惠券
		PagedAllCoupon(start, end int) (total int, rows []*SimpleCouponQueryObject)

		// 过期的优惠券
		PagedExpiresCoupon(start, end int) (total int, rows []*SimpleCouponQueryObject)
	}

	SimpleCouponQueryObject struct {
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
