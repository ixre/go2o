package order

import (
	"bytes"
	"errors"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

var _ order.IOrder = new(wholesaleOrderImpl)
var _ order.IWholesaleOrder = new(wholesaleOrderImpl)

type wholesaleOrderImpl struct {
	*baseOrderImpl
	value        *order.WholesaleOrder
	items        []*orderItem
	realItems    []*order.WholesaleItem
	paymentOrder payment.IPaymentOrder
	orderRepo    order.IOrderRepo
	expressRepo  express.IExpressRepo
	payRepo      payment.IPaymentRepo
	shipRepo     shipment.IShipmentRepo
	itemRepo     item.IGoodsItemRepo
	mchRepo      merchant.IMerchantRepo
	valueRepo    valueobject.IValueRepo
	registryRepo registry.IRegistryRepo
}

func newWholesaleOrder(base *baseOrderImpl,
	shoppingRepo order.IOrderRepo, goodsRepo item.IGoodsItemRepo,
	expressRepo express.IExpressRepo, payRepo payment.IPaymentRepo,
	shipRepo shipment.IShipmentRepo, mchRepo merchant.IMerchantRepo,
	valueRepo valueobject.IValueRepo, registryRepo registry.IRegistryRepo) order.IOrder {
	o := &wholesaleOrderImpl{
		baseOrderImpl: base,
		orderRepo:     shoppingRepo,
		itemRepo:      goodsRepo,
		expressRepo:   expressRepo,
		payRepo:       payRepo,
		shipRepo:      shipRepo,
		mchRepo:       mchRepo,
		valueRepo:     valueRepo,
		registryRepo:  registryRepo,
	}
	return o.init()
}

func (o *wholesaleOrderImpl) init() order.IOrder {
	if o.GetAggregateRootId() <= 0 {
		o.value = &order.WholesaleOrder{
			ID:          0,
			OrderNo:     "",
			OrderId:     0,
			BuyerId:     o.baseValue.BuyerId,
			VendorId:    0,
			ShopId:      0,
			ItemAmount:  0,
			ExpressFee:  0,
			PackageFee:  0,
			FinalAmount: 0,
			State:       o.baseValue.State,
		}
	}
	o.getValue()
	return o
}

func (o *wholesaleOrderImpl) getValue() *order.WholesaleOrder {
	if o.value == nil {
		id := o.GetAggregateRootId()
		if id > 0 {
			o.value = o.repo.GetWholesaleOrder("order_id= $1", id)
		}
	}
	return o.value
}

// 设置商品项
func (o *wholesaleOrderImpl) SetItems(items []*cart.ItemPair) {
	if o.GetAggregateRootId() > 0 {
		panic("wholesale has created. can't use SetItems!")
	}
	o.parseOrder(items)
	// 计算折扣
	o.applyGroupDiscount()
	// 均摊优惠折扣到商品
	o.avgDiscountForItem()
}

// 转换为订单相关对象
func (o *wholesaleOrderImpl) parseOrder(items []*cart.ItemPair) {
	if o.GetAggregateRootId() > 0 {
		panic("订单已经生成，无法解析")
	}
	o.items = []*orderItem{}
	for _, v := range items {
		o.items = append(o.items, o.createItem(v))
	}
	// 获取运营商和商铺编号
	o.value.VendorId = o.items[0].VendorId
	o.value.ShopId = o.items[0].ShopId
	// 运费计算器
	ue := o.expressRepo.GetUserExpress(o.value.VendorId)
	ec := ue.CreateCalculator()
	// 计算订单金额及运费
	for _, it := range o.items {
		o.value.ItemAmount += it.Amount
		o.value.DiscountAmount += it.Amount - it.FinalAmount
		o.appendToExpressCalculator(ue, it, ec)
	}
	ec.Calculate("") //todo:??暂不支持区域
	o.value.ExpressFee = ec.Total()
	o.value.PackageFee = 0
	//计算最终金额
	o.fixFinalAmount()
}

// 创建商品信息,并读取价格及运费信息
func (o *wholesaleOrderImpl) createItem(i *cart.ItemPair) *orderItem {
	// 获取商品信息
	it := o.itemRepo.GetItem(i.ItemId)
	sku := it.GetSku(i.SkuId)
	iv := it.GetValue()
	// 获取商品已销售快照
	snap := o.itemRepo.SnapshotService().GetLatestSalesSnapshot(
		i.ItemId, i.SkuId)
	if snap == nil {
		domain.HandleError(errors.New("商品快照生成失败："+
			strconv.Itoa(int(i.SkuId))), "domain")
		return nil
	}
	// 计算价格
	ws := it.Wholesale()
	wsPrice := ws.GetWholesalePrice(i.SkuId, i.Quantity)
	price := float32(wsPrice) * float32(i.Quantity)
	// 计算重量及体积
	weight := sku.Weight * i.Quantity
	bulk := sku.Bulk * i.Quantity
	return &orderItem{
		ID:             0,
		OrderId:        0,
		ItemId:         i.ItemId,
		SkuId:          i.SkuId,
		SnapshotId:     snap.Id,
		Quantity:       i.Quantity,
		ReturnQuantity: 0,
		Amount:         price,
		FinalAmount:    price,
		VendorId:       iv.VendorId,
		ShopId:         iv.ShopId,
		Weight:         weight,
		Bulk:           bulk,
		ExpressTplId:   iv.ExpressTid,
	}
}

// 加入运费计算器
func (o *wholesaleOrderImpl) appendToExpressCalculator(ue express.IUserExpress,
	item *orderItem, cul express.IExpressCalculator) {
	tpl := ue.GetTemplate(item.ExpressTplId)
	if tpl != nil {
		var err error
		v := tpl.Value()
		switch v.Basis {
		case express.BasisByNumber:
			err = cul.Add(item.ExpressTplId, item.Quantity)
		case express.BasisByWeight:
			err = cul.Add(item.ExpressTplId, item.Weight)
		case express.BasisByVolume:
			err = cul.Add(item.ExpressTplId, item.Weight)
		}
		if err != nil {
			log.Println("[ Wholesale Order][ Express][ Error]:", err)
		}
	}
}

// 转换订单商品
func (o *wholesaleOrderImpl) parseComplexItem(i *order.WholesaleItem) *order.ComplexItem {
	it := &order.ComplexItem{
		ID:             i.ID,
		OrderId:        i.OrderId,
		ItemId:         int64(i.ItemId),
		SkuId:          int64(i.SkuId),
		SnapshotId:     int64(i.SnapshotId),
		Quantity:       i.Quantity,
		ReturnQuantity: i.ReturnQuantity,
		Amount:         float64(i.Amount),
		FinalAmount:    float64(i.FinalAmount),
		IsShipped:      i.IsShipped,
		Data:           make(map[string]string),
	}
	o.baseOrderImpl.bindItemInfo(it)
	return it
}

// 复合的订单信息
func (o *wholesaleOrderImpl) Complex() *order.ComplexOrder {
	v := o.getValue()
	co := o.baseOrderImpl.Complex()
	co.SubOrderId = 0
	co.VendorId = v.VendorId
	co.ShopId = v.ShopId
	co.Subject = ""
	co.ConsigneePerson = v.ConsigneePerson
	co.ConsigneePhone = v.ConsigneePhone
	co.ShippingAddress = v.ShippingAddress
	co.DiscountAmount = float64(v.DiscountAmount)
	co.ItemAmount = float64(v.ItemAmount)
	co.ExpressFee = float64(v.ExpressFee)
	co.PackageFee = float64(v.PackageFee)
	co.FinalAmount = float64(v.FinalAmount)
	co.BuyerComment = v.BuyerComment
	co.IsBreak = 0
	co.UpdateTime = v.UpdateTime
	co.Items = []*order.ComplexItem{}
	for _, v := range o.Items() {
		co.Items = append(co.Items, o.parseComplexItem(v))
	}
	return co
}

// 提交订单。如遇拆单,需均摊优惠抵扣金额到商品
func (o *wholesaleOrderImpl) Submit() error {
	if o.GetAggregateRootId() > 0 {
		return errors.New("订单不允许重复提交")
	}
	err := o.checkBuyer()
	if err == nil {
		err = o.takeItemStock(o.items)
	}
	if err != nil {
		return err
	}
	// 提交订单
	err = o.baseOrderImpl.Submit()
	if err == nil {
		// 保存订单信息到常规订单
		o.value.OrderId = o.GetAggregateRootId()
		o.value.OrderNo = o.OrderNo()
		o.value.State = int32(order.StatAwaitingPayment)
		o.value.CreateTime = o.baseValue.CreateTime
		o.value.UpdateTime = o.baseValue.CreateTime
		// 保存订单
		o.value.ID, err = util.I64Err(o.repo.SaveWholesaleOrder(o.value))
		if err == nil {
			// 存储Items
			err = o.saveOrderItemsOnSubmit()
			// 生成支付单
			err = o.createPaymentForOrder()
		}
	}

	return err
}

// 检查买家及收货地址
func (o *wholesaleOrderImpl) checkBuyer() error {
	buyer := o.Buyer()
	if buyer == nil {
		return member.ErrNoSuchMember
	}
	if buyer.GetValue().State == 0 {
		return member.ErrMemberLocked
	}
	if o.value.ShippingAddress == "" ||
		o.value.ConsigneePhone == "" ||
		o.value.ConsigneePerson == "" {
		return order.ErrMissingShipAddress
	}
	return nil
}

// 扣除库存
func (o *wholesaleOrderImpl) takeItemStock(items []*orderItem) (err error) {
	okIndex := 0
	// 占用库存，并记录库存占用成功索引
	for _, v := range items {
		it := o.itemRepo.GetItem(v.ItemId)
		if it == nil {
			err = item.ErrNoSuchItem
		} else {
			err = it.TakeStock(v.SkuId, v.Quantity)
		}
		if err != nil {
			break
		}
		okIndex++
	}
	// 如果库存占用失败，则释放库存
	if err != nil {
		for i := 0; i < okIndex; i++ {
			v := items[i]
			it := o.itemRepo.GetItem(v.ItemId)
			it.FreeStock(v.SkuId, v.Quantity)
		}
	}
	return err
}

// 计算折扣
func (o *wholesaleOrderImpl) applyGroupDiscount() {
	var groupId int32 = 1
	mch := o.mchRepo.GetMerchant(o.value.VendorId)
	if mch != nil {
		basisAmount := int32(o.value.ItemAmount)
		ws := mch.Wholesaler()
		rate := ws.GetRebateRate(groupId, basisAmount)
		disAmount := rate * float64(basisAmount)
		if disAmount > 0 {
			o.value.DiscountAmount += float32(disAmount)
			o.fixFinalAmount()
		}
	}
}

// 平均优惠抵扣金额到商品
func (o *wholesaleOrderImpl) avgDiscountForItem() {
	if o.items == nil {
		panic(errors.New("仅能在下单时进行商品抵扣平均"))
	}
	if o.value.DiscountAmount > 0 {
		totalFee := o.value.ItemAmount
		disFee := o.value.DiscountAmount
		for _, v := range o.items {
			b := v.Amount / totalFee
			v.FinalAmount = v.Amount - b*disFee
		}
	}
}

// 修正订单实际金额
func (o *wholesaleOrderImpl) fixFinalAmount() {
	o.value.FinalAmount = o.value.ItemAmount - o.value.DiscountAmount +
		o.value.ExpressFee + o.value.PackageFee
}

// 保存商品项
func (o *wholesaleOrderImpl) saveOrderItemsOnSubmit() (err error) {
	orderId := o.GetAggregateRootId()
	for _, v := range o.items {
		v.OrderId = orderId
		it := o.parseOrderItem(v)
		_, err = o.repo.SaveWholesaleItem(it)
		if err != nil {
			break
		}
	}
	return err
}

// 保存商品项
func (o *wholesaleOrderImpl) saveOrderItems() (err error) {
	orderId := o.GetAggregateRootId()
	if o.realItems != nil {
		for _, v := range o.realItems {
			v.OrderId = orderId
			_, err = o.repo.SaveWholesaleItem(v)
			if err != nil {
				break
			}
		}
	}
	return err
}

// 转换订单商品
func (o *wholesaleOrderImpl) parseOrderItem(i *orderItem) *order.WholesaleItem {
	return &order.WholesaleItem{
		ID:             0,
		OrderId:        i.OrderId,
		ItemId:         int64(i.ItemId),
		SkuId:          int64(i.SkuId),
		SnapshotId:     int64(i.SnapshotId),
		Quantity:       i.Quantity,
		ReturnQuantity: i.ReturnQuantity,
		Amount:         i.Amount,
		FinalAmount:    i.FinalAmount,
		IsShipped:      0,
		UpdateTime:     i.UpdateTime,
	}
}

// 设置配送地址
func (o *wholesaleOrderImpl) SetAddress(addressId int64) error {
	if addressId <= 0 {
		return order.ErrNoSuchAddress
	}
	buyer := o.Buyer()
	if buyer == nil {
		return member.ErrNoSuchMember
	}
	addr := buyer.Profile().GetAddress(addressId)
	if addr == nil {
		return order.ErrNoSuchAddress
	}
	d := addr.GetValue()
	o.value.ShippingAddress = strings.Replace(d.Area, " ", "", -1) + d.Address
	o.value.ConsigneePerson = d.RealName
	o.value.ConsigneePhone = d.Phone
	return nil
}

// 设置或添加买家留言，如已经提交订单，将在原留言后附加
func (o *wholesaleOrderImpl) SetComment(comment string) {
	if o.GetAggregateRootId() > 0 {
		o.value.BuyerComment += "$break$" + comment
	} else {
		o.value.BuyerComment = comment
	}
}

// 生成支付单
func (o *wholesaleOrderImpl) createPaymentForOrder() error {
	v := o.baseOrderImpl.createPaymentOrder()
	v.SellerId = int(o.value.VendorId)
	v.ItemAmount = int(o.value.FinalAmount * 100)
	o.paymentOrder = o.payRepo.CreatePaymentOrder(v)
	return o.paymentOrder.Submit()
}

// 获取商品项
func (o *wholesaleOrderImpl) Items() []*order.WholesaleItem {
	if o.realItems == nil {
		id := o.GetAggregateRootId()
		o.realItems = o.repo.SelectWholesaleItem("order_id= $1", id)
	}
	return o.realItems
}

// 在线支付交易完成
func (o *wholesaleOrderImpl) OnlinePaymentTradeFinish() error {
	if o.value.IsPaid == 1 {
		return order.ErrOrderPayed
	}
	if o.value.State == order.StatAwaitingPayment {
		o.value.IsPaid = 1
		o.value.State = order.StatAwaitingConfirm
		err := o.AppendLog(order.LogSetup, true, "{finish_pay}")
		if err == nil {
			err = o.saveWholesaleOrder()
		}
		return err
	}
	return order.ErrUnusualOrderStat
}

// 记录订单日志
func (o *wholesaleOrderImpl) AppendLog(logType order.LogType,
	system bool, message string) error {
	return nil
	//todo: ???
	if o.GetAggregateRootId() <= 0 {
		panic("order not created.")
	}
	var systemInt int
	if system {
		systemInt = 1
	} else {
		systemInt = 0
	}
	l := &order.OrderLog{
		OrderId:    o.GetAggregateRootId(),
		Type:       int(logType),
		IsSystem:   systemInt,
		OrderState: int(o.value.State),
		Message:    message,
		RecordTime: time.Now().Unix(),
	}
	return o.repo.SaveNormalSubOrderLog(l)
}

// 添加备注
func (o *wholesaleOrderImpl) AddRemark(remark string) {
	o.value.BuyerComment = remark
}

// 保存订单
func (o *wholesaleOrderImpl) saveWholesaleOrder() error {
	unix := time.Now().Unix()
	o.value.UpdateTime = unix
	if o.getValue().ID <= 0 {
		panic("please use Submit() to create new wholesale order!")
	}
	_, err := o.repo.SaveWholesaleOrder(o.value)
	if err == nil {
		o.syncOrderState()
	}
	return err
}

// 同步订单状态
func (o *wholesaleOrderImpl) syncOrderState() {
	if o.State() != order.StatBreak {
		o.saveOrderState(order.OrderState(o.value.State))
	}
}

// 确认订单
func (o *wholesaleOrderImpl) Confirm() error {
	if o.value.State < order.StatAwaitingConfirm {
		return order.ErrOrderNotPayed
	}
	if o.value.State >= order.StatAwaitingPickup {
		return order.ErrOrderHasConfirm
	}
	o.value.State = order.StatAwaitingPickup
	o.value.UpdateTime = time.Now().Unix()
	err := o.saveWholesaleOrder()
	if err == nil {
		go o.addItemSalesNum()
		err = o.AppendLog(order.LogSetup, false, "{confirm}")
	}
	return err
}

// 增加商品的销售数量
func (o *wholesaleOrderImpl) addItemSalesNum() {
	for _, v := range o.Items() {
		it := o.itemRepo.GetItem(v.ItemId)
		err := it.AddSalesNum(v.SkuId, v.Quantity)
		if err != nil {
			log.Println("---增加销售数量：", v.ItemId,
				" sku:", v.SkuId, " error:", err.Error())
		}
	}
}

// 捡货(备货)
func (o *wholesaleOrderImpl) PickUp() error {
	if o.value.State < order.StatAwaitingPickup {
		return order.ErrOrderNotConfirm
	}
	if o.value.State >= order.StatAwaitingShipment {
		return order.ErrOrderHasPickUp
	}
	o.value.State = order.StatAwaitingShipment
	o.value.UpdateTime = time.Now().Unix()
	err := o.saveWholesaleOrder()
	if err == nil {
		err = o.AppendLog(order.LogSetup, true, "{pickup}")
	}
	return err
}

// 创建发货单
func (o *wholesaleOrderImpl) createShipmentOrder(items []*order.WholesaleItem) shipment.IShipmentOrder {
	if items == nil || len(items) == 0 {
		return nil
	}
	unix := time.Now().Unix()
	so := &shipment.ShipmentOrder{
		ID:          0,
		OrderId:     o.GetAggregateRootId(),
		SubOrderId:  0,
		ShipmentLog: "",
		ShipTime:    unix,
		State:       shipment.StatAwaitingShipment,
		UpdateTime:  unix,
		Items:       []*shipment.Item{},
	}
	for _, v := range items {
		if v.IsShipped == 1 {
			continue
		}
		so.Amount += float64(v.Amount)
		so.FinalAmount += float64(v.FinalAmount)
		so.Items = append(so.Items, &shipment.Item{
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

// 发货
func (o *wholesaleOrderImpl) Ship(spId int32, spOrder string) error {
	if o.value.State < order.StatAwaitingShipment {
		return order.ErrOrderNotPickUp
	}
	if o.value.State >= order.StatShipped {
		return order.ErrOrderShipped
	}
	id := o.GetAggregateRootId()
	if list := o.shipRepo.GetShipOrders(id, false); len(list) > 0 {
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
		o.value.State = order.StatShipped
		o.value.UpdateTime = time.Now().Unix()
		err = o.saveWholesaleOrder()
		if err == nil {
			// 保存商品的发货状态
			err = o.saveOrderItems()
			o.AppendLog(order.LogSetup, true, "{shipped}")
		}
	}
	return err
}

// 已收货
func (o *wholesaleOrderImpl) BuyerReceived() error {
	if o.value.State < order.StatShipped {
		return order.ErrOrderNotShipped
	}
	if o.value.State >= order.StatCompleted {
		return order.ErrIsCompleted
	}
	dt := time.Now()
	o.value.State = order.StatCompleted
	o.value.UpdateTime = dt.Unix()
	err := o.saveWholesaleOrder()
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

func (o *wholesaleOrderImpl) getOrderAmount() (amount float32, refund float32) {
	items := o.Items()
	for _, it := range items {
		if it.ReturnQuantity > 0 {
			a := it.Amount / float32(it.Quantity) * float32(it.ReturnQuantity)
			if it.ReturnQuantity != it.Quantity {
				amount += it.Amount - a
			}
			refund += a
		} else {
			amount += it.Amount
		}
	}
	//如果非全部退货、退款,则加上运费及包装费
	if amount > 0 {
		amount += o.value.ExpressFee + o.value.PackageFee
	}
	return amount, refund
}

// 获取订单的成本
func (o *wholesaleOrderImpl) getOrderCost() float32 {
	var cost float32
	items := o.Items()
	for _, it := range items {
		snap := o.itemRepo.GetSalesSnapshot(it.SnapshotId)
		cost += snap.Cost * float32(it.Quantity-it.ReturnQuantity)
	}
	//如果非全部退货、退款,则加上运费及包装费
	if cost > 0 {
		cost += o.value.ExpressFee + o.value.PackageFee
	}
	return cost
}

// 商户结算
func (o *wholesaleOrderImpl) vendorSettle() error {
	vendor := o.mchRepo.GetMerchant(o.value.VendorId)
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
func (o *wholesaleOrderImpl) vendorSettleByCost(vendor merchant.IMerchant) error {
	_, refund := o.getOrderAmount()
	sAmount := o.getOrderCost()
	if sAmount > 0 {
		totalAmount := int(sAmount * float32(enum.RATE_AMOUNT))
		refundAmount := int(refund * float32(enum.RATE_AMOUNT))
		tradeFee, _ := vendor.SaleManager().MathTradeFee(
			merchant.TKWholesaleOrder, totalAmount)
		return vendor.Account().SettleOrder(o.OrderNo(),
			totalAmount, tradeFee, refundAmount, "批发订单结算")
	}
	return nil
}

// 根据比例进行商户结算
func (o *wholesaleOrderImpl) vendorSettleByRate(vendor merchant.IMerchant) error {
	rate := o.registryRepo.Get(registry.MchOrderSettleRate).FloatValue()
	amount, refund := o.getOrderAmount()
	sAmount := amount * float32(rate)
	if sAmount > 0 {
		totalAmount := int(sAmount * float32(enum.RATE_AMOUNT))
		refundAmount := int(refund * float32(enum.RATE_AMOUNT))
		tradeFee, _ := vendor.SaleManager().MathTradeFee(
			merchant.TKWholesaleOrder, totalAmount)
		return vendor.Account().SettleOrder(o.OrderNo(),
			totalAmount, tradeFee, refundAmount, "批发订单结算")
	}
	return nil
}
func (o *wholesaleOrderImpl) vendorSettleByOrderQuantity(vendor merchant.IMerchant) error {
	fee := o.registryRepo.Get(registry.MchSingleOrderServiceFee).FloatValue()
	amount, refund := o.getOrderAmount()
	if fee > 0 {
		totalAmount := int(math.Min(float64(amount), fee) * float64(enum.RATE_AMOUNT))
		refundAmount := int(refund * float32(enum.RATE_AMOUNT))
		tradeFee, _ := vendor.SaleManager().MathTradeFee(
			merchant.TKWholesaleOrder, totalAmount)
		return vendor.Account().SettleOrder(o.value.OrderNo,
			totalAmount, tradeFee, refundAmount, "零售订单结算")

	}
	return nil
}

// 完成订单
func (o *wholesaleOrderImpl) onOrderComplete() error {
	id := o.GetAggregateRootId()
	// 更新发货单
	soList := o.shipRepo.GetShipOrders(id, false)
	for _, v := range soList {
		domain.HandleError(v.Completed(), "domain")
	}
	// 更新会员账户
	err := o.updateAccountForOrder()
	if err == nil {
		// 处理返现
		//err = o.handleCashBack()
	}

	return err
}

// 更新账户
func (o *wholesaleOrderImpl) updateAccountForOrder() error {
	if o.value.State != order.StatCompleted {
		return order.ErrUnusualOrderStat
	}
	m := o.Buyer()
	var err error
	ov := o.value
	amount := ov.FinalAmount
	acc := m.GetAccount()

	// 增加经验
	expEnabled := o.registryRepo.Get(registry.ExperienceEnabled).BoolValue()
	if expEnabled {
		rate := o.registryRepo.Get(registry.ExperienceRateByWholesaleOrder).FloatValue()
		if exp := int(float64(amount) * rate); exp > 0 {
			if err = m.AddExp(exp); err != nil {
				return err
			}
		}
	}

	// 增加积分
	//todo: 增加阶梯的返积分,比如订单满30送100积分, 不考虑额外赠送,额外的当做补贴
	rate := o.registryRepo.Get(registry.IntegralRateByWholesaleOrder).FloatValue()
	integral := int64(float64(amount) * rate)
	// 赠送积分
	if integral > 0 {
		err = m.GetAccount().Charge(member.AccountIntegral,
			"购物消费赠送积分", float32(integral), o.OrderNo(), "sys")
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

// 获取订单的日志
func (o *wholesaleOrderImpl) LogBytes() []byte {
	buf := bytes.NewBufferString("")
	orderId := o.GetAggregateRootId()
	list := o.repo.GetSubOrderLogs(orderId)
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

func (o *wholesaleOrderImpl) getLogStringByStat(stat int) string {
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

// 取消订单/退款
func (o *wholesaleOrderImpl) Cancel(reason string) error {
	if o.value.State == order.StatCancelled {
		return order.ErrOrderCancelled
	}
	// 已发货订单无法取消
	if o.value.State >= order.StatShipped {
		return order.ErrOrderShippedCancel
	}
	o.value.State = order.StatCancelled
	o.value.UpdateTime = time.Now().Unix()
	err := o.saveWholesaleOrder()
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
func (o *wholesaleOrderImpl) cancelGoods() error {
	for _, v := range o.Items() {
		snapshot := o.itemRepo.GetSalesSnapshot(v.SnapshotId)
		if snapshot == nil {
			return item.ErrNoSuchSnapshot
		}
		gds := o.itemRepo.GetItem(snapshot.SkuId)
		if gds != nil {
			// 释放库存
			gds.FreeStock(v.SkuId, v.Quantity)
			// 如果订单已付款，则取消销售数量
			if o.value.IsPaid == 1 {
				gds.CancelSale(v.SkuId, v.Quantity, o.value.OrderNo)
			}
		}
	}
	return nil
}

// 获取支付单
func (o *wholesaleOrderImpl) GetPaymentOrder() payment.IPaymentOrder {
	if o.paymentOrder == nil {
		id := o.GetAggregateRootId()
		if id <= 0 {
			panic(" Get payment order error ; because of order no yet created!")
		}
		o.paymentOrder = o.payRepo.GetPaymentBySalesOrderId(id)
	}
	return o.paymentOrder
}

// 取消支付单
func (o *wholesaleOrderImpl) cancelPaymentOrder() error {
	po := o.GetPaymentOrder()
	if po != nil {
		return po.Cancel()
	}
	return nil
}

// 谢绝订单
func (o *wholesaleOrderImpl) Decline(reason string) error {
	if o.value.State == order.StatAwaitingPayment {
		return o.Cancel("商户取消,原因:" + reason)
	}
	if o.value.State >= order.StatShipped ||
		o.value.State >= order.StatCancelled {
		return order.ErrOrderCancelled
	}
	o.value.State = order.StatDeclined
	o.value.UpdateTime = time.Now().Unix()
	return o.saveWholesaleOrder()
}
