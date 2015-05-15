/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
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
