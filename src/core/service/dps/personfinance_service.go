/**
 * Copyright 2015 @ z3q.net.
 * name : personfinance_service
 * author : jarryliu
 * date : 2016-04-01 09:41
 * description :
 * history :
 */
package dps

import (
	"errors"
	"fmt"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/personfinance"
	"go2o/src/core/infrastructure/domain"
)

const (
	TransferFromBalance = 1 + iota //通过余额转入
	TransferFromPresent            //通过创业金转入
)

type personFinanceService struct {
	_rep    personfinance.IPersonFinanceRepository
	_accRep member.IMemberRep
}

func NewPersonFinanceService(rep personfinance.IPersonFinanceRepository,
	accRep member.IMemberRep) *personFinanceService {
	return &personFinanceService{
		_rep:    rep,
		_accRep: accRep,
	}
}

func (this *personFinanceService) GetRiseInfo(personId int) (personfinance.RiseInfoValue, error) {
	pf := this._rep.GetPersonFinance(personId)
	return pf.GetRiseInfo().Value()
}

// 开通增利服务
func (this *personFinanceService) OpenRiseService(personId int) error {
	pf := this._rep.GetPersonFinance(personId)
	return pf.CreateRiseInfo()
}

func (this *personFinanceService) RiseTransferIn(personId int, transferFrom int, amount float32) error {
	r := this._rep.GetPersonFinance(personId).GetRiseInfo()
	if amount < personfinance.RiseMinTransferInAmount {
		//金额不足最低转入金额
		return errors.New(fmt.Sprintf(personfinance.ErrLessThanMinTransferIn.Error(),
			personfinance.RiseMinTransferInAmount))
	}

	m := this._accRep.GetMember(personId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	acc := m.GetAccount()
	if transferFrom == TransferFromBalance { //从余额转入
		if err := acc.DiscountBalance("理财转入", domain.NewTradeNo(10000), amount); err != nil {
			return err
		}
		return r.TransferIn(amount)
	}

	if transferFrom == TransferFromPresent { //从奖金转入
		if err := acc.DiscountPresent("理财转入", domain.NewTradeNo(10000), amount, true); err != nil {
			return err
		}
		return r.TransferIn(amount)
	}

	return errors.New("未知的转入方式")
}

func (this *personFinanceService) RiseTransferOut(personId int, amount float32) error {
	r := this._rep.GetPersonFinance(personId).GetRiseInfo()
	return r.TransferOut(amount)
}
