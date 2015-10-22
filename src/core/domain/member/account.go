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
	"go2o/src/core/domain/interface/member"
	"time"
)

var _ member.IAccount = new(Account)

type Account struct {
	_value *member.AccountValue
	_rep   member.IMemberRep
}

func NewAccount(value *member.AccountValue, rep member.IMemberRep) member.IAccount {
	return &Account{
		_value: value,
		_rep:   rep,
	}
}

// 获取领域对象编号
func (this *Account) GetDomainId() int {
	return this._value.MemberId
}

// 获取账户值
func (this *Account) GetValue() *member.AccountValue {
	return this._value
}

// 保存
func (this *Account) Save() (int, error) {
	return this._rep.SaveAccount(this._value)
}

// 根据编号获取余额变动信息
func (this *Account) GetBalanceInfo(id int) *member.BalanceInfoValue {
	return this._rep.GetBalanceInfo(id)
}

// 根据号码获取余额变动信息
func (this *Account) GetBalanceInfoByNo(no string) *member.BalanceInfoValue {
	return this._rep.GetBalanceInfoByNo(no)
}

// 保存余额变动信息
func (this *Account) SaveBalanceInfo(v *member.BalanceInfoValue) (int, error) {
	v.MemberId = this.GetDomainId()
	v.UpdateTime = time.Now().Unix()
	if v.CreateTime == 0 {
		v.CreateTime = v.UpdateTime
	}
	return this._rep.SaveBalanceInfo(v)
}

// 充值
// @title 充值标题说明
// @no    充值订单编号
// @amount 金额
func (this *Account) Charge(chargeType int, title string, tradeNo string, amount float32) error {
	//todo: 客服充值需记录操作人

	if chargeType == member.TypeBalanceNetPayCharge || chargeType == member.TypeBalanceSystemCharge ||
		chargeType == member.TypeBalanceServiceCharge {

		v := &member.BalanceInfoValue{
			Kind:    member.KindBalanceCharge,
			Type:    chargeType,
			Title:   title,
			TradeNo: tradeNo,
			Amount:  amount,
			State:   1,
		}
		_, err := this.SaveBalanceInfo(v)
		if err == nil {
			this._value.Balance += amount
			_, err = this.Save()
		}
		return err
	}
	return errors.New("error charge type")
}

// 订单抵扣消费
func (this *Account) OrderDiscount(tradeNo string, amount float32) error {
	if amount < 0 || len(tradeNo) == 0 {
		return errors.New("amount error or missing trade no")
	}

	if amount > this._value.Balance {
		return member.ErrOutOfBalance
	}

	v := &member.BalanceInfoValue{
		Kind:    member.KindBalanceShopping,
		Type:    1,
		Title:   "订单抵扣",
		TradeNo: tradeNo,
		Amount:  amount,
		State:   1,
	}
	_, err := this.SaveBalanceInfo(v)
	if err == nil {
		this._value.Balance -= amount
		_, err = this.Save()
	}
	return err
}

// 退款
func (this *Account) RequestBackBalance(backType int, title string, amount float32) error {

	if amount > this._value.Balance {
		return member.ErrOutOfBalance
	}

	v := &member.BalanceInfoValue{
		Kind:   member.KindBalanceBack,
		Type:   backType,
		Title:  title,
		Amount: amount,
		State:  0,
	}
	_, err := this.SaveBalanceInfo(v)
	if err == nil {
		this._value.Balance -= amount
		_, err = this.Save()
	}
	return err
}

// 完成退款
func (this *Account) FinishBackBalance(id int, tradeNo string) error {
	v := this.GetBalanceInfo(id)
	if v.Kind == member.KindBalanceBack {
		v.TradeNo = tradeNo
		v.State = 1
		_, err := this.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}

// 请求提现
func (this *Account) RequestApplyCash(applyType int, title string, amount float32) error {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if this._value.PresentBalance < amount {
		return member.ErrOutOfBalance
	}

	v := &member.BalanceInfoValue{
		Kind:   member.KindBalanceApplyCash,
		Type:   applyType,
		Title:  title,
		Amount: amount,
		State:  member.StateApplySubmitted,
	}

	// 提现至余额
	if applyType == member.TypeApplyCashToCharge {
		this._value.Balance += amount
		v.State = member.StateApplyOver
	}

	_, err := this.SaveBalanceInfo(v)
	if err == nil {
		this._value.PresentBalance -= amount
		_, err = this.Save()
	}
	return err
}

// 确认提现
func (this *Account) ConfirmApplyCash(id int, pass bool, remark string) error {
	//todo: remark
	v := this.GetBalanceInfo(id)
	if v.Kind == member.KindBalanceApplyCash {
		if pass {
			v.State = member.StateApplyConfirmed
		} else {
			v.State = member.StateApplyNotPass
			this._value.PresentBalance += v.Amount
			if _, err := this.Save(); err != nil {
				return err
			}
		}
		_, err := this.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}

// 完成提现
func (this *Account) FinishApplyCash(id int, tradeNo string) error {
	v := this.GetBalanceInfo(id)
	if v.Kind == member.KindBalanceApplyCash {
		v.TradeNo = tradeNo
		v.State = member.StateApplyOver
		_, err := this.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}

// 冻结余额
func (this *Account) Freezes(amount float32) error {
	if this._value.Balance < amount {
		return member.ErrNotEnoughAmount
	}
	v := &member.BalanceInfoValue{
		Kind:   member.KindBalanceFreezes,
		Title:  "资金冻结",
		Amount: amount,
		State:  member.StatusOK,
	}
	this._value.Balance -= amount
	this._value.FreezesFee += amount
	_, err := this.Save()
	if err == nil {
		_, err = this.SaveBalanceInfo(v)
	}
	return err
}

// 解冻金额
func (this *Account) Unfreezes(amount float32) error {
	if this._value.FreezesFee < amount {
		return member.ErrNotEnoughAmount
	}
	v := &member.BalanceInfoValue{
		Kind:   member.KindBalanceUnfreezes,
		Title:  "资金解结",
		Amount: amount,
		State:  member.StatusOK,
	}
	this._value.Balance += amount
	this._value.FreezesFee -= amount
	_, err := this.Save()
	if err == nil {
		_, err = this.SaveBalanceInfo(v)
	}
	return err

}

// 冻结赠送金额
func (this *Account) FreezesPresent(amount float32) error {
	if this._value.PresentBalance < amount {
		return member.ErrNotEnoughAmount
	}
	v := &member.BalanceInfoValue{
		Kind:   member.KindBalanceFreezesPresent,
		Title:  "(赠送)资金冻结",
		Amount: amount,
		State:  member.StatusOK,
	}
	this._value.PresentBalance -= amount
	this._value.FreezesPresent += amount
	_, err := this.Save()
	if err == nil {
		_, err = this.SaveBalanceInfo(v)
	}
	return err
}

// 解冻赠送金额
func (this *Account) UnfreezesPresent(amount float32) error {
	if this._value.FreezesPresent < amount {
		return member.ErrNotEnoughAmount
	}
	v := &member.BalanceInfoValue{
		Kind:   member.KindBalanceUnfreezesPresent,
		Title:  "(赠送)资金解冻",
		Amount: amount,
		State:  member.StatusOK,
	}
	this._value.PresentBalance += amount
	this._value.FreezesPresent -= amount
	_, err := this.Save()
	if err == nil {
		_, err = this.SaveBalanceInfo(v)
	}
	return err
}
