package cart

// 零售购物车
type IRetailCart interface {
	// 设置买家编号
	SetBuyer(buyerId int) error
	// 重置项目
	ResetItems(items []RetailCartItem) error
	// 获取购物车值
	Value() RetailCart
	// 获取商品集合
	Items() []RetailCartItem
}

// 零售购物车
type RetailCart struct {
	// 编号
	ID int64 `db:"id" pk:"yes" auto:"yes"`
	// 购物车编码
	Code string `db:"code"`
	// 运营商编号
	VendorId int64 `db:"vendor_id"`
	// 店铺编号
	ShopId int64 `db:"shop_id"`
	// 买家编号
	BuyerId int64 `db:"buyer_id"`
	// 买家姓名
	BuyerName string `db:"buyer_name"`
	// 购物车标志
	CartFlag int64 `db:"cart_flag"`
	// 优先级
	Priority int `db:"priority"`
	// 购物车序号
	SortNum int `db:"sort_num"`
	// 购物车金额
	TotalAmount int64 `db:"total_amount"`
	// 购物车描述
	CartDesc string `db:"cart_desc"`
	// 客位编号
	PlaceId int64 `db:"place_id"`
	// 创建时间
	CreateTime int64 `db:"create_time"`
	// 修改时间
	UpdateTime int64 `db:"update_time"`
}

// 零售购物车项目
type RetailCartItem struct {
	// 编号
	ID int64 `db:"id" pk:"yes" auto:"yes"`
	// 购物车编号
	CartId int64 `db:"cart_id"`
	// 商品编号
	ItemId int64 `db:"item_id"`
	// 商品种类
	ItemKind int `db:"item_kind"`
	// 商品标题
	ItemTitle string `db:"item_title"`
	// SKU编号
	SkuId int64 `db:"sku_id"`
	// SKU文本
	SkuText string `db:"sku_text"`
	// 卡编号
	CardId int64 `db:"card_id"`
	// 卡项目编号
	CardItemId int64 `db:"card_item_id"`
	// 数量
	Quantity int `db:"quantity"`
	// 商品单价金额
	UnitPrice int `db:"unit_price"`
	// 商品实际金额
	FinalPrice int `db:"final_price"`
	// 调整金额
	AdjustAmount int64 `db:"adjust_amount"`
	// 最终价
	FinalFee int64 `db:"final_fee"`
	// 是否勾选结算
	Checked int `db:"checked"`
}

// 零售购物车提成
type RetailCartAward struct {
	// 编号
	ID int64 `db:"id" pk:"yes" auto:"yes"`
	// 购物车编号
	CartId int64 `db:"cart_id"`
	// 购物车项目编号
	CartItemId int64 `db:"cart_item_id"`
	// 提成人员编号
	OperatorId int64 `db:"operator_id"`
	// 提成金额
	AwardValue int64 `db:"award_value"`
}
