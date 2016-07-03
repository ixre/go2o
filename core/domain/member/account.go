/**
 * Copyright 2015 @ z3q.net.
 * name : account
 * author : jarryliu
 * date : 2015-07-24 08:50
 * description :
 * history :
 */
package member

import (
	"errors"
	dm "go2o/core/domain"
	"go2o/core/domain/interface/member"
	"go2o/core/infrastructure/domain"
	"time"
)

var _ member.IAccount = new(accountImpl)

type accountImpl struct {
	_value *member.Account
	_rep   member.IMemberRep
}

func NewAccount(value *member.Account,
	rep member.IMemberRep) member.IAccount {
	return &accountImpl{
		_value: value,
		_rep:   rep,
	}
}

// 获取领域对象编号
func (this *accountImpl) GetDomainId() int {
	return this._value.MemberId
}

// 获取账户值
func (this *accountImpl) GetValue() *member.Account {
	return this._value
}

// 保存
func (this *accountImpl) Save() (int, error) {
	this._value.UpdateTime = time.Now().Unix()
	return this._rep.SaveAccount(this._value)
}

// 保存积分记录
func (this *accountImpl) SaveIntegralLog(l *member.IntegralLog) error {
	l.MemberId = this._value.MemberId
	return this._rep.SaveIntegralLog(l)
}

//　增加积分
// todo:merchantId 不需要
func (this *accountImpl) AddIntegral(merchantId int, backType int,
	integral int, log string) error {
	inLog := &member.IntegralLog{
		MerchantId: merchantId,
		MemberId:   this._value.MemberId,
		Type:       backType,
		Integral:   integral,
		Log:        log,
		RecordTime: time.Now().Unix(),
	}

	err := this._rep.SaveIntegralLog(inLog)
	if err == nil {
		this._value.Integral += integral
		_, err = this.Save()
	}
	return err
}

// 根据编号获取余额变动信息
func (this *accountImpl) GetBalanceInfo(id int) *member.BalanceInfo {
	return this._rep.GetBalanceInfo(id)
}

// 根据号码获取余额变动信息
func (this *accountImpl) GetBalanceInfoByNo(no string) *member.BalanceInfo {
	return this._rep.GetBalanceInfoByNo(no)
}

// 保存余额变动信息
func (this *accountImpl) SaveBalanceInfo(v *member.BalanceInfo) (int, error) {
	v.MemberId = this.GetDomainId()
	v.UpdateTime = time.Now().Unix()
	if v.CreateTime == 0 {
		v.CreateTime = v.UpdateTime
	}
	return this._rep.SaveBalanceInfo(v)
}

// 充值
// @title 充值标题说明
// @no    充值订单编号
// @amount 金额
func (this *accountImpl) ChargeBalance(chargeType int, title string, tradeNo string, amount float32) error {
	//todo: 客服充值需记录操作人
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}

	if chargeType == member.TypeBalanceNetPayCharge || chargeType == member.TypeBalanceSystemCharge ||
		chargeType == member.TypeBalanceServiceCharge || chargeType == member.TypeBalanceOrderRefund {

		v := &member.BalanceInfo{
			Kind:    member.KindBalanceCharge,
			Type:    chargeType,
			Title:   title,
			TradeNo: tradeNo,
			Amount:  amount,
			State:   1,
		}
		_, err := this.SaveBalanceInfo(v)
		if err == nil {
			this._value.Balance += amount
			_, err = this.Save()
		}
		return err
	}
	return errors.New("error charge type")
}

// 扣减余额
func (this *accountImpl) DiscountBalance(title string, tradeNo string, amount float32) (err error) {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if this._value.Balance < amount {
		return member.ErrNotEnoughAmount
	}
	v := &member.BalanceInfo{
		Kind:    member.KindBalanceDiscount,
		Title:   title,
		TradeNo: tradeNo,
		Amount:  -amount,
		State:   1,
	}
	_, err = this.SaveBalanceInfo(v)
	if err == nil {
		this._value.Balance -= amount
		_, err = this.Save()
	}
	return err
}

// 赠送金额
func (this *accountImpl) PresentBalance(title string, tradeNo string, amount float32) error {
	//todo:??客服调整
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if len(title) == 0 {
		if amount < 0 {
			title = "赠送账户出账"
		} else {
			title = "赠送账户入账"
		}
	}

	v := &member.BalanceInfo{
		Kind:    member.KindBalancePresent,
		Title:   title,
		TradeNo: tradeNo,
		Amount:  amount,
		State:   1,
	}
	_, err := this.SaveBalanceInfo(v)
	if err == nil {
		this._value.PresentBalance += amount
		if amount > 0 {
			this._value.TotalPresentFee += amount
		}
		_, err = this.Save()
	}
	return err
}

// 扣减奖金
func (this *accountImpl) DiscountPresent(title string, tradeNo string, amount float32, mustLargeZero bool) error {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if mustLargeZero && this._value.PresentBalance < amount {
		return member.ErrNotEnoughAmount
	}

	if len(title) == 0 {
		title = "出账"
	}

	v := &member.BalanceInfo{
		Kind:    member.KindPresentDiscount,
		Title:   title,
		TradeNo: tradeNo,
		Amount:  -amount,
		State:   1,
	}
	_, err := this.SaveBalanceInfo(v)
	if err == nil {
		this._value.PresentBalance -= amount
		_, err = this.Save()
	}
	return err
}

// 流通账户余额充值，如扣除,amount传入负数金额
func (this *accountImpl) ChargeFlowBalance(title string, tradeNo string, amount float32) error {
	if len(title) == 0 {
		if amount > 0 {
			title = "流动账户入账"
		} else {
			title = "流动账户出账"
		}
	}
	v := &member.BalanceInfo{
		Kind:    member.KindBalanceFlow,
		Title:   title,
		TradeNo: tradeNo,
		Amount:  amount,
		State:   1,
	}
	_, err := this.SaveBalanceInfo(v)
	if err == nil {
		this._value.FlowBalance += amount
		_, err = this.Save()
	}
	return err
}

// 支付单抵扣消费,tradeNo为支付单单号
func (this *accountImpl) PaymentDiscount(tradeNo string, amount float32) error {
	if amount < 0 || len(tradeNo) == 0 {
		return errors.New("amount error or missing trade no")
	}

	if amount > this._value.Balance {
		return member.ErrOutOfBalance
	}

	v := &member.BalanceInfo{
		Kind:    member.KindBalanceShopping,
		Type:    1,
		Title:   "订单抵扣",
		TradeNo: tradeNo,
		Amount:  -amount,
		State:   1,
	}
	_, err := this.SaveBalanceInfo(v)
	if err == nil {
		this._value.Balance -= amount
		_, err = this.Save()
	}
	return err
}

// 支付单抵扣积分,integral为积分,exchangeFee
func (this *accountImpl) DiscountIntegral(tradeNo string, integral int, exchangeFee float32) error {
	panic("未实现")
}

// 退款
func (this *accountImpl) RequestBackBalance(backType int, title string,
	amount float32) error {
	if amount > this._value.Balance {
		return member.ErrOutOfBalance
	}
	v := &member.BalanceInfo{
		Kind:   member.KindBalanceBack,
		Type:   backType,
		Title:  title,
		Amount: amount,
		State:  0,
	}
	_, err := this.SaveBalanceInfo(v)
	if err == nil {
		this._value.Balance -= amount
		_, err = this.Save()
	}
	return err
}

// 完成退款
func (this *accountImpl) FinishBackBalance(id int, tradeNo string) error {
	v := this.GetBalanceInfo(id)
	if v.Kind == member.KindBalanceBack {
		v.TradeNo = tradeNo
		v.State = 1
		_, err := this.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}

// 请求提现,返回info_id,交易号及错误
func (this *accountImpl) RequestApplyCash(applyType int, title string,
	amount float32, commission float32) (int, string, error) {
	if amount <= 0 {
		return 0, "", member.ErrIncorrectAmount
	}
	if this._value.PresentBalance < amount {
		return 0, "", member.ErrOutOfBalance
	}

	tradeNo := domain.NewTradeNo(00000)

	csnAmount := amount * commission
	finalAmount := amount - csnAmount
	if finalAmount > 0 {
		finalAmount = -finalAmount
	}
	v := &member.BalanceInfo{
		Kind:      member.KindBalanceApplyCash,
		Type:      applyType,
		Title:     title,
		TradeNo:   tradeNo,
		Amount:    finalAmount,
		CsnAmount: csnAmount,
		State:     member.StateApplySubmitted,
	}

	// 提现至余额
	if applyType == member.TypeApplyCashToCharge {
		this._value.Balance += amount
		v.State = member.StateApplyOver
	}

	id, err := this.SaveBalanceInfo(v)
	if err == nil {
		this._value.PresentBalance -= amount
		_, err = this.Save()
	}
	return id, tradeNo, err
}

// 确认提现
func (this *accountImpl) ConfirmApplyCash(id int, pass bool, remark string) error {
	//todo: remark
	v := this.GetBalanceInfo(id)
	if v.Kind == member.KindBalanceApplyCash {
		if pass {
			v.State = member.StateApplyConfirmed
		} else {
			if v.State == member.StateApplyNotPass {
				return dm.ErrState
			}
			v.State = member.StateApplyNotPass
			this._value.PresentBalance += v.CsnAmount + (-v.Amount)
			if _, err := this.Save(); err != nil {
				return err
			}
		}
		_, err := this.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}

// 完成提现
func (this *accountImpl) FinishApplyCash(id int, tradeNo string) error {
	v := this.GetBalanceInfo(id)
	if v.Kind == member.KindBalanceApplyCash {
		v.TradeNo = tradeNo
		v.State = member.StateApplyOver
		_, err := this.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}

// 冻结余额
func (this *accountImpl) Freezes(title string, tradeNo string, amount float32, referId int) error {
	if this._value.Balance < amount {
		return member.ErrNotEnoughAmount
	}
	if len(title) == 0 {
		title = "资金冻结"
	}
	v := &member.BalanceInfo{
		Kind:    member.KindBalanceFreezes,
		Title:   title,
		RefId:   referId,
		Amount:  amount,
		TradeNo: tradeNo,
		State:   member.StatusOK,
	}
	this._value.Balance -= amount
	this._value.FreezesFee += amount
	_, err := this.Save()
	if err == nil {
		_, err = this.SaveBalanceInfo(v)
	}
	return err
}

// 解冻金额
func (this *accountImpl) Unfreezes(title string, tradeNo string, amount float32, referId int) error {
	if this._value.FreezesFee < amount {
		return member.ErrNotEnoughAmount
	}
	if len(title) == 0 {
		title = "资金解结"
	}
	v := &member.BalanceInfo{
		Kind:    member.KindBalanceUnfreezes,
		Title:   title,
		RefId:   referId,
		Amount:  amount,
		TradeNo: tradeNo,
		State:   member.StatusOK,
	}
	this._value.Balance += amount
	this._value.FreezesFee -= amount
	_, err := this.Save()
	if err == nil {
		_, err = this.SaveBalanceInfo(v)
	}
	return err

}

// 冻结赠送金额
func (this *accountImpl) FreezesPresent(title string, tradeNo string, amount float32, referId int) error {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if this._value.PresentBalance < amount {
		return member.ErrNotEnoughAmount
	}
	if len(title) == 0 {
		title = "(赠送)资金冻结"
	}
	v := &member.BalanceInfo{
		Kind:    member.KindBalanceFreezesPresent,
		Title:   title,
		RefId:   referId,
		Amount:  -amount,
		TradeNo: tradeNo,
		State:   member.StatusOK,
	}
	this._value.PresentBalance -= amount
	this._value.FreezesPresent += amount
	_, err := this.Save()
	if err == nil {
		_, err = this.SaveBalanceInfo(v)
	}
	return err
}

// 解冻赠送金额
func (this *accountImpl) UnfreezesPresent(title string, tradeNo string, amount float32, referId int) error {
	if amount <= 0 {
		return member.ErrIncorrectAmount
	}
	if this._value.FreezesPresent < amount {
		return member.ErrNotEnoughAmount
	}
	if len(title) == 0 {
		title = "(赠送)资金解冻"
	}
	v := &member.BalanceInfo{
		Kind:    member.KindBalanceUnfreezesPresent,
		Title:   title,
		RefId:   referId,
		Amount:  amount,
		TradeNo: tradeNo,
		State:   member.StatusOK,
	}
	this._value.PresentBalance += amount
	this._value.FreezesPresent -= amount
	_, err := this.Save()
	if err == nil {
		_, err = this.SaveBalanceInfo(v)
	}
	return err
}

// 转账余额到其他账户
func (this *accountImpl) TransferBalance(kind int, amount float32,
	tradeNo string, toTitle, fromTitle string) error {
	var err error
	if kind == member.KindBalanceFlow {
		if this._value.Balance < amount {
			return member.ErrNotEnoughAmount
		}
		this._value.Balance -= amount
		this._value.FlowBalance += amount
		if _, err = this.Save(); err == nil {
			this.SaveBalanceInfo(&member.BalanceInfo{
				Kind:    member.KindBalanceTransfer,
				Title:   toTitle,
				Amount:  -amount,
				TradeNo: tradeNo,
				State:   member.StatusOK,
			})

			this.SaveBalanceInfo(&member.BalanceInfo{
				Kind:    member.KindBalanceTransfer,
				Title:   fromTitle,
				Amount:  amount,
				TradeNo: tradeNo,
				State:   member.StatusOK,
			})
		}
		return err
	}
	return member.ErrNotSupportTransfer
}

// 转账返利账户,kind为转账类型，如 KindBalanceTransfer等
// commission手续费
func (this *accountImpl) TransferPresent(kind int, amount float32, commission float32,
	tradeNo string, toTitle string, fromTitle string) error {
	var err error
	if kind == member.KindBalanceFlow {
		if this._value.Balance < amount {
			return member.ErrNotEnoughAmount
		}
		this._value.Balance -= amount
		this._value.FlowBalance += amount
		if _, err = this.Save(); err == nil {
			this.SaveBalanceInfo(&member.BalanceInfo{
				Kind:    member.KindBalanceTransfer,
				Title:   toTitle,
				Amount:  -amount,
				TradeNo: tradeNo,
				State:   member.StatusOK,
			})

			this.SaveBalanceInfo(&member.BalanceInfo{
				Kind:    kind,
				Title:   fromTitle,
				Amount:  amount,
				TradeNo: tradeNo,
				State:   member.StatusOK,
			})
		}
		return err
	}
	return member.ErrNotSupportTransfer
}

// 转账活动账户,kind为转账类型，如 KindBalanceTransfer等
// commission手续费
func (this *accountImpl) TransferFlow(kind int, amount float32, commission float32,
	tradeNo string, toTitle string, fromTitle string) error {
	var err error

	csnAmount := commission * amount
	finalAmount := amount - csnAmount

	if kind == member.KindBalancePresent {
		if this._value.FlowBalance < finalAmount {
			return member.ErrNotEnoughAmount
		}

		this._value.FlowBalance -= amount
		this._value.PresentBalance += finalAmount
		this._value.TotalPresentFee += finalAmount

		if _, err = this.Save(); err == nil {
			this.SaveBalanceInfo(&member.BalanceInfo{
				Kind:      member.KindBalanceTransfer,
				Title:     toTitle,
				Amount:    -amount,
				TradeNo:   tradeNo,
				CsnAmount: csnAmount,
				State:     member.StatusOK,
			})

			this.SaveBalanceInfo(&member.BalanceInfo{
				Kind:    kind,
				Title:   fromTitle,
				Amount:  finalAmount,
				TradeNo: tradeNo,
				State:   member.StatusOK,
			})
		}
		return err
	}

	return member.ErrNotSupportTransfer
}

// 将活动金转给其他人
func (this *accountImpl) TransferFlowTo(memberId int, kind int,
	amount float32, commission float32, tradeNo string,
	toTitle string, fromTitle string) error {

	var err error
	csnAmount := commission * amount
	finalAmount := amount + csnAmount // 转账方付手续费

	m := this._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	acc2 := m.GetAccount()

	if kind == member.KindBalanceFlow {
		if this._value.FlowBalance < finalAmount {
			return member.ErrNotEnoughAmount
		}

		this._value.FlowBalance -= finalAmount
		acc2.GetValue().FlowBalance += amount

		if _, err = this.Save(); err == nil {

			this.SaveBalanceInfo(&member.BalanceInfo{
				Kind:      member.KindBalanceTransfer,
				Title:     toTitle,
				Amount:    -finalAmount,
				CsnAmount: csnAmount,
				RefId:     memberId,
				TradeNo:   tradeNo,
				State:     member.StatusOK,
			})

			if _, err = acc2.Save(); err == nil {
				acc2.SaveBalanceInfo(&member.BalanceInfo{
					Kind:    kind,
					Title:   fromTitle,
					Amount:  amount,
					RefId:   this._value.MemberId,
					TradeNo: tradeNo,
					State:   member.StatusOK,
				})
			}
		}
		return err
	}

	return member.ErrNotSupportTransfer
}
