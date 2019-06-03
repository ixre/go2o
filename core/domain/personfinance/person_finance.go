/**
 * Copyright 2015 @ to2.net.
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
	personId int64
	accRepo  member.IMemberRepo
	rep      personfinance.IPersonFinanceRepository
}

func NewPersonFinance(personId int64, rep personfinance.IPersonFinanceRepository,
	accRepo member.IMemberRepo) personfinance.IPersonFinance {
	return &PersonFinance{
		personId: personId,
		accRepo:  accRepo,
		rep:      rep,
	}
}

// 获取聚合根
func (p *PersonFinance) GetAggregateRootId() int64 {
	return p.personId
}

// 获取账号
func (p *PersonFinance) GetMemberAccount() member.IAccount {
	return p.accRepo.GetMember(p.personId).GetAccount()
}

// 获取增利账户信息(类:余额宝)
func (p *PersonFinance) GetRiseInfo() personfinance.IRiseInfo {
	return newRiseInfo(p.GetAggregateRootId(), p, p.rep, p.accRepo)
}

// 创建增利账户信息
func (p *PersonFinance) CreateRiseInfo() error {
	_, err := p.GetRiseInfo().Value()
	if err != nil {
		v := &personfinance.RiseInfoValue{
			PersonId:   p.GetAggregateRootId(),
			UpdateTime: time.Now().Unix(),
		}
		_, err = p.rep.SaveRiseInfo(v)
		return err
	}
	return errors.New("rise info exists!")
}

// 同步到会员账户理财数据
func (p *PersonFinance) SyncToAccount() error {
	var balance float32
	var totalAmount float32
	var growEarnings float32      // 当前收益
	var totalGrowEarnings float32 // 总收益
	r := p.GetRiseInfo()
	if r, err := r.Value(); err != nil {
		return err
	} else {
		balance += r.Balance
		totalAmount += r.TotalAmount
		growEarnings += r.Rise
		totalGrowEarnings += r.TotalRise
	}
	return p.accRepo.SaveGrowAccount(p.GetAggregateRootId(),
		balance, totalAmount, growEarnings, totalGrowEarnings,
		time.Now().Unix())
}
