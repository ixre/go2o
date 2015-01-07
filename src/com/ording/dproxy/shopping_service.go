/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-05 17:53
 * description :
 * history :
 */

package dproxy

import (
	"com/domain/interface/enum"
	"com/domain/interface/shopping"
	"errors"
)

type shoppingService struct {
	spRep shopping.IShoppingRep
}

func (this *shoppingService) BuildOrder(partnerId int, memberId int,
	cart string, couponCode string) (shopping.IOrder, error) {
	var sp shopping.IShopping = this.spRep.GetShopping(partnerId)
	return sp.BuildOrder(memberId, cart, couponCode)
}

func (this *shoppingService) SubmitOrder(partnerId, memberId, shopId int, payMethod int,
	deliverAddrId int, cart string, couponCode string, note string) (
	orderNo string, err error) {
	var sp shopping.IShopping = this.spRep.GetShopping(partnerId)
	return sp.SubmitOrder(memberId, shopId, payMethod,
		deliverAddrId, cart, couponCode, note)
}

func (this *shoppingService) SetDeliverShop(partnerId int, orderNo string,
	shopId int) error {
	var sp shopping.IShopping = this.spRep.GetShopping(partnerId)
	order, err := sp.GetOrderByNo(orderNo)
	if err == nil {
		err = order.SetShop(shopId)
	}
	err = order.SetShop(shopId)
	if err == nil {
		err = order.Save()
	}
	return err
}

func (this *shoppingService) HandleOrder(partnerId int, orderNo string) error {
	var sp shopping.IShopping = this.spRep.GetShopping(partnerId)
	order, err := sp.GetOrderByNo(orderNo)
	if err == nil {
		b := order.IsOver()
		if b {
			return errors.New("订单已经完成!")
		}

		status := order.GetValue().Status
		switch status + 1 {
		case enum.ORDER_CONFIRMED:
			err = order.Confirm()
		case enum.ORDER_PROCESSING:
			err = order.Process()
		case enum.ORDER_SENDING:
			err = order.Deliver()
		case enum.ORDER_COMPLETED:
			err = order.Complete()
		}
	}
	return err
}

func (this *shoppingService) GetOrderByNo(partnerId int,
	orderNo string) *shopping.ValueOrder {
	var sp shopping.IShopping = this.spRep.GetShopping(partnerId)
	order, err := sp.GetOrderByNo(orderNo)
	if err != nil {
		Context.Log().PrintErr(err)
		return nil
	}
	v := order.GetValue()
	return &v
}

func (this *shoppingService) CancelOrder(partnerId int, orderNo string, reason string) error {
	var sp shopping.IShopping = this.spRep.GetShopping(partnerId)
	order, err := sp.GetOrderByNo(orderNo)
	if err == nil {
		err = order.Cancel(reason)
	}
	return err
}

func (this *shoppingService) ParseShoppingCart(partnerId int, cartData string) (shopping.ICart, error) {
	sp := this.spRep.GetShopping(partnerId)
	return sp.GetCart(cartData)
}
