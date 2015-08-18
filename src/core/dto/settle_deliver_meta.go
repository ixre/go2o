/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-09 17:53
 * description :
 * history :
 */
package dto

type SettleDeliverMeta struct {
	Id         int    `db:"id"`
	PersonName string `db:"personName"`
	Phone      string `db:"phone"`
	Address    string `db:"address"`
}
