/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:52
 * description :
 * history :
 */

package domain

//聚合根
type IAggregateRoot interface {
	// 获取聚合根编号
	GetAggregateRootId() int32
}

//领域对象
type IDomain interface {
	// 获取领域对象编号
	GetDomainId() int32
}

//值对象
type IValueObject interface {
	Equal(interface{}) bool
}
