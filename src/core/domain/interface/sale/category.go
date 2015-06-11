/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
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

	// 获取子栏目的编号
	GetChildId()[]int
}
