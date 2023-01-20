/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:49
 * description :
 * history :
 */

package order

import (
	"errors"

	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/core/domain/interface/delivery"
	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/domain/interface/promotion"
	"github.com/ixre/go2o/core/domain/interface/shipment"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/infrastructure/domain"
)

var _ order.IOrderManager = new(orderManagerImpl)

type orderManagerImpl struct {
	repo         order.IOrderRepo
	caller       order.IUnifiedOrderAdapter
	productRepo  product.IProductRepo
	cartRepo     cart.ICartRepo
	goodsRepo    item.IItemRepo
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
	goodsRepo item.IItemRepo, promRepo promotion.IPromotionRepo,
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

// Unified 统一调用
func (t *orderManagerImpl) Unified(orderNo string, sub bool) order.IUnifiedOrderAdapter {
	u := &unifiedOrderAdapterImpl{
		repo:    t.repo,
		manager: t,
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

// PrepareNormalOrder 预创建普通订单
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
		return nil, errors.New("not support cart kind parse to order")
	}
	val := &order.Order{
		BuyerId:   c.BuyerId(),
		OrderType: int(orderType),
	}
	o := t.repo.CreateOrder(val)
	if o.Type() != order.TRetail {
		return nil, errors.New("only support normal order")
	}
	io := o.(order.INormalOrder)
	err = io.RequireCart(c)
	if err == nil {
		io.GetByVendor()
	}
	return o, err
}

// PrepareWholesaleOrder 预创建批发订单
func (t *orderManagerImpl) PrepareWholesaleOrder(c cart.ICart) ([]order.IOrder, error) {
	if c.Kind() != cart.KWholesale {
		return nil, cart.ErrKindNotMatch
	}
	return t.breaker.BreakUp(c, nil)
}

// SubmitWholesaleOrder 提交批发订单
func (t *orderManagerImpl) submitWholesaleOrder(data order.SubmitOrderData) (order.IOrder, *order.SubmitReturnData, error) {
	rd := &order.SubmitReturnData{}

	ic := t.cartRepo.GetMyCart(data.BuyerId, cart.KWholesale)

	addressId := data.PostedData.AddressId()
	if addressId <= 0 {
		return nil, nil, order.ErrNoSuchAddress
	}
	checked := data.PostedData.CheckedData()

	list, err := t.breaker.BreakUp(ic, data.PostedData)
	for i, v := range list {
		err = t.submitSellerWholesaleOrder(v)
		if err != nil {
			return nil, nil, err
		}
		okOrder := t.GetOrderById(v.GetAggregateRootId())
		//返回订单号
		if i > 0 {
			rd.OrderNo += ","
		}
		rd.OrderNo += okOrder.OrderNo()
	}
	// 清空购物车
	if err == nil {
		if ic.Release(checked) {
			ic.Destroy()
		}
	}
	return nil, rd, err
}

func (t *orderManagerImpl) submitSellerWholesaleOrder(v order.IOrder) error {
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

// SubmitTradeOrder 提交交易类订单
func (t *orderManagerImpl) submitTradeOrder(data order.SubmitOrderData) (order.IOrder, *order.SubmitReturnData, error) {
	rd := &order.SubmitReturnData{}
	val := &order.Order{
		BuyerId:   int64(data.BuyerId),
		OrderType: int(order.TTrade),
	}
	o := t.repo.CreateOrder(val)
	io := o.(order.ITradeOrder)
	c := &order.TradeOrderValue{
		BuyerId:        int(data.BuyerId),
		StoreId:        int(data.PostedData.TradeOrderStoreId()),
		Subject:        data.Subject,
		ItemAmount:     int(data.PostedData.TradeOrderAmount()),
		DiscountAmount: 0, //todo: 需要支持用券码抵扣
	}
	err := io.Set(c, float64(data.PostedData.TradeOrderDiscount()))
	if err == nil {
		err = o.Submit()
		if err == nil {
			rd.OrderNo = o.OrderNo()
			rd.PaymentOrderNo = o.GetPaymentOrder().TradeNo()
		}
	}

	return o, rd, err
}

func (t *orderManagerImpl) GetFreeOrderNo(vendorId int64) string {
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
		m.GetAggregateRootId(), couponCode)
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

func (t *orderManagerImpl) SubmitOrder(data order.SubmitOrderData) (order.IOrder, *order.SubmitReturnData, error) {
	switch data.Type {
	case order.TRetail:
		return t.submitNormalOrder(data)
	case order.TWholesale:
		return t.submitWholesaleOrder(data)
	case order.TTrade:
		return t.submitTradeOrder(data)
	}
	return nil, nil, errors.New("not support order type")
}

func (t *orderManagerImpl) submitNormalOrder(data order.SubmitOrderData) (order.IOrder, *order.SubmitReturnData, error) {
	rd := &order.SubmitReturnData{}
	ic := t.cartRepo.GetMyCart(data.BuyerId, cart.KNormal)
	o, err := t.PrepareNormalOrder(ic)
	if err != nil {
		return nil, rd, err
	}
	buyer := o.Buyer()
	// 设置收货地址
	if err = o.SetShipmentAddress(data.AddressId); err != nil {
		return o, rd, err
	} else {
		_ = buyer.Profile().SetDefaultAddress(data.AddressId) // 更新默认收货地址为本地使用地址
	}
	// 使用返利用户代码
	no := o.(order.INormalOrder)
	if no != nil {
		if len(data.AffiliateCode) > 0 {
			_ = no.ApplyTraderCode(data.AffiliateCode)
		}
	}
	// 提交订单
	if err = o.Submit(); err != nil {
		return o, rd, err
	}
	// 合并支付
	ip := o.GetPaymentOrder()
	ipv := ip.Get()
	if len(data.CouponCode) != 0 { // 使用优惠码
		if err = t.applyCoupon(buyer, o, ip, data.CouponCode); err != nil {
			return o, rd, err
		}
	}
	// 使用余额抵扣,如果余额抵扣失败,仍然应该继续结算
	if data.BalanceDiscount {
		_ = ip.BalanceDiscount("")
	}
	// 如果全部支付成功
	if ip.State() > payment.StateAwaitingPayment {

	}

	rd.TradeNo = ipv.TradeNo
	rd.TradeAmount = ipv.FinalFee
	rd.OrderNo = ipv.OutOrderNo
	rd.PaymentOrderNo = o.GetPaymentOrder().TradeNo()
	return o, rd, err
	// 剩下单个订单未支付

	// // 合并支付
	// mergeTradeNo, fee, err := arr[0].MergePay(arr[1:])
	// if err != nil {
	// 	return o, rd, err
	// }
	// //println("----", len(arr), "个订单已合并", fee, mergeTradeNo)
	// rd.MergePay = true
	// rd.TradeAmount = int64(fee)
	// rd.TradeNo = mergeTradeNo
	// for i, v := range arr {
	// 	if i > 0 { // 拼接订单号
	// 		rd.OrderNo += ","
	// 	}
	// 	rd.OrderNo += v.Get().OutOrderNo
	// }
	//return o, rd, err
}

// 根据订单编号获取订单
func (t *orderManagerImpl) GetOrderById(orderId int64) order.IOrder {
	val := t.repo.GetOrder("id= $1 LIMIT 1", orderId)
	if val != nil {
		return t.repo.CreateOrder(val)
	}
	return nil
}

// 根据订单号获取订单
func (t *orderManagerImpl) GetOrderByNo(orderNo string) order.IOrder {
	val := t.repo.GetOrder("order_no = $1", orderNo)
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
	// 主动调用订单支付完成
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

// ChangeShipmentAddress 更改收货人信息
func (u *unifiedOrderAdapterImpl) ChangeShipmentAddress(addressId int64) error {
	//todo: 子订单改一个全部都改了
	if u.sub {
		return u.subOrder.ChangeShipmentAddress(addressId)
	}
	return u.bigOrder.ChangeShipmentAddress(addressId)
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

// Complex 复合的订单信息
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
func (u *unifiedOrderAdapterImpl) Cancel(buyerCancel bool,reason string) error {
	if err := u.check(); err != nil {
		return err
	}
	if u.sub {
		return u.subOrder.Cancel(buyerCancel,reason)
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

// Forbid implements order.IUnifiedOrderAdapter
func (u *unifiedOrderAdapterImpl) Forbid() error {
	err := u.check()
	if err == nil {
		return u.subOrder.Forbid()
	}
	return errors.New("not implemented")
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
