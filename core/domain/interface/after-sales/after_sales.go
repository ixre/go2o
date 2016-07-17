/**
 * Copyright 2015 @ z3q.net.
 * name : after_sales
 * author : jarryliu
 * date : 2016-07-16 14:41
 * description :
 * history :
 */
package afterSales

import (
	"go2o/core/domain/interface/order"
	"go2o/core/infrastructure/domain"
)

const (
	// 等待商户确认
	StatAwaitingVendor = 1 + iota
	// 商户拒绝售后
	StatDeclined
	// 调解状态
	StatIntercede
	// 同意,等待退货
	StatAwaitingReturnShip
	// 已发货
	StatReturnShipped
	// 已收货,等待系统确认
	StatAwaitingConfirm
	// 售后单已完成
	StatCompleted
	// 售后单已取消
	StatCancelled
)

const (
	// 退款申请(部分退款)
	TypeRefund = 1 + iota
	// 退货
	TypeReturn
	// 换货
	TypeExchange
	// 服务/维修
	TypeService
)

var (
	// 不需要平台确认的状态
	IgnoreConfirmStats = []int{
		TypeExchange,
		TypeService,
	}
	ErrAfterSalesOrderCompleted *domain.DomainError = domain.NewDomainError(
		"err_after_sales_order_completed", "售后单已完成,无法进行操作!")

	ErrUnusualStat *domain.DomainError = domain.NewDomainError(
		"err_after_sales_order_unusual_stat", "不合法的售后单状态")

	ErrNoSuchOrderItem *domain.DomainError = domain.NewDomainError(
		"err_after_sales_order_no_such_order_item", "订单中不包括该商品")

	ErrItemOutOfQuantity *domain.DomainError = domain.NewDomainError(
		"err_after_sales_order_out_of_quantity", "商品超出最大数量")
)

type (

	// 售后单
	IAfterSalesOrder interface {
		// 获取领域编号
		GetDomainId() int

		// 获取订单
		GetOrder() order.ISubOrder

		// 设置要退回货物信息
		SetItem(itemId int, quantity int) error

		// 提交售后申请
		Submit() error

		// 取消申请
		Cancel() error

		// 拒绝售后服务
		Decline(reason string) error

		// 同意售后服务
		Agree() error

		// 退回商品
		ReturnShip(spName string, spOrder string, image string) error

		// 收货, 在商品已退回或尚未发货情况下(线下退货),可以执行此操作
		ReturnReceive() error

		// 系统确认
		Confirm() error

		// 申请调解,只有在商户拒绝后才能申请
		RequestIntercede() error
	}

	IAfterSalesRep interface {
		// 创建退款单
		CreateRefundOrder(v *RefundOrder) IRefundOrder

		// 获取退款单
		GetRefundOrder(id int) IRefundOrder

		// 获取订单的退款单
		GetRefundOrders(orderId int) []IRefundOrder
	}

	// 售后单
	AfterSalesOrder struct {
		// 编号
		Id int `db:"id"`
		// 订单编号
		OrderId int `db:"order_id"`
		// 运营商编号
		VendorId int `db:"vendor_id"`
		// 类型，退货、换货、维修
		Type int `db:"type"`
		// 退货的商品项编号
		ItemId int `db:"item_id"`
		// 商品数量
		Quantity int `db:"quantity"`
		// 售后原因
		Reason string `db:"reason"`
		// 联系人
		PersonName string `db:"person_name"`
		// 联系电话
		PersonPhone string `db:"person_phone"`
		// 退货的快递服务商编号
		ReturnSpName string `db:"rsp_name"`
		// 退货的快递单号
		ReturnSpOrder string `db:"rsp_order"`
		// 退货凭证
		ReturnSpImage string `db:"rsp_image"`
		// 备注(系统)
		Remark string `db:"remark"`
		// 运营商备注
		VendorRemark string `db:"vendor_remark"`
		// 售后单状态
		State int `db:"state"`
		// 提交时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}
)
