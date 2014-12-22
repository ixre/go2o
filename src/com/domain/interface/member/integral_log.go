/**
 * Copyright 2014 @ ops.
 * name :
 * author : newmin
 * date : 2013-11-11 20:42
 * description :
 * history :
 */

package member

import (
	"time"
)

type IntegralLog struct {
	Id         int       `db:"id" pk:"yes" auto:"yes"`
	PtId       int       `db:"pt_id"`
	MemberId   int       `db:"member_id"`
	Type       int       `db:"type"`
	Integral   int       `db:"integral"`
	Log        string    `db:"log"`
	RecordTime time.Time `db:"record_time"`
}
