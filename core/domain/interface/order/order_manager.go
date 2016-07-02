/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:19
 * description :
 * history :
 */

package order

import (
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant/shop"
)

//　购物聚合根
type IOrderManager interface {
	// 创建订单,如果为已存在的订单则没有Cart.
	// todo:需重构为单独的类型
	CreateOrder(*ValueOrder, cart.ICart) IOrder

	// 根据数据获取购物车,
	// 如果member的cart与key不一致，则合并购物车；
	// 如果会员没有购物车，则绑定为key的购物车
	// 如果都没有，则创建一个购物车
	//GetCart(key string,memberId int) (ICart, error)

	// 根据数据获取购物车
	//GetCartByKey(key string) (ICart, error)

	// 获取购物车
	//GetShoppingCart(cartKey string) ICart

	// 将购物车转换为订单
	ParseToOrder(c cart.ICart) (IOrder, member.IMember, error)

	// 组装订单
	BuildOrder(c cart.ICart, subject string, couponCode string) (IOrder, error)

	// 提交订单
	SubmitOrder(c cart.ICart, subject string, couponCode string,
		balanceDiscount bool) (string, error)

	// 获取可用的订单号, 系统直营传入vendor为0
	GetFreeOrderNo(vendor int) string

	// 根据订单号获取订单
	GetOrderByNo(orderNo string) IOrder

	// 自动设置订单
	OrderAutoSetup(f func(error))

	// 智能选择门店
	SmartChoiceShop(address string) (shop.IShop, error)

	// 智能确定订单
	SmartConfirmOrder(order IOrder) error
}
