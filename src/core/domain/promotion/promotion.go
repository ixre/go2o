/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-04 13:48
 * description :
 * history :
 */

package promotion

import (
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/promotion"
	"time"
	"errors"
)


var _ promotion.IPromotion = new(Promotion)

type Promotion struct {
	promRep   promotion.IOldPromotionRep
	memberRep member.IMemberRep
	_partnerId int
	_promRep promotion.IPromotionRep
	_value *promotion.ValuePromotion
}

func NewPromotion(partnerId int,promRep promotion.IOldPromotionRep,
	rep promotion.IPromotionRep,
	memberRep member.IMemberRep) promotion.IPromotion {
	return &Promotion{
		_partnerId:  partnerId,
		_promRep:rep,
		promRep:   promRep,
		memberRep: memberRep,
	}
}

func newPromotion(rep promotion.IPromotionRep,v *promotion.ValuePromotion)*Promotion{
	return &Promotion{
		_promRep:rep,
		_value :v,
	}
}

// 获取聚合根编号
func (this *Promotion) GetAggregateRootId() int {
	return this._partnerId
}

// 获取值
func (this *Promotion)GetValue()*promotion.ValuePromotion{
	return this._value
}


// 获取相关的值
func (this *Promotion) GetRelationValue()interface{}{
	panic(errors.New("not implement!"))
}

// 设置值
func (this *Promotion) SetValue(v *promotion.ValuePromotion)error{
	this._value = v
	return nil
}


// 应用类型
func (this *Promotion) ApplyFor() int {
	if this._value.GoodsId > 0 {
		return promotion.ApplyForGoods
	}
	return promotion.ApplyForOrder
}

// 促销类型
func (this *Promotion) Type() int {
	return this._value.TypeFlag
}

// 保存
func (this *Promotion) Save()(int,error){
	return this._promRep.SaveValuePromotion(this._value)
}


func (this *Promotion) GetCoupon(id int) promotion.ICoupon {
	var val *promotion.ValueCoupon = this.promRep.GetCoupon(id)
	return this.CreateCoupon(val)
}

func (this *Promotion) CreateCoupon(val *promotion.ValueCoupon) promotion.ICoupon {
	val.PartnerId = this.GetAggregateRootId()
	val.CreateTime = time.Now().Unix()
	val.Amount = val.TotalAmount
	return newCoupon(val, this.promRep, this.memberRep)
}
