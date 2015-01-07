/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-05 17:49
 * description :
 * history :
 */

package shopping

import (
	"com/domain/interface/member"
	"com/domain/interface/partner"
	"com/domain/interface/promotion"
	"com/domain/interface/sale"
	"com/domain/interface/shopping"
	"com/infrastructure/log"
	"errors"
	"regexp"
	"strconv"
)

var (
	//data example : 16*1|12*2|80
	cartRegex = regexp.MustCompile("(\\d+)\\s*\\*\\s*(\\d+)")
)

type Shopping struct {
	_rep        shopping.IShoppingRep
	_saleRep    sale.ISaleRep
	_promRep    promotion.IPromotionRep
	_memberRep  member.IMemberRep
	_partnerRep partner.IPartnerRep
	_partnerId  int
}

func NewShopping(partnerId int, partnerRep partner.IPartnerRep,
	rep shopping.IShoppingRep, saleRep sale.ISaleRep,
	promRep promotion.IPromotionRep,
	memberRep member.IMemberRep) shopping.IShopping {
	return &Shopping{
		_rep:        rep,
		_saleRep:    saleRep,
		_promRep:    promRep,
		_memberRep:  memberRep,
		_partnerId:  partnerId,
		_partnerRep: partnerRep,
	}
}

func (this *Shopping) GetAggregateRootId() int {
	return this._partnerId
}

func (this *Shopping) CreateOrder(val *shopping.ValueOrder, cart shopping.ICart) shopping.IOrder {
	return newOrder(this, val, cart, this._partnerRep, this._rep, this._memberRep)
}

func (this *Shopping) CreateCart(value *shopping.ValueCart) shopping.ICart {
	return newCart(value)
}

// 根据数据获取购物车
func (this *Shopping) GetCart(s string) (shopping.ICart, error) {
	var cart *shopping.ValueCart
	var err error

	if !cartRegex.MatchString(s) {
		log.PrintErr(errors.New("Error Cart:" + s))
		return nil, errors.New("Code 103:购物车异常")
	}
	cart = new(shopping.ValueCart)
	matches := cartRegex.FindAllStringSubmatch(s, -1)

	length := len(matches)
	var ids []int = make([]int, length) //ID数组
	cart.Quantities = make(map[int]int, length)

	var id int
	var qua int

	for i, v := range matches {
		id, err = strconv.Atoi(v[1])
		if err != nil {
			continue
		}
		qua, err = strconv.Atoi(v[2])
		if err != nil {
			continue
		}

		ids[i] = id
		cart.Quantities[id] = qua
	}

	cart.Items, err = this._saleRep.GetProductByIds(this._partnerId, ids...)
	if err != nil {
		log.PrintErr(err)
		err = errors.New("Code 102:商品异常")
		return nil, err
	}

	return this.CreateCart(cart), nil
}

// 将购物车转换为订单
func (this *Shopping) GetOrderByCart(cartStr string) (*shopping.ValueOrder, error) {
	var orderVal shopping.ValueOrder
	cart, err := this.GetCart(cartStr)
	if err != nil {
		return nil, err
	}
	tf, of := cart.GetFee()
	orderVal.TotalFee = tf //总金额
	orderVal.Fee = of      //实际金额
	orderVal.PayFee = of
	orderVal.DiscountFee = tf - of //优惠金额
	return &orderVal, err
}

func (this *Shopping) GetFreeOrderNo() string {
	return this._rep.GetFreeOrderNo(this._partnerId)
}

func (this *Shopping) BuildOrder(memberId int, cartStr string,
	couponCode string) (shopping.IOrder, error) {
	var order shopping.IOrder
	cart, err := this.GetCart(cartStr)
	if err != nil {
		return nil, err
	}
	val, err := this.GetOrderByCart(cartStr)
	if err != nil {
		return nil, err
	}

	m, err := this._memberRep.GetMember(memberId)
	if err != nil {
		return order, errors.New("Code 101:登录超时")
	}

	val.MemberId = memberId
	val.PartnerId = this._partnerId
	val.Items = cartStr
	val.Status = 1
	order = this.CreateOrder(val, cart)

	if len(couponCode) != 0 {
		var coupon promotion.ICoupon
		var result bool
		coupon, err = this._promRep.GetCouponByCode(
			this._partnerId, couponCode)

		// 如果优惠券不存在
		if err != nil || coupon == nil {
			log.PrintErr(err)
			return order, errors.New("优惠券无效")
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
				return order, errors.New("优惠券无效")
			}
			order.ApplyCoupon(coupon) //应用优惠券
		}
	}

	return order, err
}

func (this *Shopping) SubmitOrder(memberId, shopId int, payMethod int,
	deliverAddrId int, cart string, couponCode string, note string) (string, error) {
	order, err := this.BuildOrder(memberId, cart, couponCode)
	if err != nil {
		order, err = this.BuildOrder(memberId, cart, "")
	}

	if err == nil {
		order.AddRemark(note)
		err = order.SetShop(shopId)
		if err == nil {
			order.SetPayment(payMethod)
			err = order.SetDeliver(deliverAddrId)
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
		log.PrintErr(err)
		return nil, errors.New("订单不存在")
	}
	return this.CreateOrder(val, nil), err
}
