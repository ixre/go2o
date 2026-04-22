package member

import (
	"errors"
	"fmt"
	"go2o/core/domain/interface/domain/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/wallet"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"math"
	"time"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : account_bak.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-04 17:15
 * description :
 * history :
 */

func (a *accountImpl) chargeWallet_(title string, amount float32, outerNo string, remark string) error {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if title == "" {
		if amount < 0 {
			title = "钱包账户出账"
		} else {
			title = "钱包账户入账"
		}
	}
	l, err := a.createWalletLog(member.KindCharge, title, amount, outerNo, true)
	if err == nil {
		l.Remark = remark
		_, err = a.rep.SaveWalletAccountLog(l)
		if err == nil {
			a.value.WalletBalance += amount
			a.value.TotalWalletAmount += amount
			_, err = a.Save()
		}
	}
	return err
}

func (a *accountImpl) walletConsume_(title string, amount float32, outerNo string, remark string) error {
	if a.value.WalletBalance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	l, err := a.createWalletLog(member.KindConsume, title, -amount, outerNo, true)
	if err == nil {
		l.Remark = remark
		_, err = a.rep.SaveWalletAccountLog(l)
		if err == nil {
			a.value.WalletBalance -= amount
			_, err = a.Save()
		}
	}
	return err
}

// 赠送金额(指定业务类型)
func (a *accountImpl) walletRefund_(kind int, title string,
	outerNo string, amount float32, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if title == "" {
		if amount < 0 {
			title = "钱包账户出账"
		} else {
			title = "钱包账户入账"
		}
	}
	unix := time.Now().Unix()
	v := &member.WalletAccountLog{
		MemberId:    a.GetDomainId(),
		Kind:        kind,
		Title:       title,
		OuterNo:     outerNo,
		Amount:      amount,
		ReviewStatus: enum.ReviewPass,
		RelateUser:  relateUser,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveWalletAccountLog(v)
	if err == nil {
		a.value.WalletBalance += amount
		// 退款不能加入到累计赠送金额
		if kind != wallet.KWithdrawRefund &&
			kind != member.KindRefund &&
			amount > 0 {
			a.value.TotalWalletAmount += amount
		}
		_, err = a.Save()
	}
	return err
}

// 扣减奖金,mustLargeZero是否必须大于0, 赠送金额存在扣为负数的情况
func (a *accountImpl) walletDiscount_(title string, amount float32, outerNo string, remark string) error {
	mustLargeZero := false
	if mustLargeZero && a.value.WalletBalance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	l, err := a.createWalletLog(member.KindDiscount, title, -amount, outerNo, true)
	if err == nil {
		l.Remark = remark
		a.value.WalletBalance -= amount
		_, err = a.Save()
	}
	return err
}

// 请求提现,返回info_id,交易号及错误
func (a *accountImpl) RequestWithdrawal_(takeKind int, title string,
	amount2 int, transactionFee int, bankAccountNo string) (int32, string, error) {
	amount := float32(amount2) / 100
	if takeKind != wallet.KWithdrawExchange &&
		takeKind != wallet.KWithdrawToBankCard &&
		takeKind != wallet.KWithdrawToPayWallet {
		return 0, "", member.ErrNotSupportTakeOutBusinessKind
	}
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return 0, "", member.ErrIncorrectAmount
	}
	// 检测是否开启提现
	takOutOn := a.registryRepo.Get(registry.MemberTakeOutOn).BoolValue()
	if !takOutOn {
		msg, _ := a.registryRepo.GetValue(registry.MemberTakeOutMessage)
		return 0, "", errors.New(msg)
	}

	// 检测是否实名
	mustTrust, _ := a.registryRepo.GetValue(registry.MemberWithdrawalMustTrust)
	if mustTrust == "true" {
		trust := a.member.Profile().GetCertificationInfo()
		if trust.ReviewStatus != int(enum.ReviewPass) {
			return 0, "", member.ErrTakeOutNotTrust
		}
	}

	// 检测非正式会员提现
	lv := a.mm.LevelManager().GetLevelById(a.member.GetValue().Level)
	if lv != nil && lv.IsOfficial == 0 {
		return 0, "", errors.New(fmt.Sprintf(
			member.ErrTakeOutLevelNoPerm.Error(), lv.Name))
	}
	// 检测余额
	if a.value.WalletBalance < amount {
		return 0, "", member.ErrOutOfBalance
	}
	// 检测提现金额是否超过限制
	minAmount := float32(a.registryRepo.Get(registry.MemberMinTakeOutAmount).FloatValue())
	if amount < minAmount {
		return 0, "", errors.New(fmt.Sprintf(member.ErrLessTakeAmount.Error(),
			format.FormatFloat(minAmount)))
	}
	maxAmount := float32(a.registryRepo.Get(registry.MemberMaxTakeOutAmount).FloatValue())
	if maxAmount > 0 && amount > maxAmount {
		return 0, "", errors.New(fmt.Sprintf(member.ErrOutTakeAmount.Error(),
			format.FormatFloat(maxAmount)))
	}
	// 检测是否超过限制
	maxTimes := a.registryRepo.Get(registry.MemberMaxTakeOutTimesOfDay).IntValue()
	if maxTimes > 0 {
		takeTimes := a.rep.GetTodayTakeOutTimes(a.GetDomainId())
		if takeTimes >= maxTimes {
			return 0, "", member.ErrAccountOutOfTakeOutTimes
		}
	}

	tradeNo := domain.NewTradeNo(8, int(a.member.GetAggregateRootId()))
	csnAmount := float32(transactionFee) / 100
	finalAmount := amount - csnAmount
	if finalAmount > 0 {
		finalAmount = -finalAmount
	}
	unix := time.Now().Unix()
	v := &member.WalletAccountLog{
		MemberId:    a.GetDomainId(),
		Kind:        takeKind,
		Title:       title,
		OuterNo:     tradeNo,
		Amount:      finalAmount,
		CsnFee:      csnAmount,
		ReviewStatus: enum.ReviewPending,
		RelateUser:  member.DefaultRelateUser,
		Remark:      "",
		CreateTime:  unix,
		UpdateTime:  unix,
	}

	// 提现至余额
	if takeKind == wallet.KWithdrawExchange {
		a.value.Balance += amount
		v.ReviewStatus = enum.ReviewPass
	}
	a.value.WalletBalance -= amount
	_, err := a.Save()
	if err == nil {
		go a.rep.AddTodayTakeOutTimes(a.GetDomainId())
		id, err := a.rep.SaveWalletAccountLog(v)
		return id, tradeNo, err
	}
	return 0, tradeNo, err
}

// 确认提现
func (a *accountImpl) ReviewWithdrawal_(id int32, pass bool, remark string) error {
	v := a.GetWalletLog(id)
	if v == nil || v.MemberId != a.value.MemberId {
		return member.ErrIncorrectInfo
	}
	if v.ReviewStatus != enum.ReviewPending {
		return member.ErrTakeOutState
	}
	// todo: 应该先冻结, 再扣除
	if pass {
		v.ReviewStatus = enum.ReviewPass
	} else {
		v.Remark += "失败:" + remark
		v.ReviewStatus = enum.ReviewReject
		//err := a.Refund(member.AccountWallet,
		//	member.KindWalletTakeOutRefund,
		//	"提现退回",  v.CsnFee+(-v.Amount),v.OuterNo,"")
		//
		// 将手续费修改到提现金额上
		v.Amount -= v.CsnFee
		v.CsnFee = 0
	}
	v.UpdateTime = time.Now().Unix()
	_, err := a.rep.SaveWalletAccountLog(v)
	if err == nil && !pass {
		a.value.WalletBalance += v.CsnFee + (-v.Amount)
	}
	return err
}

// 完成提现
func (a *accountImpl) CompleteTransaction_(id int32, tradeNo string) error {
	v := a.GetWalletLog(id)
	if v == nil || v.MemberId != a.value.MemberId {
		return member.ErrIncorrectInfo
	}
	if v.ReviewStatus != enum.ReviewPass {
		return member.ErrTakeOutState
	}
	v.OuterNo = tradeNo
	v.ReviewStatus = enum.ReviewCompleted
	v.Remark = "转款凭证:" + tradeNo
	_, err := a.rep.SaveWalletAccountLog(v)
	return err

	//if v.Kind == member.KindWalletTakeOutToBankCard {
	//    v.OuterNo = tradeNo
	//    v.State = enum.ReviewCompleted
	//    v.Remark = "银行凭证:" + tradeNo
	//    _, err := a.repo.SaveWalletAccountLog(v)
	//    return err
	//}
	//return member.ErrNotSupportTakeOutBusinessKind
}
