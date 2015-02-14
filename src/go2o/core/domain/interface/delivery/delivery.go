/**
 * Copyright 2014 @ ops Inc.
 * name : 
 * author : newmin
 * date : 2014-02-12 16:21
 * description :
 * history :
 */
package delivery

type IDelivery interface{
	// 返回聚合编号
	GetAggregateRootId() int

	// 等同于GetAggregateRootId()
	GetPartnerId()int

	//　获取覆盖区域
	GetCoverageArea(id int)ICoverageArea

	// 查看单个所在的区域
	FindSingleCoverageArea(lng,lat float32)ICoverageArea

	// 查找所有所在的区域
	FindCoverageAreas(lng,lat float32)[]ICoverageArea
}
