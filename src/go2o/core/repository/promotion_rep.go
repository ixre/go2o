/**
 * Copyright 2014 @ Ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-03 23:12
 * description :
 * history :
 */

package repository

import (
	"errors"
	"fmt"
	"github.com/atnet/gof/db"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/promotion"
	promImpl "go2o/core/domain/promotion"
	"go2o/core/infrastructure/log"
	"time"
)

var _ promotion.IPromotionRep = new(PromotionRep)

type PromotionRep struct {
	db.Connector
	_memberRep member.IMemberRep
}

func NewPromotionRep(c db.Connector, memberRep member.IMemberRep) promotion.IPromotionRep {
	return &PromotionRep{
		Connector:  c,
		_memberRep: memberRep,
	}
}

func (this *PromotionRep) GetPromotion(partnerId int) promotion.IPromotion {
	return promImpl.NewPromotion(partnerId, this, this._memberRep)
}

func (this *PromotionRep) GetCoupon(id int) *promotion.ValueCoupon {
	var e promotion.ValueCoupon
	if err := this.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

func (this *PromotionRep) SaveCoupon(c promotion.ValueCoupon) (id int, err error) {
	orm := this.Connector.GetOrm()
	var rowId int64
	if c.Id > 0 {
		_, rowId, err = orm.Save(c.Id, c)
	} else {
		_, rowId, err = orm.Save(nil, c)
	}
	return int(rowId), err
}

func (this *PromotionRep) GetCouponTake(couponId, takeId int) *promotion.ValueCouponTake {
	var v promotion.ValueCouponTake
	err := this.Connector.GetOrm().Get(takeId, &v)
	if err != nil || v.CouponId != couponId {
		return nil
	}
	return &v
}

func (this *PromotionRep) SaveCouponTake(v *promotion.ValueCouponTake) error {
	var err error
	//var n int64
	if v.Id > 0 {
		_, _, err = this.Connector.GetOrm().Save(v.Id, v)
	} else {
		_, _, err = this.Connector.GetOrm().Save(nil, v)
	}
	return err
}

func (this *PromotionRep) GetCouponTakes(couponId int) []promotion.ValueCouponTake {
	var arr []promotion.ValueCouponTake = []promotion.ValueCouponTake{}

	err := this.Connector.GetOrm().SelectByQuery(&arr,
		"SELECT * FROM pm_coupon_take WHERE coupon_id =?", couponId)
	if err != nil {
		log.PrintErr(err)
	}

	return arr
}

func (this *PromotionRep) GetCouponBind(couponId, bindId int) *promotion.ValueCouponBind {
	var v promotion.ValueCouponBind
	err := this.Connector.GetOrm().Get(bindId, &v)
	if err != nil || v.CouponId != couponId {
		return nil
	}
	return &v
}

func (this *PromotionRep) GetCouponBinds(couponId int) []promotion.ValueCouponBind {
	var arr []promotion.ValueCouponBind = []promotion.ValueCouponBind{}
	err := this.Connector.GetOrm().SelectByQuery(arr,
		"SELECT * FROM pm_coupon_bind WHERE coupon_id = ?", couponId)
	if err != nil {
		log.PrintErr(err)
	}
	return arr
}

func (this *PromotionRep) SaveCouponBind(v *promotion.ValueCouponBind) error {
	var err error
	var n int64
	if v.Id > 0 {
		_, _, err = this.Connector.GetOrm().Save(v.Id, v)
	} else {
		n, _, err = this.Connector.GetOrm().Save(nil, v)
		if n != 1 {
			err = errors.New("绑定失败")
		}
	}
	return err
}

// 获取会员的优惠券绑定
func (this *PromotionRep) GetCouponBindByMemberId(couponId, memberId int) (
	*promotion.ValueCouponBind, error) {
	var bind promotion.ValueCouponBind
	err := this.Connector.GetOrm().GetByQuery(&bind,
		fmt.Sprintf(`SELECT * FROM pm_coupon_bind WHERE is_used=0
		  AND coupon_id = %d AND member_id=%d`,
			couponId, memberId),
	)
	if err == nil {
		return &bind, err
	}
	return nil, err
}

// 获取会员的优惠券占用
func (this *PromotionRep) GetCouponTakeByMemberId(couponId, memberId int) (*promotion.ValueCouponTake, error) {
	var take promotion.ValueCouponTake
	unix := time.Now().Unix()
	err := this.Connector.GetOrm().GetByQuery(&take,
		fmt.Sprintf(`SELECT * FROM pm_coupon_take WHERE
		 is_apply=0 AND extra_time>%d
		  AND coupon_id = %d AND member_id=%d`,
			unix, couponId, memberId),
	)
	if err == nil {
		return &take, err
	}
	return nil, err
}

// 根据优惠券编号获取优惠券
func (this *PromotionRep) GetCouponByCode(partnerId int, couponCode string) (
	promotion.ICoupon, error) {
	var e promotion.ValueCoupon
	err := this.Connector.GetOrm().GetByQuery(&e,
		fmt.Sprintf(`SELECT * FROM pm_coupon WHERE pt_id=%d AND code='%s'`,
			partnerId, couponCode))

	if err == nil {
		coupon := this.GetPromotion(partnerId).CreateCoupon(&e)
		return coupon, err
	}
	return nil, err
}
