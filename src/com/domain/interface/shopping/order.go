/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-05 17:37
 * description :
 * history :
 */

package shopping

import "com/domain/interface/promotion"

type IOrder interface {
	GetDomainId() int
	// 获生成值
	GetValue() ValueOrder
	// 应用优惠券
	ApplyCoupon(coupon promotion.ICoupon) error
	// 获取应用的优惠券
	GetCoupons() []promotion.ICoupon

	// 设置Shop
	SetShop(shopId int) error
	// 设置支付方式
	SetPayment(payment int)
	// 设置配送地址
	SetDeliver(deliverAddrId int) error
	// 添加备注
	AddRemark(string)
	// 提交订单，返回订单号。如有错误则返回
	Submit() (string, error)
	// 保存订单
	Save() error
	// 订单是否结束
	IsOver() bool
	// 处理订单
	Process() error
	// 配送订单
	Deliver() error
	// 完成订单
	Finish() error

	// 取消订单
	Cancel() error
}
