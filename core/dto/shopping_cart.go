/**
 * Copyright 2013 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-22 21:56
 * description :
 * history :
 */

package dto

type ShoppingCart struct {
	Id         int         `json:"-"`
	CartKey    string      `json:"key"`
	BuyerId    int         `json:"buyer"`
	Summary    string      `json:"summary"`
	UpdateTime int64       `json:"update_time"`
	//Items      []*CartItem `json:"items"`
	TotalNum   int         `json:"total_num"` // 总数量
	TotalFee   float32     `json:"total"`
	OrderFee   float32     `json:"fee"`
	// 运营商
	Vendors    []*CartVendorGroup `json:"vendors"`
}

type CartVendorGroup struct{
	VendorId  int  `json:"vendorId"`
	VendorName string `json:"vendorName"`
	ShopId    int   `json:"shopId"`
	ShopName  string  `json:"shopName"`
	Items   []*CartItem  `json:"items"`
	//结算商品项数目
	SettleNum  int   `json:"settleNum"`
}

type CartItem struct {
	GoodsId    int     `json:"id"`
	GoodsName  string  `json:"name"`
	GoodsNo    string  `json:"no"`
	SmallTitle string  `json:"title"`
	GoodsImage string  `json:"image"`
	Quantity   int     `json:"num"`
	Price      float32 `json:"price"`
	SalePrice  float32 `json:"sale_price"`
	// 是否结算
	IsSettle   bool    `json:"isSettle"`
}
