package order

type TradeOrderValue struct {
	// 商户编号
	MerchantId int
	// 门店编号
	StoreId int
	// 买家编号
	BuyerId int
	// 商品金额
	ItemAmount int
	// 抵扣金额
	DiscountAmount int
	// 订单主题
	Subject string
}
