/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-01-09 21:42
 * description :
 * history :
 */

package shopping

type OrderLog struct {
	//Id int `db:"id" auto:"yes" pk:"yes"`
	OrderId    int    `db:"order_id"`
	Type       int    `db:"type"`
	IsSystem   int    `db:"is_system"`
	Message    string `db:"message"`
	RecordTime int64  `db:"record_time"`
}
