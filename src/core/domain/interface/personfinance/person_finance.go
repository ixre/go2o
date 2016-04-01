/**
 * Copyright 2015 @ z3q.net.
 * name : person_finance
 * author : jarryliu
 * date : 2016-03-31 10:46
 * description :
 * history :
 */
package personfinance

import (
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/infrastructure/domain"
)

const (
	MinRiseTransferInAmount  float32 = 100.00 //最低转入金额为100
	MinRiseTransferOutAmount float32 = 0.00   //最低转出金额
)

type (
	// 在此聚合下, 会员抽象为Person, PersonId 对应 MemberId
	IPersonFinance interface {
		// 获取聚合根
		GetAggregateRootId() int
		// 获取账号
		GetMemberAccount() member.IAccount
		// 获取增利账户信息(类:余额宝)
		GetRiseInfo() IRiseInfo
		// 创建增利账户信息
		CreateRiseInfo() error
	}

	// 现金增利
	IRiseInfo interface {
		GetDomainId() int

		// 获取值
		Value() (RiseInfoValue, error)

		// 设置值
		//Set(*RiseInfoValue)error

		// 转入
		TransferIn(amount float32) error

		// 转出
		TransferOut(amount float32) error

		// 根据日志记录提交转入转出,如果已经确认操作,则返回错误
		// 通常是由系统计划任务来完成此操作,转入和转出必须经过提交!
		CommitTransfer(logId int) error

		// 结算增利信息,dayRatio 为每天的收益比率
		RiseSettleForToday(dayRatio float32) error

		// 获取时间段内的增利信息
		GetRiseByTime(begin, end int64) []*RiseDayInfo

		// 保存
		Save() error
	}

	// 收益总记录
	RiseInfoValue struct {
		//Id  int `db:"id" pk:"yes" auto:"no"`
		PersonId    int     `db:"person_id" pk:"yes" auto:"no"` //人员编号
		Balance     float32 `db:"balance"`                 //本金及收益的余额
		Rise       float32 `db:"rise"`                          //当前的收益
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
		UnixDate   int64   `db:"unix_date"`
		UpdateTime int64   `db:"update_time"`
	}

	// 收益日志
	RiseLog struct {
		Id         int     `db:"id" pk:"yes" auto:"yes"`
		PersonId   int     `db:"person_id"`   //会员编号
		Amount     float32 `db:"amount"`      //金额
		Type       int     `db:"type"`        //类型
		State      int     `db:"state"`       //状态
		UnixDate   int64   `db:"unix_date"`   // 日期
		LogTime    int64   `db:"log_time"`    //日志时间
		UpdateTime int64   `db:"update_time"` //更新时间
	}

	IPersonFinanceRepository interface {
		GetPersonFinance(personId int) IPersonFinance
		GetRiseByTime(personId int, begin, end int64) []*RiseDayInfo
		GetRiseValueByPersonId(id int) (v *RiseInfoValue, err error)
		SaveRiseInfo(*RiseInfoValue) (id int, err error)

		// 获取日志
		GetRiseLog(personId, logId int) *RiseLog

		// 保存日志
		SaveRiseLog(*RiseLog) (int, error)

		// 获取日志
		GetRiseLogs(personId int, date int64, riseType int) []*RiseLog
	}
)

var (
	RiseStateOk      int = 1 //OK
	RiseStateDefault int = 0 //默认状态

	ErrIncorrectAmount *domain.DomainError = domain.NewDomainError(
		"err_balance_amount", "金额错误!")
	ErrIncorrectTransfer *domain.DomainError = domain.NewDomainError(
		"err_incorrent_transfer", "已确认或非转入(转出)操作!")
	ErrNoSuchRiseInfo *domain.DomainError = domain.NewDomainError(
		"err_no_such_rise_info", "未开通该功能!")
	ErrHasSettled *domain.DomainError = domain.NewDomainError(
		"err_has_settled", "已经结算!")
	ErrRatio *domain.DomainError = domain.NewDomainError(
		"err_ratio", "利率不正确!")

	ErrLessThanMinTransferIn *domain.DomainError = domain.NewDomainError(
		"err_less_than_min_transfer_in", "转入金额最低%d!")

	ErrLessThanMinTransferOut *domain.DomainError = domain.NewDomainError(
		"err_less_than_min_transfer_out", "转出金额最低%d!")

	ErrOutOfBalance *domain.DomainError = domain.NewDomainError(
		"err_out_of_balance", "超出帐户最大金额!")
)

const (
	RiseTypeTransferIn  int = 1 + iota //转入
	RiseTypeTransferOut                //转出
	RiseTypeSettle                     //结算
	RiseTypeAdjust                     //人工调整
)
