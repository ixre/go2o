/**
 * Copyright 2015 @ to2.net.
 * name : cash_back
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package promotion

import (
	"errors"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/promotion"
	"regexp"
)

var (
	tagRegexp = regexp.MustCompile("\\s*([^:\\|]+):([^:\\|]+)\\s*\\|*")
)

var _ promotion.ICashBackPromotion = new(CashBackPromotion)

type CashBackPromotion struct {
	*promotionImpl
	cashBackValue *promotion.ValueCashBack
	dataTag       map[string]string
}

// 获取领域编号
func (c *CashBackPromotion) GetDomainId() int32 {
	return c.cashBackValue.Id
}

// 设置详细的促销信息
func (c *CashBackPromotion) SetDetailsValue(v *promotion.ValueCashBack) error {
	g := c.goodsRepo.GetValueGoodsById(c.value.GoodsId)
	if g == nil {
		return item.ErrNoSuchItem
	}

	//todo: 商品SKU的原因，获取的价格为0，有BUG

	//	if v.BackFee > int(g.SalePrice){
	//		return item.ErrOutOfSalePrice
	//	}

	if len(v.DataTag) != 0 {
		if !tagRegexp.MatchString(v.DataTag) {
			return errors.New("自定义数据格式错误！正确的格式如：\\\"K:V | K:V\\\"")
		}
	}

	c.cashBackValue = v
	return nil
}

// 获取相关的值
func (c *CashBackPromotion) GetRelationValue() interface{} {
	return c.cashBackValue
}

// 促销类型
func (c *CashBackPromotion) TypeName() string {
	return "返现"
}

// 获取自定义数据
func (c *CashBackPromotion) GetDataTag() map[string]string {
	if c.dataTag == nil {
		c.dataTag = make(map[string]string)
		if len(c.cashBackValue.DataTag) != 0 {
			matches := tagRegexp.FindAllStringSubmatch(c.cashBackValue.DataTag, -1)
			for i := 0; i < len(matches); i++ {
				c.dataTag[matches[i][1]] = matches[i][2]
			}
		}
	}
	return c.dataTag
}

// 保存
func (c *CashBackPromotion) Save() (int32, error) {

	if c.GetRelationValue() == nil {
		return c.GetAggregateRootId(), promotion.ErrCanNotApplied
	}

	var isCreate bool = c.GetAggregateRootId() == 0
	c.value.TypeFlag = promotion.TypeFlagCashBack
	id, err := c.promotionImpl.Save()
	if err == nil {
		c.value.Id = id
		if c.cashBackValue == nil {
			c.cashBackValue = new(promotion.ValueCashBack)
		}
		c.cashBackValue.Id = c.GetAggregateRootId()
		_, err = c.promRepo.SaveValueCashBack(c.cashBackValue, isCreate)
	}
	return id, err
}
