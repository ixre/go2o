/**
 * Copyright 2015 @ S1N1 Team.
 * name : cash_back
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package promotion
import "go2o/src/core/domain/interface/promotion"

var _ promotion.ICashBackPromotion = new(CashBackPromotion)
type CashBackPromotion struct {
	*Promotion
	_cashBackValue *promotion.ValueCashBack
}


// 设置详细的促销信息
func (this *CashBackPromotion) SetDetailsValue(v *promotion.ValueCashBack)error{
	this._cashBackValue = v
	return nil
}


// 获取相关的值
func (this *CashBackPromotion) GetRelationValue()interface{}{
	return this._cashBackValue
}

// 保存
func (this *CashBackPromotion) Save()(int,error){
	var isCreate bool = this.GetAggregateRootId() == 0
	this._value.TypeFlag = promotion.TypeFlagCashBack
	id,err := this.Promotion.Save()
	if err == nil {
		_,err = this._promRep.SaveValueCashBack(this._cashBackValue,isCreate)
	}
	return id,err
}