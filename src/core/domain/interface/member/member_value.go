/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 10:11
 * description :
 * history :
 */

package member

type ValueMember struct {
	Id  int    `db:"id" auto:"yes" pk:"yes"`
	Usr string `db:"usr"`
	Pwd string `db:"Pwd"`
	// 交易密码
	TradePwd string `db:"trade_pwd"`
	// 姓名
	Name string `db:"name"`
	// 经验值
	Exp int `db:"exp"`
	// 等级
	Level int `db:"level"`

	Sex      int    `db:"sex"`
	Avatar   string `db:"avatar"`
	BirthDay string `db:"birthday"`
	Phone    string `db:"phone"`
	Address  string `db:"address"`
	Im       string `db:"im"`
	Email    string `db:"email"`
	// 邀请码
	InvitationCode string `db:"invitation_code"`
	RegFrom        string `db:"reg_from"`
	RegIp          string `db:"reg_ip"`
	State          int    `db:"state"`
	RegTime        int64  `db:"reg_time"`
	Ext1           string `db:"ext_1"` // 扩展1
	Ext2           string `db:"ext_2"` // 扩展2
	Ext3           string `db:"ext_3"` // 扩展3
	Ext4           string `db:"ext_4"` // 扩展4
	Ext5           string `db:"ext_4"` // 扩展5
	Ext6           string `db:"ext_4"` // 扩展6
	LastLoginTime  int64  `db:"last_login_time"`
	UpdateTime     int64  `db:"update_time"`
	DynamicToken   string `db:"-"` // 动态令牌，用于登陆或API调用
	TimeoutTime    int64  `db:"-"` // 超时时间
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
