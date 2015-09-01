/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:45
 * description :
 * history :
 */

package sale

// 销售仓库
type ISaleRep interface {
	GetSale(partnerId int) ISale

	// 获取货品
	GetValueItem(partnerId, itemId int) *ValueItem

	// 根据id获取货品
	GetItemByIds(ids ...int) ([]*ValueItem, error)

	SaveValueItem(*ValueItem) (int, error)

	// 获取在货架上的商品
	GetPagedOnShelvesItem(partnerId int, catIds []int, start, end int) (total int, goods []*ValueItem)

	// 获取货品销售总数
	GetItemSaleNum(partnerId int, id int) int

	// 删除货品
	DeleteItem(partnerId, goodsId int) error

	/*********** Category ************/

	// 保存分类
	SaveCategory(*ValueCategory) (int, error)

	DeleteCategory(partnerId, id int) error

	GetCategory(partnerId, id int) *ValueCategory

	GetCategories(partnerId int) []*ValueCategory

	// 获取与栏目相关的栏目
	GetRelationCategories(partnerId, categoryId int) []*ValueCategory

	// 获取子栏目
	GetChildCategories(partnerId, categoryId int) []*ValueCategory

	// 保存快照
	SaveSnapshot(*GoodsSnapshot) (int, error)

	// 获取最新的商品快照
	GetLatestGoodsSnapshot(goodsId int) *GoodsSnapshot

	// 获取指定的商品快照
	GetGoodsSnapshot(id int) *GoodsSnapshot

	// 根据Key获取商品快照
	GetGoodsSnapshotByKey(key string) *GoodsSnapshot
}
