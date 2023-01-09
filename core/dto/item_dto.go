package dto

// 商品销售记录
type ItemSalesHistoryDto struct {
	// 买家用户代码
	BuyerUserCode string 
	// 买家昵称
	BuyerName string
	// 购买时间
	BuyTime int64
	// 订单状态
	OrderState int
}
