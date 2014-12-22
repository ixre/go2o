/**
 * Copyright 2014 @ ops.
 * name :
 * author : newmin
 * date : 2013-11-11 14:43
 * description :
 * history :
 */

package entity

import (
	"time"
)

//页面传输需要的
type PartnerDto struct {
	Id      int       `db:"id" json:"id"`
	Secret  string    `db:"secret" json:"secret"`
	Name    string    `db:"name" json:"name"`
	Expires time.Time `db:"expires" json:"expires"`
	Tel     string    `db:"tel" json:"tel"`
	Phone   string    `db:"phone" json:"phone"`
	Address string    `db:"address" json:"address"`
}
