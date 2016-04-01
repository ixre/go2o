/**
 * Copyright 2015 @ z3q.net.
 * name : personfinance_service
 * author : jarryliu
 * date : 2016-04-01 09:41
 * description :
 * history :
 */
package dps

import "go2o/src/core/domain/interface/personfinance"

type personFinanceService struct {
	_rep personfinance.IPersonFinanceRepository
}

func NewPersonFinanceService(rep personfinance.IPersonFinanceRepository) *personFinanceService {
	return &personFinanceService{
		_rep: rep,
	}
}

func (this *personFinanceService) GetRiseInfo(personId int)(personfinance.RiseInfoValue,error){
	pf := this._rep.GetPersonFinance(personId)
	return pf.GetRiseInfo().Value()
}

// 开通增利服务
func (this *personFinanceService) OpenRiseService(personId int) error {
	pf := this._rep.GetPersonFinance(personId)
	return pf.CreateRiseInfo()
}

func (this *personFinanceService) RiseTransferIn(personId int, amount float32) error {
	r := this._rep.GetPersonFinance(personId).GetRiseInfo()
	return r.TransferIn(amount)
}

func (this *personFinanceService) RiseTransferOut(personId int, amount float32) error {
	r := this._rep.GetPersonFinance(personId).GetRiseInfo()
	return r.TransferOut(amount)
}
