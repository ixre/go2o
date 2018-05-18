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

// 支付通道
const (
	// 余额抵扣通道
	ChanBalance = 1
	// 钱包支付通道
	ChanWallet = 2
	// 积分兑换通道
	ChanIntegral = 3
	// 用户卡通道
	ChanUserCard = 4
	// 用户券通道
	ChanUserCoupon = 5
	// 现金支付通道
	ChanCash = 6
	// 银行卡支付通道
	ChanBankCard = 7
	// 卖家支付通道
	ChanSellerPay = 8
	// 系统支付通道
	ChanSystemPay = 9
)

// 支付标志
const (
	// 余额抵扣
	FlagBalance = 1 << iota
	// 钱包支付
	FlagWallet
	// 积分兑换
	FlagIntegral
	// 用户卡
	FlagUserCard
	// 用户券
	FlagUserCoupon
	// 现金支付
	FlagCash
	// 银行卡支付
	FlagBankCard
	// 第三方支付,如支付宝等
	FlagOutSp
	// 卖家支付通道
	FlagSellerPay
	// 系统支付通道
	FlagSystemPay
)

// 所有支付方式
const PAllFlag = FlagBalance | FlagWallet | FlagIntegral |
	FlagCash | FlagBankCard | FlagOutSp | FlagSellerPay | FlagSystemPay

// 支付单状态
const (
	// 待支付
	StateAwaitingPayment = 1
	// 已支付
	StateFinished = 2
	// 已取消
	StateCancelled = 3
	// 已终止（超时关闭）
	StateAborted = 4
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

	ErrNotSupportPaymentChannel = domain.NewDomainError(
		"err_payment_not_support_channel", "不支持此支付方式,无法完成付款")
	ErrItemAmount    = domain.NewDomainError("err_payment_item_amount", "支付单金额不能为零")
	ErrOutOfFinalFee = domain.NewDomainError("err_out_of_final_fee",
		"超出支付单金额")
	ErrNotMatchFinalFee = domain.NewDomainError("err_not_match_final_fee",
		"金额与实际金额不符，无法完成付款")
	ErrTradeNoPrefix = domain.NewDomainError(
		"err_payment_trade_no_prefix", "支付单号前缀不正确")
	ErrTradeNoExistsPrefix = domain.NewDomainError(
		"err_payment_trade_no_exists_prefix", "支付单号已存在前缀")

	ErrOrderCommitted = domain.NewDomainError(
		"err_payment_order_committed", "支付单已提交")

	ErrOrderPayed = domain.NewDomainError(
		"err_payment_order_payed", "订单已支付")

	ErrOrderCancelled = domain.NewDomainError("err_payment_order_has_cancel", "订单已经取消")

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
		GetAggregateRootId() int
		// 获取支付单的值
		Get() Order
		// 获取交易号
		TradeNo() string
		// 支付单状态
		State() int
		/** 支付方式 */
		Flag() int
		// 支付途径支付信息
		Channels() []*TradeChan
		// 在支付之前检查订单状态
		CheckPaymentState() error
		// 提交支付单
		Submit() error
		// 合并支付
		MergePay(orders []IPaymentOrder) (mergeTradeNo string, finalFee int, err error)
		// 取消支付
		Cancel() error
		// 线下现金/刷卡支付,cash:现金,bank:刷卡金额,finalZero:是否金额必须为零
		OfflineDiscount(cash int, bank int, finalZero bool) error
		// 交易完成
		TradeFinish() error
		// 支付完成并保存,传入第三名支付名称,以及外部的交易号
		PaymentFinish(spName string, outTradeNo string) error
		// 优惠券抵扣
		CouponDiscount(coupon promotion.ICouponPromotion) (int, error)
		// 使用会员的余额抵扣
		BalanceDiscount(remark string) error
		// 使用会员积分抵扣,返回抵扣的金额及错误,ignoreOut:是否忽略超出订单金额的积分
		IntegralDiscount(integral int, ignoreOut bool) (amount int, err error)
		// 系统支付金额
		SystemPayment(amount int) error
		// 钱包账户支付
		PaymentByWallet(remark string) error
		// 余额钱包混合支付，优先扣除余额。
		HybridPayment(remark string) error
		// 设置支付方式
		SetTradeSP(spName string) error

		// 调整金额,如调整金额与实付金额相加小于等于零,则支付成功。
		Adjust(amount int) error
		// 退款
		Refund(amount int) error
	}

	// 支付仓储
	IPaymentRepo interface {
		// 根据编号获取支付单
		GetPaymentOrderById(id int) IPaymentOrder
		// 根据支付单号获取支付单
		GetPaymentOrder(tradeNo string) IPaymentOrder
		// 根据订单号获取支付单
		GetPaymentBySalesOrderId(orderId int64) IPaymentOrder
		// 根据支付单号获取支付单
		GetPaymentOrderByOrderNo(orderType int, orderNo string) IPaymentOrder
		// 创建支付单
		CreatePaymentOrder(p *Order) IPaymentOrder
		// 保存支付单
		SavePaymentOrder(v *Order) (int, error)
		// 检查支付单号是否匹配
		CheckTradeNoMatch(tradeNo string, id int) bool
		// 获取交易途径支付信息
		GetTradeChannelItems(tradeNo string) []*TradeChan
		// 保存支付途径支付信息
		SavePaymentTradeChan(tradeNo string, tradeChan *TradeChan) (int, error)
		// 获取合并支付的订单
		GetMergePayOrders(mergeTradeNo string) []IPaymentOrder
		// 清除欲合并的支付单
		ResetMergePaymentOrders(tradeNos []string) error
		//  保存合并的支付单
		SaveMergePaymentOrders(s string, tradeNos []string) error
	}

	// 支付通道
	PayChannel struct {
		// 编号
		ID int `db:"id" pk:"yes" auto:"yes"`
		// 支付渠道编码
		Code string `db:"code"`
		// 支付渠道名称
		Name int `db:"name"`
		// 支付渠道门户地址
		PortalUrl string `db:"portal_url"`
	}

	// 支付单
	Order struct {
		// 编号
		ID int `db:"id" pk:"yes" auto:"yes"`
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
		PayUid int64 `db:"pay_uid"`
		// 商品金额
		ItemAmount int `db:"item_amount"`
		// 优惠金额
		DiscountAmount int `db:"discount_amount"`
		// 调整金额
		AdjustAmount int `db:"adjust_amount"`
		// 共计金额，包含抵扣金额
		TotalAmount int `db:"total_amount"`
		// 抵扣金额
		DeductAmount int `db:"deduct_amount"`
		// 手续费
		ProcedureFee int `db:"procedure_fee"`
		// 最终支付金额，包含手续费，不包含抵扣金额
		FinalFee int `db:"final_fee"`
		// 可⽤支付方式
		PaymentFlag int `db:"pay_flag"`
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
		// 过期时间
		ExpiresTime int64 `db:"expires_time"`
		// 支付时间
		PaidTime int64 `db:"paid_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
		// 交易途径支付信息
		TradeChannels []*TradeChan `db:"-"`
	}

	// 支付单项
	TradeChan struct {
		// 编号
		ID int `db:"id" pk:"yes" auto:"yes"`
		// 交易单号
		TradeNo string `db:"trade_no"`
		// 支付途径
		PayChan int `db:"pay_chan"`
		// 是否为内置支付途径
		InternalChan int `db:"internal_chan"`
		// 支付金额
		PayAmount int `db:"pay_amount"`
	}

	// 合并的支付单
	MergeOrder struct {
		// 编号
		ID int `db:"id"`
		// 合并交易单号
		MergeTradeNo string `db:"merge_trade_no"`
		// 交易号
		OrderTradeNo string `db:"order_trade_no"`
		// 提交时间
		SubmitTime int64 `db:"submit_time"`
	}

	// SP支付交易
	SpTrade struct {
		// 编号
		ID int `db:"id"`
		// 交易SP
		TradeSp string `db:"trade_sp"`
		// 交易号
		TradeNo string `db:"trade_no"`
		// 合并的订单号,交易号用"|"分割
		TradeOrders string `db:"trade_orders"`
		// 交易状态
		TradeState int `db:"trade_state"`
		// 交易结果
		TradeResult int `db:"trade_result"`
		// 交易备注
		TradeRemark string `db:"trade_remark"`
		// 交易时间
		TradeTime int `db:"trade_time"`
	}
)
