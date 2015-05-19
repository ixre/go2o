/**
 * Copyright 2015 @ S1N1 Team.
 * name : sale_tag_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

type ISaleTagRep interface {
	// 创建销售标签
	CreateSaleTag(v *ValueSaleTag) ISaleTag

	// 获取所有的销售标签
	GetAllValueSaleTags(partnerId int)[]*ValueSaleTag

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
	GetValueGoods(partnerId, tagId, begin, end int) []*ValueGoods
}
