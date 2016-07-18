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
    // 不需要平台确认的售后类型
    IgnoreConfirmTypes = []int{
        TypeExchange,
        TypeService,
    }
    ErrAfterSalesOrderCompleted *domain.DomainError = domain.NewDomainError(
        "err_after_sales_order_completed", "售后单已完成,无法进行操作!")

    ErrUnusualStat *domain.DomainError = domain.NewDomainError(
        "err_after_sales_order_unusual_stat", "不合法的售后单状态")

    ErrNoSuchOrderItem *domain.DomainError = domain.NewDomainError(
        "err_after_sales_order_no_such_order_item", "订单中不包括该商品")

    ErrOutOfQuantity *domain.DomainError = domain.NewDomainError(
        "err_after_sales_order_out_of_quantity", "超出数量")

    ErrReasonLength *domain.DomainError = domain.NewDomainError(
        "err_after_sales_order_reason_length", "原因不能少于10字")

    ErrNotConfirm *domain.DomainError = domain.NewDomainError(
        "err_after_sales_order_not_confirm", "售后单尚未确认")

    ErrHasCancelled *domain.DomainError = domain.NewDomainError(
        "err_after_sales_order_has_cancelled", "售后单已取消")

    ErrOrderAmount *domain.DomainError = domain.NewDomainError(
        "err_after_sales_order_amount", "售后单金额不能为零")
)

type (
    // 售后单状态
    Stat int

    // 售后单
    IAfterSalesOrder interface {
        // 获取领域编号
        GetDomainId() int

        // 获取售后单数据
        Value() AfterSalesOrder

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

        // 同意售后服务,部分操作在同意后,无需确认
        Agree() error

        // 退回商品
        ReturnShip(spName string, spOrder string, image string) error

        // 收货, 在商品已退回或尚未发货情况下(线下退货),可以执行此操作
        ReturnReceive() error

        // 系统确认,泛化应有不同的实现
        Confirm() error

        // 申请调解,只有在商户拒绝后才能申请
        RequestIntercede() error
    }

    IAfterSalesRep interface {
        // 创建售后单
        CreateAfterSalesOrder(v *AfterSalesOrder) IAfterSalesOrder

        // 获取售后单
        GetAfterSalesOrder(id int) IAfterSalesOrder

        // 获取订单的售后单
        GetAllOfSaleOrder(orderId int) []IAfterSalesOrder
    }

    // 售后单
    AfterSalesOrder struct {
        // 编号
        Id            int `db:"id" pk:"yes" auto:"yes"`
        // 订单编号
        OrderId       int `db:"order_id"`
        // 运营商编号
        VendorId      int `db:"vendor_id"`
        // 购买者编号
        BuyerId       int `db:"buyer_id"`
        // 类型，退货、换货、维修
        Type          int `db:"type"`
        // 退货的商品项编号
        SnapshotId    int `db:"snap_id"`
        // 商品数量
        Quantity      int `db:"quantity"`
        // 售后原因
        Reason        string `db:"reason"`
        // 联系人
        PersonName    string `db:"person_name"`
        // 联系电话
        PersonPhone   string `db:"person_phone"`
        // 退货的快递服务商编号
        ReturnSpName  string `db:"rsp_name"`
        // 退货的快递单号
        ReturnSpOrder string `db:"rsp_order"`
        // 退货凭证
        ReturnSpImage string `db:"rsp_image"`
        // 备注(系统)
        Remark        string `db:"remark"`
        // 运营商备注
        VendorRemark  string `db:"vendor_remark"`
        // 售后单状态
        State         int `db:"state"`
        // 提交时间
        CreateTime    int64 `db:"create_time"`
        // 更新时间
        UpdateTime    int64 `db:"update_time"`

        // 售后单数据,如退款单、退货单、换货单等
        Data          interface{} `db:"-"`
        // 订单状态
        StateText     string `db:"-"`
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
