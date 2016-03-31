/**
 * Copyright 2015 @ z3q.net.
 * name : person_finance
 * author : jarryliu
 * date : 2016-03-31 10:46
 * description :
 * history :
 */
package personfinance

import "go2o/src/core/domain/interface/member"

type (
	// 在此聚合下, 会员抽象为Person, PersonId 对应 MemberId
	IPersonFinance interface {
		// 获取聚合根
		GetAggregateRootId() int

		// 获取账号
		GetMemberAccount() *member.IAccount

		// 转入
		TransferIn(amount float32) error

		// 转出
		TransferOut(amount float32) error

		// 获取增利账户信息(类:余额宝)
		GetRiseInfo() *RiseInfo

		// 结算增利信息
		RiseSettleForToday() error

		// 获取时间段内的增利信息
		GetRiseByTime(begin, end int64) []*RiseDayInfo
	}

	// 收益总记录
	RiseInfo struct {
		//Id  int `db:"id" pk:"yes" auto:"no"`
		PersonId    int     `db:"person_id" pk:"yes" auto:"no"` //人员编号
		Balance     float32 `db:"base_balance"`                 //本金及收益的余额
		TransferIn  float32 `db:"transfer_in"`                  //今日转入
		TotalAmount float32 `db:"total_amount"`                 //总金额
		TotalRise   float32 `db:"total_rise"`                   //总收益
		UpdateTime  int64   `db:"update_time"`
	}

	// 收益每日结算数据
	RiseDayInfo struct {
		Id         int     `db:"id" pk:"yes" auto:"yes"`
		PersonId   int     `db:"person_id"`
		Date       string  `db:"date"`
		BaseAmount float32 `db:"base_amount"` //本金
		RiseAmount string  `db:"rise_amount"` //增加金额
		IntDate    int64   `db:"unix_date"`
		UpdateTime int64   `db:"update_time"`
	}
)
