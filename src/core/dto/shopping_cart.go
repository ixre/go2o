/**
 * Copyright 2013 @ ops Inc.
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
	UpdateTime int64       `json:"updateTime"`
	Items      []*CartItem `json:"items"`
	TotalFee   float32     `json:"total"`
	OrderFee   float32     `json:"fee"`
	IsBought	int 		`json:"isBought"`		//是否已经购买
}

type CartItem struct {
	GoodsId    int     `json:"id"`
	GoodsName  string  `json:"name"`
	GoodsNo    string  `json:"no"`
	SmallTitle string  `json:"title"`
	GoodsImage string  `json:"image"`
	Num        int     `json:"num"`
	Price      float32 `json:"price"`
	SalePrice  float32 `json:"salePrice"`
}
