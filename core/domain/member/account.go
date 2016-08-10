/**
 * Copyright 2015 @ z3q.net.
 * name : account
 * author : jarryliu
 * date : 2015-07-24 08:50
 * description :
 * history :
 */
package member

import (
	"errors"
	dm "go2o/core/domain"
	"go2o/core/domain/interface/member"
	"go2o/core/infrastructure/domain"
	"time"
)

var _ member.IAccount = new(accountImpl)

type accountImpl struct {
	_value *member.Account
	_rep   member.IMemberRep
}

func NewAccount(value *member.Account,
	rep member.IMemberRep) member.IAccount {
	return &accountImpl{
		_value: value,
		_rep:   rep,
	}
}

// 获取领域对象编号
func (a *accountImpl) GetDomainId() int {
	return a._value.MemberId
}

// 获取账户值
func (a *accountImpl) GetValue() *member.Account {
	return a._value
}

// 保存
func (a *accountImpl) Save() (int, error) {
	a._value.UpdateTime = time.Now().Unix()
	return a._rep.SaveAccount(a._value)
}

// 根据编号获取余额变动信息
func (a *accountImpl) GetBalanceInfo(id int) *member.BalanceInfo {
	return a._rep.GetBalanceInfo(id)
}

// 根据号码获取余额变动信息
func (a *accountImpl) GetBalanceInfoByNo(no string) *member.BalanceInfo {
	return a._rep.GetBalanceInfoByNo(no)
}

// 保存余额变动信息
func (a *accountImpl) SaveBalanceInfo(v *member.BalanceInfo) (int, error) {
	v.MemberId = a.GetDomainId()
	v.UpdateTime = time.Now().Unix()
	if v.CreateTime == 0 {
		v.CreateTime = v.UpdateTime
	}
	return a._rep.SaveBalanceInfo(v)
}

// 充值
// @title 充值标题说明
// @no    充值订单编号
// @amount 金额
func (a *accountImpl) ChargeBalance(chargeType int, title string, tradeNo string, amount float32) error {
	//todo: 客服充值需记录操作人
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}

	if chargeType == member.TypeBalanceNetPayCharge || chargeType == member.TypeBalanceSystemCharge ||
		chargeType == member.TypeBalanceServiceCharge || chargeType == member.TypeBalanceOrderRefund {

		v := &member.BalanceInfo{
			Kind:    member.KindBalanceCharge,
			Type:    chargeType,
			Title:   title,
			TradeNo: tradeNo,
			Amount:  amount,
			State:   1,
		}
		_, err := a.SaveBalanceInfo(v)
		if err == nil {
			a._value.Balance += amount
			_, err = a.Save()
		}
		return err
	}
	return errors.New("error charge type")
}

// 扣减余额
func (a *accountImpl) DiscountBalance(title string, tradeNo string, amount float32) (err error) {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if a._value.Balance < amount {
		return member.ErrNotEnoughAmount
	}
	v := &member.BalanceInfo{
		Kind:    member.KindBalanceDiscount,
		Title:   title,
		TradeNo: tradeNo,
		Amount:  -amount,
		State:   1,
	}
	_, err = a.SaveBalanceInfo(v)
	if err == nil {
		a._value.Balance -= amount
		_, err = a.Save()
	}
	return err
}

// 赠送金额
func (a *accountImpl) PresentBalance(title string, tradeNo string, amount float32) error {
	//todo:??客服调整
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if len(title) == 0 {
		if amount < 0 {
			title = "赠送账户出账"
		} else {
			title = "赠送账户入账"
		}
	}

	v := &member.BalanceInfo{
		Kind:    member.KindBalancePresent,
		Title:   title,
		TradeNo: tradeNo,
		Amount:  amount,
		State:   1,
	}
	_, err := a.SaveBalanceInfo(v)
	if err == nil {
		a._value.PresentBalance += amount
		if amount > 0 {
			a._value.TotalPresentFee += amount
		}
		_, err = a.Save()
	}
	return err
}

// 扣减奖金
func (a *accountImpl) DiscountPresent(title string, tradeNo string,
	amount float32, mustLargeZero bool) error {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if mustLargeZero && a._value.PresentBalance < amount {
		return member.ErrNotEnoughAmount
	}

	if len(title) == 0 {
		title = "出账"
	}

	v := &member.BalanceInfo{
		Kind:    member.KindPresentDiscount,
		Title:   title,
		TradeNo: tradeNo,
		Amount:  -amount,
		State:   1,
	}
	_, err := a.SaveBalanceInfo(v)
	if err == nil {
		a._value.PresentBalance -= amount
		_, err = a.Save()
	}
	return err
}

// 流通账户余额充值，如扣除,amount传入负数金额
func (a *accountImpl) ChargeFlowBalance(title string, tradeNo string, amount float32) error {
	if len(title) == 0 {
		if amount > 0 {
			title = "流动账户入账"
		} else {
			title = "流动账户出账"
		}
	}
	v := &member.BalanceInfo{
		Kind:    member.KindBalanceFlow,
		Title:   title,
		TradeNo: tradeNo,
		Amount:  amount,
		State:   1,
	}
	_, err := a.SaveBalanceInfo(v)
	if err == nil {
		a._value.FlowBalance += amount
		_, err = a.Save()
	}
	return err
}

// 支付单抵扣消费,tradeNo为支付单单号
func (a *accountImpl) PaymentDiscount(tradeNo string, amount float32) error {
	if amount < 0 || len(tradeNo) == 0 {
		return errors.New("amount error or missing trade no")
	}

	if amount > a._value.Balance {
		return member.ErrOutOfBalance
	}

	v := &member.BalanceInfo{
		Kind:    member.KindBalanceShopping,
		Type:    1,
		Title:   "订单抵扣",
		TradeNo: tradeNo,
		Amount:  -amount,
		State:   1,
	}
	_, err := a.SaveBalanceInfo(v)
	if err == nil {
		a._value.Balance -= amount
		_, err = a.Save()
	}
	return err
}

//　增加积分
func (a *accountImpl) AddIntegral(logType int, outerNo string, value int, remark string) error {
	if value <= 0 {
		return member.ErrIncorrectQuota
	}
	if logType <= 0 {
		logType = member.TypeIntegralPresent
	}
	if logType == member.TypeIntegralShoppingPresent && outerNo == "" {
		return member.ErrMissingOuterNo
	}
	l := &member.IntegralLog{
		MemberId:   a._value.MemberId,
		Type:       logType,
		OuterNo:    outerNo,
		Value:      value,
		Remark:     remark,
		CreateTime: time.Now().Unix(),
	}
	err := a._rep.SaveIntegralLog(l)
	if err == nil {
		a._value.Integral += value
		_, err = a.Save()
	}
	return err
}

// 积分抵扣
func (a *accountImpl) IntegralDiscount(logType int, outerNo string,
	value int, remark string) error {
	if value <= 0 {
		return member.ErrIncorrectQuota
	}
	if a._value.Integral < value {
		return member.ErrNoSuchIntegral
	}

	if logType == member.TypeIntegralPaymentDiscount && outerNo == "" {
		return member.ErrMissingOuterNo
	}

	if logType <= 0 {
		logType = member.TypeIntegralDiscount
	}

	l := &member.IntegralLog{
		MemberId:   a._value.MemberId,
		Type:       logType,
		Value:      -value,
		OuterNo:    outerNo,
		Remark:     remark,
		CreateTime: time.Now().Unix(),
	}
	err := a._rep.SaveIntegralLog(l)
	if err == nil {
		a._value.Integral -= value
		_, err = a.Save()
	}
	return err
}

// 冻结积分,当new为true不扣除积分,反之扣除积分
func (a *accountImpl) FreezesIntegral(value int, new bool, remark string) error {
	if !new {
		if a._value.Integral < value {
			return member.ErrNoSuchIntegral
		}
		a._value.Integral -= value
	}
	a._value.FreezesIntegral += value
	_, err := a.Save()
	if err == nil {
		l := &member.IntegralLog{
			MemberId:   a._value.MemberId,
			Type:       member.TypeIntegralFreeze,
			Value:      -value,
			Remark:     remark,
			CreateTime: time.Now().Unix(),
		}
		err = a._rep.SaveIntegralLog(l)
	}
	return err
}

// 解冻积分
func (a *accountImpl) UnfreezesIntegral(value int, remark string) error {
	if a._value.FreezesIntegral < value {
		return member.ErrNoSuchIntegral
	}
	a._value.FreezesIntegral -= value
	a._value.Integral += value
	_, err := a.Save()
	if err == nil {
		l := &member.IntegralLog{
			MemberId:   a._value.MemberId,
			Type:       member.TypeIntegralUnfreeze,
			Value:      value,
			Remark:     remark,
			CreateTime: time.Now().Unix(),
		}
		err = a._rep.SaveIntegralLog(l)
	}
	return err
}

// 退款
func (a *accountImpl) RequestBackBalance(backType int, title string,
	amount float32) error {
	if amount > a._value.Balance {
		return member.ErrOutOfBalance
	}
	v := &member.BalanceInfo{
		Kind:   member.KindBalanceBack,
		Type:   backType,
		Title:  title,
		Amount: amount,
		State:  0,
	}
	_, err := a.SaveBalanceInfo(v)
	if err == nil {
		a._value.Balance -= amount
		_, err = a.Save()
	}
	return err
}

// 完成退款
func (a *accountImpl) FinishBackBalance(id int, tradeNo string) error {
	v := a.GetBalanceInfo(id)
	if v.Kind == member.KindBalanceBack {
		v.TradeNo = tradeNo
		v.State = 1
		_, err := a.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}

// 请求提现,返回info_id,交易号及错误
func (a *accountImpl) RequestApplyCash(applyType int, title string,
	amount float32, commission float32) (int, string, error) {
	if amount <= 0 {
		return 0, "", member.ErrIncorrectAmount
	}
	if a._value.PresentBalance < amount {
		return 0, "", member.ErrOutOfBalance
	}

	tradeNo := domain.NewTradeNo(00000)

	csnAmount := amount * commission
	finalAmount := amount - csnAmount
	if finalAmount > 0 {
		finalAmount = -finalAmount
	}
	v := &member.BalanceInfo{
		Kind:      member.KindBalanceApplyCash,
		Type:      applyType,
		Title:     title,
		TradeNo:   tradeNo,
		Amount:    finalAmount,
		CsnAmount: csnAmount,
		State:     member.StateApplySubmitted,
	}

	// 提现至余额
	if applyType == member.TypeApplyCashToCharge {
		a._value.Balance += amount
		v.State = member.StateApplyOver
	}

	id, err := a.SaveBalanceInfo(v)
	if err == nil {
		a._value.PresentBalance -= amount
		_, err = a.Save()
	}
	return id, tradeNo, err
}

// 确认提现
func (a *accountImpl) ConfirmApplyCash(id int, pass bool, remark string) error {
	//todo: remark
	v := a.GetBalanceInfo(id)
	if v.Kind == member.KindBalanceApplyCash {
		if pass {
			v.State = member.StateApplyConfirmed
		} else {
			if v.State == member.StateApplyNotPass {
				return dm.ErrState
			}
			v.State = member.StateApplyNotPass
			a._value.PresentBalance += v.CsnAmount + (-v.Amount)
			if _, err := a.Save(); err != nil {
				return err
			}
		}
		_, err := a.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}

// 完成提现
func (a *accountImpl) FinishApplyCash(id int, tradeNo string) error {
	v := a.GetBalanceInfo(id)
	if v.Kind == member.KindBalanceApplyCash {
		v.TradeNo = tradeNo
		v.State = member.StateApplyOver
		_, err := a.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}

// 冻结余额
func (a *accountImpl) Freezes(title string, tradeNo string,
	amount float32, referId int) error {
	if a._value.Balance < amount {
		return member.ErrNotEnoughAmount
	}
	if len(title) == 0 {
		title = "资金冻结"
	}
	v := &member.BalanceInfo{
		Kind:    member.KindBalanceFreezes,
		Title:   title,
		RefId:   referId,
		Amount:  amount,
		TradeNo: tradeNo,
		State:   member.StatusOK,
	}
	a._value.Balance -= amount
	a._value.FreezesFee += amount
	_, err := a.Save()
	if err == nil {
		_, err = a.SaveBalanceInfo(v)
	}
	return err
}

// 解冻金额
func (a *accountImpl) Unfreezes(title string, tradeNo string,
	amount float32, referId int) error {
	if a._value.FreezesFee < amount {
		return member.ErrNotEnoughAmount
	}
	if len(title) == 0 {
		title = "资金解结"
	}
	v := &member.BalanceInfo{
		Kind:    member.KindBalanceUnfreezes,
		Title:   title,
		RefId:   referId,
		Amount:  amount,
		TradeNo: tradeNo,
		State:   member.StatusOK,
	}
	a._value.Balance += amount
	a._value.FreezesFee -= amount
	_, err := a.Save()
	if err == nil {
		_, err = a.SaveBalanceInfo(v)
	}
	return err

}

// 冻结赠送金额
func (a *accountImpl) FreezesPresent(title string, tradeNo string,
	amount float32, referId int) error {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if a._value.PresentBalance < amount {
		return member.ErrNotEnoughAmount
	}
	if len(title) == 0 {
		title = "(赠送)资金冻结"
	}
	v := &member.BalanceInfo{
		Kind:    member.KindBalanceFreezesPresent,
		Title:   title,
		RefId:   referId,
		Amount:  -amount,
		TradeNo: tradeNo,
		State:   member.StatusOK,
	}
	a._value.PresentBalance -= amount
	a._value.FreezesPresent += amount
	_, err := a.Save()
	if err == nil {
		_, err = a.SaveBalanceInfo(v)
	}
	return err
}

// 解冻赠送金额
func (a *accountImpl) UnfreezesPresent(title string, tradeNo string,
	amount float32, referId int) error {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if a._value.FreezesPresent < amount {
		return member.ErrNotEnoughAmount
	}
	if len(title) == 0 {
		title = "(赠送)资金解冻"
	}
	v := &member.BalanceInfo{
		Kind:    member.KindBalanceUnfreezesPresent,
		Title:   title,
		RefId:   referId,
		Amount:  amount,
		TradeNo: tradeNo,
		State:   member.StatusOK,
	}
	a._value.PresentBalance += amount
	a._value.FreezesPresent -= amount
	_, err := a.Save()
	if err == nil {
		_, err = a.SaveBalanceInfo(v)
	}
	return err
}

// 转账余额到其他账户
func (a *accountImpl) TransferBalance(kind int, amount float32,
	tradeNo string, toTitle, fromTitle string) error {
	var err error
	if kind == member.KindBalanceFlow {
		if a._value.Balance < amount {
			return member.ErrNotEnoughAmount
		}
		a._value.Balance -= amount
		a._value.FlowBalance += amount
		if _, err = a.Save(); err == nil {
			a.SaveBalanceInfo(&member.BalanceInfo{
				Kind:    member.KindBalanceTransfer,
				Title:   toTitle,
				Amount:  -amount,
				TradeNo: tradeNo,
				State:   member.StatusOK,
			})

			a.SaveBalanceInfo(&member.BalanceInfo{
				Kind:    member.KindBalanceTransfer,
				Title:   fromTitle,
				Amount:  amount,
				TradeNo: tradeNo,
				State:   member.StatusOK,
			})
		}
		return err
	}
	return member.ErrNotSupportTransfer
}

// 转账返利账户,kind为转账类型，如 KindBalanceTransfer等
// commission手续费
func (a *accountImpl) TransferPresent(kind int, amount float32, commission float32,
	tradeNo string, toTitle string, fromTitle string) error {
	var err error
	if kind == member.KindBalanceFlow {
		if a._value.Balance < amount {
			return member.ErrNotEnoughAmount
		}
		a._value.Balance -= amount
		a._value.FlowBalance += amount
		if _, err = a.Save(); err == nil {
			a.SaveBalanceInfo(&member.BalanceInfo{
				Kind:    member.KindBalanceTransfer,
				Title:   toTitle,
				Amount:  -amount,
				TradeNo: tradeNo,
				State:   member.StatusOK,
			})

			a.SaveBalanceInfo(&member.BalanceInfo{
				Kind:    kind,
				Title:   fromTitle,
				Amount:  amount,
				TradeNo: tradeNo,
				State:   member.StatusOK,
			})
		}
		return err
	}
	return member.ErrNotSupportTransfer
}

// 转账活动账户,kind为转账类型，如 KindBalanceTransfer等
// commission手续费
func (a *accountImpl) TransferFlow(kind int, amount float32, commission float32,
	tradeNo string, toTitle string, fromTitle string) error {
	var err error

	csnAmount := commission * amount
	finalAmount := amount - csnAmount

	if kind == member.KindBalancePresent {
		if a._value.FlowBalance < finalAmount {
			return member.ErrNotEnoughAmount
		}

		a._value.FlowBalance -= amount
		a._value.PresentBalance += finalAmount
		a._value.TotalPresentFee += finalAmount

		if _, err = a.Save(); err == nil {
			a.SaveBalanceInfo(&member.BalanceInfo{
				Kind:      member.KindBalanceTransfer,
				Title:     toTitle,
				Amount:    -amount,
				TradeNo:   tradeNo,
				CsnAmount: csnAmount,
				State:     member.StatusOK,
			})

			a.SaveBalanceInfo(&member.BalanceInfo{
				Kind:    kind,
				Title:   fromTitle,
				Amount:  finalAmount,
				TradeNo: tradeNo,
				State:   member.StatusOK,
			})
		}
		return err
	}

	return member.ErrNotSupportTransfer
}

// 将活动金转给其他人
func (a *accountImpl) TransferFlowTo(memberId int, kind int,
	amount float32, commission float32, tradeNo string,
	toTitle string, fromTitle string) error {

	var err error
	csnAmount := commission * amount
	finalAmount := amount + csnAmount // 转账方付手续费

	m := a._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	acc2 := m.GetAccount()

	if kind == member.KindBalanceFlow {
		if a._value.FlowBalance < finalAmount {
			return member.ErrNotEnoughAmount
		}

		a._value.FlowBalance -= finalAmount
		acc2.GetValue().FlowBalance += amount

		if _, err = a.Save(); err == nil {

			a.SaveBalanceInfo(&member.BalanceInfo{
				Kind:      member.KindBalanceTransfer,
				Title:     toTitle,
				Amount:    -finalAmount,
				CsnAmount: csnAmount,
				RefId:     memberId,
				TradeNo:   tradeNo,
				State:     member.StatusOK,
			})

			if _, err = acc2.Save(); err == nil {
				acc2.SaveBalanceInfo(&member.BalanceInfo{
					Kind:    kind,
					Title:   fromTitle,
					Amount:  amount,
					RefId:   a._value.MemberId,
					TradeNo: tradeNo,
					State:   member.StatusOK,
				})
			}
		}
		return err
	}

	return member.ErrNotSupportTransfer
}
