/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-23 22:06
 * description :
 * history :
 */

package member

type IDeliver interface {
	GetDomainId() int

	GetValue() DeliverAddress

	SetValue(*DeliverAddress) error

	Save() (int, error)
}
