/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 17:39
 * description :
 * history :
 */
package user

type IPerson interface {
	//获取领域对象编号
	GetDomainId() int32

	GetValue() PersonValue

	SetValue(*PersonValue) error

	Save() (int32, error)
}
