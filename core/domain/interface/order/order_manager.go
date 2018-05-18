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
	"go2o/core/dto"
)

type (
	//　订单服务
	IOrderManager interface {
		// 统一调用
		Unified(orderNo string, sub bool) IUnifiedOrderAdapter
		// 预创建普通订单
		PrepareNormalOrder(c cart.ICart) (IOrder, error)
		// 预创建批发订单
		PrepareWholesaleOrder(c cart.ICart) ([]IOrder, error)
		// 提交批发订单
		SubmitWholesaleOrder(c cart.ICart, data IPostedData) (map[string]string, error)
		// 提交交易类订单
		SubmitTradeOrder(o *ComplexOrder, tradeRate float64) (IOrder, error)
		// 接收在线交易支付的通知，不主动调用
		NotifyOrderTradeSuccess(orderNo string, subOrder bool) error
		// 提交订单
		SubmitOrder(c cart.ICart, addressId int64, couponCode string,
			balanceDiscount bool) (IOrder, *SubmitReturnData, error)
		// 获取可用的订单号, 系统直营传入vendor为0
		GetFreeOrderNo(vendor int32) string
		// 根据订单编号获取订单
		GetOrderById(orderId int64) IOrder
		// 根据订单号获取订单
		GetOrderByNo(orderNo string) IOrder

		// 获取子订单
		GetSubOrder(id int64) ISubOrder
	}

	// 订单提交附带的数据
	IPostedData interface {
		// 获取勾选的商品和SKU数据
		CheckedData() map[int64][]int64
		// 获取收货地址编号
		AddressId() int64
		// 获取订单留言
		GetComment(sellerId int32) string
	}

	// 统一订单适配器
	IUnifiedOrderAdapter interface {
		// 复合的订单信息
		Complex() *ComplexOrder
		// 取消订单
		Cancel(reason string) error
		// 订单确认
		Confirm() error
		// 备货完成
		PickUp() error
		// 订单发货,并记录配送服务商编号及单号
		Ship(spId int32, spOrder string) error
		// 消费者收货
		BuyerReceived() error
		// 获取订单日志
		LogBytes() []byte
		// 物流日志
	}

	IOrderRepo interface {
		// 获取订单服务
		Manager() IOrderManager
		// 创建订单
		CreateOrder(*Order) IOrder
		// 获取可用的订单号, 系统直营传入vendor为0
		GetFreeOrderNo(vendorId int32) string
		// 获取订单编号
		GetOrderId(orderNo string, subOrder bool) int64
		// Get OrderList
		GetOrder(where string, arg ...interface{}) *Order
		// Save OrderList
		SaveOrder(v *Order) (int, error)

		// 保存订单优惠券绑定
		SaveOrderCouponBind(*OrderCoupon) error
		// 获取订单的促销绑定
		GetOrderPromotionBinds(orderNo string) []*OrderPromotionBind
		// 保存订单的促销绑定
		SavePromotionBindForOrder(*OrderPromotionBind) (int32, error)

		// 生成空白订单,并保存返回对象
		CreateNormalSubOrder(*NormalSubOrder) ISubOrder
		// 根据编号获取订单
		GetNormalOrderById(orderId int64) *NormalOrder
		// 保存订单,返回订单编号
		SaveNormalOrder(v *NormalOrder) (int, error)
		// 获取订单的所有子订单
		GetNormalSubOrders(orderId int64) []*NormalSubOrder
		// 保存订单日志
		SaveNormalSubOrderLog(*OrderLog) error
		// 获取子订单,todo: 删除
		GetSubOrder(id int64) *NormalSubOrder
		// 根据子订单
		GetSubOrderByOrderNo(s string) ISubOrder
		// 保存子订单
		SaveSubOrder(value *NormalSubOrder) (int, error)
		// 保存子订单的商品项,并返回编号和错误
		SaveOrderItem(subOrderId int64, value *SubOrderItem) (int32, error)

		// 获取订单项
		GetSubOrderItems(orderId int64) []*SubOrderItem
		// 获取订单的操作记录
		GetSubOrderLogs(orderId int64) []*OrderLog
		// 根据商品快照获取订单项
		GetOrderItemBySnapshotId(orderId int64, snapshotId int32) *SubOrderItem
		// 根据商品快照获取订单项数据传输对象
		GetOrderItemDtoBySnapshotId(orderId int64, snapshotId int32) *dto.OrderItem

		// Get WholesaleOrder
		GetWholesaleOrder(where string, v ...interface{}) *WholesaleOrder
		// Save WholesaleOrder
		SaveWholesaleOrder(v *WholesaleOrder) (int, error)
		// Save WholesaleItem
		SaveWholesaleItem(v *WholesaleItem) (int, error)
		// Select WholesaleItem
		SelectWholesaleItem(where string, v ...interface{}) []*WholesaleItem

		// Get TradeOrder
		GetTradeOrder(where string, v ...interface{}) *TradeOrder
		// Save TradeOrder
		SaveTradeOrder(v *TradeOrder) (int, error)
	}
)
