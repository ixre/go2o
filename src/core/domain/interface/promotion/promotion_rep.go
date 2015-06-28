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
	// 保存促销
	SaveValuePromotion(*ValuePromotion)(int,error)
}
