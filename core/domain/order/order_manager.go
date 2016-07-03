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
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/valueobject"
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
	_partnerRep  merchant.IMerchantRep
	_deliveryRep delivery.IDeliveryRep
	_valRep      valueobject.IValueRep
	_merchant    merchant.IMerchant
}

func NewOrderManager(cartRep cart.ICartRep, partnerRep merchant.IMerchantRep,
	rep order.IOrderRep, saleRep sale.ISaleRep, goodsRep goods.IGoodsRep,
	promRep promotion.IPromotionRep, memberRep member.IMemberRep,
	deliveryRep delivery.IDeliveryRep, valRep valueobject.IValueRep) order.IOrderManager {

	return &orderManagerImpl{
		_rep:         rep,
		_cartRep:     cartRep,
		_saleRep:     saleRep,
		_goodsRep:    goodsRep,
		_promRep:     promRep,
		_memberRep:   memberRep,
		_partnerRep:  partnerRep,
		_deliveryRep: deliveryRep,
		_valRep:      valRep,
	}
}

func (this *orderManagerImpl) CreateOrder(val *order.ValueOrder,
	cart cart.ICart) order.IOrder {
	return newOrder(this, val, cart, this._partnerRep,
		this._rep, this._goodsRep, this._saleRep, this._promRep,
		this._memberRep, this._valRep)
}

// 将购物车转换为订单
func (this *orderManagerImpl) ParseToOrder(c cart.ICart) (order.IOrder,
	member.IMember, error) {
	val := &order.ValueOrder{}
	var m member.IMember
	var err error

	if c == nil {
		return nil, m, cart.ErrEmptyShoppingCart
	}
	if err = c.Check(); err != nil {
		return nil, m, err
	}
	// 判断购买会员
	val.BuyerId = c.GetValue().BuyerId
	if val.BuyerId > 0 {
		m = this._memberRep.GetMember(val.BuyerId)
	}
	if m == nil {
		return nil, m, member.ErrSessionTimeout
	}

	val.VendorId = -1

	tf, of := c.GetFee()
	val.TotalFee = tf //总金额
	val.Fee = of      //实际金额
	val.PayFee = of
	val.DiscountFee = tf - of //优惠金额
	val.VendorId = -1
	val.Status = 1

	o := this.CreateOrder(val, c)
	return o, m, nil
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

// 生成订单
func (this *orderManagerImpl) BuildOrder(c cart.ICart, subject string, couponCode string) (
	order.IOrder, error) {
	order, m, err := this.ParseToOrder(c)
	if err != nil {
		return order, err
	}
	var val = order.GetValue()
	if len(subject) > 0 {
		val.Subject = subject
		order.SetValue(&val)
	}

	if len(couponCode) != 0 {
		var coupon promotion.ICouponPromotion
		var result bool
		cp := this._promRep.GetCouponByCode(
			m.GetAggregateRootId(), couponCode)

		// 如果优惠券不存在
		if cp == nil {
			log.Error(err)
			return order, errors.New("优惠券无效")
		}

		coupon = cp.(promotion.ICouponPromotion)
		result, err = coupon.CanUse(m, val.Fee)
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
				log.Error(err)
				return order, errors.New("优惠券无效")
			}
			err = order.ApplyCoupon(coupon) //应用优惠券
		}
	}

	return order, err
}

func (this *orderManagerImpl) SubmitOrder(c cart.ICart, subject string,
	couponCode string, useBalanceDiscount bool) (string, error) {
	order, err := this.BuildOrder(c, subject, couponCode)
	if err != nil {
		return "", err
	}
	var cv = c.GetValue()
	if err == nil {
		err = order.SetShop(cv.ShopId)
		if err == nil {
			order.SetPayment(cv.PaymentOpt)
			err = order.SetDeliver(cv.DeliverId)
			if useBalanceDiscount {
				order.UseBalanceDiscount()
			}
			if err == nil {
				return order.Submit()
			}
		}
	}
	return "", err
}

func (this *orderManagerImpl) GetOrderByNo(orderNo string) order.IOrder {
	val := this._rep.GetValueOrderByNo(orderNo)
	if val != nil {
		val.Items = this._rep.GetOrderItems(val.Id)
		return this.CreateOrder(val, nil)
	}
	return nil
}

var (
	shopLocker sync.Mutex
	biShops    []shop.IShop
)

// 自动设置订单
func (this *orderManagerImpl) OrderAutoSetup(f func(error)) {
	var orders []*order.ValueOrder
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

func (this *orderManagerImpl) SmartConfirmOrder(order order.IOrder) error {
	var err error
	v := order.GetValue()
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
		sp, err = this.SmartChoiceShop(v.DeliverAddress)
		if err != nil {
			order.Suspend("智能分配门店失败！原因：" + err.Error())
			return err
		}
	}

	if sp != nil && sp.Type() == shop.TypeOfflineShop {
		sv := sp.GetValue()
		order.SetShop(sp.GetDomainId())
		err = order.Confirm()
		//err = order.Process()
		ofs := sp.(shop.IOfflineShop).GetShopValue()
		order.AppendLog(enum.ORDER_LOG_SETUP, false, fmt.Sprintf(
			"自动分配门店:%s,电话：%s", sv.Name, ofs.Tel))
	}
	return err
}

func (this *orderManagerImpl) setupOrder(v *order.ValueOrder,
	conf *merchant.SaleConf, t time.Time, f func(error)) {
	var err error
	order := this.CreateOrder(v, nil)
	dur := time.Duration(t.Unix()-v.CreateTime) * time.Second

	switch v.Status {
	case enum.ORDER_WAIT_PAYMENT:
		if v.IsPaid == 0 && dur > time.Minute*time.Duration(conf.OrderTimeOutMinute) {
			order.Cancel("超时未付款，系统取消")
			log.Printf("[ AUTO][OrderSetup]:%s - Payment Timeout\n", v.OrderNo)
		}

	case enum.ORDER_WAIT_CONFIRM:
		if dur > time.Minute*time.Duration(conf.OrderConfirmAfterMinute) {
			err = this.SmartConfirmOrder(order)
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
			err = order.SignReceived()

			log.Printf("[ AUTO][OrderSetup]:%s - Received \n", v.OrderNo)
			if err == nil {
				err = order.Complete()
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
