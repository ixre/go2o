/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-12-09 10:11
 * description :
 * history :
 */

package member

type ValueMember struct {
	Id   int    `db:"id" auto:"yes" pk:"yes"`
	Usr  string `db:"usr"`
	Pwd  string `db:"Pwd"`
	Name string `db:"name"`
	// 经验值
	Exp int `db:"exp"`
	// 等级
	Level int `db:"level"`

	Sex      int    `db:"sex"`
	Avatar   string `db:"avatar"`
	Birthday string `db:"birthday"`
	Phone    string `db:"phone"`
	Address  string `db:"address"`
	Qq       string `db:"qq"`
	Email    string `db:"email"`
	RegIp    string `db:"reg_ip"`
	State    int    `db:"state"`
	//登录密钥
	LoginToken string `db:"-"`

	RegTime       int64 `db:"reg_time"`
	LastLoginTime int64 `db:"last_login_time"`
}

type BankInfo struct {
	MemberId    int    `db:"member_id" pk:"yes"`
	Name        string `db:"name"`
	Account     string `db:"account"`
	AccountName string `db:"account_name"`
	Network     string `db:"network"`
	State       int    `db:"state"`
	UpdateTime  int64  `db:"update_time"`
}
