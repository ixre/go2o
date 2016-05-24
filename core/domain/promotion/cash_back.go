/**
 * Copyright 2015 @ z3q.net.
 * name : cash_back
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package promotion

import (
	"errors"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"regexp"
)

var (
	tagRegexp = regexp.MustCompile("\\s*([^:\\|]+):([^:\\|]+)\\s*\\|*")
)

var _ promotion.ICashBackPromotion = new(CashBackPromotion)

type CashBackPromotion struct {
	*Promotion
	_cashBackValue *promotion.ValueCashBack
	_dataTag       map[string]string
}

// 获取领域编号
func (this *CashBackPromotion) GetDomainId() int {
	return this._cashBackValue.Id
}

// 设置详细的促销信息
func (this *CashBackPromotion) SetDetailsValue(v *promotion.ValueCashBack) error {
	g := this._goodsRep.GetValueGoodsById(this._value.GoodsId)
	if g == nil {
		return sale.ErrNoSuchGoods
	}

	//todo: 商品SKU的原因，获取的价格为0，有BUG

	//	if v.BackFee > int(g.SalePrice){
	//		return sale.ErrOutOfSalePrice
	//	}

	if len(v.DataTag) != 0 {
		if !tagRegexp.MatchString(v.DataTag) {
			return errors.New("自定义数据格式错误！正确的格式如：\\\"K:V | K:V\\\"")
		}
	}

	this._cashBackValue = v
	return nil
}

// 获取相关的值
func (this *CashBackPromotion) GetRelationValue() interface{} {
	return this._cashBackValue
}

// 促销类型
func (this *CashBackPromotion) TypeName() string {
	return "返现"
}

// 获取自定义数据
func (this *CashBackPromotion) GetDataTag() map[string]string {
	if this._dataTag == nil {
		this._dataTag = make(map[string]string)
		if len(this._cashBackValue.DataTag) != 0 {
			matches := tagRegexp.FindAllStringSubmatch(this._cashBackValue.DataTag, -1)
			for i := 0; i < len(matches); i++ {
				this._dataTag[matches[i][1]] = matches[i][2]
			}
		}
	}
	return this._dataTag
}

// 保存
func (this *CashBackPromotion) Save() (int, error) {

	if this.GetRelationValue() == nil {
		return this.GetAggregateRootId(), promotion.ErrCanNotApplied
	}

	var isCreate bool = this.GetAggregateRootId() == 0
	this._value.TypeFlag = promotion.TypeFlagCashBack
	id, err := this.Promotion.Save()
	if err == nil {
		this._value.Id = id
		if this._cashBackValue == nil {
			this._cashBackValue = new(promotion.ValueCashBack)
		}
		this._cashBackValue.Id = this.GetAggregateRootId()
		_, err = this._promRep.SaveValueCashBack(this._cashBackValue, isCreate)
	}
	return id, err
}
