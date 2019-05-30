/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 16:30
 * description :
 * history :
 */
package user

// 用户凭据
type CredentialValue struct {
	Id       int32  `db:"id" auto:"yes" pk:"yes"`
	Usr      string `db:"user"`
	Pwd      string `db:"pwd"`
	PersonId int    `db:"person_id"`
	// 标记凭据的类型
	Sign    string `db:"sign"`
	Enabled int    `db:"enabled"`
}
