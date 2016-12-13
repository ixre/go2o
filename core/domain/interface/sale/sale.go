/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 11:44
 * description :
 * history :
 */

package sale

import "go2o/core/domain/interface/product"

type ISale interface {
	GetAggregateRootId() int32

	// 类目服务
	CategoryManager() product.IGlobCatService

	// 标签服务
	LabelManager() ILabelManager

	// 货品服务
	ItemManager() IItemManager

	// 商品服务
	GoodsManager() IGoodsManager
}
