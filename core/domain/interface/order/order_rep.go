/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:50
 * description :
 * history :
 */

package order

type IOrderRep interface {
	// 获取订单服务
	Manager() IOrderManager

	// 保存订单,返回订单编号
	SaveOrder(v *ValueOrder) (int, error)

	// 保存订单优惠券绑定
	SaveOrderCouponBind(*OrderCoupon) error

	// 获取订单的促销绑定
	GetOrderPromotionBinds(orderNo string) []*OrderPromotionBind

	// 保存订单的促销绑定
	SavePromotionBindForOrder(*OrderPromotionBind) (int, error)

	// 获取可用的订单号, 系统直营传入vendor为0
	GetFreeOrderNo(vendorId int) string

	// 根据编号获取订单
	GetOrderById(id int) *ValueOrder

	// 根据订单号获取订单
	GetValueOrderByNo(orderNo string) *ValueOrder

	// 获取等待处理的订单
	GetWaitingSetupOrders(vendorId int) ([]*ValueOrder, error)

	// 获取订单项
	GetOrderItems(orderId int) []*OrderItem

	// 保存订单日志
	SaveOrderLog(*OrderLog) error
}
