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
	member    *memberImpl
	mm        member.IMemberManager
	value     *member.Account
	rep       member.IMemberRepo
	valueRepo valueobject.IValueRepo
}

func NewAccount(m *memberImpl, value *member.Account,
	rep member.IMemberRepo, mm member.IMemberManager,
	valueRepo valueobject.IValueRepo) member.IAccount {
	return &accountImpl{
		member:    m,
		value:     value,
		rep:       rep,
		mm:        mm,
		valueRepo: valueRepo,
	}
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
	a.value.UpdateTime = time.Now().Unix()
	return a.rep.SaveAccount(a.value)
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

// 根据编号获取余额变动信息
func (a *accountImpl) GetBalanceInfo(id int32) *member.BalanceInfo {
	return a.rep.GetBalanceInfo(id)
}

// 根据号码获取余额变动信息
func (a *accountImpl) GetBalanceInfoByNo(no string) *member.BalanceInfo {
	return a.rep.GetBalanceInfoByNo(no)
}

// 保存余额变动信息
func (a *accountImpl) SaveBalanceInfo(v *member.BalanceInfo) (int32, error) {
	v.MemberId = a.GetDomainId()
	v.UpdateTime = time.Now().Unix()
	if v.CreateTime == 0 {
		v.CreateTime = v.UpdateTime
	}
	return a.rep.SaveBalanceInfo(v)
}

// 充值
func (a *accountImpl) Charge(account int32, kind int32, title, outerNo string,
	amount float32, relateUser int64) error {
	switch account {
	case member.AccountBalance:
		return a.chargeBalance(kind, title, outerNo, amount, relateUser)
	case member.AccountWallet:
		return a.chargeWallet(kind, title, outerNo, amount, relateUser)
	case member.AccountFlow:
		return a.chargeFlowBalance(title, outerNo, amount)
	}
	panic(errors.New("不支持的账户类型操作"))
}

func (a *accountImpl) Adjust(account int, title string, amount float32, remark string, relateUser int64) error {
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
	case member.AccountBalance:
		return a.adjustBalanceAccount(title, amount, remark, relateUser)
		//todo: 支持其他账户的调整
	}
	panic("not support other account adjust")
}

// 调整账户余额
func (a *accountImpl) adjustBalanceAccount(title string, amount float32, remark string, relateUser int64) error {
	unix := time.Now().Unix()
	v := &member.BalanceLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindAdjust,
		Title:        title,
		OuterNo:      "",
		Amount:       amount,
		Remark:       remark,
		RelateUser:   relateUser,
		State:        1,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	_, err := a.rep.SaveBalanceLog(v)
	if err == nil {
		a.value.Balance += amount
		_, err = a.Save()
	}
	return err
}

// 退款
func (a *accountImpl) Refund(account int, kind int32, title string,
	outerNo string, amount float32, relateUser int64) error {
	switch account {
	case member.AccountBalance:
		if kind != member.KindBalanceRefund {
			return member.ErrBusinessKind
		}
		return a.chargeBalanceNoLimit(kind, title, outerNo, amount, relateUser)
	case member.AccountWallet:
		if kind != member.KindWalletPaymentRefund &&
			kind != member.KindWalletTakeOutRefund {
			return member.ErrBusinessKind
		}
		return a.chargePresentNoLimit(kind, title, outerNo, amount, relateUser)
	}
	panic(errors.New("不支持的账户类型操作"))
}

// 充值余额
func (a *accountImpl) chargeBalance(kind int32, title string, outerNo string,
	amount float32, relateUser int64) error {
	switch kind {
	case member.ChargeByUser:
		kind = member.KindBalanceCharge
	case member.ChargeBySystem:
		kind = member.KindBalanceSystemCharge
	case member.ChargeByService:
		kind = member.KindBalanceServiceCharge
	}

	switch kind {
	case member.KindBalanceCharge,
		member.KindBalanceSystemCharge,
		member.KindBalanceServiceCharge:
		return a.chargeBalanceNoLimit(kind, title, outerNo,
			amount, relateUser)
	}
	return member.ErrNotSupportChargeMethod

}

// 充值,客服充值时,需提供操作人(relateUser)
func (a *accountImpl) chargeBalanceNoLimit(kind int32, title string, outerNo string,
	amount float32, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if relateUser <= 0 && kind == member.KindBalanceServiceCharge {
		return member.ErrNoSuchRelateUser
	}
	unix := time.Now().Unix()
	v := &member.BalanceLog{
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
	_, err := a.rep.SaveBalanceLog(v)
	if err == nil {
		a.value.Balance += amount
		_, err = a.Save()
	}
	return err
}

func (a *accountImpl) chargeWallet(kind int32, title string,
	outerNo string, amount float32, relateUser int64) error {
	switch kind {
	case member.ChargeBySystem:
		kind = member.KindWalletAdd
	case member.ChargeByService:
		kind = member.KindWalletServiceAdd
	}
	if kind < member.KindMine &&
		kind != member.KindWalletServiceAdd &&
		kind != member.KindWalletAdd {
		return member.ErrBusinessKind
	}
	return a.chargePresentNoLimit(kind, title, outerNo,
		amount, relateUser)
}

// 赠送金额(指定业务类型)
func (a *accountImpl) chargePresentNoLimit(kind int32, title string,
	outerNo string, amount float32, relateUser int64) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	// 客服操作
	if relateUser == 0 && (kind == member.KindWalletServiceAdd) {
		return member.ErrNoSuchRelateUser
	}

	if title == "" {
		if amount < 0 {
			title = "钱包账户出账"
		} else {
			title = "钱包账户入账"
		}
	}
	unix := time.Now().Unix()
	v := &member.MWalletLog{
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
	_, err := a.rep.SavePresentLog(v)
	if err == nil {
		a.value.WalletBalance += amount
		// 退款不能加入到累计赠送金额
		if kind != member.KindWalletTakeOutRefund &&
			kind != member.KindWalletPaymentRefund &&
			amount > 0 {
			a.value.TotalPresentFee += amount
		}
		_, err = a.Save()
	}
	return err
}

// 根据编号获取余额变动信息
func (a *accountImpl) GetWalletLog(id int32) *member.MWalletLog {
	e := member.MWalletLog{}
	if tmp.Db().GetOrm().Get(id, &e) == nil {
		return &e
	}
	return nil
}

// 扣减余额
func (a *accountImpl) DiscountBalance(title string, outerNo string,
	amount float32, relateUser int64) (err error) {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if a.value.Balance < amount {
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
	_, err = a.rep.SaveBalanceLog(v)
	if err == nil {
		a.value.Balance -= amount
		_, err = a.Save()
	}
	return err
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
	a.value.Balance += amount
	a.value.FreezeBalance -= amount
	_, err := a.Save()
	if err == nil {
		_, err = a.rep.SaveBalanceLog(v)
	}
	return err

}

// 扣减奖金,mustLargeZero是否必须大于0, 赠送金额存在扣为负数的情况
func (a *accountImpl) DiscountWallet(title string, outerNo string, amount float32,
	relateUser int64, mustLargeZero bool) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	if mustLargeZero && a.value.WalletBalance < amount {
		return member.ErrAccountNotEnoughAmount
	}

	if len(title) == 0 {
		title = "出账"
	}
	kind := member.KindWalletDiscount
	if relateUser > 0 {
		kind = member.KindWalletServiceDiscount
	}

	unix := time.Now().Unix()
	v := &member.MWalletLog{
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
	_, err := a.rep.SavePresentLog(v)
	if err == nil {
		a.value.WalletBalance -= amount
		_, err = a.Save()
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
	v := &member.MWalletLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindWalletFreeze,
		Title:        title,
		RelateUser:   relateUser,
		Amount:       -amount,
		OuterNo:      outerNo,
		State:        member.StatusOK,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	a.value.WalletBalance -= amount
	a.value.FreezeWallet += amount
	_, err := a.Save()
	if err == nil {
		_, err = a.rep.SavePresentLog(v)
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
	v := &member.MWalletLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindWalletUnfreeze,
		Title:        title,
		RelateUser:   relateUser,
		Amount:       amount,
		OuterNo:      outerNo,
		State:        member.StatusOK,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
	a.value.WalletBalance += amount
	a.value.FreezeWallet -= amount
	_, err := a.Save()
	if err == nil {
		_, err = a.rep.SavePresentLog(v)
	}
	return err
}

// 流通账户余额充值，如扣除,amount传入负数金额
func (a *accountImpl) chargeFlowBalance(title string,
	tradeNo string, amount float32) error {
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
		a.value.FlowBalance += amount
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
	if amount > a.value.Balance {
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
	_, err := a.rep.SaveBalanceLog(v)
	if err == nil {
		a.value.Balance -= amount
		_, err = a.Save()
	}
	return err
}

//　增加积分
func (a *accountImpl) AddIntegral(logType int, outerNo string,
	value int64, remark string) error {
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
		MemberId:   a.value.MemberId,
		Type:       logType,
		OuterNo:    outerNo,
		Value:      value,
		Remark:     remark,
		CreateTime: time.Now().Unix(),
	}
	err := a.rep.SaveIntegralLog(l)
	if err == nil {
		a.value.Integral += value
		_, err = a.Save()
	}
	return err
}

// 积分抵扣
func (a *accountImpl) IntegralDiscount(logType int, outerNo string,
	value int64, remark string) error {
	if value <= 0 {
		return member.ErrIncorrectQuota
	}
	if a.value.Integral < value {
		return member.ErrNoSuchIntegral
	}

	if logType == member.TypeIntegralPaymentDiscount && outerNo == "" {
		return member.ErrMissingOuterNo
	}

	if logType <= 0 {
		logType = member.TypeIntegralDiscount
	}

	l := &member.IntegralLog{
		MemberId:   a.value.MemberId,
		Type:       logType,
		Value:      -value,
		OuterNo:    outerNo,
		Remark:     remark,
		CreateTime: time.Now().Unix(),
	}
	err := a.rep.SaveIntegralLog(l)
	if err == nil {
		a.value.Integral -= value
		_, err = a.Save()
	}
	return err
}

// 冻结积分,当new为true不扣除积分,反之扣除积分
func (a *accountImpl) FreezesIntegral(value int64, new bool, remark string) error {
	if !new {
		if a.value.Integral < value {
			return member.ErrNoSuchIntegral
		}
		a.value.Integral -= value
	}
	a.value.FreezeIntegral += value
	_, err := a.Save()
	if err == nil {
		l := &member.IntegralLog{
			MemberId:   a.value.MemberId,
			Type:       member.TypeIntegralFreeze,
			Value:      -value,
			Remark:     remark,
			CreateTime: time.Now().Unix(),
		}
		err = a.rep.SaveIntegralLog(l)
	}
	return err
}

// 解冻积分
func (a *accountImpl) UnfreezesIntegral(value int64, remark string) error {
	if a.value.FreezeIntegral < value {
		return member.ErrNoSuchIntegral
	}
	a.value.FreezeIntegral -= value
	a.value.Integral += value
	_, err := a.Save()
	if err == nil {
		l := &member.IntegralLog{
			MemberId:   a.value.MemberId,
			Type:       member.TypeIntegralUnfreeze,
			Value:      value,
			Remark:     remark,
			CreateTime: time.Now().Unix(),
		}
		err = a.rep.SaveIntegralLog(l)
	}
	return err
}

// 退款
func (a *accountImpl) RequestBackBalance(backType int, title string,
	amount float32) error {
	if amount > a.value.Balance {
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
		a.value.Balance -= amount
		_, err = a.Save()
	}
	return err
}

// 完成退款
func (a *accountImpl) FinishBackBalance(id int32, tradeNo string) error {
	v := a.GetBalanceInfo(id)
	if v == nil || v.MemberId != a.value.MemberId {
		return member.ErrIncorrectInfo
	}
	if v.Kind == member.KindBalanceRefund {
		v.TradeNo = tradeNo
		v.State = 1
		_, err := a.SaveBalanceInfo(v)
		return err
	}
	return errors.New("kind not match")
}

// 请求提现,返回info_id,交易号及错误
func (a *accountImpl) RequestTakeOut(takeKind int32, title string,
	amount float32, commission float32) (int32, string, error) {
	if takeKind != member.KindWalletTakeOutToBalance &&
		takeKind != member.KindWalletTakeOutToBankCard &&
		takeKind != member.KindWalletTakeOutToThirdPart {
		return 0, "", member.ErrNotSupportTakeOutBusinessKind
	}
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return 0, "", member.ErrIncorrectAmount
	}
	// 检测是否开启提现
	conf := a.valueRepo.GetRegistry()
	if !conf.MemberTakeOutOn {
		return 0, "", errors.New(conf.MemberTakeOutMessage)
	}

	// 检测是否实名
	if conf.TakeOutMustTrust {
		trust := a.member.Profile().GetTrustedInfo()
		if trust.ReviewState != enum.ReviewPass {
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
	conf2 := a.valueRepo.GetGlobNumberConf()
	if amount < conf2.MinTakeOutAmount {
		return 0, "", errors.New(fmt.Sprintf(member.ErrLessTakeAmount.Error(),
			format.FormatFloat(conf2.MinTakeOutAmount)))
	}
	if amount > conf2.MaxTakeOutAmount {
		return 0, "", errors.New(fmt.Sprintf(member.ErrOutTakeAmount.Error(),
			format.FormatFloat(conf2.MaxTakeOutAmount)))
	}
	// 检测是否超过限制
	if maxTimes := conf2.MaxTakeOutTimesOfDay; maxTimes > 0 {
		takeTimes := a.rep.GetTodayTakeOutTimes(a.GetDomainId())
		if takeTimes >= maxTimes {
			return 0, "", member.ErrAccountOutOfTakeOutTimes
		}
	}

	tradeNo := domain.NewTradeNo(8, int(a.member.GetAggregateRootId()))
	csnAmount := amount * commission
	finalAmount := amount - csnAmount
	if finalAmount > 0 {
		finalAmount = -finalAmount
	}
	unix := time.Now().Unix()
	v := &member.MWalletLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: takeKind,
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
	if takeKind == member.KindWalletTakeOutToBalance {
		a.value.Balance += amount
		v.State = enum.ReviewPass
	}
	a.value.WalletBalance -= amount
	_, err := a.Save()
	if err == nil {
		go a.rep.AddTodayTakeOutTimes(a.GetDomainId())
		id, err := a.rep.SavePresentLog(v)
		return id, tradeNo, err
	}
	return 0, tradeNo, err
}

// 确认提现
func (a *accountImpl) ConfirmTakeOut(id int32, pass bool, remark string) error {
	v := a.GetWalletLog(id)
	if v == nil || v.MemberId != a.value.MemberId {
		return member.ErrIncorrectInfo
	}
	if v.State != enum.ReviewAwaiting {
		return member.ErrTakeOutState
	}
	if pass {
		v.State = enum.ReviewPass
	} else {
		v.Remark += "失败:" + remark
		v.State = enum.ReviewReject
		err := a.Refund(member.AccountWallet,
			member.KindWalletTakeOutRefund,
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
	_, err := a.rep.SavePresentLog(v)
	return err
	//if v.BusinessKind == member.KindWalletTakeOutToBankCard {
	//	if pass {
	//		v.State = enum.ReviewPass
	//	} else {
	//		if v.State == enum.ReviewReject {
	//			return dm.ErrState
	//		}
	//		v.Remark += "失败:" + remark
	//		v.State = enum.ReviewReject
	//		err := a.chargePresentByKind(member.KindWalletTakOutRefund,
	//			"提现退回", v.OuterNo, v.CsnFee+(-v.Amount),
	//			member.DefaultRelateUser)
	//		if err != nil {
	//			return err
	//		}
	//		// 将手续费修改到提现金额上
	//		v.Amount -= v.CsnFee
	//		v.CsnFee = 0
	//	}
	//	v.UpdateTime = time.Now().Unix()
	//	_, err := a.rep.SavePresentLog(v)
	//	return err
	//}
	//return member.ErrNotSupportTakeOutBusinessKind
}

// 完成提现
func (a *accountImpl) FinishTakeOut(id int32, tradeNo string) error {
	v := a.GetWalletLog(id)
	if v == nil || v.MemberId != a.value.MemberId {
		return member.ErrIncorrectInfo
	}
	if v.State != enum.ReviewPass {
		return member.ErrTakeOutState
	}
	v.OuterNo = tradeNo
	v.State = enum.ReviewConfirm
	v.Remark = "转款凭证:" + tradeNo
	_, err := a.rep.SavePresentLog(v)
	return err

	//if v.BusinessKind == member.KindWalletTakeOutToBankCard {
	//    v.OuterNo = tradeNo
	//    v.State = enum.ReviewConfirm
	//    v.Remark = "银行凭证:" + tradeNo
	//    _, err := a.rep.SavePresentLog(v)
	//    return err
	//}
	//return member.ErrNotSupportTakeOutBusinessKind
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
	a.value.ExpiredPresent += amount
	a.value.UpdateTime = unix
	l := &member.MWalletLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindWalletExpired,
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
	_, err := a.rep.SavePresentLog(l)
	if err == nil {
		_, err = a.Save()
	}
	return err
}

// 获取会员名称
func (a *accountImpl) getMemberName(m member.IMember) string {
	if tr := m.Profile().GetTrustedInfo(); tr.RealName != "" &&
		tr.ReviewState == enum.ReviewPass {
		return tr.RealName
	} else {
		return m.GetValue().Usr
	}
}

// 转账
func (a *accountImpl) TransferAccount(accountKind int, toMember int64, amount float32,
	csnRate float32, remark string) error {
	if amount <= 0 || math.IsNaN(float64(amount)) {
		return member.ErrIncorrectAmount
	}
	tm := a.rep.GetMember(toMember)
	if tm == nil {
		return member.ErrNoSuchMember
	}

	tradeNo := domain.NewTradeNo(8, int(a.member.GetAggregateRootId()))
	csnFee := amount * csnRate

	// 检测是否开启转账
	keys := []string{valueobject.RKMemberTransferAccountsOn,
		valueobject.RKMemberTransferAccountsMessage}
	registry := a.valueRepo.GetsRegistryMap(keys)
	if b := registry[keys[0]]; b != "true" && b != "1" {
		return errors.New(registry[keys[1]])
	}

	switch accountKind {
	case member.AccountWallet:
		return a.transferPresent(tm, tradeNo, amount, csnFee, remark)
	case member.AccountBalance:
		return a.transferBalance(tm, tradeNo, amount, csnFee, remark)
	}
	return nil
}

func (a *accountImpl) transferBalance(tm member.IMember, tradeNo string,
	amount, csnFee float32, remark string) error {
	if a.value.Balance < amount+csnFee {
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

func (a *accountImpl) transferPresent(tm member.IMember, tradeNo string,
	amount, csnFee float32, remark string) error {
	// 检测非正式会员转账
	lv := a.mm.LevelManager().GetLevelById(a.member.GetValue().Level)
	if lv != nil && lv.IsOfficial == 0 {
		return errors.New(fmt.Sprintf(
			member.ErrTransferAccountsLevelNoPerm.Error(), lv.Name))
	}
	if a.value.WalletBalance < amount+csnFee {
		return member.ErrAccountNotEnoughAmount
	}
	unix := time.Now().Unix()
	// 扣款
	toName := a.getMemberName(tm)
	l := &member.MWalletLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindWalletTransferOut,
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
	_, err := a.rep.SavePresentLog(l)
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
	tl := &member.MWalletLog{
		MemberId:     a.GetDomainId(),
		BusinessKind: member.KindWalletTransferIn,
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
	_, err := a.rep.SavePresentLog(tl)
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
	_, err := a.rep.SaveBalanceLog(tl)
	if err == nil {
		a.value.Balance += amount
		a.value.UpdateTime = unix
		_, err = a.Save()
	}
	return err
}

// 转账余额到其他账户
func (a *accountImpl) TransferBalance(kind int32, amount float32,
	tradeNo string, toTitle, fromTitle string) error {
	var err error
	if kind == member.KindBalanceFlow {
		if a.value.Balance < amount {
			return member.ErrAccountNotEnoughAmount
		}
		a.value.Balance -= amount
		a.value.FlowBalance += amount
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
func (a *accountImpl) TransferWallet(kind int32, amount float32, commission float32,
	tradeNo string, toTitle string, fromTitle string) error {
	var err error
	if kind == member.KindBalanceFlow {
		if a.value.Balance < amount {
			return member.ErrAccountNotEnoughAmount
		}
		a.value.Balance -= amount
		a.value.FlowBalance += amount
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
func (a *accountImpl) TransferFlow(kind int32, amount float32, commission float32,
	tradeNo string, toTitle string, fromTitle string) error {
	var err error

	csnAmount := commission * amount
	finalAmount := amount - csnAmount

	if kind == member.KindWalletTransferIn {
		if a.value.FlowBalance < finalAmount {
			return member.ErrAccountNotEnoughAmount
		}

		a.value.FlowBalance -= amount
		a.value.WalletBalance += finalAmount
		a.value.TotalPresentFee += finalAmount

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
func (a *accountImpl) TransferFlowTo(memberId int64, kind int32,
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

	if kind == member.KindBalanceFlow {
		if a.value.FlowBalance < finalAmount {
			return member.ErrAccountNotEnoughAmount
		}

		a.value.FlowBalance -= finalAmount
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
					RefId:   a.value.MemberId,
					TradeNo: tradeNo,
					State:   member.StatusOK,
				})
			}
		}
		return err
	}

	return member.ErrNotSupportTransfer
}
