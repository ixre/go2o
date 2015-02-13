/**
 * Copyright 2014 @ ops Inc.
 * name : 
 * author : newmin
 * date : 2015-02-12 16:21
 * description :
 * history :
 */
package delivery

type IDelivery interface{
	// 返回聚合编号
	GetAggregateRootId() int

	// 等同于GetAggregateRootId()
	GetPartnerId()int
}
