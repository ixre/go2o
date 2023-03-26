/**
 * Copyright 2015 @ 56x.net.
 * name : payment
 * author : jarryliu
 * date : 2016-07-02 23:06
 * description : 支付单据
 * history :
 */

// Package payment 支付单,不限于订单,可以生成支付单,即一个支付请求
package payment

import (
	"github.com/ixre/go2o/core/domain/interface/promotion"
	"github.com/ixre/go2o/core/infrastructure/domain"
)

// 支付通道
const (
	// MBalance 余额抵扣通道
	MBalance = 1 << 0
	// MWallet 钱包支付通道
	MWallet = 1 << 1
	// MIntegral 积分兑换通道
	MIntegral = 1 << 2
	// MUserCard 用户卡通道
	MUserCard = 1 << 3
	// MUserCoupon 用户券通道
	MUserCoupon = 1 << 4
	// MCash 现金支付通道
	MCash = 1 << 5
	// MBankCard 银行卡支付通道(64)
	MBankCard = 1 << 6
	// MPaySP 第三方支付(128)
	MPaySP = 1 << 7
	// MSellerPay 卖家支付通道
	MSellerPay = 1 << 8
	// MSystemPay 系统支付通道
	MSystemPay = 1 << 9
)

// PAllFlag 所有支付方式
const PAllFlag = MBalance | MWallet | MIntegral | MUserCard |
	MUserCoupon | MCash | MBankCard | MPaySP | MSellerPay | MSystemPay

// 支付单状态
const (
	// StateAwaitingPayment 待支付
	StateAwaitingPayment = 1
	// StateFinished 已支付
	StateFinished = 2
	// StateCancelled 已关闭
	StateClosed = 3
	// StateRefunded 已退款
	StateRefunded = 4
)

var (
	ErrNoSuchPaymentOrder = domain.NewError(
		"err_no_such_payment_order", "支付单不存在")

	ErrExistsTradeNo = domain.NewError(
		"err_payment_exists_trade_no", "支付单号重复")

	ErrPaymentNotSave = domain.NewError(
		"err_payment_not_save", "支付单需存后才能执行操作")

	ErrFinalAmount = domain.NewError(
		"err_final_amount", "支付单金额有误")

	ErrNotSupportPaymentChannel = domain.NewError(
		"err_payment_not_support_channel", "不支持此支付方式,无法完成付款")
	ErrItemAmount       = domain.NewError("err_payment_item_amount", "支付单金额不能为零")
	ErrOutOfFinalAmount = domain.NewError("err_out_of_final_amount",
		"超出支付单金额")
	ErrNotMatchFinalAmount = domain.NewError("err_not_match_final_amount",
		"金额与实际金额不符，无法完成付款")
	ErrTradeNoPrefix = domain.NewError(
		"err_payment_trade_no_prefix", "支付单号前缀不正确")
	ErrTradeNoExistsPrefix = domain.NewError(
		"err_payment_trade_no_exists_prefix", "支付单号已存在前缀")

	ErrOrderCommitted = domain.NewError(
		"err_payment_order_committed", "支付单已提交")

	ErrOrderPayed = domain.NewError(
		"err_payment_order_payed", "订单已支付")

	ErrOrderClosed = domain.NewError("err_payment_order_has_closed", "订单已经取消")

	ErrOrderRefunded = domain.NewError("err_payment_order_has_refunded", "订单已退款")

	ErrOrderNotPayed = domain.NewError("err_payment_order_not_payed", "订单未支付")

	ErrCanNotUseBalance = domain.NewError("err_can_not_use_balance", "不能使用余额支付")

	ErrNotEnoughAmount = domain.NewError("err_payment_not_enough_amount", "余额不足,无法完成支付")

	ErrCanNotUseIntegral = domain.NewError("err_can_not_use_integral", "不能使用积分抵扣")

	ErrCanNotUseCoupon = domain.NewError("err_can_not_use_coupon", "不能使用优惠券")

	ErrCanNotSystemDiscount = domain.NewError("err_can_not_system_discount", "不允许系统支付")

	ErrOuterNo = domain.NewError("err_outer_no", "第三方交易号错误")
)

type (
	// IPaymentOrder 支付单接口
	IPaymentOrder interface {
		// GetAggregateRootId 获取聚合根编号
		GetAggregateRootId() int
		// Get 获取支付单的值
		Get() Order
		// TradeNo 获取交易号
		TradeNo() string
		// State 支付单状态
		State() int
		// Flag 支付方式
		Flag() int
		// TradeMethods 支付途径支付信息
		TradeMethods() []*TradeMethodData
		// CheckPaymentState 在支付之前检查订单状态
		CheckPaymentState() error
		// Submit 提交支付单
		Submit() error
		// MergePay 合并支付
		MergePay(orders []IPaymentOrder) (mergeTradeNo string, finalAmount int, err error)
		// Cancel 取消支付/退款
		Cancel() error
		// OfflineDiscount 线下现金/刷卡支付,cash:现金,bank:刷卡金额,finalZero:是否金额必须为零
		OfflineDiscount(cash int, bank int, finalZero bool) error
		// TradeFinish 交易完成
		TradeFinish() error
		// PaymentFinish 支付完成并保存,传入第三名支付名称,以及外部的交易号
		PaymentFinish(spName string, outTradeNo string) error
		// CouponDiscount 优惠券抵扣
		CouponDiscount(coupon promotion.ICouponPromotion) (int, error)
		// BalanceDeduct 使用会员的余额抵扣
		BalanceDeduct(remark string) error
		// WalletDeduct 使用会员的钱包抵扣
		WalletDeduct(remark string) error
		// IntegralDiscount 使用会员积分抵扣,返回抵扣的金额及错误,ignoreOut:是否忽略超出订单金额的积分
		IntegralDiscount(integral int, ignoreOut bool) (amount int, err error)
		// SystemPayment SystemPayment 系统支付金额
		SystemPayment(amount int) error
		// PaymentByWallet PaymentByWallet 钱包账户支付
		PaymentByWallet(remark string) error
		// PaymentWithCard 使用会员卡支付,cardCode:会员卡编码,amount:支付金额
		PaymentWithCard(cardCode string, amount int) error
		// HybridPayment 余额钱包混合支付，优先扣除余额。
		HybridPayment(remark string) error

		// Adjust 调整金额,如调整金额与实付金额相加小于等于零,则支付成功。
		Adjust(amount int) error
		// Refund 退款
		Refund(amount int) error
		// ChanName 获取支付通道字符串
		ChanName(method int) string
	}

	// IPaymentRepo 支付仓储
	IPaymentRepo interface {
		// GetPaymentOrderById 根据编号获取支付单
		GetPaymentOrderById(id int) IPaymentOrder
		// DeletePaymentOrder 拆分后删除父支付单
		DeletePaymentOrder(id int) error
		// DeletePaymentTradeData 删除支付单的支付数据
		DeletePaymentTradeData(orderId int)error
		// GetPaymentOrder 根据支付单号获取支付单
		GetPaymentOrder(paymenOrderNo string) IPaymentOrder
		// GetPaymentBySalesOrderId 根据订单号获取支付单
		GetPaymentBySalesOrderId(orderId int64) IPaymentOrder
		// GetPaymentOrderByOrderNo 根据支付单号获取支付单
		GetPaymentOrderByOrderNo(orderType int, orderNo string) IPaymentOrder
		// CreatePaymentOrder 创建支付单
		CreatePaymentOrder(p *Order) IPaymentOrder
		// SavePaymentOrder 保存支付单
		SavePaymentOrder(v *Order) (int, error)
		// CheckTradeNoMatch 检查支付单号是否匹配
		CheckTradeNoMatch(tradeNo string, id int) bool
		// GetTradeChannelItems 获取交易途径支付信息
		GetTradeChannelItems(tradeNo string) []*TradeMethodData
		// SavePaymentTradeChan 保存支付途径支付信息
		SavePaymentTradeChan(tradeNo string, tradeChan *TradeMethodData) (int, error)
		// GetMergePayOrders 获取合并支付的订单
		GetMergePayOrders(mergeTradeNo string) []IPaymentOrder
		// ResetMergePaymentOrders 清除欲合并的支付单
		ResetMergePaymentOrders(tradeNos []string) error
		// SaveMergePaymentOrders 保存合并的支付单
		SaveMergePaymentOrders(s string, tradeNos []string) error
		// FindAllIntegrateApp 集成支付应用
		FindAllIntegrateApp() []*IntegrateApp
		// SaveIntegrateApp Save 集成支付应用
		SaveIntegrateApp(v *IntegrateApp) (int, error)
		// DeleteIntegrateApp Delete 集成支付应用
		DeleteIntegrateApp(primary interface{}) error
		// GetAwaitCloseOrders 获取支付超时待关闭的订单
		GetAwaitCloseOrders(lastId int, size int) []IPaymentOrder
	}

	// Order 支付单
	Order struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 卖家编号
		SellerId int `db:"seller_id"`
		// 交易类型
		TradeType string `db:"trade_type"`
		// 交易号
		TradeNo string `db:"trade_no"`
		// 支付单的类型，如购物或其他
		OrderType int `db:"order_type"`
		// 是否为子订单
		SubOrder int `db:"sub_order"`
		// 外部订单号
		OutOrderNo string `db:"out_order_no"`
		// 支付单详情
		Subject string `db:"subject"`
		// 买家编号
		BuyerId int64 `db:"buyer_id"`
		// 支付用户编号
		PayerId int64 `db:"payer_id"`
		// 优惠金额,todo: 删除但目前依赖于优惠券
		DiscountAmount int64 `db:"-"`
		// 调整金额
		AdjustAmount int64 `db:"adjust_amount"`
		// 共计金额，包含抵扣金额
		TotalAmount int64 `db:"total_amount"`
		// 抵扣金额
		DeductAmount int64 `db:"deduct_amount"`
		// 手续费
		ProcedureFee int64 `db:"procedure_fee"`
		// 最终支付金额，包含手续费，不包含抵扣金额
		FinalAmount int64 `db:"final_amount"`
		// 实付金额
		PaidAmount int64 `db:"paid_amount"`
		// 可⽤支付方式
		PayFlag int `db:"pay_flag"`
		// 实际支付方式
		FinalFlag int `db:"final_flag"`
		// 其他支付信息
		ExtraData string `db:"extra_data"`
		// 交易支付渠道
		TradeChannel int `db:"trade_channel"`
		// 外部交易提供商
		OutTradeSp string `db:"out_trade_sp"`
		// 外部交易订单号
		OutTradeNo string `db:"out_trade_no"`
		// 订单状态
		State int `db:"state"`
		// 提交时间
		SubmitTime int64 `db:"submit_time"`
		// 支付过期时间
		ExpiresTime int64 `db:"expires_time"`
		// 支付时间
		PaidTime int64 `db:"paid_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
		// 交易途径支付信息
		TradeMethods []*TradeMethodData `db:"-"`
	}

	// TradeMethodData 支付单项
	TradeMethodData struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 交易单号
		TradeNo string `db:"trade_no"`
		// 支付途径
		Method int `db:"pay_method"`
		// 支付代码
		Code string `db:"pay_code"`
		// 是否为内置支付途径
		Internal int `db:"internal"`
		// 支付金额
		Amount int64 `db:"pay_amount"`
		// 外部交易单号
		OutTradeNo string `db:"out_trade_no"`
		// 支付时间
		PayTime int64 `db:"pay_time"`
	}

	// MergeOrder 合并的支付单
	MergeOrder struct {
		// 编号
		ID int `db:"id" pk:"yes" auto:"yes"`
		// 合并交易单号
		MergeTradeNo string `db:"merge_trade_no"`
		// 交易号
		OrderTradeNo string `db:"order_trade_no"`
		// 提交时间
		SubmitTime int64 `db:"submit_time"`
	}


	// IntegrateApp 集成支付应用
	IntegrateApp struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 支付应用名称
		AppName string `db:"app_name"`
		// 支付应用接口
		AppUrl string `db:"app_url"`
		// 集成方式: 1:API调用 2: 跳转
		IntegrateType int `db:"integrate_type"`
		// 显示顺序
		SortNumber int `db:"sort_number"`
		// 是否启用
		Enabled int `db:"enabled"`
		// 支付提示信息
		Hint string `db:"hint"`
		// 是否高亮显示
		Highlight int `db:"highlight"`
	}
)
