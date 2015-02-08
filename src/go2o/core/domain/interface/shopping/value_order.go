/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-05 17:19
 * description :
 * history :
 */

package shopping

type ValueOrder struct {
	Id        int    `db:"id" pk:"yes" auto:"yes" json:"id"`
	OrderNo   string `db:"order_no" json:"orderNo"`
	MemberId  int    `db:"member_id" json:"memberId"`
	PartnerId int    `db:"pt_id" json:"partnerId"`
	ShopId    int    `db:"shop_id" json:"shopId"`
	ItemsInfo string `db:"items_info" json:"itemsInfo"`
	// 总金额
	TotalFee float32 `db:"total_fee" json:"totalFee"`
	// 实际金额
	Fee float32 `db:"fee" json:"fee"`
	// 支付金额
	PayFee float32 `db:"pay_fee" json:"payFee"`
	//减免金额(包含优惠券金额)
	DiscountFee float32 `db:"discount_fee" json:"discountFee"`
	//优惠券优惠金额
	CouponFee float32 `db:"coupon_fee" json:"couponFee"`

	PayMethod int `db:"pay_method" json:"payMethod"`

	IsPaid int `db:"is_paid" json:"isPaid"`

	// 是否挂起，如遇到无法自动进行的时挂起，来提示人工确认。
	IsSuspend int `db:"is_suspend" json:"is_suspend"`

	Note string `db:"note" json:"note"`

	Remark string `db:"note" json:"remark"`

	// 支付时间
	PaidTime int64 `db:"paid_time" json:"paidTime"`

	DeliverName    string `db:"deliver_name" json:"deliverName"`
	DeliverPhone   string `db:"deliver_phone" json:"deliverPhone"`
	DeliverAddress string `db:"deliver_address" json:"deliverAddress"`
	DeliverTime    int64  `db:"deliver_time" json:"deliverTime"`
	CreateTime     int64  `db:"create_time" json:"createTime"`

	// 订单状态
	Status int `db:"status" json:"status"`

	UpdateTime int64 `db:"update_time" json:"updateTime"`
}
