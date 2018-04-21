package order

import (
	"bytes"
	"errors"
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"log"
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
}

func newWholesaleOrder(base *baseOrderImpl,
	shoppingRepo order.IOrderRepo, goodsRepo item.IGoodsItemRepo,
	expressRepo express.IExpressRepo, payRepo payment.IPaymentRepo,
	shipRepo shipment.IShipmentRepo, mchRepo merchant.IMerchantRepo,
	valueRepo valueobject.IValueRepo) order.IOrder {
	o := &wholesaleOrderImpl{
		baseOrderImpl: base,
		orderRepo:     shoppingRepo,
		itemRepo:      goodsRepo,
		expressRepo:   expressRepo,
		payRepo:       payRepo,
		shipRepo:      shipRepo,
		mchRepo:       mchRepo,
		valueRepo:     valueRepo,
	}
	return o.init()
}

func (s *wholesaleOrderImpl) init() order.IOrder {
	if s.GetAggregateRootId() <= 0 {
		s.value = &order.WholesaleOrder{
			ID:          0,
			OrderNo:     "",
			OrderId:     0,
			BuyerId:     s.baseValue.BuyerId,
			VendorId:    0,
			ShopId:      0,
			ItemAmount:  0,
			ExpressFee:  0,
			PackageFee:  0,
			FinalAmount: 0,
			State:       s.baseValue.State,
		}
	}
	s.getValue()
	return s
}

func (s *wholesaleOrderImpl) getValue() *order.WholesaleOrder {
	if s.value == nil {
		id := s.GetAggregateRootId()
		if id > 0 {
			s.value = s.repo.GetWholesaleOrder("order_id=?", id)
		}
	}
	return s.value
}

// 设置商品项
func (s *wholesaleOrderImpl) SetItems(items []*cart.ItemPair) {
	if s.GetAggregateRootId() > 0 {
		panic("wholesale has created. can't use SetItems!")
	}
	s.parseOrder(items)
	// 计算折扣
	s.applyGroupDiscount()
	// 均摊优惠折扣到商品
	s.avgDiscountForItem()
}

// 转换为订单相关对象
func (s *wholesaleOrderImpl) parseOrder(items []*cart.ItemPair) {
	if s.GetAggregateRootId() > 0 {
		panic("订单已经生成，无法解析")
	}
	s.items = []*orderItem{}
	for _, v := range items {
		s.items = append(s.items, s.createItem(v))
	}
	// 获取运营商和商铺编号
	s.value.VendorId = s.items[0].VendorId
	s.value.ShopId = s.items[0].ShopId
	// 运费计算器
	ue := s.expressRepo.GetUserExpress(s.value.VendorId)
	ec := ue.CreateCalculator()
	// 计算订单金额及运费
	for _, item := range s.items {
		s.value.ItemAmount += item.Amount
		s.value.DiscountAmount += item.Amount - item.FinalAmount
		s.appendToExpressCalculator(ue, item, ec)
	}
	ec.Calculate("") //todo:??暂不支持区域
	s.value.ExpressFee = ec.Total()
	s.value.PackageFee = 0
	//计算最终金额
	s.fixFinalAmount()
}

// 创建商品信息,并读取价格及运费信息
func (s *wholesaleOrderImpl) createItem(i *cart.ItemPair) *orderItem {
	// 获取商品信息
	it := s.itemRepo.GetItem(i.ItemId)
	sku := it.GetSku(i.SkuId)
	iv := it.GetValue()
	// 获取商品已销售快照
	snap := s.itemRepo.SnapshotService().GetLatestSalesSnapshot(
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
func (s *wholesaleOrderImpl) appendToExpressCalculator(ue express.IUserExpress,
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
func (s *wholesaleOrderImpl) parseComplexItem(i *order.WholesaleItem) *order.ComplexItem {
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
	s.baseOrderImpl.bindItemInfo(it)
	return it
}

// 复合的订单信息
func (s *wholesaleOrderImpl) Complex() *order.ComplexOrder {
	v := s.getValue()
	co := s.baseOrderImpl.Complex()
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
	for _, v := range s.Items() {
		co.Items = append(co.Items, s.parseComplexItem(v))
	}
	return co
}

// 提交订单。如遇拆单,需均摊优惠抵扣金额到商品
func (s *wholesaleOrderImpl) Submit() error {
	if s.GetAggregateRootId() > 0 {
		return errors.New("订单不允许重复提交")
	}
	err := s.checkBuyer()
	if err == nil {
		err = s.takeItemStock(s.items)
	}
	if err != nil {
		return err
	}
	// 提交订单
	err = s.baseOrderImpl.Submit()
	if err == nil {
		// 保存订单信息到常规订单
		s.value.OrderId = s.GetAggregateRootId()
		s.value.OrderNo = s.OrderNo()
		s.value.State = int32(order.StatAwaitingPayment)
		s.value.CreateTime = s.baseValue.CreateTime
		s.value.UpdateTime = s.baseValue.CreateTime
		// 保存订单
		s.value.ID, err = util.I64Err(s.repo.SaveWholesaleOrder(s.value))
		if err == nil {
			// 存储Items
			err = s.saveOrderItemsOnSubmit()
			// 生成支付单
			err = s.createPaymentForOrder()
		}
	}

	return err
}

// 检查买家及收货地址
func (s *wholesaleOrderImpl) checkBuyer() error {
	buyer := s.Buyer()
	if buyer == nil {
		return member.ErrNoSuchMember
	}
	if buyer.GetValue().State == 0 {
		return member.ErrMemberDisabled
	}
	if s.value.ShippingAddress == "" ||
		s.value.ConsigneePhone == "" ||
		s.value.ConsigneePerson == "" {
		return order.ErrMissingShipAddress
	}
	return nil
}

// 扣除库存
func (s *wholesaleOrderImpl) takeItemStock(items []*orderItem) (err error) {
	okIndex := 0
	// 占用库存，并记录库存占用成功索引
	for _, v := range items {
		it := s.itemRepo.GetItem(v.ItemId)
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
			it := s.itemRepo.GetItem(v.ItemId)
			it.FreeStock(v.SkuId, v.Quantity)
		}
	}
	return err
}

// 计算折扣
func (s *wholesaleOrderImpl) applyGroupDiscount() {
	var groupId int32 = 1
	mch := s.mchRepo.GetMerchant(s.value.VendorId)
	if mch != nil {
		basisAmount := int32(s.value.ItemAmount)
		ws := mch.Wholesaler()
		rate := ws.GetRebateRate(groupId, basisAmount)
		disAmount := rate * float64(basisAmount)
		if disAmount > 0 {
			s.value.DiscountAmount += float32(disAmount)
			s.fixFinalAmount()
		}
	}
}

// 平均优惠抵扣金额到商品
func (s *wholesaleOrderImpl) avgDiscountForItem() {
	if s.items == nil {
		panic(errors.New("仅能在下单时进行商品抵扣平均"))
	}
	if s.value.DiscountAmount > 0 {
		totalFee := s.value.ItemAmount
		disFee := s.value.DiscountAmount
		for _, v := range s.items {
			b := v.Amount / totalFee
			v.FinalAmount = v.Amount - b*disFee
		}
	}
}

// 修正订单实际金额
func (s *wholesaleOrderImpl) fixFinalAmount() {
	s.value.FinalAmount = s.value.ItemAmount - s.value.DiscountAmount +
		s.value.ExpressFee + s.value.PackageFee
}

// 保存商品项
func (s *wholesaleOrderImpl) saveOrderItemsOnSubmit() (err error) {
	orderId := s.GetAggregateRootId()
	for _, v := range s.items {
		v.OrderId = orderId
		item := s.parseOrderItem(v)
		_, err = s.repo.SaveWholesaleItem(item)
		if err != nil {
			break
		}
	}
	return err
}

// 保存商品项
func (s *wholesaleOrderImpl) saveOrderItems() (err error) {
	orderId := s.GetAggregateRootId()
	if s.realItems != nil {
		for _, v := range s.realItems {
			v.OrderId = orderId
			_, err = s.repo.SaveWholesaleItem(v)
			if err != nil {
				break
			}
		}
	}
	return err
}

// 转换订单商品
func (s *wholesaleOrderImpl) parseOrderItem(i *orderItem) *order.WholesaleItem {
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
func (s *wholesaleOrderImpl) SetAddress(addressId int64) error {
	if addressId <= 0 {
		return order.ErrNoSuchAddress
	}
	buyer := s.Buyer()
	if buyer == nil {
		return member.ErrNoSuchMember
	}
	addr := buyer.Profile().GetAddress(addressId)
	if addr == nil {
		return order.ErrNoSuchAddress
	}
	d := addr.GetValue()
	s.value.ShippingAddress = strings.Replace(d.Area, " ", "", -1) + d.Address
	s.value.ConsigneePerson = d.RealName
	s.value.ConsigneePhone = d.Phone
	return nil
}

// 设置或添加买家留言，如已经提交订单，将在原留言后附加
func (s *wholesaleOrderImpl) SetComment(comment string) {
	if s.GetAggregateRootId() > 0 {
		s.value.BuyerComment += "$break$" + comment
	} else {
		s.value.BuyerComment = comment
	}
}

// 生成支付单
func (s *wholesaleOrderImpl) createPaymentForOrder() error {
	v := s.baseOrderImpl.createPaymentOrder()
	v.VendorId = s.value.VendorId
	v.TotalAmount = s.value.FinalAmount
	v.CouponDiscount = 0
	v.IntegralDiscount = 0
	v.FinalFee = v.TotalAmount - v.SubAmount - v.SystemDiscount -
		v.IntegralDiscount - v.BalanceDiscount
	s.paymentOrder = s.payRepo.CreatePaymentOrder(v)
	return s.paymentOrder.Commit()
}

// 获取商品项
func (s *wholesaleOrderImpl) Items() []*order.WholesaleItem {
	if s.realItems == nil {
		id := s.GetAggregateRootId()
		s.realItems = s.repo.SelectWholesaleItem("order_id=?", id)
	}
	return s.realItems
}

// 在线支付交易完成
func (s *wholesaleOrderImpl) OnlinePaymentTradeFinish() error {
	if s.value.IsPaid == 1 {
		return order.ErrOrderPayed
	}
	if s.value.State == order.StatAwaitingPayment {
		s.value.IsPaid = 1
		s.value.State = order.StatAwaitingConfirm
		err := s.AppendLog(order.LogSetup, true, "{finish_pay}")
		if err == nil {
			err = s.saveWholesaleOrder()
		}
		return err
	}
	return order.ErrUnusualOrderStat
}

// 记录订单日志
func (s *wholesaleOrderImpl) AppendLog(logType order.LogType,
	system bool, message string) error {
	return nil
	//todo: ???
	if s.GetAggregateRootId() <= 0 {
		return errors.New("order not created.")
	}
	var systemInt int
	if system {
		systemInt = 1
	} else {
		systemInt = 0
	}
	l := &order.OrderLog{
		OrderId:    s.GetAggregateRootId(),
		Type:       int(logType),
		IsSystem:   systemInt,
		OrderState: int(s.value.State),
		Message:    message,
		RecordTime: time.Now().Unix(),
	}
	return s.repo.SaveNormalSubOrderLog(l)
}

// 添加备注
func (s *wholesaleOrderImpl) AddRemark(remark string) {
	s.value.BuyerComment = remark
}

// 保存订单
func (s *wholesaleOrderImpl) saveWholesaleOrder() error {
	unix := time.Now().Unix()
	s.value.UpdateTime = unix
	if s.getValue().ID <= 0 {
		panic("please use Submit() to create new wholesale order!")
	}
	_, err := s.repo.SaveWholesaleOrder(s.value)
	if err == nil {
		s.syncOrderState()
	}
	return err
}

// 同步订单状态
func (s *wholesaleOrderImpl) syncOrderState() {
	if s.State() != order.StatBreak {
		s.saveOrderState(order.OrderState(s.value.State))
	}
}

// 确认订单
func (s *wholesaleOrderImpl) Confirm() error {
	if s.value.State < order.StatAwaitingConfirm {
		return order.ErrOrderNotPayed
	}
	if s.value.State >= order.StatAwaitingPickup {
		return order.ErrOrderHasConfirm
	}
	s.value.State = order.StatAwaitingPickup
	s.value.UpdateTime = time.Now().Unix()
	err := s.saveWholesaleOrder()
	if err == nil {
		go s.addItemSalesNum()
		err = s.AppendLog(order.LogSetup, false, "{confirm}")
	}
	return err
}

// 增加商品的销售数量
func (s *wholesaleOrderImpl) addItemSalesNum() {
	for _, v := range s.Items() {
		it := s.itemRepo.GetItem(v.ItemId)
		err := it.AddSalesNum(v.SkuId, v.Quantity)
		if err != nil {
			log.Println("---增加销售数量：", v.ItemId,
				" sku:", v.SkuId, " error:", err.Error())
		}
	}
}

// 捡货(备货)
func (s *wholesaleOrderImpl) PickUp() error {
	if s.value.State < order.StatAwaitingPickup {
		return order.ErrOrderNotConfirm
	}
	if s.value.State >= order.StatAwaitingShipment {
		return order.ErrOrderHasPickUp
	}
	s.value.State = order.StatAwaitingShipment
	s.value.UpdateTime = time.Now().Unix()
	err := s.saveWholesaleOrder()
	if err == nil {
		err = s.AppendLog(order.LogSetup, true, "{pickup}")
	}
	return err
}

// 创建发货单
func (s *wholesaleOrderImpl) createShipmentOrder(items []*order.WholesaleItem) shipment.IShipmentOrder {
	if items == nil || len(items) == 0 {
		return nil
	}
	unix := time.Now().Unix()
	so := &shipment.ShipmentOrder{
		ID:          0,
		OrderId:     s.GetAggregateRootId(),
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
	return s.shipRepo.CreateShipmentOrder(so)
}

// 发货
func (s *wholesaleOrderImpl) Ship(spId int32, spOrder string) error {
	if s.value.State < order.StatAwaitingShipment {
		return order.ErrOrderNotPickUp
	}
	if s.value.State >= order.StatShipped {
		return order.ErrOrderShipped
	}
	id := s.GetAggregateRootId()
	if list := s.shipRepo.GetShipOrders(id, false); len(list) > 0 {
		return order.ErrPartialShipment
	}
	if spId <= 0 || spOrder == "" {
		return shipment.ErrMissingSpInfo
	}

	so := s.createShipmentOrder(s.Items())
	if so == nil {
		return order.ErrUnusualOrder
	}
	// 生成发货单并发货
	err := so.Ship(spId, spOrder)
	if err == nil {
		s.value.State = order.StatShipped
		s.value.UpdateTime = time.Now().Unix()
		err = s.saveWholesaleOrder()
		if err == nil {
			// 保存商品的发货状态
			err = s.saveOrderItems()
			s.AppendLog(order.LogSetup, true, "{shipped}")
		}
	}
	return err
}

// 已收货
func (s *wholesaleOrderImpl) BuyerReceived() error {
	if s.value.State < order.StatShipped {
		return order.ErrOrderNotShipped
	}
	if s.value.State >= order.StatCompleted {
		return order.ErrIsCompleted
	}
	dt := time.Now()
	s.value.State = order.StatCompleted
	s.value.UpdateTime = dt.Unix()
	err := s.saveWholesaleOrder()
	if err == nil {
		err = s.AppendLog(order.LogSetup, true, "{completed}")
		if err == nil {
			go s.vendorSettle()
			// 执行其他的操作
			if err2 := s.onOrderComplete(); err != nil {
				domain.HandleError(err2, "domain")
			}
		}
	}
	return err
}

func (s *wholesaleOrderImpl) getOrderAmount() (amount float32, refund float32) {
	items := s.Items()
	for _, item := range items {
		if item.ReturnQuantity > 0 {
			a := item.Amount / float32(item.Quantity) * float32(item.ReturnQuantity)
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
		amount += s.value.ExpressFee + s.value.PackageFee
	}
	return amount, refund
}

// 获取订单的成本
func (s *wholesaleOrderImpl) getOrderCost() float32 {
	var cost float32
	items := s.Items()
	for _, item := range items {
		snap := s.itemRepo.GetSalesSnapshot(item.SnapshotId)
		cost += snap.Cost * float32(item.Quantity-item.ReturnQuantity)
	}
	//如果非全部退货、退款,则加上运费及包装费
	if cost > 0 {
		cost += s.value.ExpressFee + s.value.PackageFee
	}
	return cost
}

// 商户结算
func (s *wholesaleOrderImpl) vendorSettle() error {
	vendor := s.mchRepo.GetMerchant(s.value.VendorId)
	if vendor != nil {
		conf := s.valueRepo.GetGlobMchSaleConf()
		switch conf.MchOrderSettleMode {
		case enum.MchModeSettleByCost:
			return s.vendorSettleByCost(vendor)
		case enum.MchModeSettleByRate:
			return s.vendorSettleByRate(vendor, conf.MchOrderSettleRate)
		}

	}
	return nil
}

// 根据供货价进行商户结算
func (s *wholesaleOrderImpl) vendorSettleByCost(vendor merchant.IMerchant) error {
	_, refund := s.getOrderAmount()
	sAmount := s.getOrderCost()
	if sAmount > 0 {
		totalAmount := int(sAmount * float32(enum.RATE_Amount))
		refundAmount := int(refund * float32(enum.RATE_Amount))
		tradeFee, _ := vendor.SaleManager().MathTradeFee(
			merchant.TKWholesaleOrder, totalAmount)
		return vendor.Account().SettleOrder(s.OrderNo(),
			totalAmount, tradeFee, refundAmount, "批发订单结算")
	}
	return nil
}

// 根据比例进行商户结算
func (s *wholesaleOrderImpl) vendorSettleByRate(vendor merchant.IMerchant, rate float32) error {
	amount, refund := s.getOrderAmount()
	sAmount := amount * rate
	if sAmount > 0 {
		totalAmount := int(sAmount * float32(enum.RATE_Amount))
		refundAmount := int(refund * float32(enum.RATE_Amount))
		tradeFee, _ := vendor.SaleManager().MathTradeFee(
			merchant.TKWholesaleOrder, totalAmount)
		return vendor.Account().SettleOrder(s.OrderNo(),
			totalAmount, tradeFee, refundAmount, "批发订单结算")
	}
	return nil
}

// 完成订单
func (s *wholesaleOrderImpl) onOrderComplete() error {
	id := s.GetAggregateRootId()
	// 更新发货单
	soList := s.shipRepo.GetShipOrders(id, false)
	for _, v := range soList {
		domain.HandleError(v.Completed(), "domain")
	}
	// 更新会员账户
	err := s.updateAccountForOrder()
	if err == nil {
		// 处理返现
		//err = o.handleCashBack()
	}

	return err
}

// 更新账户
func (s *wholesaleOrderImpl) updateAccountForOrder() error {
	if s.value.State != order.StatCompleted {
		return order.ErrUnusualOrderStat
	}
	m := s.Buyer()
	var err error
	ov := s.value
	conf := s.valueRepo.GetGlobNumberConf()
	registry := s.valueRepo.GetRegistry()
	amount := ov.FinalAmount
	acc := m.GetAccount()

	// 增加经验
	if registry.MemberExperienceEnabled {
		rate := conf.ExperienceRateByOrder
		if exp := int32(amount * rate); exp > 0 {
			if err = m.AddExp(exp); err != nil {
				return err
			}
		}
	}

	// 增加积分
	//todo: 增加阶梯的返积分,比如订单满30送100积分
	integral := int64(amount*conf.IntegralRateByConsumption) + conf.IntegralBackExtra
	// 赠送积分
	if integral > 0 {
		err = m.GetAccount().AddIntegral(member.TypeIntegralShoppingPresent,
			s.value.OrderNo, integral, "")
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
func (s *wholesaleOrderImpl) LogBytes() []byte {
	buf := bytes.NewBufferString("")
	orderId := s.GetAggregateRootId()
	list := s.repo.GetSubOrderLogs(orderId)
	for _, v := range list {
		buf.WriteString(time.Unix(v.RecordTime, 0).Format("2006-01-02 15:04:05"))
		buf.WriteString("  ")
		if v.Message[:1] == "{" {
			if msg := s.getLogStringByStat(v.OrderState); len(msg) > 0 {
				v.Message = msg
			}
		}
		buf.WriteString(v.Message)
		buf.Write([]byte("\n"))
	}
	return buf.Bytes()
}

func (s *wholesaleOrderImpl) getLogStringByStat(stat int) string {
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
func (s *wholesaleOrderImpl) Cancel(reason string) error {
	if s.value.State == order.StatCancelled {
		return order.ErrOrderCancelled
	}
	// 已发货订单无法取消
	if s.value.State >= order.StatShipped {
		return order.ErrOrderShippedCancel
	}
	s.value.State = order.StatCancelled
	s.value.UpdateTime = time.Now().Unix()
	err := s.saveWholesaleOrder()
	if err == nil {
		domain.HandleError(s.AppendLog(order.LogSetup, true, reason), "domain")
		// 取消支付单
		err = s.cancelPaymentOrder()
		if err == nil {
			// 取消商品
			err = s.cancelGoods()
		}
	}
	return err
}

// 取消商品
func (s *wholesaleOrderImpl) cancelGoods() error {
	for _, v := range s.Items() {
		snapshot := s.itemRepo.GetSalesSnapshot(v.SnapshotId)
		if snapshot == nil {
			return item.ErrNoSuchSnapshot
		}
		gds := s.itemRepo.GetItem(snapshot.SkuId)
		if gds != nil {
			// 释放库存
			gds.FreeStock(v.SkuId, v.Quantity)
			// 如果订单已付款，则取消销售数量
			if s.value.IsPaid == 1 {
				gds.CancelSale(v.SkuId, v.Quantity, s.value.OrderNo)
			}
		}
	}
	return nil
}

// 获取支付单
func (s *wholesaleOrderImpl) GetPaymentOrder() payment.IPaymentOrder {
	if s.paymentOrder == nil {
		id := s.GetAggregateRootId()
		if id <= 0 {
			panic(" Get payment order error ; because of order no yet created!")
		}
		s.paymentOrder = s.payRepo.GetPaymentBySalesOrderId(id)
	}
	return s.paymentOrder
}

// 取消支付单
func (s *wholesaleOrderImpl) cancelPaymentOrder() error {
	po := s.GetPaymentOrder()
	if po != nil {
		v := po.GetValue()
		//if true {
		//	log.Println("支付单号为：", v.TradeNo, "; 金额：", v.FinalFee,
		//		"; 订单金额:", o.value.FinalFee)
		//}
		// 订单金额为0,则取消订单
		if v.FinalFee-s.value.FinalAmount <= 0 {
			return po.Cancel()
		}
		return po.Adjust(-s.value.FinalAmount)
	}
	return nil
}

// 谢绝订单
func (s *wholesaleOrderImpl) Decline(reason string) error {
	if s.value.State == order.StatAwaitingPayment {
		return s.Cancel("商户取消,原因:" + reason)
	}
	if s.value.State >= order.StatShipped ||
		s.value.State >= order.StatCancelled {
		return order.ErrOrderCancelled
	}
	s.value.State = order.StatDeclined
	s.value.UpdateTime = time.Now().Unix()
	return s.saveWholesaleOrder()
}
