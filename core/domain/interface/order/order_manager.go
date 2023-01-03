/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:19
 * description :
 * history :
 */

package order

import (
	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/core/dto"
)

type (
	// 订单服务
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
		SubmitTradeOrder(o *TradeOrderValue, tradeRate float64) (IOrder, error)
		// 接收在线交易支付的通知，不主动调用
		NotifyOrderTradeSuccess(orderNo string, subOrder bool) error
		// 提交订单
		SubmitOrder(data SubmitOrderData) (IOrder, *SubmitReturnData, error)
		// 获取可用的订单号, 系统直营传入vendor为0
		GetFreeOrderNo(vendor int64) string
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
		GetComment(sellerId int64) string
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
		// 更改收货人信息
		ChangeShipmentAddress(addressId int64) error
		// 消费者收货
		BuyerReceived() error
		// 删除订单
		Forbid() error
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
		GetFreeOrderNo(vendorId int64) string
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

		// CreateNormalSubOrder 生成空白订单,并保存返回对象
		CreateNormalSubOrder(*NormalSubOrder) ISubOrder
		// GetNormalSubOrders 获取订单的所有子订单
		GetNormalSubOrders(orderId int64) []*NormalSubOrder
		// SaveNormalSubOrderLog 保存订单日志
		SaveNormalSubOrderLog(*OrderLog) error
		// GetSubOrder 获取子订单,todo: 删除
		GetSubOrder(id int64) *NormalSubOrder
		// GetSubOrderByOrderNo 根据子订单
		GetSubOrderByOrderNo(s string) ISubOrder
		// SaveOrderItem 保存子订单
		SaveSubOrder(value *NormalSubOrder) (int, error)
		// SaveOrderItem 保存子订单的商品项,并返回编号和错误
		SaveOrderItem(subOrderId int64, value *SubOrderItem) (int32, error)
		// UpdateSubOrderId 更新子订单商品项订单编号为买家订单编号
		UpdateSubOrderId(subOrderId int64) error
		// UpdateSubOrderId 删除子订单
		DeleteSubOrder(subOrderId int64) error
		// DeleteSubOrderItems 删除子订单商品
		DeleteSubOrderItems(subOrderId int64) error

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

		// SaveRebateList 保存订单返利
		SaveOrderRebate(v *AffliteRebate) (int, error)
	}

	// SubmitData 订单提交数据
	SubmitOrderData struct {
		// 买家编号
		BuyerId int64
		// 订单类型
		Type OrderType
		// 收货地址编号
		AddressId int64
		// 优惠券
		CouponCode string
		// 是否余额支付
		BalanceDiscount bool
		// 返利推广人代码
		AffliteCode string
		// 提交的订单数据
		PostedData IPostedData
	}
)
