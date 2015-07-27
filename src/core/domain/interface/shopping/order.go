/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:37
 * description :
 * history :
 */

package shopping

import (
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/promotion"
)

type IOrder interface {
	GetDomainId() int

	// 获生成值
	GetValue() ValueOrder

	// 应用优惠券
	ApplyCoupon(coupon promotion.ICouponPromotion) error

	// 获取应用的优惠券
	GetCoupons() []promotion.ICouponPromotion

	// 获取可用的促销,不包含优惠券
	GetAvailableOrderPromotions() []promotion.IPromotion

	// 获取最省的促销
	GetBestSavePromotion() (p promotion.IPromotion, saveFee float32, integral int)

	// 获取促销绑定
	GetPromotionBinds() []*OrderPromotionBind

	// 设置Shop,如果不需要记录日志，则remark传递空
	SetShop(shopId int) error

	// 设置支付方式
	SetPayment(payment int)

	// 使用余额支付
	PaymentWithBalance() error

	// 在线交易支付
	PaymentOnlineTrade(serverProvider string,tradeNo string)error

	// 设置配送地址
	SetDeliver(deliverAddressId int) error

	// 添加备注
	AddRemark(string)

	// 应用余额支付
	UseBalanceDiscount()

	// 提交订单，返回订单号。如有错误则返回
	Submit() (string, error)

	// 保存订单
	Save() (int, error)

	// 添加日志,system表示为系统日志
	AppendLog(t enum.OrderLogType, system bool, message string) error

	// 订单是否结束
	IsOver() bool

	// 处理订单
	Process() error

	// 确认订单
	Confirm() error

	// 配送订单
	Deliver() error

	// 挂起
	Suspend(reason string) error

	// 标记收货
	SignReceived() error

	// 获取支付金额
	GetPaymentFee() float32

	// 完成订单
	Complete() error

	// 取消订单
	Cancel(reason string) error
}
