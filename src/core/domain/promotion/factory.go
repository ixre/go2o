/**
 * Copyright 2015 @ S1N1 Team.
 * name : factor
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package promotion
import "go2o/src/core/domain/interface/promotion"

func FactoryPromotion(rep promotion.IPromotionRep,v *promotion.ValuePromotion)promotion.IPromotion{
	prom := newPromotion(rep,v)

	if prom.Type() == promotion.TypeFlagCashBack{
		pv := rep.GetValueCashBack(prom.GetAggregateRootId())
		if pv == nil{
			pv = &promotion.ValueCashBack{
				Id:prom.GetAggregateRootId(),
			}
		}
		cp := &CashBackPromotion{
			Promotion:prom,
			_cashBackValue:pv,
		}
		return cp
	}

	//todo:
	return prom
}