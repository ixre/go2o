/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package dto

// 列表页商品
type ListGoods struct {
	Id         int
	Name       string
	SmallTitle string
	Image      string
	Price      float32
	SalePrice  float32
}
