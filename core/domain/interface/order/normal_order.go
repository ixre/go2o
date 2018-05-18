package order

import (
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/promotion"
)

type (
	// 普通订单
	INormalOrder interface {
		// 读取购物车数据,用于预生成订单
		RequireCart(c cart.ICart) error
		// 根据运营商获取商品和运费信息,限未生成的订单
		GetByVendor() (items map[int32][]*SubOrderItem, expressFee map[int32]float32)
		// 在线支付交易完成
		OnlinePaymentTradeFinish() error
		// 设置配送地址
		SetAddress(addressId int64) error
		// 提交订单。如遇拆单,需均摊优惠抵扣金额到商品
		Submit() error

		//根据运营商拆单,返回拆单结果,及拆分的订单数组
		//BreakUpByVendor() ([]IOrder, error)

		// 获取子订单列表
		GetSubOrders() []ISubOrder
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
	}

	// 子订单(普通订单拆分)
	ISubOrder interface {
		// 获取领域对象编号
		GetDomainId() int64
		// 获取值对象
		GetValue() *NormalSubOrder
		// 复合的订单信息
		Complex() *ComplexOrder

		// 获取商品项
		Items() []*SubOrderItem
		// 在线支付交易完成
		PaymentFinishByOnlineTrade() error
		// 记录订单日志
		AppendLog(logType LogType, system bool, message string) error
		// 添加备注
		AddRemark(string)
		// 确认订单
		Confirm() error
		// 捡货(备货)
		PickUp() error
		// 发货
		Ship(spId int32, spOrder string) error
		// 已收货
		BuyerReceived() error
		// 获取订单的日志
		LogBytes() []byte
		// 挂起
		Suspend(reason string) error
		// 取消订单/退款
		Cancel(reason string) error
		// 退回商品
		Return(snapshotId int64, quantity int32) error
		// 撤销退回商品
		RevertReturn(snapshotId int64, quantity int32) error
		// 谢绝订单
		Decline(reason string) error
		// 提交子订单
		Submit() (int64, error)
		// 获取支付单
		GetPaymentOrder() payment.IPaymentOrder
	}

	// 普通订单
	NormalOrder struct {
		// 编号
		ID int64 `db:"id" pk:"yes" auto:"yes"`
		// 订单编号
		OrderId int64 `db:"order_id"`
		// 商品金额
		ItemAmount float32 `db:"item_amount"`
		// 优惠减免金额
		DiscountAmount float32 `db:"discount_amount" json:"discountFee"`
		// 运费
		ExpressFee float32 `db:"express_fee"`
		// 包装费用
		PackageFee float32 `db:"package_fee"`
		// 实际金额
		FinalAmount float32 `db:"final_amount" json:"fee"`
		// 收货人
		ConsigneePerson string `db:"consignee_person" json:"deliverName"`
		// 收货人联系电话
		ConsigneePhone string `db:"consignee_phone" json:"deliverPhone"`
		// 收货地址
		ShippingAddress string `db:"shipping_address" json:"deliverAddress"`
		// 订单是否拆分
		IsBreak int32 `db:"is_break"`
		// 更新时间
		UpdateTime int64 `db:"update_time" json:"updateTime"`
	}

	// 子订单
	NormalSubOrder struct {
		// 编号
		ID int64 `db:"id" pk:"yes" auto:"yes"`
		// 订单号
		OrderNo string `db:"order_no"`
		// 订单编号
		OrderId int64 `db:"order_id"`
		// 购买人编号(冗余,便于商户处理数据)
		BuyerId int64 `db:"buyer_id"`
		// 运营商编号
		VendorId int32 `db:"vendor_id" json:"vendorId"`
		// 店铺编号
		ShopId int32 `db:"shop_id" json:"shopId"`
		// 订单标题
		Subject string `db:"subject" json:"subject"`
		// 商品金额
		ItemAmount float32 `db:"item_amount"`
		// 优惠减免金额
		DiscountAmount float32 `db:"discount_amount" json:"discountFee"`
		// 运费
		ExpressFee float32 `db:"express_fee"`
		// 包装费用
		PackageFee float32 `db:"package_fee"`
		// 实际金额
		FinalAmount float32 `db:"final_amount" json:"fee"`
		// 是否支付
		IsPaid int `db:"is_paid"`
		// 是否挂起，如遇到无法自动进行的时挂起，来提示人工确认。
		IsSuspend int `db:"is_suspend" json:"is_suspend"`
		// 顾客备注
		BuyerComment string `db:"buyer_comment"`
		// 系统备注
		Remark string `db:"remark" json:"remark"`
		// 订单状态
		State int32 `db:"state" json:"state"`
		// 下单时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time" json:"updateTime"`
		// 订单项
		Items []*SubOrderItem `db:"-"`
	}

	// 订单商品项
	SubOrderItem struct {
		// 编号
		ID int64 `db:"id" pk:"yes" auto:"yes" json:"id"`
		// 订单编号
		OrderId int64 `db:"order_id"`
		// 商品编号
		ItemId int64 `db:"item_id"`
		// 商品SKU编号
		SkuId int64 `db:"sku_id"`
		// 快照编号
		SnapshotId int64 `db:"snap_id"`
		// 数量
		Quantity int32 `db:"quantity"`
		// 退回数量(退货)
		ReturnQuantity int32 `db:"return_quantity"`
		// 金额
		Amount float32 `db:"amount"`
		// 最终金额, 可能会有优惠均摊抵扣的金额
		FinalAmount float32 `db:"final_amount"`
		// 是否发货
		IsShipped int32 `db:"is_shipped"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
		// 运营商编号
		VendorId int32 `db:"-"`
		// 商店编号
		ShopId int32 `db:"-"`
		// 重量,用于生成订单时存储数据
		Weight int32 `db:"-"`
		// 体积:毫升(ml)
		Bulk int32 `db:"-"`
		// 快递模板编号
		ExpressTplId int32 `db:"-"`
	}
)
