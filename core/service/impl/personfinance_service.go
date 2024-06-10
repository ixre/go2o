/**
 * Copyright 2015 @ 56x.net.
 * name : personfinance_service
 * author : jarryliu
 * date : 2016-04-01 09:41
 * description :
 * history :
 */
package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/personfinance"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/go2o/core/variable"
	"github.com/ixre/gof/log"
)

var _ proto.FinanceServiceServer = new(personFinanceService)

type personFinanceService struct {
	repo    personfinance.IPersonFinanceRepository
	memRepo member.IMemberRepo
	serviceUtil
	proto.UnimplementedFinanceServiceServer
}

func NewPersonFinanceService(rep personfinance.IPersonFinanceRepository,
	accRepo member.IMemberRepo) proto.FinanceServiceServer {
	return &personFinanceService{
		repo:    rep,
		memRepo: accRepo,
	}
}

func (p *personFinanceService) GetRiseInfo(_ context.Context, id *proto.PersonId) (*proto.SRiseInfo, error) {
	pf := p.repo.GetPersonFinance(id.Value)
	if pf != nil {
		v, err := pf.GetRiseInfo().Value()
		if err == nil {
			return p.parseRiseInfoDto(v), nil
		}
	}
	return nil, fmt.Errorf("no such rise info")
}

// 转入(业务放在service,是为person_finance解耦)
// Parameters:
//   - PersonId
//   - TransferWith
//   - Amount
func (p *personFinanceService) RiseTransferIn(_ context.Context, r *proto.TransferInRequest) (*proto.Result, error) {
	//return errors.New("服务暂时不可用")
	pf := p.repo.GetPersonFinance(r.PersonId)
	err := pf.GetRiseInfo().TransferIn(float32(r.Amount),
		personfinance.TransferWith(r.TransferWith))
	return p.result(err), nil
}

// 转出
func (p *personFinanceService) RiseTransferOut(_ context.Context, r *proto.RiseTransferOutRequest) (*proto.Result, error) {
	//return errors.New("系统正在升级，暂停服务!")

	pf := p.repo.GetPersonFinance(r.PersonId)
	ir := pf.GetRiseInfo()

	m := p.memRepo.GetMember(r.PersonId)
	if m == nil {
		return p.error(member.ErrNoSuchMember), nil
	}
	acc := m.GetAccount()
	tradeNo := domain.NewTradeNo(8, int(r.PersonId))
	if r.TransferWith == personfinance.TransferOutWithBalance {
		//转入余额
		err := ir.TransferOut(float32(r.Amount),
			personfinance.TransferWith(r.TransferWith),
			personfinance.RiseStateOk)
		if err == nil {
			_, err = acc.CarryTo(member.AccountBalance,
				member.AccountOperateData{
					Title:   variable.AliasGrowthAccount + "转出",
					Amount:  int(r.Amount * 100),
					OuterNo: tradeNo,
					Remark:  "sys",
				}, false, 0)
			if err != nil {
				log.Println("[ TransferOut][ Error]:", err.Error())
			}
			err = pf.SyncToAccount()
		}
		return p.error(err), nil
	}

	if r.TransferWith == personfinance.TransferFromWithWallet {
		//转入钱包
		err := ir.TransferOut(float32(r.Amount),
			personfinance.TransferWith(r.TransferWith),
			personfinance.RiseStateOk)
		if err == nil {
			_, err = acc.CarryTo(member.AccountWallet,
				member.AccountOperateData{
					Title:   variable.AliasGrowthAccount + "转出",
					Amount:  int(r.Amount * 100),
					OuterNo: tradeNo,
					Remark:  "sys",
				}, false, 0)
			if err != nil {
				log.Println("[ TransferOut][ Error]:", err.Error())
			}
			err = pf.SyncToAccount()
		}
		return p.error(err), nil
	}

	if r.TransferWith == personfinance.TransferOutWithBank {
		if b := m.Profile().GetBankCard(r.BankAccountNo); b == nil {
			return p.error(member.ErrBankNoSuchCard), nil
		}
		err := ir.TransferOut(float32(r.Amount),
			personfinance.TransferWith(r.TransferWith),
			personfinance.RiseStateOk)
		if err == nil {
			err = pf.SyncToAccount()
		}
		return p.error(err), nil
	}
	return p.error(errors.New("暂时无法提供服务")), nil
}

// 结算收益(按日期每天结息)
func (p *personFinanceService) RiseSettleByDay(_ context.Context, r *proto.RiseSettleRequest) (*proto.Result, error) {
	pf := p.repo.GetPersonFinance(r.PersonId)
	ir := pf.GetRiseInfo()
	err := ir.RiseSettleByDay(r.SettleDay, float32(r.Ratio))
	if err == nil {
		//同步到会员账户
		err = pf.SyncToAccount()
	}
	return p.error(err), nil
}

// 提交转入/转出日志
func (p *personFinanceService) CommitTransfer(_ context.Context, r *proto.CommitTransferRequest) (*proto.Result, error) {
	pf := p.repo.GetPersonFinance(r.PersonId)
	var err error
	if pf == nil {
		err = personfinance.ErrNoSuchRiseInfo
	} else {
		rs := pf.GetRiseInfo()
		err = rs.CommitTransfer(int32(r.LogId))
	}
	return p.error(err), nil
}

// 开通服务
func (p *personFinanceService) OpenRiseService(_ context.Context, id *proto.PersonId) (*proto.Result, error) {
	m := p.memRepo.GetMember(id.Value)
	var err error
	if m == nil {
		err = member.ErrNoSuchMember
	} else {
		if m.GetValue().Level < int(variable.PersonFinanceMinLevelLimit) {
			err = errors.New("会员等级不够,请升级后再开通理财账户！")
		} else {
			pf := p.repo.GetPersonFinance(id.Value)
			err = pf.CreateRiseInfo()
		}
	}
	return p.error(err), nil
}

func (p *personFinanceService) parseRiseInfoDto(v personfinance.RiseInfoValue) *proto.SRiseInfo {
	return &proto.SRiseInfo{
		PersonId:         v.PersonId,
		Balance:          float64(v.Balance),
		SettlementAmount: float64(v.SettlementAmount),
		Rise:             float64(v.Rise),
		TransferIn:       float64(v.TransferIn),
		TotalAmount:      float64(v.TotalAmount),
		TotalRise:        float64(v.TotalRise),
		SettledDate:      v.SettledDate,
		UpdateTime:       v.UpdateTime,
	}
}
