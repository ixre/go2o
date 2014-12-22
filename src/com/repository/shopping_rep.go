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
	sp "com/domain/interface/shopping"
	"com/domain/shopping"
	"com/infrastructure/domain"
	"errors"
	"ops/cf/db"
)

type ShoppingRep struct {
	db.Connector
	SaleRep    sale.ISaleRep
	PromRep    promotion.IPromotionRep
	MemberRep  member.IMemberRep
	PartnerRep partner.IPartnerRep
}

func (this *ShoppingRep) GetShopping(partnerId int) sp.IShopping {
	if this.SaleRep == nil {
		panic("saleRep uninitialize!")
	}
	return shopping.NewShopping(partnerId, this.PartnerRep,
		this, this.SaleRep, this.PromRep, this.MemberRep)
}

// 获取可用的订单号
func (this *ShoppingRep) GetFreeOrderNo(partnerId int) string {
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
func (this *ShoppingRep) SaveOrder(partnerId int, v *sp.ValueOrder) (int, error) {
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
func (this *ShoppingRep) SaveOrderCouponBind(val *sp.OrderCoupon) error {
	_, _, err := this.Connector.GetOrm().Save(nil, val)
	return err
}

func (this *ShoppingRep) GetOrderByNo(partnerId int, orderNo string) (
	*sp.ValueOrder, error) {
	var v = new(sp.ValueOrder)
	err := this.Connector.GetOrm().GetBy(v, "order_no='"+orderNo+"'")
	if err != nil {
		return nil, err
	}
	return v, err
}
