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
	GetProductByIds(partnerId int, ids ...int) ([]IProduct, error)
}
