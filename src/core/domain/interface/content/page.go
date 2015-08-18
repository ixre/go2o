/**
 * Copyright 2015 @ z3q.net.
 * name : page.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

type IPage interface {
	// 获取领域编号
	GetDomainId() int

	// 获取值
	GetValue() *ValuePage

	// 设置值
	SetValue(*ValuePage) error

	// 保存
	Save() (int, error)
}
