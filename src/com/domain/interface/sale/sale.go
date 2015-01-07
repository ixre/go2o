/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-08 11:44
 * description :
 * history :
 */

package sale

type ISale interface {
	GetAggregateRootId() int

	CreateProduct(*ValueProduct) IProduct

	// 根据产品编号获取产品
	GetProduct(int) IProduct

	// 删除商品
	DeleteProduct(int) error
}
