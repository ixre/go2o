/**
 * Copyright 2015 @ to2.net.
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
	// 已收货,等待系统确认
	StatAwaitingConfirm
	// 已发货
	StatReturnShipped
	// 已退回
	StateRejected
	// 等待处理
	StateAwaitingProcess
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
	// 不需要平台确认的售后类型
	IgnoreConfirmTypes = []int{
		TypeExchange,
		TypeService,
	}

	ErrNoSuchOrder *domain.DomainError = domain.NewError(
		"err_after_sales_no_such_order", "售后单不存在")

	ErrAutoProcess *domain.DomainError = domain.NewError(
		"err_after_sales_auto_process", "售后单不能人工处理")

	ErrAfterSalesOrderCompleted *domain.DomainError = domain.NewError(
		"err_after_sales_order_completed", "售后单已完成,无法进行操作!")

	ErrAfterSalesRejected *domain.DomainError = domain.NewError(
		"err_after_sales_rejected", "售后单已被系统退回")

	ErrUnusualStat *domain.DomainError = domain.NewError(
		"err_after_sales_order_unusual_stat", "不合法的售后单状态")

	ErrNoSuchOrderItem *domain.DomainError = domain.NewError(
		"err_after_sales_order_no_such_order_item", "订单中不包括该商品")

	ErrOutOfQuantity *domain.DomainError = domain.NewError(
		"err_after_sales_order_out_of_quantity", "超出数量")

	ErrReasonLength *domain.DomainError = domain.NewError(
		"err_after_sales_order_reason_length", "原因不能少于10字")

	ErrNotConfirm *domain.DomainError = domain.NewError(
		"err_after_sales_order_not_confirm", "售后单尚未确认")

	ErrHasCanceled *domain.DomainError = domain.NewError(
		"err_after_sales_order_has_canceled", "售后单已取消")

	ErrOrderAmount *domain.DomainError = domain.NewError(
		"err_after_sales_order_amount", "售后单金额不能为零")

	ErrExchangeNotReceiveItems *domain.DomainError = domain.NewError(
		"err_exhange_order_no_shipping", "订单尚未收货，无法进行换货")

	ErrExchangeOrderNoShipping *domain.DomainError = domain.NewError(
		"err_exhange_order_no_shipping", "换货单尚未发货给顾客")

	ErrNotReceive *domain.DomainError = domain.NewError(
		"err_after_sales_order_not_receive", "尚未收货")

	ErrRefundAfterShipped *domain.DomainError = domain.NewError(
		"err_after_sales_refund_after_shipped", "订单已经发货,只能进行退货或换货操作")

	ErrReturnAfterReceived *domain.DomainError = domain.NewError(
		"err_after_sales_return_after_received", "订单已确认收货,只能进行换货操作")
)

type (
	// 售后单状态
	Stat int

	// 提交售后  -》 商户同意(拒绝) -》 系统确认(取消),或客服调解 -》
	// 开始发货  -》 确认收货       -》 处理售后     -》 完成

	// 售后单
	IAfterSalesOrder interface {
		// 获取领域编号
		GetDomainId() int32

		// 获取售后单数据
		Value() AfterSalesOrder

		// 获取订单
		GetOrder() order.ISubOrder

		// 设置要退回货物信息
		SetItem(snapshotId int64, quantity int32) error

		// 提交售后申请
		Submit() (int32, error)

		// 取消申请
		Cancel() error

		// 同意售后服务,部分操作在同意后,无需确认
		Agree() error

		// 拒绝售后服务
		Decline(reason string) error

		// 申请调解,只有在商户拒绝后才能申请
		RequestIntercede() error

		// 系统确认,泛化应有不同的实现
		Confirm() error

		// 退回售后单
		Reject(remark string) error

		// 退回商品
		ReturnShip(spName string, spOrder string, image string) error

		// 收货, 在商品已退回或尚未发货情况下(线下退货),可以执行此操作
		ReturnReceive() error

		// 处理售后单,处理完成后将变为已完成
		Process() error
	}

	IAfterSalesRepo interface {
		// 创建售后单
		CreateAfterSalesOrder(v *AfterSalesOrder) IAfterSalesOrder

		// 获取售后单
		GetAfterSalesOrder(id int32) IAfterSalesOrder

		// 获取订单的售后单
		GetAllOfSaleOrder(orderId int64) []IAfterSalesOrder
	}

	// 售后单
	AfterSalesOrder struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 订单编号
		OrderId int64 `db:"order_id"`
		// 运营商编号
		VendorId int32 `db:"vendor_id"`
		// 购买者编号
		BuyerId int64 `db:"buyer_id"`
		// 类型，退货、换货、维修
		Type int `db:"type"`
		// 退货的商品项编号
		SnapshotId int64 `db:"snap_id"`
		// 商品数量
		Quantity int32 `db:"quantity"`
		// 售后原因
		Reason string `db:"reason"`
		// 商品售后图片凭证
		ImageUrl string `db:"image_url"`
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
		// 售后单数据,如退款单、退货单、换货单等
		Data interface{} `db:"-"`
		// 订单状态
		StateText string `db:"-"`
	}

	AfterSalesOrderNew struct {
		OrderNo         string `db:"order_no"`
		ItemsInfo       string `db:"items_info"`
		BuyerName       string `db:"name"`
		BuyerId         string `db:"user"`
		ConsigneePerson string `db:"consignee_person"` //收货人
		ConsigneePhone  string `db:"consignee_phone"`  //电话
		ShippingAddress string `db:"shipping_address"` //地址
		ItemAmount      int    `db:"item_amount"`
		FinalAmount     int    `db:"final_amount"`
		IsPaid          int    `db:"is_paid"`
	}
)

// 返回售后状态说明
func (s Stat) String() string {
	switch s {
	case StatAwaitingVendor:
		return "等待商户确认"
	case StatDeclined:
		return "商户拒绝"
	case StatIntercede:
		return "客服处理中"
	case StatAwaitingReturnShip:
		return "等待退货"
	case StatReturnShipped:
		return "等待商户收货"
	case StatAwaitingConfirm:
		return "系统处理中"
	case StatCompleted:
		return "已完成"
	case StatCancelled:
		return "已取消"
	}
	return "-"
}
