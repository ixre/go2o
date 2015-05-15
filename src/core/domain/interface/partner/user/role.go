/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-14 16:44
 * description :
 * history :
 */
package user

type IRole interface {
	//获取领域对象编号
	GetDomainId() int

	GetValue() RoleValue

	SetValue(*RoleValue) error

	Save() (int, error)
}
