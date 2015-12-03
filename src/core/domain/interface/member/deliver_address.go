/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-11-25 21:43
 * description :
 * history :
 */

package member

// 会员收货地址
type DeliverAddress struct {
	Id        int    `db:"id" pk:"yes" auto:"yes"`
	MemberId  int    `db:"member_id"`
	RealName  string `db:"real_name"`
	Phone     string `db:"phone"`
	Address   string `db:"address"`
	IsDefault int    `db:"is_default"`
}
