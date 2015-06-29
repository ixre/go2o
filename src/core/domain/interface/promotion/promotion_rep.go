/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-03 14:11
 * description :
 * history :
 */

package promotion

type IPromotionRep interface {
	// 获取促销
	GetPromotion(id int) IPromotion

	// 获取促销
	CreatePromotion(*ValuePromotion) IPromotion

	// 获取促销
	GetValuePromotion(id int) *ValuePromotion

	// 保存促销
	SaveValuePromotion(*ValuePromotion) (int, error)

	// 删除促销
	DeletePromotion(id int)error

	// 保存返现促销
	SaveValueCashBack(v *ValueCashBack, create bool) (int, error)

	// 获取返现促销
	GetValueCashBack(int) *ValueCashBack

	// 删除现金返现促销
	DeleteValueCashBack(id int)error

	// 获取商品的促销编号
	GetGoodsPromotionId(goodsId int, promFlag int) int
}
