/**
 * Copyright 2015 @ z3q.net.
 * name : sale_tag_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

import "go2o/src/core/domain/interface/valueobject"

type ISaleTagRep interface {
	// 创建销售标签
	CreateSaleTag(v *ValueSaleTag) ISaleTag

	// 获取所有的销售标签
	GetAllValueSaleTags(partnerId int) []*ValueSaleTag

	// 获取销售标签值
	GetValueSaleTag(partnerId int, tagId int) *ValueSaleTag

	// 根据Code获取销售标签
	GetSaleTagByCode(partnerId int, code string) *ValueSaleTag

	// 删除销售标签
	DeleteSaleTag(partnerId int, id int) error

	// 获取销售标签
	GetSaleTag(partnerId int, tagId int) ISaleTag

	// 保存销售标签
	SaveSaleTag(partnerId int, v *ValueSaleTag) (int, error)

	// 获取商品
	GetValueGoodsBySaleTag(partnerId, tagId int, sortBy string,
		begin, end int) []*valueobject.Goods

	// 获取分页商品
	GetPagedValueGoodsBySaleTag(partnerId, tagId int, sortBy string,
		begin, end int) (int, []*valueobject.Goods)

	// 获取商品的销售标签
	GetItemSaleTags(itemId int) []*ValueSaleTag

	// 清理商品的销售标签
	CleanItemSaleTags(itemId int) error

	// 保存商品的销售标签
	SaveItemSaleTags(itemId int, tagIds []int) error
}
