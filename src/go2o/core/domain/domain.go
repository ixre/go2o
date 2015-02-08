/**
 * Copyright 2014 @ Ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-03 23:52
 * description :
 * history :
 */

package domain

//聚合根
type IAggregateRoot interface {
	GetAggregateRootId() int
}

//领域对象
type IDomain interface {
	//获取领域对象编号
	GetDomainId() int
}

//值对象
type IValueObject interface {
	Equal(interface{}) bool
}
