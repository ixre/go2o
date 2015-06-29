/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:12
 * description :
 * history :
 */

package repository

import (
	"errors"
	"fmt"
	"github.com/atnet/gof/db"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/domain/interface/sale"
	promImpl "go2o/src/core/domain/promotion"
	"go2o/src/core/infrastructure/log"
	"time"
)

var _ promotion.IPromotionRep = new(promotionRep)

type promotionRep struct {
	db.Connector
	_memberRep member.IMemberRep
	_goodsRep  sale.IGoodsRep
}

func NewPromotionRep(c db.Connector, goodsRep sale.IGoodsRep,
	memberRep member.IMemberRep) promotion.IPromotionRep {
	return &promotionRep{
		Connector:  c,
		_memberRep: memberRep,
		_goodsRep:  goodsRep,
	}
}

// 获取促销
func (this *promotionRep) GetValuePromotion(id int) *promotion.ValuePromotion {
	var e promotion.ValuePromotion
	if err := this.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 获取促销
func (this *promotionRep) GetPromotion(id int) promotion.IPromotion {
	v := this.GetValuePromotion(id)
	if v != nil {
		return this.CreatePromotion(v)
	}
	return nil
}

// 获取促销
func (this *promotionRep) CreatePromotion(v *promotion.ValuePromotion) promotion.IPromotion {
	return promImpl.FactoryPromotion(this, this._goodsRep, this._memberRep, v)
}

// 保存促销
func (this *promotionRep) SaveValuePromotion(v *promotion.ValuePromotion) (int, error) {
	var err error
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM pm_info WHERE partner_id=?", &v.Id, v.PartnerId)
	}
	return v.Id, err
}

// 删除促销
func (this *promotionRep) DeletePromotion(id int) error {
	return this.Connector.GetOrm().DeleteByPk(promotion.ValuePromotion{}, id)
}

// 保存返现促销
func (this *promotionRep) SaveValueCashBack(v *promotion.ValueCashBack, create bool) (int, error) {
	var err error
	var orm = this.Connector.GetOrm()
	if !create {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
	}
	return v.Id, err
}

// 获取返现促销
func (this *promotionRep) GetValueCashBack(id int) *promotion.ValueCashBack {
	var e promotion.ValueCashBack
	if err := this.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 删除返现促销
func (this *promotionRep) DeleteValueCashBack(id int) error {
	return this.Connector.GetOrm().DeleteByPk(promotion.ValueCashBack{}, id)
}

// 获取商品的促销编号
func (this *promotionRep) GetGoodsPromotionId(goodsId int, promFlag int) int {
	var id int
	this.Connector.ExecScalar("SELECT id FROM pm_info WHERE goods_id=? AND type_flag=? AND enabled=1", &id, goodsId, promFlag)
	return id
}

// 获取商品的促销
func (this *promotionRep) GetPromotionOfGoods(goodsId int) []*promotion.ValuePromotion {
	var arr []*promotion.ValuePromotion = []*promotion.ValuePromotion{}
	err := this.Connector.GetOrm().Select(&arr, "goods_id=? AND enabled=1 ORDER BY id", goodsId)
	if err == nil {
		return arr
	}

	return make([]*promotion.ValuePromotion, 0)
}

/*****   OLD ******/

func (this *promotionRep) GetValueCoupon(id int) *promotion.ValueCoupon {
	var e promotion.ValueCoupon
	if err := this.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

func (this *promotionRep) SaveValueCoupon(v *promotion.ValueCoupon, isCreate bool) (id int, err error) {
	orm := this.Connector.GetOrm()
	if isCreate {
		_, _, err = orm.Save(nil, v)
	} else {
		_, _, err = orm.Save(v.Id, v)
	}
	return v.Id, err
}

// 删除优惠券
func (this *promotionRep) DeleteValueCoupon(id int) error {
	return this.Connector.GetOrm().DeleteByPk(promotion.ValueCoupon{}, id)
}

func (this *promotionRep) GetCouponTake(couponId, takeId int) *promotion.ValueCouponTake {
	var v promotion.ValueCouponTake
	err := this.Connector.GetOrm().Get(takeId, &v)
	if err != nil || v.CouponId != couponId {
		return nil
	}
	return &v
}

func (this *promotionRep) SaveCouponTake(v *promotion.ValueCouponTake) error {
	var err error
	//var n int64
	if v.Id > 0 {
		_, _, err = this.Connector.GetOrm().Save(v.Id, v)
	} else {
		_, _, err = this.Connector.GetOrm().Save(nil, v)
	}
	return err
}

func (this *promotionRep) GetCouponTakes(couponId int) []promotion.ValueCouponTake {
	var arr []promotion.ValueCouponTake = []promotion.ValueCouponTake{}

	err := this.Connector.GetOrm().SelectByQuery(&arr,
		"SELECT * FROM pm_coupon_take WHERE coupon_id =?", couponId)
	if err != nil {
		log.PrintErr(err)
	}

	return arr
}

func (this *promotionRep) GetCouponBind(couponId, bindId int) *promotion.ValueCouponBind {
	var v promotion.ValueCouponBind
	err := this.Connector.GetOrm().Get(bindId, &v)
	if err != nil || v.CouponId != couponId {
		return nil
	}
	return &v
}

func (this *promotionRep) GetCouponBinds(couponId int) []promotion.ValueCouponBind {
	var arr []promotion.ValueCouponBind = []promotion.ValueCouponBind{}
	err := this.Connector.GetOrm().SelectByQuery(arr,
		"SELECT * FROM pm_coupon_bind WHERE coupon_id = ?", couponId)
	if err != nil {
		log.PrintErr(err)
	}
	return arr
}

func (this *promotionRep) SaveCouponBind(v *promotion.ValueCouponBind) error {
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
func (this *promotionRep) GetCouponBindByMemberId(couponId, memberId int) (
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
func (this *promotionRep) GetCouponTakeByMemberId(couponId, memberId int) (*promotion.ValueCouponTake, error) {
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

// 根据优惠券代码获取优惠券
func (this *promotionRep) GetValueCouponByCode(partnerId int, couponCode string) *promotion.ValueCoupon {
	var e promotion.ValueCoupon
	err := this.Connector.GetOrm().GetByQuery(&e,
		fmt.Sprintf(`SELECT * FROM pm_info INNER JOIN pm_coupon ON pm_info.id=pm_coupon.id
		 		WHERE partner_id=%d AND code='%s'`,
			partnerId, couponCode))

	if err == nil {
		return &e
	}
	return nil
}

// 根据代码获取优惠券
func (this *promotionRep) GetCouponByCode(partnerId int, code string) promotion.IPromotion {
	v := this.GetValueCouponByCode(partnerId, code)
	if v != nil {
		p := this.GetValuePromotion(v.Id)
		prom := this.CreatePromotion(p)
		cp := prom.(promotion.ICouponPromotion)
		cp.SetDetailsValue(v)
		return prom
	}
	return nil
}
