/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-12-22 17:59
 * description :
 * history :
 */

package partner

import (
	"time"
)

//门店
type ValueShop struct {
	Id         int       `db:"id" pk:"yes" auto:"yes"`
	PartnerId  int       `db:"pt_id"`
	Name       string    `db:"name"`
	Address    string    `db:"address"`
	Phone      string    `db:"phone"`
	OrderIndex int       `db:"order_index"`
	State      int       `db:"state"`
	CreateTime time.Time `db:"create_time"`
}
