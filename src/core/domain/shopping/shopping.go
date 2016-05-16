/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:49
 * description :
 * history :
 */

package shopping

import (
	"errors"
	"fmt"
	"go2o/src/core/domain/interface/delivery"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/merchant"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/infrastructure/lbs"
	"go2o/src/core/infrastructure/log"
	"sync"
	"time"
)

type Shopping struct {
	_rep         shopping.IShoppingRep
	_saleRep     sale.ISaleRep
	_goodsRep    sale.IGoodsRep
	_promRep     promotion.IPromotionRep
	_memberRep   member.IMemberRep
	_partnerRep  merchant.IMerchantRep
	_deliveryRep delivery.IDeliveryRep
	_partnerId   int
	_partner     merchant.IMerchant
}

func NewShopping(partnerId int, partnerRep merchant.IMerchantRep,
	rep shopping.IShoppingRep, saleRep sale.ISaleRep, goodsRep sale.IGoodsRep,
	promRep promotion.IPromotionRep, memberRep member.IMemberRep,
	deliveryRep delivery.IDeliveryRep) shopping.IShopping {

	pt, _ := partnerRep.GetMerchant(partnerId)

	return &Shopping{
		_rep:         rep,
		_saleRep:     saleRep,
		_goodsRep:    goodsRep,
		_promRep:     promRep,
		_memberRep:   memberRep,
		_partnerId:   partnerId,
		_partnerRep:  partnerRep,
		_deliveryRep: deliveryRep,
		_partner:     pt,
	}
}

func (this *Shopping) GetAggregateRootId() int {
	return this._partnerId
}

func (this *Shopping) CreateOrder(val *shopping.ValueOrder, cart shopping.ICart) shopping.IOrder {
	return newOrder(this, val, cart, this._partnerRep, this._rep, this._saleRep, this._promRep, this._memberRep)
}

//创建购物车
// @buyerId 为购买会员ID,0表示匿名购物车
func (this *Shopping) NewCart(buyerId int) shopping.ICart {
	var cart shopping.ICart = newCart(this._partnerRep, this._memberRep, this._saleRep,
		this._goodsRep, this._rep, this._partnerId, buyerId)
	cart.Save()
	return cart
}

// 检查购物车
func (this *Shopping) CheckCart(cart shopping.ICart) error {
	if cart == nil || len(cart.GetValue().Items) == 0 {
		return shopping.ErrEmptyShoppingCart
	}

	sl := this._saleRep.GetSale(this._partnerId)
	for _, v := range cart.GetValue().Items {
		goods := sl.GetGoods(v.GoodsId)
		if goods == nil {
			return sale.ErrNoSuchGoods // 没有商品
		}
		stockNum := goods.GetValue().StockNum
		if stockNum == 0 {
			return sale.ErrFullOfStock // 已经卖完了
		}
		if stockNum < v.Quantity {
			return sale.ErrOutOfStock // 超出库存
		}
	}
	return nil
}

// 根据数据获取购物车
func (this *Shopping) GetCartByKey(key string) (shopping.ICart, error) {
	cart, error := this._rep.GetShoppingCart(key)
	if error == nil {
		return createCart(this._partnerRep, this._memberRep, this._saleRep,
			this._goodsRep, this._rep, this._partnerId, cart), nil
	}
	return nil, error
}

func (this *Shopping) GetShoppingCart(buyerId int, cartKey string) shopping.ICart {

	var hasOutCart = len(cartKey) != 0
	var hasBuyer = buyerId != 0

	var memCart shopping.ICart = nil // 消费者的购物车
	var outCart shopping.ICart = nil // 通过cartKey传入的购物车

	if hasBuyer {
		// 如果没有传递cartKey ，或者传递的cart和会员绑定的购物车相同，直接返回
		if memCart, _ = this.GetCurrentCart(buyerId); memCart != nil {
			if !hasOutCart || memCart.GetValue().CartKey == cartKey {
				return memCart
			}
		} else {
			memCart = this.NewCart(buyerId)
		}
	}

	if hasOutCart {
		outCart, _ = this.GetCartByKey(cartKey)
	}

	// 合并购物车
	if outCart != nil && hasBuyer {
		if bid := outCart.GetValue().BuyerId; bid <= 0 || bid == buyerId {
			memCart, _ = memCart.Combine(outCart)
			outCart.Destroy()
			memCart.Save()
		}
	}

	if memCart != nil {
		return memCart
	}

	if outCart != nil {
		return outCart
	}

	return this.NewCart(buyerId)

	//	if !hasOutCart {
	//		if c == nil {
	//			// 新的购物车不存在，直接返回会员的购物车
	//			if mc != nil {
	//				return mc
	//			}
	//		} else {
	//			cv := c.GetValue()
	//			//合并购物车
	//			if cv.BuyerId <= 0 {
	//				// 设置购买者
	//				if hasBuyer {
	//					c.SetBuyer(buyerId)
	//				}
	//			} else if mc != nil && cv.BuyerId == buyerId {
	//				// 合并购物车
	//				nc, err := mc.Combine(c)
	//				if err == nil {
	//					nc.Save()
	//					return nc
	//				}
	//				return mc
	//			}
	//
	//			// 如果没有购买，则返回
	//			return c
	//		}
	//	}

	// 返回一个新的购物车
	//	return this.NewCart(buyerId)
}

// 获取没有结算的购物车
func (this *Shopping) GetCurrentCart(buyerId int) (shopping.ICart, error) {
	cart, error := this._rep.GetLatestCart(buyerId)
	if error == nil {
		return createCart(this._partnerRep, this._memberRep, this._saleRep,
			this._goodsRep, this._rep, this._partnerId, cart), nil
	}
	return nil, error
}

// 绑定购物车会员编号
func (this *Shopping) BindCartBuyer(cartKey string, buyerId int) error {
	cart, err := this.GetCartByKey(cartKey)
	if err != nil {
		return err
	}
	return cart.SetBuyer(buyerId)
}

// 将购物车转换为订单
func (this *Shopping) ParseShoppingCart(memberId int) (shopping.IOrder,
	member.IMember, shopping.ICart, error) {
	var order shopping.IOrder
	var val shopping.ValueOrder
	var cart shopping.ICart
	var m member.IMember
	var err error

	m = this._memberRep.GetMember(memberId)
	if m == nil {
		return nil, m, nil, member.ErrSessionTimeout
	}

	cart, err = this.GetCurrentCart(memberId)

	if err != nil || cart == nil || len(cart.GetValue().Items) == 0 {
		return nil, m, cart, shopping.ErrEmptyShoppingCart
	}

	val.MemberId = memberId
	val.MerchantId = this._partnerId

	tf, of := cart.GetFee()
	val.TotalFee = tf //总金额
	val.Fee = of      //实际金额
	val.PayFee = of
	val.DiscountFee = tf - of //优惠金额
	val.MerchantId = this._partnerId
	val.Status = 1

	order = this.CreateOrder(&val, cart)
	return order, m, cart, nil
}

func (this *Shopping) GetFreeOrderNo() string {
	return this._rep.GetFreeOrderNo(this._partnerId)
}

// 智能选择门店
func (this *Shopping) SmartChoiceShop(address string) (merchant.IShop, error) {
	dly := this._deliveryRep.GetDelivery(this.GetAggregateRootId())

	lng, lat, err := lbs.GetLocation(address)
	if err != nil {
		return nil, errors.New("无法识别的地址：" + address)
	}
	var cov delivery.ICoverageArea = dly.GetNearestCoverage(lng, lat)
	if cov == nil {
		return nil, delivery.ErrNotCoveragedArea
	}
	shopId, _, err := dly.GetDeliveryInfo(cov.GetDomainId())
	return this._partner.GetShop(shopId), err
}

// 生成订单
func (this *Shopping) BuildOrder(memberId int, subject string, couponCode string) (shopping.IOrder, shopping.ICart, error) {
	order, m, cart, err := this.ParseShoppingCart(memberId)
	if err != nil {
		return order, cart, err
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
			this._partnerId, couponCode)

		// 如果优惠券不存在
		if cp == nil {
			log.Error(err)
			return order, cart, errors.New("优惠券无效")
		}

		coupon = cp.(promotion.ICouponPromotion)
		result, err = coupon.CanUse(m, val.Fee)
		if result {
			if coupon.CanTake() {
				_, err = coupon.GetTake(memberId)
				//如果未占用，则占用
				if err != nil {
					err = coupon.Take(memberId)
				}
			} else {
				_, err = coupon.GetBind(memberId)
			}
			if err != nil {
				log.Error(err)
				return order, cart, errors.New("优惠券无效")
			}
			err = order.ApplyCoupon(coupon) //应用优惠券
		}
	}

	return order, cart, err
}

func (this *Shopping) SubmitOrder(memberId int, subject string, couponCode string, useBalanceDiscount bool) (string, error) {
	order, cart, err := this.BuildOrder(memberId, subject, couponCode)
	if err != nil {
		return "", err
	}
	var cv = cart.GetValue()
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

func (this *Shopping) GetOrderByNo(orderNo string) (shopping.IOrder, error) {
	val, err := this._rep.GetOrderByNo(this._partnerId, orderNo)
	if err != nil {
		return nil, errors.New("订单不存在")
	}

	val.Items = this._rep.GetOrderItems(val.Id)
	return this.CreateOrder(val, nil), err
}

var (
	shopLocker sync.Mutex
	biShops    []merchant.IShop
)

// 自动设置订单
func (this *Shopping) OrderAutoSetup(f func(error)) {
	var orders []*shopping.ValueOrder
	var err error

	shopLocker.Lock()
	defer func() {
		shopLocker.Unlock()
	}()
	biShops = nil
	log.Println("[SETUP] start auto setup")

	saleConf := this._partner.GetSaleConf()
	if saleConf.AutoSetupOrder == 1 {
		orders, err = this._rep.GetWaitingSetupOrders(this._partnerId)
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

func (this *Shopping) SmartConfirmOrder(order shopping.IOrder) error {
	var err error
	v := order.GetValue()
	log.Printf("[ AUTO][OrderSetup]:%s - Confirm \n", v.OrderNo)
	var shop merchant.IShop
	if biShops == nil {
		biShops = this._partner.GetBusinessInShops()
	}
	if len(biShops) == 1 {
		shop = biShops[0]
	} else {
		shop, err = this.SmartChoiceShop(v.DeliverAddress)
		if err != nil {
			order.Suspend("智能分配门店失败！原因：" + err.Error())
			return err
		}
	}

	if shop != nil {
		sv := shop.GetValue()
		order.SetShop(shop.GetDomainId())
		err = order.Confirm()
		//err = order.Process()
		order.AppendLog(enum.ORDER_LOG_SETUP, false, fmt.Sprintf(
			"自动分配门店:%s,电话：%s", sv.Name, sv.Phone))
	}
	return err
}

func (this *Shopping) setupOrder(v *shopping.ValueOrder,
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
