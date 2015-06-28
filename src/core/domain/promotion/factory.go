/**
 * Copyright 2015 @ S1N1 Team.
 * name : factor
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package promotion
import (
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/domain/interface/sale"
)

func FactoryPromotion(rep promotion.IPromotionRep,saleRep sale.ISaleRep,
	v *promotion.ValuePromotion)promotion.IPromotion{
	prom := newPromotion(rep,saleRep,v)

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