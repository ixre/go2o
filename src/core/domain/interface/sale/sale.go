/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-08 11:44
 * description :
 * history :
 */

package sale

type ISale interface {
	GetAggregateRootId() int

	CreateGoods(*ValueItem) IItem

	// 根据产品编号获取产品
	GetGoods(int) IItem

	// 创建分类
	CreateCategory(*ValueCategory) ICategory

	// 获取分类
	GetCategory(int) ICategory

	// 获取所有分类
	GetCategories() []ICategory

	// 删除分类
	DeleteCategory(int) error

	// 获取所有的销售标签
	GetAllSaleTags() []ISaleTag

	// 初始化销售标签
	InitSaleTags() error

	// 获取销售标签
	GetSaleTag(id int) ISaleTag

	// 根据Code获取销售标签
	GetSaleTagByCode(code string) ISaleTag

	// 创建销售标签
	CreateSaleTag(v *ValueSaleTag) ISaleTag

	// 删除销售标签
	DeleteSaleTag(id int) error

	// 删除商品
	DeleteGoods(int) error

	// 获取指定的商品快照
	GetGoodsSnapshot(id int) *GoodsSnapshot

	// 根据Key获取商品快照
	GetGoodsSnapshotByKey(key string) *GoodsSnapshot
}
