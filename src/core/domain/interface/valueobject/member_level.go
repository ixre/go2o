/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 10:28
 * description :
 * history :
 */

package valueobject

type MemberLevel struct {
	Id int `db:"id" auto:"yes" pk:"yes"`

	PartnerId int `db:"partner_id"`

	// 等级值(1,2,4,8,16)
	Value int `db:"value" `

	Name       string `db:"name"`
	RequireExp int    `db:"require_exp"`
	Enabled    int    `db:"enabled"`
}
