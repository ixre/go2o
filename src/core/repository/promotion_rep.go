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
	"github.com/atnet/gof/db"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/promotion"
	promImpl "go2o/src/core/domain/promotion"
)

var _ promotion.IPromotionRep = new(promotionRep)

type promotionRep struct {
	db.Connector
	_memberRep member.IMemberRep
}

func NewPromotionRep(c db.Connector, memberRep member.IMemberRep) promotion.IPromotionRep {
	return &promotionRep{
		Connector:  c,
		_memberRep: memberRep,
	}
}

// 获取促销
func (this *promotionRep) GetValuePromotion(id int)*promotion.ValuePromotion{
	var e promotion.ValuePromotion
	if err := this.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}


// 获取促销
func (this *promotionRep) GetPromotion(id int)promotion.IPromotion{
	v := this.GetValuePromotion(id)
	if v!= nil{
		return this.CreatePromotion(v)
	}
	return nil
}

// 获取促销
func (this *promotionRep) CreatePromotion(v *promotion.ValuePromotion)promotion.IPromotion{
	return promImpl.FactoryPromotion(this,v)
}

// 保存促销
func (this *promotionRep) SaveValuePromotion(v *promotion.ValuePromotion)(int,error){
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


// 保存返现促销
func (this *promotionRep) SaveValueCashBack(v *promotion.ValueCashBack,create bool)(int,error){
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
func (this *promotionRep)  GetValueCashBack(id int)*promotion.ValueCashBack{
	var e promotion.ValueCashBack
	if err := this.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}