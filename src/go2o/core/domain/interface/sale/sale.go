/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-08 11:44
 * description :
 * history :
 */

package sale

type ISale interface {
	GetAggregateRootId() int

	CreateGoods(*ValueGoods) IGoods

	// 根据产品编号获取产品
	GetGoods(int) IGoods

	// 创建分类
	CreateCategory(*ValueCategory) ICategory

	// 获取分类
	GetCategory(int) ICategory

	// 获取所有分类
	GetCategories() []ICategory

	// 删除分类
	DeleteCategory(int) error

	// 删除商品
	DeleteGoods(int) error

	// 获取指定的商品快照
	GetGoodsSnapshot(id int) *GoodsSnapshot

	// 根据Key获取商品快照
	GetGoodsSnapshotByKey(key string) *GoodsSnapshot
}
