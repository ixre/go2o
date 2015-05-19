/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:49
 * description :
 * history :
 */

package sale

type IGoods interface {
	GetDomainId() int

	// 获取商品的值
	GetValue() ValueGoods

	// 是否上架
	IsOnShelves() bool

	// 获取销售标签
	GetSaleTags() []*ValueSaleTag

	// 保存销售标签
	SaveSaleTags([]int) error

	// 设置商品值
	SetValue(*ValueGoods) error

	// 保存
	Save() (int, error)

	// 生成快照
	GenerateSnapshot() (int, error)

	// 获取最新的快照
	GetLatestSnapshot() *GoodsSnapshot
}
