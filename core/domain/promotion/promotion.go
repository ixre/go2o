/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-12-04 13:48
 * description :
 * history :
 */

package promotion

import (
	"errors"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/promotion"
	"time"
)

var _ promotion.IPromotion = new(promotionImpl)

type promotionImpl struct {
	memberRepo member.IMemberRepo
	mchId      int
	promRepo   promotion.IPromotionRepo
	value      *promotion.PromotionInfo
	goodsRepo  item.IGoodsItemRepo
}

func newPromotion(rep promotion.IPromotionRepo, goodsRepo item.IGoodsItemRepo,
	memRepo member.IMemberRepo, v *promotion.PromotionInfo) *promotionImpl {
	return &promotionImpl{
		promRepo:   rep,
		memberRepo: memRepo,
		goodsRepo:  goodsRepo,
		value:      v,
	}
}

// 获取聚合根编号
func (p *promotionImpl) GetAggregateRootId() int32 {
	if p.value != nil {
		return p.value.Id
	}
	return 0
}

// 获取值
func (p *promotionImpl) GetValue() *promotion.PromotionInfo {
	return p.value
}

// 获取相关的值
func (p *promotionImpl) GetRelationValue() interface{} {
	panic(errors.New("not implement!"))
}

// 设置值
func (p *promotionImpl) SetValue(v *promotion.PromotionInfo) error {

	// 一种促销只能有一个?
	//todo: 每个商户设置不一样
	if false {
		if p.GetAggregateRootId() == 0 && p.value.GoodsId > 0 {
			if p.promRepo.GetGoodsPromotionId(p.value.GoodsId, p.value.TypeFlag) > 0 {
				return promotion.ErrExistsSamePromotionFlag
			}
		}
	}

	p.value = v
	return nil
}

// 应用类型
func (p *promotionImpl) ApplyFor() int {
	if p.value.GoodsId > 0 {
		return promotion.ApplyForGoods
	}
	return promotion.ApplyForOrder
}

// 促销类型
func (p *promotionImpl) Type() int {
	return p.value.TypeFlag
}

// 促销类型
func (p *promotionImpl) TypeName() string {
	panic(errors.New("not implement"))
}

// 保存
func (p *promotionImpl) Save() (int32, error) {
	p.value.UpdateTime = time.Now().Unix()
	return p.promRepo.SaveValuePromotion(p.value)
}
