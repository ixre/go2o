/**
 * Copyright 2015 @ z3q.net.
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
	       vendor_id,o.shop_id,s.name as shop_name,
	       o.goods_amount,o.discount_amount,o.express_fee,
	       o.package_fee,o.final_fee,o.status
	*/
	// 会员分页子订单
	PagedMemberSubOrder struct {
		Id             int
		OrderNo        string
		ParentNo       string
		VendorId       int
		ShopId         int
		ShopName       string
		GoodsAmount    float32
		DiscountAmount float32
		ExpressFee     float32
		PackageFee     float32
		IsPaid         bool
		FinalAmount    float32
		State          int
		StateText      string
		CreateTime     int64
		Items          []*OrderItem
	}

	PagedVendorOrder struct {
		Id        int
		OrderNo   string
		ParentNo  string
		BuyerId   int
		BuyerName string

		//VendorId    int
		//ShopId      int
		//ShopName    string
		GoodsAmount    float32
		DiscountAmount float32
		ExpressFee     float32
		PackageFee     float32
		IsPaid         bool
		FinalAmount    float32
		State          int
		StateText      string
		CreateTime     int64
		Items          []*OrderItem
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
		OrderId int
		// 商品快照编号
		SnapshotId int
		// 商品SKU编号
		SkuId int
		// 商品标题
		GoodsTitle string
		// 商品图片
		Image string
		// 商品单价
		Price float32
		// 商品实际单价
		FinalPrice float32
		// 商品数量
		Quantity int
		// 商品总金额
		Amount float32
		// 商品实际总金额
		FinalAmount float32
	}
)
