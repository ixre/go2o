/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */
package delivery

type PartnerDeliverBind struct {
	Id           int `db:"id" auto:"yes" pk:"yes"`
	PartnerId    int `db:"partner_id"`
	CoverageId   int `db:"coverage_id"`
	ShopId       int `db:"shop_id"`
	DeliverUsrId int `db:"delivery_usr_id"`
	Enabled      int `db:"enabled"`
}
