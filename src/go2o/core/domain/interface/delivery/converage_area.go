/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2015-02-12 16:38
 * description :
 * history :
 */
package delivery

type IConverageArea interface {
	GetDomainId() int

	GetValue() ConverageValue

	SetValue(*ConverageValue) error

	Save() (int, error)
}
