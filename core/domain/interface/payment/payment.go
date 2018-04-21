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

	StateAwaitingPayment = 0 // 等待支付
	StateFinishPayment   = 1 // 已支付
	StateHasCancel       = 2 // 已经取消

	TypeShopping = 1 //购物
)

const (
	// 允许余额抵扣
	OptBalanceDiscount = 1 << iota
	// 允许积分抵扣
	OptIntegralDiscount
	// 允许钱包支付
	OptWalletPayment
	// 允许系统支付
	OptSystemPayment
	// 允许使用优惠券
	OptUseCoupon

	// 全部支付权限
	OptPerm = OptBalanceDiscount | OptWalletPayment | OptIntegralDiscount |
		OptSystemPayment | OptUseCoupon
)

const (
	// 线上支付
	SignOnlinePay int32 = 1
	// 线下支付
	SignOfflinePay int32 = 2
	// 钱包账户支付
	SignWalletAccount int32 = 3
)

var (
	ErrNoSuchPaymentOrder = domain.NewDomainError(
		"err_no_such_payment_order", "支付单不存在")

	ErrExistsTradeNo = domain.NewDomainError(
		"err_payment_exists_trade_no", "支付单号重复")

	ErrPaymentNotSave = domain.NewDomainError(
		"err_payment_not_save", "支付单需存后才能执行操作")

	ErrFinalFee = domain.NewDomainError(
		"err_final_fee", "支付单金额有误")

	ErrNotSupportPaymentOpt = domain.NewDomainError(
		"err_payment_not_support_opt", "不支持此支付方式,无法完成付款")

	ErrTradeNoPrefix = domain.NewDomainError(
		"err_payment_trade_no_prefix", "支付单号前缀不正确")

	ErrTradeNoExistsPrefix = domain.NewDomainError(
		"err_payment_trade_no_exists_prefix", "支付单号已存在前缀")

	ErrOrderCommitted = domain.NewDomainError(
		"err_payment_order_committed", "支付单已提交")

	ErrOrderPayed = domain.NewDomainError(
		"err_payment_order_payed", "订单已支付")

	ErrOrderHasCancel = domain.NewDomainError("err_payment_order_has_cancel", "订单已经取消")

	ErrOrderNotPayed = domain.NewDomainError("err_payment_order_not_payed", "订单未支付")

	ErrCanNotUseBalance = domain.NewDomainError("err_can_not_use_balance", "不能使用余额支付")

	ErrNotEnoughAmount = domain.NewDomainError("err_payment_not_enough_amount", "余额不足,无法完成支付")

	ErrCanNotUseIntegral = domain.NewDomainError("err_can_not_use_integral", "不能使用积分抵扣")

	ErrCanNotUseCoupon = domain.NewDomainError("err_can_not_use_coupon", "不能使用优惠券")

	ErrCanNotSystemDiscount = domain.NewDomainError("err_can_not_system_discount", "不允许系统支付")

	ErrOuterNo = domain.NewDomainError("err_outer_no", "第三方交易号错误")
)

type (
	// 支付单接口
	IPaymentOrder interface {
		// 获取聚合根编号
		GetAggregateRootId() int32
		// 获取交易号
		GetTradeNo() string
		// 为交易号增加一个2位的前缀
		TradeNoPrefix(prefix string) error
		// 优惠券抵扣
		CouponDiscount(coupon promotion.ICouponPromotion) (float32, error)
		// 使用会员的余额抵扣
		BalanceDiscount(remark string) error
		// 使用会员积分抵扣,返回抵扣的金额及错误,ignoreOut:是否忽略超出订单金额的积分
		IntegralDiscount(integral int64, ignoreOut bool) (float32, error)
		// 系统支付金额
		SystemPayment(fee float32) error
		// 钱包账户支付
		PaymentByWallet(remark string) error
		// 余额钱包混合支付，优先扣除余额。
		HybridPayment(remark string) error
		// 设置支付方式
		SetPaymentSign(paymentSign int32) error
		// 绑定订单号,如果交易号为空则绑定参数中传递的交易号,
		// 支付单的交易号,可能是与订单号一样的
		BindOrder(orderId int64, tradeNo string) error
		// 提交支付单
		Commit() error
		// 支付完成并保存,传入第三名支付名称,以及外部的交易号
		PaymentFinish(spName string, outerNo string) error
		// 获取支付单的值
		GetValue() PaymentOrder
		// 取消支付
		Cancel() error
		// 调整金额,如调整金额与实付金额相加小于等于零,则支付成功。
		Adjust(amount float32) error
		// 退款
		Refund(amount float64) error
	}

	// 支付仓储
	IPaymentRepo interface {
		// 根据编号获取支付单
		GetPaymentOrderById(id int32) IPaymentOrder
		// 根据支付单号获取支付单
		GetPaymentOrder(paymentNo string) IPaymentOrder
		// 根据订单号获取支付单
		GetPaymentBySalesOrderId(orderId int64) IPaymentOrder
		// 创建支付单
		CreatePaymentOrder(p *PaymentOrder) IPaymentOrder
		// 保存支付单
		SavePaymentOrder(v *PaymentOrder) (int32, error)
		// 检查支付单号是否匹配
		CheckTradeNoMatch(tradeNo string, id int32) bool
		// 通知支付单完成
		//NotifyPaymentFinish(paymentOrderId int32) error
	}

	// 支付单实体
	PaymentOrder struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 支付单号
		TradeNo string `db:"trade_no"`
		// 交易类型
		TradeType string `db:"trade_type"`
		// 运营商编号，0表示无
		VendorId int32 `db:"vendor_id"`
		// 支付单类型,如果购物或其他
		Type int32 `db:"order_type"`
		// 订单编号,0表示无
		OrderId int32 `db:"order_id"`
		// 支付单主题
		Subject string `db:"subject"`
		// 购买用户
		BuyUser int64 `db:"buy_user"`
		// 支付用户
		PaymentUser int64 `db:"payment_user"`
		// 支付单金额
		TotalAmount float32 `db:"total_amount"`
		// 余额抵扣
		BalanceDiscount float32 `db:"balance_discount"`
		// 积分抵扣
		IntegralDiscount float32 `db:"integral_discount"`
		// 系统支付抵扣金额
		SystemDiscount float32 `db:"system_discount"`
		// 优惠券金额
		CouponDiscount float32 `db:"coupon_discount"`
		// 立减金额
		SubAmount float32 `db:"sub_amount"`
		// 调整的金额
		AdjustmentAmount float32 `db:"adjustment_amount"`
		// 最终支付金额
		FinalFee float32 `db:"final_fee"`
		// 支付选项，位运算。可用优惠券，积分抵扣等运算
		PayFlag int32 `db:"payment_opt"`
		// 支付方式
		PaymentSign int32 `db:"payment_sign"`
		// 在线支付的交易单号
		OuterNo string `db:"outer_no"`
		//创建时间
		CreateTime int64 `db:"create_time"`
		//支付时间
		PaidTime int64 `db:"paid_time"`
		// 状态:  0为未付款，1为已付款，2为已取消
		State int32 `db:"state"`
	}
)
