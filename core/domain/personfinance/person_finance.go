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
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/personfinance"
	"time"
)

var _ personfinance.IPersonFinance = new(PersonFinance)

type PersonFinance struct {
	personId int
	accRep   member.IMemberRep
	rep      personfinance.IPersonFinanceRepository
}

func NewPersonFinance(personId int, rep personfinance.IPersonFinanceRepository,
	accRep member.IMemberRep) personfinance.IPersonFinance {
	return &PersonFinance{
		personId: personId,
		accRep:   accRep,
		rep:      rep,
	}
}

// 获取聚合根
func (this *PersonFinance) GetAggregateRootId() int {
	return this.personId
}

// 获取账号
func (this *PersonFinance) GetMemberAccount() member.IAccount {
	return this.accRep.GetMember(this.personId).GetAccount()
}

// 获取增利账户信息(类:余额宝)
func (this *PersonFinance) GetRiseInfo() personfinance.IRiseInfo {
	return newRiseInfo(this.GetAggregateRootId(), this.rep, this.accRep)
}

// 创建增利账户信息
func (this *PersonFinance) CreateRiseInfo() error {
	_, err := this.GetRiseInfo().Value()
	if err != nil {
		v := &personfinance.RiseInfoValue{
			PersonId:   this.GetAggregateRootId(),
			UpdateTime: time.Now().Unix(),
		}
		_, err = this.rep.SaveRiseInfo(v)
		return err
	}
	return errors.New("rise info exists!")
}

// 同步到会员账户理财数据
func (this *PersonFinance) SyncToAccount() error {
	var balance float32
	var totalAmount float32
	var growEarnings float32      // 当前收益
	var totalGrowEarnings float32 // 总收益
	r := this.GetRiseInfo()
	if r, err := r.Value(); err != nil {
		return err
	} else {
		balance += r.Balance
		totalAmount += r.TotalAmount
		growEarnings += r.Rise
		totalGrowEarnings += r.TotalRise
	}
	return this.accRep.SaveGrowAccount(this.GetAggregateRootId(),
		balance, totalAmount, growEarnings, totalGrowEarnings, time.Now().Unix())
}
