/**
 * Copyright 2015 @ z3q.net.
 * name : person_finance
 * author : jarryliu
 * date : 2016-03-31 17:17
 * description :
 * history :
 */
package personfinance

import (
	"errors"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/personfinance"
	"time"
)

var _ personfinance.IPersonFinance = new(PersonFinance)

type PersonFinance struct {
	_personId int
	_accRep   member.IMemberRep
	_rep      personfinance.IPersonFinanceRepository
}

func NewPersonFinance(personId int, rep personfinance.IPersonFinanceRepository,
	accRep member.IMemberRep) personfinance.IPersonFinance {
	return &PersonFinance{
		_personId: personId,
		_accRep:   accRep,
		_rep:      rep,
	}
}

// 获取聚合根
func (this *PersonFinance) GetAggregateRootId() int {
	return this._personId
}

// 获取账号
func (this *PersonFinance) GetMemberAccount() member.IAccount {
	return this._accRep.GetMember(this._personId).GetAccount()
}

// 获取增利账户信息(类:余额宝)
func (this *PersonFinance) GetRiseInfo() personfinance.IRiseInfo {
	return newRiseInfo(this.GetAggregateRootId(), this._rep, this._accRep)
}

// 创建增利账户信息
func (this *PersonFinance) CreateRiseInfo() error {
	_, err := this.GetRiseInfo().Value()
	if err != nil {
		v := &personfinance.RiseInfoValue{
			PersonId:   this.GetAggregateRootId(),
			UpdateTime: time.Now().Unix(),
		}
		_, err = this._rep.SaveRiseInfo(v)
		return err
	}
	return errors.New("rise info exists!")
}
