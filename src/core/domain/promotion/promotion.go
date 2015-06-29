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
	"errors"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/domain/interface/sale"
	"time"
)

var _ promotion.IPromotion = new(Promotion)

type Promotion struct {
	_memberRep member.IMemberRep
	_partnerId int
	_promRep   promotion.IPromotionRep
	_value     *promotion.ValuePromotion
	_goodsRep  sale.IGoodsRep
}

func newPromotion(rep promotion.IPromotionRep, goodsRep sale.IGoodsRep,
	memRep member.IMemberRep, v *promotion.ValuePromotion) *Promotion {
	return &Promotion{
		_promRep:   rep,
		_memberRep: memRep,
		_goodsRep:  goodsRep,
		_value:     v,
	}
}

// 获取聚合根编号
func (this *Promotion) GetAggregateRootId() int {
	if this._value != nil {
		return this._value.Id
	}
	return 0
}

// 获取值
func (this *Promotion) GetValue() *promotion.ValuePromotion {
	return this._value
}

// 获取相关的值
func (this *Promotion) GetRelationValue() interface{} {
	panic(errors.New("not implement!"))
}

// 设置值
func (this *Promotion) SetValue(v *promotion.ValuePromotion) error {

	// 一种促销只能有一个?
	//todo: 每个商户设置不一样
	if false {
		if this.GetAggregateRootId() == 0 && this._value.GoodsId > 0 {
			if this._promRep.GetGoodsPromotionId(this._value.GoodsId, this._value.TypeFlag) > 0 {
				return promotion.ErrExistsSamePromotionFlag
			}
		}
	}

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
func (this *Promotion) Save() (int, error) {
	this._value.UpdateTime = time.Now().Unix()
	return this._promRep.SaveValuePromotion(this._value)
}
