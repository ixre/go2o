package wallet

import (
	"bytes"
	"fmt"
	"github.com/ixre/gof/algorithm"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/wallet"
	"go2o/core/infrastructure/domain"
	"strconv"
	"strings"
	"time"
)

var _ wallet.IWallet = new(WalletImpl)

func NewWallet(v *wallet.Wallet, repo wallet.IWalletRepo) wallet.IWallet {
	return &WalletImpl{
		_value: v,
		_repo:  repo,
	}
}

type WalletImpl struct {
	_value *wallet.Wallet
	_repo  wallet.IWalletRepo
}

func (w *WalletImpl) GetAggregateRootId() int64 {
	return w._value.ID
}

func (w *WalletImpl) Hash() string {
	if w._value.HashCode == "" {
		str := fmt.Sprintf("%d&%d*%d", w._value.UserId, w._value.WalletType, time.Now().Unix)
		hash := algorithm.DJBHash([]byte(str))
		w._value.HashCode = strconv.Itoa(hash)
	}
	return w._value.HashCode
}

func (w *WalletImpl) NodeId() int {
	if w._value.NodeId <= 0 {
		w._value.NodeId = int(w._value.UserId * int64(w._value.WalletType) % 10)
	}
	return w._value.NodeId
}

func (w *WalletImpl) Get() wallet.Wallet {
	return *w._value
}

func (w *WalletImpl) State() int {
	return w._value.State
}

func (w *WalletImpl) getLog(logId int64) *wallet.WalletLog {
	return w._repo.GetLog(w.GetAggregateRootId(), logId)
}

func (w *WalletImpl) GetLog(logId int64) wallet.WalletLog {
	wl := w.getLog(logId)
	if wl != nil {
		return *wl
	}
	return wallet.WalletLog{}
}

func (w *WalletImpl) Save() (int64, error) {
	unix := time.Now().Unix()
	// 初始化
	if w.GetAggregateRootId() <= 0 {
		w.initWallet(unix)
	}
	err := w.checkWallet()
	// 保存
	if err == nil {
		w._value.UpdateTime = unix
		id, err2 := util.I64Err(w._repo.SaveWallet_(w._value))
		if err2 != nil {
			return id, err
		}
		w._value.ID = id
	}
	return w.GetAggregateRootId(), err
}
func (w *WalletImpl) initWallet(unix int64) {
	w._value.CreateTime = unix
	w._value.State = wallet.StatNormal
	w._value.WalletFlag = wallet.FlagCharge | wallet.FlagDiscount
	if w._value.WalletType <= 0 {
		w._value.WalletType = wallet.TPerson
	} else if w._value.WalletType != wallet.TMerchant &&
		w._value.WalletType != wallet.TPerson {
		panic("not support wallet type" + strconv.Itoa(w._value.WalletType))
	}
	if w._value.HashCode == "" {
		w.Hash()
	}
	if w._value.NodeId <= 0 {
		w.NodeId()
	}
}

// 检查钱包
func (w *WalletImpl) checkWallet() error {
	if w._value.UserId <= 0 {
		panic("incorrect wallet user id")
	}
	if flag := w._value.WalletFlag; flag <= 0 {
		panic("incorrect wallet flag:" + strconv.Itoa(flag))
	}
	if len(w._value.Remark) > 40 {
		return wallet.ErrRemarkLength
	}
	// 判断是否存在
	match := w._repo.CheckWalletUserMatch(w._value.UserId,
		w._value.WalletType, w.GetAggregateRootId())
	if !match {
		return wallet.ErrSingletonWallet
	}
	return nil
}

// 检查数值、操作人员信息
func (w *WalletImpl) checkValueOpu(value int, checkOpu bool, opuId int, opuName string) error {
	if value == 0 {
		return wallet.ErrAmountZero
	}
	if checkOpu && (opuId <= 0 || len(opuName) == 0) {
		return wallet.ErrMissingOperator
	}
	return w.checkWalletState(w, false)
}

// 检查钱包状态
func (w *WalletImpl) checkWalletState(iw wallet.IWallet, target bool) error {
	if iw == nil {
		if target {
			return wallet.ErrNoSuchTargetWalletAccount
		}
		return wallet.ErrNoSuchWalletAccount
	}
	switch iw.State() {
	case wallet.StatNormal:
		return nil
	case wallet.StatDisabled:
		if target {
			return wallet.ErrTargetWalletAccountNotService
		}
		return wallet.ErrWalletDisabled
	case wallet.StatClosed:
		if target {
			return wallet.ErrTargetWalletAccountNotService
		}
		return wallet.ErrWalletClosed
	}
	panic("unknown wallet state")
}

// 创建钱包日志
func (w *WalletImpl) createWalletLog(kind int, value int, title string, opuId int,
	opuName string) *wallet.WalletLog {
	unix := time.Now().Unix()
	return &wallet.WalletLog{
		WalletId:     w.GetAggregateRootId(),
		Kind:         kind,
		Title:        strings.TrimSpace(title),
		OuterChan:    "",
		OuterNo:      "",
		Value:        value,
		OperatorId:   opuId,
		OperatorName: strings.TrimSpace(opuName),
		Remark:       "",
		ReviewState:  wallet.ReviewPass,
		ReviewRemark: "",
		ReviewTime:   0,
		CreateTime:   unix,
		UpdateTime:   unix,
	}
}

// 保存钱包日志
func (w *WalletImpl) saveWalletLog(l *wallet.WalletLog) error {
	if l.Kind <= 0 {
		panic("wallet log kind error")
	}
	if l.Value == 0 {
		panic("incorrect value")
	}
	l.Title = strings.TrimSpace(l.Title)
	if l.Title == "" {
		panic("wallet log title can't empty")
	}
	id, err := util.I64Err(w._repo.SaveWalletLog_(l))
	if err == nil {
		l.ID = id
	}
	return err
}

// 调整余额，可能存在扣为负数的情况，需传入操作人员编号或操作人员名称
func (w *WalletImpl) Adjust(value int, title, outerNo string, opuId int, opuName string) error {
	err := w.checkValueOpu(value, true, opuId, opuName)
	if err == nil {
		w._value.AdjustAmount += value
		w._value.Balance += value
		l := w.createWalletLog(wallet.KAdjust, value, title, opuId, opuName)
		l.OuterNo = outerNo
		l.Balance = w._value.Balance
		err = w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
}

// 支付抵扣,must是否必须大于0
func (w *WalletImpl) Discount(value int, title, outerNo string, must bool) error {
	err := w.checkValueOpu(value, false, 0, "")
	if err == nil {
		if value > 0 {
			value = -value
		}
		if must && w._value.Balance < -value {
			return wallet.ErrOutOfAmount
		}
		w._value.Balance += value
		w._value.TotalPay += -value
		l := w.createWalletLog(wallet.KDiscount, value, title, 0, "")
		l.OuterNo = outerNo
		l.Balance = w._value.Balance
		err := w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
}

func (w *WalletImpl) Freeze(value int, title, outerNo string, opuId int, opuName string) error {
	err := w.checkValueOpu(value, false, opuId, opuName)
	if err == nil {
		if value > 0 {
			value = -value
		}
		if w._value.Balance < -value {
			return wallet.ErrOutOfAmount
		}
		w._value.Balance += value
		w._value.FreezeAmount += -value
		l := w.createWalletLog(wallet.KFreeze, value, title, opuId, opuName)
		l.OuterNo = outerNo
		l.Balance = w._value.Balance
		err := w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
}

func (w *WalletImpl) Unfreeze(value int, title, outerNo string, opuId int, opuName string) error {
	err := w.checkValueOpu(value, false, opuId, opuName)
	if err == nil {
		if value < 0 {
			value = -value
		}
		if w._value.FreezeAmount < value {
			return wallet.ErrOutOfAmount
		}
		w._value.Balance += value
		w._value.FreezeAmount += -value
		l := w.createWalletLog(wallet.KUnfreeze, value, title, opuId, opuName)
		l.OuterNo = outerNo
		l.Balance = w._value.Balance
		err := w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
}

func (w *WalletImpl) FreezeExpired(value int, remark string) error {
	if value == 0 {
		return wallet.ErrAmountZero
	}
	if value < 0 {
		value = -value
	}
	if w._value.FreezeAmount < value {
		return wallet.ErrOutOfAmount
	}
	w._value.FreezeAmount -= value
	w._value.ExpiredAmount += value
	l := w.createWalletLog(wallet.KExpired, -value, "过期失效", 0, "")
	err := w.saveWalletLog(l)
	if err == nil {
		_, err = w.Save()
	}
	return err
}

func (w *WalletImpl) Income(value int, tradeFee int, title, outerNo string) error {
	err := w.checkValueOpu(value, false, 0, "")
	if err == nil {
		if value < 0 {
			value = -value
		}
		if tradeFee < 0 {
			tradeFee = -tradeFee
		}
		w._value.Balance += value
		// 保存日志
		l := w.createWalletLog(wallet.KIncome, value, title, 0, "")
		l.OuterNo = outerNo
		l.TradeFee = -tradeFee
		l.ReviewState = wallet.ReviewPass
		l.ReviewTime = time.Now().Unix()
		l.Balance = w._value.Balance
		err = w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
}

func (w *WalletImpl) Charge(value int, by int, title, outerNo string, opuId int, opuName string) error {
	needOpuId := by == wallet.CServiceAgentCharge
	err := w.checkValueOpu(value, needOpuId, opuId, opuName)
	if err == nil {
		if value < 0 {
			value = -value
		}
		var kind int = by
		// 用户或客服充值、才会计入累计充值记录
		switch by {
		case wallet.CUserCharge:
			kind = wallet.KCharge
			w._value.TotalCharge += value
		case wallet.CServiceAgentCharge:
			kind = wallet.KServiceAgentCharge
			w._value.TotalCharge += value
		case wallet.CSystemCharge:
			kind = wallet.KIncome
		case wallet.CRefundCharge:
			kind = wallet.KPaymentOrderRefund
		default:
			if by < 20 {
				panic("wallet can't charge by internal defined kind")
			}
		}
		w._value.Balance += value
		// 保存日志
		l := w.createWalletLog(kind, value, title, opuId, opuName)
		l.OuterNo = outerNo
		l.Balance = w._value.Balance
		err := w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
}

func (w *WalletImpl) Refund(value int, kind int, title, outerNo string, opuId int, opuName string) error {
	err := w.checkValueOpu(value, false, opuId, opuName)
	if err == nil {
		if value < 0 {
			value = -value
		}
		if !(kind == wallet.KPaymentOrderRefund ||
			kind == wallet.KTransferRefund ||
			kind == wallet.KTakeOutRefund) {
			panic("not support refund kind")
		}
		switch kind {
		// 扣减总支付金额
		case wallet.KPaymentOrderRefund:
			w._value.TotalPay -= value
		}
		w._value.Balance += value
		// 保存日志
		l := w.createWalletLog(kind, value, title, opuId, opuName)
		l.OuterNo = outerNo
		l.Balance = w._value.Balance
		err := w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
}

func (w *WalletImpl) Transfer(toWalletId int64, value int, tradeFee int, title, toTitle, remark string) error {
	if value == 0 {
		return wallet.ErrAmountZero
	}
	if value < 0 {
		value = -value
	}
	if tradeFee < 0 {
		tradeFee = -tradeFee
	}
	var tw wallet.IWallet = w._repo.GetWallet(toWalletId)
	err := w.checkWalletState(tw, true)
	if err != nil {
		return err
	}
	// 验证金额
	if w._value.Balance < value+tradeFee {
		return wallet.ErrOutOfAmount
	}
	w._value.Balance -= value + tradeFee
	tradeNo := domain.NewTradeNo(8, int(w._value.UserId))
	l := w.createWalletLog(wallet.KTransferOut, -value, title, 0, "")
	l.TradeFee = -tradeFee
	l.OuterNo = tradeNo
	l.Remark = remark
	l.Balance = w._value.Balance
	err = w.saveWalletLog(l)
	if err == nil {
		if _, err = w.Save(); err == nil {
			err = tw.ReceiveTransfer(w.GetAggregateRootId(), value, tradeNo, toTitle, remark)
		}
	}
	return err
}

func (w *WalletImpl) ReceiveTransfer(fromWalletId int64, value int, tradeNo, title, remark string) error {
	if value == 0 {
		return wallet.ErrAmountZero
	}
	if value < 0 {
		value = -value
	}
	w._value.Balance += value
	l := w.createWalletLog(wallet.KTransferIn, value, title, 0, "")
	l.OuterNo = tradeNo
	l.Remark = remark
	l.Balance = w._value.Balance
	err := w.saveWalletLog(l)
	if err == nil {
		_, err = w.Save()
	}
	return err
}

// 申请提现,kind：提现方式,返回info_id,交易号 及错误,value为提现金额,tradeFee为手续费
func (w *WalletImpl) RequestTakeOut(value int, tradeFee int, kind int, title string) (int64, string, error) {
	if value == 0 {
		return 0, "", wallet.ErrAmountZero
	}
	if value < 0 {
		value = -value
	}
	if kind != wallet.KTakeOutToBankCard &&
		kind != wallet.KTakeOutToThirdPart && kind < 20 {
		return 0, "", wallet.ErrNotSupportTakeOutBusinessKind
	}
	// 判断是否暂停提现
	if wallet.TakeOutPause {
		return 0, "", wallet.ErrTakeOutPause
	}
	// 判断最低提现和最高提现
	if value < wallet.MinTakeOutAmount {
		return 0, "", wallet.ErrLessThanMinTakeAmount
	}
	if value > wallet.MaxTakeOutAmount {
		return 0, "", wallet.ErrMoreThanMinTakeAmount
	}
	// 余额是否不足
	if w._value.Balance < value+tradeFee {
		return 0, "", wallet.ErrOutOfAmount
	}
	tradeNo := domain.NewTradeNo(8, int(w._value.UserId))
	w._value.Balance -= value
	l := w.createWalletLog(kind, -(value - tradeFee), title, 0, "")
	l.TradeFee = -tradeFee
	l.OuterNo = tradeNo
	l.ReviewState = wallet.ReviewAwaiting
	l.ReviewRemark = ""
	l.Balance = w._value.Balance
	err := w.saveWalletLog(l)
	if err == nil {
		_, err = w.Save()
	}
	return l.ID, l.OuterNo, err
}

func (w *WalletImpl) ReviewTakeOut(takeId int64, pass bool, remark string, opuId int, opuName string) error {
	if err := w.checkValueOpu(1, true, opuId, opuName); err != nil {
		return err
	}
	l := w.getLog(takeId)
	if l == nil {
		return wallet.ErrNoSuchTakeOutLog
	}
	if l.ReviewState != wallet.ReviewAwaiting {
		return wallet.ErrTakeOutState
	}
	if pass {
		l.ReviewState = wallet.ReviewPass
	} else {
		l.ReviewRemark = remark
		l.ReviewState = wallet.ReviewReject
		err := w.Refund(-(l.TradeFee + l.Value), wallet.KTakeOutRefund, "提现退回",
			l.OuterNo, 0, "")
		if err != nil {
			return err
		}
	}
	l.OperatorId = opuId
	l.OperatorName = opuName
	l.UpdateTime = time.Now().Unix()
	return w.saveWalletLog(l)
}

func (w *WalletImpl) FinishTakeOut(takeId int64, outerNo string) error {
	l := w.getLog(takeId)
	if l == nil {
		return wallet.ErrNoSuchTakeOutLog
	}
	if l.ReviewState != wallet.ReviewPass {
		return wallet.ErrTakeOutState
	}
	l.OuterNo = outerNo
	l.ReviewState = wallet.ReviewConfirm
	l.Remark = "转款凭证:" + outerNo
	return w.saveWalletLog(l)
}

func (w *WalletImpl) PagingLog(begin int, over int, opt map[string]string, sort string) (int, []*wallet.WalletLog) {
	where := bytes.NewBuffer(nil)
	// 添加业务类型筛选
	if kind, ok := opt["kind"]; ok {
		where.WriteString(" AND kind IN (")
		where.WriteString(kind)
		where.WriteString(")")
	}
	// 添加审核状态条件
	if reviewState, ok := opt["review_state"]; ok {
		where.WriteString(" AND review_state=")
		where.WriteString(reviewState)
	}
	// 添加金额
	minAmount, ok1 := opt["min_amount"]
	maxAmount, ok2 := opt["max_amount"]
	if ok1 && ok2 {
		where.WriteString(" AND value BETWEEN ")
		where.WriteString(minAmount)
		where.WriteString(" AND ")
		where.WriteString(maxAmount)
	} else {
		if ok1 {
			where.WriteString(" AND value >= ")
			where.WriteString(minAmount)
		} else {
			where.WriteString(" AND value <= ")
			where.WriteString(minAmount)
		}
	}
	// 添加时间
	beginTime, ok1 := opt["begin_time"]
	overTime, ok2 := opt["over_time"]
	if ok1 && ok2 {
		where.WriteString(" AND create_time BETWEEN ")
		where.WriteString(beginTime)
		where.WriteString(" AND ")
		where.WriteString(overTime)
	} else {
		if ok1 {
			where.WriteString(" AND create_time >= ")
			where.WriteString(beginTime)
		} else {
			where.WriteString(" AND create_time <= ")
			where.WriteString(overTime)
		}
	}
	// 添加操作人员筛选
	if opuName, ok := opt["op_name"]; ok {
		where.WriteString(" AND op_name='")
		where.WriteString(opuName)
		where.WriteString("'")
	}
	if opuId, ok := opt["op_uid"]; ok {
		where.WriteString(" AND op_uid='")
		where.WriteString(opuId)
	}
	return w._repo.PagingWalletLog(w.GetAggregateRootId(), w.NodeId(),
		begin, over, where.String(), sort)
}
