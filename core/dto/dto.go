/**
 * Copyright 2015 @ to2.net.
 * name : message_result
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dto

type (
	//操作消息结果
	MessageResult struct {
		Result  bool   `json:"result"`
		Message string `json:"message"`
		Tag     int    `json:"tag"`
	}

	// 站内信
	SiteMessage struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 消息类型
		Type int `db:"msg_type"`
		// 消息用途
		UseFor       int `db:"use_for"`
		SenderUserId int
		SenderName   string
		// 是否只能阅读
		Readonly int `db:"read_only"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 数据
		Data interface{}
		// 接收者编号
		ToId int32 `db:"to_id"`
		// 接收者角色
		ToRole int `db:"to_role"`
		// 是否阅读
		HasRead int `db:"has_read"`
		// 阅读时间
		ReadTime int64 `db:"read_time"`
	}

	// 商品查询复合信息
	GoodsComplex struct {
		GoodsId int    `db:"id"`
		ItemId  int    `db:"item_id"`
		MchId   int    `db:"mch_id"`
		MchName string `db:"mch_name"`
	}

	PagedMemberAfterSalesOrder struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 订单编号
		OrderNo string `db:"order_id"`
		// 运营商编号
		VendorId int32 `db:"vendor_id"`
		// 运营商名称
		VendorName string `db:"vendor_name"`
		// 类型，退货、换货、维修
		Type       int `db:"type"`
		SkuId      int32
		GoodsTitle string
		GoodsImage string
		// 退货的商品项编号
		SnapshotId int32 `db:"snap_id"`
		// 商品数量
		Quantity int `db:"quantity"`
		// 售后单状态
		State int `db:"state"`
		// 提交时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
		// 订单状态
		StateText string `db:"-"`
	}

	// 分页商户售后单
	PagedVendorAfterSalesOrder struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 订单编号
		OrderNo string `db:"order_id"`
		// 会员编号
		BuyerId int32 `db:"vendor_id"`
		// 会员名称
		BuyerName string `db:"buyer_name"`
		// 类型，退货、换货、维修
		Type       int `db:"type"`
		SkuId      int32
		GoodsTitle string
		GoodsImage string
		// 退货的商品项编号
		SnapshotId int32 `db:"snap_id"`
		// 商品数量
		Quantity int `db:"quantity"`
		// 售后单状态
		State int `db:"state"`
		// 提交时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
		// 订单状态
		StateText string `db:"-"`
	}

	// 店铺收藏
	PagedShopFav struct {
		Id         int32  `db:"id"`
		ShopId     int32  `db:"shop_id"`
		ShopName   string `db:"shop_name"`
		MchId      int32  `db:"mch_id"`
		Logo       string `db:"logo"`
		UpdateTime int64  `db:"update_time"`
	}

	// 商品收藏
	PagedGoodsFav struct {
		Id         int32  `db:"id"`
		SkuId      int32  `db:"sku_id"`
		GoodsName  string `db:"goods_name"`
		Image      string `db:"image"`
		OnShelves  int    `db:"on_shelves"`
		StockNum   int    `db:"stock_num"`
		SalePrice  string `db:"sale_price"`
		UpdateTime int64  `db:"update_time"`
	}
	// 分类
	//Category struct {
	//	ID    int32
	//	Name  string
	//	Icon  string
	//	Url   string
	//	Level int
	//	Child []Category
	//}
	ListOnlineShop struct {
		Id         int32  `db:"sp.id"`
		Name       string `db:"sp.name"`
		Alias      string `db:"alias"`
		Host       string `db:"ol.host"`
		Logo       string `db:"logo"`
		CreateTime int64  `db:"sp.create_time" json:"-"`
	}
)
