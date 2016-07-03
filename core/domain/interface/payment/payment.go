/**
 * Copyright 2015 @ z3q.net.
 * name : payment
 * author : jarryliu
 * date : 2016-07-02 23:06
 * description : 支付单据
 * history :
 */

// 支付单,不限于订单,可以生成支付单,即一个支付请求
package payment

import (
	"go2o/core/domain/interface/promotion"
	"go2o/core/infrastructure/domain"
)

const (
	PaymentByBuyer = 1 // 购买者支付
	PaymentByCM    = 2 // 客服人工支付

	StateNotYetPayment = 0 // 尚未支付
	StateFinishPayment = 1 // 已支付
	StateHasCancel     = 2 // 已经取消
)

const (
	// 允许余额抵扣
	OptBalanceDiscount = 1 << iota
	// 允许积分抵扣
	OptIntegralDiscount
	// 允许系统支付
	OptSystemPayment
	// 允许使用优惠券
	OptUseCoupon

	// 全部支付权限
	OptPerm = OptBalanceDiscount | OptIntegralDiscount |
		OptSystemPayment | OptUseCoupon
)

var (
	ErrPaymentNotSave *domain.DomainError = domain.NewDomainError(
		"err_payment_not_save", "支付单需保存后才能执行操作")

	ErrOrderPayed *domain.DomainError = domain.NewDomainError(
		"err_payment_order_payed", "订单已支付")

	ErrOrderHasCancel *domain.DomainError = domain.NewDomainError(
		"err_payment_order_has_cancel", "订单已经取消")

	ErrOrderNotPayed *domain.DomainError = domain.NewDomainError(
		"err_payment_order_not_payed", "订单未支付")

	ErrCanNotUseBalance *domain.DomainError = domain.NewDomainError(
		"err_can_not_use_balance", "不能使用余额支付")

	ErrCanNotUseIntegral *domain.DomainError = domain.NewDomainError(
		"err_can_not_use_integral", "不能使用积分抵扣")

	ErrCanNotUseCoupon *domain.DomainError = domain.NewDomainError(
		"err_can_not_use_coupon", "不能使用优惠券")

	ErrCanNotSystemDiscount *domain.DomainError = domain.NewDomainError(
		"err_can_not_system_discount", "不允许系统支付")
)

type (
	/// <summary>
	/// 支付单接口
	/// </summary>
	IPaymentOrder interface {
		// 获取聚合根编号
		GetAggregateRootId() int
		/// <summary>
		/// 优惠券抵扣
		/// </summary>
		CouponDiscount(coupon promotion.ICouponPromotion) (float32, error)
		BalanceDiscount(fee float32) error
		IntegralDiscount(integral int) error
		/// <summary>
		/// 系统支付金额
		/// </summary>
		SystemPayment(fee float32) error
		BindOrder(orderId int) error
		Save() (int, error)
		PaymentFinish(tradeNo string) error
		GetValue() PaymentOrderBean
		/// <summary>
		/// 取消支付
		/// </summary>
		Cancel() error
	}

	IPaymentRep interface {
		// 根据编号获取支付单
		GetPaymentOrder(id int) IPaymentOrder
		// 根据支付单号获取支付单
		GetPaymentOrderByNo(paymentNo string) IPaymentOrder
		// 创建支付单
		CreatePaymentOrder(p *PaymentOrderBean) IPaymentOrder
		// 保存支付单
		SavePaymentOrder(v *PaymentOrderBean) (int, error)
	}

	PaymentOrderBean struct {
		Id int
		/// <summary>
		/// 支付单号
		/// </summary>
		PaymentNo string
		/// <summary>
		/// 运营商编号，0表示无
		/// </summary>
		VendorId int
		/// <summary>
		/// 订单编号,0表示无
		/// </summary>
		OrderId int
		/// <summary>
		/// 购买用户
		/// </summary>
		BuyUser int
		/// <summary>
		/// 支付用户
		/// </summary>
		PaymentUser int
		/// <summary>
		/// 支付单金额
		/// </summary>
		TotalFee float32
		/// <summary>
		/// 余额抵扣
		/// </summary>
		BalanceDiscount float32
		/// <summary>
		/// 积分抵扣
		/// </summary>
		IntegralDiscount float32
		/// <summary>
		/// 系统支付抵扣金额
		/// </summary>
		SystemDiscount float32
		/// <summary>
		/// 优惠券金额
		/// </summary>
		CouponFee float32
		/// <summary>
		/// 立减金额
		/// </summary>
		SubFee float32
		/// <summary>
		/// 最终支付金额
		/// </summary>
		FinalFee float32
		/// <summary>
		/// 支付选项，位运算。可用优惠券，积分抵扣等运算
		/// </summary>
		PaymentOpt int
		/// 支付方式
		PaymentSign int
		//创建时间
		CreateTime int64
		/// 在线支付的交易单号
		TradeNo string
		//支付时间
		PaidTime int64
		/// <summary>
		/// 状态:  0为未付款，1为已付款，2为已取消
		/// </summary>
		State int
	}
)
