/**
 * Copyright 2015 @ S1N1 Team.
 * name : cash_back
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package promotion

// 返现促销
type ICashBackPromotion interface{
	// 设置详细的促销信息
	SetDetailsValue(*ValueCashBack)error
}