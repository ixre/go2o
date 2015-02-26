/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-14 15:37
 * description :
 * history :
 */
package delivery

import (
	"go2o/core/domain/interface/delivery"
)

var _ delivery.IDelivery = new(Delivery)

type Delivery struct {
	id  int
	rep delivery.IDeliveryRep
}

func NewDelivery(id int, dlvRep delivery.IDeliveryRep) delivery.IDelivery {
	return &Delivery{
		id:  id,
		rep: dlvRep,
	}
}

// 返回聚合编号
func (this *Delivery) GetAggregateRootId() int {
	return this.id
}

// 等同于GetAggregateRootId()
func (this *Delivery) GetPartnerId() int {
	return this.id
}
