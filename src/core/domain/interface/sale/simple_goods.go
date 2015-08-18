/**
 * Copyright 2015 @ z3q.net.
 * name : simple_goods.go
 * author : jarryliu
 * date : 2015-08-18 09:24
 * description :
 * history :
 */
package sale

// 简单商品信息
type SimpleGoods struct{
	GoodsId  int 	`json:"id"`
	GoodsImage string `json:"img"`
	Name string  `json:"name"`
	Quantity string `json:"qty"`
}