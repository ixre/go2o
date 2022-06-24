/**
 * Copyright 2015 @ 56x.net.
 * name : after_sales
 * author : jarryliu
 * date : 2016-07-17 11:42
 * description :
 * history :
 */
package afterSales

import (
	"errors"
	"strings"
	"time"

	afterSales "github.com/ixre/go2o/core/domain/interface/after-sales"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/tmp"
	"github.com/ixre/gof/db/orm"
)

var _ afterSales.IAfterSalesOrder = new(afterSalesOrderImpl)

type afterSalesOrderImpl struct {
	value       *afterSales.AfterSalesOrder
	repo        afterSales.IAfterSalesRepo
	order       order.ISubOrder
	orderRepo   order.IOrderRepo
	paymentRepo payment.IPaymentRepo
}

func NewAfterSalesOrder(v *afterSales.AfterSalesOrder,
	rep afterSales.IAfterSalesRepo, orderRepo order.IOrderRepo,
	memberRepo member.IMemberRepo, paymentRepo payment.IPaymentRepo) afterSales.IAfterSalesOrder {
	as := newAfterSalesOrder(v, rep, orderRepo, paymentRepo)
	switch v.Type {
	case afterSales.TypeReturn:
		return newReturnOrderImpl(as, memberRepo, paymentRepo)
	case afterSales.TypeExchange:
		return newExchangeOrderImpl(as)
	case afterSales.TypeRefund:
		return newRefundOrder(as, memberRepo, paymentRepo)
	}
	panic(errors.New("不支持的售后单类型"))
}

func newAfterSalesOrder(v *afterSales.AfterSalesOrder,
	rep afterSales.IAfterSalesRepo, orderRepo order.IOrderRepo,
	paymentRepo payment.IPaymentRepo) *afterSalesOrderImpl {
	return &afterSalesOrderImpl{
		value:       v,
		repo:        rep,
		orderRepo:   orderRepo,
		paymentRepo: paymentRepo,
	}
}

// 获取领域编号
func (a *afterSalesOrderImpl) GetDomainId() int32 {
	return a.value.Id
}

// 获取售后单数据
func (a *afterSalesOrderImpl) Value() afterSales.AfterSalesOrder {
	return *a.value
}

func (a *afterSalesOrderImpl) saveAfterSalesOrder() error {
	if a.value.OrderId <= 0 {
		panic(errors.New("售后单没有绑定订单"))
	}
	if a.value.SnapshotId <= 0 || a.value.Quantity <= 0 {
		panic(errors.New("售后单缺少商品"))
	}
	a.value.OrderNo = a.repo.GetFreeOrderNo(a.value.OrderId)
	a.value.UpdateTime = time.Now().Unix()
	id, err := orm.I32(orm.Save(tmp.Orm, a.value, int(a.GetDomainId())))
	if err == nil {
		a.value.Id = id
	}
	return err
}

// 获取订单
func (a *afterSalesOrderImpl) GetOrder() order.ISubOrder {
	if a.order == nil {
		if a.value.OrderId > 0 {
			a.order = a.orderRepo.Manager().GetSubOrder(int64(a.value.OrderId))
		}
		if a.order == nil {
			panic(errors.New("售后单对应的订单不存在"))
		}
	}
	return a.order
}

// 设置要退回货物信息
func (a *afterSalesOrderImpl) SetItem(snapshotId int64, quantity int32) error {
	for _, v := range a.GetOrder().Items() {
		if v.SnapshotId == snapshotId {
			// 判断是否超过数量
			if v.Quantity < quantity {
				return afterSales.ErrOutOfQuantity
			}
			// 设置退回商品
			a.value.SnapshotId = snapshotId
			a.value.Quantity = quantity
			return nil
		}
	}
	return afterSales.ErrNoSuchOrderItem
}

// 提交售后申请
func (a *afterSalesOrderImpl) Submit() (int32, error) {
	if a.GetDomainId() > 0 {
		panic(errors.New("售后单已提交"))
	}
	// 售后单未包括商品项
	if a.value.SnapshotId <= 0 || a.value.Quantity <= 0 {
		return 0, afterSales.ErrNoSuchOrderItem
	}
	a.value.Reason = strings.TrimSpace(a.value.Reason)
	if len(a.value.Reason) < 10 {
		return 0, afterSales.ErrReasonLength
	}
	ov := a.GetOrder().GetValue()
	a.value.VendorId = ov.VendorId
	a.value.BuyerId = ov.BuyerId
	a.value.Status = afterSales.StatAwaitingVendor
	a.value.CreateTime = time.Now().Unix()
	err := a.saveAfterSalesOrder()
	return a.GetDomainId(), err
}

// 取消申请
func (a *afterSalesOrderImpl) Cancel() error {
	if a.value.Status == afterSales.StatCompleted {
		return afterSales.ErrAfterSalesOrderCompleted
	}
	if a.value.Status == afterSales.StatCancelled {
		return afterSales.ErrHasCanceled
	}
	a.value.Status = afterSales.StatCancelled
	return a.saveAfterSalesOrder()
}

// 同意售后服务
func (a *afterSalesOrderImpl) Agree() error {
	a.value.Status = afterSales.StatAwaitingVendor
	if a.value.Status != afterSales.StatAwaitingVendor {
		return afterSales.ErrUnusualStat
	}
	// 判断是否需要审核
	needConfirm := true
	for _, v := range afterSales.IgnoreConfirmTypes {
		if a.value.Type == v {
			needConfirm = false
			break
		}
	}
	// 设置为待审核状态
	a.value.Status = afterSales.StatAwaitingConfirm
	// 需要审核
	if needConfirm {
		return a.saveAfterSalesOrder()
	}
	// 不需要审核,直接审核通过
	return a.Confirm()
}

// 拒绝售后服务
func (a *afterSalesOrderImpl) Decline(reason string) error {
	if a.value.Status != afterSales.StatAwaitingVendor {
		return afterSales.ErrUnusualStat
	}
	a.value.Status = afterSales.StatDeclined
	a.value.VendorRemark = reason
	return a.saveAfterSalesOrder()
}

// 申请调解,只有在商户拒绝后才能申请
func (a *afterSalesOrderImpl) RequestIntercede() error {
	if a.value.Status != afterSales.StatDeclined {
		return afterSales.ErrUnusualStat
	}
	a.value.Status = afterSales.StatIntercede
	return a.saveAfterSalesOrder()
}

// 系统确认
func (a *afterSalesOrderImpl) Confirm() error {
	if a.value.Status == afterSales.StatCompleted {
		return afterSales.ErrAfterSalesOrderCompleted
	}
	if a.value.Status == afterSales.StatRejected {
		return afterSales.ErrAfterSalesRejected
	}
	if a.value.Status != afterSales.StatAwaitingConfirm {
		return afterSales.ErrUnusualStat
	}
	// 退款,不需要退货,直接进入完成状态
	if a.value.Type == afterSales.TypeRefund {
		return a.awaitingProcess()
	}
	a.value.Status = afterSales.StatAwaitingReturnShip
	return a.saveAfterSalesOrder()
}

// 退回售后单
func (a *afterSalesOrderImpl) Reject(remark string) error {
	if a.value.Status == afterSales.StatCompleted {
		return afterSales.ErrAfterSalesOrderCompleted
	}
	if a.value.Status != afterSales.StatAwaitingConfirm {
		return afterSales.ErrUnusualStat
	}
	a.value.Remark = remark
	a.value.Status = afterSales.StatRejected
	return a.saveAfterSalesOrder()
}

// 退回商品
func (a *afterSalesOrderImpl) ReturnShip(spName string, spOrder string, image string) error {
	if a.value.Status != afterSales.StatAwaitingReturnShip {
		return afterSales.ErrUnusualStat
	}
	a.value.ShipmentExpress = spName
	a.value.ShipmentOrderNo = spOrder
	a.value.ShipmentImage = image
	a.value.Status = afterSales.StatReturnShipped
	return a.saveAfterSalesOrder()
}

// 收货, 在商品已退回或尚未发货情况下(线下退货),可以执行此操作
func (a *afterSalesOrderImpl) ReturnReceive() error {
	if a.value.Status != afterSales.StatAwaitingReturnShip &&
		a.value.Status != afterSales.StatReturnShipped {
		return afterSales.ErrUnusualStat
	}
	return a.awaitingProcess()
}

// 等待处理
func (a *afterSalesOrderImpl) awaitingProcess() error {
	if a.value.Status == afterSales.StatCompleted {
		return afterSales.ErrAfterSalesOrderCompleted
	}
	if a.value.Status == afterSales.StatRejected {
		return afterSales.ErrAfterSalesRejected
	}

	// 判断状态是否正确
	statOK := a.value.Status == afterSales.StatAwaitingReturnShip ||
		a.value.Status == afterSales.StatReturnShipped
	if !statOK && a.value.Type == afterSales.TypeRefund {
		statOK = a.value.Status == afterSales.StatAwaitingConfirm
	}
	if !statOK {
		return afterSales.ErrNotConfirm
	}

	// 等待处理
	a.value.Status = afterSales.StatAwaitingProcess
	return a.saveAfterSalesOrder()
}

// 处理售后单,处理完成后将变为已完成
func (a *afterSalesOrderImpl) Process() error {
	if a.value.Status != afterSales.StatAwaitingProcess {
		return afterSales.ErrUnusualStat
	}
	a.value.Status = afterSales.StatCompleted
	return a.saveAfterSalesOrder()
}
