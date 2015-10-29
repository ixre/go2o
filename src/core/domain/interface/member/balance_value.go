/**
 * Copyright 2015 @ z3q.net.
 * name : charge_log
 * author : jarryliu
 * date : 2015-07-24 09:19
 * description :
 * history :
 */
package member

// 余额变动信息
type BalanceInfoValue struct {
	Id         int     `db:"id" auto:"yes" pk:"yes"`
	MemberId   int     `db:"member_id"`
	TradeNo    string  `db:"trade_no"`
	Kind       int     `db:"kind"`
	Type       int     `db:"type"`
	Title      string  `db:"title"`
	Amount     float32 `db:"amount"` 		// 金额
	CsnAmount  float32 `db:"csn_amount"` 	// 手续费
	RefId      int     `db:"ref_id"` 		// 引用编号
	State      int     `db:"state"`
	CreateTime int64   `db:"create_time"`
	UpdateTime int64   `db:"update_time"`
}
