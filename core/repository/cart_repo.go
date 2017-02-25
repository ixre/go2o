/**
 * Copyright 2015 @ z3q.net.
 * name : cart_repo.go
 * author : jarryliu
 * date : 2016-06-29 22:54
 * description :
 * history :
 */
package repository

import (
	"database/sql"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	cartImpl "go2o/core/domain/cart"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"log"
)

var _ cart.ICartRepo = new(cartRepo)

type cartRepo struct {
	db.Connector
	_orm        orm.Orm
	_goodsRepo  item.IGoodsItemRepo
	_memberRepo member.IMemberRepo
}

func NewCartRepo(conn db.Connector, memberRepo member.IMemberRepo,
	goodsRepo item.IGoodsItemRepo) cart.ICartRepo {
	return &cartRepo{
		Connector:   conn,
		_orm:        conn.GetOrm(),
		_memberRepo: memberRepo,
		_goodsRepo:  goodsRepo,
	}
}

// 创建购物车对象
func (c *cartRepo) CreateCart(v *cart.ValueCart) cart.ICart {
	return cartImpl.CreateCart(v, c, c._memberRepo, c._goodsRepo)
}

// 创建一个购物车
func (c *cartRepo) NewCart() cart.ICart {
	return cartImpl.NewCart(-1, c, c._memberRepo, c._goodsRepo)
}

// 获取购物车
func (c *cartRepo) GetCart(id int32) cart.ICart {
	v := c.GetSaleCart(id)
	if v != nil {
		v.Items = c.SelectSaleCartItem("cart_id=?", id)
		return c.CreateCart(v)
	}
	return nil
}

// Get SaleCart
func (s *cartRepo) GetSaleCart(primary interface{}) *cart.ValueCart {
	e := cart.ValueCart{}
	err := s._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCart")
	}
	return nil
}

// Select SaleCart
func (s *cartRepo) SelectSaleCart(where string, v ...interface{}) []*cart.ValueCart {
	list := []*cart.ValueCart{}
	err := s._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCart")
	}
	return list
}

// Save SaleCart
func (s *cartRepo) SaveSaleCart(v *cart.ValueCart) (int, error) {
	id, err := orm.Save(s._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCart")
	}
	return id, err
}

// Delete SaleCart
func (s *cartRepo) DeleteSaleCart(primary interface{}) error {
	err := s._orm.DeleteByPk(cart.ValueCart{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCart")
	}
	return err
}

// Get SaleCartItem
func (s *cartRepo) GetSaleCartItem(primary interface{}) *cart.CartItem {
	e := cart.CartItem{}
	err := s._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCartItem")
	}
	return nil
}

// Select SaleCartItem
func (s *cartRepo) SelectSaleCartItem(where string, v ...interface{}) []*cart.CartItem {
	list := []*cart.CartItem{}
	err := s._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCartItem")
	}
	return list
}

// Save SaleCartItem
func (s *cartRepo) SaveSaleCartItem(v *cart.CartItem) (int, error) {
	id, err := orm.Save(s._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCartItem")
	}
	return id, err
}

// Delete SaleCartItem
func (s *cartRepo) DeleteSaleCartItem(primary interface{}) error {
	err := s._orm.DeleteByPk(cart.CartItem{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCartItem")
	}
	return err
}

// Batch Delete SaleCartItem
func (s *cartRepo) BatchDeleteSaleCartItem(where string, v ...interface{}) (int64, error) {
	r, err := s._orm.Delete(cart.CartItem{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCartItem")
	}
	return r, err
}

// 获取购物车
func (c *cartRepo) GetShoppingCartByKey(key string) cart.ICart {
	ca := c.GetShoppingCart(key)
	if ca != nil {
		return c.CreateCart(ca)
	}
	return nil
}

// 获取会员没有结算的购物车
func (c *cartRepo) GetMemberCurrentCart(buyerId int32) cart.ICart {
	ca := c.GetLatestCart(buyerId)
	if ca != nil {
		return c.CreateCart(ca)
	}
	return nil
}

// 获取购物车
func (c *cartRepo) GetShoppingCart(key string) *cart.ValueCart {
	var v = &cart.ValueCart{}
	if c.Connector.GetOrm().GetBy(v, "code=?", key) == nil {
		items := []*cart.CartItem{}
		c.Connector.GetOrm().Select(&items, "cart_id=?", v.Id)
		v.Items = items
		return v
	}
	return nil
}

// 获取最新的购物车
func (c *cartRepo) GetLatestCart(buyerId int32) *cart.ValueCart {
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
func (c *cartRepo) SaveShoppingCart(v *cart.ValueCart) (int32, error) {
	return orm.I32(orm.Save(c.GetOrm(), v, int(v.Id)))
}

// 移出购物车项
func (c *cartRepo) RemoveCartItem(id int32) error {
	return c.Connector.GetOrm().DeleteByPk(cart.CartItem{}, id)
}

// 保存购物车项
func (c *cartRepo) SaveCartItem(v *cart.CartItem) (int32, error) {
	return orm.I32(orm.Save(c.GetOrm(), v, int(v.Id)))
}

// 清空购物车项
func (c *cartRepo) EmptyCartItems(cartId int32) error {
	_, err := c.Connector.GetOrm().Delete(cart.CartItem{}, "cart_id=?", cartId)
	return err
}

// 删除购物车
func (c *cartRepo) DeleteCart(cartId int32) error {
	return c.Connector.GetOrm().DeleteByPk(cart.ValueCart{}, cartId)
}

// 获取会员购物车(批发)
func (c *cartRepo) GetWholesaleCart(buyerId int32) cart.ICart {
	//todo: ???
	v := c.GetSaleCart(0)
	if v != nil {
		v.Items = c.SelectSaleCartItem("cart_id=?", 0)
		return c.CreateCart(v)
	}
	return nil
}
