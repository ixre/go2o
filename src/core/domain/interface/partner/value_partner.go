/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:52
 * description :
 * history :
 */

package partner

//合作商
type ValuePartner struct {
	Id            int    `db:"id" pk:"yes" auto:"yes"`
	Usr           string `db:"usr"`
	Pwd           string `db:"pwd"`
	Name          string `db:"name"`
	Logo          string `db:"logo"`
	Tel           string `db:"tel"`
	Phone         string `db:"phone"`
	Address       string `db:"address"`
	ExpiresTime   int64  `db:"expires_time"`
	JoinTime      int64  `db:"join_time"`
	UpdateTime    int64  `db:"update_time"`
	LoginTime     int64  `db:"login_time"`
	LastLoginTime int64  `db:"last_login_time"`
}
