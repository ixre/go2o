/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-08 10:45
 * description :
 * history :
 */

package sale

// 销售仓库
type ISaleRep interface {
	GetSale(partnerId int) ISale

	GetValueProduct(partnerId, productId int) *ValueProduct

	GetProductByIds(partnerId int, ids ...int) ([]IProduct, error)

	SaveProduct(*ValueProduct) (int, error)

	GetProductsByCid(partnerId, categoryId, num int) []*ValueProduct

	DeleteProduct(partnerId, productId int) error
}
