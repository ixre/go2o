/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:50
 * description :
 * history :
 */

package shopping

type IShoppingRep interface {
	GetShopping(partnerId int) IShopping

	// 保存订单,返回订单编号
	SaveOrder(partnerId int, v *ValueOrder) (int, error)

	// 保存订单优惠券绑定
	SaveOrderCouponBind(*OrderCoupon) error

	// 获取订单的促销绑定
	GetOrderPromotionBinds(orderNo string) []*OrderPromotionBind

	// 保存订单的促销绑定
	SavePromotionBindForOrder(*OrderPromotionBind) (int, error)

	// 获取可用的订单号
	GetFreeOrderNo(partnerId int) string

	// 根据订单号获取订单
	GetOrderByNo(partnerId int, orderNo string) (*ValueOrder, error)

	// 获取等待处理的订单
	GetWaitingSetupOrders(partnerId int) ([]*ValueOrder, error)

	// 保存订单日志
	SaveOrderLog(*OrderLog) error

	// 获取购物车
	GetShoppingCart(key string) (*ValueCart, error)

	// 获取最新的购物车
	GetLatestCart(buyerId int) (*ValueCart, error)

	// 保存购物车
	SaveShoppingCart(*ValueCart) (int, error)

	// 移出购物车项
	RemoveCartItem(int) error

	// 保存购物车项
	SaveCartItem(*ValueCartItem) (int, error)

	// 清空购物车项
	EmptyCartItems(id int) error

	// 删除购物车
	DeleteCart(id int) error
}
