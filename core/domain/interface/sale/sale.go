/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 11:44
 * description :
 * history :
 */

package sale

type ISale interface {
	GetAggregateRootId() int64

	// 类目服务
	CategoryManager() ICategoryManager

	// 标签服务
	LabelManager() ILabelManager

	// 货品服务
	ItemManager() IItemManager

	// 商品服务
	GoodsManager() IGoodsManager
}
