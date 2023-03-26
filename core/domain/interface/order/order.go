/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:37
 * description :
 * history :
 */

package order

import (
	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/infrastructure/domain"
)

// 自动拆单应在下单前完成
// 用户拆单,则重新生成子订单
// 参考:
//http://www.pmcaff.com/discuss?id=1000000000138488
//http://www.zhihu.com/question/31640837

type OrderStatus int
type OrderType int32

const (
	// 零售订单(线上/线下)
	TRetail OrderType = 1
	// 零售子订单
	TRetailSubOrder OrderType = 9
	// 批发订单
	TWholesale OrderType = 2
	// 虚拟订单,如：手机充值
	TVirtual OrderType = 3
	// 交易订单,如：线下支付。
	TTrade OrderType = 4
	// 服务订单
	TService OrderType = 5
)

const (
	/****** 在履行前,订单可以取消申请退款  ******/

	// StatAwaitingPayment 等待支付
	StatAwaitingPayment = 1
	// StatAwaitingConfirm 等待确认
	StatAwaitingConfirm = 2
	// StatAwaitingPickup 等待备货
	StatAwaitingPickup = 3
	// StatAwaitingShipment 等待发货
	StatAwaitingShipment = 4

	/****** 订单取消 ******/

	// StatCancelled 系统取消
	StatCancelled = 11
	// StatAwaitingCancel 买家申请取消,等待卖家确认
	StatAwaitingCancel = 12
	// StatDeclined 卖家谢绝订单,由于无货等原因
	StatDeclined = 13
	// StatRefunded 已退款,完成取消
	StatRefunded = 14

	/****** 履行后订单只能退货或换货 ******/

	// PartiallyShipped 部分发货(将订单商品分多个包裹发货)
	PartiallyShipped = 5
	// StatShipped 完成发货
	StatShipped = 6
	// StatBreak 订单已拆分
	StatBreak = 7
	// StatCompleted 订单完成
	StatCompleted = 8

	/****** 售后状态 ******/

	// StatGoodsRefunded 已退货
	StatGoodsRefunded = 15
)

const (
	// BreakDefault 默认
	BreakDefault = 0
	// BreakAwaitBreak 待拆分
	BreakAwaitBreak = 1
	// BreakNoBreak 无需拆分
	BreakNoBreak = 2
	// Breaked 已拆分
	Breaked = 3
)

func (t OrderStatus) String() string {
	switch t {
	case StatAwaitingPayment:
		return "待付款"
	case StatAwaitingConfirm:
		return "待确认"
	case StatAwaitingPickup:
		return "正在备货"
	case StatAwaitingShipment:
		return "等待发货"
	case StatCancelled:
		return "交易关闭"
	case StatDeclined:
		return "卖家关闭"
	case StatAwaitingCancel:
		return "等待退款"
	case StatRefunded:
		return "已退款"
	case PartiallyShipped:
		return "已部分发货"
	case StatShipped:
		return "待收货"
	case StatCompleted:
		return "交易完成"
	case StatBreak:
		return "已拆单"
	case StatGoodsRefunded:
		return "已退货"
	}
	return "Error State"
}

// BackEndString 后端状态描述
func (t OrderStatus) BackEndString() string {
	return t.String()
}

const (
	LogSetup       LogType = 1
	LogChangePrice LogType = 2
)

type LogType int

func (o LogType) String() string {
	switch o {
	case LogSetup:
		return "流程"
	case LogChangePrice:
		return "调价"
	}
	return ""
}

var (
	ErrNoCheckedItem = domain.NewError(
		"err_order_no_checked_item", "没有可结算的商品")

	ErrRequireCart = domain.NewError(
		"err_require_cart", "订单已生成,无法引入购物车")

	ErrNoSuchOrder = domain.NewError(
		"err_no_such_order", "订单不存在")

	ErrOrderPayed = domain.NewError(
		"err_order_payed ", "订单已支付")

	ErrNoYetCreated = domain.NewError(
		"err_order_not_yet_created ", "订单尚未生成")

	ErrUnusualOrder = domain.NewError(
		"err_unusual_order", "订单异常")

	ErrMissingShipAddress = domain.NewError(
		"err_missing_ship_address", "未设置收货地址")

	ErrUnusualOrderStat = domain.NewError(
		"err_except_order_stat", "订单状态不匹配、无法执行此操作!")

	ErrPartialShipment = domain.NewError(
		"err_order_partial_shipment", "订单部分商品已经发货")

	ErrOrderNotPayed = domain.NewError(
		"err_order_not_payed ", "订单未支付")

	ErrOutOfQuantity = domain.NewError(
		"err_order_out_of_quantity", "超出数量")
	ErrNoSuchGoodsOfOrder = domain.NewError(
		"err_order_no_such_goods_of_order", "订单中不包括该商品")
	ErrOrderHasConfirm = domain.NewError(
		"err_order_has_confirm", "订单已经确认")

	ErrOrderNotConfirm = domain.NewError(
		"err_order_not_confirm", "请等待系统确认")

	ErrOrderHasPickUp = domain.NewError(
		"err_order_has_pick_up", "订单已经备货")

	ErrOrderHasShipment = domain.NewError(
		"err_order_has_shipment", "订单已经发货")

	ErrNoSuchAddress = domain.NewError(
		"err_order_no_address", "请选择收货地址")

	ErrOrderShipped = domain.NewError(
		"err_order_shipped", "订单已经发货")

	ErrOrderNotShipped = domain.NewError(
		"err_order_not_shipped", "订单尚未发货")

	ErrIsCompleted = domain.NewError(
		"err_order_is_completed", "订单已经完成")

	ErrOrderBreakUpFail = domain.NewError(
		"err_order_break_up_fail", "拆分订单操作失败")

	ErrPromotionApplied = domain.NewError(
		"err_promotion_applied", "已经使用相同的促销")

	ErrEmptyReason = domain.NewError(
		"err_order_empty_reason", "原因不能为空")

	ErrOrderCancelled = domain.NewError(
		"err_order_can_not_cancel", "订单已经取消")

	ErrOrderShippedCancel = domain.NewError(
		"err_order_shipped_cancel", "订单已发货，无法取消")

	ErrHasRefund = domain.NewError(
		"err_order_has_refund", "订单已经退款")

	ErrDisallowRefund = domain.NewError(
		"err_order_disallow_refund", "订单不允许退款")
	ErrDisallowCancel = domain.NewError(
		"err_order_disallow_cancel", "会员无法取消此订单")
	ErrTradeRateLessZero = domain.NewError(
		"err_order_trade_rate_less_zero", "交易类订单结算比例不能小于零")

	ErrTradeRateMoreThan100 = domain.NewError(
		"err_order_trade_rate_more_than_100", "交易类订单结算比例必须小于或等于100%")

	ErrMissingSubject = domain.NewError(
		"err_order_missing_subject", "缺少订单标题")

	ErrTicketImage = domain.NewError(
		"err_order_ticket_image", "请上传正确的发票凭证")

	ErrForbidStatus = domain.NewError(
		"err_order_forbid_status", "仅已取消或完成的订单才能删除")
)

type (
	IOrder interface {
		// GetAggregateRootId 获取编号
		GetAggregateRootId() int64
		// Type 订单类型
		Type() OrderType
		// State 获取订单状态
		State() OrderStatus
		// Buyer 获取购买的会员
		Buyer() member.IMember
		// SetShipmentAddress 设置配送地址
		SetShipmentAddress(addressId int64) error
		// 更改收货人信息
		ChangeShipmentAddress(addressId int64) error
		// OrderNo 获取订单号
		OrderNo() string
		// Complex 复合的订单信息
		Complex() *ComplexOrder
		// Submit 提交订单。如遇拆单,需均摊优惠抵扣金额到商品
		Submit() error
		// BuildCart 通过订单创建购物车 */
		BuildCart() cart.ICart
		// GetPaymentOrder 获取支付单
		GetPaymentOrder() payment.IPaymentOrder
	}

	// IWholesaleOrder 批发订单
	IWholesaleOrder interface {
		// SetItems 设置商品项
		SetItems(items []*cart.ItemPair)
		// SetComment 设置或添加买家留言，如已经提交订单，将在原留言后附加
		SetComment(comment string)
		// Items 获取商品项
		Items() []*WholesaleItem
		// OnlinePaymentTradeFinish 在线支付交易完成
		OnlinePaymentTradeFinish() error

		// AppendLog 记录订单日志
		AppendLog(logType LogType, system bool, message string) error
		// AddRemark 添加备注
		AddRemark(string)
		// Confirm 确认订单
		Confirm() error
		// PickUp 捡货(备货)
		PickUp() error
		// Ship 发货
		Ship(spId int32, spOrder string) error
		// BuyerReceived 已收货
		BuyerReceived() error
		// LogBytes 获取订单的日志
		LogBytes() []byte
		// Cancel 取消订单/退款
		Cancel(reason string) error
		// Decline 谢绝订单
		Decline(reason string) error
	}

	// ITradeOrder 交易订单
	ITradeOrder interface {
		// Set 从订单信息中拷贝相应的数据,并设置订单结算比例
		Set(o *TradeOrderValue, rate float64) error
		// CashPay 现金支付
		CashPay() error
		// TradePaymentFinish 交易支付完成
		TradePaymentFinish() error
		// UpdateTicket 更新发票数据
		UpdateTicket(img string) error
	}
)

type (
	// SubmitReturnData 订单提交返回数据
	SubmitReturnData struct {
		// 订单号，多个订单号，用","分割
		OrderNo string
		// 交易金额
		TradeAmount int64
		// 合并支付
		IsMergePay bool
		// 支付单号
		PaymentOrderNo string
		// 支付状态
		PaymentState int
	}

	// Order 订单
	Order struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes"`
		// 订单号
		OrderNo string `db:"order_no"`
		// 订单类型1:普通 2:批发 3:线下
		OrderType int `db:"order_type"`
		// 订单主题
		Subject string `db:"subject"`
		// 买家
		BuyerId int64 `db:"buyer_id"`
		// 买家用户名
		BuyerUser string `db:"buyer_user"`
		// 商品数量
		ItemCount int `db:"item_count"`
		// 商品金额
		ItemAmount int64 `db:"item_amount"`
		// 抵扣金额
		DiscountAmount int64 `db:"discount_amount"`
		// 物流费
		ExpressFee int64 `db:"express_fee"`
		// 包装费
		PackageFee int64 `db:"package_fee"`
		// 订单最终金额
		FinalAmount int64 `db:"final_amount"`
		// 收货人姓名
		ConsigneeName string `db:"consignee_name"`
		// 收货人电话
		ConsigneePhone string `db:"consignee_phone"`
		// 收货人地址
		ShippingAddress string `db:"shipping_address"`
		// 地址下单后是否修改
		ConsigneeModified int `db:"consignee_modified"`
		// 是否拆分
		IsBreak int `db:"is_break"`
		// 是否支付
		IsPaid int `db:"is_paid"`
		// 订单状态
		Status int `db:"status"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// MinifyItem 订单商品项
	MinifyItem struct {
		ItemId   int32
		SkuId    int32
		Quantity int32
	}

	// WholesaleOrder 批发订单
	WholesaleOrder struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes"`
		// 订单号
		OrderNo string `db:"order_no"`
		// 订单编号
		OrderId int64 `db:"order_id"`
		// 买家
		BuyerId int64 `db:"buyer_id"`
		// 供货商
		VendorId int64 `db:"vendor_id"`
		// 店铺编号
		ShopId int64 `db:"shop_id"`
		// 店铺名称
		ShopName string `db:"shop_name"`
		// 买家留言
		BuyerComment string `db:"buyer_comment"`
		// 备注
		Remark string `db:"remark"`
		// 订单状态
		Status int `db:"status"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// WholesaleItem 批发订单商品
	WholesaleItem struct {
		// 编号
		ID int64 `db:"id" pk:"yes" auto:"yes"`
		// 订单编号
		OrderId int64 `db:"order_id"`
		// 商品编号
		ItemId int64 `db:"item_id"`
		// SKU编号
		SkuId int64 `db:"sku_id"`
		// 商品快照编号
		SnapshotId int64 `db:"snapshot_id"`
		// 商品销售价格(不含优惠抵扣)
		Price int64 `db:"-"` //todo
		// 销售数量
		Quantity int32 `db:"quantity"`
		// 退货数量
		ReturnQuantity int32 `db:"return_quantity"`
		// 商品总金额
		Amount int64 `db:"amount"`
		// 商品实际金额
		FinalAmount int64 `db:"final_amount"`
		// 是否已发货
		IsShipped int32 `db:"is_shipped"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// TradeOrder 交易类订单
	TradeOrder struct {
		// 编号
		ID int64 `db:"id" pk:"yes" auto:"yes"`
		// 订单编号
		OrderId int64 `db:"order_id"`
		// 商家编号
		VendorId int64 `db:"vendor_id"`
		// 店铺编号
		ShopId int64 `db:"shop_id"`
		// 订单标题
		Subject string `db:"subject"`
		// 订单金额
		OrderAmount int64 `db:"order_amount"`
		// 抵扣金额
		DiscountAmount int64 `db:"discount_amount"`
		// 订单最终金额
		FinalAmount int64 `db:"final_amount"`
		// 交易结算比例（商户)，允许为0和1
		TradeRate float64 `db:"trade_rate"`
		// 是否现金支付
		CashPay int32 `db:"cash_pay"`
		// 发票图片
		TicketImage string `db:"ticket_image"`
		// 订单备注
		Remark string `db:"remark"`
		// 订单状态
		Status int `db:"status"`
		// 订单创建时间
		CreateTime int64 `db:"create_time"`
		// 订单更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// OrderLog 订单变动日志
	OrderLog struct {
		Id      int32 `db:"id" auto:"yes" pk:"yes"`
		OrderId int64 `db:"order_id"`
		Type    int   `db:"type"`
		// 订单状态
		OrderState int    `db:"order_state"`
		IsSystem   int    `db:"is_system"`
		Message    string `db:"message"`
		RecordTime int64  `db:"record_time"`
	}
	OrderPromotionBind struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 订单号
		OrderId int32 `db:"order_id"`
		// 促销编号
		PromotionId int32 `db:"promotion_id"`
		// 促销类型
		PromotionType int `db:"promotion_type"`
		// 标题
		Title string `db:"title"`
		// 节省金额
		SaveFee int64 `db:"save_fee"`
		// 赠送积分
		PresentIntegral int `db:"present_integral"`
		// 是否应用
		IsApply int `db:"is_apply"`
		// 是否确认
		IsConfirm int `db:"is_confirm"`
	}

	// OrderCoupon 应用到订单的优惠券
	OrderCoupon struct {
		OrderId      int32  `db:"order_id"`
		CouponId     int32  `db:"coupon_id"`
		CouponCode   string `db:"coupon_code"`
		Fee          int64  `db:"coupon_fee"`
		Describe     string `db:"coupon_describe"`
		SendIntegral int    `db:"send_integral"`
	}
)
