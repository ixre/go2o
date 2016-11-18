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
	c := c.GetShoppingCart(key)
	if c != nil {
		return c.CreateCart(c)
	}
	return nil
}

// 获取会员没有结算的购物车
func (c *cartRep) GetMemberCurrentCart(buyerId int64) cart.ICart {
	c := c.GetLatestCart(buyerId)
	if c != nil {
		return c.CreateCart(c)
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
func (c *cartRep) GetLatestCart(buyerId int64) *cart.ValueCart {
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
func (c *cartRep) SaveShoppingCart(v *cart.ValueCart) (int64, error) {
	var err error
	_orm := c.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = _orm.Save(v.Id, v)
	} else {
		_, _, err = _orm.Save(nil, v)
		c.Connector.ExecScalar(`SELECT MAX(id) FROM sale_cart`, &v.Id)
	}
	return v.Id, err
}

// 移出购物车项
func (c *cartRep) RemoveCartItem(id int64) error {
	return c.Connector.GetOrm().DeleteByPk(cart.CartItem{}, id)
}

// 保存购物车项
func (c *cartRep) SaveCartItem(v *cart.CartItem) (int64, error) {
	_orm := c.Connector.GetOrm()
	var err error
	if v.Id > 0 {
		_, _, err = _orm.Save(v.Id, v)
	} else {
		_, _, err = _orm.Save(nil, v)
		c.Connector.ExecScalar(`SELECT MAX(id) FROM sale_cart_item where cart_id=?`, &v.Id, v.CartId)
	}

	return v.Id, err
}

// 清空购物车项
func (c *cartRep) EmptyCartItems(cartId int64) error {
	_, err := c.Connector.GetOrm().Delete(cart.CartItem{}, "cart_id=?", cartId)
	return err
}

// 删除购物车
func (c *cartRep) DeleteCart(cartId int64) error {
	return c.Connector.GetOrm().DeleteByPk(cart.ValueCart{}, cartId)
}
