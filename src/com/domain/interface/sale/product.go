/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-08 10:49
 * description :
 * history :
 */

package sale

type IProduct interface {
	GetDomainId() int
	GetValue() ValueProduct
}
