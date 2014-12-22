/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-05 17:19
 * description :
 * history :
 */

package shopping

//　购物聚合根
type IShopping interface {
	GetAggregateRootId() int
	// 创建订单,如果为已存在的订单则没有Cart.
	// todo:需重构为单独的类型
	CreateOrder(*ValueOrder, ICart) IOrder
	//创建购物车
	CreateCart(value *ValueCart) ICart
	// 根据数据获取购物车
	GetCart(cart string) (ICart, error)
	// 将购物车转换为订单
	GetOrderByCart(cartStr string) (*ValueOrder, error)
	// 组装订单
	BuildOrder(memberId int, cart string, couponCode string) (IOrder, error)
	// 提交订单
	SubmitOrder(memberId, shopId int, payMethod int, deliverAddrId int,
		cart string, couponCode string, note string) (string, error)

	// 获取可用的订单号
	GetFreeOrderNo() string

	// 根据订单号获取订单
	GetOrderByNo(orderNo string) (IOrder, error)
}
