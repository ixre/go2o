/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:37
 * description :
 * history :
 */

package shopping

import (
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/promotion"
)

const (
	PaymentNotYet  = 0 // 订单尚未支付
	PaymentByBuyer = 1 // 购买者支付
	PaymentByCM    = 2 // 客服人工支付
)

type (
	IOrder interface {
		GetDomainId() int

		// 获取订单号
		GetOrderNo() string

		// 获生成值
		GetValue() ValueOrder

		// 设置订单值
		SetValue(*ValueOrder) error

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

		// 客服使用余额支付
		CmPaymentWithBalance() error

		// 在线交易支付
		PaymentForOnlineTrade(serverProvider string, tradeNo string) error

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

		// 配送订单,并记录配送服务商编号及单号
		Deliver(spId int, spNo string) error

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

	// 订单商品项
	OrderItem struct {
		Id         int     `db:"id" pk:"yes" auto:"yes" json:"id"`
		OrderId    int     `db:"order_id"`
		SnapshotId int     `db:"snapshot_id"`
		Quantity   int     `db:"quantity"`
		Sku        string  `db:"sku"`
		Fee        float32 `db:"fee"`
		UpdateTime int64   `db:"update_time"`
	}

	// 简单商品信息
	OrderGoods struct {
		GoodsId    int    `json:"id"`
		GoodsImage string `json:"img"`
		Name       string `json:"name"`
		Quantity   int    `json:"qty"`
	}

	OrderLog struct {
		//Id int `db:"id" auto:"yes" pk:"yes"`
		OrderId    int    `db:"order_id"`
		Type       int    `db:"type"`
		IsSystem   int    `db:"is_system"`
		Message    string `db:"message"`
		RecordTime int64  `db:"record_time"`
	}
	OrderPromotionBind struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`

		// 促销编号
		PromotionId int `db:"promotion_id"`

		// 促销类型
		PromotionType int `db:"promotion_type"`

		// 订单号
		OrderNo string `db:"order_no"`

		// 标题
		Title string `db:"title"`

		// 节省金额
		SaveFee float32 `db:"save_fee"`

		// 赠送积分
		PresentIntegral int `db:"present_integral"`

		// 是否应用
		IsApply int `db:"is_apply"`

		// 是否确认
		IsConfirm int `db:"is_confirm"`
	}

	// 应用到订单的优惠券
	OrderCoupon struct {
		OrderId      int     `db:"order_id"`
		CouponId     int     `db:"coupon_id"`
		CouponCode   string  `db:"coupon_code"`
		Fee          float32 `db:"coupon_fee"`
		Describe     string  `db:"coupon_describe"`
		SendIntegral int     `db:"send_integral"`
	}

	ValueOrder struct {
		Id         int    `db:"id" pk:"yes" auto:"yes" json:"id"`
		OrderNo    string `db:"order_no" json:"orderNo"`
		MemberId   int    `db:"member_id" json:"memberId"`
		MerchantId int    `db:"merchant_id" json:"merchantId"`
		// 订单标题
		Subject   string `db:"subject" json:"subject"`
		ShopId    int    `db:"shop_id" json:"shopId"`
		ItemsInfo string `db:"items_info" json:"itemsInfo"`
		// 总金额
		TotalFee float32 `db:"total_fee" json:"totalFee"`
		// 实际金额
		Fee float32 `db:"fee" json:"fee"`
		// 支付金额
		PayFee float32 `db:"pay_fee" json:"payFee"`
		// 减免金额(包含优惠券金额)
		DiscountFee float32 `db:"discount_fee" json:"discountFee"`
		// 余额抵扣
		BalanceDiscount float32 `db:"balance_discount" json:"balaceDiscount"`
		// 优惠券优惠金额
		CouponFee float32 `db:"coupon_fee" json:"couponFee"`
		// 支付方式
		PaymentOpt int `db:"payment_opt" json:"payMethod"`

		IsPaid int `db:"is_paid" json:"isPaid"`

		// 是否为顾客付款
		PaymentSign int `db:"payment_sign" json:"paymentSign"`

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

		// 订单项
		Items []*OrderItem `db:"-"`
	}
)

func (this *OrderCoupon) Clone(coupon promotion.ICouponPromotion,
	orderId int, orderFee float32) *OrderCoupon {
	v := coupon.GetDetailsValue()
	this.CouponCode = v.Code
	this.CouponId = v.Id
	this.OrderId = orderId
	this.Fee = coupon.GetCouponFee(orderFee)
	this.Describe = coupon.GetDescribe()
	this.SendIntegral = v.Integral
	return this
}
