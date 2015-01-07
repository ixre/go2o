/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-11-22 20:01
 * description :
 * history :
 */

package partner

type IShop interface {
	GetDomainId() int

	GetValue() ValueShop

	SetValue(*ValueShop) error

	Save() (int, error)
}
