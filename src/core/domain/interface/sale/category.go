/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-08 10:49
 * description :
 * history :
 */

package sale

type ICategory interface {
	GetDomainId() int

	GetValue() ValueCategory

	SetValue(*ValueCategory) error

	Save() (int, error)
}
