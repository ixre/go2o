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
	// ComplexOrder 订单复合信息
	ComplexOrder struct {
		// 订单编号
		OrderId int64
		// 订单类型
		OrderType int32
		// 订单号
		OrderNo string
		// 购买人编号
		BuyerId int64
		// 买家用户名
		BuyerUser string
		// 订单标题
		Subject string
		// 商品数量
		ItemCount int
		// 商品金额
		ItemAmount int64
		// 优惠减免金额
		DiscountAmount int64
		// 运费
		ExpressFee int64
		// 包装费用
		PackageFee int64
		// 实际金额
		FinalAmount int64
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
		// 扩展数据
		Data map[string]string
		// 订单详情
		Details []*ComplexOrderDetails
	}

	// ComplexOrderDetails 订单详情
	ComplexOrderDetails struct {
		// 编号
		Id int64
		// 订单号
		OrderNo string
		// 店铺编号
		ShopId int64
		// 店铺名称
		ShopName string
		// 商品金额
		ItemAmount int64
		// 优惠减免金额
		DiscountAmount int64
		// 运费
		ExpressFee int64
		// 包装费用
		PackageFee int64
		// 实际金额
		FinalAmount int64
		// 买家留言
		BuyerComment string
		// 订单状态
		State int32
		// 状态文本
		StateText string
		// 商品项
		Items []*ComplexItem
		// 更新时间
		UpdateTime int64
	}

	// ComplexConsignee 收货人信息
	ComplexConsignee struct {
		// 收货人
		ConsigneeName string
		// 收货人联系电话
		ConsigneePhone string
		// 收货地址
		ShippingAddress string
	}

	// ComplexItem 符合的订单项
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
		Price int64 `db:"-"`
		// 商品实际单价
		FinalPrice int64 `db:"-"`
		// 数量
		Quantity int32 `db:"quantity"`
		// 退回数量(退货)
		ReturnQuantity int32 `db:"return_quantity"`
		// 金额
		Amount int64 `db:"amount"`
		// 最终金额, 可能会有优惠均摊抵扣的金额
		FinalAmount int64 `db:"final_amount"`
		// 是否发货
		IsShipped int32 `db:"is_shipped"`
		// 其他数据
		Data map[string]string
	}
)
