/**
 * Copyright 2014 @ z3q.net.
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
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale/goods"
	promImpl "go2o/core/domain/promotion"
	"go2o/core/infrastructure/log"
	"time"
)

var _ promotion.IPromotionRep = new(promotionRep)

type promotionRep struct {
	db.Connector
	_memberRep member.IMemberRep
	_goodsRep  goods.IGoodsRep
}

func NewPromotionRep(c db.Connector, goodsRep goods.IGoodsRep,
	memberRep member.IMemberRep) promotion.IPromotionRep {
	return &promotionRep{
		Connector:  c,
		_memberRep: memberRep,
		_goodsRep:  goodsRep,
	}
}

// 获取促销
func (this *promotionRep) GetValuePromotion(id int32) *promotion.PromotionInfo {
	var e promotion.PromotionInfo
	if err := this.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 获取促销
func (this *promotionRep) GetPromotion(id int32) promotion.IPromotion {
	v := this.GetValuePromotion(id)
	if v != nil {
		return this.CreatePromotion(v)
	}
	return nil
}

// 获取促销
func (this *promotionRep) CreatePromotion(v *promotion.PromotionInfo) promotion.IPromotion {
	return promImpl.FactoryPromotion(this, this._goodsRep, this._memberRep, v)
}

// 保存促销
func (this *promotionRep) SaveValuePromotion(v *promotion.PromotionInfo) (int32, error) {
	return orm.I32(orm.Save(this.GetOrm(), v, int(v.Id)))
}

// 删除促销
func (this *promotionRep) DeletePromotion(id int32) error {
	return this.Connector.GetOrm().DeleteByPk(promotion.PromotionInfo{}, id)
}

// 保存返现促销
func (this *promotionRep) SaveValueCashBack(v *promotion.ValueCashBack, create bool) (int32, error) {
	return orm.I32(orm.Save(this.GetOrm(), v, int(v.Id)))
}

// 获取返现促销
func (this *promotionRep) GetValueCashBack(id int32) *promotion.ValueCashBack {
	var e promotion.ValueCashBack
	if err := this.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 删除返现促销
func (this *promotionRep) DeleteValueCashBack(id int32) error {
	return this.Connector.GetOrm().DeleteByPk(promotion.ValueCashBack{}, id)
}

// 获取商品的促销编号
func (this *promotionRep) GetGoodsPromotionId(goodsId int32, promFlag int) int {
	var id int
	this.Connector.ExecScalar("SELECT id FROM pm_info WHERE goods_id=? AND type_flag=? AND enabled=1", &id, goodsId, promFlag)
	return id
}

// 获取商品的促销
func (this *promotionRep) GetPromotionOfGoods(goodsId int32) []*promotion.PromotionInfo {
	var arr []*promotion.PromotionInfo = []*promotion.PromotionInfo{}
	err := this.Connector.GetOrm().Select(&arr, "goods_id=? AND enabled=1 ORDER BY id", goodsId)
	if err == nil {
		return arr
	}

	return make([]*promotion.PromotionInfo, 0)
}

// 获取商户订单可用的促销
func (this *promotionRep) GetPromotionOfMerchantOrder(mchId int32) []*promotion.PromotionInfo {
	var arr []*promotion.PromotionInfo = []*promotion.PromotionInfo{}
	err := this.Connector.GetOrm().Select(&arr, "mch_id=? AND goods_id=0 AND enabled=1 ORDER BY id", mchId)
	if err == nil {
		return arr
	}
	return make([]*promotion.PromotionInfo, 0)
}

/*****   OLD ******/

func (this *promotionRep) GetValueCoupon(id int32) *promotion.ValueCoupon {
	var e promotion.ValueCoupon
	if err := this.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

func (this *promotionRep) SaveValueCoupon(v *promotion.ValueCoupon, isCreate bool) (id int32, err error) {
	return orm.I32(orm.Save(this.GetOrm(), v, int(v.Id)))
}

// 删除优惠券
func (this *promotionRep) DeleteValueCoupon(id int32) error {
	return this.Connector.GetOrm().DeleteByPk(promotion.ValueCoupon{}, id)
}

func (this *promotionRep) GetCouponTake(couponId, takeId int32) *promotion.ValueCouponTake {
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

func (this *promotionRep) GetCouponTakes(couponId int32) []promotion.ValueCouponTake {
	var arr []promotion.ValueCouponTake = []promotion.ValueCouponTake{}

	err := this.Connector.GetOrm().SelectByQuery(&arr,
		"SELECT * FROM pm_coupon_take WHERE coupon_id =?", couponId)
	if err != nil {
		log.Error(err)
	}

	return arr
}

func (this *promotionRep) GetCouponBind(couponId, bindId int32) *promotion.ValueCouponBind {
	var v promotion.ValueCouponBind
	err := this.Connector.GetOrm().Get(bindId, &v)
	if err != nil || v.CouponId != couponId {
		return nil
	}
	return &v
}

func (this *promotionRep) GetCouponBinds(couponId int32) []promotion.ValueCouponBind {
	var arr []promotion.ValueCouponBind = []promotion.ValueCouponBind{}
	err := this.Connector.GetOrm().SelectByQuery(arr,
		"SELECT * FROM pm_coupon_bind WHERE coupon_id = ?", couponId)
	if err != nil {
		log.Error(err)
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
func (this *promotionRep) GetCouponBindByMemberId(couponId, memberId int32) (
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
func (this *promotionRep) GetCouponTakeByMemberId(couponId, memberId int32) (*promotion.ValueCouponTake, error) {
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
func (this *promotionRep) GetValueCouponByCode(mchId int32, couponCode string) *promotion.ValueCoupon {
	var e promotion.ValueCoupon
	err := this.Connector.GetOrm().GetByQuery(&e,
		fmt.Sprintf(`SELECT * FROM pm_info INNER JOIN pm_coupon ON pm_info.id=pm_coupon.id
		 		WHERE mch_id=%d AND code='%s'`,
			mchId, couponCode))

	if err == nil {
		return &e
	}
	return nil
}

// 根据代码获取优惠券
func (this *promotionRep) GetCouponByCode(mchId int32, code string) promotion.IPromotion {
	v := this.GetValueCouponByCode(mchId, code)
	if v != nil {
		p := this.GetValuePromotion(v.Id)
		prom := this.CreatePromotion(p)
		cp := prom.(promotion.ICouponPromotion)
		cp.SetDetailsValue(v)
		return prom
	}
	return nil
}
