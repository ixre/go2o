package order

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : complex_order
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-02 07:21
 * description :
 * history :
 */

type (
	// 订单复合信息
	ComplexOrder struct {
		// 订单编号
		OrderId int64
		// 子订单编号
		SubOrderId int64
		// 订单类型
		OrderType int32
		// 是否为子订单
		SubOrder bool
		// 订单号
		OrderNo string
		// 购买人编号
		BuyerId int64
		// 运营商编号
		VendorId int64
		// 店铺编号
		ShopId int64
		// 订单标题
		Subject string
		// 商品金额
		ItemAmount float64
		// 优惠减免金额
		DiscountAmount float64
		// 运费
		ExpressFee float64
		// 包装费用
		PackageFee float64
		// 实际金额
		FinalAmount float64
		// 买家留言
		BuyerComment string
		// 收货人信息
		Consignee *ComplexConsignee
		// 订单是否拆分
		IsBreak int32
		// 订单状态
		State int32
		// 状态文本
		StateText string
		// 订单生成时间
		CreateTime int64
		// 更新时间
		UpdateTime int64
		// 商品项
		Items []*ComplexItem
		// 扩展数据
		Data map[string]string
	}

	// 收货人信息
	ComplexConsignee struct {
		// 收货人
		ConsigneePerson string
		// 收货人联系电话
		ConsigneePhone string
		// 收货地址
		ShippingAddress string
	}

	// 符合的订单项
	ComplexItem struct {
		// 编号
		ID int64 `db:"id" pk:"yes" auto:"yes" json:"id"`
		// 商品编号
		ItemId int64 `db:"item_id"`
		// 商品SKU编号
		SkuId int64 `db:"sku_id"`
		// SKU名称
		SkuWord string `db:"-"`
		// 快照编号
		SnapshotId int64 `db:"snap_id"`
		// 商品标题
		ItemTitle string `db:"item_title"`
		// 商品图片
		MainImage string `db:"image"`
		// 商品单价
		Price float32 `db:"-"`
		// 商品实际单价
		FinalPrice float32 `db:"-"`
		// 数量
		Quantity int32 `db:"quantity"`
		// 退回数量(退货)
		ReturnQuantity int32 `db:"return_quantity"`
		// 金额
		Amount float64 `db:"amount"`
		// 最终金额, 可能会有优惠均摊抵扣的金额
		FinalAmount float64 `db:"final_amount"`
		// 是否发货
		IsShipped int32 `db:"is_shipped"`
		// 其他数据
		Data map[string]string
	}
)