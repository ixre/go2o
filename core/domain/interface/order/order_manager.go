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
	"go2o/core/domain/interface/payment"
	"go2o/core/dto"
)

type (
	//　购物聚合根
	IOrderManager interface {
		// 创建订单,如果为已存在的订单则没有Cart.
		// todo:需重构为单独的类型
		CreateOrder(*Order) IOrder

		// 生成空白订单,并保存返回对象
		CreateSubOrder(*SubOrder) ISubOrder

		// 将购物车转换为订单
		ParseToOrder(c cart.ICart) (IOrder, member.IMember, error)

		// 预生成订单及支付单
		PrepareOrder(c cart.ICart, subject string, couponCode string) (IOrder,
			payment.IPaymentOrder, error)

		// 提交订单
		SubmitOrder(c cart.ICart, subject string, couponCode string,
			balanceDiscount bool) (IOrder, payment.IPaymentOrder, error)

		// 获取可用的订单号, 系统直营传入vendor为0
		GetFreeOrderNo(vendor int32) string

		// 根据订单编号获取订单
		GetOrderById(orderId int32) IOrder

		// 根据订单号获取订单
		GetOrderByNo(orderNo string) IOrder

		// 接收在线交易支付的通知，不主动调用
		ReceiveNotifyOfOnlineTrade(orderId int32) error

		// 自动设置订单
		OrderAutoSetup(f func(error))

		// 智能选择门店
		SmartChoiceShop(address string) (shop.IShop, error)

		// 智能确定订单
		SmartConfirmOrder(order IOrder) error

		// 根据父订单编号获取购买的商品项
		GetItemsByParentOrderId(orderId int32) []*OrderItem

		//*********  子订单  *********//

		// 获取子订单
		GetSubOrder(id int32) ISubOrder
	}

	IOrderRepo interface {
		// 获取订单服务
		Manager() IOrderManager

		// 保存订单,返回订单编号
		SaveOrder(v *Order) (int32, error)

		// 保存订单优惠券绑定
		SaveOrderCouponBind(*OrderCoupon) error

		// 获取订单的促销绑定
		GetOrderPromotionBinds(orderNo string) []*OrderPromotionBind

		// 保存订单的促销绑定
		SavePromotionBindForOrder(*OrderPromotionBind) (int32, error)

		// 获取可用的订单号, 系统直营传入vendor为0
		GetFreeOrderNo(vendorId int32) string

		// 根据编号获取订单
		GetOrderById(id int32) *Order

		// 根据订单号获取订单
		GetValueOrderByNo(orderNo string) *Order

		// 获取等待处理的订单
		GetWaitingSetupOrders(vendorId int32) ([]*Order, error)

		// 保存订单日志
		SaveSubOrderLog(*OrderLog) error

		// 获取订单的所有子订单
		GetSubOrdersByParentId(orderId int32) []*SubOrder

		// 获取订单编号
		GetOrderId(orderNo string, subOrder bool) int32

		// 获取子订单
		GetSubOrder(id int32) *SubOrder

		// 根据订单号获取子订单
		GetSubOrderByNo(orderNo string) *SubOrder

		// 保存子订单
		SaveSubOrder(value *SubOrder) (int32, error)

		// 保存子订单的商品项,并返回编号和错误
		SaveOrderItem(subOrderId int32, value *OrderItem) (int32, error)

		// 获取订单项
		GetSubOrderItems(orderId int32) []*OrderItem

		// 根据父订单编号获取购买的商品项
		GetItemsByParentOrderId(orderId int32) []*OrderItem

		// 获取订单的操作记录
		GetSubOrderLogs(orderId int32) []*OrderLog

		// 根据商品快照获取订单项
		GetOrderItemBySnapshotId(orderId int32, snapshotId int32) *OrderItem

		// 根据商品快照获取订单项数据传输对象
		GetOrderItemDtoBySnapshotId(orderId int32, snapshotId int32) *dto.OrderItem
	}
)
