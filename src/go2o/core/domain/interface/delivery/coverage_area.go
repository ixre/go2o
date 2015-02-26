/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-12 16:38
 * description :
 * history :
 */
package delivery

type ICoverageArea interface {
	GetDomainId() int

	GetValue() CoverageValue

	SetValue(*CoverageValue) error

	Save() (int, error)
}
