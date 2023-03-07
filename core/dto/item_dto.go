package dto

// 商品销售记录
type ItemSalesHistoryDto struct {
	// 买家用户代码
	BuyerUserCode string
	// 买家昵称
	BuyerName string
	// 买家头像
	BuyerPortrait string
	// 购买时间
	BuyTime int64
	// 订单状态
	OrderState int
}

// SearchItemResultDto 搜索商品返回数据
type SearchItemResultDto struct {
	// 商品编号
	ItemId int64
	// 商品标志
	ItemFlag int64
	// 商品编码
	Code string
	// 供货商编号
	SellerId int64
	// 商品标题
	Title string
	// 商品图片
	Image string
	// 价格区间
	PriceRange string
	// 库存数量
	StockNum int32
}
