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
	"fmt"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/delivery"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/lbs"
	"go2o/core/infrastructure/log"
	"sync"
	"time"
)

var _ order.IOrderManager = new(orderManagerImpl)

type orderManagerImpl struct {
	_rep         order.IOrderRep
	_saleRep     sale.ISaleRep
	_cartRep     cart.ICartRep
	_goodsRep    goods.IGoodsRep
	_promRep     promotion.IPromotionRep
	_memberRep   member.IMemberRep
	_mchRep      merchant.IMerchantRep
	_deliveryRep delivery.IDeliveryRep
	_valRep      valueobject.IValueRep
	_payRep      payment.IPaymentRep
	_expressRep  express.IExpressRep
	_merchant    merchant.IMerchant
	_shipRep     shipment.IShipmentRep
}

func NewOrderManager(cartRep cart.ICartRep, mchRep merchant.IMerchantRep,
	rep order.IOrderRep, payRep payment.IPaymentRep, saleRep sale.ISaleRep,
	goodsRep goods.IGoodsRep, promRep promotion.IPromotionRep,
	memberRep member.IMemberRep, deliveryRep delivery.IDeliveryRep,
	expressRep express.IExpressRep, shipRep shipment.IShipmentRep,
	valRep valueobject.IValueRep) order.IOrderManager {

	return &orderManagerImpl{
		_rep:         rep,
		_cartRep:     cartRep,
		_saleRep:     saleRep,
		_goodsRep:    goodsRep,
		_promRep:     promRep,
		_memberRep:   memberRep,
		_payRep:      payRep,
		_mchRep:      mchRep,
		_deliveryRep: deliveryRep,
		_valRep:      valRep,
		_expressRep:  expressRep,
		_shipRep:     shipRep,
	}
}

// 生成订单
func (this *orderManagerImpl) CreateOrder(val *order.Order) order.IOrder {
	return newOrder(this, val, this._mchRep,
		this._rep, this._goodsRep, this._saleRep, this._promRep,
		this._memberRep, this._expressRep, this._valRep)
}

// 生成空白订单,并保存返回对象
func (this *orderManagerImpl) CreateSubOrder(v *order.SubOrder) order.ISubOrder {
	return NewSubOrder(v, this, this._rep, this._memberRep,
		this._goodsRep, this._shipRep, this._saleRep,
		this._valRep, this._mchRep)
}

// 在下单前检查购物车
func (this *orderManagerImpl) checkCartForOrder(c cart.ICart) error {
	if c == nil {
		return cart.ErrEmptyShoppingCart
	}
	return c.Check()
}

// 将购物车转换为订单
func (this *orderManagerImpl) ParseToOrder(c cart.ICart) (order.IOrder,
	member.IMember, error) {
	var m member.IMember
	err := this.checkCartForOrder(c)
	if err != nil {
		return nil, m, err
	}
	val := &order.Order{}

	// 判断购买会员
	buyerId := c.GetValue().BuyerId
	if buyerId > 0 {
		val.BuyerId = buyerId
		m = this._memberRep.GetMember(val.BuyerId)
	}
	if m == nil {
		return nil, m, member.ErrNoSuchMember
	}
	val.State = order.StatAwaitingPayment
	o := this.CreateOrder(val)
	err = o.RequireCart(c)
	o.GetByVendor()
	return o, m, err
}

// 预生成订单及支付单
func (this *orderManagerImpl) PrepareOrder(c cart.ICart, subject string,
	couponCode string) (order.IOrder, payment.IPaymentOrder, error) {
	//todo: subject 或备注先不理会,可能是多个note。且在下单后再提交备注
	order, m, err := this.ParseToOrder(c)
	var py payment.IPaymentOrder
	if err == nil {
		py = this.createPaymentOrder(m, order)
		//todo:
		//val := order.GetValue()
		//if len(subject) > 0 {
		//	val.Subject = subject
		//	order.SetValue(val)
		//}
		if len(couponCode) != 0 {
			err = this.applyCoupon(m, order, py, couponCode)
		}
	}
	return order, py, err
}

func (this *orderManagerImpl) GetFreeOrderNo(vendorId int) string {
	return this._rep.GetFreeOrderNo(vendorId)
}

// 智能选择门店
func (this *orderManagerImpl) SmartChoiceShop(address string) (shop.IShop, error) {
	//todo: 应只选择线下实体店
	//todo: AggregateRootId
	dly := this._deliveryRep.GetDelivery(-1)

	lng, lat, err := lbs.GetLocation(address)
	if err != nil {
		return nil, errors.New("无法识别的地址：" + address)
	}
	var cov delivery.ICoverageArea = dly.GetNearestCoverage(lng, lat)
	if cov == nil {
		return nil, delivery.ErrNotCoveragedArea
	}
	shopId, _, err := dly.GetDeliveryInfo(cov.GetDomainId())
	return this._merchant.ShopManager().GetShop(shopId), err
}

// 生成支付单
func (this *orderManagerImpl) createPaymentOrder(m member.IMember,
	o order.IOrder) payment.IPaymentOrder {
	val := o.GetValue()
	v := &payment.PaymentOrderBean{
		BuyUser:     m.GetAggregateRootId(),
		PaymentUser: m.GetAggregateRootId(),
		VendorId:    0,
		OrderId:     0,
		// 支付单金额
		TotalFee: val.FinalAmount,
		// 余额抵扣
		BalanceDiscount: 0,
		// 积分抵扣
		IntegralDiscount: 0,
		// 系统支付抵扣金额
		SystemDiscount: 0,
		// 优惠券金额
		CouponDiscount: 0,
		// 立减金额
		SubFee: 0,
		// 支付选项
		PaymentOpt: payment.OptPerm,
		// 支付方式
		PaymentSign: enum.PaymentOnlinePay,
		//创建时间
		CreateTime: time.Now().Unix(),
		// 在线支付的交易单号
		OuterNo: "",
		//支付时间
		PaidTime: 0,
		// 状态:  0为未付款，1为已付款，2为已取消
		State: payment.StateNotYetPayment,
	}
	v.FinalFee = v.TotalFee - v.SubFee - v.SystemDiscount -
		v.IntegralDiscount - v.BalanceDiscount
	return this._payRep.CreatePaymentOrder(v)
}

// 应用优惠券
func (this *orderManagerImpl) applyCoupon(m member.IMember, order order.IOrder,
	py payment.IPaymentOrder, couponCode string) error {
	po := py.GetValue()
	cp := this._promRep.GetCouponByCode(
		m.GetAggregateRootId(), couponCode)
	// 如果优惠券不存在
	if cp == nil {
		return errors.New("优惠券无效")
	}
	// 获取优惠券
	coupon := cp.(promotion.ICouponPromotion)
	result, err := coupon.CanUse(m, po.TotalFee)
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
			if err = order.ApplyCoupon(coupon); err == nil {
				_, err = py.CouponDiscount(coupon)
			}
		}
	}
	return err
}

func (this *orderManagerImpl) SubmitOrder(c cart.ICart, subject string,
	couponCode string, useBalanceDiscount bool) (order.IOrder,
	payment.IPaymentOrder, error) {
	order, py, err := this.PrepareOrder(c, subject, couponCode)
	if err != nil {
		return order, py, err
	}
	orderNo, err := order.Submit()
	tradeNo := orderNo
	if err == nil {
		cv := c.GetValue()
		cv.PaymentOpt = enum.PaymentOnlinePay
		pyUpdate := false
		//todo: 设置配送门店
		//err = order.SetShop(cv.ShopId)
		//err = order.SetDeliver(cv.DeliverId)

		// 设置支付方式
		if err = py.SetPaymentSign(cv.PaymentOpt); err != nil {
			return order, py, err
		}

		// 处理支付单
		py.BindOrder(order.GetAggregateRootId(), tradeNo)
		if _, err = py.Save(); err != nil {
			err = errors.New("下单出错:" + err.Error())
			//todo: 取消订单
			//order.Cancel(err.Error())
			domain.HandleError(err, "domain")
			return order, py, err
		}

		// 使用余额支付
		if useBalanceDiscount {
			err = py.BalanceDiscount()
			pyUpdate = true
		}

		// 如果已支付完成,则将订单设为支付完成
		if v := py.GetValue(); v.FinalFee == 0 &&
			v.State == payment.StateFinishPayment {
			for _, sub := range order.GetSubOrders() {
				sub.PaymentFinishByOnlineTrade()
			}
		}

		// 更新支付单
		if err == nil && pyUpdate {
			_, err = py.Save()
		}
	}
	return order, py, err
}

// 根据订单编号获取订单
func (this *orderManagerImpl) GetOrderById(orderId int) order.IOrder {
	val := this._rep.GetOrderById(orderId)
	if val != nil {
		return this.CreateOrder(val)
	}
	return nil
}

// 根据订单号获取订单
func (this *orderManagerImpl) GetOrderByNo(orderNo string) order.IOrder {
	val := this._rep.GetValueOrderByNo(orderNo)
	if val != nil {
		return this.CreateOrder(val)
	}
	return nil
}

// 在线交易支付
func (this *orderManagerImpl) PaymentForOnlineTrade(orderId int) error {
	o := this.GetOrderById(orderId)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.OnlinePaymentTradeFinish()
}

// 获取子订单
func (this *orderManagerImpl) GetSubOrder(id int) order.ISubOrder {
	if v := this._rep.GetSubOrder(id); v != nil {
		return this.CreateSubOrder(v)
	}
	return nil
}

// 根据父订单编号获取购买的商品项
func (this *orderManagerImpl) GetItemsByParentOrderId(orderId int) []*order.OrderItem {
	return this._rep.GetItemsByParentOrderId(orderId)
}

var (
	shopLocker sync.Mutex
	biShops    []shop.IShop
)

// 自动设置订单
func (this *orderManagerImpl) OrderAutoSetup(f func(error)) {
	var orders []*order.Order
	var err error

	shopLocker.Lock()
	defer func() {
		shopLocker.Unlock()
	}()
	biShops = nil
	log.Println("[SETUP] start auto setup")

	saleConf := this._merchant.ConfManager().GetSaleConf()
	if saleConf.AutoSetupOrder == 1 {
		orders, err = this._rep.GetWaitingSetupOrders(-1)
		if err != nil {
			f(err)
			return
		}

		dt := time.Now()
		for _, v := range orders {
			this.setupOrder(v, &saleConf, dt, f)
		}
	}
}

const (
	order_timeout_hour   = 24
	order_confirm_minute = 4
	order_process_minute = 11
	order_sending_minute = 31
	order_receive_hour   = 5
	order_complete_hour  = 11
)

func (this *orderManagerImpl) SmartConfirmOrder(o order.IOrder) error {

	return nil

	//todo:  自动确认订单
	var err error
	v := o.GetValue()
	log.Printf("[ AUTO][OrderSetup]:%s - Confirm \n", v.OrderNo)
	var sp shop.IShop
	if biShops == nil {
		// /pay/return_alipay?out_trade_no=ZY1607375766&request_token=requestToken&result=success&trade_no
		// =2016070221001004880246862127&sign=75a18ca0d75750ac22fedbbe6468c187&sign_type=MD5
		//todo:  拆分订单
		biShops = this._merchant.ShopManager().GetBusinessInShops()
	}
	if len(biShops) == 1 {
		sp = biShops[0]
	} else {
		sp, err = this.SmartChoiceShop(v.ShippingAddress)
		if err != nil {
			//todo:
			panic("not impl")
			//order.Suspend("智能分配门店失败！原因：" + err.Error())
			return err
		}
	}

	if sp != nil && sp.Type() == shop.TypeOfflineShop {
		sv := sp.GetValue()
		//todo: set shop
		panic("not impl")
		//order.SetShop(sp.GetDomainId())
		err = o.Confirm()
		//err = order.Process()
		ofs := sp.(shop.IOfflineShop).GetShopValue()
		o.AppendLog(&order.OrderLog{
			Type:     int(order.LogSetup),
			IsSystem: 1,
			Message:  fmt.Sprintf("自动分配门店:%s,电话：%s", sv.Name, ofs.Tel),
		})
	}
	return err
}

func (this *orderManagerImpl) setupOrder(v *order.Order,
	conf *merchant.SaleConf, t time.Time, f func(error)) {
	var err error
	od := this.CreateOrder(v)
	dur := time.Duration(t.Unix()-v.CreateTime) * time.Second

	switch v.State {
	case order.StatAwaitingPayment:
		if v.IsPaid == 0 && dur > time.Minute*time.Duration(conf.OrderTimeOutMinute) {
			//todo: del

			//order.Cancel("超时未付款，系统取消")
			log.Printf("[ AUTO][OrderSetup]:%s - Payment Timeout\n", v.OrderNo)
		}

	case enum.ORDER_WAIT_CONFIRM:
		if dur > time.Minute*time.Duration(conf.OrderConfirmAfterMinute) {
			err = this.SmartConfirmOrder(od)
		}

	//		case enum.ORDER_WAIT_DELIVERY:
	//			if dur > time.Minute*order_process_minute {
	//				err = order.Process()
	//				if ctx.Debug() {
	//					ctx.Log().Printf("[ AUTO][OrderSetup]:%s - Processing \n", v.OrderNo)
	//				}
	//			}

	//		case enum.ORDER_WAIT_RECEIVE:
	//			if dur > time.Hour * conf.OrderTimeOutReceiveHour {
	//				err = order.Deliver()
	//				if ctx.Debug() {
	//					ctx.Log().Printf("[ AUTO][OrderSetup]:%s - Sending \n", v.OrderNo)
	//				}
	//			}
	case enum.ORDER_WAIT_RECEIVE:
		if dur > time.Hour*time.Duration(conf.OrderTimeOutReceiveHour) {
			//todo:
			panic("not impl")
			//err = order.SignReceived()

			log.Printf("[ AUTO][OrderSetup]:%s - Received \n", v.OrderNo)
			if err == nil {
				//todo: del
				panic("not impl")
				//err = order.Complete()
				log.Printf("[ AUTO][OrderSetup]:%s - Complete \n", v.OrderNo)
			}
		}

		//		case enum.ORDER_COMPLETED:
		//			if dur > time.Hour*order_complete_hour {
		//				err = order.Complete()
		//				if ctx.Debug() {
		//					ctx.Log().Printf("[ AUTO][OrderSetup]:%s - Complete \n", v.OrderNo)
		//				}
		//			}
	}

	if err != nil {
		f(err)
	}
}
