// create for src 30/11/2017 ( jarrysix@gmail.com )
package merchant

import (
	"math"
	"time"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/domain/tmp"
	"github.com/ixre/go2o/core/variable"
	"github.com/ixre/gof/db/orm"
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
func (a *accountImpl) GetBalanceLog(id int32) *merchant.BalanceLog {
	e := merchant.BalanceLog{}
	if tmp.Orm.Get(id, &e) == nil {
		return &e
	}
	return nil
	//return a.mchImpl._repo.GetBalanceLog(id)
}

// 根据号码获取余额变动信息
func (a *accountImpl) GetBalanceLogByOuterNo(outerNo string) *merchant.BalanceLog {
	e := merchant.BalanceLog{}
	if tmp.Orm.GetBy(&e, "outer_no= $1", outerNo) == nil {
		return &e
	}
	return nil
	//return a.mchImpl._repo.GetBalanceLogByOuterNo(outerNo)
}

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
func (a *accountImpl) SaveBalanceLog(v *merchant.BalanceLog) (int32, error) {
	return orm.I32(orm.Save(tmp.Orm, v, int(v.Id)))
	//return a.mchImpl._repo.SaveBalanceLog(v)
}

// 支出
func (a *accountImpl) TakePayment(outerNo string, amount int, csn int, remark string) error {
	if amount <= 0 {
		return merchant.ErrAmount
	}
	if a.value.Balance < int64(amount) {
		return merchant.ErrNoMoreAmount
	}
	l := a.createBalanceLog(merchant.KindAccountTakePayment,
		remark, outerNo, -int64(amount), int64(csn), 1)
	_, err := a.SaveBalanceLog(l)
	if err == nil {
		a.value.Balance -= int64(amount)
		a.value.UpdateTime = time.Now().Unix()
		err = a.Save()
		if err == nil {
			iw := a.getWallet()
			err = iw.Discount(int(amount), remark, outerNo, true)
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
	fAmount := int64(amount / 100)
	fTradeFee := int64(tradeFee / 100)
	fRefund := int64(refundAmount / 100)
	a.value.Balance += fAmount
	a.value.SalesAmount += fTradeFee
	a.value.RefundAmount += fRefund
	a.value.UpdateTime = time.Now().Unix()
	err := a.Save()
	if err == nil {
		iw := a.getWallet()
		_, err = iw.CarryTo(wallet.OperateData{
			Title:  "订单结算",
			Amount: amount - tradeFee, OuterNo: orderNo, Remark: remark}, false, tradeFee)
		// 记录旧日志,todo:可能去掉
		l := a.createBalanceLog(merchant.KindAccountSettleOrder,
			remark, orderNo, fAmount, fTradeFee, 1)
		a.SaveBalanceLog(l)
	}
	return err
}

func (a *accountImpl) getWallet() wallet.IWallet {
	iw := a.walletRepo.GetWalletByUserId(a.GetValue().MchId, wallet.TMerchant)
	if iw == nil {
		iw = a.walletRepo.CreateWallet(a.GetValue().MchId,
			a.mchImpl._value.LoginUser,
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
	m := a.memberRepo.GetMember(a.mchImpl._value.MemberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	l := a.createBalanceLog(merchant.KindAccountTransferToMember,
		"提取到会员"+variable.AliasWalletAccount, "", -int64(amount), 0, 1)
	_, err := a.SaveBalanceLog(l)
	if err == nil {
		_, err = m.GetAccount().CarryTo(member.AccountWallet,
			member.AccountOperateData{
				Title:   variable.AliasMerchantBalanceAccount + "提现",
				Amount:  amount * 100,
				OuterNo: "",
				Remark:  "sys",
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
						Title:   "返还商户提现手续费",
						Amount:  int(csn * 100),
						OuterNo: "",
						Remark:  "",
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
	m := a.memberRepo.GetMember(a.mchImpl._value.MemberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	l := a.createBalanceLog(merchant.KindAccountTransferToMember,
		"提取到会员"+variable.AliasWalletAccount, "", -int64(amount), 0, 1)
	_, err := a.SaveBalanceLog(l)
	if err == nil {
		_, err = m.GetAccount().CarryTo(member.AccountWallet,
			member.AccountOperateData{
				Title:   variable.AliasMerchantBalanceAccount + "提现",
				Amount:  int(amount * 100),
				OuterNo: "",
				Remark:  "sys",
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
						Title:   "返还商户提现手续费",
						Amount:  int(csn * 100),
						OuterNo: "",
						Remark:  "",
					}, false, 0)
			}
		}
	}

	return err
}

// 赠送
func (a *accountImpl) Present(amount int, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return merchant.ErrAmount
	}
	l := a.createBalanceLog(merchant.KindAccountPresent,
		remark, "", int64(amount), 0, 1)
	_, err := a.SaveBalanceLog(l)
	if err == nil {
		a.value.PresentAmount += int64(amount)
		a.value.UpdateTime = time.Now().Unix()
		err = a.Save()
	}
	return err
}

// 充值
func (a *accountImpl) Charge(kind int32, amount int,
	title, outerNo string, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return merchant.ErrAmount
	}
	l := a.createBalanceLog(merchant.KindAccountCharge,
		title, outerNo, int64(amount), 0, 1)
	_, err := a.SaveBalanceLog(l)
	if err == nil {
		a.value.Balance += int64(amount)
		a.value.UpdateTime = time.Now().Unix()
		err = a.Save()
		if err != nil {
			return err
		}
	}
	return err
}
