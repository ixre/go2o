/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */
package delivery

type MerchantDeliverBind struct {
	Id           int32 `db:"id" auto:"yes" pk:"yes"`
	MerchantId   int32 `db:"merchant_id"`
	CoverageId   int32 `db:"coverage_id"`
	ShopId       int32 `db:"shop_id"`
	DeliverUsrId int32 `db:"delivery_user_id"`
	Enabled      int   `db:"enabled"`
}
