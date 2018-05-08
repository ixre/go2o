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
	       o.vendor_id,o.shop_id,s.name as shop_name,
	       o.item_amount,o.discount_amount,o.express_fee,
	       o.package_fee,o.final_fee,o.status
	*/
	// 会员分页子订单
	PagedMemberSubOrder struct {
		Id             int64
		OrderNo        string
		ParentNo       string
		VendorId       int
		ShopId         int
		ShopName       string
		ItemAmount     float32
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
		Id        int64
		OrderNo   string
		ParentNo  string
		BuyerId   int
		BuyerName string

		//VendorId    int
		//ShopId      int
		//ShopName    string
		ItemAmount     float32
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
		Price float32
		// 商品实际单价
		FinalPrice float32
		// 商品数量
		Quantity int
		// 退货数量
		ReturnQuantity int
		// 商品总金额
		Amount float32
		// 商品实际总金额
		FinalAmount float32
		// 是否已发货
		IsShipped int
	}
)

/*
func ParseSubOrder(v *order.SubOrder) *SubOrder {
    return &SubOrder{
        ID:v.ID,
        OrderNo:v.OrderNo,
        ParentId:v.ParentId,
        BuyerId:v.BuyerId,
        VendorId:v.VendorId,
        ShopId:v.ShopId,
        Subject:v.Subject,
        ItemsInfo:v.ItemsInfo,
        ItemAmount:v.ItemAmount,
        DiscountAmount:v.DiscountAmount,
        ExpressFee:v.ExpressFee,
        PackageFee:v.PackageFee,
        FinalFee:v.FinalFee,
        IsPaid:v.IsPaid,
        IsSuspend:v.IsSuspend,
        Note:v.BuyerComment,
        Remark:v.Remark,
        UpdateTime:v.UpdateTime,
        State:v.State,
        Items:v.Items,
        StateText:order.OrderState(v.State).String(),
    }
}*/
