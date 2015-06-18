/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:49
 * description :
 * history :
 */

package sale

// 商品
type IItem interface {
	GetDomainId() int

	// 获取商品的值
	GetValue() ValueItem

	// 是否上架
	IsOnShelves() bool

	// 获取销售标签
	GetSaleTags() []*ValueSaleTag

	// 保存销售标签
	SaveSaleTags([]int) error

	// 设置商品值
	SetValue(*ValueItem) error

	// 保存
	Save() (int, error)
}
