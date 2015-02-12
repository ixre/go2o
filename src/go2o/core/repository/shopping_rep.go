/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-05 17:53
 * description :
 * history :
 */

package repository

import (
	"errors"
	"github.com/atnet/gof/db"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/partner"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/shopping"
	sp "go2o/core/domain/shopping"
	"go2o/core/infrastructure/domain"
)

var _ shopping.IShoppingRep = new(shoppingRep)

type shoppingRep struct {
	db.Connector
	saleRep    sale.ISaleRep
	promRep    promotion.IPromotionRep
	memberRep  member.IMemberRep
	partnerRep partner.IPartnerRep
	cache      map[int]shopping.IShopping
}

func NewShoppingRep(c db.Connector, ptRep partner.IPartnerRep,
	saleRep sale.ISaleRep, promRep promotion.IPromotionRep,
	memRep member.IMemberRep) shopping.IShoppingRep {
	return (&shoppingRep{
		Connector:  c,
		saleRep:    saleRep,
		promRep:    promRep,
		memberRep:  memRep,
		partnerRep: ptRep,
	}).init()
}

func (this *shoppingRep) init() shopping.IShoppingRep {
	this.cache = make(map[int]shopping.IShopping)
	return this
}

func (this *shoppingRep) GetShopping(partnerId int) shopping.IShopping {
	if this.saleRep == nil {
		panic("saleRep uninitialize!")
	}
	v, ok := this.cache[partnerId]
	if !ok {
		v = sp.NewShopping(partnerId, this.partnerRep,
			this, this.saleRep, this.promRep, this.memberRep)
		this.cache[partnerId] = v
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
		if d.ExecScalar(`SELECT COUNT(0) FROM mm_relation WHERE member_id=? AND reg_ptid=?`,
			&num, v.MemberId, v.PartnerId); num != 1 {
			return v.Id, errors.New("error partner and member.")
		}
		_, _, err = d.GetOrm().Save(nil, v)
		if err == nil {
			err = d.ExecScalar(`SELECT MAX(id) FROM pt_order WHERE pt_id=? AND member_id=?`, &v.Id,
				partnerId, v.MemberId)
			return v.Id, err
		}
		return -1, err
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
	err := this.Connector.GetOrm().GetBy(v, "pt_id=? AND order_no=?", partnerId, orderNo)
	if err != nil {
		return nil, err
	}
	return v, err
}

// 获取等待处理的订单
func (this *shoppingRep) GetWaitingSetupOrders(partnerId int) ([]*shopping.ValueOrder, error) {
	dst := []*shopping.ValueOrder{}
	err := this.Connector.GetOrm().Select(&dst,
		"pt_id=? AND is_suspend=0 AND status IN("+enum.ORDER_SETUP_STATE+")",
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

	this.setAttachGoodsInfo(v.Items)

	return v, err
}

// 获取未结算的购物车
func (this *shoppingRep) GetNotBoughtCart(buyerId int) (*shopping.ValueCart, error) {
	var v = new(shopping.ValueCart)
	err := this.Connector.GetOrm().GetBy(v, "is_bought=0 AND buyer_id=?", buyerId)
	if v == nil || err != nil {
		return nil, err
	}

	var items = []*shopping.ValueCartItem{}
	err = this.Connector.GetOrm().Select(&items, "cart_id=?", v.Id)
	if err == nil {
		v.Items = items
	}
	this.setAttachGoodsInfo(v.Items)

	return v, err
}

func (this *shoppingRep) setAttachGoodsInfo(items []*shopping.ValueCartItem) {
	if items != nil {
		l := len(items)
		if l == 0 {
			return
		}
		var ids []int = make([]int, l)
		for i, v := range items {
			ids[i] = v.GoodsId
		}

		// 设置附加的值
		goods, err := this.saleRep.GetGoodsByIds(ids...)
		if err == nil {
			var goodsMap = make(map[int]*sale.ValueGoods, len(goods))
			for _, v := range goods {
				goodsMap[v.Id] = v
			}

			for _, v := range items {
				gv, ok := goodsMap[v.GoodsId]
				if ok {
					v.Name = gv.Name
					v.SmallTitle = gv.SmallTitle
					v.Price = gv.Price
					v.GoodsNo = gv.GoodsNo
					v.Image = gv.Image
					v.SalePrice = gv.SalePrice
				}
			}
		}
	}
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
	return this.Connector.GetOrm().DeleteByPk(shopping.ValueCart{}, id)
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
