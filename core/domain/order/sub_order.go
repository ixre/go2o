package order

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/shipment"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/gof/util"
)

//todo: ?? 自动收货功能

var _ order.ISubOrder = new(subOrderImpl)

// 子订单实现
type subOrderImpl struct {
	value           *order.NormalSubOrder
	parent          order.IOrder
	buyer           member.IMember
	internalSuspend bool //内部挂起
	paymentRepo     payment.IPaymentRepo
	repo            order.IOrderRepo
	memberRepo      member.IMemberRepo
	itemRepo        item.IItemRepo
	productRepo     product.IProductRepo
	manager         order.IOrderManager
	shipRepo        shipment.IShipmentRepo
	valRepo         valueobject.IValueRepo
	mchRepo         merchant.IMerchantRepo
	registryRepo    registry.IRegistryRepo
}

func NewSubNormalOrder(v *order.NormalSubOrder,
	manager order.IOrderManager, rep order.IOrderRepo,
	mmRepo member.IMemberRepo, goodsRepo item.IItemRepo,
	shipRepo shipment.IShipmentRepo, productRepo product.IProductRepo,
	paymentRepo payment.IPaymentRepo, valRepo valueobject.IValueRepo,
	mchRepo merchant.IMerchantRepo, registryRepo registry.IRegistryRepo) order.ISubOrder {
	return &subOrderImpl{
		value:        v,
		manager:      manager,
		repo:         rep,
		memberRepo:   mmRepo,
		itemRepo:     goodsRepo,
		productRepo:  productRepo,
		shipRepo:     shipRepo,
		paymentRepo:  paymentRepo,
		valRepo:      valRepo,
		mchRepo:      mchRepo,
		registryRepo: registryRepo,
	}
}

// GetDomainId 获取领域对象编号
func (o *subOrderImpl) GetDomainId() int64 {
	return o.value.Id
}

// GetValue 获取值对象
func (o *subOrderImpl) GetValue() *order.NormalSubOrder {
	return o.value
}

func parseDetailValue(subOrder order.ISubOrder) *order.ComplexOrderDetails {
	v := subOrder.GetValue()
	dst := &order.ComplexOrderDetails{
		Id:             subOrder.GetDomainId(),
		OrderNo:        v.OrderNo,
		ShopId:         v.ShopId,
		ShopName:       v.ShopName,
		ItemAmount:     v.ItemAmount,
		DiscountAmount: v.DiscountAmount,
		ExpressFee:     v.ExpressFee,
		PackageFee:     v.PackageFee,
		FinalAmount:    v.FinalAmount,
		BuyerComment:   v.BuyerComment,
		Status:         v.Status,
		StatusText:     "",
		Items:          []*order.ComplexItem{},
		UpdateTime:     v.UpdateTime,
	}
	impl := subOrder.(*subOrderImpl)
	for _, v := range subOrder.Items() {
		dst.Items = append(dst.Items, impl.parseComplexItem(v))
	}
	return dst
}

// Complex 复合的订单信息
func (o *subOrderImpl) Complex() *order.ComplexOrder {
	bo := o.baseOrder()
	if bo != nil {
		co := o.baseOrder().Complex()
		co.Details = []*order.ComplexOrderDetails{parseDetailValue(o)}
		return co
	}
	return nil
}

// 转换订单商品
func (o *subOrderImpl) parseComplexItem(i *order.SubOrderItem) *order.ComplexItem {
	snap := o.itemRepo.GetSalesSnapshot(i.SnapshotId)
	it := &order.ComplexItem{
		ID:             i.ID,
		ItemId:         i.ItemId,
		SkuId:          i.SkuId,
		SkuWord:        snap.Sku,
		SnapshotId:     i.SnapshotId,
		ItemTitle:      snap.GoodsTitle,
		MainImage:      snap.Image,
		Price:          snap.Price,
		FinalPrice:     snap.Price,
		Quantity:       i.Quantity,
		ReturnQuantity: i.ReturnQuantity,
		Amount:         i.Amount,
		FinalAmount:    i.FinalAmount,
		IsShipped:      i.IsShipped,
		Data:           make(map[string]string),
	}
	base := o.baseOrder().(*normalOrderImpl)
	base.baseOrderImpl.bindItemInfo(it)
	return it
}

// 获取商品项
func (o *subOrderImpl) Items() []*order.SubOrderItem {
	if (o.value.Items == nil || len(o.value.Items) == 0) &&
		o.GetDomainId() > 0 {
		o.value.Items = o.repo.GetSubOrderItems(o.GetDomainId())
	}
	return o.value.Items
}

// 获取订单
func (o *subOrderImpl) baseOrder() order.IOrder {
	if o.parent == nil {
		o.parent = o.manager.GetOrderById(o.value.OrderId)
	}
	return o.parent
}

// 获取购买的会员
func (o *subOrderImpl) getBuyer() member.IMember {
	return o.baseOrder().Buyer()
}

// 添加备注
func (o *subOrderImpl) AddRemark(remark string) {
	o.value.Remark = remark
}

func (o *subOrderImpl) saveOrderItems() error {
	unix := time.Now().Unix()
	id := o.GetDomainId()
	for _, v := range o.Items() {
		if v.OrderId == 0 {
			v.OrderId = id
		}
		v.SellerOrderId = id
		v.UpdateTime = unix
		_, err := o.repo.SaveOrderItem(id, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// 提交子订单
func (o *subOrderImpl) Submit() (int64, error) {
	if o.GetDomainId() > 0 {
		panic("suborder is created!")
	}
	if o.value.CreateTime <= 0 {
		unix := time.Now().Unix()
		o.value.CreateTime = unix
		o.value.UpdateTime = unix
	}
	id, err := util.I64Err(o.repo.SaveSubOrder(o.value))
	if err == nil {
		o.value.Id = id
		err = o.saveOrderItems()
		o.AppendLog(order.LogSetup, true, "{created}")
	}
	return id, err
}

// 保存订单
func (o *subOrderImpl) saveSubOrder() error {
	unix := time.Now().Unix()
	o.value.UpdateTime = unix
	if o.GetDomainId() <= 0 {
		panic("please use Submit() to create new suborder!")
	}
	_, err := o.repo.SaveSubOrder(o.value)
	if err == nil {
		o.syncOrderState()
	}
	return err
}

// 同步订单状态
func (o *subOrderImpl) syncOrderState() {
	if bo := o.baseOrder(); bo != nil {
		oi := bo.(*normalOrderImpl).baseOrderImpl
		if oi.State() != order.StatBreak {
			oi.saveOrderState(order.OrderStatus(o.value.Status))
		}
	}

}

// 订单完成支付
func (o *subOrderImpl) orderFinishPaid() error {
	if o.value.Status > order.StatAwaitingPayment {
		return order.ErrOrderPayed
	}
	if o.value.Status == order.StatAwaitingPayment {
		o.value.Status = order.StatAwaitingPickup
		// 更新拆分状态
		if o.value.BreakStatus == order.BreakAwaitBreak {
			o.value.BreakStatus = order.Breaked
		}
		err := o.AppendLog(order.LogSetup, true, "{finish_pay}")
		if err == nil {
			err = o.saveSubOrder()
		}
		return err
	}
	return order.ErrUnusualOrderStat
}

// 在线支付交易完成
func (o *subOrderImpl) PaymentFinishByOnlineTrade() error {
	return o.orderFinishPaid()
}

// 挂起
func (o *subOrderImpl) Suspend(reason string) error {
	o.value.IsSuspend = 1
	o.internalSuspend = true
	o.value.UpdateTime = time.Now().Unix()
	err := o.saveSubOrder()
	if err == nil {
		err = o.AppendLog(order.LogSetup, true, "订单已锁定"+reason)
	}
	return err
}

// 添加日志
func (o *subOrderImpl) AppendLog(logType order.LogType, system bool, message string) error {
	if o.GetDomainId() <= 0 {
		return errors.New("order not created.")
	}
	var systemInt int
	if system {
		systemInt = 1
	} else {
		systemInt = 0
	}
	l := &order.OrderLog{
		OrderId:    o.GetDomainId(),
		Type:       int(logType),
		IsSystem:   systemInt,
		OrderState: int(o.value.Status),
		Message:    message,
		RecordTime: time.Now().Unix(),
	}
	return o.repo.SaveNormalSubOrderLog(l)
}

// 确认订单
func (o *subOrderImpl) Confirm() (err error) {
	//todo: 线下交易,自动确认
	//if o._value.PaymentOpt == enum.PaymentOnlinePay &&
	//o._value.IsPaid == enum.FALSE {
	//    return order.ErrOrderNotPayed
	//}
	if o.value.Status < order.StatAwaitingConfirm {
		return order.ErrOrderNotPayed
	}
	if o.value.Status >= order.StatAwaitingPickup {
		return order.ErrOrderHasConfirm
	}
	o.value.Status = order.StatAwaitingPickup
	o.value.UpdateTime = time.Now().Unix()
	err = o.saveSubOrder()
	if err == nil {
		go o.addItemSalesNum()
		err = o.AppendLog(order.LogSetup, false, "{confirm}")
	}
	return err
}

// 增加商品的销售数量
func (o *subOrderImpl) addItemSalesNum() {
	//log.Println("---订单：",o.value.OrderNo," 商品：",len(o.Items()))
	for _, v := range o.Items() {
		it := o.itemRepo.GetItem(v.ItemId)
		err := it.AddSalesNum(v.SkuId, v.Quantity)
		if err != nil {
			log.Println("---增加销售数量：", v.ItemId,
				" sku:", v.SkuId, " error:", err.Error())
		}
	}
}

// PickUp 捡货(备货)
func (o *subOrderImpl) PickUp() error {
	if o.value.Status < order.StatAwaitingPickup {
		return order.ErrOrderNotConfirm
	}
	if o.value.Status >= order.StatAwaitingShipment {
		return order.ErrOrderHasPickUp
	}
	o.value.Status = order.StatAwaitingShipment
	o.value.UpdateTime = time.Now().Unix()
	err := o.saveSubOrder()
	if err == nil {
		err = o.AppendLog(order.LogSetup, true, "{pickup}")
	}
	return err
}

// Ship 发货
func (o *subOrderImpl) Ship(spId int32, spOrder string) error {
	//so := o._shipRepo.GetOrders()
	//todo: 可进行发货修改
	if o.value.Status >= order.StatShipped {
		return order.ErrOrderShipped
	}
	// 如果没有备货完成,则发货前自动完成备货
	if o.value.Status < order.StatAwaitingShipment {
		o.value.Status = order.StatAwaitingShipment
		o.value.UpdateTime = time.Now().Unix()
		_ = o.AppendLog(order.LogSetup, true, "{pickup}")
		//return order.ErrOrderNotPickUp
	}
	if list := o.shipRepo.GetShipOrders(o.GetDomainId(), true); len(list) > 0 {
		return order.ErrPartialShipment
	}
	if spId <= 0 || spOrder == "" {
		return shipment.ErrMissingSpInfo
	}

	so := o.createShipmentOrder(o.Items())
	if so == nil {
		return order.ErrUnusualOrder
	}
	// 生成发货单并发货
	err := so.Ship(spId, spOrder)
	if err == nil {
		o.value.Status = order.StatShipped
		o.value.UpdateTime = time.Now().Unix()
		err = o.saveSubOrder()
		if err == nil {
			// 保存商品的发货状态
			err = o.saveOrderItems()
			_ = o.AppendLog(order.LogSetup, true, "{shipped}")
		}
	}
	return err
}

func (o *subOrderImpl) createShipmentOrder(items []*order.SubOrderItem) shipment.IShipmentOrder {
	if items == nil || len(items) == 0 {
		return nil
	}
	unix := time.Now().Unix()
	orderId := o.baseOrder().GetAggregateRootId()
	subOrderId := o.GetDomainId()
	so := &shipment.ShipmentOrder{
		ID:          0,
		OrderId:     orderId,
		SubOrderId:  subOrderId,
		ShipmentLog: "",
		ShipTime:    unix,
		State:       shipment.StatAwaitingShipment,
		UpdateTime:  unix,
		Items:       []*shipment.ShipmentItem{},
	}
	for _, v := range items {
		if v.IsShipped == 1 {
			continue
		}
		so.Amount += float64(v.Amount)
		so.FinalAmount += float64(v.FinalAmount)
		so.Items = append(so.Items, &shipment.ShipmentItem{
			ID:          0,
			SnapshotId:  v.SnapshotId,
			Quantity:    v.Quantity,
			Amount:      float64(v.Amount),
			FinalAmount: float64(v.FinalAmount),
		})
		v.IsShipped = 1
	}
	return o.shipRepo.CreateShipmentOrder(so)
}

// 已收货
func (o *subOrderImpl) BuyerReceived() error {
	if o.value.Status < order.StatShipped {
		return order.ErrOrderNotShipped
	}
	if o.value.Status >= order.StatCompleted {
		return order.ErrIsCompleted
	}
	dt := time.Now()
	o.value.Status = order.StatCompleted
	o.value.UpdateTime = dt.Unix()
	o.value.IsSuspend = 0
	err := o.saveSubOrder()
	if err == nil {
		err = o.AppendLog(order.LogSetup, true, "{completed}")
		if err == nil {
			go o.vendorSettle()
			// 执行其他的操作
			if err2 := o.onOrderComplete(); err != nil {
				domain.HandleError(err2, "domain")
			}
		}
	}
	return err
}

func (o *subOrderImpl) getOrderAmount() (amount int64, refund int64) {
	items := o.Items()
	for _, item := range items {
		if item.ReturnQuantity > 0 {
			a := item.Amount / int64(item.Quantity) * int64(item.ReturnQuantity)
			if item.ReturnQuantity != item.Quantity {
				amount += item.Amount - a
			}
			refund += a
		} else {
			amount += item.Amount
		}
	}
	//如果非全部退货、退款,则加上运费及包装费
	if amount > 0 {
		amount += o.value.ExpressFee + o.value.PackageFee
	}
	return amount, refund
}

// 获取订单的成本
func (o *subOrderImpl) getOrderCost() int64 {
	var cost int64
	items := o.Items()
	for _, item := range items {
		snap := o.itemRepo.GetSalesSnapshot(item.SnapshotId)
		cost += snap.Cost * int64(item.Quantity-item.ReturnQuantity)
	}
	//如果非全部退货、退款,则加上运费及包装费
	if cost > 0 {
		cost += o.value.ExpressFee + o.value.PackageFee
	}
	return cost
}

// 商户结算
func (o *subOrderImpl) vendorSettle() error {
	vendor := o.mchRepo.GetMerchant(int(o.value.VendorId))
	if vendor != nil {
		settleMode := o.registryRepo.Get(registry.MchOrderSettleMode).IntValue()
		switch enum.MchSettleMode(settleMode) {
		case enum.MchModeSettleByCost:
			return o.vendorSettleByCost(vendor)
		case enum.MchModeSettleByRate:
			return o.vendorSettleByRate(vendor)
		case enum.MchModeSettleByOrderQuantity:
			return o.vendorSettleByOrderQuantity(vendor)
		}
	}
	return nil
}

// 根据供货价进行商户结算
func (o *subOrderImpl) vendorSettleByCost(vendor merchant.IMerchant) error {
	_, refund := o.getOrderAmount()
	sAmount := o.getOrderCost()
	if sAmount > 0 {
		totalAmount := int(float32(sAmount) * float32(enum.RATE_AMOUNT))
		refundAmount := int(float32(refund) * float32(enum.RATE_AMOUNT))
		tradeFee, _ := vendor.SaleManager().MathTradeFee(
			merchant.TKNormalOrder, totalAmount)
		return vendor.Account().SettleOrder(o.value.OrderNo,
			totalAmount, tradeFee, refundAmount, "零售订单结算")
	}
	return nil
}

// 根据比例进行商户结算
func (o *subOrderImpl) vendorSettleByRate(vendor merchant.IMerchant) error {
	rate := o.registryRepo.Get(registry.MchOrderSettleRate).FloatValue()
	amount, refund := o.getOrderAmount()
	sAmount := int64(float64(amount) * rate)
	if sAmount > 0 {
		totalAmount := int(float32(sAmount) * float32(enum.RATE_AMOUNT))
		refundAmount := int(float32(refund) * float32(enum.RATE_AMOUNT))
		tradeFee, _ := vendor.SaleManager().MathTradeFee(
			merchant.TKNormalOrder, totalAmount)
		return vendor.Account().SettleOrder(o.value.OrderNo,
			totalAmount, tradeFee, refundAmount, "零售订单结算")

	}
	return nil
}
func (o *subOrderImpl) vendorSettleByOrderQuantity(vendor merchant.IMerchant) error {
	fee := o.registryRepo.Get(registry.MchSingleOrderServiceFee).FloatValue()
	amount, refund := o.getOrderAmount()
	if fee > 0 {
		totalAmount := int(math.Min(float64(amount), fee) * float64(enum.RATE_AMOUNT))
		refundAmount := int(float32(refund) * float32(enum.RATE_AMOUNT))
		tradeFee, _ := vendor.SaleManager().MathTradeFee(
			merchant.TKNormalOrder, totalAmount)
		return vendor.Account().SettleOrder(o.value.OrderNo,
			totalAmount, tradeFee, refundAmount, "零售订单结算")

	}
	return nil
}

// 获取订单的日志
func (o *subOrderImpl) LogBytes() []byte {
	buf := bytes.NewBufferString("")
	list := o.repo.GetSubOrderLogs(o.GetDomainId())
	for _, v := range list {
		buf.WriteString(time.Unix(v.RecordTime, 0).Format("2006-01-02 15:04:05"))
		buf.WriteString("  ")
		if v.Message[:1] == "{" {
			if msg := o.getLogStringByStat(v.OrderState); len(msg) > 0 {
				v.Message = msg
			}
		}
		buf.WriteString(v.Message)
		buf.Write([]byte("\n"))
	}
	return buf.Bytes()
}

func (o *subOrderImpl) getLogStringByStat(stat int) string {
	switch stat {
	case order.StatAwaitingPayment:
		return "订单已提交..."
	case order.StatAwaitingConfirm:
		return "订单已支付,等待商户确认。"
	case order.StatAwaitingPickup:
		return "订单已确认,备货中..."
	case order.StatAwaitingShipment:
		return "备货完成,即将发货。"
	case order.StatShipped:
		return "订单已发货,请等待收货。"
	case order.StatCompleted:
		return "已收货,订单完成。"
	}
	return ""
}

// 更新账户
func (o *subOrderImpl) updateAccountForOrder(m member.IMember) error {
	if o.value.Status != order.StatCompleted {
		return order.ErrUnusualOrderStat
	}
	var err error
	ov := o.value
	amount := ov.FinalAmount
	acc := m.GetAccount()

	// 增加经验
	expEnabled := o.registryRepo.Get(registry.ExperienceEnabled).BoolValue()
	if expEnabled {
		rate := o.registryRepo.Get(registry.ExperienceRateByOrder).FloatValue()
		if exp := int(float64(amount) * rate); exp > 0 {
			if err = m.AddExp(exp); err != nil {
				return err
			}
		}
	}

	// 增加积分
	//todo: 增加阶梯的返积分,比如订单满30送100积分, 不考虑额外赠送,额外的当做补贴
	rate := o.registryRepo.Get(registry.IntegralRateByWholesaleOrder).FloatValue()
	integral := int(float64(amount) * rate)
	// 赠送积分
	if integral > 0 {
		_, err = m.GetAccount().CarryTo(member.AccountIntegral,
			member.AccountOperateData{
				Title:   "购物消费赠送积分",
				Amount:  integral,
				OuterNo: o.value.OrderNo,
				Remark:  "sys",
			}, false, 0)
		if err != nil {
			return err
		}
	}
	acv := acc.GetValue()
	acv.TotalExpense += ov.ItemAmount
	acv.TotalPay += ov.FinalAmount
	acv.UpdateTime = time.Now().Unix()
	_, err = acc.Save()
	return err
}

// 取消订单
func (o *subOrderImpl) Cancel(reason string) error {
	if o.value.Status == order.StatCancelled {
		return order.ErrOrderCancelled
	}
	// 已发货订单无法取消
	if o.value.Status >= order.StatShipped {
		return order.ErrOrderShippedCancel
	}

	o.value.Status = order.StatCancelled
	o.value.UpdateTime = time.Now().Unix()
	err := o.saveSubOrder()
	if err == nil {
		domain.HandleError(o.AppendLog(order.LogSetup, true, reason), "domain")
		// 取消支付单
		err = o.cancelPaymentOrder()
		if err == nil {
			// 取消商品
			err = o.cancelGoods()
		}
	}
	return err
}

// 取消商品
func (o *subOrderImpl) cancelGoods() error {
	for _, v := range o.Items() {
		snapshot := o.itemRepo.GetSalesSnapshot(v.SnapshotId)
		if snapshot == nil {
			return item.ErrNoSuchSnapshot
		}
		gds := o.itemRepo.GetItem(snapshot.SkuId)
		if gds != nil {
			// 释放库存
			gds.ReleaseStock(v.SkuId, v.Quantity)
			// 如果订单已付款，则取消销售数量
			if o.value.Status > order.StatAwaitingPayment {
				gds.CancelSale(v.SkuId, v.Quantity, o.value.OrderNo)
			}
		}
	}
	return nil
}

// 取消支付单
func (o *subOrderImpl) cancelPaymentOrder() error {
	od := o.baseOrder()
	if od.Type() != order.TRetail {
		panic("not support order type")
	}
	ip := od.GetPaymentOrder()
	if ip != nil {
		return ip.Cancel() //todo: there have a bug, when other order has shipmented. all sub order will be cancelled.
	}
	return nil
}

// 退回商品
func (o *subOrderImpl) Return(snapshotId int64, quantity int32) error {
	for _, v := range o.Items() {
		if v.SnapshotId == snapshotId {
			if v.Quantity-v.ReturnQuantity < quantity {
				return order.ErrOutOfQuantity
			}
			v.ReturnQuantity += quantity
			_, err := o.repo.SaveOrderItem(o.GetDomainId(), v)
			return err
		}
	}
	return order.ErrNoSuchGoodsOfOrder
}

// 撤销退回商品
func (o *subOrderImpl) RevertReturn(snapshotId int64, quantity int32) error {
	for _, v := range o.Items() {
		if v.SnapshotId == snapshotId {
			if v.ReturnQuantity < quantity {
				return order.ErrOutOfQuantity
			}
			v.ReturnQuantity -= quantity
			_, err := o.repo.SaveOrderItem(o.GetDomainId(), v)
			return err
		}
	}
	return order.ErrNoSuchGoodsOfOrder
}

// 申请退款
func (o *subOrderImpl) SubmitRefund(reason string) error {
	if o.value.Status == order.StatAwaitingPayment {
		return o.Cancel("订单主动申请取消,原因:" + reason)
	}
	if o.value.Status >= order.StatShipped ||
		o.value.Status >= order.StatCancelled {
		return order.ErrOrderCancelled
	}
	o.value.Status = order.StatAwaitingCancel
	o.value.UpdateTime = time.Now().Unix()
	return o.saveSubOrder()
}

// 谢绝订单
func (o *subOrderImpl) Decline(reason string) error {
	if o.value.Status == order.StatAwaitingPayment {
		return o.Cancel("商户取消,原因:" + reason)
	}
	if o.value.Status >= order.StatShipped ||
		o.value.Status >= order.StatCancelled {
		return order.ErrOrderCancelled
	}
	o.value.Status = order.StatDeclined
	o.value.UpdateTime = time.Now().Unix()
	return o.saveSubOrder()
}

// 退款 todo: will delete,代码供取消订单参考
func (o *subOrderImpl) refund() error {
	// 已退款
	if o.value.Status == order.StatRefunded ||
		o.value.Status == order.StatCancelled {
		return order.ErrHasRefund
	}
	// 不允许退款
	if o.value.Status != order.StatAwaitingCancel &&
		o.value.Status != order.StatDeclined {
		return order.ErrDisallowRefund
	}
	o.value.Status = order.StatRefunded
	o.value.UpdateTime = time.Now().Unix()
	err := o.saveSubOrder()
	if err == nil {
		err = o.cancelPaymentOrder()
	}
	return err
}

// 取消退款申请
func (o *subOrderImpl) CancelRefund() error {
	if o.value.Status != order.StatAwaitingCancel {
		panic(errors.New("订单已经取消,不允许再退款"))
	}
	o.value.Status = order.StatAwaitingConfirm
	o.value.UpdateTime = time.Now().Unix()
	return o.saveSubOrder()
}

// 完成订单
func (o *subOrderImpl) onOrderComplete() error {
	// 更新发货单
	soList := o.shipRepo.GetShipOrders(o.GetDomainId(), true)
	for _, v := range soList {
		domain.HandleError(v.Completed(), "domain")
	}

	// 获取消费者消息
	m := o.getBuyer()
	if m == nil {
		return member.ErrNoSuchMember
	}

	// 更新会员账户
	err := o.updateAccountForOrder(m)
	if err != nil {
		return err
	}

	// 处理返现
	err = o.handleCashBack()

	return err
}

// Destory 销毁订单
func (o *subOrderImpl) Destory() error {
	if o.value.BreakStatus == order.BreakDefault {
		err := o.repo.UpdateSubOrderId(o.GetDomainId())
		if err == nil {
			err = o.repo.DeleteSubOrder(o.GetDomainId())
		}
		return err
	}
	err := o.repo.DeleteSubOrderItems(o.GetDomainId())
	if err == nil {
		err = o.repo.DeleteSubOrder(o.GetDomainId())
	}
	return err
}

// 更新返现到会员账户
func (o *subOrderImpl) updateShoppingMemberBackFee(mchName string,
	m member.IMember, fee int64, unixTime int64) error {
	if fee <= 0 || math.IsNaN(float64(fee)) {
		return nil
	}
	v := o.GetValue()

	//更新账户
	acc := m.GetAccount()
	//给自己返现
	tit := fmt.Sprintf("订单:%s(商户:%s)返现￥%.2f元", v.OrderNo, mchName, fee)
	_, err := acc.CarryTo(member.AccountWallet,
		member.AccountOperateData{
			Title:   tit,
			Amount:  int(fee * 100),
			OuterNo: o.value.OrderNo,
			Remark:  "sys",
		}, false, 0)
	return err
}
