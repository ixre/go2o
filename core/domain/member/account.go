/**
 * Copyright 2015 @ 56x.net.
 * name : account
 * author : jarryliu
 * date : 2015-07-24 08:50
 * description :
 * history :
 */
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
	"go2o/core/msq"
	"math"
	"strconv"
	"strings"
	"time"
)

var _ member.IAccount = new(accountImpl)

type accountImpl struct {
	member       *memberImpl
	wallet       wallet.IWallet
	mm           member.IMemberManager
	value        *member.Account
	rep          member.IMemberRepo
	walletRepo   wallet.IWalletRepo
	registryRepo registry.IRegistryRepo
}

func NewAccount(m *memberImpl, value *member.Account,
	rep member.IMemberRepo, mm member.IMemberManager,
	walletRepo wallet.IWalletRepo,
	registryRepo registry.IRegistryRepo) member.IAccount {
	var wal wallet.IWallet
	if wc := value.WalletCode; len(wc) > 0 {
		wal = walletRepo.GetWalletByCode(value.WalletCode)
	}
	impl := &accountImpl{
		member:       m,
		value:        value,
		wallet:       wal,
		rep:          rep,
		walletRepo:   walletRepo,
		mm:           mm,
		registryRepo: registryRepo,
	}
	if value.MemberId > 0 && wal == nil {
		impl.createWallet()
	}
	return impl
}

func (a *accountImpl) createWallet() {
	flag := wallet.FlagCharge | wallet.FlagDiscount
	a.wallet = a.walletRepo.CreateWallet(
		a.member.GetAggregateRootId(),
		a.member.value.User,
		1, "MemberWallet", flag)
	a.wallet.Save()
	// 绑定钱包
	a.value.WalletCode = a.wallet.Get().HashCode
}

// 获取领域对象编号
func (a *accountImpl) GetDomainId() int64 {
	return a.value.MemberId
}

// 获取账户值
func (a *accountImpl) GetValue() *member.Account {
	return a.value
}

// 保存
func (a *accountImpl) Save() (int64, error) {
	// 判断是否新建账号
	origin := a.rep.GetAccount(a.member.GetAggregateRootId())
	isCreate := origin == nil
	// 更新账户
	a.value.MemberId = a.member.GetAggregateRootId()
	a.value.UpdateTime = time.Now().Unix()
	n, err := a.rep.SaveAccount(a.value)
	if err == nil {
		// 创建钱包
		if isCreate {
			a.createWallet()
			a.rep.SaveAccount(a.value)
		}
		// 推送钱包更新消息
		if !isCreate {
			go msq.PushDelay(msq.MemberAccountUpdated, strconv.Itoa(int(a.value.MemberId)), 500)
		}
	}
	return n, err
}

func (a *accountImpl) Wallet() wallet.IWallet {
	return a.wallet
}

// 设置优先(默认)支付方式, account 为账户类型
func (a *accountImpl) SetPriorityPay(account int, enabled bool) error {
	if enabled {
		support := false
		if account == member.AccountBalance ||
			account == member.AccountWallet ||
			account == member.AccountIntegral {
			support = true
		}
		if support {
			a.value.PriorityPay = account
			_, err := a.Save()
			return err
		}
		// 不支持支付的账号类型
		return member.ErrNotSupportPaymentAccountType
	}

	// 关闭默认支付
	a.value.PriorityPay = 0
	_, err := a.Save()
	return err
}

// 充值
func (a *accountImpl) Charge(account int32, title string,
	amount int, outerNo string, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectQuota
	}
	switch account {
	case member.AccountIntegral:
		return a.integralCharge(title, amount, outerNo, remark)
	case member.AccountBalance:
		return a.chargeBalance(title, float32(amount)/100, outerNo, remark)
	case member.AccountWallet:
		return a.chargeWallet(title, amount, outerNo, remark)
	case member.AccountFlow:
		return a.chargeFlow(title, float32(amount)/100, outerNo, remark)
	}
	return member.ErrNotSupportAccountType
}

func (a *accountImpl) Adjust(account int, title string, amount int, remark string, relateUser int64) error {
	if amount == 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if title == "" || remark == "" {
		return member.ErrNoSuchLogTitleOrRemark
	}
	if relateUser <= 0 {
		return member.ErrNoSuchRelateUser
	}
	switch account {
	case member.AccountIntegral:
		return a.adjustIntegralAccount(title, amount, remark, relateUser)
	case member.AccountBalance:
		return a.adjustBalanceAccount(title, float32(amount)/100, remark, relateUser)
	case member.AccountWallet:
		return a.walletAdjust(title, amount, remark, relateUser)
	case member.AccountFlow:
		return a.adjustFlowAccount(title, float32(amount)/100, remark, relateUser)
	}
	panic("not support other account adjust")
}

// 消耗
func (a *accountImpl) Consume(account int, title string, amount int, outerNo string, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectQuota
	}
	switch account {
	case member.AccountIntegral:
		return a.integralConsume(title, amount, outerNo, remark)
	case member.AccountBalance:
		return a.balanceConsume(title, float32(amount)/100, outerNo, remark)
	case member.AccountWallet:
		return a.walletConsume(title, amount, outerNo, remark)
	case member.AccountFlow:
		return a.flowAccountConsume(title, float32(amount)/100, outerNo, remark)
	}
	return member.ErrNotSupportAccountType
}

func (a *accountImpl) Discount(account int, title string, amount int, outerNo string, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectQuota
	}
	switch account {
	case member.AccountIntegral:
		return a.integralDiscount(title, amount, outerNo, remark)
	case member.AccountBalance:
		return a.discountBalance(title, float32(amount)/100, outerNo, remark)
	case member.AccountWallet:
		return a.walletDiscount(title, amount, outerNo, remark)
	case member.AccountFlow:
		//return a(title,  float32(amount)/100, outerNo, remark)
	}
	return member.ErrNotSupportAccountType
}

// 创建积分日志
func (a *accountImpl) createIntegralLog(kind int, title string, value int, outerNo string, checkTrade bool) (*member.IntegralLog, error) {
	title = strings.TrimSpace(title)
	if len(title) == 0 {
		return nil, member.ErrNoSuchLogTitleOrRemark
	}
	if checkTrade && len(outerNo) == 0 {
		return nil, member.ErrMissingOuterNo
	}
	unix := time.Now().Unix()
	return &member.IntegralLog{
		Id:          0,
		MemberId:    int(a.value.MemberId),
		Kind:        kind,
		Title:       title,
		OuterNo:     outerNo,
		Value:       value,
		Remark:      "",
		RelateUser:  0,
		ReviewState: int16(enum.ReviewPass),
		CreateTime:  unix,
		UpdateTime:  unix,
	}, nil

}

// 创建余额日志
func (a *accountImpl) createBalanceLog(kind int, title string, amount float32, outerNo string, checkTrade bool) (*member.BalanceLog, error) {
	title = strings.TrimSpace(title)
	outerNo = strings.TrimSpace(outerNo)
	if len(title) == 0 {
		return nil, member.ErrNoSuchLogTitleOrRemark
	}
	if math.IsNaN(float64(amount)) {
		return nil, member.ErrIncorrectAmount
	}
	if checkTrade && len(outerNo) == 0 {
		return nil, member.ErrMissingOuterNo
	}
	unix := time.Now().Unix()
	return &member.BalanceLog{
		MemberId:    a.value.MemberId,
		Kind:        kind,
		Title:       title,
		OuterNo:     outerNo,
		Amount:      amount,
		ReviewState: enum.ReviewPass,
		RelateUser:  0,
		CreateTime:  unix,
		UpdateTime:  unix,
	}, nil

}

// 创建钱包日志
func (a *accountImpl) createWalletLog(kind int, title string, amount float32, outerNo string, checkTrade bool) (*member.WalletAccountLog, error) {
	title = strings.TrimSpace(title)
	outerNo = strings.TrimSpace(outerNo)
	if len(title) == 0 {
		return nil, member.ErrNoSuchLogTitleOrRemark
	}
	if math.IsNaN(float64(amount)) {
		return nil, member.ErrIncorrectAmount
	}
	if checkTrade && len(outerNo) == 0 {
		return nil, member.ErrMissingOuterNo
	}
	unix := time.Now().Unix()
	return &member.WalletAccountLog{
		MemberId:    a.value.MemberId,
		Kind:        kind,
		Title:       title,
		OuterNo:     outerNo,
		Amount:      amount,
		ReviewState: enum.ReviewPass,
		RelateUser:  0,
		CreateTime:  unix,
		UpdateTime:  unix,
	}, nil

}

// 创建活动账户日志
func (a *accountImpl) createFlowAccountLog(kind int, title string, amount float32, outerNo string, checkTrade bool) (*member.FlowAccountLog, error) {
	title = strings.TrimSpace(title)
	outerNo = strings.TrimSpace(outerNo)
	if len(title) == 0 {
		return nil, member.ErrNoSuchLogTitleOrRemark
	}
	if math.IsNaN(float64(amount)) {
		return nil, member.ErrIncorrectAmount
	}
	if checkTrade && len(outerNo) == 0 {
		return nil, member.ErrMissingOuterNo
	}
	unix := time.Now().Unix()
	return &member.FlowAccountLog{
		MemberId:    a.value.MemberId,
		Kind:        kind,
		Title:       title,
		OuterNo:     outerNo,
		Amount:      amount,
		ReviewState: int(enum.ReviewPass),
		RelateUser:  0,
		CreateTime:  unix,
		UpdateTime:  unix,
	}, nil

}

// 充值积分
func (a *accountImpl) integralCharge(title string, value int, outerNo string, remark string) error {
	if value <= 0 {
		return member.ErrIncorrectAmount
	}
	l, err := a.createIntegralLog(member.KindCharge, title, value, outerNo, true)
	if err == nil {
		l.Remark = remark
		err = a.rep.SaveIntegralLog(l)
		if err == nil {
			a.value.Integral += value
			_, err = a.Save()
		}
	}
	return err
}

// 调整账户余额
func (a *accountImpl) adjustBalanceAccount(title string, amount float32, remark string, relateUser int64) error {
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindAdjust,
		Title:       title,
		OuterNo:     "",
		Amount:      amount,
		Remark:      remark,
		RelateUser:  relateUser,
		ReviewState: enum.ReviewPass,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveBalanceLog(v)
	if err == nil {
		a.value.Balance += amount
		_, err = a.Save()
	}
	return err
}

// 充值余额
func (a *accountImpl) chargeBalance(title string, amount float32, outerNo string, remark string) error {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	l, err := a.createBalanceLog(member.KindCharge, title, amount, outerNo, true)
	if err == nil {
		l.Remark = remark
		_, err = a.rep.SaveBalanceLog(l)
		if err == nil {
			a.value.Balance += amount
			_, err = a.Save()
		}
	}
	return err
}

// 充值,客服充值时,需提供操作人(relateUser)
func (a *accountImpl) chargeBalanceNoLimit(kind int, title string, outerNo string,
	amount float32, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:    a.GetDomainId(),
		Kind:        kind,
		Title:       title,
		OuterNo:     outerNo,
		Amount:      amount,
		ReviewState: enum.ReviewPass,
		RelateUser:  relateUser,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveBalanceLog(v)
	if err == nil {
		a.value.Balance += amount
		_, err = a.Save()
	}
	return err
}

// 调整钱包余额
func (a *accountImpl) adjustFlowAccount(title string, amount float32, remark string, relateUser int64) error {
	l, err := a.createFlowAccountLog(member.KindAdjust, title, amount, "", false)
	if err == nil {
		l.Remark = remark
		_, err = a.rep.SaveFlowAccountInfo(l)
		if err == nil {
			a.value.FlowBalance += amount
			_, err = a.Save()
		}
	}
	return err
}

// 调整积分余额
func (a *accountImpl) adjustIntegralAccount(title string, value int, remark string, relateUser int64) error {
	l, err := a.createIntegralLog(member.KindAdjust, title, value, "", false)
	if err == nil {
		l.Remark = remark
		l.RelateUser = int(relateUser)
		err = a.rep.SaveIntegralLog(l)
		if err == nil {
			a.value.Integral += value
			_, err = a.Save()
		}
	}
	return err
}

// 退款
func (a *accountImpl) Refund(account int, title string,
	amount int, outerNo string, remark string) error {
	switch account {
	case member.AccountIntegral:
		return a.integralRefund(title, outerNo, amount, remark)
	case member.AccountBalance:
		return a.chargeBalanceNoLimit(member.KindRefund, title, outerNo, float32(amount)/100, 0)
	case member.AccountWallet:
		//if kind != member.KindRefund &&
		//	kind != member.KindWalletTakeOutRefund {
		//	return member.ErrBusinessKind
		//}
		return a.walletRefund(member.KindRefund, title, outerNo, amount, 1)
	}
	panic(errors.New("不支持的账户类型操作"))
}

func (a *accountImpl) chargeWallet(title string, amount int, outerNo string, remark string) error {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if title == "" {
		title = "钱包账户入账"
	}
	err := a.wallet.Charge(amount, member.KindCharge,
		title, outerNo, remark, 1, "")
	if err == nil {
		a.value.TotalWalletAmount += float32(amount) / 100
		err = a.asyncWallet()
	}
	return err
}

// 钱包消费
func (a *accountImpl) walletConsume(title string, amount int, outerNo string, remark string) error {
	if a.wallet.Get().Balance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	err := a.wallet.Consume(-amount, title, outerNo)
	if err == nil {
		err = a.asyncWallet()
	}
	return err
}

// 扣减奖金,mustLargeZero是否必须大于0, 赠送金额存在扣为负数的情况
func (a *accountImpl) walletDiscount(title string, amount int, outerNo string, remark string) error {
	mustLargeZero := false
	if mustLargeZero && a.wallet.Get().Balance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	err := a.wallet.Discount(-amount, title, outerNo, false)
	if err == nil {
		err = a.asyncWallet()
	}
	return err
}

// 赠送金额(指定业务类型)
func (a *accountImpl) walletRefund(kind int, title string,
	outerNo string, amount int, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if title == "" {
		title = "钱包退款入账"
	}
	err := a.wallet.Refund(amount, kind, title, outerNo, int(relateUser), "")
	if err == nil {
		err = a.asyncWallet()
	}
	return err
}

// 调整钱包余额
func (a *accountImpl) walletAdjust(title string, amount int, remark string, relateUser int64) error {
	err := a.wallet.Adjust(amount, title, "-", remark, int(relateUser), "-")
	if err == nil {
		err = a.asyncWallet()
	}
	return err
}

// 检查交易日志所需要的信息是否完整
func (a *accountImpl) checkTradeLog(amount int, outerNo string) error {
	if math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if len(outerNo) == 0 {
		return member.ErrMissingOuterNo
	}
	return nil
}

// 同步到账户余额
func (a *accountImpl) asyncWallet() error {
	a.value.WalletBalance = float32(a.wallet.Get().Balance) / 100
	_, err := a.Save()
	return err
}

// 根据编号获取余额变动信息
func (a *accountImpl) GetWalletLog(id int64) wallet.WalletLog {
	return a.wallet.GetLog(id)
}

// 冻结余额
func (a *accountImpl) Freeze(title string, outerNo string,
	amount float32, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if a.value.Balance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	if len(title) == 0 {
		title = "资金冻结"
	}
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindFreeze,
		Title:       title,
		Amount:      -amount,
		OuterNo:     outerNo,
		RelateUser:  relateUser,
		ReviewState: enum.ReviewPass,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	a.value.Balance -= amount
	a.value.FreezeBalance += amount
	_, err := a.Save()
	if err == nil {
		_, err = a.rep.SaveBalanceLog(v)
	}
	return err
}

// 解冻金额
func (a *accountImpl) Unfreeze(title string, outerNo string,
	amount float32, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if a.value.FreezeBalance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	if len(title) == 0 {
		title = "资金解结"
	}
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindUnfreeze,
		Title:       title,
		RelateUser:  relateUser,
		Amount:      amount,
		OuterNo:     outerNo,
		ReviewState: enum.ReviewPass,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	a.value.Balance += amount
	a.value.FreezeBalance -= amount
	_, err := a.Save()
	if err == nil {
		_, err = a.rep.SaveBalanceLog(v)
	}
	return err

}

// 冻结赠送金额
func (a *accountImpl) FreezeWallet(title string, outerNo string,
	amount float32, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if a.value.WalletBalance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	if len(title) == 0 {
		title = "(赠送)资金冻结"
	}
	unix := time.Now().Unix()
	v := &member.WalletAccountLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindFreeze,
		Title:       title,
		RelateUser:  relateUser,
		Amount:      -amount,
		OuterNo:     outerNo,
		ReviewState: enum.ReviewPass,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	a.value.WalletBalance -= amount
	a.value.FreezeWallet += amount
	_, err := a.Save()
	if err == nil {
		_, err = a.rep.SaveWalletAccountLog(v)
	}
	return err
}

// 解冻赠送金额
func (a *accountImpl) UnfreezeWallet(title string, outerNo string,
	amount float32, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if a.value.FreezeWallet < amount {
		return member.ErrAccountNotEnoughAmount
	}
	if len(title) == 0 {
		title = "(赠送)资金解冻"
	}
	unix := time.Now().Unix()
	v := &member.WalletAccountLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindUnfreeze,
		Title:       title,
		RelateUser:  relateUser,
		Amount:      amount,
		OuterNo:     outerNo,
		ReviewState: enum.ReviewPass,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	a.value.WalletBalance += amount
	a.value.FreezeWallet -= amount
	_, err := a.Save()
	if err == nil {
		_, err = a.rep.SaveWalletAccountLog(v)
	}
	return err
}

// 流通账户余额充值，如扣除,amount传入负数金额
func (a *accountImpl) chargeFlow(title string, amount float32, outerNo string, remark string) error {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	l, err := a.createFlowAccountLog(member.KindCharge, title, amount, outerNo, true)
	if err == nil {
		l.Remark = remark
		_, err = a.rep.SaveFlowAccountInfo(l)
		if err == nil {
			a.value.FlowBalance += amount
			_, err = a.Save()
		}
	}
	return err
}

// 支付单抵扣消费,tradeNo为支付单单号
func (a *accountImpl) PaymentDiscount(tradeNo string,
	amount float32, remark string) error {
	if amount < 0 || len(tradeNo) == 0 {
		return errors.New("amount error or missing trade no")
	}
	if amount > a.value.Balance {
		return member.ErrOutOfBalance
	}
	if remark == "" {
		remark = "支付抵扣"
	}

	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindDiscount,
		Title:       remark,
		OuterNo:     tradeNo,
		Amount:      -amount,
		ReviewState: enum.ReviewPass,
		RelateUser:  member.DefaultRelateUser,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveBalanceLog(v)
	if err == nil {
		a.value.Balance -= amount
		_, err = a.Save()
	}
	return err
}

// 冻结积分,当new为true不扣除积分,反之扣除积分
func (a *accountImpl) FreezesIntegral(title string, value int, new bool, relateUser int64) error {
	if !new {
		if a.value.Integral < value {
			return member.ErrNoSuchIntegral
		}
		a.value.Integral -= value
	}
	a.value.FreezeIntegral += value
	_, err := a.Save()
	if err == nil {
		unix := time.Now().Unix()
		l := &member.IntegralLog{
			Id:          0,
			MemberId:    int(a.value.MemberId),
			Kind:        member.TypeIntegralFreeze,
			Title:       title,
			OuterNo:     "",
			Value:       -value,
			Remark:      "",
			RelateUser:  int(relateUser),
			ReviewState: int16(enum.ReviewPass),
			CreateTime:  unix,
			UpdateTime:  unix,
		}
		err = a.rep.SaveIntegralLog(l)
	}
	return err
}

// 解冻积分
func (a *accountImpl) UnfreezesIntegral(title string, value int) error {
	if a.value.FreezeIntegral < value {
		return member.ErrNoSuchIntegral
	}
	a.value.FreezeIntegral -= value
	a.value.Integral += value
	_, err := a.Save()
	if err == nil {
		unix := time.Now().Unix()
		var l = &member.IntegralLog{
			Id:          0,
			MemberId:    int(a.value.MemberId),
			Kind:        member.TypeIntegralUnfreeze,
			Title:       title,
			OuterNo:     "",
			Value:       value,
			Remark:      "",
			RelateUser:  0,
			ReviewState: int16(enum.ReviewPass),
			CreateTime:  unix,
			UpdateTime:  unix,
		}
		err = a.rep.SaveIntegralLog(l)
	}
	return err
}

// 请求提现,返回info_id,交易号及错误
func (a *accountImpl) RequestWithdrawal(takeKind int, title string,
	amount int, tradeFee int, accountNo string) (int64, string, error) {
	if takeKind != wallet.KWithdrawExchange &&
		takeKind != wallet.KWithdrawToBankCard &&
		takeKind != wallet.KWithdrawToThirdPart {
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
		trust := a.member.Profile().GetTrustedInfo()
		if trust.ReviewState != int(enum.ReviewPass) {
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
	if a.wallet.Get().Balance < amount {
		return 0, "", member.ErrOutOfBalance
	}
	// 检测提现金额是否超过限制
	minAmountStr, _ := a.registryRepo.GetValue(registry.MemberMinTakeOutAmount)
	minAmount, err := strconv.Atoi(minAmountStr)
	if amount < minAmount*100 {
		return 0, "", errors.New(fmt.Sprintf(member.ErrLessTakeAmount.Error(),
			format.FormatFloat(float32(minAmount))))
	}
	maxAmountStr, _ := a.registryRepo.GetValue(registry.MemberMaxTakeOutAmount)
	maxAmount, err := strconv.Atoi(maxAmountStr)
	if maxAmount > 0 && amount > maxAmount*100 {
		return 0, "", errors.New(fmt.Sprintf(member.ErrOutTakeAmount.Error(),
			format.FormatFloat(float32(maxAmount))))
	}
	// 检测是否超过限制
	maxTimes := a.registryRepo.Get(registry.MemberMaxTakeOutTimesOfDay).IntValue()
	if maxTimes > 0 {
		takeTimes := a.rep.GetTodayTakeOutTimes(a.GetDomainId())
		if takeTimes >= maxTimes {
			return 0, "", member.ErrAccountOutOfTakeOutTimes
		}
	}
	// 检测银行卡
	accountName := ""
	bankName := ""
	if takeKind == wallet.KWithdrawToBankCard {
		if len(accountNo) == 0 {
			return 0, "", errors.New("未指定提现的银行卡号")
		}
		bank := a.member.Profile().GetBankCard(accountNo)
		if bank == nil {
			return 0, "", member.ErrBankNoSuchCard
		}
		accountName = bank.AccountName
		bankName = bank.BankName
	}
	tradeNo := domain.NewTradeNo(8, int(a.member.GetAggregateRootId()))
	finalAmount := amount - tradeFee
	if finalAmount > 0 {
		finalAmount = -finalAmount
	}
	id, tradeNo, err := a.wallet.RequestWithdrawal(finalAmount, tradeFee, takeKind,
		title, accountNo, accountName, bankName)
	if err == nil {
		err = a.asyncWallet()
		if err == nil {
			go a.rep.AddTodayTakeOutTimes(a.GetDomainId())
		}
	}
	return id, tradeNo, err
}

// 确认提现
func (a *accountImpl) ReviewWithdrawal(id int64, pass bool, remark string) error {
	//todo: opr_uid
	err := a.wallet.ReviewWithdrawal(id, pass, remark, 1, "系统")
	if err == nil {
		err = a.asyncWallet()
	}
	return err
}

// 完成提现
func (a *accountImpl) FinishWithdrawal(id int64, tradeNo string) error {
	//todo: opr_uid
	err := a.wallet.FinishWithdrawal(id, tradeNo)
	if err == nil {
		err = a.asyncWallet()
	}
	return err
}

// 将冻结金额标记为失效
func (a *accountImpl) FreezeExpired(accountKind int, amount float32, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	switch accountKind {
	case member.AccountBalance:
		return a.balanceFreezeExpired(amount, remark)
	case member.AccountWallet:
		return a.presentFreezeExpired(amount, remark)
	}
	return nil
}

func (a *accountImpl) balanceFreezeExpired(amount float32, remark string) error {
	if a.value.FreezeBalance < amount {
		return member.ErrIncorrectAmount
	}
	unix := time.Now().Unix()
	a.value.FreezeBalance -= amount
	a.value.ExpiredBalance += amount
	a.value.UpdateTime = unix
	l := &member.BalanceLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindExpired,
		Title:       "过期失效",
		OuterNo:     "",
		Amount:      amount,
		CsnFee:      0,
		ReviewState: enum.ReviewPass,
		RelateUser:  member.DefaultRelateUser,
		Remark:      remark,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveBalanceLog(l)
	if err == nil {
		_, err = a.Save()
	}
	return err
}

func (a *accountImpl) presentFreezeExpired(amount float32, remark string) error {
	if a.value.FreezeWallet < amount {
		return member.ErrIncorrectAmount
	}
	unix := time.Now().Unix()
	a.value.FreezeWallet -= amount
	a.value.ExpiredWallet += amount
	a.value.UpdateTime = unix
	l := &member.WalletAccountLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindExpired,
		Title:       "过期失效",
		OuterNo:     "",
		Amount:      amount,
		CsnFee:      0,
		ReviewState: enum.ReviewPass,
		RelateUser:  member.DefaultRelateUser,
		Remark:      remark,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveWalletAccountLog(l)
	if err == nil {
		_, err = a.Save()
	}
	return err
}

// 获取会员名称
func (a *accountImpl) getMemberName(m member.IMember) string {
	if tr := m.Profile().GetTrustedInfo(); tr.RealName != "" &&
		tr.ReviewState == int(enum.ReviewPass) {
		return tr.RealName
	} else {
		return m.GetValue().User
	}
}

// 转账
func (a *accountImpl) TransferAccount(accountKind int, toMember int64, amount int,
	tradeFee int, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	tm := a.rep.GetMember(toMember)
	if tm == nil {
		return member.ErrNoSuchMember
	}

	tradeNo := domain.NewTradeNo(8, int(a.member.GetAggregateRootId()))

	// 检测是否开启转账
	transferOn := a.registryRepo.Get(registry.MemberTransferAccountsOn).BoolValue()
	if !transferOn {
		msg := a.registryRepo.Get(registry.MemberTransferAccountsMessage).StringValue()
		return errors.New(msg)
	}

	switch accountKind {
	case member.AccountWallet:
		return a.transferWalletAccount(tm, tradeNo, amount, tradeFee, remark)
	case member.AccountBalance:
		return a.transferBalance(tm, tradeNo, amount, tradeFee, remark)
	}
	return nil
}

func (a *accountImpl) transferBalance(tm member.IMember, tradeNo string,
	tradeAmount, tradeFee int, remark string) error {
	csnFee := float32(tradeFee) / 100
	amount := float32(tradeAmount) / 100
	if a.value.Balance < amount+csnFee {
		return member.ErrAccountNotEnoughAmount
	}
	unix := time.Now().Unix()
	// 扣款
	toName := a.getMemberName(tm)
	l := &member.BalanceLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindTransferOut,
		Title:       "转账给" + toName,
		OuterNo:     tradeNo,
		Amount:      -amount,
		CsnFee:      csnFee,
		ReviewState: enum.ReviewPass,
		RelateUser:  member.DefaultRelateUser,
		Remark:      remark,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveBalanceLog(l)
	if err == nil {
		a.value.Balance -= amount + csnFee
		a.value.UpdateTime = unix
		_, err = a.Save()
		if err == nil {
			err = tm.GetAccount().ReceiveTransfer(member.AccountBalance,
				a.GetDomainId(), tradeNo, amount, remark)
		}
	}
	return err
}

func (a *accountImpl) transferWalletAccount(tm member.IMember, tradeNo string,
	tradeAmount, tradeFee int, remark string) error {
	csnFee := float32(tradeFee) / 100
	amount := float32(tradeAmount) / 100
	// 检测非正式会员转账
	lv := a.mm.LevelManager().GetLevelById(a.member.GetValue().Level)
	if lv != nil && lv.IsOfficial == 0 {
		return errors.New(fmt.Sprintf(
			member.ErrTransferAccountSMemberLevelNoPerm.Error(), lv.Name))
	}
	if a.value.WalletBalance < amount+csnFee {
		return member.ErrAccountNotEnoughAmount
	}
	unix := time.Now().Unix()
	// 扣款
	toName := a.getMemberName(tm)
	l := &member.WalletAccountLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindTransferOut,
		Title:       "转账给" + toName,
		OuterNo:     tradeNo,
		Amount:      -amount,
		CsnFee:      csnFee,
		ReviewState: enum.ReviewPass,
		RelateUser:  member.DefaultRelateUser,
		Remark:      remark,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveWalletAccountLog(l)
	if err == nil {
		a.value.WalletBalance -= amount + csnFee
		a.value.UpdateTime = unix
		_, err = a.Save()
		if err == nil {
			err = tm.GetAccount().ReceiveTransfer(member.AccountWallet,
				a.GetDomainId(), tradeNo, amount, remark)
		}
	}
	return err
}

// 接收转账
func (a *accountImpl) ReceiveTransfer(accountKind int, fromMember int64,
	tradeNo string, amount float32, remark string) error {
	switch accountKind {
	case member.AccountWallet:
		return a.receivePresentTransfer(fromMember, tradeNo, amount, remark)
	case member.AccountBalance:
		return a.receiveBalanceTransfer(fromMember, tradeNo, amount, remark)
	}
	return member.ErrNotSupportTransfer
}

func (a *accountImpl) receivePresentTransfer(fromMember int64, tradeNo string,
	amount float32, remark string) error {
	fm := a.rep.GetMember(fromMember)
	if fm == nil {
		return member.ErrNoSuchMember
	}
	fromName := a.getMemberName(fm)
	unix := time.Now().Unix()
	tl := &member.WalletAccountLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindTransferIn,
		Title:       "转账收款（" + fromName + "）",
		OuterNo:     tradeNo,
		Amount:      amount,
		CsnFee:      0,
		ReviewState: enum.ReviewPass,
		RelateUser:  member.DefaultRelateUser,
		Remark:      remark,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveWalletAccountLog(tl)
	if err == nil {
		a.value.WalletBalance += amount
		a.value.UpdateTime = unix
		_, err = a.Save()
	}
	return err
}

func (a *accountImpl) receiveBalanceTransfer(fromMember int64, tradeNo string,
	amount float32, remark string) error {
	fromName := a.getMemberName(a.rep.GetMember(a.GetDomainId()))
	unix := time.Now().Unix()
	tl := &member.BalanceLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindTransferIn,
		Title:       "转账收款（" + fromName + "）",
		OuterNo:     tradeNo,
		Amount:      amount,
		CsnFee:      0,
		ReviewState: enum.ReviewPass,
		RelateUser:  member.DefaultRelateUser,
		Remark:      remark,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveBalanceLog(tl)
	if err == nil {
		a.value.Balance += amount
		a.value.UpdateTime = unix
		_, err = a.Save()
	}
	return err
}

// 转账余额到其他账户
func (a *accountImpl) TransferBalance(kind int, amount float32,
	tradeNo string, toTitle, fromTitle string) error {
	var err error
	if kind == member.AccountFlow {
		if a.value.Balance < amount {
			return member.ErrAccountNotEnoughAmount
		}
		a.value.Balance -= amount
		a.value.FlowBalance += amount
		if _, err = a.Save(); err == nil {
			a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
				Kind:        member.KindTransferOut,
				Title:       toTitle,
				Amount:      -amount,
				OuterNo:     tradeNo,
				ReviewState: int(enum.ReviewPass),
			})

			a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
				Kind:        member.KindTransferIn,
				Title:       fromTitle,
				Amount:      amount,
				OuterNo:     tradeNo,
				ReviewState: int(enum.ReviewPass),
			})
		}
		return err
	}
	return member.ErrNotSupportTransfer
}

// 转账活动账户,kind为转账类型，如 KindBalanceTransfer等
// commission手续费
func (a *accountImpl) TransferFlow(kind int, amount float32, commission float32,
	tradeNo string, toTitle string, fromTitle string) error {
	var err error

	csnAmount := commission * amount
	finalAmount := amount - csnAmount

	if a.value.FlowBalance < finalAmount {
		return member.ErrAccountNotEnoughAmount
	}

	a.value.FlowBalance -= amount
	a.value.WalletBalance += finalAmount
	a.value.TotalWalletAmount += finalAmount

	if _, err = a.Save(); err == nil {
		a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
			Kind:        member.KindTransferOut,
			Title:       toTitle,
			Amount:      -amount,
			OuterNo:     tradeNo,
			CsnFee:      csnAmount,
			ReviewState: int(enum.ReviewPass),
		})

		a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
			Kind:        member.KindTransferIn,
			Title:       fromTitle,
			Amount:      finalAmount,
			OuterNo:     tradeNo,
			ReviewState: int(enum.ReviewPass),
		})
	}
	return err
}

// 将活动金转给其他人
func (a *accountImpl) TransferFlowTo(memberId int64, kind int,
	amount float32, commission float32, tradeNo string,
	toTitle string, fromTitle string) error {

	var err error
	csnAmount := commission * amount
	finalAmount := amount + csnAmount // 转账方付手续费

	m := a.rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	acc2 := m.GetAccount()

	if a.value.FlowBalance < finalAmount {
		return member.ErrAccountNotEnoughAmount
	}

	a.value.FlowBalance -= finalAmount
	acc2.GetValue().FlowBalance += amount

	if _, err = a.Save(); err == nil {
		_, err = a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
			Kind:        member.KindTransferOut,
			Title:       toTitle,
			Amount:      -finalAmount,
			CsnFee:      csnAmount,
			RelateUser:  memberId,
			OuterNo:     tradeNo,
			ReviewState: member.StatusOK,
		})

		if _, err = acc2.Save(); err == nil {
			_, err = a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
				Kind:        member.KindTransferIn,
				Title:       fromTitle,
				Amount:      amount,
				RelateUser:  a.value.MemberId,
				OuterNo:     tradeNo,
				ReviewState: member.StatusOK,
			})
		}
	}
	return err
}

// 充值积分
func (a *accountImpl) integralConsume(title string, value int, outerNo string, remark string) error {
	if a.value.Integral < value {
		return member.ErrNoSuchIntegral
	}
	l, err := a.createIntegralLog(member.KindConsume, title, -value, outerNo, true)
	if err == nil {
		l.Remark = remark
		err = a.rep.SaveIntegralLog(l)
		if err == nil {
			a.value.Integral -= value
			_, err = a.Save()
		}
	}
	return err
}

func (a *accountImpl) balanceConsume(title string, amount float32, outerNo string, remark string) error {
	if a.value.Balance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	l, err := a.createBalanceLog(member.KindConsume, title, -amount, outerNo, true)
	if err == nil {
		l.Remark = remark
		_, err = a.rep.SaveBalanceLog(l)
		if err == nil {
			a.value.Balance -= amount
			_, err = a.Save()
		}
	}
	return err
}

func (a *accountImpl) flowAccountConsume(title string, amount float32, outerNo string, remark string) error {
	if a.value.FlowBalance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	l, err := a.createFlowAccountLog(member.KindConsume, title, -amount, outerNo, true)
	if err == nil {
		l.Remark = remark
		_, err = a.rep.SaveFlowAccountInfo(l)
		if err == nil {
			a.value.FlowBalance -= amount
			_, err = a.Save()
		}
	}
	return err
}

// 积分抵扣
func (a *accountImpl) integralDiscount(title string, value int, outerNo string, remark string) error {
	if a.value.Integral < value {
		return member.ErrNoSuchIntegral
	}
	l, err := a.createIntegralLog(member.KindDiscount, title, -value, outerNo, true)
	if err == nil {
		l.Remark = remark
		if err == nil {
			a.value.Integral -= value
			_, err = a.Save()
		}
	}
	return err
}

// 扣减余额
func (a *accountImpl) discountBalance(title string, amount float32, outerNo string, remark string) (err error) {
	if a.value.Balance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	l, err := a.createBalanceLog(member.KindDiscount, title, -amount, outerNo, true)
	if err == nil {
		l.Remark = remark
		if err == nil {
			a.value.Balance -= amount
			_, err = a.Save()
		}
	}
	return err
}

// 充值积分
func (a *accountImpl) integralRefund(title string, outerNo string,
	value int, remark string) error {
	if value <= 0 {
		return member.ErrIncorrectAmount
	}
	l, err := a.createIntegralLog(member.KindRefund, title, value, outerNo, true)
	if err == nil {
		l.Remark = remark
		err = a.rep.SaveIntegralLog(l)
		if err == nil {
			a.value.Integral += value
			_, err = a.Save()
		}
	}
	return err
}
