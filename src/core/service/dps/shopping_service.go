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
	cartKey string, couponCode string) (shopping.IOrder, shopping.ICart, error) {
	var sp shopping.IShopping = this._rep.GetShopping(partnerId)
	return sp.BuildOrder(memberId, couponCode)
}

func (this *shoppingService) SubmitOrder(partnerId, memberId int, couponCode string) (
	orderNo string, err error) {
	var sp shopping.IShopping = this._rep.GetShopping(partnerId)
	return sp.SubmitOrder(memberId, couponCode)
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
		case enum.ORDER_CONFIRMED:
			err = order.Confirm()
		case enum.ORDER_PROCESSING:
			err = order.Process()
		case enum.ORDER_SENDING:
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
func (this *shoppingService) getShoppingCart(partnerId int, memberId int, cartKey string) shopping.ICart {
	sp := this._rep.GetShopping(partnerId)
	var c shopping.ICart
	var mc shopping.ICart

	var skIsNil = len(cartKey) == 0
	var mmNotNil = memberId != 0

	if mmNotNil {
		mc, _ = sp.GetNotBoughtCart(memberId)
		if mc != nil && (skIsNil || mc.GetValue().CartKey == cartKey) {
			return mc
		}
	}

	if !skIsNil {
		// 根据Key获取购物车
		c, _ = sp.GetCart(cartKey)
		if c == nil {
			// 新的购物车不存在，直接返回会员的购物车
			if mc != nil {
				return mc
			}
		} else {

			cv := c.GetValue()
			//合并购物车
			if cv.BuyerId <= 0 {
				// 设置购买者
				if mmNotNil {
					c.SetBuyer(memberId)
				}
			} else if mc != nil {
				// 合并购物车
				nc, err := mc.Combine(c)
				if err == nil {
					nc.Save()
					return nc
				}
				return mc
			}

			// 如果没有购买，则返回
			if cv.IsBought == 0 {
				return c
			}
		}
	}

	// 返回一个新的购物车
	return sp.NewCart(memberId)
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
	cart.IsBought = v.IsBought

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
					Num:        v.Num,
					Price:      v.Price,
					SalePrice:  v.SalePrice,
				}
				cart.TotalNum += cart.Items[i].Num
			}
		}
	}

	return cart
}

func (this *shoppingService) AddCartItem(partnerId, memberId int, cartKey string, goodsId, num int) *dto.CartItem {
	cart := this.getShoppingCart(partnerId, memberId, cartKey)
	item := cart.AddItem(goodsId, num)
	if item != nil {
		cart.Save()
		return &dto.CartItem{
			GoodsId:    item.GoodsId,
			GoodsName:  item.Name,
			SmallTitle: item.SmallTitle,
			GoodsImage: item.Image,
			Num:        num,
			Price:      item.Price,
			SalePrice:  item.SalePrice,
		}
	}
	return nil
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

func (this *shoppingService) PayForOrder(partnerId int, orderNo string) error {
	var sp shopping.IShopping = this._rep.GetShopping(partnerId)
	order, err := sp.GetOrderByNo(orderNo)
	if err == nil {
		err = order.SignPaid()
	}
	return err
}
