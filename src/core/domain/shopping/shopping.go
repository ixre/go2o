/**
 * Copyright 2014 @ ops Inc.
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
	"github.com/atnet/gof"
	"go2o/src/core/domain/interface/delivery"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/infrastructure"
	"go2o/src/core/infrastructure/lbs"
	"go2o/src/core/infrastructure/log"
	"regexp"
	"time"
)

var (
	//data example : 16*1|12*2|80
	cartRegex = regexp.MustCompile("(\\d+)\\s*\\*\\s*(\\d+)")
)

type Shopping struct {
	_rep         shopping.IShoppingRep
	_saleRep     sale.ISaleRep
	_promRep     promotion.IPromotionRep
	_memberRep   member.IMemberRep
	_partnerRep  partner.IPartnerRep
	_deliveryRep delivery.IDeliveryRep
	_partnerId   int
	_partner     partner.IPartner
}

func NewShopping(partnerId int, partnerRep partner.IPartnerRep,
	rep shopping.IShoppingRep, saleRep sale.ISaleRep,
	promRep promotion.IPromotionRep, memberRep member.IMemberRep,
	deliveryRep delivery.IDeliveryRep) shopping.IShopping {

	pt, _ := partnerRep.GetPartner(partnerId)

	return &Shopping{
		_rep:         rep,
		_saleRep:     saleRep,
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
	return newOrder(this, val, cart, this._partnerRep, this._rep, this._memberRep)
}

//创建购物车
// @buyerId 为购买会员ID,0表示匿名购物车
func (this *Shopping) NewCart(buyerId int) shopping.ICart {
	var cart shopping.ICart = newCart(this._partnerRep, this._memberRep, this._saleRep,
		this._rep, this._partnerId, buyerId)
	cart.Save()
	return cart
}

// 根据数据获取购物车
func (this *Shopping) GetCart(key string) (shopping.ICart, error) {
	cart, error := this._rep.GetShoppingCart(key)
	if error == nil {
		return createCart(this._partnerRep, this._memberRep, this._saleRep,
			this._rep, this._partnerId, cart), nil
	}
	return nil, error
}

// 获取没有结算的购物车
func (this *Shopping) GetNotBoughtCart(buyerId int) (shopping.ICart, error) {
	cart, error := this._rep.GetNotBoughtCart(buyerId)
	if error == nil {
		return createCart(this._partnerRep, this._memberRep, this._saleRep,
			this._rep, this._partnerId, cart), nil
	}
	return nil, error
}

// 绑定购物车会员编号
func (this *Shopping) BindCartBuyer(cartKey string, buyerId int) error {
	cart, err := this.GetCart(cartKey)
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

	m, err = this._memberRep.GetMember(memberId)
	if m == nil {
		return nil, m, nil, member.ErrSessionTimeout
	}

	cart, err = this.GetNotBoughtCart(memberId)
	if err != nil || cart == nil || len(cart.GetValue().Items) == 0 {
		return nil, m, cart, shopping.ErrEmptyShoppingCart
	}

	val.MemberId = memberId
	val.PartnerId = this._partnerId

	tf, of := cart.GetFee()
	val.TotalFee = tf //总金额
	val.Fee = of      //实际金额
	val.PayFee = of
	val.DiscountFee = tf - of //优惠金额
	val.PartnerId = this._partnerId
	val.Status = 1

	order = this.CreateOrder(&val, cart)
	return order, m, cart, nil
}

func (this *Shopping) GetFreeOrderNo() string {
	return this._rep.GetFreeOrderNo(this._partnerId)
}

// 智能选择门店
func (this *Shopping) SmartChoiceShop(address string) (partner.IShop, error) {
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
func (this *Shopping) BuildOrder(memberId int, couponCode string) (shopping.IOrder, shopping.ICart, error) {
	order, m, cart, err := this.ParseShoppingCart(memberId)
	if err != nil {
		return order, cart, err
	}

	if len(couponCode) != 0 {
		var coupon promotion.ICoupon
		var result bool
		var val = order.GetValue()
		coupon, err = this._promRep.GetCouponByCode(
			this._partnerId, couponCode)

		// 如果优惠券不存在
		if err != nil || coupon == nil {
			log.PrintErr(err)
			return order, cart, errors.New("优惠券无效")
		}

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
				log.PrintErr(err)
				return order, cart, errors.New("优惠券无效")
			}
			order.ApplyCoupon(coupon) //应用优惠券
		}
	}

	return order, cart, err
}

func (this *Shopping) SubmitOrder(memberId int, couponCode string) (string, error) {
	order, cart, err := this.BuildOrder(memberId, couponCode)
	if err != nil {
		order, cart, err = this.BuildOrder(memberId, "")
	}
	var cv = cart.GetValue()

	if err == nil {
		err = order.SetShop(cv.ShopId)
		if err == nil {
			order.SetPayment(cv.PaymentOpt)
			err = order.SetDeliver(cv.DeliverId)
			if err == nil {
				var orderNo string
				orderNo, err = order.Submit()
				if err == nil {
					err = cart.BindOrder(orderNo)
				}
				return orderNo, err
			}
		}
	}

	return "", err
}

func (this *Shopping) GetOrderByNo(orderNo string) (shopping.IOrder, error) {
	val, err := this._rep.GetOrderByNo(this._partnerId, orderNo)

	if err != nil {
		log.PrintErr(err)
		return nil, errors.New("订单不存在")
	}
	return this.CreateOrder(val, nil), err
}

// 自动设置订单
func (this *Shopping) OrderAutoSetup(f func(error)) {
	var orders []*shopping.ValueOrder
	var err error

	log.Println("[SETUP] start auto setup")

	saleConf := this._partner.GetSaleConf()
	if saleConf.AutoSetupOrder == 1 {
		ctx := infrastructure.GetApp()
		orders, err = this._rep.GetWaitingSetupOrders(this._partnerId)
		if err != nil {
			f(err)
			return
		}

		dt := time.Now()
		for _, v := range orders {
			this.setupOrder(ctx, v, &saleConf, dt, f)
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

func (this *Shopping) setupOrder(ctx gof.App, v *shopping.ValueOrder,
	conf *partner.SaleConf, t time.Time, f func(error)) {
	var err error
	order := this.CreateOrder(v, nil)
	dur := time.Duration(t.Unix()-v.CreateTime) * time.Second

	if v.PaymentOpt == enum.PAY_ONLINE {
		if v.IsPaid == 0 && dur > time.Hour*order_timeout_hour {
			order.Cancel("超时未付款，系统取消")
			if ctx.Debug() {
				ctx.Log().Printf("[ AUTO][OrderSetup]:%s - Payment Timeout\n", v.OrderNo)
			}
		}
	} else if v.PaymentOpt == enum.PAY_OFFLINE {
		switch v.Status + 1 {
		case enum.ORDER_CONFIRMED:
			if dur > time.Minute*order_confirm_minute {
				err = order.Confirm()
				if ctx.Debug() {
					ctx.Log().Printf("[ AUTO][OrderSetup]:%s - Confirm \n", v.OrderNo)
				}

				shop, err := this.SmartChoiceShop(v.DeliverAddress)
				if err != nil {
					log.Println(err)
					order.Suspend("智能分配门店失败！原因：" + err.Error())
				} else {
					sv := shop.GetValue()
					order.SetShop(shop.GetDomainId())
					order.AppendLog(enum.ORDER_LOG_SETUP, false, fmt.Sprintf(
						"自动分配门店:%s,电话：%s", sv.Name, sv.Phone))
				}
			}
		case enum.ORDER_PROCESSING:
			if dur > time.Minute*order_process_minute {
				err = order.Process()
				if ctx.Debug() {
					ctx.Log().Printf("[ AUTO][OrderSetup]:%s - Processing \n", v.OrderNo)
				}
			}

		case enum.ORDER_SENDING:
			if dur > time.Minute*order_sending_minute {
				err = order.Deliver()
				if ctx.Debug() {
					ctx.Log().Printf("[ AUTO][OrderSetup]:%s - Sending \n", v.OrderNo)
				}
			}
		case enum.ORDER_RECEIVED:
			if dur > time.Hour*order_receive_hour {
				err = order.SignReceived()
				if ctx.Debug() {
					ctx.Log().Printf("[ AUTO][OrderSetup]:%s - Received \n", v.OrderNo)
				}
			}
		case enum.ORDER_COMPLETED:
			if dur > time.Hour*order_complete_hour {
				err = order.Complete()
				if ctx.Debug() {
					ctx.Log().Printf("[ AUTO][OrderSetup]:%s - Complete \n", v.OrderNo)
				}
			}
		}
	}

	if err != nil {
		f(err)
	}
}
