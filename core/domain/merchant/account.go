// create for src 30/11/2017 ( jarrysix@gmail.com )
package merchant

import (
	"math"
	"time"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/variable"
)

var _ merchant.IAccount = new(accountImpl)

type accountImpl struct {
	mchImpl    *merchantImpl
	value      *merchant.Account
	memberRepo member.IMemberRepo
	walletRepo wallet.IWalletRepo
}

func newAccountImpl(mchImpl *merchantImpl, a *merchant.Account,
	memberRepo member.IMemberRepo, walletRepo wallet.IWalletRepo) merchant.IAccount {
	return &accountImpl{
		mchImpl:    mchImpl,
		value:      a,
		memberRepo: memberRepo,
		walletRepo: walletRepo,
	}
}

// 获取领域对象编号
func (a *accountImpl) GetDomainId() int64 {
	return a.value.MchId
}

// 获取账户值
func (a *accountImpl) GetValue() *merchant.Account {
	return a.value
}

// 保存
func (a *accountImpl) Save() error {
	_, err := a.mchImpl._repo.SaveAccount(a.value)
	return err
}

// 根据编号获取余额变动信息
func (a *accountImpl) GetBalanceLog(id int) *merchant.BalanceLog {
	return a.mchImpl._repo.GetBalanceAccountLog(id)
}

// // 根据号码获取余额变动信息
// func (a *accountImpl) GetBalanceLogByOuterNo(outerNo string) *merchant.BalanceLog {
// 	e := merchant.BalanceLog{}
// 	if tmp.Orm.GetBy(&e, "outer_no= $1", outerNo) == nil {
// 		return &e
// 	}
// 	return a.mchImpl._repo.GetBalanceLogByOuterNo(outerNo)
// }

func (a *accountImpl) createBalanceLog(kind int, title string, outerNo string,
	amount int64, csn int64, state int) *merchant.BalanceLog {
	unix := time.Now().Unix()
	return &merchant.BalanceLog{
		// 编号
		Id: 0,
		// 商户编号
		MchId: a.GetDomainId(),
		// 日志类型
		Kind: kind,
		// 标题
		Title: title,
		// 外部订单号
		OuterNo: outerNo,
		// 金额
		Amount: amount,
		// 手续费
		CsnAmount: csn,
		// 状态
		State: state,
		// 创建时间
		CreateTime: unix,
		// 更新时间
		UpdateTime: unix,
	}
}

// 保存余额变动信息
func (a *accountImpl) SaveBalanceLog(v *merchant.BalanceLog) (int, error) {
	return a.mchImpl._repo.SaveBalanceAccountLog(v)
}

// Consume 消耗商户支出，例如广告费、提现等
func (a *accountImpl) Consume(transactionTitle string, amount int, outerTxNo string, transactionRemark string) error {
	if amount <= 0 {
		return merchant.ErrAmount
	}
	if a.value.Balance < int64(amount) {
		return merchant.ErrNoMoreAmount
	}
	l := a.createBalanceLog(merchant.KindAccountTakePayment,
		transactionRemark, outerTxNo, -int64(amount), 0, 1)
	_, err := a.SaveBalanceLog(l)
	if err == nil {
		a.value.Balance -= int64(amount)
		a.value.UpdateTime = time.Now().Unix()
		err = a.Save()
		if err == nil {
			iw := a.getWallet()
			err = iw.Consume(int(amount),
				transactionTitle,
				outerTxNo,
				transactionRemark)
		}
	}
	return err
}

// 订单结账
func (a *accountImpl) Carry(p merchant.CarryParams) (txId int, err error) {
	if p.Amount <= 0 || math.IsNaN(float64(p.Amount)) {
		return 0, merchant.ErrAmount
	}
	fAmount := int64(p.Amount / 100)
	fTradeFee := int64(p.TransactionFee / 100)
	fRefund := int64(p.RefundAmount / 100)
	a.value.Balance += fAmount
	a.value.SalesAmount += fTradeFee
	a.value.RefundAmount += fRefund
	a.value.UpdateTime = time.Now().Unix()
	err = a.Save()
	if err == nil {
		iw := a.getWallet()
		return iw.CarryTo(wallet.TransactionData{
			TransactionTitle:  p.TransactionTitle,
			Amount:            p.Amount,
			TransactionFee:    p.TransactionFee,
			OuterTxNo:         p.OuterTxNo,
			TransactionRemark: p.TransactionRemark,
		}, p.Freeze)
		// 记录旧日志,todo:可能去掉
		// l := a.createBalanceLog(merchant.KindAccountCarry,
		// 	remark, orderNo, fAmount, fTradeFee, 1)
		// a.SaveBalanceLog(l)
	}
	return 0, err
}

func (a *accountImpl) getWallet() wallet.IWallet {
	iw := a.walletRepo.GetWalletByUserId(a.GetValue().MchId, wallet.TMerchant)
	if iw == nil {
		iw = a.walletRepo.CreateWallet(a.GetValue().MchId,
			a.mchImpl._value.Username,
			wallet.TMerchant,
			"MchWallet",
			wallet.FlagCharge|wallet.FlagDiscount)
		iw.Save()
	}
	return iw
}

//todo: 转入到奖金，手续费又被用于消费。这是一个bug

// 提现
//todo:???

// 转到会员账户
func (a *accountImpl) TransferToMember(amount int) error {
	panic("TransferToMember需重构或移除")
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return merchant.ErrAmount
	}
	if a.value.Balance < int64(amount) || a.value.Balance <= 0 {
		return merchant.ErrNoMoreAmount
	}
	if a.mchImpl._value.MemberId <= 0 {
		return member.ErrNoSuchMember
	}
	m := a.memberRepo.GetMember(int64(a.mchImpl._value.MemberId))
	if m == nil {
		return member.ErrNoSuchMember
	}
	l := a.createBalanceLog(merchant.KindAccountTransferToMember,
		"提取到会员"+variable.AliasWalletAccount, "", -int64(amount), 0, 1)
	_, err := a.SaveBalanceLog(l)
	if err == nil {
		_, err = m.GetAccount().CarryTo(member.AccountWallet,
			member.AccountOperateData{
				TransactionTitle:   variable.AliasMerchantBalanceAccount + "提现",
				Amount:             amount * 100,
				OuterTransactionNo: "",
				TransactionRemark:  "sys",
			}, false, 0)
		if err != nil {
			return err
		}
		a.value.Balance -= int64(amount)
		a.value.WithdrawAmount += int64(amount)
		a.value.UpdateTime = time.Now().Unix()
		err = a.Save()
		if err != nil {
			return err
		}
		// 判断是否提现免费,如果免费,则赠送手续费
		takeFee := a.mchImpl._registryRepo.Get(registry.MerchantTakeOutCashFree).BoolValue()
		if takeFee {
			takeRate := a.mchImpl._registryRepo.Get(registry.MerchantTakeOutCsn).FloatValue()
			if takeRate > 0 {
				csn := float64(amount) * takeRate
				_, err = m.GetAccount().CarryTo(member.AccountWallet,
					member.AccountOperateData{
						TransactionTitle:   "返还商户提现手续费",
						Amount:             int(csn * 100),
						OuterTransactionNo: "",
						TransactionRemark:  "",
					}, false, 0)
			}
		}
	}
	return err
}

func (a *accountImpl) TransferToMember1(amount float32) error {
	panic("TransferToMember2需重构或移除")
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return merchant.ErrAmount
	}
	if a.value.Balance < int64(amount) || a.value.Balance <= 0 {
		return merchant.ErrNoMoreAmount
	}
	if a.mchImpl._value.MemberId <= 0 {
		return member.ErrNoSuchMember
	}
	m := a.memberRepo.GetMember(int64(a.mchImpl._value.MemberId))
	if m == nil {
		return member.ErrNoSuchMember
	}
	l := a.createBalanceLog(merchant.KindAccountTransferToMember,
		"提取到会员"+variable.AliasWalletAccount, "", -int64(amount), 0, 1)
	_, err := a.SaveBalanceLog(l)
	if err == nil {
		_, err = m.GetAccount().CarryTo(member.AccountWallet,
			member.AccountOperateData{
				TransactionTitle:   variable.AliasMerchantBalanceAccount + "提现",
				Amount:             int(amount * 100),
				OuterTransactionNo: "",
				TransactionRemark:  "sys",
			}, false, 0)
		if err != nil {
			return err
		}
		a.value.Balance -= int64(amount)
		a.value.WithdrawAmount += int64(amount)
		a.value.UpdateTime = time.Now().Unix()
		err = a.Save()
		if err != nil {
			return err
		}

		// 判断是否提现免费,如果免费,则赠送手续费
		takeFree := a.mchImpl._registryRepo.Get(registry.MerchantTakeOutCashFree).BoolValue()
		if takeFree {
			rate := a.mchImpl._registryRepo.Get(registry.MerchantTakeOutCsn).FloatValue()
			if rate > 0 {
				csn := float32(rate) * amount
				_, err = m.GetAccount().CarryTo(member.AccountWallet,
					member.AccountOperateData{
						TransactionTitle:   "返还商户提现手续费",
						Amount:             int(csn * 100),
						OuterTransactionNo: "",
						TransactionRemark:  "",
					}, false, 0)
			}
		}
	}

	return err
}

// FreezeWallet 冻结钱包
func (a *accountImpl) Freeze(p wallet.TransactionData, relateUser int64) (int, error) {
	id, err := a.getWallet().Freeze(wallet.TransactionData{
		TransactionTitle:  p.TransactionTitle,
		Amount:            p.Amount,
		OuterTxNo:         p.OuterTxNo,
		TransactionRemark: p.TransactionRemark,
		TransactionId:     p.TransactionId,
	}, wallet.Operator{
		OperatorUid:  int(relateUser),
		OperatorName: "",
	})
	return id, err
}

// UnfreezeWallet 解冻赠送金额
func (a *accountImpl) Unfreeze(d wallet.TransactionData, relateUser int64) error {
	return a.getWallet().Unfreeze(d.Amount, d.TransactionTitle, d.OuterTxNo, int(relateUser), "")
}

// 调整钱包余额
func (a *accountImpl) Adjust(title string, amount int, remark string, relateUser int64) error {
	return a.getWallet().Adjust(amount, title, "", remark, int(relateUser), "-")
}

// FinishWithdrawal implements merchant.IAccount.
func (a *accountImpl) FinishWithdrawal(transactionId int, outerTransactionNo string) error {
	//todo: opr_uid
	return a.getWallet().FinishWithdrawal(transactionId, outerTransactionNo)
}

// RequestWithdrawal implements merchant.IAccount.
func (a *accountImpl) RequestWithdrawal(w *wallet.WithdrawTransaction) (int64, string, error) {
	return a.getWallet().RequestWithdrawal(
		wallet.WithdrawTransaction{
			Amount:           w.Amount,
			TransactionFee:   w.TransactionFee,
			Kind:             w.Kind,
			TransactionTitle: w.TransactionTitle,
			BankName:         w.BankName,
			AccountNo:        w.AccountNo,
			AccountName:      w.AccountName,
		})
}

// ReviewWithdrawal implements merchant.IAccount.
func (a *accountImpl) ReviewWithdrawal(transactionId int, pass bool, remark string) error {
	return a.getWallet().ReviewWithdrawal(transactionId, pass, remark, 1, "系统")
}
