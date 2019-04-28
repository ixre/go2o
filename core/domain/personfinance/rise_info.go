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
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"go2o/core/infrastructure/tool"
	"math"
	"time"
)

var _ personfinance.IRiseInfo = new(riseInfo)

type riseInfo struct {
	personId int64
	value    *personfinance.RiseInfoValue
	rep      personfinance.IPersonFinanceRepository
	mmRepo   member.IMemberRepo
	pf       *PersonFinance
}

func newRiseInfo(personId int64, pf *PersonFinance,
	rep personfinance.IPersonFinanceRepository,
	mmRepo member.IMemberRepo) personfinance.IRiseInfo {
	return &riseInfo{
		personId: personId,
		rep:      rep,
		mmRepo:   mmRepo,
		pf:       pf,
	}
}

func (r *riseInfo) GetDomainId() int64 {
	return r.personId
}

// 根据日志记录提交转入转出,如果已经确认操作,则返回错误
// 通常是由系统计划任务来完成此操作,转入和转出必须经过提交!
func (r *riseInfo) CommitTransfer(logId int32) (err error) {
	if r.value == nil {
		//判断会员是否存在
		if _, err = r.Value(); err != nil {
			return err
		}
	}
	l := r.rep.GetRiseLog(r.GetDomainId(), logId)
	if l == nil || l.State != personfinance.RiseStateDefault || ( // 状态应用未确认
		l.Type != personfinance.RiseTypeTransferIn && // 类型应为转入或转出
			l.Type != personfinance.RiseTypeTransferOut) {
		return personfinance.ErrIncorrectTransfer
	}
	if r.value.TransferIn < l.Amount {
		return personfinance.ErrIncorrectAmount
	}
	switch l.Type {
	case personfinance.RiseTypeTransferIn:
		r.value.Balance += l.Amount
		r.value.SettlementAmount += l.Amount // 计入结算金额
		r.value.TransferIn -= l.Amount
		err = r.Save()
		//todo: 记录开使计算收益的日志
	case personfinance.RiseTypeTransferOut:
		//todo: 处理打款, 转出成功日志
	}
	l.State = personfinance.RiseStateOk
	l.UpdateTime = time.Now().Unix()
	if err == nil {
		_, err = r.rep.SaveRiseLog(l)
	}
	return err
}

// 转入扣款
func (r *riseInfo) transferInPayment(amount float32,
	transferWith personfinance.TransferWith) (err error) {
	m := r.mmRepo.GetMember(r.personId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	acc := m.GetAccount()
	tradeNo := domain.NewTradeNo(8, int(r.value.PersonId))
	switch transferWith {
	//从余额转入
	case personfinance.TransferFromWithBalance:
		err = acc.DiscountBalance("理财转入", tradeNo,
			amount, member.DefaultRelateUser)
		if err != nil {
			return err
		}
		//从钱包转入
	case personfinance.TransferFromWithWallet:
		err = acc.DiscountWallet("理财转入", tradeNo,
			amount, member.DefaultRelateUser, true)
		if err != nil {
			return err
		}
		//其他方式转入
	default:
		return errors.New("暂时无法提供服务")
	}
	return nil
}

// 转入
func (r *riseInfo) TransferIn(amount float32,
	transferWith personfinance.TransferWith) (err error) {
	if r.value == nil {
		//判断会员是否存在
		if _, err = r.Value(); err != nil {
			return err
		}
	}
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return personfinance.ErrIncorrectAmount
	}
	if amount < personfinance.RiseMinTransferInAmount {
		return errors.New(fmt.Sprintf(personfinance.ErrLessThanMinTransferIn.Error(),
			format.FormatFloat(personfinance.RiseMinTransferInAmount)))
	}

	if amount < personfinance.RiseMinTransferInAmount {
		//金额不足最低转入金额
		return errors.New(fmt.Sprintf(personfinance.ErrLessThanMinTransferIn.Error(),
			format.FormatFloat(personfinance.RiseMinTransferInAmount)))
	}
	err = r.transferInPayment(amount, transferWith)
	if err == nil {
		r.value.TransferIn += amount
		r.value.TotalAmount += amount
		dt := time.Now()
		r.value.UpdateTime = dt.Unix()
		if err = r.Save(); err == nil {
			//保存并记录日志
			_, err = r.rep.SaveRiseLog(&personfinance.RiseLog{
				PersonId:     r.GetDomainId(),
				Title:        "[转入]从" + personfinance.TransferInWithText(transferWith) + "转入",
				Amount:       amount,
				Type:         personfinance.RiseTypeTransferIn,
				TransferWith: int(transferWith),
				State:        personfinance.RiseStateDefault,
				UnixDate:     tool.GetStartDate(dt).Unix(),
				LogTime:      r.value.UpdateTime,
				UpdateTime:   r.value.UpdateTime,
			})
			if err == nil {
				return r.pf.SyncToAccount() //同步到会员账户
			}
		}
	}
	return err
}

// 转出,w为转出方式(如银行,余额等),state为日志的状态,某些操作
// 需要确认,有些不需要.通过state来传入
func (r *riseInfo) TransferOut(amount float32,
	w personfinance.TransferWith, state int) (err error) {
	if r.value == nil {
		//判断会员是否存在
		if _, err = r.Value(); err != nil {
			return err
		}
	}
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return personfinance.ErrIncorrectAmount
	}

	if amount > r.value.Balance {
		//超出账户金额
		return personfinance.ErrOutOfBalance
	}

	// 低于最低转出金额,且不是全部转出.返回错误. 若转出到余额则无限制
	if amount != r.value.Balance && //非全部转出
		w != personfinance.TransferOutWithBalance && //非转出余额
		amount < personfinance.RiseMinTransferOutAmount {
		if r.value.Balance > personfinance.RiseMinTransferOutAmount {
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
	r.value.UpdateTime = dt.Unix()
	r.value.Balance -= amount
	r.value.SettlementAmount -= amount

	if r.value.SettlementAmount < 0 {
		//提现超出结算金额
		r.value.SettlementAmount = 0
	}

	if r.value.Balance == 0 {
		//若全部提出,则理财收益和结算金额清零
		r.value.Rise = 0
		r.value.SettlementAmount = 0
	}

	if err = r.Save(); err == nil {
		//保存并记录日志
		_, err = r.rep.SaveRiseLog(&personfinance.RiseLog{
			PersonId:     r.GetDomainId(),
			Title:        "[转出]转出到" + personfinance.TransferOutWithText(w),
			Amount:       amount,
			Type:         personfinance.RiseTypeTransferOut,
			TransferWith: int(w),
			State:        state,
			UnixDate:     tool.GetStartDate(dt).Unix(),
			LogTime:      r.value.UpdateTime,
			UpdateTime:   r.value.UpdateTime,
		})
	}
	//todo: 新增操作记录,如审核,打款,完成等
	return err
}

// 结算收益(按天结息),settleUnix:结算日期的时间戳(不含时间),
// dayRatio 为每天的收益比率
func (r *riseInfo) RiseSettleByDay(settleDateUnix int64, dayRatio float32) error {
	_, err := r.Value()
	if err != nil {
		return err
	}
	// 错误的收益率
	if dayRatio < 0 {
		return personfinance.ErrRatio
	}
	// 不是日期
	if settleDateUnix%100 != 0 {
		return personfinance.ErrUnixDate
	}
	// 结算日期年月日
	settleDateStr := time.Unix(settleDateUnix, 0).Format("2006-01-02")
	// 判断是否已结算
	if b, err := r.daySettled(settleDateUnix); b {
		if err != nil {
			return err
		}
		return personfinance.ErrHasSettled
	}
	// 开始结算
	if r.value.SettlementAmount > 0 {
		//按2位小数精度
		amount := float32(format.FixedDecimal(float64(
			r.value.SettlementAmount * dayRatio)))

		// 有可能出现金额太小,收益为0的情况,这时应标记结算日期为最新;
		// 但不增加收益日志
		if amount > 0.00 {
			_, err = r.monthSettle(r.value, settleDateUnix)
			if err != nil {
				return err
			}
			r.value.Balance += amount
			r.value.Rise += amount
			r.value.TotalRise += amount
		}
		r.value.SettledDate = settleDateUnix //结算日为昨日
		r.value.UpdateTime = time.Now().Unix()
		err = r.Save()
		if err == nil && amount > 0.00 {
			// 保存计息日志
			_, err = r.rep.SaveRiseLog(&personfinance.RiseLog{
				PersonId:   r.GetDomainId(),
				Title:      "收益",
				Amount:     amount,
				Type:       personfinance.RiseTypeGenerateInterest,
				State:      personfinance.RiseStateOk,
				UnixDate:   settleDateUnix,
				LogTime:    r.value.UpdateTime,
				UpdateTime: r.value.UpdateTime,
			})
			// 存储每日收益
			_, err = r.rep.SaveRiseDayInfo(&personfinance.RiseDayInfo{
				PersonId:         r.GetDomainId(),
				Date:             settleDateStr,
				SettlementAmount: r.value.SettlementAmount,
				RiseAmount:       amount,
				UnixDate:         settleDateUnix,
				UpdateTime:       r.value.UpdateTime,
			})
		}
	}
	return err
}

// 月结
func (r *riseInfo) monthSettle(v *personfinance.RiseInfoValue,
	settleDateUnix int64) (settled bool, err error) {
	if settleDateUnix%100 != 0 {
		return false, personfinance.ErrUnixDate
	}
	y, m, d := time.Unix(settleDateUnix, 0).Date()
	isBonusDay := tool.LastDay(y, m) == d // 是否为分红结算月结日期
	if isBonusDay {
		//月结分红时,将余额作为新的投资金额
		//分红 = 余额 - 结算结算金额
		bonusAmount := r.value.Balance - r.value.SettlementAmount //分红金额
		if bonusAmount > 0 {
			r.value.SettlementAmount = r.value.Balance
			r.value.UpdateTime = time.Now().Unix()
			err = r.Save()
			if err == nil {
				// 保存分红再投资日志
				_, err = r.rep.SaveRiseLog(&personfinance.RiseLog{
					PersonId:   r.GetDomainId(),
					Title:      "[月结]红利再投资",
					Amount:     bonusAmount,
					Type:       personfinance.RiseTypeMonthSettle,
					State:      personfinance.RiseStateOk,
					UnixDate:   settleDateUnix,
					LogTime:    r.value.UpdateTime,
					UpdateTime: r.value.UpdateTime,
				})
			}
		}
		return true, err
	}
	return false, err
}

// 是否已经结算
func (r *riseInfo) daySettled(dateUnix int64) (bool, error) {
	if dateUnix%100 != 0 {
		return true, personfinance.ErrUnixDate
	}
	arr := r.rep.GetRiseLogs(r.GetDomainId(), dateUnix,
		personfinance.RiseTypeGenerateInterest)
	return len(arr) > 0, nil
}

// 获取时间段内的增利信息
func (r *riseInfo) GetRiseByTime(begin, end int64) []*personfinance.RiseDayInfo {
	return r.rep.GetRiseByTime(r.GetDomainId(), begin, end)
}

// 获取值
func (r *riseInfo) Value() (personfinance.RiseInfoValue, error) {
	if r.value == nil {
		if v, err := r.rep.GetRiseValueByPersonId(r.GetDomainId()); err != nil {
			return personfinance.RiseInfoValue{}, err
		} else {
			r.value = v
		}
	}
	return *r.value, nil
}

// 保存
func (r *riseInfo) Save() error {
	_, err := r.rep.SaveRiseInfo(r.value)
	return err
}
