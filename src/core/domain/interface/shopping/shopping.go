/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:19
 * description :
 * history :
 */

package shopping

import (
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
)

//　购物聚合根
type IShopping interface {
	GetAggregateRootId() int

	// 创建订单,如果为已存在的订单则没有Cart.
	// todo:需重构为单独的类型
	CreateOrder(*ValueOrder, ICart) IOrder

	//创建购物车
	// @buyerId 为购买会员ID,0表示匿名购物车
	NewCart(buyerId int) ICart

	// 检查购物车
	CheckCart(cart ICart) error

	// 根据数据获取购物车,
	// 如果member的cart与key不一致，则合并购物车；
	// 如果会员没有购物车，则绑定为key的购物车
	// 如果都没有，则创建一个购物车
	//GetCart(key string,memberId int) (ICart, error)

	// 根据数据获取购物车
	GetCartByKey(key string) (ICart, error)

	// 获取购物车
	GetShoppingCart(buyerId int, cartKey string) ICart

	// 获取没有结算的购物车
	GetCurrentCart(buyerId int) (ICart, error)

	// 绑定购物车会员编号
	BindCartBuyer(cartKey string, buyerId int) error

	// 将购物车转换为订单
	ParseShoppingCart(memberId int) (IOrder, member.IMember, ICart, error)

	// 组装订单
	BuildOrder(memberId int, couponCode string) (IOrder, ICart, error)

	// 提交订单
	SubmitOrder(memberId int, couponCode string, useBalanceDiscount bool) (string, error)

	// 获取可用的订单号
	GetFreeOrderNo() string

	// 根据订单号获取订单
	GetOrderByNo(orderNo string) (IOrder, error)

	// 自动设置订单
	OrderAutoSetup(f func(error))

	// 智能选择门店
	SmartChoiceShop(address string) (partner.IShop, error)
}
