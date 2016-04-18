/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:49
 * description :
 * history :
 */

package sale

import "go2o/src/core/infrastructure/domain"

type (

	ICategory interface {
		GetDomainId() int

		GetValue() *ValueCategory

		GetOption()domain.IOptionStore

		SetValue(*ValueCategory) error

		Save() (int, error)

		// 获取子栏目的编号
		GetChildId() []int
	}
)

var(
	C_OptionViewName string = "viewName" //显示的视图名称
	C_OptionDescribe string = "describe" //描述
)

