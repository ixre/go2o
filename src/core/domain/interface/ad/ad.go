/**
 * Copyright 2015 @ z3q.net.
 * name : IAdvertisement
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

type IAdvertisement interface {
	// 获取领域对象编号
	GetDomainId() int

	// 是否为系统内置的广告
	System() bool

	// 广告类型
	Type() int

	// 广告名称
	Name() string

	// 设置值
	SetValue(*ValueAdvertisement) error

	// 获取值
	GetValue() *ValueAdvertisement

	// 保存广告
	Save() (int, error)
}
