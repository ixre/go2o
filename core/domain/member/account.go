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
	"fmt"
	"github.com/jsix/gof/db/orm"
	dm "go2o/core/domain"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/domain/tmp"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"math"
	"time"
)

var _ member.IAccount = new(accountImpl)

type accountImpl struct {
	_member  *memberImpl
	mm       member.IMemberManager
	_value   *member.Account
	_rep     member.IMemberRep
	valueRep valueobject.IValueRep
}

func NewAccount(m *memberImpl, value *member.Account,
	rep member.IMemberRep, mm member.IMemberManager,
	valueRep valueobject.IValueRep) member.IAccount {
	return &accountImpl{
		_member:  m,
		_value:   value,
		_rep:     rep,
		mm:       mm,
		valueRep: valueRep,
	}
}

// 获取领域对象编号
func (a *accountImpl) GetDomainId() int {
	return a._value.MemberId
}

// 获取账户值
func (a *accountImpl) GetValue() *member.Account {
	return a._value
}

// 保存
func (a *accountImpl) Save() (int, error) {
	a._value.UpdateTime = time.Now().Unix()
	return a._rep.SaveAccount(a._value)
}

// 设置优先(默认)支付方式, account 为账户类型
func (a *accountImpl) SetPriorityPay(account int, enabled bool) error {
	if enabled {
		support := false
		if account == member.AccountBalance ||
			account == member.AccountPresent ||
			account == member.AccountIntegral {
			support = true
		}
		if support {
			a._value.PriorityPay = account
			_, err := a.Save()
			return err
		}
		// 不支持支付的账号类型
		return member.ErrNotSupportPaymentAccountType
	}

	// 关闭默认支付
	a._value.PriorityPay = 0
	_, err := a.Save()
	return err
}

// 根据编号获取余额变动信息
func (a *accountImpl) GetBalanceInfo(id int) *member.BalanceInfo {
	return a._rep.GetBalanceInfo(id)
}

// 根据号码获取余额变动信息
func (a *accountImpl) GetBalanceInfoByNo(no string) *member.BalanceInfo {
	return a._rep.GetBalanceInfoByNo(no)
}

// 保存余额变动信息
func (a *accountImpl) SaveBalanceInfo(v *member.BalanceInfo) (int, error) {
	v.MemberId = a.GetDomainId()
	v.UpdateTime = time.Now().Unix()
	if v.CreateTime == 0 {
		v.CreateTime = v.UpdateTime
	}
	return a._rep.SaveBalanceInfo(v)
}

// 充值,客服充值时,需提供操作人(relateUser)
func (a *accountImpl) ChargeForBalance(chargeType int, title string, outerNo string,
	amount float32, relateUser int) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if chargeType == member.ChargeByService && relateUser <= 0 {
		return member.ErrNoSuchRelateUser
	}
	busKind := member.KindBalanceCharge
	switch chargeType {
	default:
		return member.ErrNotSupportChargeMethod
	case member.ChargeByUser:
		busKind = member.KindBalanceCharge
	case member.ChargeBySystem:
		busKind = member.KindBalanceSystemCharge
	case member.ChargeByService:
		busKind = member.KindBalanceServiceCharge
	case member.ChargeByRefund:
		busKind = member.KindBalanceRefund
	}
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: busKind,
		Title:        title,
		OuterNo:      outerNo,
		Amount:       amount,
		State:        1,
		RelateUser:   relateUser,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.saveBalanceLog(v)
	if err == nil {
		a._value.Balance += amount
		_, err = a.Save()
	}
	return err
}

// 保存余额日志
func (a *accountImpl) saveBalanceLog(v *member.BalanceLog) (int, error) {
	return orm.Save(tmp.Db().GetOrm(), v, v.Id)
}

// 保存赠送账户日志
func (a *accountImpl) savePresentLog(v *member.PresentLog) (int, error) {
	return orm.Save(tmp.Db().GetOrm(), v, v.Id)
}

// 根据编号获取余额变动信息
func (a *accountImpl) GetPresentLog(id int) *member.PresentLog {
	e := member.PresentLog{}
	if tmp.Db().GetOrm().Get(id, &e) == nil {
		return &e
	}
	return nil
}

// 扣减余额
func (a *accountImpl) DiscountBalance(title string, outerNo string,
	amount float32, relateUser int) (err error) {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if a._value.Balance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	kind := member.KindBalanceDiscount
	if relateUser > 0 {
		kind = member.KindBalanceServiceDiscount
	}
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: kind,
		Title:        title,
		OuterNo:      outerNo,
		Amount:       -amount,
		State:        1,
		RelateUser:   relateUser,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err = a.saveBalanceLog(v)
	if err == nil {
		a._value.Balance -= amount
		_, err = a.Save()
	}
	return err
}

// 冻结余额
func (a *accountImpl) Freeze(title string, outerNo string,
	amount float32, relateUser int) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if a._value.Balance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	if len(title) == 0 {
		title = "资金冻结"
	}
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindBalanceFreeze,
		Title:        title,
		Amount:       -amount,
		OuterNo:      outerNo,
		RelateUser:   relateUser,
		State:        member.StatusOK,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	a._value.Balance -= amount
	a._value.FreezeBalance += amount
	_, err := a.Save()
	if err == nil {
		_, err = a.saveBalanceLog(v)
	}
	return err
}

// 解冻金额
func (a *accountImpl) Unfreeze(title string, outerNo string,
	amount float32, relateUser int) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if a._value.FreezeBalance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	if len(title) == 0 {
		title = "资金解结"
	}
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindBalanceUnfreeze,
		Title:        title,
		RelateUser:   relateUser,
		Amount:       amount,
		OuterNo:      outerNo,
		State:        member.StatusOK,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	a._value.Balance += amount
	a._value.FreezeBalance -= amount
	_, err := a.Save()
	if err == nil {
		_, err = a.saveBalanceLog(v)
	}
	return err

}

// 赠送金额,客服操作时,需提供操作人(relateUser)
func (a *accountImpl) ChargeForPresent(title string, outerNo string,
	amount float32, relateUser int) error {
	kind := member.KindPresentAdd
	if relateUser > 0 {
		kind = member.KindPresentServiceAdd
	}
	return a.ChargePresentByKind(kind, title, outerNo, amount, relateUser)
}

// 赠送金额(指定业务类型)
func (a *accountImpl) ChargePresentByKind(kind int, title string,
	outerNo string, amount float32, relateUser int) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if title == "" {
		if amount < 0 {
			title = "赠送账户出账"
		} else {
			title = "赠送账户入账"
		}
	}
	unix := time.Now().Unix()
	v := &member.PresentLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: kind,
		Title:        title,
		OuterNo:      outerNo,
		Amount:       amount,
		State:        1,
		RelateUser:   relateUser,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.savePresentLog(v)
	if err == nil {
		a._value.PresentBalance += amount
		if amount > 0 {
			a._value.TotalPresentFee += amount
		}
		_, err = a.Save()
	}
	return err
}

// 扣减奖金,mustLargeZero是否必须大于0, 赠送金额存在扣为负数的情况
func (a *accountImpl) DiscountPresent(title string, outerNo string, amount float32,
	relateUser int, mustLargeZero bool) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if mustLargeZero && a._value.PresentBalance < amount {
		return member.ErrAccountNotEnoughAmount
	}

	if len(title) == 0 {
		title = "出账"
	}
	kind := member.KindPresentDiscount
	if relateUser > 0 {
		kind = member.KindPresentServiceDiscount
	}

	unix := time.Now().Unix()
	v := &member.PresentLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: kind,
		Title:        title,
		OuterNo:      outerNo,
		Amount:       -amount,
		State:        1,
		RelateUser:   relateUser,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.savePresentLog(v)
	if err == nil {
		a._value.PresentBalance -= amount
		_, err = a.Save()
	}
	return err
}

// 冻结赠送金额
func (a *accountImpl) FreezePresent(title string, outerNo string,
	amount float32, relateUser int) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if a._value.PresentBalance < amount {
		return member.ErrAccountNotEnoughAmount
	}
	if len(title) == 0 {
		title = "(赠送)资金冻结"
	}
	unix := time.Now().Unix()
	v := &member.PresentLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindPresentFreeze,
		Title:        title,
		RelateUser:   relateUser,
		Amount:       -amount,
		OuterNo:      outerNo,
		State:        member.StatusOK,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	a._value.PresentBalance -= amount
	a._value.FreezePresent += amount
	_, err := a.Save()
	if err == nil {
		_, err = a.savePresentLog(v)
	}
	return err
}

// 解冻赠送金额
func (a *accountImpl) UnfreezePresent(title string, outerNo string,
	amount float32, relateUser int) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if a._value.FreezePresent < amount {
		return member.ErrAccountNotEnoughAmount
	}
	if len(title) == 0 {
		title = "(赠送)资金解冻"
	}
	unix := time.Now().Unix()
	v := &member.PresentLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindPresentUnfreeze,
		Title:        title,
		RelateUser:   relateUser,
		Amount:       amount,
		OuterNo:      outerNo,
		State:        member.StatusOK,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	a._value.PresentBalance += amount
	a._value.FreezePresent -= amount
	_, err := a.Save()
	if err == nil {
		_, err = a.savePresentLog(v)
	}
	return err
}

// 流通账户余额充值，如扣除,amount传入负数金额
func (a *accountImpl) ChargeFlowBalance(title string, tradeNo string, amount float32) error {
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
	_, err := a.SaveBalanceInfo(v)
	if err == nil {
		a._value.FlowBalance += amount
		_, err = a.Save()
	}
	return err
}

// 支付单抵扣消费,tradeNo为支付单单号
func (a *accountImpl) PaymentDiscount(tradeNo string,
	amount float32, remark string) error {
	if amount < 0 || len(tradeNo) == 0 {
		return errors.New("amount error or missing trade no")
	}
	if amount > a._value.Balance {
		return member.ErrOutOfBalance
	}
	if remark == "" {
		remark = "支付抵扣"
	}

	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindBalanceDiscount,
		Title:        remark,
		OuterNo:      tradeNo,
		Amount:       -amount,
		State:        1,
		RelateUser:   member.DefaultRelateUser,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.saveBalanceLog(v)
	if err == nil {
		a._value.Balance -= amount
		_, err = a.Save()
	}
	return err
}

//　增加积分
func (a *accountImpl) AddIntegral(logType int, outerNo string, value int, remark string) error {
	if value <= 0 || math.IsNaN(float64(value)) {
		return member.ErrIncorrectQuota
	}
	if logType <= 0 {
		logType = member.TypeIntegralPresent
	}
	if logType == member.TypeIntegralShoppingPresent && outerNo == "" {
		return member.ErrMissingOuterNo
	}
	l := &member.IntegralLog{
		MemberId:   a._value.MemberId,
		Type:       logType,
		OuterNo:    outerNo,
		Value:      value,
		Remark:     remark,
		CreateTime: time.Now().Unix(),
	}
	err := a._rep.SaveIntegralLog(l)
	if err == nil {
		a._value.Integral += value
		_, err = a.Save()
	}
	return err
}

// 积分抵扣
func (a *accountImpl) IntegralDiscount(logType int, outerNo string,
	value int, remark string) error {
	if value <= 0 {
		return member.ErrIncorrectQuota
	}
	if a._value.Integral < value {
		return member.ErrNoSuchIntegral
	}

	if logType == member.TypeIntegralPaymentDiscount && outerNo == "" {
		return member.ErrMissingOuterNo
	}

	if logType <= 0 {
		logType = member.TypeIntegralDiscount
	}

	l := &member.IntegralLog{
		MemberId:   a._value.MemberId,
		Type:       logType,
		Value:      -value,
		OuterNo:    outerNo,
		Remark:     remark,
		CreateTime: time.Now().Unix(),
	}
	err := a._rep.SaveIntegralLog(l)
	if err == nil {
		a._value.Integral -= value
		_, err = a.Save()
	}
	return err
}

// 冻结积分,当new为true不扣除积分,反之扣除积分
func (a *accountImpl) FreezesIntegral(value int, new bool, remark string) error {
	if !new {
		if a._value.Integral < value {
			return member.ErrNoSuchIntegral
		}
		a._value.Integral -= value
	}
	a._value.FreezeIntegral += value
	_, err := a.Save()
	if err == nil {
		l := &member.IntegralLog{
			MemberId:   a._value.MemberId,
			Type:       member.TypeIntegralFreeze,
			Value:      -value,
			Remark:     remark,
			CreateTime: time.Now().Unix(),
		}
		err = a._rep.SaveIntegralLog(l)
	}
	return err
}

// 解冻积分
func (a *accountImpl) UnfreezesIntegral(value int, remark string) error {
	if a._value.FreezeIntegral < value {
		return member.ErrNoSuchIntegral
	}
	a._value.FreezeIntegral -= value
	a._value.Integral += value
	_, err := a.Save()
	if err == nil {
		l := &member.IntegralLog{
			MemberId:   a._value.MemberId,
			Type:       member.TypeIntegralUnfreeze,
			Value:      value,
			Remark:     remark,
			CreateTime: time.Now().Unix(),
		}
		err = a._rep.SaveIntegralLog(l)
	}
	return err
}

// 退款
func (a *accountImpl) RequestBackBalance(backType int, title string,
	amount float32) error {
	if amount > a._value.Balance {
		return member.ErrOutOfBalance
	}
	v := &member.BalanceInfo{
		Kind:   member.KindBalanceRefund,
		Type:   backType,
		Title:  title,
		Amount: amount,
		State:  0,
	}
	_, err := a.SaveBalanceInfo(v)
	if err == nil {
		a._value.Balance -= amount
		_, err = a.Save()
	}
	return err
}

// 完成退款
func (a *accountImpl) FinishBackBalance(id int, tradeNo string) error {
	v := a.GetBalanceInfo(id)
	if v.Kind == member.KindBalanceRefund {
		v.TradeNo = tradeNo
		v.State = 1
		_, err := a.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}

// 请求提现,返回info_id,交易号及错误
func (a *accountImpl) RequestTakeOut(businessKind int, title string,
	amount float32, commission float32) (int, string, error) {
	if businessKind != member.KindPresentTakeOutToBalance &&
		businessKind != member.KindPresentTakeOutToBankCard &&
		businessKind != member.KindPresentTakeOutToThirdPart {
		return 0, "", member.ErrNotSupportTakeOutBusinessKind
	}
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return 0, "", member.ErrIncorrectAmount
	}
	// 检测是否开启提现
	conf := a.valueRep.GetRegistry()
	if !conf.MemberTakeOutOn {
		return 0, "", errors.New(conf.MemberTakeOutMessage)
	}
	// 检测非正式会员提现
	lv := a.mm.LevelManager().GetLevelById(a._member.GetValue().Level)
	if lv != nil && lv.IsOfficial == 0 {
		return 0, "", errors.New(fmt.Sprintf(
			member.ErrTakeOutLevelNoPerm.Error(), lv.Name))
	}
	// 检测余额
	if a._value.PresentBalance < amount {
		return 0, "", member.ErrOutOfBalance
	}
	// 检测提现金额是否超过限制
	conf2 := a.valueRep.GetGlobNumberConf()
	if amount < conf2.MinTakeAmount {
		return 0, "", errors.New(fmt.Sprintf(member.ErrLessTakeAmount.Error(),
			format.FormatFloat(conf2.MinTakeAmount)))
	}
	if amount > conf2.MaxTakeAmount {
		return 0, "", errors.New(fmt.Sprintf(member.ErrOutTakeAmount.Error(),
			format.FormatFloat(conf2.MaxTakeAmount)))
	}

	tradeNo := domain.NewTradeNo(00000)
	csnAmount := amount * commission
	finalAmount := amount - csnAmount
	if finalAmount > 0 {
		finalAmount = -finalAmount
	}
	unix := time.Now().Unix()
	v := &member.PresentLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: businessKind,
		Title:        title,
		OuterNo:      tradeNo,
		Amount:       finalAmount,
		CsnFee:       csnAmount,
		State:        enum.ReviewAwaiting,
		RelateUser:   member.DefaultRelateUser,
		Remark:       "",
		CreateTime:   unix,
		UpdateTime:   unix,
	}

	// 提现至余额
	if businessKind == member.KindPresentTakeOutToBalance {
		a._value.Balance += amount
		v.State = enum.ReviewPass
	}

	id, err := a.savePresentLog(v)
	if err == nil {
		a._value.PresentBalance -= amount
		_, err = a.Save()
	}
	return id, tradeNo, err
}

// 确认提现
func (a *accountImpl) ConfirmTakeOut(id int, pass bool, remark string) error {
	v := a.GetPresentLog(id)
	if v.BusinessKind == member.KindPresentTakeOutToBankCard {
		if pass {
			v.State = enum.ReviewPass
		} else {
			if v.State == enum.ReviewReject {
				return dm.ErrState
			}
			v.Remark += "失败:" + remark
			v.State = enum.ReviewReject
			err := a.ChargePresentByKind(member.KindPresentTakOutRefund,
				"提现退回", v.OuterNo, v.CsnFee+(-v.Amount),
				member.DefaultRelateUser)
			if err != nil {
				return err
			}
			// 将手续费修改到提现金额上
			v.Amount -= v.CsnFee
			v.CsnFee = 0
		}
		v.UpdateTime = time.Now().Unix()
		_, err := a.savePresentLog(v)
		return err
	}
	return member.ErrNotSupportTakeOutBusinessKind
}

// 完成提现
func (a *accountImpl) FinishTakeOut(id int, tradeNo string) error {
	v := a.GetPresentLog(id)
	if v.BusinessKind == member.KindPresentTakeOutToBankCard {
		v.OuterNo = tradeNo
		v.State = enum.ReviewConfirm
		v.Remark = "银行凭证:" + tradeNo
		_, err := a.savePresentLog(v)
		return err
	}
	return member.ErrNotSupportTakeOutBusinessKind
}

// 将冻结金额标记为失效
func (a *accountImpl) FreezeExpired(accountKind int, amount float32, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	switch accountKind {
	case member.AccountBalance:
		return a.balanceFreezeExpired(amount, remark)
	case member.AccountPresent:
		return a.presentFreezeExpired(amount, remark)
	}
	return nil
}

func (a *accountImpl) balanceFreezeExpired(amount float32, remark string) error {
	if a._value.FreezeBalance < amount {
		return member.ErrIncorrectAmount
	}
	unix := time.Now().Unix()
	a._value.FreezeBalance -= amount
	a._value.ExpiredBalance += amount
	a._value.UpdateTime = unix
	l := &member.BalanceLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindBalanceExpired,
		Title:        "过期失效",
		OuterNo:      "",
		Amount:       amount,
		CsnFee:       0,
		State:        enum.ReviewPass,
		RelateUser:   member.DefaultRelateUser,
		Remark:       remark,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.saveBalanceLog(l)
	if err == nil {
		_, err = a.Save()
	}
	return err
}

func (a *accountImpl) presentFreezeExpired(amount float32, remark string) error {
	if a._value.FreezePresent < amount {
		return member.ErrIncorrectAmount
	}
	unix := time.Now().Unix()
	a._value.FreezePresent -= amount
	a._value.ExpiredPresent += amount
	a._value.UpdateTime = unix
	l := &member.PresentLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindPresentExpired,
		Title:        "过期失效",
		OuterNo:      "",
		Amount:       amount,
		CsnFee:       0,
		State:        enum.ReviewPass,
		RelateUser:   member.DefaultRelateUser,
		Remark:       remark,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.savePresentLog(l)
	if err == nil {
		_, err = a.Save()
	}
	return err
}

// 获取会员名称
func (a *accountImpl) getMemberName(m member.IMember) string {
	if tr := m.Profile().GetTrustedInfo(); tr.RealName != "" &&
		tr.Reviewed == enum.ReviewPass {
		return tr.RealName
	} else {
		return m.GetValue().Usr
	}
}

// 转账
func (a *accountImpl) TransferAccounts(accountKind int, toMember int, amount float32,
	csnRate float32, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	tm := a._rep.GetMember(toMember)
	if tm == nil {
		return member.ErrNoSuchMember
	}

	tradeNo := domain.NewTradeNo(00000)
	csnFee := amount * csnRate

	// 检测是否开启转账
	conf := a.valueRep.GetRegistry()
	if !conf.MemberTransferAccountsOn {
		return errors.New(conf.MemberTransferAccountsMessage)
	}

	switch accountKind {
	case member.AccountPresent:
		return a.transferPresent(tm, tradeNo, amount, csnFee, remark)
	case member.AccountBalance:
		return a.transferBalance(tm, tradeNo, amount, csnFee, remark)
	}
	return nil
}

func (a *accountImpl) transferBalance(tm member.IMember, tradeNo string,
	amount, csnFee float32, remark string) error {
	if a._value.Balance < amount+csnFee {
		return member.ErrAccountNotEnoughAmount
	}
	unix := time.Now().Unix()
	// 扣款
	toName := a.getMemberName(tm)
	l := &member.BalanceLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindBalanceTransferOut,
		Title:        "转账给" + toName,
		OuterNo:      tradeNo,
		Amount:       -amount,
		CsnFee:       csnFee,
		State:        enum.ReviewPass,
		RelateUser:   member.DefaultRelateUser,
		Remark:       remark,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.saveBalanceLog(l)
	if err == nil {
		a._value.Balance -= amount + csnFee
		a._value.UpdateTime = unix
		_, err = a.Save()
		if err == nil {
			err = tm.GetAccount().ReceiveTransfer(member.AccountBalance,
				a.GetDomainId(), tradeNo, amount, remark)
		}
	}
	return err
}

func (a *accountImpl) transferPresent(tm member.IMember, tradeNo string,
	amount, csnFee float32, remark string) error {
	// 检测非正式会员转账
	lv := a.mm.LevelManager().GetLevelById(a._member.GetValue().Level)
	if lv != nil && lv.IsOfficial == 0 {
		return errors.New(fmt.Sprintf(
			member.ErrTransferAccountsLevelNoPerm.Error(), lv.Name))
	}
	if a._value.PresentBalance < amount+csnFee {
		return member.ErrAccountNotEnoughAmount
	}
	unix := time.Now().Unix()
	// 扣款
	toName := a.getMemberName(tm)
	l := &member.PresentLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindPresentTransferOut,
		Title:        "转账给" + toName,
		OuterNo:      tradeNo,
		Amount:       -amount,
		CsnFee:       csnFee,
		State:        enum.ReviewPass,
		RelateUser:   member.DefaultRelateUser,
		Remark:       remark,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.savePresentLog(l)
	if err == nil {
		a._value.PresentBalance -= amount + csnFee
		a._value.UpdateTime = unix
		_, err = a.Save()
		if err == nil {
			err = tm.GetAccount().ReceiveTransfer(member.AccountPresent,
				a.GetDomainId(), tradeNo, amount, remark)
		}
	}
	return err
}

// 接收转账
func (a *accountImpl) ReceiveTransfer(accountKind int, fromMember int,
	tradeNo string, amount float32, remark string) error {
	switch accountKind {
	case member.AccountPresent:
		return a.receivePresentTransfer(fromMember, tradeNo, amount, remark)
	case member.AccountBalance:
		return a.receiveBalanceTransfer(fromMember, tradeNo, amount, remark)
	}
	return member.ErrNotSupportTransfer
}

func (a *accountImpl) receivePresentTransfer(fromMember int, tradeNo string,
	amount float32, remark string) error {
	fm := a._rep.GetMember(fromMember)
	if fm == nil {
		return member.ErrNoSuchMember
	}
	fromName := a.getMemberName(fm)
	unix := time.Now().Unix()
	tl := &member.PresentLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindPresentTransferIn,
		Title:        "转账收款（" + fromName + "）",
		OuterNo:      tradeNo,
		Amount:       amount,
		CsnFee:       0,
		State:        enum.ReviewPass,
		RelateUser:   member.DefaultRelateUser,
		Remark:       remark,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.savePresentLog(tl)
	if err == nil {
		a._value.PresentBalance += amount
		a._value.UpdateTime = unix
		_, err = a.Save()
	}
	return err
}

func (a *accountImpl) receiveBalanceTransfer(fromMember int, tradeNo string,
	amount float32, remark string) error {
	fromName := a.getMemberName(a._rep.GetMember(a.GetDomainId()))
	unix := time.Now().Unix()
	tl := &member.BalanceLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindBalanceTransferIn,
		Title:        "转账收款（" + fromName + "）",
		OuterNo:      tradeNo,
		Amount:       amount,
		CsnFee:       0,
		State:        enum.ReviewPass,
		RelateUser:   member.DefaultRelateUser,
		Remark:       remark,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.saveBalanceLog(tl)
	if err == nil {
		a._value.Balance += amount
		a._value.UpdateTime = unix
		_, err = a.Save()
	}
	return err
}

// 转账余额到其他账户
func (a *accountImpl) TransferBalance(kind int, amount float32,
	tradeNo string, toTitle, fromTitle string) error {
	var err error
	if kind == member.KindBalanceFlow {
		if a._value.Balance < amount {
			return member.ErrAccountNotEnoughAmount
		}
		a._value.Balance -= amount
		a._value.FlowBalance += amount
		if _, err = a.Save(); err == nil {
			a.SaveBalanceInfo(&member.BalanceInfo{
				Kind:    member.KindBalanceTransfer,
				Title:   toTitle,
				Amount:  -amount,
				TradeNo: tradeNo,
				State:   member.StatusOK,
			})

			a.SaveBalanceInfo(&member.BalanceInfo{
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
func (a *accountImpl) TransferPresent(kind int, amount float32, commission float32,
	tradeNo string, toTitle string, fromTitle string) error {
	var err error
	if kind == member.KindBalanceFlow {
		if a._value.Balance < amount {
			return member.ErrAccountNotEnoughAmount
		}
		a._value.Balance -= amount
		a._value.FlowBalance += amount
		if _, err = a.Save(); err == nil {
			a.SaveBalanceInfo(&member.BalanceInfo{
				Kind:    member.KindBalanceTransfer,
				Title:   toTitle,
				Amount:  -amount,
				TradeNo: tradeNo,
				State:   member.StatusOK,
			})

			a.SaveBalanceInfo(&member.BalanceInfo{
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
func (a *accountImpl) TransferFlow(kind int, amount float32, commission float32,
	tradeNo string, toTitle string, fromTitle string) error {
	var err error

	csnAmount := commission * amount
	finalAmount := amount - csnAmount

	if kind == member.KindPresentTransferIn {
		if a._value.FlowBalance < finalAmount {
			return member.ErrAccountNotEnoughAmount
		}

		a._value.FlowBalance -= amount
		a._value.PresentBalance += finalAmount
		a._value.TotalPresentFee += finalAmount

		if _, err = a.Save(); err == nil {
			a.SaveBalanceInfo(&member.BalanceInfo{
				Kind:      member.KindBalanceTransfer,
				Title:     toTitle,
				Amount:    -amount,
				TradeNo:   tradeNo,
				CsnAmount: csnAmount,
				State:     member.StatusOK,
			})

			a.SaveBalanceInfo(&member.BalanceInfo{
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
func (a *accountImpl) TransferFlowTo(memberId int, kind int,
	amount float32, commission float32, tradeNo string,
	toTitle string, fromTitle string) error {

	var err error
	csnAmount := commission * amount
	finalAmount := amount + csnAmount // 转账方付手续费

	m := a._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	acc2 := m.GetAccount()

	if kind == member.KindBalanceFlow {
		if a._value.FlowBalance < finalAmount {
			return member.ErrAccountNotEnoughAmount
		}

		a._value.FlowBalance -= finalAmount
		acc2.GetValue().FlowBalance += amount

		if _, err = a.Save(); err == nil {

			a.SaveBalanceInfo(&member.BalanceInfo{
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
					RefId:   a._value.MemberId,
					TradeNo: tradeNo,
					State:   member.StatusOK,
				})
			}
		}
		return err
	}

	return member.ErrNotSupportTransfer
}
