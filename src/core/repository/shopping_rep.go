/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:53
 * description :
 * history :
 */

package repository

import (
	"errors"
	"github.com/atnet/gof/db"
	"go2o/src/core/domain/interface/delivery"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/domain/interface/shopping"
	shoppingImpl "go2o/src/core/domain/shopping"
	"go2o/src/core/infrastructure/domain"
	"github.com/atnet/gof"
	"go2o/src/core/variable"
)

var _ shopping.IShoppingRep = new(shoppingRep)

type shoppingRep struct {
	db.Connector
	_saleRep    sale.ISaleRep
	_goodsRep   sale.IGoodsRep
	_promRep    promotion.IPromotionRep
	_memberRep  member.IMemberRep
	_partnerRep partner.IPartnerRep
	_deliverRep delivery.IDeliveryRep
	_cache      map[int]shopping.IShopping
}

func NewShoppingRep(c db.Connector, ptRep partner.IPartnerRep,
	saleRep sale.ISaleRep, goodsRep sale.IGoodsRep, promRep promotion.IPromotionRep,
	memRep member.IMemberRep, deliverRep delivery.IDeliveryRep) shopping.IShoppingRep {
	return (&shoppingRep{
		Connector:   c,
		_saleRep:    saleRep,
		_goodsRep:   goodsRep,
		_promRep:    promRep,
		_memberRep:  memRep,
		_partnerRep: ptRep,
		_deliverRep: deliverRep,
	}).init()
}

func (this *shoppingRep) init() shopping.IShoppingRep {
	this._cache = make(map[int]shopping.IShopping)
	return this
}

func (this *shoppingRep) GetShopping(partnerId int) shopping.IShopping {
	if this._saleRep == nil {
		panic("saleRep uninitialize!")
	}
	v, ok := this._cache[partnerId]
	if !ok {
		v = shoppingImpl.NewShopping(partnerId, this._partnerRep,
			this, this._saleRep, this._goodsRep, this._promRep, this._memberRep, this._deliverRep)
		this._cache[partnerId] = v
	}
	return v
}

// 获取可用的订单号
func (this *shoppingRep) GetFreeOrderNo(partnerId int) string {
	//todo:实际应用需要预先生成订单号
	d := this.Connector
	var order_no string
	for {
		order_no = domain.NewOrderNo(partnerId)
		var rec int
		if d.ExecScalar(`SELECT COUNT(0) FROM pt_order where order_no=?`,
			&rec, order_no); rec == 0 {
			break
		}
	}
	return order_no
}
func (this *shoppingRep) SaveOrder(partnerId int, v *shopping.ValueOrder) (int, error) {
	var err error
	d := this.Connector
	v.PartnerId = partnerId

	if v.Id > 0 {
		_, _, err = d.GetOrm().Save(v.Id, v)
	} else {
		//验证Partner和Member是否有绑定关系
		var num int
		if d.ExecScalar(`SELECT COUNT(0) FROM mm_relation WHERE member_id=? AND reg_partner_id=?`,
			&num, v.MemberId, v.PartnerId); num != 1 {
			return v.Id, errors.New("error partner and member.")
		}
		_, _, err = d.GetOrm().Save(nil, v)
		if err == nil {
			err = d.ExecScalar(`SELECT MAX(id) FROM pt_order WHERE partner_id=? AND member_id=?`, &v.Id,
				partnerId, v.MemberId)
		}

		// Sign new order
		gof.CurrentApp.Storage().Set(variable.KvHaveNewOrder, 1)
	}

	// 保存订单项
	if err == nil && v.Items != nil {
		orm := d.GetOrm()
		for _, v1 := range v.Items {
			if v1.Id > 0 {
				orm.Save(v1.Id, v1)
			} else {
				orm.Save(nil, v1)
			}
		}
	}

	return v.Id, err
}

//　保存订单优惠券绑定
func (this *shoppingRep) SaveOrderCouponBind(val *shopping.OrderCoupon) error {
	_, _, err := this.Connector.GetOrm().Save(nil, val)
	return err
}

func (this *shoppingRep) GetOrderByNo(partnerId int, orderNo string) (
	*shopping.ValueOrder, error) {
	var v = new(shopping.ValueOrder)
	err := this.Connector.GetOrm().GetBy(v, "partner_id=? AND order_no=?", partnerId, orderNo)
	if err != nil {
		return nil, err
	}
	return v, err
}

// 获取等待处理的订单
func (this *shoppingRep) GetWaitingSetupOrders(partnerId int) ([]*shopping.ValueOrder, error) {
	dst := []*shopping.ValueOrder{}
	err := this.Connector.GetOrm().Select(&dst,
		"partner_id=? AND is_suspend=0 AND status IN("+enum.ORDER_SETUP_STATE+")",
		partnerId)
	if err != nil {
		return nil, err
	}
	return dst, err
}

// 保存订单日志
func (this *shoppingRep) SaveOrderLog(v *shopping.OrderLog) error {
	_, _, err := this.Connector.GetOrm().Save(nil, v)
	return err
}

// 获取购物车
func (this *shoppingRep) GetShoppingCart(key string) (*shopping.ValueCart, error) {
	var v = new(shopping.ValueCart)
	err := this.Connector.GetOrm().GetBy(v, "cart_key=?", key)
	if v == nil || err != nil {
		return nil, err
	}

	var items = []*shopping.ValueCartItem{}
	err = this.Connector.GetOrm().Select(&items, "cart_id=?", v.Id)
	if err == nil {
		v.Items = items
	}

	return v, err
}

// 获取最新的购物车
func (this *shoppingRep) GetLatestCart(buyerId int) (*shopping.ValueCart, error) {
	var v = new(shopping.ValueCart)
	err := this.Connector.GetOrm().GetBy(v, "buyer_id=? ORDER BY id DESC", buyerId)
	if v == nil || err != nil {
		return nil, err
	}

	var items = []*shopping.ValueCartItem{}
	err = this.Connector.GetOrm().Select(&items, "cart_id=?", v.Id)
	if err == nil {
		v.Items = items
	}
	return v, err
}

// 保存购物车
func (this *shoppingRep) SaveShoppingCart(v *shopping.ValueCart) (int, error) {
	var err error
	_orm := this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = _orm.Save(v.Id, v)
	} else {
		_, _, err = _orm.Save(nil, v)
		this.Connector.ExecScalar(`SELECT MAX(id) FROM sale_cart`, &v.Id)
	}
	return v.Id, err
}

// 移出购物车项
func (this *shoppingRep) RemoveCartItem(id int) error {
	return this.Connector.GetOrm().DeleteByPk(shopping.ValueCartItem{}, id)
}

// 保存购物车项
func (this *shoppingRep) SaveCartItem(v *shopping.ValueCartItem) (int, error) {
	_orm := this.Connector.GetOrm()
	var err error
	if v.Id > 0 {
		_, _, err = _orm.Save(v.Id, v)
	} else {
		_, _, err = _orm.Save(nil, v)
		this.Connector.ExecScalar(`SELECT MAX(id) FROM sale_cart_item where cart_id=?`, &v.Id, v.CartId)
	}

	return v.Id, err
}

// 清空购物车项
func (this *shoppingRep) EmptyCartItems(id int) error {
	_, err := this.Connector.GetOrm().Delete(shopping.ValueCartItem{}, "cart_id=?", id)
	return err
}

// 删除购物车
func (this *shoppingRep) DeleteCart(id int) error {
	return this.Connector.GetOrm().DeleteByPk(shopping.ValueCart{}, id)
}

// 获取订单的促销绑定
func (this *shoppingRep) GetOrderPromotionBinds(orderNo string) []*shopping.OrderPromotionBind {
	var arr []*shopping.OrderPromotionBind = []*shopping.OrderPromotionBind{}
	err := this.Connector.GetOrm().Select(&arr, "order_no=?", orderNo)
	if err == nil {
		return arr
	}
	return make([]*shopping.OrderPromotionBind, 0)
}

// 保存订单的促销绑定
func (this *shoppingRep) SavePromotionBindForOrder(v *shopping.OrderPromotionBind) (int, error) {
	var err error
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM pt_order_pb WHERE order_no=? AND promotion_id=?",
			&v.Id, v.OrderNo, v.PromotionId)
	}
	return v.Id, err
}
