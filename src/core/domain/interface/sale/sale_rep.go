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
	GetSale(merchantId int) ISale

	// 获取货品
	GetValueItem(merchantId, itemId int) *ValueItem

	// 根据id获取货品
	GetItemByIds(ids ...int) ([]*ValueItem, error)

	SaveValueItem(*ValueItem) (int, error)

	// 获取在货架上的商品
	GetPagedOnShelvesItem(merchantId int, catIds []int, start, end int) (total int, goods []*ValueItem)

	// 获取货品销售总数
	GetItemSaleNum(merchantId int, id int) int

	// 删除货品
	DeleteItem(merchantId, goodsId int) error

	/*********** Category ************/

	// 保存分类
	SaveCategory(*ValueCategory) (int, error)

	DeleteCategory(merchantId, id int) error

	GetCategory(merchantId, id int) *ValueCategory

	GetCategories(merchantId int) CategoryList

	// 获取与栏目相关的栏目
	GetRelationCategories(merchantId, categoryId int) CategoryList

	// 获取子栏目
	GetChildCategories(merchantId, categoryId int) CategoryList

	// 保存快照
	SaveSnapshot(*GoodsSnapshot) (int, error)

	// 获取最新的商品快照
	GetLatestGoodsSnapshot(goodsId int) *GoodsSnapshot

	// 获取指定的商品快照
	GetGoodsSnapshot(id int) *GoodsSnapshot

	// 根据Key获取商品快照
	GetGoodsSnapshotByKey(key string) *GoodsSnapshot
}
