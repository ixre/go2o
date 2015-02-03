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
	"com/domain/interface/member"
	"com/domain/interface/partner"
	"com/domain/interface/promotion"
	"com/domain/interface/sale"
	"com/domain/interface/shopping"
	sp "com/domain/shopping"
	"com/infrastructure/domain"
	"errors"
	"github.com/atnet/gof/db"
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
			err = d.ExecScalar(`SELECT MAX(id) FROM pt_order WHERE pt_id=? AND mid=?`, &v.Id,
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
	err := this.Connector.GetOrm().GetBy(v, "order_no='"+orderNo+"'")
	if err != nil {
		return nil, err
	}
	return v, err
}
