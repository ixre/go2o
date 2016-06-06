/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 11:44
 * description :
 * history :
 */

package sale

import (
	"go2o/core/domain/interface/valueobject"
)

type ISale interface {
	GetAggregateRootId() int

	// 分类服务
	CategoryManager() ICategoryManager

	// 创建产品
	CreateItem(*ValueItem) IItem

	// 根据产品编号获取货品
	GetItem(int) IItem

	// 删除货品
	DeleteItem(int) error

	// 创建商品
	CreateGoodsByItem(IItem, *ValueGoods) IGoods

	// 创建商品
	CreateGoods(*ValueGoods) IGoods

	// 根据产品编号获取商品
	GetGoods(int) IGoods

	// 根据产品SKU获取商品
	GetGoodsBySku(itemId, sku int) IGoods

	// 获取所有的销售标签
	GetAllSaleLabels() []ISaleLabel

	// 初始化销售标签
	InitSaleLabels() error

	// 获取销售标签
	GetSaleLabel(id int) ISaleLabel

	// 根据Code获取销售标签
	GetSaleLabelByCode(code string) ISaleLabel

	// 创建销售标签
	CreateSaleLabel(v *SaleLabel) ISaleLabel

	// 删除销售标签
	DeleteSaleLabel(id int) error

	// 删除商品
	DeleteGoods(int) error

	// 获取指定的商品快照
	GetGoodsSnapshot(id int) *GoodsSnapshot

	// 根据Key获取商品快照
	GetGoodsSnapshotByKey(key string) *GoodsSnapshot

	// 获取指定数量已上架的商品
	GetOnShelvesGoods(start, end int, sortBy string) []*valueobject.Goods
}
