package order

import (
	"errors"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"log"
	"strconv"
)

var _ order.IOrder = new(wholesaleOrderImpl)
var _ order.IWholesaleOrder = new(wholesaleOrderImpl)

type wholesaleOrderImpl struct {
	*baseOrderImpl
	manager order.IOrderManager
	value   *order.WholesaleOrder
	items   []*orderItem

	cart            cart.ICart //购物车,仅在订单生成时设置
	paymentOrder    payment.IPaymentOrder
	coupons         []promotion.ICouponPromotion
	availPromotions []promotion.IPromotion
	orderPbs        []*order.OrderPromotionBind
	orderRepo       order.IOrderRepo
	expressRepo     express.IExpressRepo
	payRepo         payment.IPaymentRepo
	goodsRepo       item.IGoodsItemRepo
	productRepo     product.IProductRepo
	promRepo        promotion.IPromotionRepo
	valRepo         valueobject.IValueRepo
	// 运营商商品映射,用于整理购物车
	vendorItemsMap map[int32][]*order.SubOrderItem
	// 运营商与邮费的MAP
	vendorExpressMap map[int32]float32
	// 是否为内部挂起
	internalSuspend bool
	subList         []order.ISubOrder
}

func newWholesaleOrder(shopping order.IOrderManager, base *baseOrderImpl,
	shoppingRepo order.IOrderRepo, goodsRepo item.IGoodsItemRepo,
	productRepo product.IProductRepo, promRepo promotion.IPromotionRepo,
	expressRepo express.IExpressRepo, payRepo payment.IPaymentRepo,
	valRepo valueobject.IValueRepo) order.IOrder {
	return &wholesaleOrderImpl{
		baseOrderImpl: base,
		manager:       shopping,
		//value:       value,
		promRepo:    promRepo,
		orderRepo:   shoppingRepo,
		goodsRepo:   goodsRepo,
		productRepo: productRepo,
		valRepo:     valRepo,
		expressRepo: expressRepo,
		payRepo:     payRepo,
	}
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
	it := w.goodsRepo.GetItem(i.ItemId)
	sku := it.GetSku(i.SkuId)
	iv := it.GetValue()
	// 获取商品已销售快照
	snap := w.goodsRepo.SnapshotService().GetLatestSalesSnapshot(
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
func (w *wholesaleOrderImpl) Complex() *order.ComplexOrder {
	v := w.value
	co := w.baseOrderImpl.Complex()
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
func (w *wholesaleOrderImpl) Submit() error {
	err := w.baseOrderImpl.Submit()
	return err
}

// 获取商品项
func (w *wholesaleOrderImpl) Items() []*order.OrderWholesaleItem {
	panic("not implement")
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
