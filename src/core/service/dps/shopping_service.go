/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:53
 * description :
 * history :
 */

package dps

import (
	"bytes"
	"errors"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/dto"
)

type shoppingService struct {
	_rep shopping.IShoppingRep
}

func NewShoppingService(r shopping.IShoppingRep) *shoppingService {
	return &shoppingService{
		_rep: r,
	}
}

func (this *shoppingService) BuildOrder(partnerId int, memberId int,
	cartKey string, couponCode string) (map[string]interface{}, error) {
	var sp shopping.IShopping = this._rep.GetShopping(partnerId)
	order, _, err := sp.BuildOrder(memberId, couponCode)
	if err != nil {
		return nil, err
	}

	v := order.GetValue()
	buf := bytes.NewBufferString("")

	for _, v := range order.GetCoupons() {
		buf.WriteString(v.GetDescribe())
		buf.WriteString("\n")
	}

	var data map[string]interface{}
	data = make(map[string]interface{})
	if couponCode != "" {
		if v.CouponFee == 0 {
			data["result"] = v.CouponFee != 0
			data["message"] = "优惠券无效"
		} else {
			// 成功应用优惠券
			data["totalFee"] = v.TotalFee
			data["fee"] = v.Fee
			data["payFee"] = v.PayFee
			data["discountFee"] = v.DiscountFee
			data["couponFee"] = v.CouponFee
			data["couponDescribe"] = buf.String()
		}
	} else {
		//　取消优惠券
		data["totalFee"] = v.TotalFee
		data["fee"] = v.Fee
		data["payFee"] = v.PayFee
		data["discountFee"] = v.DiscountFee
	}
	return data, err
}

func (this *shoppingService) SubmitOrder(partnerId, memberId int, couponCode string, useBalanceDiscount bool) (
	orderNo string, err error) {
	var sp shopping.IShopping = this._rep.GetShopping(partnerId)
	return sp.SubmitOrder(memberId, couponCode, useBalanceDiscount)
}

func (this *shoppingService) SetDeliverShop(partnerId int, orderNo string,
	shopId int) error {
	var sp shopping.IShopping = this._rep.GetShopping(partnerId)
	order, err := sp.GetOrderByNo(orderNo)
	if err == nil {
		err = order.SetShop(shopId)
	}
	err = order.SetShop(shopId)
	if err == nil {
		_, err = order.Save()
	}
	return err
}

func (this *shoppingService) HandleOrder(partnerId int, orderNo string) error {
	var sp shopping.IShopping = this._rep.GetShopping(partnerId)
	order, err := sp.GetOrderByNo(orderNo)
	if err == nil {
		b := order.IsOver()
		if b {
			return errors.New("订单已经完成!")
		}

		status := order.GetValue().Status
		switch status + 1 {
		case enum.ORDER_WAIT_CONFIRM:
			err = order.Confirm()
		case enum.ORDER_WAIT_DELIVERY:
			err = order.Process()
		case enum.ORDER_WAIT_RECEIVE:
			err = order.Deliver()
		case enum.ORDER_RECEIVED:
			err = order.SignReceived()
		case enum.ORDER_COMPLETED:
			err = order.Complete()
		}
	}
	return err
}

func (this *shoppingService) GetOrderByNo(partnerId int,
	orderNo string) *shopping.ValueOrder {
	var sp shopping.IShopping = this._rep.GetShopping(partnerId)
	order, err := sp.GetOrderByNo(orderNo)
	if err != nil {
		Context.Log().PrintErr(err)
		return nil
	}
	v := order.GetValue()
	return &v
}

func (this *shoppingService) CancelOrder(partnerId int, orderNo string, reason string) error {
	var sp shopping.IShopping = this._rep.GetShopping(partnerId)
	order, err := sp.GetOrderByNo(orderNo)
	if err == nil {
		err = order.Cancel(reason)
	}
	return err
}

//  获取购物车
func (this *shoppingService) getShoppingCart(partnerId int, buyerId int, cartKey string) shopping.ICart {
	sp := this._rep.GetShopping(partnerId)
	return sp.GetShoppingCart(buyerId, cartKey)
}
func (this *shoppingService) GetShoppingCart(partnerId int, memberId int, cartKey string) *dto.ShoppingCart {
	cart := this.getShoppingCart(partnerId, memberId, cartKey)
	return this.parseDtoCart(cart)
}

func (this *shoppingService) parseDtoCart(c shopping.ICart) *dto.ShoppingCart {
	var cart = new(dto.ShoppingCart)
	v := c.GetValue()
	cart.Id = c.GetDomainId()
	cart.BuyerId = v.BuyerId
	cart.CartKey = v.CartKey
	cart.UpdateTime = v.UpdateTime
	t, f := c.GetFee()
	cart.TotalFee = t
	cart.OrderFee = f
	cart.Summary = c.GetSummary()

	if v.Items != nil {
		if l := len(v.Items); l != 0 {
			cart.Items = make([]*dto.CartItem, l)
			for i, v := range v.Items {
				cart.Items[i] = &dto.CartItem{
					GoodsId:    v.GoodsId,
					GoodsName:  v.Name,
					GoodsNo:    v.GoodsNo,
					SmallTitle: v.SmallTitle,
					GoodsImage: v.Image,
					Num:        v.Quantity,
					Price:      v.Price,
					SalePrice:  v.SalePrice,
				}
				cart.TotalNum += cart.Items[i].Num
			}
		}
	}

	return cart
}

func (this *shoppingService) AddCartItem(partnerId, memberId int, cartKey string, goodsId, num int) (*dto.CartItem, error) {
	cart := this.getShoppingCart(partnerId, memberId, cartKey)
	item, err := cart.AddItem(goodsId, num)
	if err == nil {
		cart.Save()
		return &dto.CartItem{
			GoodsId:    item.GoodsId,
			GoodsName:  item.Name,
			SmallTitle: item.SmallTitle,
			GoodsImage: item.Image,
			Num:        num,
			Price:      item.Price,
			SalePrice:  item.SalePrice,
		}, nil
	}
	return nil, err
}
func (this *shoppingService) SubCartItem(partnerId, memberId int, cartKey string, goodsId, num int) error {
	cart := this.getShoppingCart(partnerId, memberId, cartKey)
	err := cart.RemoveItem(goodsId, num)
	if err == nil {
		_, err = cart.Save()
	}
	return err
}

// 更新购物车结算
func (this *shoppingService) PrepareSettlePersist(partnerId, memberId, shopId, paymentOpt, deliverOpt, deliverId int) error {
	var cart = this.getShoppingCart(partnerId, memberId, "")
	err := cart.SettlePersist(shopId, paymentOpt, deliverOpt, deliverId)
	if err == nil {
		_, err = cart.Save()
	}
	return err
}

func (this *shoppingService) GetCartSettle(partnerId, memberId int, cartKey string) *dto.SettleMeta {
	var cart = this.getShoppingCart(partnerId, memberId, cartKey)
	shop, deliver, payOpt, dlvOpt := cart.GetSettleData()
	var st *dto.SettleMeta = new(dto.SettleMeta)
	st.PaymentOpt = payOpt
	st.DeliverOpt = dlvOpt
	if shop != nil {
		v := shop.GetValue()

		st.Shop = &dto.SettleShopMeta{
			Id:   v.Id,
			Name: v.Name,
			Tel:  v.Phone,
		}
	}

	if deliver != nil {
		v := deliver.GetValue()
		st.Deliver = &dto.SettleDeliverMeta{
			Id:         v.Id,
			PersonName: v.RealName,
			Phone:      v.Phone,
			Address:    v.Address,
		}
	}

	return st
}

func (this *shoppingService) OrderAutoSetup(partnerId int, f func(error)) {
	sp := this._rep.GetShopping(partnerId)
	sp.OrderAutoSetup(f)
}

// 使用余额为订单付款
func (this *shoppingService) PayForOrderWithBalance(partnerId int, orderNo string) error {
	var sp shopping.IShopping = this._rep.GetShopping(partnerId)
	order, err := sp.GetOrderByNo(orderNo)
	if err == nil {
		err = order.PaymentWithBalance()
	}
	return err
}

// 确认付款
func (this *shoppingService) PayForOrderOnlineTrade(partnerId int, orderNo string, spName string, tradeNo string) error {
	var sp shopping.IShopping = this._rep.GetShopping(partnerId)
	order, err := sp.GetOrderByNo(orderNo)
	if err == nil {
		err = order.PaymentOnlineTrade(spName, tradeNo)
	}
	return err
}

// 确定订单
func (this *shoppingService) ConfirmOrder(partnerId int, orderNo string) error {
	var sp shopping.IShopping = this._rep.GetShopping(partnerId)
	order, err := sp.GetOrderByNo(orderNo)
	if err == nil {
		err = order.Confirm()
	}
	return err
}
