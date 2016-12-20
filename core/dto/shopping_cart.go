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
	Id         int32  `json:"-"`
	CartKey    string `json:"key"`
	BuyerId    int32  `json:"buyer"`
	Summary    string `json:"summary"`
	UpdateTime int64  `json:"update_time"`
	//Items      []*CartItem `json:"items"`
	TotalNum int32   `json:"total_num"` // 总数量
	TotalFee float32 `json:"total"`
	OrderFee float32 `json:"fee"`
	// 运营商
	Vendors []*CartVendorGroup `json:"vendors"`
}

type CartVendorGroup struct {
	VendorId   int32       `json:"vendorId"`
	VendorName string      `json:"vendorName"`
	ShopId     int32       `json:"shopId"`
	ShopName   string      `json:"shopName"`
	Items      []*CartItem `json:"items"`
	//结算商品项数目
	CheckedNum int `json:"checkedNum"`
}

type CartItem struct {
	GoodsId    int32   `json:"skuId"`
	GoodsName  string  `json:"name"`
	GoodsNo    string  `json:"no"`
	SmallTitle string  `json:"title"`
	GoodsImage string  `json:"image"`
	Quantity   int32   `json:"num"`
	SpecWord   string  `json:"specWord"`
	Price      float32 `json:"price"`
	SalePrice  float32 `json:"salePrice"`
	// 是否结算
	Checked bool `json:"checked"`
}
