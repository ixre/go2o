/**
 * Copyright 2015 @ z3q.net.
 * name : rise_info
 * author : jarryliu
 * date : 2016-03-31 18:06
 * description :
 * history :
 */
package personfinance

import (
	"errors"
	"fmt"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/personfinance"
	"go2o/src/core/infrastructure/tool"
	"time"
)

var _ personfinance.IRiseInfo = new(riseInfo)

type riseInfo struct {
	_personId int
	_v        *personfinance.RiseInfoValue
	_rep      personfinance.IPersonFinanceRepository
	_accRep   member.IMemberRep
}

func newRiseInfo(personId int, rep personfinance.IPersonFinanceRepository,
	accRep member.IMemberRep) personfinance.IRiseInfo {
	return &riseInfo{
		_personId: personId,
		_rep:      rep,
		_accRep:   accRep,
	}
}

func (this *riseInfo) GetDomainId() int {
	return this._personId
}

// 根据日志记录提交转入转出,如果已经确认操作,则返回错误
// 通常是由系统计划任务来完成此操作,转入和转出必须经过提交!
func (this *riseInfo) CommitTransfer(logId int) (err error) {
	if this._v == nil { //判断会员是否存在
		if _, err = this.Value(); err != nil {
			return err
		}
	}
	l := this._rep.GetRiseLog(this.GetDomainId(), logId)
	if l == nil || l.State != personfinance.RiseStateDefault || ( // 状态应用未确认
	l.Type != personfinance.RiseTypeTransferIn &&                 // 类型应为转入或转出
		l.Type != personfinance.RiseTypeTransferOut) {
		return personfinance.ErrIncorrectTransfer
	}
	if this._v.TransferIn < l.Amount {
		return personfinance.ErrIncorrectAmount
	}
	switch l.Type {
	case personfinance.RiseTypeTransferIn:
		this._v.Balance += l.Amount
		this._v.TransferIn -= l.Amount
		err = this.Save()
		//todo: 记录开使计算收益的日志
	case personfinance.RiseTypeTransferOut:
		//todo: 处理打款, 转出成功日志
	}
	l.State = personfinance.RiseStateOk
	l.UpdateTime = time.Now().Unix()
	if err == nil {
		_, err = this._rep.SaveRiseLog(l)
	}
	return err
}

// 转入
func (this *riseInfo) TransferIn(amount float32) (err error) {
	if this._v == nil { //判断会员是否存在
		if _, err = this.Value(); err != nil {
			return err
		}
	}
	if amount < personfinance.RiseMinTransferInAmount {
		return errors.New(fmt.Sprintf(personfinance.ErrLessThanMinTransferIn.Error(),
			personfinance.RiseMinTransferInAmount))
	}

	dt := time.Now()
	this._v.TransferIn += amount
	this._v.TotalAmount += amount
	this._v.UpdateTime = dt.Unix()
	if err = this.Save(); err == nil { //保存并记录日志
		_, err = this._rep.SaveRiseLog(&personfinance.RiseLog{
			PersonId:   this.GetDomainId(),
			Amount:     amount,
			Type:       personfinance.RiseTypeTransferIn,
			State:      personfinance.RiseStateDefault,
			UnixDate:   tool.GetStartDate(dt).Unix(),
			LogTime:    this._v.UpdateTime,
			UpdateTime: this._v.UpdateTime,
		})
	}
	return err
}

// 转出
func (this *riseInfo) TransferOut(amount float32) (err error) {
	if this._v == nil { //判断会员是否存在
		if _, err = this.Value(); err != nil {
			return err
		}
	}
	if amount > this._v.Balance { //超出账户金额
		return personfinance.ErrOutOfBalance
	}

	if amount != this._v.Balance && amount < personfinance.RiseMinTransferOutAmount {
		return errors.New(fmt.Sprintf(personfinance.ErrLessThanMinTransferOut.Error(),
			personfinance.RiseMinTransferOutAmount))
	}

	dt := time.Now()
	this._v.UpdateTime = dt.Unix()
	this._v.Balance -= amount
	if this._v.Balance == 0 { //若全部提出,则理财金额清0
		this._v.Rise = 0
	}
	if err = this.Save(); err == nil { //保存并记录日志
		_, err = this._rep.SaveRiseLog(&personfinance.RiseLog{
			PersonId:   this.GetDomainId(),
			Amount:     amount,
			Type:       personfinance.RiseTypeTransferOut,
			State:      personfinance.RiseStateDefault,
			UnixDate:   tool.GetStartDate(dt).Unix(),
			LogTime:    this._v.UpdateTime,
			UpdateTime: this._v.UpdateTime,
		})
	}
	//todo: 新增操作记录,如审核,打款,完成等
	return err
}

// 结算增利信息(前一日)
func (this *riseInfo) RiseSettleForToday(dayRatio float32) (err error) {
	if this._v == nil { //判断会员是否存在
		if _, err = this.Value(); err != nil {
			return err
		}
	}

	if dayRatio < 0 {
		return personfinance.ErrRatio
	}

	dt := time.Now().Add(time.Hour * -24) //计算昨日的收益
	if this.daySettled(dt) {
		return personfinance.ErrHasSettled
	}

	if this._v.Balance > 0 {
		amount := this._v.Balance * dayRatio
		if amount > 0.01 {
			this._v.Balance += amount
			this._v.Rise += amount
			this._v.TotalRise += amount
			err = this.Save()
			if err == nil {
				_, err = this._rep.SaveRiseLog(&personfinance.RiseLog{
					PersonId:   this.GetDomainId(),
					Amount:     amount,
					Type:       personfinance.RiseTypeSettle,
					State:      personfinance.RiseStateOk,
					UnixDate:   tool.GetStartDate(dt).Unix(),
					LogTime:    this._v.UpdateTime,
					UpdateTime: this._v.UpdateTime,
				})
			}
		}
	}
	return err
}

// 是否已经结算
func (this *riseInfo) daySettled(t time.Time) bool {
	unix := tool.GetStartDate(t).Unix()
	arr := this._rep.GetRiseLogs(this.GetDomainId(), unix,
		personfinance.RiseTypeSettle)
	return len(arr) > 0
}

// 获取时间段内的增利信息
func (this *riseInfo) GetRiseByTime(begin, end int64) []*personfinance.RiseDayInfo {
	return this._rep.GetRiseByTime(this.GetDomainId(), begin, end)
}

// 获取值
func (this *riseInfo) Value() (personfinance.RiseInfoValue, error) {
	if this._v == nil {
		if v, err := this._rep.GetRiseValueByPersonId(this.GetDomainId()); err != nil {
			return personfinance.RiseInfoValue{}, personfinance.ErrNoSuchRiseInfo
		} else {
			this._v = v
		}
	}
	return *this._v, nil
}

// 保存
func (this *riseInfo) Save() error {
	_, err := this._rep.SaveRiseInfo(this._v)
	return err
}
