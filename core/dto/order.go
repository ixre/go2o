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
	       o.package_fee,o.final_fee,o.status
	*/
	// 会员分页子订单
	PagedMemberSubOrder struct {
		Id             int64
		OrderNo        string
		ParentNo       string
		VendorId       int64
		ShopId         int64
		ShopName       string
		ItemAmount     int64
		DiscountAmount int64
		ExpressFee     int64
		PackageFee     int64
		IsPaid         bool
		FinalAmount    int64
		State          int
		StateText      string
		CreateTime     int64
		Items          []*OrderItem
	}

	PagedVendorOrder struct {
		Id        int64
		OrderNo   string
		ParentNo  string
		BuyerId   int
		BuyerName string
		// 订单详情,主要描述订单的内容
		Details string
		//VendorId    int
		//ShopId      int
		//ShopName    string
		ItemAmount     int64
		DiscountAmount int64
		ExpressFee     int64
		PackageFee     int64
		IsPaid         bool
		FinalAmount    int64
		State          int
		StateText      string
		CreateTime     int64
		Items          []*OrderItem
		Data           map[string]string
	}

	/*
	   SELECT si.id,si.order_id,si.snap_id,sn.sku_id,sn.goods_title,sn.img,
	           si.quantity,si.fee,si.final_fee
	*/

	// 订单商品项
	OrderItem struct {
		// 编号
		Id int
		// 订单编号
		OrderId int64
		// 商品快照编号
		SnapshotId int
		// Sku编号
		SkuId int
		// 商品编号
		ItemId int32
		// 商品标题
		GoodsTitle string
		// 商品图片
		Image string
		// 商品单价
		Price int64
		// 商品实际单价
		FinalPrice int64
		// 商品数量
		Quantity int
		// 退货数量
		ReturnQuantity int
		// 商品总金额
		Amount int64
		// 商品实际总金额
		FinalAmount int64
		// 是否已发货
		IsShipped int
	}
)
