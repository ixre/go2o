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
	Id           int `db:"id" auto:"yes" pk:"yes"`
	MerchantId   int `db:"merchant_id"`
	CoverageId   int `db:"coverage_id"`
	ShopId       int `db:"shop_id"`
	DeliverUsrId int `db:"delivery_usr_id"`
	Enabled      int `db:"enabled"`
}
