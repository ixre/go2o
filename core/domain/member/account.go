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
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/go2o/core/msq"
	"log"
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


func newAccount(m *memberImpl, value *member.Account,
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
		impl.initWallet()
	}
	return impl
}

func (a *accountImpl) initWallet() {
	flag := wallet.FlagCharge | wallet.FlagDiscount
	// 存在钱包,但关联失败
	a.wallet = a.walletRepo.GetWalletByUserId(a.member.GetAggregateRootId(), 1)
	if a.wallet != nil {
		a.value.WalletCode = a.wallet.Get().HashCode
		a.Save()
		return
	}
	//　创建新的钱包
	a.wallet = a.walletRepo.CreateWallet(
		a.member.GetAggregateRootId(),
		a.member.value.User,
		1, "MemberWallet", flag)
	if _, err := a.wallet.Save(); err != nil {
		log.Println("[ go2o][ member]: create wallet failed,error", err.Error())
	}
	a.value.WalletCode = a.wallet.Get().HashCode 		// 绑定钱包
}

// GetDomainId 获取领域对象编号
func (a *accountImpl) GetDomainId() int64 {
	return a.value.MemberId
}

// GetValue 获取账户值
func (a *accountImpl) GetValue() *member.Account {
	return a.value
}

// Save 保存
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
			a.initWallet()
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

// SetPriorityPay 设置优先(默认)支付方式, account 为账户类型
func (a *accountImpl) SetPriorityPay(account member.AccountType, enabled bool) error {
	if enabled {
		support := false
		if account == member.AccountBalance ||
			account == member.AccountWallet ||
			account == member.AccountIntegral {
			support = true
		}
		if support {
			a.value.PriorityPay = int(account)
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

// Charge 充值
func (a *accountImpl) Charge(account member.AccountType, title string,
	amount int, outerNo string, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectQuota
	}
	switch account {
	case member.AccountIntegral:
		return a.integralCharge(title, amount, outerNo, remark)
	case member.AccountBalance:
		return a.chargeBalance(title, amount, outerNo, remark)
	case member.AccountWallet:
		return a.chargeWallet(title, amount, outerNo, remark)
	case member.AccountFlow:
		return a.chargeFlow(title, amount, outerNo, remark)
	}
	return member.ErrNotSupportAccountType
}

func (a *accountImpl) Adjust(account member.AccountType, title string, amount int, remark string, relateUser int64) error {
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
		return a.adjustBalanceAccount(title, amount, remark, relateUser)
	case member.AccountWallet:
		return a.walletAdjust(title, amount, remark, relateUser)
	case member.AccountFlow:
		return a.adjustFlowAccount(title, amount, remark, relateUser)
	}
	panic("not support other account adjust")
}

// 消耗
func (a *accountImpl) Consume(account member.AccountType, title string, amount int, outerNo string, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectQuota
	}
	switch account {
	case member.AccountIntegral:
		return a.integralConsume(title, amount, outerNo, remark)
	case member.AccountBalance:
		return a.balanceConsume(title, amount, outerNo, remark)
	case member.AccountWallet:
		return a.walletConsume(title, amount, outerNo, remark)
	case member.AccountFlow:
		return a.flowAccountConsume(title, amount, outerNo, remark)
	}
	return member.ErrNotSupportAccountType
}

func (a *accountImpl) Discount(account member.AccountType, title string, amount int, outerNo string, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectQuota
	}
	switch account {
	case member.AccountIntegral:
		return a.integralDiscount(title, amount, outerNo, remark)
	case member.AccountBalance:
		return a.discountBalance(title, amount, outerNo, remark)
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
func (a *accountImpl) createBalanceLog(kind int, title string, amount int, outerNo string, checkTrade bool) (*member.BalanceLog, error) {
	title = strings.TrimSpace(title)
	outerNo = strings.TrimSpace(outerNo)
	if len(title) == 0 {
		return nil, member.ErrNoSuchLogTitleOrRemark
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
		Amount:      int64(amount),
		ReviewState: enum.ReviewPass,
		RelateUser:  0,
		CreateTime:  unix,
		UpdateTime:  unix,
	}, nil

}

// 创建钱包日志
func (a *accountImpl) createWalletLog(kind int, title string, amount int, outerNo string, checkTrade bool) (*member.WalletAccountLog, error) {
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
		Amount:      int64(amount),
		ReviewState: enum.ReviewPass,
		RelateUser:  0,
		CreateTime:  unix,
		UpdateTime:  unix,
	}, nil

}

// 创建活动账户日志
func (a *accountImpl) createFlowAccountLog(kind int, title string, amount int, outerNo string, checkTrade bool) (*member.FlowAccountLog, error) {
	title = strings.TrimSpace(title)
	outerNo = strings.TrimSpace(outerNo)
	if len(title) == 0 {
		return nil, member.ErrNoSuchLogTitleOrRemark
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
		Amount:      int64(amount),
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
func (a *accountImpl) adjustBalanceAccount(title string, amount int, remark string, relateUser int64) error {
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindAdjust,
		Title:       title,
		OuterNo:     "",
		Amount:      int64(amount),
		Remark:      remark,
		RelateUser:  relateUser,
		ReviewState: enum.ReviewPass,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveBalanceLog(v)
	if err == nil {
		a.value.Balance += int64(amount)
		_, err = a.Save()
	}
	return err
}

// 充值余额
func (a *accountImpl) chargeBalance(title string, amount int, outerNo string, remark string) error {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	l, err := a.createBalanceLog(member.KindCharge, title, amount, outerNo, true)
	if err == nil {
		l.Remark = remark
		_, err = a.rep.SaveBalanceLog(l)
		if err == nil {
			a.value.Balance += int64(amount)
			_, err = a.Save()
		}
	}
	return err
}

// 充值,客服充值时,需提供操作人(relateUser)
func (a *accountImpl) chargeBalanceNoLimit(kind int, title string, outerNo string,
	amount int, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:    a.GetDomainId(),
		Kind:        kind,
		Title:       title,
		OuterNo:     outerNo,
		Amount:      int64(amount),
		ReviewState: enum.ReviewPass,
		RelateUser:  relateUser,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveBalanceLog(v)
	if err == nil {
		a.value.Balance += int64(amount)
		_, err = a.Save()
	}
	return err
}

// 调整钱包余额
func (a *accountImpl) adjustFlowAccount(title string, amount int, remark string, relateUser int64) error {
	l, err := a.createFlowAccountLog(member.KindAdjust, title, amount, "", false)
	if err == nil {
		l.Remark = remark
		_, err = a.rep.SaveFlowAccountInfo(l)
		if err == nil {
			a.value.FlowBalance += int64(amount)
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

// Refund 账户退款
func (a *accountImpl) Refund(account member.AccountType, title string,
	amount int, outerNo string, remark string) error {
	switch account {
	case member.AccountIntegral:
		return a.integralRefund(title, outerNo, amount, remark)
	case member.AccountBalance:
		return a.chargeBalanceNoLimit(member.KindRefund, title, outerNo, amount, 0)
	case member.AccountWallet:
		//if kind != member.KindRefund &&
		//	kind != member.KindWalletTakeOutRefund {
		//	return member.ErrBusinessKind
		//}
		return a.walletRefund(wallet.KPaymentOrderRefund, title, outerNo, amount, 1)
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
		a.value.TotalWalletAmount += int64(amount)
		err = a.asyncWallet()
	}
	return err
}

// 钱包消费
func (a *accountImpl) walletConsume(title string, amount int, outerNo string, remark string) error {
	if a.wallet.Get().Balance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	err := a.wallet.Consume(-amount, title, outerNo, remark)
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

//  钱包退款(指定业务类型)
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
	a.value.WalletBalance = int64(a.wallet.Get().Balance)
	_, err := a.Save()
	return err
}

// GetWalletLog 根据编号获取余额变动信息
func (a *accountImpl) GetWalletLog(id int64) wallet.WalletLog {
	return a.wallet.GetLog(id)
}

// Freeze 冻结余额
func (a *accountImpl) Freeze(account member.AccountType,p member.AccountOperateData, relateUser int64) error {
	switch account {
	case member.AccountBalance:
		return a.freezeBalance(p,relateUser)
	case member.AccountWallet:
		return a.freezeWallet(p,relateUser)
	case member.AccountIntegral:
		return a.freezesIntegral(p,relateUser)
	}
	panic("not support account type")
}


// Freeze 冻结余额
func (a *accountImpl) freezeBalance(p member.AccountOperateData, relateUser int64) error {
	if p.Amount <= 0 || math.IsNaN(float64(p.Amount)) {
		return member.ErrIncorrectAmount
	}
	if a.value.Balance < int64(p.Amount) {
		return member.ErrAccountNotEnoughAmount
	}
	if len(p.Title) == 0 {
		p.Title = "资金冻结"
	}
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindFreeze,
		Title:       p.Title,
		Amount:      -int64(p.Amount),
		OuterNo:     p.OuterNo,
		RelateUser:  relateUser,
		Remark: p.Remark,
		ReviewState: enum.ReviewPass,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	a.value.Balance -= int64(p.Amount)
	a.value.FreezeBalance += int64(p.Amount)
	_, err := a.Save()
	if err == nil {
		_, err = a.rep.SaveBalanceLog(v)
	}
	return err
}
// FreezeWallet 冻结钱包
func (a *accountImpl) freezeWallet(p member.AccountOperateData, relateUser int64) error {
	err := a.wallet.Freeze(wallet.OperateData{
		Title:  p.Title,
		Amount:  p.Amount,
		OuterNo: p.OuterNo,
		Remark:  p.Remark,
	},wallet.Operator{
		OperatorUid:  int(relateUser),
		OperatorName: "",
	})
	if err == nil{
		err = a.asyncWallet()
	}
	return err

	//if p.Amount <= 0 || math.IsNaN(float64(p.Amount)) {
	//	return member.ErrIncorrectAmount
	//}
	//if a.value.WalletBalance < int64( p.Amount) {
	//	return member.ErrAccountNotEnoughAmount
	//}
	//if len(p.Title) == 0 {
	//	p.Title = "(赠送)资金冻结"
	//}
	//unix := time.Now().Unix()
	//v := &member.WalletAccountLog{
	//	MemberId:    a.GetDomainId(),
	//	Kind:        member.KindFreeze,
	//	Title:       p.Title,
	//	RelateUser:  relateUser,
	//	Amount:      -int64( p.Amount),
	//	OuterNo:     p.OuterNo,
	//	Remark: p.Remark,
	//	ReviewState: enum.ReviewPass,
	//	CreateTime:  unix,
	//	UpdateTime:  unix,
	//}
	//a.value.WalletBalance -= int64( p.Amount)
	//a.value.FreezeWallet += int64( p.Amount)
	//_, err := a.Save()
	//if err == nil {
	//	_, err = a.rep.SaveWalletAccountLog(v)
	//}
	//return err
}
// 冻结余额
func (a *accountImpl) freezesIntegral(p member.AccountOperateData, relateUser int64) error {
	//if !new {
	//	if a.value.Integral < value {
	//		return member.ErrNoSuchIntegral
	//	}
	//	a.value.Integral -= value
	//}
	a.value.FreezeIntegral += p.Amount
	_, err := a.Save()
	if err == nil {
		unix := time.Now().Unix()
		l := &member.IntegralLog{
			Id:          0,
			MemberId:    int(a.value.MemberId),
			Kind:        member.TypeIntegralFreeze,
			Title:       p.Title,
			OuterNo:     p.OuterNo,
			Value:       -p.Amount,
			Remark:      p.Remark,
			RelateUser:  int(relateUser),
			ReviewState: int16(enum.ReviewPass),
			CreateTime:  unix,
			UpdateTime:  unix,
		}
		err = a.rep.SaveIntegralLog(l)
	}
	return err
}

// Unfreeze 解冻金额
func (a *accountImpl) Unfreeze(title string, outerNo string, amount int, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if a.value.FreezeBalance < int64(amount) {
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
		Amount:      int64(amount),
		OuterNo:     outerNo,
		ReviewState: enum.ReviewPass,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	a.value.Balance += int64(amount)
	a.value.FreezeBalance -= int64(amount)
	_, err := a.Save()
	if err == nil {
		_, err = a.rep.SaveBalanceLog(v)
	}
	return err

}


// UnfreezeWallet 解冻赠送金额
func (a *accountImpl) UnfreezeWallet(title string, outerNo string,
	amount int, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if a.value.FreezeWallet < int64(amount) {
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
		Amount:      int64(amount),
		OuterNo:     outerNo,
		ReviewState: enum.ReviewPass,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	a.value.WalletBalance += int64(amount)
	a.value.FreezeWallet -= int64(amount)
	_, err := a.Save()
	if err == nil {
		_, err = a.rep.SaveWalletAccountLog(v)
	}
	return err
}

// 流通账户余额充值，如扣除,amount传入负数金额
func (a *accountImpl) chargeFlow(title string, amount int, outerNo string, remark string) error {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	l, err := a.createFlowAccountLog(member.KindCharge, title, amount, outerNo, true)
	if err == nil {
		l.Remark = remark
		_, err = a.rep.SaveFlowAccountInfo(l)
		if err == nil {
			a.value.FlowBalance += int64(amount)
			_, err = a.Save()
		}
	}
	return err
}

// PaymentDiscount 支付单抵扣消费,tradeNo为支付单单号
func (a *accountImpl) PaymentDiscount(tradeNo string, amount int, remark string) error {
	if amount < 0 || len(tradeNo) == 0 {
		return errors.New("amount error or missing trade no")
	}
	if int64(amount) > a.value.Balance {
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
		Amount:      -int64(amount),
		ReviewState: enum.ReviewPass,
		RelateUser:  member.DefaultRelateUser,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveBalanceLog(v)
	if err == nil {
		a.value.Balance -= int64(amount)
		_, err = a.Save()
	}
	return err
}

// UnfreezesIntegral 解冻积分
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

// RequestWithdrawal 请求提现,返回info_id,交易号及错误
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
	takOutOn := a.registryRepo.Get(registry.MemberWithdrawEnabled).BoolValue()
	if !takOutOn {
		msg, _ := a.registryRepo.GetValue(registry.MemberWithdrawMessage)
		return 0, "", errors.New(msg)
	}
	// 检测是否实名
	mustTrust, _ := a.registryRepo.GetValue(registry.MemberWithdrawalMustVerification)
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
	minAmountStr, _ := a.registryRepo.GetValue(registry.MemberWithdrawMinAmount)
	minAmount, err := strconv.Atoi(minAmountStr)
	if amount < minAmount*100 {
		return 0, "", errors.New(fmt.Sprintf(member.ErrLessTakeAmount.Error(),
			format.FormatFloat(float32(minAmount))))
	}
	maxAmountStr, _ := a.registryRepo.GetValue(registry.MemberWithdrawMaxAmount)
	maxAmount, err := strconv.Atoi(maxAmountStr)
	if maxAmount > 0 && amount > maxAmount*100 {
		return 0, "", errors.New(fmt.Sprintf(member.ErrOutTakeAmount.Error(),
			format.FormatFloat(float32(maxAmount))))
	}
	// 检测是否超过限制
	maxTimes := a.registryRepo.Get(registry.MemberWithdrawMaxTimeOfDay).IntValue()
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

// ReviewWithdrawal 确认提现
func (a *accountImpl) ReviewWithdrawal(id int64, pass bool, remark string) error {
	//todo: opr_uid
	err := a.wallet.ReviewWithdrawal(id, pass, remark, 1, "系统")
	if err == nil {
		err = a.asyncWallet()
	}
	return err
}

// FinishWithdrawal 完成提现
func (a *accountImpl) FinishWithdrawal(id int64, tradeNo string) error {
	//todo: opr_uid
	err := a.wallet.FinishWithdrawal(id, tradeNo)
	if err == nil {
		err = a.asyncWallet()
	}
	return err
}

// FreezeExpired 将冻结金额标记为失效
func (a *accountImpl) FreezeExpired(account member.AccountType, amount int, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	switch account {
	case member.AccountBalance:
		return a.balanceFreezeExpired(amount, remark)
	case member.AccountWallet:
		return a.walletFreezeExpired(amount, remark)
	}
	panic("not support account type")
}

func (a *accountImpl) balanceFreezeExpired(amount int, remark string) error {
	if a.value.FreezeBalance < int64(amount) {
		return member.ErrIncorrectAmount
	}
	unix := time.Now().Unix()
	a.value.FreezeBalance -= int64(amount)
	a.value.ExpiredBalance += int64(amount)
	a.value.UpdateTime = unix
	l := &member.BalanceLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindExpired,
		Title:       "过期失效",
		OuterNo:     "",
		Amount:      int64(amount),
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

func (a *accountImpl) walletFreezeExpired(amount int, remark string) error {
	if a.value.FreezeWallet < int64(amount) {
		return member.ErrIncorrectAmount
	}
	unix := time.Now().Unix()
	a.value.FreezeWallet -= int64(amount)
	a.value.ExpiredWallet += int64(amount)
	a.value.UpdateTime = unix
	l := &member.WalletAccountLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindExpired,
		Title:       "过期失效",
		OuterNo:     "",
		Amount:      int64(amount),
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

// TransferAccount 转账
func (a *accountImpl) TransferAccount(account member.AccountType, toMember int64, amount int,
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
	transferOn := a.registryRepo.Get(registry.MemberAccountTransferEnabled).BoolValue()
	if !transferOn {
		msg := a.registryRepo.Get(registry.MemberAccountTransferMessage).StringValue()
		return errors.New(msg)
	}

	switch account {
	case member.AccountWallet:
		return a.transferWalletAccount(tm, tradeNo, amount, tradeFee, remark)
	case member.AccountBalance:
		return a.transferBalance(tm, tradeNo, amount, tradeFee, remark)
	}
	return nil
}

func (a *accountImpl) transferBalance(tm member.IMember, tradeNo string,
	tradeAmount, tradeFee int, remark string) error {
	csnFee := tradeFee
	amount := tradeAmount
	if a.value.Balance < int64(amount+csnFee) {
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
		Amount:      -int64(amount),
		CsnFee:      int64(csnFee),
		ReviewState: enum.ReviewPass,
		RelateUser:  member.DefaultRelateUser,
		Remark:      remark,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveBalanceLog(l)
	if err == nil {
		a.value.Balance -= int64(amount + csnFee)
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
	csnFee := tradeFee
	amount := tradeAmount
	// 检测非正式会员转账
	lv := a.mm.LevelManager().GetLevelById(a.member.GetValue().Level)
	if lv != nil && lv.IsOfficial == 0 {
		return errors.New(fmt.Sprintf(
			member.ErrTransferAccountSMemberLevelNoPerm.Error(), lv.Name))
	}
	if a.value.WalletBalance < int64(amount+csnFee) {
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
		Amount:      -int64(amount),
		CsnFee:      int64(csnFee),
		ReviewState: enum.ReviewPass,
		RelateUser:  member.DefaultRelateUser,
		Remark:      remark,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveWalletAccountLog(l)
	if err == nil {
		a.value.WalletBalance -= int64(amount + csnFee)
		a.value.UpdateTime = unix
		_, err = a.Save()
		if err == nil {
			err = tm.GetAccount().ReceiveTransfer(member.AccountWallet,
				a.GetDomainId(), tradeNo, amount, remark)
		}
	}
	return err
}

// ReceiveTransfer 接收转账
func (a *accountImpl) ReceiveTransfer(account member.AccountType, fromMember int64,
	tradeNo string, amount int, remark string) error {
	switch account {
	case member.AccountWallet:
		return a.receivePresentTransfer(fromMember, tradeNo, amount, remark)
	case member.AccountBalance:
		return a.receiveBalanceTransfer(fromMember, tradeNo, amount, remark)
	}
	return member.ErrNotSupportTransfer
}

func (a *accountImpl) receivePresentTransfer(fromMember int64, tradeNo string,
	amount int, remark string) error {
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
		Amount:      int64(amount),
		CsnFee:      0,
		ReviewState: enum.ReviewPass,
		RelateUser:  member.DefaultRelateUser,
		Remark:      remark,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveWalletAccountLog(tl)
	if err == nil {
		a.value.WalletBalance += int64(amount)
		a.value.UpdateTime = unix
		_, err = a.Save()
	}
	return err
}

func (a *accountImpl) receiveBalanceTransfer(fromMember int64, tradeNo string,
	amount int, remark string) error {
	fromName := a.getMemberName(a.rep.GetMember(a.GetDomainId()))
	unix := time.Now().Unix()
	tl := &member.BalanceLog{
		MemberId:    a.GetDomainId(),
		Kind:        member.KindTransferIn,
		Title:       "转账收款（" + fromName + "）",
		OuterNo:     tradeNo,
		Amount:      int64(amount),
		CsnFee:      0,
		ReviewState: enum.ReviewPass,
		RelateUser:  member.DefaultRelateUser,
		Remark:      remark,
		CreateTime:  unix,
		UpdateTime:  unix,
	}
	_, err := a.rep.SaveBalanceLog(tl)
	if err == nil {
		a.value.Balance += int64(amount)
		a.value.UpdateTime = unix
		_, err = a.Save()
	}
	return err
}

// TransferBalance 转账余额到其他账户
func (a *accountImpl) TransferBalance(account member.AccountType, amount int,
	tradeNo string, toTitle, fromTitle string) error {
	var err error
	if account == member.AccountFlow {
		if a.value.Balance < int64(amount) {
			return member.ErrAccountNotEnoughAmount
		}
		a.value.Balance -= int64(amount)
		a.value.FlowBalance += int64(amount)
		if _, err = a.Save(); err == nil {
			a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
				Kind:        member.KindTransferOut,
				Title:       toTitle,
				Amount:      -int64(amount),
				OuterNo:     tradeNo,
				ReviewState: int(enum.ReviewPass),
			})

			a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
				Kind:        member.KindTransferIn,
				Title:       fromTitle,
				Amount:      int64(amount),
				OuterNo:     tradeNo,
				ReviewState: int(enum.ReviewPass),
			})
		}
		return err
	}
	return member.ErrNotSupportTransfer
}

// TransferFlow TransferFlow 转账活动账户,kind为转账类型，如 KindBalanceTransfer等
// commission手续费
func (a *accountImpl) TransferFlow(kind int, amount int, commission float32,
	tradeNo string, toTitle string, fromTitle string) error {
	var err error

	csnAmount := int64(commission * float32(amount))
	finalAmount := int64(amount) - csnAmount

	if a.value.FlowBalance < int64(finalAmount) {
		return member.ErrAccountNotEnoughAmount
	}

	a.value.FlowBalance -= int64(amount)
	a.value.WalletBalance += int64(finalAmount)
	a.value.TotalWalletAmount += int64(finalAmount)

	if _, err = a.Save(); err == nil {
		a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
			Kind:        kind,
			Title:       toTitle,
			Amount:      -int64(amount),
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

// TransferFlowTo 将活动金转给其他人
func (a *accountImpl) TransferFlowTo(memberId int64, kind int,
	amount int, commission float32, tradeNo string,
	toTitle string, fromTitle string) error {

	var err error
	csnAmount := int64(commission * float32(amount))
	finalAmount := int64(amount) + csnAmount // 转账方付手续费

	m := a.rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	acc2 := m.GetAccount()

	if a.value.FlowBalance < finalAmount {
		return member.ErrAccountNotEnoughAmount
	}

	a.value.FlowBalance -= finalAmount
	acc2.GetValue().FlowBalance += int64(amount)

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
				Amount:      int64(amount),
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

func (a *accountImpl) balanceConsume(title string, amount int, outerNo string, remark string) error {
	if a.value.Balance < int64(amount) {
		return member.ErrAccountNotEnoughAmount
	}
	l, err := a.createBalanceLog(member.KindConsume, title, -amount, outerNo, true)
	if err == nil {
		l.Remark = remark
		_, err = a.rep.SaveBalanceLog(l)
		if err == nil {
			a.value.Balance -= int64(amount)
			_, err = a.Save()
		}
	}
	return err
}

func (a *accountImpl) flowAccountConsume(title string, amount int, outerNo string, remark string) error {
	if a.value.FlowBalance < int64(amount) {
		return member.ErrAccountNotEnoughAmount
	}
	l, err := a.createFlowAccountLog(member.KindConsume, title, -amount, outerNo, true)
	if err == nil {
		l.Remark = remark
		_, err = a.rep.SaveFlowAccountInfo(l)
		if err == nil {
			a.value.FlowBalance -= int64(amount)
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
func (a *accountImpl) discountBalance(title string, amount int, outerNo string, remark string) (err error) {
	if a.value.Balance < int64(amount) {
		return member.ErrAccountNotEnoughAmount
	}
	l, err := a.createBalanceLog(member.KindDiscount, title, -amount, outerNo, true)
	if err == nil {
		l.Remark = remark
		if err == nil {
			a.value.Balance -= int64(amount)
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
