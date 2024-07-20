package order

import (
	"bytes"
	"errors"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/shipment"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/gof/util"
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
	itemRepo     item.IItemRepo
	mchRepo      merchant.IMerchantRepo
	valueRepo    valueobject.IValueRepo
	shopRepo     shop.IShopRepo
	registryRepo registry.IRegistryRepo
}

func newWholesaleOrder(base *baseOrderImpl,
	shoppingRepo order.IOrderRepo, goodsRepo item.IItemRepo,
	expressRepo express.IExpressRepo, payRepo payment.IPaymentRepo,
	shipRepo shipment.IShipmentRepo, mchRepo merchant.IMerchantRepo,
	shopRepo shop.IShopRepo, valueRepo valueobject.IValueRepo,
	registryRepo registry.IRegistryRepo) order.IOrder {
	o := &wholesaleOrderImpl{
		baseOrderImpl: base,
		orderRepo:     shoppingRepo,
		itemRepo:      goodsRepo,
		expressRepo:   expressRepo,
		payRepo:       payRepo,
		shipRepo:      shipRepo,
		mchRepo:       mchRepo,
		shopRepo:      shopRepo,
		valueRepo:     valueRepo,
		registryRepo:  registryRepo,
	}
	return o.init()
}

func (o *wholesaleOrderImpl) init() order.IOrder {
	if o.GetAggregateRootId() <= 0 {
		o.value = &order.WholesaleOrder{
			Id:       0,
			OrderNo:  "",
			OrderId:  0,
			BuyerId:  o.baseValue.BuyerId,
			VendorId: 0,
			ShopId:   0,
			Status:   o.baseValue.Status,
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
	// 获取运营商和店铺编号
	o.value.VendorId = o.items[0].VendorId
	o.value.ShopId = o.items[0].ShopId

	// 店铺名称
	isp := o.shopRepo.GetShop(o.items[0].ShopId).(shop.IOnlineShop)
	o.value.ShopName = isp.GetShopValue().ShopName

	// 运费计算器
	ue := o.expressRepo.GetUserExpress(int(o.value.VendorId))
	ec := ue.CreateCalculator()
	// 计算订单金额及运费
	for _, it := range o.items {
		o.baseValue.ItemCount += int(it.Quantity)
		o.baseValue.ItemAmount += it.Amount
		o.baseValue.DiscountAmount += it.Amount - it.FinalAmount
		o.appendToExpressCalculator(ue, it, ec)
	}
	ec.Calculate("") //todo:??暂不支持区域
	o.baseValue.ExpressFee = ec.Total()
	o.baseValue.PackageFee = 0
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
	price := int64(float32(wsPrice) * float32(i.Quantity))
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
	tid := int(item.ExpressTplId)
	tpl := ue.GetTemplate(tid)
	if tpl != nil {
		var err error
		v := tpl.Value()
		switch v.Basis {
		case express.BasisByNumber:
			err = cul.Add(tid, int(item.Quantity))
		case express.BasisByWeight:
			err = cul.Add(tid, int(item.Weight))
		case express.BasisByVolume:
			err = cul.Add(tid, int(item.Weight))
		}
		if err != nil {
			log.Println("[ Wholesale Order][ Express][ Error]:", err)
		}
	}
}

// 转换订单商品
func (o *wholesaleOrderImpl) parseComplexItem(i *order.WholesaleItem) *order.ComplexItem {
	snap := o.itemRepo.GetSalesSnapshot(i.SnapshotId)
	it := &order.ComplexItem{
		ID:             i.ID,
		ItemId:         i.ItemId,
		SkuId:          i.SkuId,
		SkuWord:        snap.Sku,
		SnapshotId:     i.SnapshotId,
		ItemTitle:      snap.GoodsTitle,
		MainImage:      snap.Image,
		Price:          i.Price,
		FinalPrice:     0,
		Quantity:       i.Quantity,
		ReturnQuantity: i.ReturnQuantity,
		Amount:         i.Amount,
		FinalAmount:    i.FinalAmount,
		IsShipped:      i.IsShipped,
		Data:           make(map[string]string),
	}
	o.baseOrderImpl.bindItemInfo(it)
	return it
}

// Complex 复合的订单信息
func (o *wholesaleOrderImpl) Complex() *order.ComplexOrder {
	co := o.baseOrderImpl.Complex()
	dt := &order.ComplexOrderDetails{
		Id:             o.GetAggregateRootId(),
		OrderNo:        o.value.OrderNo,
		ShopId:         o.value.ShopId,
		ShopName:       o.value.ShopName,
		ItemAmount:     co.ItemAmount,
		DiscountAmount: co.DiscountAmount,
		ExpressFee:     co.ExpressFee,
		PackageFee:     co.PackageFee,
		FinalAmount:    co.FinalAmount,
		BuyerComment:   o.value.BuyerComment,
		Status:         o.value.Status,
		StatusText:     "",
		Items:          []*order.ComplexItem{},
		UpdateTime:     o.value.UpdateTime,
	}
	for _, v := range o.Items() {
		dt.Items = append(dt.Items, o.parseComplexItem(v))
	}
	co.Details = append(co.Details, dt)
	return co
}

// Submit 提交订单。如遇拆单,需均摊优惠抵扣金额到商品
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
		o.value.Status = order.StatAwaitingPayment
		o.value.CreateTime = o.baseValue.CreateTime
		o.value.UpdateTime = o.baseValue.CreateTime
		// 保存订单
		o.value.Id, err = util.I64Err(o.repo.SaveWholesaleOrder(o.value))
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
	if buyer.TestFlag(member.FlagLocked) {
		return member.ErrMemberLocked
	}
	if o.baseValue.ShippingAddress == "" ||
		o.baseValue.ConsigneePhone == "" ||
		o.baseValue.ConsigneeName == "" {
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
			it.ReleaseStock(v.SkuId, v.Quantity)
		}
	}
	return err
}

// 计算折扣
func (o *wholesaleOrderImpl) applyGroupDiscount() {
	var groupId int32 = 1
	mch := o.mchRepo.GetMerchant(int(o.value.VendorId))
	if mch != nil {
		basisAmount := int32(o.baseValue.ItemAmount)
		ws := mch.Wholesaler()
		rate := ws.GetRebateRate(groupId, basisAmount)
		disAmount := rate * float64(basisAmount)
		if disAmount > 0 {
			o.baseValue.DiscountAmount += int64(disAmount)
			o.fixFinalAmount()
		}
	}
}

// 平均优惠抵扣金额到商品
func (o *wholesaleOrderImpl) avgDiscountForItem() {
	if o.items == nil {
		panic(errors.New("仅能在下单时进行商品抵扣平均"))
	}
	if o.baseValue.DiscountAmount > 0 {
		totalFee := o.baseValue.ItemAmount
		disFee := o.baseValue.DiscountAmount
		for _, v := range o.items {
			b := v.Amount / totalFee
			v.FinalAmount = v.Amount - b*disFee
		}
	}
}

// 修正订单实际金额
func (o *wholesaleOrderImpl) fixFinalAmount() {
	o.baseValue.FinalAmount = o.baseValue.ItemAmount - o.baseValue.DiscountAmount +
		o.baseValue.ExpressFee + o.baseValue.PackageFee
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
		ItemId:         i.ItemId,
		SkuId:          i.SkuId,
		SnapshotId:     i.SnapshotId,
		Quantity:       i.Quantity,
		ReturnQuantity: i.ReturnQuantity,
		Amount:         i.Amount,
		FinalAmount:    i.FinalAmount,
		IsShipped:      0,
		UpdateTime:     i.UpdateTime,
	}
}

// SetComment 设置或添加买家留言，如已经提交订单，将在原留言后附加
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
	v.TotalAmount = o.baseValue.FinalAmount
	o.paymentOrder = o.payRepo.CreatePaymentOrder(v)
	return o.paymentOrder.Submit()
}

// Items 获取商品项
func (o *wholesaleOrderImpl) Items() []*order.WholesaleItem {
	if o.realItems == nil {
		id := o.GetAggregateRootId()
		o.realItems = o.repo.SelectWholesaleItem("order_id= $1", id)
	}
	return o.realItems
}

// OnlinePaymentTradeFinish 在线支付交易完成
func (o *wholesaleOrderImpl) OnlinePaymentTradeFinish() error {
	if o.value.Status > order.StatAwaitingPayment {
		return order.ErrOrderPayed
	}
	if o.value.Status == order.StatAwaitingPayment {
		o.value.Status = order.StatAwaitingPickup
		err := o.AppendLog(order.LogSetup, true, "{finish_pay}")
		if err == nil {
			err = o.saveWholesaleOrder()
			if err == nil {
				o.baseValue.IsPaid = 1
				o.baseOrderImpl.saveOrder()
			}
		}
		return err
	}
	return order.ErrUnusualOrderStat
}

// AppendLog 记录订单日志
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
		OrderState: int(o.value.Status),
		Message:    message,
		RecordTime: time.Now().Unix(),
	}
	return o.repo.SaveNormalSubOrderLog(l)
}

// AddRemark 添加备注
func (o *wholesaleOrderImpl) AddRemark(remark string) {
	o.value.BuyerComment = remark
}

// 保存订单
func (o *wholesaleOrderImpl) saveWholesaleOrder() error {
	unix := time.Now().Unix()
	o.value.UpdateTime = unix
	if o.getValue().Id <= 0 {
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
		o.saveOrderState(order.OrderStatus(o.value.Status))
	}
}

// Confirm 确认订单
func (o *wholesaleOrderImpl) Confirm() error {
	if o.value.Status < order.StatAwaitingConfirm {
		return order.ErrOrderNotPayed
	}
	if o.value.Status >= order.StatAwaitingPickup {
		return order.ErrOrderHasConfirm
	}
	o.value.Status = order.StatAwaitingPickup
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
	if o.value.Status < order.StatAwaitingPickup {
		return order.ErrOrderNotConfirm
	}
	if o.value.Status >= order.StatAwaitingShipment {
		return order.ErrOrderHasPickUp
	}
	o.value.Status = order.StatAwaitingShipment
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

// Ship 发货
func (o *wholesaleOrderImpl) Ship(spId int32, spOrder string) error {
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
		o.value.Status = order.StatShipped
		o.value.UpdateTime = time.Now().Unix()
		err = o.saveWholesaleOrder()
		if err == nil {
			// 保存商品的发货状态
			err = o.saveOrderItems()
			_ = o.AppendLog(order.LogSetup, true, "{shipped}")
		}
	}
	return err
}

// BuyerReceived 已收货
func (o *wholesaleOrderImpl) BuyerReceived() error {
	if o.value.Status < order.StatShipped {
		return order.ErrOrderNotShipped
	}
	if o.value.Status >= order.StatCompleted {
		return order.ErrIsCompleted
	}
	dt := time.Now()
	o.value.Status = order.StatCompleted
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

func (o *wholesaleOrderImpl) getOrderAmount() (amount int, refund int) {
	items := o.Items()
	for _, it := range items {
		if it.ReturnQuantity > 0 {
			a := int64(float32(it.Amount) / float32(it.Quantity) * float32(it.ReturnQuantity))
			if it.ReturnQuantity != it.Quantity {
				amount += int(it.Amount - a)
			}
			refund += int(a)
		} else {
			amount += int(it.Amount)
		}
	}
	//如果非全部退货、退款,则加上运费及包装费
	if amount > 0 {
		amount += int(o.baseValue.ExpressFee + o.baseValue.PackageFee)
	}
	return amount, refund
}

// 获取订单的成本
func (o *wholesaleOrderImpl) getOrderCost() int64 {
	var cost int64
	items := o.Items()
	for _, it := range items {
		snap := o.itemRepo.GetSalesSnapshot(it.SnapshotId)
		cost += snap.Cost * int64(it.Quantity-it.ReturnQuantity)
	}
	//如果非全部退货、退款,则加上运费及包装费
	if cost > 0 {
		cost += o.baseValue.ExpressFee + o.baseValue.PackageFee
	}
	return cost
}

// 商户结算
func (o *wholesaleOrderImpl) vendorSettle() error {
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
func (o *wholesaleOrderImpl) vendorSettleByCost(vendor merchant.IMerchantAggregateRoot) error {
	_, refund := o.getOrderAmount()
	sAmount := o.getOrderCost()
	if sAmount > 0 {
		totalAmount := sAmount
		refundAmount := refund
		tradeFee, _ := vendor.SaleManager().MathTradeFee(
			merchant.TKWholesaleOrder, int(totalAmount))
		return vendor.Account().SettleOrder(o.OrderNo(),
			int(totalAmount), tradeFee, refundAmount, "批发订单结算")
	}
	return nil
}

// 根据比例进行商户结算
func (o *wholesaleOrderImpl) vendorSettleByRate(vendor merchant.IMerchantAggregateRoot) error {
	rate := o.registryRepo.Get(registry.MchOrderSettleRate).FloatValue()
	amount, refund := o.getOrderAmount()
	sAmount := float32(amount) * float32(rate)
	if sAmount > 0 {
		totalAmount := int(sAmount * float32(enum.RATE_AMOUNT))
		refundAmount := int(float32(refund) * float32(enum.RATE_AMOUNT))
		tradeFee, _ := vendor.SaleManager().MathTradeFee(
			merchant.TKWholesaleOrder, totalAmount)
		return vendor.Account().SettleOrder(o.OrderNo(),
			totalAmount, tradeFee, refundAmount, "批发订单结算")
	}
	return nil
}
func (o *wholesaleOrderImpl) vendorSettleByOrderQuantity(vendor merchant.IMerchantAggregateRoot) error {
	fee := o.registryRepo.Get(registry.MchSingleOrderServiceFee).FloatValue()
	amount, refund := o.getOrderAmount()
	if fee > 0 {
		totalAmount := int(math.Min(float64(amount), fee) * float64(enum.RATE_AMOUNT))
		refundAmount := int(float32(refund) * float32(enum.RATE_AMOUNT))
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
	if o.value.Status != order.StatCompleted {
		return order.ErrUnusualOrderStat
	}
	m := o.Buyer()
	var err error
	ov := o.baseValue
	amount := ov.FinalAmount
	acc := m.GetAccount()

	// 增加经验
	expEnabled := o.registryRepo.Get(registry.ExperienceEnabled).BoolValue()
	if expEnabled {
		rate := o.registryRepo.Get(registry.ExperienceRateByWholesaleOrder).FloatValue()
		if exp := int(float64(amount) * rate / 100); exp > 0 {
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
		_, err = acc.CarryTo(member.AccountIntegral,
			member.AccountOperateData{
				Title:   "购物消费赠送积分",
				Amount:  integral,
				OuterNo: o.OrderNo(),
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

// LogBytes 获取订单的日志
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

// Cancel 取消订单/退款
func (o *wholesaleOrderImpl) Cancel(reason string) error {
	if o.value.Status == order.StatCancelled {
		return order.ErrOrderCancelled
	}
	// 已发货订单无法取消
	if o.value.Status >= order.StatShipped {
		return order.ErrOrderShippedCancel
	}
	o.value.Status = order.StatCancelled
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
			gds.ReleaseStock(v.SkuId, v.Quantity)
			// 如果订单已付款，则取消销售数量
			if o.value.Status > order.StatAwaitingPayment {
				gds.CancelSale(v.SkuId, v.Quantity, o.value.OrderNo)
			}
		}
	}
	return nil
}

// GetPaymentOrder 获取支付单
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

// Decline 谢绝订单
func (o *wholesaleOrderImpl) Decline(reason string) error {
	if o.value.Status == order.StatAwaitingPayment {
		return o.Cancel("商户取消,原因:" + reason)
	}
	if o.value.Status >= order.StatShipped ||
		o.value.Status >= order.StatCancelled {
		return order.ErrOrderCancelled
	}
	o.value.Status = order.StatDeclined
	o.value.UpdateTime = time.Now().Unix()
	return o.saveWholesaleOrder()
}
