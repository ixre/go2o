/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-08 10:31
 * description :
 * history :
 */

package shopping

import (
	"com/domain/interface/sale"
)

type ValueCart struct {
	//购物车商品
	Items []sale.IProduct
	//购物车商品数量
	Quantities map[int]int
	//客户端计算的金额
	ClientFee float32
}
