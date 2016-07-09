/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:37
 * description :
 * history :
 */

package order

import (
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/promotion"
	"go2o/core/infrastructure/domain"
)

// 自动拆单应在下单前完成
// 用户拆单,则重新生成子订单
// 参考:
//http://www.pmcaff.com/discuss?id=1000000000138488
//http://www.zhihu.com/question/31640837

type OrderState int

const (
	/****** 在履行前,订单可以取消申请推狂  ******/

	// 等待支付
	StatAwaitingPayment = 1
	// 等待确认
	StatAwaitingConfirm = 2
	// 等待备货
	StatAwaitingPickup = 3
	// 等待发货
	StatAwaitingShipment = 4

	/****** 订单取消 ******/

	// 系统取消
	StatCancelled = 11
	// 买家申请取消,等待卖家确认
	StatAwaitingCancel = 12
	// 卖家谢绝订单,由于无货等原因
	StatDeclined = 13
	// 已退款,完成取消
	StatRefunded = 14

	/****** 履行后订单只能退货或换货 ******/

	// 部分发货(将订单商品分多个包裹发货)
	PartiallyShipped = 5
	// 完成发货
	StatShipped = 6
	// 订单完成
	StatCompleted = 7

	/****** 售后状态 ******/

	// 已退货
	StatGoodsRefunded = 15
)

func (t OrderState) String() string {
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
	case StatGoodsRefunded:
		return "已退货"
	}
	return "Error State"
}

// 后端状态描述
func (t OrderState) BackEndString() string {
	return t.String()
}

const (
	LogSetup       LogType = 1
	LogChangePrice LogType = 2
)

type LogType int

func (this LogType) String() string {
	switch this {
	case LogSetup:
		return "流程"
	case LogChangePrice:
		return "调价"
	}
	return ""
}

var (
	ErrRequireCart *domain.DomainError = domain.NewDomainError(
		"err_require_cart", "订单已生成,无法引入购物车")

	ErrNoSuchOrder *domain.DomainError = domain.NewDomainError(
		"err_no_such_order", "订单不存在")

	ErrOrderPayed *domain.DomainError = domain.NewDomainError(
		"err_order_payed ", "订单已支付")

	ErrOrderNotPayed *domain.DomainError = domain.NewDomainError(
		"err_order_not_payed ", "订单未支付")

	ErrOrderDelved *domain.DomainError = domain.NewDomainError(
		"err_order_delved ", "订单已发货")

	ErrOrderBreakUpFail *domain.DomainError = domain.NewDomainError(
		"err_order_break_up_fail", "拆分订单操作失败")

	ErrPromotionApplied *domain.DomainError = domain.NewDomainError(
		"err_promotion_applied", "已经使用相同的促销")
)

type (
	IOrder interface {
		// 获取聚合根编号
		GetAggregateRootId() int

		// 获取订单号
		GetOrderNo() string

		// 获生成值
		GetValue() *Order

		// 读取购物车数据,用于预生成订单
		RequireCart(c cart.ICart) error

		// 根据运营商获取商品和运费信息,限未生成的订单
		GetByVendor() (items map[int][]*OrderItem, expressFee map[int]float32)

		// 获取购买的会员
		GetBuyer() member.IMember

		// GetComplexValue()dto.OrderComplex
		// 设置订单值
		//SetValue(*Order) error

		// 应用优惠券
		ApplyCoupon(coupon promotion.ICouponPromotion) error

		// 获取应用的优惠券
		GetCoupons() []promotion.ICouponPromotion

		// 获取可用的促销,不包含优惠券
		GetAvailableOrderPromotions() []promotion.IPromotion

		// 获取最省的促销
		GetBestSavePromotion() (p promotion.IPromotion,
			saveFee float32, integral int)

		// 获取促销绑定
		GetPromotionBinds() []*OrderPromotionBind

		// 设置支付方式
		//SetPayment(payment int)

		// 使用余额支付
		PaymentWithBalance() error

		// 客服使用余额支付
		CmPaymentWithBalance() error

		// 在线交易支付
		PaymentForOnlineTrade(serverProvider string, tradeNo string) error

		// 设置配送地址
		SetDeliver(deliverAddressId int) error

		// 提交订单，返回订单号。如有错误则返回
		Submit() (string, error)

		// 保存订单, 在生成支付单后,应该根据实际支付金额
		// 进行拆单,并切均摊优惠抵扣金额
		Save() (int, error)

		//根据运营商拆单,返回拆单结果,及拆分的订单数组
		//BreakUpByVendor() ([]IOrder, error)

		// 添加日志,system表示为系统日志
		AppendLog(l *OrderLog) error

		// 订单是否结束
		IsOver() bool

		// 处理订单
		Process() error

		// 确认订单
		Confirm() error

		// 配送订单,并记录配送服务商编号及单号
		Deliver(spId int, spNo string) error

		// 获取支付金额
		//GetPaymentFee() float32

		//*********** 删除  *************//

		// 设置Shop,如果不需要记录日志，则remark传递空
		//SetShop(shopId int) error

		// 添加备注
		//AddRemark(string)

		// 挂起
		//Suspend(reason string) error

		// 标记收货
		//SignReceived() error

		// 完成订单
		//Complete() error

		// 取消订单
		//Cancel(reason string) error
	}

	ISubOrder interface {
		// 获取领域对象编号
		GetDomainId() int

		// 获取值对象
		GetValue() *SubOrder

		// 记录订单日志
		AppendLog(logType LogType, system bool, message string) error

		// 设置Shop,如果不需要记录日志，则remark传递空
		SetShop(shopId int) error

		// 添加备注
		AddRemark(string)

		// 挂起
		Suspend(reason string) error

		// 标记收货
		SignReceived() error

		// 完成订单
		Complete() error

		// 取消订单
		Cancel(reason string) error

		// 保存订单
		Save() (int, error)
	}

	// 简单商品信息
	OrderGoods struct {
		GoodsId    int    `json:"id"`
		GoodsImage string `json:"img"`
		Name       string `json:"name"`
		Quantity   int    `json:"qty"`
	}

	OrderLog struct {
		Id      int `db:"id" auto:"yes" pk:"yes"`
		OrderId int `db:"order_id"`
		Type    int `db:"type"`
		// 订单状态
		OrderState int    `db:"order_state"`
		IsSystem   int    `db:"is_system"`
		Message    string `db:"message"`
		RecordTime int64  `db:"record_time"`
	}
	OrderPromotionBind struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`

		// 订单号
		OrderId int `db:"order_id"`

		// 促销编号
		PromotionId int `db:"promotion_id"`

		// 促销类型
		PromotionType int `db:"promotion_type"`

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

	//todo: ??? 父订单的金额,是否可不用?

	// 订单
	Order struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 订单号
		OrderNo string `db:"order_no"`
		// 购买人编号
		BuyerId int `db:"buyer_id"`
		// 订单详情
		ItemsInfo string `db:"items_info" json:"itemsInfo"`
		// 商品金额
		GoodsFee float32 `db:"goods_fee"`
		// 优惠减免金额
		DiscountFee float32 `db:"discount_fee" json:"discountFee"`
		// 运费
		ExpressFee float32 `db:"express_fee"`
		// 包装费用
		PackageFee float32 `db:"package_fee"`
		// 实际金额
		FinalFee float32 `db:"final_fee" json:"fee"`
		// 是否支付
		IsPaid int `db:"is_paid"`
		// 支付时间
		PaidTime int64 `db:"paid_time" json:"paidTime"`
		// 收货人
		ConsigneePerson string `db:"consignee_person" json:"deliverName"`
		// 收货人联系电话
		ConsigneePhone string `db:"consignee_phone" json:"deliverPhone"`
		// 收货地址
		ShippingAddress string `db:"shipping_address" json:"deliverAddress"`
		// 发货时间
		ShippingTime int64 `db:"shipping_time" json:"deliverTime"`
		// 订单生成时间
		CreateTime int64 `db:"create_time" json:"createTime"`
		// 更新时间
		UpdateTime int64 `db:"update_time" json:"updateTime"`
		// 订单状态
		//todo: ???删除?
		State int `db:"status" json:"status"`
	}

	// 子订单
	SubOrder struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 订单号
		OrderNo string `db:"order_no"`
		// 订单编号
		ParentId int `db:"parent_order"`
		// 购买人编号(冗余,便于商户处理数据)
		BuyerId int `db:"buyer_id"`
		// 运营商编号
		VendorId int `db:"vendor_id" json:"vendorId"`
		// 店铺编号
		ShopId int `db:"shop_id" json:"shopId"`
		// 订单标题
		Subject string `db:"subject" json:"subject"`
		// 订单详情
		ItemsInfo string `db:"items_info" json:"itemsInfo"`
		// 商品金额
		GoodsFee float32 `db:"goods_fee"`
		// 优惠减免金额
		DiscountFee float32 `db:"discount_fee" json:"discountFee"`
		// 运费
		ExpressFee float32 `db:"express_fee"`
		// 包装费用
		PackageFee float32 `db:"package_fee"`
		// 实际金额
		FinalFee float32 `db:"final_fee" json:"fee"`
		// 是否挂起，如遇到无法自动进行的时挂起，来提示人工确认。
		IsSuspend int `db:"is_suspend" json:"is_suspend"`
		// 顾客备注
		Note string `db:"note" json:"note"`
		// 系统备注
		Remark string `db:"remark" json:"remark"`
		// 更新时间
		UpdateTime int64 `db:"update_time" json:"updateTime"`
		// 订单状态
		State int `db:"status" json:"status"`
		// 订单项
		Items []*OrderItem `db:"-"`
	}

	// 订单商品项
	OrderItem struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes" json:"id"`
		// 订单编号
		OrderId int `db:"order_id"`
		// 商品SKU编号
		SkuId int `db:"sku_id"`
		// 快照编号
		SnapshotId int `db:"snap_id"`
		// 数量
		Quantity int `db:"quantity"`
		// SKU描述
		//Sku string `db:"sku"`
		// 金额
		Fee float32 `db:"fee"`
		// 最终金额, 可能会有优惠均摊抵扣的金额
		FinalFee float32 `db:"final_fee"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`

		// 运营商编号
		VendorId int `db:"-"`
		// 商店编号
		ShopId int `db:"-"`
		// 重量,用于生成订单时存储数据
		Weight int `db:"-"`
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
