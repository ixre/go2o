/**
 * Copyright 2015 @ S1N1 Team.
 * name : account
 * author : jarryliu
 * date : 2015-07-24 08:50
 * description :
 * history :
 */
package member

import (
	"go2o/src/core/domain/interface/member"
	"time"
	"errors"
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
func (this *Account) GetBalanceInfo(id int) *member.BalanceInfoValue{
	return this._rep.GetBalanceInfo(id)
}

// 根据号码获取余额变动信息
func (this *Account) GetBalanceInfoByNo(no string) *member.BalanceInfoValue{
	return this._rep.GetBalanceInfoByNo(no)
}

// 保存余额变动信息
func (this *Account) SaveBalanceInfo(v *member.BalanceInfoValue)(int,error){
	v.MemberId = this.GetDomainId()
	v.UpdateTime = time.Now().Unix()
	return this._rep.SaveBalanceInfo(v)
}

// 充值
// @title 充值标题说明
// @no    充值订单编号
// @amount 金额
func (this *Account) Charge(chargeType int,title string,tradeNo string,amount float32)(error){
	v := &member.BalanceInfoValue{
		Kind:member.KindBalanceCharge,
		Type:chargeType,
		Title : title,
		TradeNo:tradeNo,
		Amount:amount,
		State : 1,
	}
	_,err := this.SaveBalanceInfo(v)
	return err
}

// 退款
func (this *Account) RequestBackBalance(backType int,title string,amount float32)(error){
	v := &member.BalanceInfoValue{
		Kind:member.KindBalanceBack,
		Type:backType,
		Title : title,
		Amount:amount,
		State : 0,
	}
	_,err := this.SaveBalanceInfo(v)
	return err
}

// 完成退款
func (this *Account) FinishBackBalance(id int,tradeNo string)(error){
	v := this.GetBalanceInfo(id)
	if v.Kind == member.KindBalanceBack{
		v.TradeNo = tradeNo
		v.State = 1
		_,err := this.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}

// 请求提现
func (this *Account) RequestApplyCash(applyType int,title string,amount float32)(error){
	v := &member.BalanceInfoValue{
		Kind:member.KindBalanceBack,
		Type:applyType,
		Title : title,
		Amount:amount,
		State : 0,
	}
	_,err := this.SaveBalanceInfo(v)
	return err
}


// 确认提现
func (this *Account) ConfirmApplyCash(id int)(error){
	v := this.GetBalanceInfo(id)
	if v.Kind == member.KindBalanceApplyCash{
		v.State = 1
		_,err := this.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}

// 完成提现
func (this *Account) FinishApplyCash(id int,tradeNo string)(error){
	v := this.GetBalanceInfo(id)
	if v.Kind == member.KindBalanceApplyCash{
		v.TradeNo = tradeNo
		v.State = 2
		_,err := this.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}