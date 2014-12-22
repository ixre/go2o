/**
 * Copyright 2014 @ Ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-04 23:46
 * description :
 * history :
 */

package promotion

// 促销聚合根
type IPromotion interface {
	GetAggregateRootId() int
	GetCoupon(id int) ICoupon
	CreateCoupon(val *ValueCoupon) ICoupon
}
