/**
 * Copyright 2014 @ ops.
 * name :
 * author : jarryliu
 * date : 2013-11-13 21:08
 * description :
 * history :
 */
package member

type IncomeLog struct {
	Id         int     `db:"id" pk:"yes" auto:"yes"`
	OrderId    int     `db:"order_id"`
	MemberId   int     `db:"member_id"`
	Type       string  `db:"type"`
	Fee        float32 `db:"fee"`
	Log        string  `db:"log"`
	State      int     `db:"state"`
	RecordTime int64   `db:"record_time"`
}
