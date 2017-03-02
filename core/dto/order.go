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
	// 子订单
	SubOrder struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 订单号
		OrderNo string `db:"order_no"`
		// 订单编号
		ParentId int `db:"order_pid"`
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
		Note string `db:"note" json:"note"`
		// 系统备注
		Remark string `db:"remark" json:"remark"`
		// 更新时间
		UpdateTime int64 `db:"update_time" json:"updateTime"`
		// 订单状态
		State int `db:"state" json:"state"`
		// 状态文本
		StateText string
		// 订单项
		Items []*OrderItem `db:"-"`
	}

	/*
	   o.order_no,po.order_no as parent_no,
	       vendor_id,o.shop_id,s.name as shop_name,
	       o.item_amount,o.discount_amount,o.express_fee,
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
		Id        int
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
		OrderId int
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
        Id:v.Id,
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
        FinalAmount:v.FinalAmount,
        IsPaid:v.IsPaid,
        IsSuspend:v.IsSuspend,
        Note:v.BuyerRemark,
        Remark:v.Remark,
        UpdateTime:v.UpdateTime,
        State:v.State,
        Items:v.Items,
        StateText:order.OrderState(v.State).String(),
    }
}*/
