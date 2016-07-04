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
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/promotion"
	"go2o/core/infrastructure/domain"
)

// 自动拆单应在下单前完成
// 用户拆单,则重新生成子订单
// 参考:
//http://www.pmcaff.com/discuss?id=1000000000138488
//http://www.zhihu.com/question/31640837

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

		// 应用余额支付
		UseBalanceDiscount()

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
		AppendLog(t enum.OrderLogType, system bool, message string) error

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
		// 订单总金额
		TotalFee float32 `db:"total_fee"`
		// 优惠减免金额
		DiscountFee float32 `db:"discount_fee" json:"discountFee"`
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
		Status int `db:"status" json:"status"`
	}

	// 子订单
	SubOrder struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 订单号
		OrderNo string `db:"order_no"`
		// 订单编号
		ParentId int `db:"parent_order"`
		// 运营商编号
		VendorId int `db:"vendor_id" json:"vendorId"`
		// 店铺编号
		ShopId int `db:"shop_id" json:"shopId"`
		// 订单标题
		Subject string `db:"subject" json:"subject"`
		// 订单详情
		ItemsInfo string `db:"items_info" json:"itemsInfo"`
		// 订单总金额
		TotalFee float32 `db:"total_fee"`
		// 优惠减免金额
		DiscountFee float32 `db:"discount_fee" json:"discountFee"`
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
		Status int `db:"status" json:"status"`
		// 订单项
		Items []*OrderItem `db:"-"`
	}

	// 订单商品项
	OrderItem struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes" json:"id"`
		// 订单编号
		OrderId int `db:"order_id"`
		// 运营商编号
		VendorId int `db:"vendor_id"`
		// 商店编号
		ShopId int `db:"shop_id"`
		// 商品SKU编号
		SkuId int `db:"sku_id"`
		// 快照编号
		SnapshotId int `db:"snap_id"`
		// 数量
		Quantity int `db:"quantity"`
		// SKU描述
		Sku string `db:"sku"`
		// 金额
		Fee float32 `db:"fee"`
		// 最终金额, 可能会有优惠均摊抵扣的金额
		FinalFee float32 `db:"final_fee"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	ValueOrder2 struct {
		Id       int    `db:"id" pk:"yes" auto:"yes" json:"id"`
		OrderNo  string `db:"order_no" json:"orderNo"`
		BuyerId  int    `db:"buyer_id" json:"memberId"`
		VendorId int    `db:"vendor_id" json:"vendorId"`
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
