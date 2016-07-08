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
	       o.goods_fee,o.discount_fee,o.express_fee,
	       o.package_fee,o.final_fee,o.status
	*/
	// 会员分页子订单
	PagedMemberSubOrder struct {
		Id          int
		OrderNo     string
		ParentNo    string
		VendorId    int
		ShopId      int
		ShopName    string
		GoodsFee    float32
		DiscountFee float32
		ExpressFee  float32
		PackageFee  float32
		IsPaid      bool
		FinalFee    float32
		Status      int
		CreateTime  int64
		Items       []*OrderItem
	}

	/*
	   SELECT si.id,si.order_id,si.snap_id,sn.sku_id,sn.goods_title,sn.img,
	           si.quantity,si.fee,si.final_fee
	*/

	// 订单商品项
	OrderItem struct {
		Id         int
		OrderId    int
		SnapshotId int
		SkuId      int
		GoodsTitle string
		Image      string
		Quantity   int
		Fee        float32
		FinalFee   float32
	}
)
