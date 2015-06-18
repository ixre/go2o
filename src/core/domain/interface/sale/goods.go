/**
 * Copyright 2015 @ S1N1 Team.
 * name : sale_goods.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

// 商品
type IGoods interface {
	// 获取领域对象编号
	GetDomainId() int

	// 获取货品
	GetGoods() IItem

	// 设置值
	GetValue() *ValueGoods

	// 设置值
	SetValue(*ValueGoods) error

	// 保存
	Save() (int, error)

	// 生成快照
	GenerateSnapshot() (int, error)

	// 获取最新的快照
	GetLatestSnapshot() *GoodsSnapshot
}
