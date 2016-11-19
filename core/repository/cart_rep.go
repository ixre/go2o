/**
 * Copyright 2015 @ z3q.net.
 * name : cart_rep.go
 * author : jarryliu
 * date : 2016-06-29 22:54
 * description :
 * history :
 */
package repository

import (
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	cartImpl "go2o/core/domain/cart"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/sale/goods"
)

var _ cart.ICartRep = new(cartRep)

type cartRep struct {
	db.Connector
	_goodsRep  goods.IGoodsRep
	_memberRep member.IMemberRep
}

func NewCartRep(conn db.Connector, memberRep member.IMemberRep,
	goodsRep goods.IGoodsRep) cart.ICartRep {
	return &cartRep{
		Connector:  conn,
		_memberRep: memberRep,
		_goodsRep:  goodsRep,
	}
}

// 创建购物车对象
func (c *cartRep) CreateCart(v *cart.ValueCart) cart.ICart {
	return cartImpl.CreateCart(v, c, c._memberRep, c._goodsRep)
}

// 创建一个购物车
func (c *cartRep) NewCart() cart.ICart {
	return cartImpl.NewCart(-1, c, c._memberRep, c._goodsRep)
}

// 获取购物车
func (c *cartRep) GetShoppingCartByKey(key string) cart.ICart {
	ca := c.GetShoppingCart(key)
	if ca != nil {
		return c.CreateCart(ca)
	}
	return nil
}

// 获取会员没有结算的购物车
func (c *cartRep) GetMemberCurrentCart(buyerId int32) cart.ICart {
	ca := c.GetLatestCart(buyerId)
	if ca != nil {
		return c.CreateCart(ca)
	}
	return nil
}

// 获取购物车
func (c *cartRep) GetShoppingCart(key string) *cart.ValueCart {
	var v = &cart.ValueCart{}
	if c.Connector.GetOrm().GetBy(v, "cart_key=?", key) == nil {
		items := []*cart.CartItem{}
		c.Connector.GetOrm().Select(&items, "cart_id=?", v.Id)
		v.Items = items
		return v
	}
	return nil
}

// 获取最新的购物车
func (c *cartRep) GetLatestCart(buyerId int32) *cart.ValueCart {
	var v = &cart.ValueCart{}
	if c.Connector.GetOrm().GetBy(v, "buyer_id=? ORDER BY id DESC", buyerId) == nil {
		var items = []*cart.CartItem{}
		c.Connector.GetOrm().Select(&items, "cart_id=?", v.Id)
		v.Items = items
		return v
	}
	return nil
}

// 保存购物车
func (c *cartRep) SaveShoppingCart(v *cart.ValueCart) (int32, error) {
	return orm.I32(orm.Save(c.GetOrm(), v, int(v.Id)))
}

// 移出购物车项
func (c *cartRep) RemoveCartItem(id int32) error {
	return c.Connector.GetOrm().DeleteByPk(cart.CartItem{}, id)
}

// 保存购物车项
func (c *cartRep) SaveCartItem(v *cart.CartItem) (int32, error) {
	return orm.I32(orm.Save(c.GetOrm(), v, int(v.Id)))
}

// 清空购物车项
func (c *cartRep) EmptyCartItems(cartId int32) error {
	_, err := c.Connector.GetOrm().Delete(cart.CartItem{}, "cart_id=?", cartId)
	return err
}

// 删除购物车
func (c *cartRep) DeleteCart(cartId int32) error {
	return c.Connector.GetOrm().DeleteByPk(cart.ValueCart{}, cartId)
}
