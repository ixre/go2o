/**
 * Copyright 2015 @ 56x.net.
 * name : cart_repo.go
 * author : jarryliu
 * date : 2016-06-29 22:54
 * description :
 * history :
 */
package repos

import (
	"database/sql"
	"log"
	"time"

	cartImpl "github.com/ixre/go2o/core/domain/cart"
	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

var _ cart.ICartRepo = new(cartRepo)

type cartRepo struct {
	db.Connector
	o           orm.Orm
	_itemRepo   item.IItemRepo
	_memberRepo member.IMemberRepo
	_mchRepo    merchant.IMerchantRepo
}

func NewCartRepo(o orm.Orm, memberRepo member.IMemberRepo,
	_mchRepo merchant.IMerchantRepo,
	itemRepo item.IItemRepo) cart.ICartRepo {
	return &cartRepo{
		Connector:   o.Connector(),
		o:           o,
		_memberRepo: memberRepo,
		_mchRepo:    _mchRepo,
		_itemRepo:   itemRepo,
	}
}

// 获取买家的购物车
func (c *cartRepo) GetMyCart(buyerId int64, k cart.Kind) cart.ICart {
	switch k {
	case cart.KNormal:
		return c.getMyNormalCart(buyerId)
	case cart.KWholesale:
		return c.getMyWholesaleCart(buyerId)
	}
	return nil
}

// 获取普通购物车
func (c *cartRepo) getMyNormalCart(buyerId int64) cart.ICart {
	v := c.getNormalCart(buyerId)
	if v != nil {
		return c.CreateNormalCart(v)
	}
	return nil
}

// 获取批发购物车
func (c *cartRepo) getMyWholesaleCart(buyerId int64) cart.ICart {
	v := c.getWholesaleCart(buyerId)
	if v == nil {
		unix := time.Now().Unix()
		v = &cart.WsCart{
			BuyerId:    buyerId,
			CreateTime: unix,
			UpdateTime: unix,
		}
	}
	return cartImpl.CreateWholesaleCart(v, c,
		c._memberRepo, c._mchRepo, c._itemRepo)
}

func (c *cartRepo) getNormalCart(buyerId int64) *cart.NormalCart {
	e := cart.NormalCart{}
	err := c.o.GetBy(&e, "buyer_id= $1", buyerId)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCart")
	}
	return nil
}

func (w *cartRepo) getWholesaleCart(buyerId int64) *cart.WsCart {
	e := cart.WsCart{}
	err := w.o.GetBy(&e, "buyer_id= $1", buyerId)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsCart")
	}
	return nil
}

// 创建购物车对象
func (c *cartRepo) CreateNormalCart(v *cart.NormalCart) cart.ICart {
	return cartImpl.NewNormalCart(v, c, c._memberRepo, c._itemRepo)
}

// 创建一个购物车
func (c *cartRepo) NewTempNormalCart(buyerId int, code string) cart.ICart {
	return cartImpl.CreateTempNormalCart(buyerId, code, c, c._memberRepo, c._itemRepo)
}

// 获取购物车
func (c *cartRepo) GetNormalCart(id int32) cart.ICart {
	v := c.getSaleCart(id)
	if v != nil {
		return c.CreateNormalCart(v)
	}
	return nil
}

// Get SaleCart
func (s *cartRepo) getSaleCart(primary interface{}) *cart.NormalCart {
	e := cart.NormalCart{}
	err := s.o.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCart")
	}
	return nil
}

// Save SaleCart
func (s *cartRepo) SaveNormalCart(v *cart.NormalCart) (int, error) {
	id, err := orm.Save(s.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCart")
	}
	return id, err
}

// Delete SaleCart
func (s *cartRepo) DeleteNormalCart(primary interface{}) error {
	err := s.o.DeleteByPk(cart.NormalCart{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCart")
	}
	return err
}

// Get SaleCartItem
func (s *cartRepo) GetSaleCartItem(primary interface{}) *cart.NormalCartItem {
	e := cart.NormalCartItem{}
	err := s.o.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCartItem")
	}
	return nil
}

// Select SaleCartItem
func (s *cartRepo) SelectNormalCartItem(where string, v ...interface{}) []*cart.NormalCartItem {
	list := []*cart.NormalCartItem{}
	err := s.o.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCartItem")
	}
	return list
}

// Save SaleCartItem
func (s *cartRepo) SaveNormalCartItem(v *cart.NormalCartItem) (int, error) {
	id, err := orm.Save(s.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCartItem")
	}
	return id, err
}

// Delete SaleCartItem
func (s *cartRepo) DeleteNormalCartItem(primary interface{}) error {
	err := s.o.DeleteByPk(cart.NormalCartItem{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCartItem")
	}
	return err
}

// Batch Delete SaleCartItem
func (s *cartRepo) BatchDeleteNormalCartItem(where string, v ...interface{}) (int64, error) {
	r, err := s.o.Delete(cart.NormalCartItem{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleCartItem")
	}
	return r, err
}

// 获取购物车
func (c *cartRepo) GetShoppingCartByKey(key string) cart.ICart {
	ca := c.GetShoppingCart(key)
	if ca != nil {
		return c.CreateNormalCart(ca)
	}
	return nil
}

// 获取购物车
func (c *cartRepo) GetShoppingCart(key string) *cart.NormalCart {
	var v = &cart.NormalCart{}
	if c.o.GetBy(v, "code= $1", key) == nil {
		var items []*cart.NormalCartItem
		c.o.Select(&items, "cart_id= $1", v.Id)
		v.Items = items
		return v
	}
	return nil
}

// 获取最新的购物车
func (c *cartRepo) GetLatestCart(buyerId int64) *cart.NormalCart {
	var v = &cart.NormalCart{}
	if c.o.GetBy(v, "buyer_id= $1 ORDER BY id DESC", buyerId) == nil {
		var items []*cart.NormalCartItem
		c.o.Select(&items, "cart_id= $1", v.Id)
		v.Items = items
		return v
	}
	return nil
}

// 保存购物车
func (c *cartRepo) SaveShoppingCart(v *cart.NormalCart) (int32, error) {
	return orm.I32(orm.Save(c.o, v, int(v.Id)))
}

// 移出购物车项
func (c *cartRepo) RemoveCartItem(id int32) error {
	return c.o.DeleteByPk(cart.NormalCartItem{}, id)
}

// 保存购物车项
func (c *cartRepo) SaveCartItem(v *cart.NormalCartItem) (int32, error) {
	return orm.I32(orm.Save(c.o, v, int(v.Id)))
}

// 清空购物车项
func (c *cartRepo) EmptyCartItems(cartId int32) error {
	_, err := c.o.Delete(cart.NormalCartItem{}, "cart_id= $1", cartId)
	return err
}

// 删除购物车
func (c *cartRepo) DeleteCart(cartId int32) error {
	return c.o.DeleteByPk(cart.NormalCart{}, cartId)
}

// Save WsCart
func (w *cartRepo) SaveWsCart(v *cart.WsCart) (int, error) {
	id, err := orm.Save(w.o, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsCart")
	}
	return id, err
}

// Delete WsCart
func (w *cartRepo) DeleteWsCart(primary interface{}) error {
	err := w.o.DeleteByPk(cart.WsCart{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsCart")
	}
	return err
}

// Select WsCartItem
func (w *cartRepo) SelectWsCartItem(where string, v ...interface{}) []*cart.WsCartItem {
	var list []*cart.WsCartItem
	err := w.o.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsCartItem")
	}
	return list
}

// Save WsCartItem
func (w *cartRepo) SaveWsCartItem(v *cart.WsCartItem) (int, error) {
	id, err := orm.Save(w.o, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsCartItem")
	}
	return id, err
}

// Batch Delete WsCartItem
func (w *cartRepo) BatchDeleteWsCartItem(where string, v ...interface{}) (int64, error) {
	r, err := w.o.Delete(cart.WsCartItem{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsCartItem")
	}
	return r, err
}
