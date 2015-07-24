/**
 * Copyright 2015 @ S1N1 Team.
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
	OperateNo  string  `db:"operate_no"`
	Kind       int     `db:"kind"`
	Type       int     `db:"type"`
	Title      string  `db:"title"`
	Amount     float32 `db:"amount"`
	State      int     `db:"state"`
	UpdateTime int64   `db:"update_time"`
}
