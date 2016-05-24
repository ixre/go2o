/**
 * Copyright 2014 @ ops.
 * name :
 * author : jarryliu
 * date : 2013-11-11 20:42
 * description :
 * history :
 */

package member

type IntegralLog struct {
	Id         int    `db:"id" pk:"yes" auto:"yes"`
	MerchantId int    `db:"merchant_id"`
	MemberId   int    `db:"member_id"`
	Type       int    `db:"type"`
	Integral   int    `db:"integral"`
	Log        string `db:"log"`
	RecordTime int64  `db:"record_time"`
}
