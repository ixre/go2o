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
	CreateProduct(*ValueProduct) IProduct
}
