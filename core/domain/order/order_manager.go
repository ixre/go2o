/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:49
 * description :
 * history :
 */

package order

import (
	"errors"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/delivery"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
)

var _ order.IOrderManager = new(orderManagerImpl)

type orderManagerImpl struct {
	repo         order.IOrderRepo
	caller       order.IUnifiedOrderAdapter
	productRepo  product.IProductRepo
	cartRepo     cart.ICartRepo
	goodsRepo    item.IGoodsItemRepo
	promRepo     promotion.IPromotionRepo
	memberRepo   member.IMemberRepo
	mchRepo      merchant.IMerchantRepo
	deliveryRepo delivery.IDeliveryRepo
	valRepo      valueobject.IValueRepo
	paymentRepo  payment.IPaymentRepo
	expressRepo  express.IExpressRepo
	mch          merchant.IMerchant
	shipRepo     shipment.IShipmentRepo
	breaker      *wholesaleOrderBreaker
}

func NewOrderManager(cartRepo cart.ICartRepo, mchRepo merchant.IMerchantRepo,
	repo order.IOrderRepo, payRepo payment.IPaymentRepo, productRepo product.IProductRepo,
	goodsRepo item.IGoodsItemRepo, promRepo promotion.IPromotionRepo,
	memberRepo member.IMemberRepo, deliveryRepo delivery.IDeliveryRepo,
	expressRepo express.IExpressRepo, shipRepo shipment.IShipmentRepo,
	valRepo valueobject.IValueRepo) order.IOrderManager {
	return &orderManagerImpl{
		repo:         repo,
		cartRepo:     cartRepo,
		productRepo:  productRepo,
		goodsRepo:    goodsRepo,
		promRepo:     promRepo,
		memberRepo:   memberRepo,
		paymentRepo:  payRepo,
		mchRepo:      mchRepo,
		deliveryRepo: deliveryRepo,
		valRepo:      valRepo,
		expressRepo:  expressRepo,
		shipRepo:     shipRepo,
		breaker:      newWholesaleOrderBreaker(repo),
	}
}

// 统一调用
func (o *orderManagerImpl) Unified(orderNo string, sub bool) order.IUnifiedOrderAdapter {
	u := &unifiedOrderAdapterImpl{
		repo:    o.repo,
		manager: o,
	}
	return u.adapter(orderNo, sub)
}

// 在下单前检查购物车
func (t *orderManagerImpl) checkCartForOrder(c cart.ICart) error {
	if c == nil {
		return cart.ErrEmptyShoppingCart
	}
	return c.Prepare()
}

// 预创建普通订单
func (t *orderManagerImpl) PrepareNormalOrder(c cart.ICart) (order.IOrder, error) {
	err := t.checkCartForOrder(c)
	if err != nil {
		return nil, err
	}
	orderType := order.TRetail
	switch c.Kind() {
	case cart.KNormal:
		orderType = order.TRetail
	case cart.KWholesale:
		orderType = order.TWholesale
	default:
		panic("not support cart kind parse to order")
	}
	val := &order.Order{
		BuyerId:   c.BuyerId(),
		OrderType: int32(orderType),
	}
	o := t.repo.CreateOrder(val)
	if o.Type() != order.TRetail {
		panic("only support normal order")
	}
	io := o.(order.INormalOrder)
	err = io.RequireCart(c)
	io.GetByVendor()
	return o, err
}

// 预创建批发订单
func (o *orderManagerImpl) PrepareWholesaleOrder(c cart.ICart) ([]order.IOrder, error) {
	if c.Kind() != cart.KWholesale {
		return nil, cart.ErrKindNotMatch
	}
	return o.breaker.BreakUp(c, nil)
}

// 提交批发订单
func (o *orderManagerImpl) SubmitWholesaleOrder(c cart.ICart,
	data order.IPostedData) (map[string]string, error) {
	if c.Kind() != cart.KWholesale {
		return nil, cart.ErrKindNotMatch
	}
	addressId := data.AddressId()
	if addressId <= 0 {
		return nil, order.ErrNoSuchAddress
	}
	checked := data.CheckedData()
	rd := map[string]string{
		"error": "",
	}

	list, err := o.breaker.BreakUp(c, data)
	for i, v := range list {
		err = o.submitSellerWholesaleOrder(v)
		if err != nil {
			return map[string]string{}, err
		}
		okOrder := o.GetOrderById(v.GetAggregateRootId())
		//返回订单号
		if i > 0 {
			rd["order_no"] += ","
		}
		rd["order_no"] += okOrder.OrderNo()
	}
	// 清空购物车
	if err == nil {
		if c.Release(checked) {
			c.Destroy()
		}
	}
	return rd, err
}

func (o *orderManagerImpl) submitSellerWholesaleOrder(v order.IOrder) error {
	err := v.Submit()
	if err == nil {
		//todo:???
		// 余额支付
		//py := io.GetPaymentOrder()
		//if useBalanceDiscount {
		//    py.BalanceDiscount("")
		//}
	}
	return err
}

// 提交交易类订单
func (t *orderManagerImpl) SubmitTradeOrder(c *order.ComplexOrder,
	tradeRate float64) (order.IOrder, error) {
	val := &order.Order{
		BuyerId:   c.BuyerId,
		OrderType: int32(order.TTrade),
	}
	o := t.repo.CreateOrder(val)
	io := o.(order.ITradeOrder)
	err := io.Set(c, tradeRate)
	if err == nil {
		err = o.Submit()
	}
	return o, err
}

func (t *orderManagerImpl) GetFreeOrderNo(vendorId int32) string {
	return t.repo.GetFreeOrderNo(vendorId)
}

// 应用优惠券
func (t *orderManagerImpl) applyCoupon(m member.IMember, o order.IOrder,
	py payment.IPaymentOrder, couponCode string) error {
	if o.Type() != order.TRetail {
		return errors.New("不支持优惠券")
	}
	io := o.(order.INormalOrder)
	po := py.Get()
	//todo: ?? 重构
	cp := t.promRepo.GetCouponByCode(
		int32(m.GetAggregateRootId()), couponCode)
	// 如果优惠券不存在
	if cp == nil {
		return errors.New("优惠券无效")
	}
	// 获取优惠券
	coupon := cp.(promotion.ICouponPromotion)
	result, err := coupon.CanUse(m, float32(po.TotalAmount/100))
	if result {
		if coupon.CanTake() {
			_, err = coupon.GetTake(m.GetAggregateRootId())
			//如果未占用，则占用
			if err != nil {
				err = coupon.Take(m.GetAggregateRootId())
			}
		} else {
			_, err = coupon.GetBind(m.GetAggregateRootId())
		}
		if err != nil {
			domain.HandleError(err, "domain")
			err = errors.New("优惠券无效")
		} else {
			//应用优惠券
			if err = io.ApplyCoupon(coupon); err == nil {
				_, err = py.CouponDiscount(coupon)
			}
		}
	}
	return err
}

func (t *orderManagerImpl) SubmitOrder(c cart.ICart, addressId int64,
	couponCode string, useBalanceDiscount bool) (order.IOrder, *order.SubmitReturnData, error) {
	rd := &order.SubmitReturnData{}
	o, err := t.PrepareNormalOrder(c)
	if err != nil {
		return nil, rd, err
	}
	if o.Type() != order.TRetail {
		panic("only support retail cart!")
	}
	io := o.(order.INormalOrder)
	buyer := o.Buyer()
	// 设置收货地址
	if err = io.SetAddress(addressId); err != nil {
		return o, rd, err
	} else {
		buyer.Profile().SetDefaultAddress(addressId) // 更新默认收货地址为本地使用地址
	}
	if err = o.Submit(); err != nil {
		return o, rd, err
	}
	// 合并支付
	iss := io.GetSubOrders()
	arr := make([]payment.IPaymentOrder, 0)
	for _, v := range iss {
		ip := v.GetPaymentOrder()
		if len(couponCode) != 0 { // 使用优惠码
			if err = t.applyCoupon(buyer, o, ip, couponCode); err != nil {
				return o, rd, err
			}
		}
		if useBalanceDiscount { // 使用余额抵扣
			if err = ip.BalanceDiscount(""); err != nil {
				return o, rd, err
			}
		}
		if ip.State() == payment.StateAwaitingPayment {
			arr = append(arr, ip)
		}
	}
	l := len(arr)
	// 如果全部支付成功
	if l == 0 {
		ip := iss[0].GetPaymentOrder().Get()
		rd.TradeNo = ip.TradeNo
		rd.TradeAmount = ip.FinalFee
		rd.OrderNo = ip.OutOrderNo
		return o, rd, err
	}
	// 剩下单个订单未支付
	if l == 1 {
		ip := arr[0].Get()
		rd.TradeNo = ip.TradeNo
		rd.TradeAmount = ip.FinalFee
		rd.OrderNo = ip.OutOrderNo
		return o, rd, err
	}
	// 合并支付
	mergeTradeNo, fee, err := arr[0].MergePay(arr[1:])
	if err != nil {
		return o, rd, err
	}
	//println("----", len(arr), "个订单已合并", fee, mergeTradeNo)
	rd.MergePay = true
	rd.TradeAmount = fee
	rd.TradeNo = mergeTradeNo
	for i, v := range arr {
		if i > 0 { // 拼接订单号
			rd.OrderNo += ","
		}
		rd.OrderNo += v.Get().OutOrderNo
	}
	return o, rd, err
}

// 根据订单编号获取订单
func (t *orderManagerImpl) GetOrderById(orderId int64) order.IOrder {
	val := t.repo.GetOrder("id=? LIMIT 1", orderId)
	if val != nil {
		return t.repo.CreateOrder(val)
	}
	return nil
}

// 根据订单号获取订单
func (t *orderManagerImpl) GetOrderByNo(orderNo string) order.IOrder {
	val := t.repo.GetOrder("order_no=?", orderNo)
	if val != nil {
		return t.repo.CreateOrder(val)
	}
	return nil
}

// 接收在线交易支付的通知，不主动调用
func (t *orderManagerImpl) NotifyOrderTradeSuccess(orderNo string, subOrder bool) error {
	if subOrder { // 处理子订单
		iso := t.repo.GetSubOrderByOrderNo(orderNo)
		return iso.PaymentFinishByOnlineTrade()
	}
	o := t.GetOrderByNo(orderNo)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	switch o.Type() {
	case order.TRetail:
		io := o.(order.INormalOrder)
		return io.OnlinePaymentTradeFinish()
	case order.TWholesale:
		io := o.(order.IWholesaleOrder)
		return io.OnlinePaymentTradeFinish()
	case order.TTrade:
		io := o.(order.ITradeOrder)
		return io.TradePaymentFinish()
	}
	panic("unknown order type")
}

// 获取子订单
func (t *orderManagerImpl) GetSubOrder(id int64) order.ISubOrder {
	if v := t.repo.GetSubOrder(id); v != nil {
		return t.repo.CreateNormalSubOrder(v)
	}
	return nil
}

var _ order.IUnifiedOrderAdapter = new(unifiedOrderAdapterImpl)

type unifiedOrderAdapterImpl struct {
	repo     order.IOrderRepo
	manager  order.IOrderManager
	bigOrder order.IOrder
	subOrder order.ISubOrder
	sub      bool
}

func (u *unifiedOrderAdapterImpl) adapter(orderNo string, sub bool) order.IUnifiedOrderAdapter {
	u.sub = sub
	if u.sub {
		u.subOrder = u.repo.GetSubOrderByOrderNo(orderNo)
	} else {
		orderId := u.repo.GetOrderId(orderNo, sub)
		u.bigOrder = u.manager.GetOrderById(orderId)
	}
	return u
}

func (u *unifiedOrderAdapterImpl) check() error {
	if u.sub && u.subOrder == nil {
		return order.ErrNoSuchOrder
	}
	if !u.sub && u.bigOrder == nil {
		return order.ErrNoSuchOrder
	}
	return nil
}

// 复合的订单信息
func (u *unifiedOrderAdapterImpl) Complex() *order.ComplexOrder {
	if err := u.check(); err == nil {
		if u.sub {
			return u.subOrder.Complex()
		}
		return u.bigOrder.Complex()
	}
	return nil
}

// 取消订单
func (u *unifiedOrderAdapterImpl) Cancel(reason string) error {
	if err := u.check(); err != nil {
		return err
	}
	if u.sub {
		return u.subOrder.Cancel(reason)
	}
	return u.cancel(reason)
}

func (u *unifiedOrderAdapterImpl) cancel(reason string) error {
	switch u.bigOrder.Type() {
	case order.TWholesale:
		return u.bigOrder.(order.IWholesaleOrder).Cancel(reason)
	}
	return nil
}

// 确定订单
func (u *unifiedOrderAdapterImpl) Confirm() error {
	if err := u.check(); err != nil {
		return err
	}
	if u.sub {
		return u.subOrder.Confirm()
	}
	return u.confirm()
}

func (u *unifiedOrderAdapterImpl) confirm() error {
	switch u.bigOrder.Type() {
	case order.TWholesale:
		return u.bigOrder.(order.IWholesaleOrder).Confirm()
	}
	return nil
}

// 备货完成
func (u *unifiedOrderAdapterImpl) PickUp() error {
	if err := u.check(); err != nil {
		return err
	}
	if u.sub {
		return u.subOrder.PickUp()
	}
	return u.pickup()
}

func (u *unifiedOrderAdapterImpl) pickup() error {
	switch u.bigOrder.Type() {
	case order.TWholesale:
		return u.bigOrder.(order.IWholesaleOrder).PickUp()
	}
	return nil
}

// 订单发货,并记录配送服务商编号及单号
func (u *unifiedOrderAdapterImpl) Ship(spId int32, spOrder string) error {
	if err := u.check(); err != nil {
		return err
	}
	if u.sub {
		return u.subOrder.Ship(spId, spOrder)
	}
	return u.ship(spId, spOrder)
}

func (u *unifiedOrderAdapterImpl) ship(spId int32, spOrder string) error {
	switch u.bigOrder.Type() {
	case order.TWholesale:
		return u.bigOrder.(order.IWholesaleOrder).Ship(spId, spOrder)
	}
	return nil
}

// 消费者收货
func (u *unifiedOrderAdapterImpl) BuyerReceived() error {
	if err := u.check(); err != nil {
		return err
	}
	if u.sub {
		return u.subOrder.BuyerReceived()
	}
	return u.buyerReceived()
}

func (u *unifiedOrderAdapterImpl) buyerReceived() error {
	switch u.bigOrder.Type() {
	case order.TWholesale:
		return u.bigOrder.(order.IWholesaleOrder).BuyerReceived()
	}
	return nil
}

// 获取订单日志
func (u *unifiedOrderAdapterImpl) LogBytes() []byte {
	if err := u.check(); err != nil {
		return []byte(nil)
	}
	if u.sub {
		return u.subOrder.LogBytes()
	}
	return u.logBytes()
}

func (u *unifiedOrderAdapterImpl) logBytes() []byte {
	switch u.bigOrder.Type() {
	case order.TWholesale:
		return u.bigOrder.(order.IWholesaleOrder).LogBytes()
	}
	return []byte(nil)
}
