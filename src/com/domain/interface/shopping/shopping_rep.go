/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-05 17:50
 * description :
 * history :
 */

package shopping

type IShoppingRep interface {
	GetShopping(partnerId int) IShopping
	// 保存订单,返回订单编号
	SaveOrder(partnerId int, v *ValueOrder) (int, error)
	//　保存订单优惠券绑定
	SaveOrderCouponBind(*OrderCoupon) error
	// 获取可用的订单号
	GetFreeOrderNo(partnerId int) string

	// 根据订单号获取订单
	GetOrderByNo(partnerId int, orderNo string) (*ValueOrder, error)
}
