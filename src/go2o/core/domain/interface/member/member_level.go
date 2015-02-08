/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-09 10:28
 * description :
 * history :
 */

package member

type MemberLevel struct {
	Id int `db:"id" auto:"yes"`
	// 等级值
	Value int `db:"value" pk:"yes"`

	Name       string `db:"name"`
	RequireExp int    `db:"require_exp"`
	Enabled    int    `db:"enabled"`
}
