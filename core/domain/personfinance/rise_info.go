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
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/personfinance"
	"go2o/core/infrastructure/format"
	"go2o/core/infrastructure/tool"
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
	if this._v == nil {
		//判断会员是否存在
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
		this._v.SettlementAmount += l.Amount // 计入结算金额
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
func (this *riseInfo) TransferIn(amount float32,
	w personfinance.TransferWith) (err error) {
	if this._v == nil {
		//判断会员是否存在
		if _, err = this.Value(); err != nil {
			return err
		}
	}

	if amount <= 0 {
		return personfinance.ErrIncorrectAmount
	}

	if amount < personfinance.RiseMinTransferInAmount {
		return errors.New(fmt.Sprintf(personfinance.ErrLessThanMinTransferIn.Error(),
			format.FormatFloat(personfinance.RiseMinTransferInAmount)))
	}

	dt := time.Now()
	this._v.TransferIn += amount
	this._v.TotalAmount += amount
	this._v.UpdateTime = dt.Unix()
	if err = this.Save(); err == nil {
		//保存并记录日志
		_, err = this._rep.SaveRiseLog(&personfinance.RiseLog{
			PersonId:     this.GetDomainId(),
			Title:        "[转入]从" + personfinance.TransferInWithText(w) + "转入",
			Amount:       amount,
			Type:         personfinance.RiseTypeTransferIn,
			TransferWith: int(w),
			State:        personfinance.RiseStateDefault,
			UnixDate:     tool.GetStartDate(dt).Unix(),
			LogTime:      this._v.UpdateTime,
			UpdateTime:   this._v.UpdateTime,
		})
	}
	return err
}

// 转出,w为转出方式(如银行,余额等),state为日志的状态,某些操作
// 需要确认,有些不需要.通过state来传入
func (this *riseInfo) TransferOut(amount float32,
	w personfinance.TransferWith, state int) (err error) {
	if this._v == nil {
		//判断会员是否存在
		if _, err = this.Value(); err != nil {
			return err
		}
	}
	if amount <= 0 {
		return personfinance.ErrIncorrectAmount
	}

	if amount > this._v.Balance {
		//超出账户金额
		return personfinance.ErrOutOfBalance
	}

	// 低于最低转出金额,且不是全部转出.返回错误. 若转出到余额则无限制
	if amount != this._v.Balance && //非全部转出
		w != personfinance.TransferOutWithBalance && //非转出余额
		amount < personfinance.RiseMinTransferOutAmount {
		if this._v.Balance > personfinance.RiseMinTransferOutAmount {
			//金额大于转出金额
			return errors.New(fmt.Sprintf(personfinance.ErrLessThanMinTransferOut.Error(),
				format.FormatFloat(personfinance.RiseMinTransferOutAmount)))
		} else {
			//金额小于转出金额
			return errors.New(fmt.Sprintf(personfinance.ErrMustAllTransferOut.Error(),
				format.FormatFloat(personfinance.RiseMinTransferOutAmount)))
		}
	}

	dt := time.Now()
	this._v.UpdateTime = dt.Unix()
	this._v.Balance -= amount
	this._v.SettlementAmount -= amount

	if this._v.SettlementAmount < 0 {
		//提现超出结算金额
		this._v.SettlementAmount = 0
	}

	if this._v.Balance == 0 {
		//若全部提出,则理财收益和结算金额清零
		this._v.Rise = 0
		this._v.SettlementAmount = 0
	}

	if err = this.Save(); err == nil {
		//保存并记录日志
		_, err = this._rep.SaveRiseLog(&personfinance.RiseLog{
			PersonId:     this.GetDomainId(),
			Title:        "[转出]转出到" + personfinance.TransferOutWithText(w),
			Amount:       amount,
			Type:         personfinance.RiseTypeTransferOut,
			TransferWith: int(w),
			State:        state,
			UnixDate:     tool.GetStartDate(dt).Unix(),
			LogTime:      this._v.UpdateTime,
			UpdateTime:   this._v.UpdateTime,
		})
	}
	//todo: 新增操作记录,如审核,打款,完成等
	return err
}

// 结算收益(按天结息),settleUnix:结算日期的时间戳(不含时间),
// dayRatio 为每天的收益比率
func (this *riseInfo) RiseSettleByDay(settleDateUnix int64, dayRatio float32) (err error) {
	if this._v == nil {
		//判断会员是否存在
		if _, err = this.Value(); err != nil {
			return err
		}
	}

	if dayRatio < 0 {
		return personfinance.ErrRatio
	}

	//dt := time.Now().Add(time.Hour * -24) //计算昨日的收益
	//dtUnix := tool.GetStartDate(dt).Unix()

	if settleDateUnix%100 != 0 {
		return personfinance.ErrUnixDate
	}

	settleDateStr := time.Unix(settleDateUnix, 0).Format("2006-01-02") //结算日期年月日

	if b, err := this.daySettled(settleDateUnix); b {
		if err != nil {
			return err
		}
		return personfinance.ErrHasSettled
	}

	if this._v.SettlementAmount > 0 {
		amount := float32(format.FixedDecimal(float64(
			this._v.SettlementAmount * dayRatio))) //按2位小数精度

		// 有可能出现金额太小,收益为0的情况,这时应标记结算日期为最新;
		// 但不增加收益日志
		if amount > 0.00 {
			if _, err = this.monthSettle(this._v, settleDateUnix); err != nil {
				return err
			}

			this._v.Balance += amount
			this._v.Rise += amount
			this._v.TotalRise += amount
		}
		this._v.SettledDate = settleDateUnix //结算日为昨日
		this._v.UpdateTime = time.Now().Unix()
		err = this.Save()
		if err == nil && amount > 0.00 {
			// 保存计息日志
			_, err = this._rep.SaveRiseLog(&personfinance.RiseLog{
				PersonId:   this.GetDomainId(),
				Title:      "收益",
				Amount:     amount,
				Type:       personfinance.RiseTypeGenerateInterest,
				State:      personfinance.RiseStateOk,
				UnixDate:   settleDateUnix,
				LogTime:    this._v.UpdateTime,
				UpdateTime: this._v.UpdateTime,
			})
			// 存储每日收益
			_, err = this._rep.SaveRiseDayInfo(&personfinance.RiseDayInfo{
				PersonId:         this.GetDomainId(),
				Date:             settleDateStr,
				SettlementAmount: this._v.SettlementAmount,
				RiseAmount:       amount,
				UnixDate:         settleDateUnix,
				UpdateTime:       this._v.UpdateTime,
			})
		}
	}
	return err
}

// 月结
func (this *riseInfo) monthSettle(v *personfinance.RiseInfoValue, settleDateUnix int64) (settled bool, err error) {
	if settleDateUnix%100 != 0 {
		return false, personfinance.ErrUnixDate
	}
	y, m, d := time.Unix(settleDateUnix, 0).Date()
	isBonusDay := tool.LastDay(y, m) == d // 是否为分红结算月结日期
	if isBonusDay {
		//月结分红时,将余额作为新的投资金额
		//分红 = 余额 - 结算结算金额
		bonusAmount := this._v.Balance - this._v.SettlementAmount //分红金额
		if bonusAmount > 0 {
			this._v.SettlementAmount = this._v.Balance
			this._v.UpdateTime = time.Now().Unix()
			err = this.Save()
			if err == nil {
				// 保存分红再投资日志
				_, err = this._rep.SaveRiseLog(&personfinance.RiseLog{
					PersonId:   this.GetDomainId(),
					Title:      "[月结]红利再投资",
					Amount:     bonusAmount,
					Type:       personfinance.RiseTypeMonthSettle,
					State:      personfinance.RiseStateOk,
					UnixDate:   settleDateUnix,
					LogTime:    this._v.UpdateTime,
					UpdateTime: this._v.UpdateTime,
				})
			}
		}
		return true, err
	}
	return false, err
}

// 是否已经结算
func (this *riseInfo) daySettled(dateUnix int64) (bool, error) {
	if dateUnix%100 != 0 {
		return true, personfinance.ErrUnixDate
	}
	arr := this._rep.GetRiseLogs(this.GetDomainId(), dateUnix,
		personfinance.RiseTypeGenerateInterest)
	return len(arr) > 0, nil
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
