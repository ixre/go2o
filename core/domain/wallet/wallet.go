package wallet

import (
	"fmt"
	"github.com/jsix/gof/algorithm"
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/wallet"
	"go2o/core/infrastructure/domain"
	"strconv"
	"strings"
	"time"
)

var _ wallet.IWallet = new(WalletImpl)

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
	err := w.checkWallet()
	if err == nil {
		unix := time.Now().Unix()
		w._value.UpdateTime = unix
		// 初始化
		if w.GetAggregateRootId() <= 0 {
			w._value.CreateTime = unix
			w._value.State = wallet.StatNormal
		}
		// 保存
		id, err2 := util.I64Err(w._repo.SaveWallet_(w._value))
		if err2 != nil {
			return id, err
		}
		w._value.ID = id
	}
	return w.GetAggregateRootId(), err
}

// 检查钱包
func (w *WalletImpl) checkWallet() error {
	if w._value.UserId <= 0 {
		panic("incorrect wallet user id")
	}
	// 判断是否存在
	match := w._repo.CheckWalletUserMatch(w._value.UserId,
		w._value.WalletType, w.GetAggregateRootId())
	if match {
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
		OpUid:        opuId,
		OpName:       strings.TrimSpace(opuName),
		Remark:       "",
		ReviewState:  int(enum.ReviewPass),
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
		err := w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
}

// 支付抵扣,must是否必须大于0
func (w *WalletImpl) Discount(value int, title, outerNo string, must bool, opuId int, opuName string) error {
	err := w.checkValueOpu(value, false, opuId, opuName)
	if err == nil {
		if value > 0 {
			value = -value
		}
		if must && w._value.Balance < -value {
			return wallet.ErrOutOfAmount
		}
		w._value.Balance += value
		w._value.TotalPay += -value
		l := w.createWalletLog(wallet.KDiscount, value, title, opuId, opuName)
		l.OuterNo = outerNo
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
		if w._value.FreezeAmount < -value {
			return wallet.ErrOutOfAmount
		}
		w._value.Balance += value
		w._value.FreezeAmount += -value
		l := w.createWalletLog(wallet.KFreeze, value, title, opuId, opuName)
		l.OuterNo = outerNo
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
		err := w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
}

func (w *WalletImpl) Charge(value int, by int, title, outerNo string, opuId int, opuName string) error {
	err := w.checkValueOpu(value, false, opuId, opuName)
	if err == nil {
		if value < 0 {
			value = -value
		}
		var kind int = by
		// 用户或客服充值、才会计入累计充值记录
		switch by {
		case wallet.CUserCharge:
			kind = wallet.KindCharge
			w._value.TotalCharge += value
		case wallet.CServiceAgentCharge:
			kind = wallet.KindServiceAgentCharge
			w._value.TotalCharge += value
		case wallet.CSystemCharge:
			kind = wallet.KindSystemCharge
		case wallet.CRefundCharge:
			kind = wallet.KindPaymentOrderRefund
		case by < 10:
			panic("wallet can't charge by internal defined kind")
		}
		w._value.Balance += value
		// 保存日志
		l := w.createWalletLog(kind, value, title, opuId, opuName)
		l.OuterNo = outerNo
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
		if !(kind == wallet.KindPaymentOrderRefund ||
			kind == wallet.KindTransferRefund ||
			kind == wallet.KindTakeOutRefund) {
			panic("not support refund kind")
		}
		switch kind {
		// 扣减总支付金额
		case wallet.KindPaymentOrderRefund:
			w._value.TotalPay -= value
		}
		w._value.Balance += value
		// 保存日志
		l := w.createWalletLog(kind, value, title, opuId, opuName)
		l.OuterNo = outerNo
		err := w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
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
	wv := tw.Get()
	if wv.Balance < value+tradeFee {
		return wallet.ErrOutOfAmount
	}
	w._value.Balance -= value + tradeFee
	tradeNo := domain.NewTradeNo(000)
	l := w.createWalletLog(wallet.KTransferOut, value, title, 0, "")
	l.OuterNo = tradeNo
	l.Remark = remark
	err = w.saveWalletLog(l)
	if err == nil {
		if _, err = w.Save(); err == nil {
			err = tw.ReceiveTransfer(w.GetAggregateRootId(), value, tradeNo, toTitle)
		}
	}
	return err
}

func (w *WalletImpl) ReceiveTransfer(fromWalletId int64, value int, tradeNo string, title string) error {
	if value == 0 {
		return wallet.ErrAmountZero
	}
	if value < 0 {
		value = -value
	}
	w._value.Balance += value
	l := w.createWalletLog(wallet.KTransferOut, value, title, 0, "")
	l.OuterNo = tradeNo
	err := w.saveWalletLog(l)
	if err == nil {
		_, err = w.Save()
	}
	return err
}

// 申请提现,kind：提现方式,返回info_id,交易号 及错误,value为提现金额,tradeFee为手续费
func (w *WalletImpl) RequestTakeOut(value int, kind int, title string, tradeFee int) (int64, string, error) {
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
	tradeNo := domain.NewTradeNo(000)
	w._value.Balance -= value
	l := w.createWalletLog(kind, -(value - tradeFee), title, 0, "")
	l.TradeFee = -tradeFee
	l.OuterNo = tradeNo
	l.ReviewState = wallet.ReviewAwaiting
	l.ReviewRemark = ""
	err := w.saveWalletLog(l)
	if err == nil {
		_, err = w.Save()
	}
	return l.ID, l.OuterNo, err
}

func (w *WalletImpl) ReviewTakeOut(takeId int64, pass bool, remark string) error {
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
		err := w.Refund(l.TradeFee - l.Value, )
	}
}

func (w *WalletImpl) FinishTakeOut(takeId int32, outerNo string) error {
	panic("implement me")
}

func (w *WalletImpl) FreezeExpired(value int, remark string) error {
	panic("implement me")
}
