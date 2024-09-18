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
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/gof/domain/eventbus"
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
	a.wallet = a.walletRepo.GetWalletByUserId(int64(a.member.GetAggregateRootId()), 1)
	if a.wallet != nil {
		a.value.WalletCode = a.wallet.Get().HashCode
		a.Save()
		return
	}
	//　创建新的钱包
	a.wallet = a.walletRepo.CreateWallet(
		int(a.member.GetAggregateRootId()),
		a.member.value.Username,
		1, "MemberWallet", flag)
	if _, err := a.wallet.Save(); err != nil {
		log.Println("[ GO2O][ member]: create wallet failed,error", err.Error())
	}
	a.value.WalletCode = a.wallet.Get().HashCode // 绑定钱包
}

// GetDomainId 获取领域对象编号
func (a *accountImpl) GetDomainId() int {
	return a.value.MemberId
}

// GetValue 获取账户值
func (a *accountImpl) GetValue() *member.Account {
	return a.value
}

// Save 保存
func (a *accountImpl) Save() (int, error) {
	// 判断是否新建账号
	origin := a.rep.GetAccount(a.member.GetAggregateRootId())
	isCreate := origin == nil
	// 更新账户
	a.value.MemberId = int(a.member.GetAggregateRootId())
	a.value.UpdateTime = int(time.Now().Unix())
	n, err := a.rep.SaveAccount(a.value)
	if err == nil {
		// 创建钱包
		if isCreate {
			a.initWallet()
			a.rep.SaveAccount(a.value)
		}
		// 推送钱包更新消息
		if !isCreate {
			eventbus.Publish(&events.MemberAccountPushEvent{
				Account: *a.value,
			})
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
	if account != member.AccountWallet {
		return errors.New("只支持钱包充值")
	}
	return a.chargeWallet(title, amount, outerNo, remark)
}

func (a *accountImpl) CarryTo(account member.AccountType, d member.AccountOperateData, freeze bool, transactionFee int) (int, error) {
	if d.Amount <= 0 || math.IsNaN(float64(d.Amount)) {
		return 0, member.ErrIncorrectQuota
	}
	switch account {
	case member.AccountIntegral:
		return 0, a.carryToIntegral(d, freeze)
	case member.AccountBalance:
		return 0, a.carryToBalance(d, freeze, transactionFee)
	case member.AccountWallet:
		return a.carryToWallet(d, freeze, transactionFee)
	case member.AccountFlow:
		return 0, a.chargeFlow(d)
	}
	return 0, member.ErrNotSupportAccountType
}

// ReviewCarryTo 审核入账
func (a *accountImpl) ReviewCarryTo(account member.AccountType, transactionId int, pass bool, reason string) error {
	switch account {
	case member.AccountIntegral:
		return a.reviewIntegralCarryTo(transactionId, pass, reason)
	case member.AccountBalance:
		return a.reviewBalanceCarryTo(transactionId, pass, reason)
	case member.AccountWallet:
		return a.reviewWalletCarryTo(transactionId, pass, reason)
	}
	return member.ErrNotSupportAccountType
}

func (a *accountImpl) carryToWallet(d member.AccountOperateData, freeze bool, transactionFee int) (int, error) {
	id, err := a.wallet.CarryTo(wallet.TransactionData{
		TransactionTitle:  d.TransactionTitle,
		Amount:            d.Amount,
		OuterTxNo:         d.OuterTransactionNo,
		TransactionFee:    transactionFee,
		TransactionRemark: d.TransactionRemark,
		OuterTxUid:        d.OuterTxUid,
	}, freeze)
	if err == nil {
		err = a.asyncWallet()
	}
	return id, err
}

func (a *accountImpl) reviewWalletCarryTo(transactionId int, pass bool, reason string) error {
	err := a.wallet.ReviewCarryTo(transactionId, pass, reason)
	if err == nil {
		err = a.asyncWallet()
	}
	return err
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

// PrefreezeConsume implements member.IAccount.
func (a *accountImpl) PrefreezeConsume(transactionId int, transactionTitle string, transactionRemark string) error {
	return a.wallet.PrefreezeConsume(wallet.TransactionData{
		TransactionTitle:  transactionTitle,
		Amount:            0,
		TransactionFee:    0,
		OuterTxNo:         "",
		TransactionRemark: transactionRemark,
		TransactionId:     transactionId,
		OuterTxUid:        0,
	})
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
		return a.walletDiscount(title, amount, outerNo)
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
		Id:           0,
		MemberId:     int(a.value.MemberId),
		Kind:         kind,
		Subject:      title,
		OuterNo:      outerNo,
		ChangeValue:  value,
		Remark:       "",
		RelateUser:   0,
		ReviewStatus: int(enum.ReviewApproved),
		CreateTime:   int(unix),
		UpdateTime:   int(unix),
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
		MemberId:     int(a.value.MemberId),
		Kind:         int16(kind),
		Subject:      title,
		OuterNo:      outerNo,
		ChangeValue:  int(amount),
		ReviewStatus: int(enum.ReviewApproved),
		RelateUser:   0,
		CreateTime:   int(unix),
		UpdateTime:   int(unix),
	}, nil

}

// // 创建钱包日志
// func (a *accountImpl) createWalletLog(kind int, title string, amount int, outerNo string, checkTrade bool) (*member.WalletAccountLog, error) {
// 	title = strings.TrimSpace(title)
// 	outerNo = strings.TrimSpace(outerNo)
// 	if len(title) == 0 {
// 		return nil, member.ErrNoSuchLogTitleOrRemark
// 	}
// 	if math.IsNaN(float64(amount)) {
// 		return nil, member.ErrIncorrectAmount
// 	}
// 	if checkTrade && len(outerNo) == 0 {
// 		return nil, member.ErrMissingOuterNo
// 	}
// 	unix := time.Now().Unix()
// 	return &member.WalletAccountLog{
// 		MemberId:    a.value.MemberId,
// 		Kind:        kind,
// 		Title:       title,
// 		OuterNo:     outerNo,
// 		Amount:      int64(amount),
// 		ReviewStatus: enum.ReviewPass,
// 		RelateUser:  0,
// 		CreateTime:  unix,
// 		UpdateTime:  unix,
// 	}, nil
// }

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
		MemberId:     int64(a.value.MemberId),
		Kind:         kind,
		Title:        title,
		OuterNo:      outerNo,
		Amount:       int64(amount),
		ReviewStatus: int(enum.ReviewApproved),
		RelateUser:   0,
		CreateTime:   unix,
		UpdateTime:   unix,
	}, nil

}

// 充值余额
func (a *accountImpl) carryToBalance(d member.AccountOperateData, freeze bool, transactionFee int) error {
	if d.Amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if transactionFee > 0 {
		d.Amount -= transactionFee
	}
	l, err := a.createBalanceLog(member.KindCarry, d.TransactionTitle, d.Amount, d.OuterTransactionNo, true)
	if freeze {
		a.value.FreezeBalance += d.Amount
		l.ReviewStatus = int(enum.ReviewPending)
		l.Remark = "待审核"
	} else {
		a.value.Balance += d.Amount
		l.ReviewStatus = int(enum.ReviewConfirm)
	}
	if err == nil {
		l.Remark = d.TransactionRemark
		l.ProcedureFee = transactionFee
		l.Balance = int(a.value.Balance)
		_, err = a.rep.SaveBalanceLog(l)
		if err == nil {
			_, err = a.Save()
		}
	}
	return err
}

func (a *accountImpl) reviewBalanceCarryTo(transactionId int, pass bool, reason string) error {
	l := a.rep.GetBalanceLog(transactionId)
	if l.ReviewStatus != int(enum.ReviewPending) {
		return wallet.ErrNotSupport
	}
	a.value.FreezeBalance -= l.ChangeValue
	l.UpdateTime = int(time.Now().Unix())
	if pass {
		a.value.Balance += l.ChangeValue
		l.ReviewStatus = int(enum.ReviewApproved)
		l.Remark = "系统审核通过"
	} else {
		l.ReviewStatus = int(enum.ReviewRejected)
		l.Remark = reason
	}
	_, err := a.rep.SaveBalanceLog(l)
	if err == nil {
		_, err = a.Save()
	}
	return err
}

// 充值积分
func (a *accountImpl) carryToIntegral(d member.AccountOperateData, freeze bool) error {
	if d.Amount <= 0 {
		return member.ErrIncorrectAmount
	}
	l, err := a.createIntegralLog(member.KindCarry, d.TransactionTitle, d.Amount, d.OuterTransactionNo, true)
	if freeze {
		a.value.FreezeIntegral += d.Amount
		l.ReviewStatus = int(enum.ReviewPending)
		l.Remark = "待审核"
	} else {
		a.value.Integral += d.Amount
		l.ReviewStatus = int(enum.ReviewApproved)
	}
	if err == nil {
		l.Remark = d.TransactionRemark
		l.Balance = a.value.Integral
		//l.ProcedureFee = transactionFee
		err = a.rep.SaveIntegralLog(l)
		if err == nil {
			_, err = a.Save()
		}
	}
	return err
}

func (a *accountImpl) reviewIntegralCarryTo(transactionId int, pass bool, reason string) error {
	l := a.rep.GetIntegralLog(transactionId)
	if l.ReviewStatus != int(enum.ReviewPending) {
		return wallet.ErrNotSupport
	}
	a.value.FreezeIntegral -= l.ChangeValue
	l.UpdateTime = int(time.Now().Unix())
	if pass {
		a.value.Integral += l.ChangeValue
		l.ReviewStatus = int(enum.ReviewApproved)
		l.Remark = "系统审核通过"
	} else {
		l.ReviewStatus = int(enum.ReviewRejected)
		l.Remark = reason
	}
	err := a.rep.SaveIntegralLog(l)
	if err == nil {
		_, err = a.Save()
	}
	return err
}

// 流通账户余额充值，如扣除,amount传入负数金额
func (a *accountImpl) chargeFlow(d member.AccountOperateData) error {
	if d.Amount <= 0 {
		return member.ErrIncorrectAmount
	}
	l, err := a.createFlowAccountLog(member.KindCharge, d.TransactionTitle, d.Amount, d.OuterTransactionNo, true)
	if err == nil {
		l.Remark = d.TransactionRemark
		_, err = a.rep.SaveFlowAccountInfo(l)
		if err == nil {
			a.value.FlowBalance += d.Amount
			_, err = a.Save()
		}
	}
	return err
}

// 调整账户余额
func (a *accountImpl) adjustBalanceAccount(title string, amount int, remark string, relateUser int64) error {
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:     int(a.value.MemberId),
		Kind:         member.KindAdjust,
		Subject:      title,
		OuterNo:      "",
		ChangeValue:  int(amount),
		Remark:       remark,
		RelateUser:   int(relateUser),
		ReviewStatus: int(enum.ReviewApproved),
		CreateTime:   int(unix),
		UpdateTime:   int(unix),
	}
	_, err := a.rep.SaveBalanceLog(v)
	if err == nil {
		a.value.Balance += amount
		_, err = a.Save()
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
		MemberId:     int(a.value.MemberId),
		Kind:         int16(kind),
		Subject:      title,
		OuterNo:      outerNo,
		ChangeValue:  int(amount),
		ReviewStatus: int(enum.ReviewApproved),
		RelateUser:   int(relateUser),
		CreateTime:   int(unix),
		UpdateTime:   int(unix),
	}
	_, err := a.rep.SaveBalanceLog(v)
	if err == nil {
		a.value.Balance += amount
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
		a.value.TotalWalletAmount += amount
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
func (a *accountImpl) walletDiscount(title string, amount int, outerNo string) error {
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

// 钱包退款(指定业务类型)
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
	a.value.WalletBalance = a.wallet.Get().Balance
	a.value.FreezeWallet = a.wallet.Get().FreezeAmount
	_, err := a.Save()
	return err
}

// GetWalletLog 根据编号获取余额变动信息
func (a *accountImpl) GetWalletLog(id int64) wallet.WalletLog {
	return a.wallet.GetLog(id)
}

// Freeze 冻结余额
func (a *accountImpl) Freeze(account member.AccountType, p member.AccountOperateData, relateUser int64) (int, error) {
	switch account {
	case member.AccountBalance:
		return a.freezeBalance(p, relateUser)
	case member.AccountWallet:
		return a.freezeWallet(p, relateUser)
	case member.AccountIntegral:
		return a.freezesIntegral(p, relateUser)
	}
	panic("not support account type")
}

// Freeze 冻结余额
func (a *accountImpl) freezeBalance(p member.AccountOperateData, relateUser int64) (int, error) {
	if p.Amount <= 0 || math.IsNaN(float64(p.Amount)) {
		return 0, member.ErrIncorrectAmount
	}
	if a.value.Balance < p.Amount {
		return 0, member.ErrAccountNotEnoughAmount
	}
	if len(p.TransactionTitle) == 0 {
		p.TransactionTitle = "资金冻结"
	}
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:     int(a.value.MemberId),
		Kind:         int16(member.KindFreeze),
		Subject:      p.TransactionTitle,
		OuterNo:      p.OuterTransactionNo,
		ChangeValue:  -int(p.Amount),
		RelateUser:   int(relateUser),
		Remark:       p.TransactionRemark,
		ReviewStatus: int(enum.ReviewApproved),
		CreateTime:   int(unix),
		UpdateTime:   int(unix),
	}
	a.value.Balance -= p.Amount
	a.value.FreezeBalance += p.Amount
	_, err := a.Save()
	if err == nil {
		_, err = a.rep.SaveBalanceLog(v)
	}
	return int(v.Id), err
}

// FreezeWallet 冻结钱包
func (a *accountImpl) freezeWallet(p member.AccountOperateData, relateUser int64) (int, error) {
	id, err := a.wallet.Freeze(wallet.TransactionData{
		TransactionTitle:  p.TransactionTitle,
		Amount:            p.Amount,
		OuterTxNo:         p.OuterTransactionNo,
		TransactionRemark: p.TransactionRemark,
		TransactionId:     p.TransactionId,
		OuterTxUid:        p.OuterTxUid,
	}, wallet.Operator{
		OperatorUid:  int(relateUser),
		OperatorName: "",
	})
	if err == nil {
		err = a.asyncWallet()
	}
	return id, err

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
	//	ReviewStatus: enum.ReviewPass,
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
func (a *accountImpl) freezesIntegral(p member.AccountOperateData, relateUser int64) (int, error) {
	//if !new {
	//	if a.value.Integral < value {
	//		return member.ErrNoSuchIntegral
	//	}
	//	a.value.Integral -= value
	//}
	a.value.FreezeIntegral += p.Amount
	unix := time.Now().Unix()
	l := &member.IntegralLog{
		Id:           0,
		MemberId:     int(a.value.MemberId),
		Kind:         member.TypeIntegralFreeze,
		Subject:      p.TransactionTitle,
		OuterNo:      p.OuterTransactionNo,
		ChangeValue:  -p.Amount,
		Remark:       p.TransactionRemark,
		Balance:      a.value.Integral,
		RelateUser:   int(relateUser),
		ReviewStatus: int(enum.ReviewApproved),
		CreateTime:   int(unix),
		UpdateTime:   int(unix),
	}
	_, err := a.Save()
	if err == nil {
		err = a.rep.SaveIntegralLog(l)
	}
	return l.Id, err
}

// Unfreeze 解冻金额

func (a *accountImpl) Unfreeze(account member.AccountType, d member.AccountOperateData, isRefundBalance bool, relateUser int64) error {
	if d.Amount <= 0 || math.IsNaN(float64(d.Amount)) {
		return member.ErrIncorrectAmount
	}
	switch account {
	case member.AccountBalance:
		return a.unfreezeBalance(d, relateUser)
	case member.AccountWallet:
		return a.unfreezeWallet(d, isRefundBalance, relateUser)
	case member.AccountIntegral:
		return a.unfreezesIntegral(d, relateUser)
	}
	panic("not support account type")
}

func (a *accountImpl) unfreezeBalance(d member.AccountOperateData, relateUser int64) error {
	if a.value.FreezeBalance < d.Amount {
		return member.ErrAccountNotEnoughAmount
	}
	if len(d.TransactionTitle) == 0 {
		d.TransactionTitle = "资金解结"
	}
	a.value.FreezeBalance -= d.Amount
	a.value.Balance += d.Amount
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:     int(a.value.MemberId),
		Kind:         int16(member.KindUnfreeze),
		Subject:      d.TransactionTitle,
		RelateUser:   int(relateUser),
		ChangeValue:  int(d.Amount),
		OuterNo:      d.OuterTransactionNo,
		Balance:      int(a.value.Balance),
		ProcedureFee: 0,
		ReviewStatus: int(enum.ReviewApproved),
		CreateTime:   int(unix),
		UpdateTime:   int(unix),
	}
	_, err := a.Save()
	if err == nil {
		_, err = a.rep.SaveBalanceLog(v)
	}
	return err

}

// UnfreezeWallet 解冻赠送金额
func (a *accountImpl) unfreezeWallet(d member.AccountOperateData, isRefundBalance bool, relateUser int64) error {
	err := a.wallet.Unfreeze(d.Amount, d.TransactionTitle, d.OuterTransactionNo, isRefundBalance, int(relateUser), "")
	if err == nil {
		err = a.asyncWallet()
	}
	return err
}

// PaymentDiscount 支付单抵扣消费,tradeNo为支付单单号
func (a *accountImpl) PaymentDiscount(tradeNo string, amount int, remark string) error {
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
		MemberId:     int(a.value.MemberId),
		Kind:         int16(member.KindDiscount),
		Subject:      remark,
		OuterNo:      tradeNo,
		ChangeValue:  -int(amount),
		ReviewStatus: int(enum.ReviewApproved),
		RelateUser:   int(member.DefaultRelateUser),
		CreateTime:   int(unix),
		UpdateTime:   int(unix),
	}
	_, err := a.rep.SaveBalanceLog(v)
	if err == nil {
		a.value.Balance -= amount
		_, err = a.Save()
	}
	return err
}

// UnfreezesIntegral 解冻积分
func (a *accountImpl) unfreezesIntegral(d member.AccountOperateData, relateUser int64) error {
	if a.value.FreezeIntegral < d.Amount {
		return member.ErrNoSuchIntegral
	}
	a.value.FreezeIntegral -= d.Amount
	a.value.Integral += d.Amount
	unix := time.Now().Unix()
	var l = &member.IntegralLog{
		Id:           0,
		MemberId:     int(a.value.MemberId),
		Kind:         member.TypeIntegralUnfreeze,
		Subject:      d.TransactionTitle,
		OuterNo:      d.OuterTransactionNo,
		ChangeValue:  d.Amount,
		Remark:       d.TransactionRemark,
		Balance:      a.value.Integral,
		RelateUser:   int(relateUser),
		ReviewStatus: int(enum.ReviewApproved),
		CreateTime:   int(unix),
		UpdateTime:   int(unix),
	}
	_, err := a.Save()
	if err == nil {
		err = a.rep.SaveIntegralLog(l)
	}
	return err
}

// RequestWithdrawal 请求提现,返回info_id,交易号及错误
func (a *accountImpl) RequestWithdrawal(w *wallet.WithdrawTransaction) (int, string, error) {
	if w.Kind != wallet.KWithdrawExchange &&
		w.Kind != wallet.KWithdrawToBankCard &&
		w.Kind != wallet.KWithdrawToPayWallet {
		return 0, "", member.ErrNotSupportTakeOutBusinessKind
	}
	if w.Amount <= 0 || math.IsNaN(float64(w.Amount)) {
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
		trust := a.member.Profile().GetCertificationInfo()
		if trust.ReviewStatus != int(enum.ReviewApproved) {
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
	if a.wallet.Get().Balance < w.Amount {
		return 0, "", member.ErrOutOfBalance
	}
	// 检测提现金额是否超过限制
	minAmountStr, _ := a.registryRepo.GetValue(registry.MemberWithdrawMinAmount)
	minAmount, _ := strconv.Atoi(minAmountStr)
	if w.Amount < minAmount*100 {
		return 0, "", errors.New(fmt.Sprintf(member.ErrLessTakeAmount.Error(),
			format.FormatIntMoney(int64(minAmount))))
	}
	maxAmountStr, _ := a.registryRepo.GetValue(registry.MemberWithdrawMaxAmount)
	maxAmount, _ := strconv.Atoi(maxAmountStr)
	if maxAmount > 0 && w.Amount > maxAmount*100 {
		return 0, "", errors.New(fmt.Sprintf(member.ErrOutTakeAmount.Error(),
			format.FormatIntMoney(int64(maxAmount))))
	}
	// 检测是否超过限制
	maxTimes := a.registryRepo.Get(registry.MemberWithdrawMaxTimeOfDay).IntValue()
	if maxTimes > 0 {
		takeTimes := a.rep.GetTodayTakeOutTimes(int64(a.GetDomainId()))
		if takeTimes >= maxTimes {
			return 0, "", member.ErrAccountOutOfTakeOutTimes
		}
	}
	// 检测银行卡
	accountName := ""
	bankName := ""
	if w.Kind == wallet.KWithdrawToBankCard {
		if len(w.AccountNo) == 0 {
			return 0, "", errors.New("未指定提现的银行卡号")
		}
		bank := a.member.Profile().GetBankCard(w.AccountNo)
		if bank == nil {
			return 0, "", member.ErrBankNoSuchCard
		}
		accountName = bank.AccountName
		bankName = bank.BankName
	}
	finalAmount := w.Amount - w.TransactionFee
	if finalAmount > 0 {
		finalAmount = -finalAmount
	}
	id, tradeNo, err := a.wallet.RequestWithdrawal(
		wallet.WithdrawTransaction{
			Amount:           finalAmount,
			TransactionFee:   w.TransactionFee,
			Kind:             w.Kind,
			TransactionTitle: w.TransactionTitle,
			BankName:         bankName,
			AccountNo:        w.AccountNo,
			AccountName:      accountName,
		})
	if err == nil {
		err = a.asyncWallet()
		if err == nil {
			go a.rep.AddTodayTakeOutTimes(int64(a.GetDomainId()))
		}

		// 推送提现申请事件
		eventbus.Publish(&events.WithdrawalPushEvent{
			MemberId:       a.value.MemberId,
			RequestId:      int(id),
			Amount:         w.Amount,
			TransactionFee: w.TransactionFee,
			IsReviewEvent:  false,
		})
	}
	return id, tradeNo, err
}

// ReviewWithdrawal 确认提现
func (a *accountImpl) ReviewWithdrawal(transactionId int, pass bool, remark string) error {
	//todo: opr_uid
	err := a.wallet.ReviewWithdrawal(transactionId, pass, remark, 1, "系统")
	if err == nil {
		err = a.asyncWallet()
		if pass {
			log := a.wallet.GetLog(int64(transactionId))
			// 推送提现申请事件
			eventbus.Publish(&events.WithdrawalPushEvent{
				MemberId:       a.value.MemberId,
				RequestId:      int(transactionId),
				Amount:         int(log.ChangeValue),
				TransactionFee: log.TransactionFee,
				ReviewResult:   log.ReviewStatus == int(enum.ReviewApproved),
				IsReviewEvent:  true,
			})
		}
	}
	return err
}

// FinishWithdrawal 完成提现
func (a *accountImpl) FinishWithdrawal(transactionId int, outerTransactionNo string) error {
	//todo: opr_uid
	err := a.wallet.FinishWithdrawal(transactionId, outerTransactionNo)
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
	if a.value.FreezeBalance < amount {
		return member.ErrIncorrectAmount
	}
	unix := time.Now().Unix()
	a.value.FreezeBalance -= amount
	a.value.ExpiredBalance += amount
	a.value.UpdateTime = int(unix)
	l := &member.BalanceLog{
		MemberId:     int(a.value.MemberId),
		Kind:         int16(member.KindExpired),
		Subject:      "过期失效",
		OuterNo:      "",
		ChangeValue:  int(amount),
		ProcedureFee: 0,
		ReviewStatus: int(enum.ReviewApproved),
		RelateUser:   int(member.DefaultRelateUser),
		Remark:       remark,
		CreateTime:   int(unix),
		UpdateTime:   int(unix),
	}
	_, err := a.rep.SaveBalanceLog(l)
	if err == nil {
		_, err = a.Save()
	}
	return err
}

func (a *accountImpl) walletFreezeExpired(amount int, remark string) error {
	if a.value.FreezeWallet < amount {
		return member.ErrIncorrectAmount
	}
	unix := time.Now().Unix()
	a.value.FreezeWallet -= amount
	a.value.ExpiredWallet += amount
	a.value.UpdateTime = int(unix)
	l := &member.WalletAccountLog{
		MemberId:     int64(a.GetDomainId()),
		Kind:         member.KindExpired,
		Title:        "过期失效",
		OuterNo:      "",
		Amount:       int64(amount),
		ProcedureFee: 0,
		ReviewStatus: enum.ReviewApproved,
		RelateUser:   member.DefaultRelateUser,
		Remark:       remark,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.rep.SaveWalletAccountLog(l)
	if err == nil {
		_, err = a.Save()
	}
	return err
}

// 获取会员名称
func (a *accountImpl) getMemberName(m member.IMemberAggregateRoot) string {
	if tr := m.Profile().GetCertificationInfo(); tr.RealName != "" &&
		tr.ReviewStatus == int(enum.ReviewApproved) {
		return tr.RealName
	} else {
		return m.GetValue().Username
	}
}

// TransferAccount 转账
func (a *accountImpl) TransferAccount(account member.AccountType, toMember int64, amount int,
	transactionFee int, remark string) error {
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
		return a.transferWalletAccount(tm, tradeNo, amount, transactionFee, remark)
	case member.AccountBalance:
		return a.transferBalance(tm, tradeNo, amount, transactionFee, remark)
	}
	return nil
}

func (a *accountImpl) transferBalance(tm member.IMemberAggregateRoot, tradeNo string,
	tradeAmount, transactionFee int, remark string) error {
	csnFee := transactionFee
	amount := tradeAmount
	if a.value.Balance < int(amount+csnFee) {
		return member.ErrAccountNotEnoughAmount
	}
	unix := time.Now().Unix()
	// 扣款
	toName := a.getMemberName(tm)
	l := &member.BalanceLog{
		MemberId:     int(a.value.MemberId),
		Kind:         int16(member.KindTransferOut),
		Subject:      "转账给" + toName,
		OuterNo:      tradeNo,
		ChangeValue:  -int(amount),
		ProcedureFee: int(csnFee),
		ReviewStatus: int(enum.ReviewApproved),
		RelateUser:   int(member.DefaultRelateUser),
		Remark:       remark,
		CreateTime:   int(unix),
		UpdateTime:   int(unix),
	}
	_, err := a.rep.SaveBalanceLog(l)
	if err == nil {
		a.value.Balance -= amount + csnFee
		a.value.UpdateTime = int(unix)
		_, err = a.Save()
		if err == nil {
			err = tm.GetAccount().ReceiveTransfer(member.AccountBalance,
				int64(a.GetDomainId()), tradeNo, amount, remark)
		}
	}
	return err
}

func (a *accountImpl) transferWalletAccount(tm member.IMemberAggregateRoot, tradeNo string,
	tradeAmount, transactionFee int, remark string) error {
	csnFee := transactionFee
	amount := tradeAmount
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
		MemberId:     int64(a.GetDomainId()),
		Kind:         member.KindTransferOut,
		Title:        "转账给" + toName,
		OuterNo:      tradeNo,
		Amount:       -int64(amount),
		ProcedureFee: int64(csnFee),
		ReviewStatus: enum.ReviewApproved,
		RelateUser:   member.DefaultRelateUser,
		Remark:       remark,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.rep.SaveWalletAccountLog(l)
	if err == nil {
		a.value.WalletBalance -= amount + csnFee
		a.value.UpdateTime = int(unix)
		_, err = a.Save()
		if err == nil {
			err = tm.GetAccount().ReceiveTransfer(member.AccountWallet,
				int64(a.GetDomainId()), tradeNo, amount, remark)
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
		MemberId:     int64(a.GetDomainId()),
		Kind:         member.KindTransferIn,
		Title:        "转账收款（" + fromName + "）",
		OuterNo:      tradeNo,
		Amount:       int64(amount),
		ProcedureFee: 0,
		ReviewStatus: enum.ReviewApproved,
		RelateUser:   member.DefaultRelateUser,
		Remark:       remark,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.rep.SaveWalletAccountLog(tl)
	if err == nil {
		a.value.WalletBalance += amount
		a.value.UpdateTime = int(unix)
		_, err = a.Save()
	}
	return err
}

func (a *accountImpl) receiveBalanceTransfer(fromMember int64, tradeNo string,
	amount int, remark string) error {
	fromName := a.getMemberName(a.rep.GetMember(int64(a.GetDomainId())))
	unix := time.Now().Unix()
	tl := &member.BalanceLog{
		MemberId:     int(a.value.MemberId),
		Kind:         int16(member.KindTransferIn),
		Subject:      "转账收款（" + fromName + "）",
		OuterNo:      tradeNo,
		ChangeValue:  int(amount),
		ProcedureFee: 0,
		ReviewStatus: int(enum.ReviewApproved),
		RelateUser:   int(member.DefaultRelateUser),
		Remark:       remark,
		CreateTime:   int(unix),
		UpdateTime:   int(unix),
	}
	_, err := a.rep.SaveBalanceLog(tl)
	if err == nil {
		a.value.Balance += amount
		a.value.UpdateTime = int(unix)
		_, err = a.Save()
	}
	return err
}

// TransferBalance 转账余额到其他账户
func (a *accountImpl) TransferBalance(account member.AccountType, amount int,
	tradeNo string, toTitle, fromTitle string) error {
	var err error
	if account == member.AccountFlow {
		if a.value.Balance < amount {
			return member.ErrAccountNotEnoughAmount
		}
		a.value.Balance -= amount
		a.value.FlowBalance += amount
		if _, err = a.Save(); err == nil {
			a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
				Kind:         member.KindTransferOut,
				Title:        toTitle,
				Amount:       -int64(amount),
				OuterNo:      tradeNo,
				ReviewStatus: int(enum.ReviewApproved),
			})

			a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
				Kind:         member.KindTransferIn,
				Title:        fromTitle,
				Amount:       int64(amount),
				OuterNo:      tradeNo,
				ReviewStatus: int(enum.ReviewApproved),
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

	csnAmount := int(commission * float32(amount))
	finalAmount := amount - csnAmount

	if a.value.FlowBalance < int(finalAmount) {
		return member.ErrAccountNotEnoughAmount
	}

	a.value.FlowBalance -= finalAmount
	a.value.WalletBalance += finalAmount
	a.value.TotalWalletAmount += finalAmount

	if _, err = a.Save(); err == nil {
		a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
			Kind:         kind,
			Title:        toTitle,
			Amount:       -int64(amount),
			OuterNo:      tradeNo,
			CsnFee:       int64(csnAmount),
			ReviewStatus: int(enum.ReviewApproved),
		})

		a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
			Kind:         member.KindTransferIn,
			Title:        fromTitle,
			Amount:       int64(finalAmount),
			OuterNo:      tradeNo,
			ReviewStatus: int(enum.ReviewApproved),
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

	if a.value.FlowBalance < int(finalAmount) {
		return member.ErrAccountNotEnoughAmount
	}

	a.value.FlowBalance -= int(finalAmount)
	acc2.GetValue().FlowBalance += int(amount)

	if _, err = a.Save(); err == nil {
		_, err = a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
			Kind:         member.KindTransferOut,
			Title:        toTitle,
			Amount:       -finalAmount,
			CsnFee:       csnAmount,
			RelateUser:   memberId,
			OuterNo:      tradeNo,
			ReviewStatus: member.StatusOK,
		})

		if _, err = acc2.Save(); err == nil {
			_, err = a.rep.SaveFlowAccountInfo(&member.FlowAccountLog{
				Kind:         member.KindTransferIn,
				Title:        fromTitle,
				Amount:       int64(amount),
				RelateUser:   int64(a.value.MemberId),
				OuterNo:      tradeNo,
				ReviewStatus: member.StatusOK,
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

func (a *accountImpl) flowAccountConsume(title string, amount int, outerNo string, remark string) error {
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
func (a *accountImpl) discountBalance(title string, amount int, outerNo string, remark string) (err error) {
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
