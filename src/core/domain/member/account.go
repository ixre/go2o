/**
 * Copyright 2015 @ S1N1 Team.
 * name : account
 * author : jarryliu
 * date : 2015-07-24 08:50
 * description :
 * history :
 */
package member

import "go2o/src/core/domain/interface/member"

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
