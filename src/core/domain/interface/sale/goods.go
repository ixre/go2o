/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:49
 * description :
 * history :
 */

package sale

type IGoods interface {
	GetDomainId() int

	GetValue() ValueGoods

	// 是否上架
	IsOnShelves() bool

	SetValue(*ValueGoods) error

	Save() (int, error)

	// 生成快照
	GenerateSnapshot() (int, error)

	// 获取最新的快照
	GetLatestSnapshot() *GoodsSnapshot
}
