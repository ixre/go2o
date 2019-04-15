// create for src 30/11/2017 ( jarrysix@gmail.com )
package merchant

import (
	"github.com/ixre/gof/db/orm"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/wallet"
	"go2o/core/domain/tmp"
	"go2o/core/variable"
	"math"
	"time"
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
func (a *accountImpl) GetDomainId() int32 {
	return a.value.MchId
}

// 获取账户值
func (a *accountImpl) GetValue() *merchant.Account {
	return a.value
}

// 保存
func (a *accountImpl) Save() error {
	_, err := orm.Save(tmp.Db().GetOrm(), a.value, int(a.GetDomainId()))
	//_, err := a.mchImpl._rep.SaveMerchantAccount(a)
	return err
}

// 根据编号获取余额变动信息
func (a *accountImpl) GetBalanceLog(id int32) *merchant.BalanceLog {
	e := merchant.BalanceLog{}
	if tmp.Db().GetOrm().Get(id, &e) == nil {
		return &e
	}
	return nil
	//return a.mchImpl._rep.GetBalanceLog(id)
}

// 根据号码获取余额变动信息
func (a *accountImpl) GetBalanceLogByOuterNo(outerNo string) *merchant.BalanceLog {
	e := merchant.BalanceLog{}
	if tmp.Db().GetOrm().GetBy(&e, "outer_no=?", outerNo) == nil {
		return &e
	}
	return nil
	//return a.mchImpl._rep.GetBalanceLogByOuterNo(outerNo)
}

func (a *accountImpl) createBalanceLog(kind int, title string, outerNo string,
	amount float32, csn float32, state int) *merchant.BalanceLog {
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
func (a *accountImpl) SaveBalanceLog(v *merchant.BalanceLog) (int32, error) {
	return orm.I32(orm.Save(tmp.Db().GetOrm(), v, int(v.Id)))
	//return a.mchImpl._rep.SaveBalanceLog(v)
}

// 支出
func (a *accountImpl) TakePayment(outerNo string, amount float64, csn float64, remark string) error {
	if amount <= 0 || math.IsNaN(amount) {
		return merchant.ErrAmount
	}
	if float64(a.value.Balance) < amount {
		return merchant.ErrNoMoreAmount
	}
	l := a.createBalanceLog(merchant.KindAccountTakePayment,
		remark, outerNo, float32(-amount), float32(csn), 1)
	_, err := a.SaveBalanceLog(l)
	if err == nil {
		a.value.Balance -= float32(amount)
		a.value.UpdateTime = time.Now().Unix()
		err = a.Save()
		if err == nil {
			iw := a.getWallet()
			err = iw.Discount(int(amount*float64(wallet.AmountRateSize)), remark, outerNo, true)
		}
	}
	return err
}

// 订单结账
func (a *accountImpl) SettleOrder(orderNo string, amount int, tradeFee int,
	refundAmount int, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return merchant.ErrAmount
	}
	fAmount := float32(amount / 100)
	fTradeFee := float32(tradeFee / 100)
	fRefund := float32(refundAmount / 100)
	a.value.Balance += fAmount
	a.value.SalesAmount += fTradeFee
	a.value.RefundAmount += fRefund
	a.value.UpdateTime = time.Now().Unix()
	err := a.Save()
	if err == nil {
		iw := a.getWallet()
		err = iw.Income(amount-tradeFee, tradeFee, remark, orderNo)
		// 记录旧日志,todo:可能去掉
		l := a.createBalanceLog(merchant.KindAccountSettleOrder,
			remark, orderNo, fAmount, fTradeFee, 1)
		a.SaveBalanceLog(l)
	}
	return err
}

func (a *accountImpl) getWallet() wallet.IWallet {
	iw := a.walletRepo.GetWalletByUserId(int64(a.GetValue().MchId), wallet.TMerchant)
	if iw == nil {
		iw = a.walletRepo.CreateWallet(&wallet.Wallet{
			UserId:     int64(a.GetValue().MchId),
			WalletType: wallet.TMerchant,
			WalletFlag: wallet.FlagCharge | wallet.FlagDiscount,
		})
		iw.Save()
	}
	return iw
}

//todo: 转入到奖金，手续费又被用于消费。这是一个bug

// 提现
//todo:???

// 转到会员账户
func (a *accountImpl) TransferToMember(amount float32) error {
	panic("TransferToMember需重构或移除")
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return merchant.ErrAmount
	}
	if a.value.Balance < amount || a.value.Balance <= 0 {
		return merchant.ErrNoMoreAmount
	}
	if a.mchImpl._value.MemberId <= 0 {
		return member.ErrNoSuchMember
	}
	m := a.memberRepo.GetMember(a.mchImpl._value.MemberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	l := a.createBalanceLog(merchant.KindAccountTransferToMember,
		"提取到会员"+variable.AliasWalletAccount, "", -amount, 0, 1)
	_, err := a.SaveBalanceLog(l)
	if err == nil {
		err = m.GetAccount().Charge(member.AccountWallet,
			member.KindWalletAdd,
			variable.AliasMerchantBalanceAccount+
				"提现", "", amount, member.DefaultRelateUser)
		if err != nil {
			return err
		}
		a.value.Balance -= amount
		a.value.TakeAmount += amount
		a.value.UpdateTime = time.Now().Unix()
		err = a.Save()
		if err != nil {
			return err
		}

		// 判断是否提现免费,如果免费,则赠送手续费
		registry := a.mchImpl._valRepo.GetRegistry()
		if registry.MerchantTakeOutCashFree {
			conf := a.mchImpl._valRepo.GetGlobNumberConf()
			if conf.TakeOutCsn > 0 {
				csn := amount * conf.TakeOutCsn
				err = m.GetAccount().Charge(member.AccountWallet,
					member.KindWalletAdd, "返还商户提现手续费", "",
					csn, member.DefaultRelateUser)
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
	if a.value.Balance < amount || a.value.Balance <= 0 {
		return merchant.ErrNoMoreAmount
	}
	if a.mchImpl._value.MemberId <= 0 {
		return member.ErrNoSuchMember
	}
	m := a.memberRepo.GetMember(a.mchImpl._value.MemberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	l := a.createBalanceLog(merchant.KindAccountTransferToMember,
		"提取到会员"+variable.AliasWalletAccount, "", -amount, 0, 1)
	_, err := a.SaveBalanceLog(l)
	if err == nil {
		err = m.GetAccount().Charge(member.AccountWallet,
			member.KindWalletAdd, variable.AliasMerchantBalanceAccount+
				"提现", "", amount, member.DefaultRelateUser)
		if err != nil {
			return err
		}
		a.value.Balance -= amount
		a.value.TakeAmount += amount
		a.value.UpdateTime = time.Now().Unix()
		err = a.Save()
		if err != nil {
			return err
		}

		// 判断是否提现免费,如果免费,则赠送手续费
		registry := a.mchImpl._valRepo.GetRegistry()
		if registry.MerchantTakeOutCashFree {
			conf := a.mchImpl._valRepo.GetGlobNumberConf()
			if conf.TakeOutCsn > 0 {
				csn := float32(0)
				err = m.GetAccount().Charge(member.AccountWallet,
					member.KindWalletAdd, "返还商户提现手续费", "",
					csn, member.DefaultRelateUser)
			}
		}
	}

	return err
}

// 赠送
func (a *accountImpl) Present(amount float32, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return merchant.ErrAmount
	}
	l := a.createBalanceLog(merchant.KindAccountPresent,
		remark, "", amount, 0, 1)
	_, err := a.SaveBalanceLog(l)
	if err == nil {
		a.value.PresentAmount += amount
		a.value.UpdateTime = time.Now().Unix()
		err = a.Save()
	}
	return err
}

// 充值
func (a *accountImpl) Charge(kind int32, amount float64,
	title, outerNo string, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return merchant.ErrAmount
	}
	l := a.createBalanceLog(merchant.KindAccountCharge,
		title, outerNo, float32(amount), 0, 1)
	_, err := a.SaveBalanceLog(l)
	if err == nil {
		a.value.Balance += float32(amount)
		a.value.UpdateTime = time.Now().Unix()
		err = a.Save()
		if err != nil {
			return err
		}
	}
	return err
}
