/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-04 13:48
 * description :
 * history :
 */

package promotion

import (
	"errors"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale/goods"
	"time"
)

var _ promotion.IPromotion = new(promotionImpl)

type promotionImpl struct {
	_memberRep  member.IMemberRep
	_merchantId int
	_promRep    promotion.IPromotionRep
	_value      *promotion.PromotionInfo
	_goodsRep   goods.IGoodsRep
}

func newPromotion(rep promotion.IPromotionRep, goodsRep goods.IGoodsRep,
	memRep member.IMemberRep, v *promotion.PromotionInfo) *promotionImpl {
	return &promotionImpl{
		_promRep:   rep,
		_memberRep: memRep,
		_goodsRep:  goodsRep,
		_value:     v,
	}
}

// 获取聚合根编号
func (this *promotionImpl) GetAggregateRootId() int {
	if this._value != nil {
		return this._value.Id
	}
	return 0
}

// 获取值
func (this *promotionImpl) GetValue() *promotion.PromotionInfo {
	return this._value
}

// 获取相关的值
func (this *promotionImpl) GetRelationValue() interface{} {
	panic(errors.New("not implement!"))
}

// 设置值
func (this *promotionImpl) SetValue(v *promotion.PromotionInfo) error {

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
func (this *promotionImpl) ApplyFor() int {
	if this._value.GoodsId > 0 {
		return promotion.ApplyForGoods
	}
	return promotion.ApplyForOrder
}

// 促销类型
func (this *promotionImpl) Type() int {
	return this._value.TypeFlag
}

// 促销类型
func (this *promotionImpl) TypeName() string {
	panic(errors.New("not implement"))
}

// 保存
func (this *promotionImpl) Save() (int, error) {
	this._value.UpdateTime = time.Now().Unix()
	return this._promRep.SaveValuePromotion(this._value)
}
