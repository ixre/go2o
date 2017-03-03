package order

import (
	"errors"
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/infrastructure/domain"
	"log"
	"strconv"
	"strings"
)

var _ order.IOrder = new(wholesaleOrderImpl)
var _ order.IWholesaleOrder = new(wholesaleOrderImpl)

type wholesaleOrderImpl struct {
	*baseOrderImpl
	value        *order.WholesaleOrder
	items        []*orderItem
	paymentOrder payment.IPaymentOrder
	orderRepo    order.IOrderRepo
	expressRepo  express.IExpressRepo
	payRepo      payment.IPaymentRepo
	itemRepo     item.IGoodsItemRepo
}

func newWholesaleOrder(base *baseOrderImpl,
	shoppingRepo order.IOrderRepo, goodsRepo item.IGoodsItemRepo,
	expressRepo express.IExpressRepo, payRepo payment.IPaymentRepo) order.IOrder {
	return &wholesaleOrderImpl{
		baseOrderImpl: base,
		orderRepo:     shoppingRepo,
		itemRepo:      goodsRepo,
		expressRepo:   expressRepo,
		payRepo:       payRepo,
	}
}

func (o *wholesaleOrderImpl) getValue() *order.WholesaleOrder {
	if o.value == nil {
		id := o.GetAggregateRootId()
		if id > 0 {
			o.value = o.repo.GetWholesaleOrder("order_id=?", id)
		}
	}
	return o.value
}

// 设置商品项
func (w *wholesaleOrderImpl) SetItems(items []*order.MinifyItem) {
	if w.GetAggregateRootId() > 0 {
		panic("wholesale has created. can't use SetItems!")
	}
	w.parseOrder(items)
}

// 转换为订单相关对象
func (w *wholesaleOrderImpl) parseOrder(items []*order.MinifyItem) {
	w.value = &order.WholesaleOrder{
		ID:          0,
		OrderNo:     "",
		OrderId:     0,
		BuyerId:     w.baseValue.BuyerId,
		VendorId:    0,
		ShopId:      0,
		ItemAmount:  0,
		ExpressFee:  0,
		PackageFee:  0,
		FinalAmount: 0,
		State:       w.baseValue.State,
	}
	w.items = []*orderItem{}
	for _, v := range items {
		w.items = append(w.items, w.createItem(v))
	}
	// 获取运营商和商铺编号
	w.value.VendorId = w.items[0].VendorId
	w.value.ShopId = w.items[0].ShopId
	// 运费计算器
	ue := w.expressRepo.GetUserExpress(w.value.VendorId)
	ec := ue.CreateCalculator()
	// 计算订单金额及运费
	for _, item := range w.items {
		w.value.ItemAmount += item.Amount
		w.value.DiscountAmount += item.Amount - item.FinalAmount
		w.appendToExpressCalculator(ue, item, ec)
	}
	ec.Calculate("")
	w.value.ExpressFee = ec.Total()
	w.value.PackageFee = 0
	//计算最终金额
	w.value.FinalAmount = w.value.ItemAmount - w.value.DiscountAmount +
		w.value.ExpressFee + w.value.PackageFee
}

// 创建商品信息,并读取价格及运费信息
func (w *wholesaleOrderImpl) createItem(i *order.MinifyItem) *orderItem {
	// 获取商品信息
	it := w.itemRepo.GetItem(i.ItemId)
	sku := it.GetSku(i.SkuId)
	iv := it.GetValue()
	// 获取商品已销售快照
	snap := w.itemRepo.SnapshotService().GetLatestSalesSnapshot(
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
func (w *wholesaleOrderImpl) appendToExpressCalculator(ue express.IUserExpress,
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

// 复合的订单信息
func (o *wholesaleOrderImpl) Complex() *order.ComplexOrder {
	v := o.getValue()
	co := o.baseOrderImpl.Complex()
	co.ConsigneePerson = v.ConsigneePerson
	co.ConsigneePhone = v.ConsigneePhone
	co.ShippingAddress = v.ShippingAddress
	co.DiscountAmount = v.DiscountAmount
	co.ItemAmount = v.ItemAmount
	co.ExpressFee = v.ExpressFee
	co.PackageFee = v.PackageFee
	co.FinalAmount = v.FinalAmount
	co.IsBreak = 0
	co.UpdateTime = v.UpdateTime
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
		// 均摊优惠折扣到商品
		o.avgDiscountForItem()
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
		return member.ErrMemberDisabled
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

// 平均优惠抵扣金额到商品
func (o *wholesaleOrderImpl) avgDiscountForItem() {
	if o.items == nil {
		panic(errors.New("仅能在下单时进行商品抵扣平均"))
	}
	if o.value.DiscountAmount > 0 {
		totalFee := o.value.ItemAmount
		disFee := o.value.DiscountAmount
		for _, v := range o.items {
			b := (v.Amount / totalFee)
			v.FinalAmount = v.Amount - b*disFee
		}
	}
}

// 保存商品项
func (o *wholesaleOrderImpl) saveOrderItemsOnSubmit() (err error) {
	orderId := o.GetAggregateRootId()
	for _, v := range o.items {
		v.OrderId = orderId
		item := o.parseOrderItem(v)
		_, err = o.repo.SaveWholesaleItem(item)
		if err != nil {
			break
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
		SnapId:         int64(i.SnapshotId),
		Quantity:       i.Quantity,
		ReturnQuantity: i.ReturnQuantity,
		Amount:         i.Amount,
		FinalAmount:    i.FinalAmount,
		IsShipped:      0,
		UpdateTime:     i.UpdateTime,
	}
}

// 获取商品项
func (w *wholesaleOrderImpl) Items() []*order.WholesaleItem {
	panic("not implement")
}

// 设置配送地址
func (o *wholesaleOrderImpl) SetAddress(addressId int32) error {
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

// 生成支付单
func (o *wholesaleOrderImpl) createPaymentForOrder() error {
	v := o.baseOrderImpl.createPaymentOrder()
	v.VendorId = o.value.VendorId
	v.TotalFee = o.value.FinalAmount
	v.CouponDiscount = 0
	v.IntegralDiscount = 0
	v.FinalAmount = v.TotalFee - v.SubAmount - v.SystemDiscount -
		v.IntegralDiscount - v.BalanceDiscount
	po := o.payRepo.CreatePaymentOrder(v)
	_, err := po.Commit()
	return err
}

// 在线支付交易完成
func (w *wholesaleOrderImpl) PaymentFinishByOnlineTrade() error {
	panic("not implement")
}

// 记录订单日志
func (w *wholesaleOrderImpl) AppendLog(logType order.LogType,
	system bool, message string) error {
	panic("not implement")
}

// 添加备注
func (w *wholesaleOrderImpl) AddRemark(string) {
	panic("not implement")
}

// 确认订单
func (w *wholesaleOrderImpl) Confirm() error {
	panic("not implement")
}

// 捡货(备货)
func (w *wholesaleOrderImpl) PickUp() error {
	panic("not implement")
}

// 发货
func (w *wholesaleOrderImpl) Ship(spId int32, spOrder string) error {
	panic("not implement")
}

// 已收货
func (w *wholesaleOrderImpl) BuyerReceived() error {
	panic("not implement")
}

// 获取订单的日志
func (w *wholesaleOrderImpl) LogBytes() []byte {
	panic("not implement")
}

// 取消订单/退款
func (w *wholesaleOrderImpl) Cancel(reason string) error {
	panic("not implement")
}

// 谢绝订单
func (w *wholesaleOrderImpl) Decline(reason string) error {
	panic("not implement")
}

// 订单拆单
type wholesaleOrderBreaker struct {
	repo order.IOrderRepo
}

func newWholesaleOrderBreaker(repo order.IOrderRepo) *wholesaleOrderBreaker {
	return &wholesaleOrderBreaker{
		repo: repo,
	}
}

// 拆分订单
func (w *wholesaleOrderBreaker) BreakUp(c cart.IWholesaleCart) ([]order.IOrder, error) {
	items := c.GetValue().Items
	if len(items) == 0 {
		return []order.IOrder{}, cart.ErrEmptyShoppingCart
	}
	// 将购物车的商品按运营商拆分
	vendorItemsMap := w.breakVendorItemMap(items)
	if l := len(vendorItemsMap); l == 0 {
		return []order.IOrder{}, cart.ErrNoChecked
	}
	list := []order.IOrder{}
	cc := c.(cart.ICart)
	buyerId := cc.BuyerId()
	for _, items := range vendorItemsMap {
		o := w.createOrder(buyerId, items)
		list = append(list, o)
	}
	return list, nil
}

// 创建批发订单
func (w *wholesaleOrderBreaker) createOrder(buyerId int32, items []*cart.WsCartItem) order.IOrder {
	v := &order.Order{
		BuyerId:   buyerId,
		OrderType: int32(order.TWholesale),
		State:     int32(order.StatAwaitingPayment),
	}
	o := w.repo.CreateOrder(v)
	wo := o.(order.IWholesaleOrder)
	list := make([]*order.MinifyItem, len(items))
	for i, v := range items {
		list[i] = &order.MinifyItem{
			ItemId:   v.ItemId,
			SkuId:    v.SkuId,
			Quantity: v.Quantity,
		}
	}
	wo.SetItems(list)
	return o
}

// 生成运营商与订单商品的映射
func (w *wholesaleOrderBreaker) breakVendorItemMap(items []*cart.WsCartItem) map[int32][]*cart.WsCartItem {
	mp := make(map[int32][]*cart.WsCartItem)
	for _, v := range items {
		//必须勾选为结算
		if v.Checked == 1 {
			list, ok := mp[v.VendorId]
			if !ok {
				list = []*cart.WsCartItem{}
			}
			mp[v.VendorId] = append(list, v)
		}
	}
	return mp
}
