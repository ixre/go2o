/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 16:44
 * description :
 * history :
 */
package user

type IRole interface {
	//获取领域对象编号
	GetDomainId() int32

	GetValue() RoleValue

	SetValue(*RoleValue) error

	Save() (int32, error)
}
