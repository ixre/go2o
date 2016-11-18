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
	Id           int64 `db:"id" auto:"yes" pk:"yes"`
	MerchantId   int64 `db:"merchant_id"`
	CoverageId   int64 `db:"coverage_id"`
	ShopId       int64 `db:"shop_id"`
	DeliverUsrId int64 `db:"delivery_usr_id"`
	Enabled      int   `db:"enabled"`
}
