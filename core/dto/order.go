/**
 * Copyright 2015 @ 56x.net.
 * name : order
 * author : jarryliu
 * date : 2016-07-08 16:46
 * description :
 * history :
 */
package dto

type (
	/*
	   o.order_no,po.order_no as parent_no,
	       o.vendor_id,o.shop_id,s.shop_name,
	       o.item_amount,o.discount_amount,o.express_fee,
	       o.package_fee,o.final_amount,o.status
	*/

	// MemberPagingOrderDto 会员订单分页对象
	MemberPagingOrderDto struct {
		// 订单编号
		OrderId int64 `json:"orderId"`
		// 买家
		BuyerId int64 `json:"buyerId"`
		// 买家用户名
		BuyerUser string `json:"BuyerUser"`
		// 店铺编号
		ShopId int64 `json:"shopId"`
		// 店铺名称
		ShopName string `json:"shopName"`
		// 订单号
		OrderNo string `json:"orderNo"`
		// 商品数量
		ItemCount int `json:"itemCount"`
		// 商品总金额
		ItemAmount int64 `json:"itemAmount"`
		// 抵扣金额
		DiscountAmount int64 `json:"discountAmount"`
		// 优惠金额
		DeductAmount int64 `json:"deductAmount"`
		// 快递费
		ExpressFee int64 `json:"expressFee"`
		// 包装费
		PackageFee int64 `json:"packageFee"`
		// 最终金额
		FinalAmount int64 `json:"finalAmount"`
		// 是否支付
		IsPaid int32 `json:"isPaid"`
		// 状态
		Status int32 `json:"status"`
		// 状态文本
		StatusText string `json:"statusText"`
		// 下单时间
		CreateTime int64 `json:"createTime"`
		// 订单商品
		Items []*OrderItem `json:"items"`
	}

	PagedVendorOrder struct {
		Id        int64  `json:"id"`
		OrderNo   string `json:"orderNo"`
		ParentNo  string `json:"parentNo"`
		BuyerId   int    `json:"buyerId"`
		BuyerName string `json:"buyerName"`
		// 订单详情,主要描述订单的内容
		Details string `json:"details"`
		//VendorId    int
		//ShopId      int
		//ShopName    string
		ItemAmount     int64             `json:"itemAmount"`
		DiscountAmount int64             `json:"discountAmount"`
		ExpressFee     int64             `json:"expressFee"`
		PackageFee     int64             `json:"packageFee"`
		IsPaid         bool              `json:"isPaid"`
		FinalAmount    int64             `json:"finalAmount"`
		Status         int               `json:"status"`
		StatusText     string            `json:"statusText"`
		CreateTime     int64             `json:"createTime"`
		Items          []*OrderItem      `json:"items"`
		Data           map[string]string `json:"data"`
	}

	/*
	   SELECT si.id,si.order_id,si.snap_id,sn.sku_id,sn.goods_title,sn.img,
	           si.quantity,si.fee,si.final_amount
	*/

	// OrderItem 订单商品项
	OrderItem struct {
		// 编号
		Id int `json:"id"`
		// 订单编号
		OrderId int64 `json:"orderId"`
		// 商品快照编号
		SnapshotId int `json:"snapshotId"`
		// Sku规格
		SpecWord string `json:"itemSpec"`
		// Sku编号
		SkuId int `json:"skuId"`
		// 商品编号
		ItemId int32 `json:"itemId"`
		// 商品标题
		ItemTitle string `json:"itemTitle"`
		// 商品图片
		Image string `json:"image"`
		// 商品单价
		Price int64 `json:"price"`
		// 商品实际单价
		FinalPrice int64 `json:"finalPrice"`
		// 商品数量
		Quantity int `json:"quantity"`
		// 退货数量
		ReturnQuantity int `json:"returnQuantity"`
		// 商品总金额
		Amount int64 `json:"amount"`
		// 商品实际总金额
		FinalAmount int64 `json:"finalAmount"`
		// 是否已发货
		IsShipped int `json:"isShipped"`
	}
)
