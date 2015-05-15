/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-04 23:46
 * description :
 * history :
 */

//整单折扣 (自动满１)
//满立减
//满就送 (送汤)
//**优惠券

package promotion

// 促销聚合根
type IPromotion interface {
	GetAggregateRootId() int
	GetCoupon(id int) ICoupon
	CreateCoupon(val *ValueCoupon) ICoupon
}
