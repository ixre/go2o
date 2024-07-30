package wallet

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/gof/algorithm"
	"github.com/ixre/gof/domain/eventbus"
	"github.com/ixre/gof/util"
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

func (w *WalletImpl) GetAggregateRootId() int {
	return w._value.Id
}

func (w *WalletImpl) Hash() string {
	if w._value.HashCode == "" {
		str := fmt.Sprintf("%d&%d*%d", w._value.UserId, w._value.WalletType, time.Now().Unix())
		hash := algorithm.DJBHash([]byte(str))
		w._value.HashCode = strconv.Itoa(hash)
	}
	return w._value.HashCode
}

func (w *WalletImpl) NodeId() int {
	if w._value.NodeId <= 0 {
		w._value.NodeId = int(w._value.UserId * w._value.WalletType % 10)
	}
	return w._value.NodeId
}

func (w *WalletImpl) Get() wallet.Wallet {
	return *w._value
}

func (w *WalletImpl) State() int {
	return int(w._value.State)
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

func (w *WalletImpl) Save() (int, error) {
	unix := time.Now().Unix()
	// 初始化
	if w.GetAggregateRootId() <= 0 {
		w.initWallet(unix)
	}
	err := w.checkWallet()
	// 保存
	if err == nil {
		w._value.UpdateTime = int(unix)
		id, err2 := util.I64Err(w._repo.SaveWallet_(w._value))
		if err2 != nil {
			return int(id), err
		}
		w._value.Id = int(id)
	}
	return w.GetAggregateRootId(), err
}
func (w *WalletImpl) initWallet(unix int64) {
	w._value.CreateTime = int(unix)
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
	if w._value.WalletName == "" || len(w._value.WalletName) > 40 {
		return wallet.ErrWalletName
	}
	// 判断是否存在
	match := w._repo.CheckWalletUserMatch(int64(w._value.UserId),
		w._value.WalletType, int64(w.GetAggregateRootId()))
	if !match {
		return wallet.ErrSingletonWallet
	}
	return nil
}

// 检查数值、操作人员信息
func (w *WalletImpl) checkValueOpu(value int, checkOpu bool, operatorUid int, operatorName string) error {
	if value == 0 {
		return wallet.ErrAmountZero
	}
	if checkOpu && (operatorUid <= 0 || len(operatorName) == 0) {
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
func (w *WalletImpl) createWalletLog(kind int, value int, title string, operatorUid int,
	operatorName string) *wallet.WalletLog {
	unix := time.Now().Unix()
	return &wallet.WalletLog{
		WalletId:     w.GetAggregateRootId(),
		WalletUser:   w._value.Username,
		Kind:         kind,
		Subject:      strings.TrimSpace(title),
		OuterChan:    "",
		OuterTxNo:    "",
		ChangeValue:  value,
		OprUid:       int(operatorUid),
		OprName:      strings.TrimSpace(operatorName),
		Remark:       "",
		ReviewStatus: wallet.ReviewPass,
		ReviewRemark: "",
		ReviewTime:   0,
		CreateTime:   int(unix),
		UpdateTime:   int(unix),
	}
}

// 保存钱包日志
func (w *WalletImpl) saveWalletLog(l *wallet.WalletLog) error {
	if l.Kind <= 0 {
		return errors.New("wallet log kind error")
	}
	if l.ChangeValue == 0 {
		return errors.New("incorrect value")
	}
	l.Subject = strings.TrimSpace(l.Subject)
	if l.Subject == "" {
		return errors.New("wallet log title can't empty")
	}
	isUpdate := l.Id > 0
	id, err := util.I64Err(w._repo.SaveWalletLog_(l))
	if err == nil {
		l.Id = int(id)
	}
	if w._value.WalletType == wallet.TPerson {
		eventbus.Publish(&events.AccountLogPushEvent{
			IsUpdateEvent: isUpdate,
			MemberId:      int(w._value.UserId),
			Account:       3,
			LogId:         int(id),
			LogKind:       l.Kind,
			Subject:       l.Subject,
			OuterNo:       l.OuterTxNo,
			ChangeValue:   int(l.ChangeValue),
			Balance:       int(l.Balance),
			ProcedureFee:  l.TransactionFee,
			ReviewStatus:  l.ReviewStatus,
			CreateTime:    int(l.CreateTime),
		})
	}
	return err
}

// Adjust 调整余额，可能存在扣为负数的情况，需传入操作人员编号或操作人员名称
func (w *WalletImpl) Adjust(value int, title, outerNo string,
	remark string, operatorUid int, operatorName string) error {
	err := w.checkValueOpu(value, true, operatorUid, operatorName)
	if err == nil {
		w._value.AdjustAmount += value
		w._value.Balance += value
		l := w.createWalletLog(wallet.KAdjust, value, title, operatorUid, operatorName)
		l.OuterTxNo = outerNo
		l.Remark = remark
		l.Balance = w._value.Balance
		err = w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
}

// Consume 消费
func (w *WalletImpl) Consume(amount int, title string, outerNo string, remark string) error {
	if amount > 0 {
		amount = -amount
	}
	if w._value.Balance < -amount {
		return wallet.ErrOutOfAmount
	}
	w._value.Balance += amount
	w._value.TotalPay += -amount
	l := w.createWalletLog(wallet.KConsume, amount, title, 0, "")
	l.OuterTxNo = outerNo
	l.Balance = w._value.Balance
	l.Remark = remark
	err := w.saveWalletLog(l)
	if err == nil {
		_, err = w.Save()
	}
	return err
}

// Discount 支付抵扣,must是否必须大于0
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
		l.OuterTxNo = outerNo
		l.Balance = w._value.Balance
		err := w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
}

func (w *WalletImpl) Freeze(data wallet.TransactionData, operator wallet.Operator) (int, error) {
	err := w.checkValueOpu(data.Amount, false, operator.OperatorUid, operator.OperatorName)
	if err == nil {
		if data.Amount > 0 {
			data.Amount = -data.Amount
		}
		if data.TransactionId > 0 {
			return w.refreeze(data)
		}
		return w.freeze(data, operator)
	}
	return 0, err
}

// 创建新的冻结记录
func (w *WalletImpl) freeze(data wallet.TransactionData, operator wallet.Operator) (int, error) {
	if w._value.Balance < -data.Amount {
		return 0, wallet.ErrOutOfAmount
	}
	w._value.Balance += data.Amount
	w._value.FreezeAmount += -data.Amount
	l := w.createWalletLog(wallet.KFreeze, data.Amount, data.TransactionTitle, operator.OperatorUid, operator.OperatorName)
	l.OuterTxNo = data.OuterTxNo
	l.Balance = w._value.Balance
	l.OuterTxUid = data.OuterTxUid
	err := w.saveWalletLog(l)
	if err == nil {
		_, err = w.Save()
	}
	return int(l.Id), err
}

// 对现有的冻结流水进行更新
func (w *WalletImpl) refreeze(data wallet.TransactionData) (int, error) {
	l := w._repo.GetWalletLog_(data.TransactionId)
	if l == nil || l.WalletId != w.GetAggregateRootId() {
		return 0, errors.New("冻结失败,交易日志不存在")
	}
	diffValue := int(math.Abs(float64(data.Amount - int(l.ChangeValue))))
	if w._value.Balance < diffValue {
		return 0, wallet.ErrOutOfAmount
	}
	w._value.Balance -= diffValue
	w._value.FreezeAmount += diffValue
	l.OuterTxNo = data.OuterTxNo
	l.ChangeValue = data.Amount
	l.Balance = w._value.Balance
	l.OuterTxUid = data.OuterTxUid
	err := w.saveWalletLog(l)
	if err == nil {
		_, err = w.Save()
	}
	return int(l.Id), err
}

func (w *WalletImpl) Unfreeze(value int, title, outerNo string, isRefundBalance bool, operatorUid int, operatorName string) error {
	err := w.checkValueOpu(value, false, operatorUid, operatorName)
	if err == nil {
		if value < 0 {
			value = -value
		}
		if w._value.FreezeAmount < value {
			return wallet.ErrOutOfAmount
		}
		if isRefundBalance {
			w._value.Balance += value
		}
		w._value.FreezeAmount += -value
		l := w.createWalletLog(wallet.KUnfreeze, value, title, operatorUid, operatorName)
		l.OuterTxNo = outerNo
		l.Balance = w._value.Balance
		l.Remark = ""
		err = w.saveWalletLog(l)
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

// CarryTo 入账
func (w *WalletImpl) CarryTo(tx wallet.TransactionData, review bool) (int, error) {
	err := w.checkValueOpu(tx.Amount, false, 0, "")
	if err == nil {
		if tx.Amount < 0 {
			tx.Amount = -tx.Amount
		}
		l := w.createWalletLog(wallet.KCarry, tx.Amount, tx.TransactionTitle, 0, "")
		// 扣除手续费，手续费无需冻结
		transactionFee := int(math.Abs(float64(tx.TransactionFee)))
		if transactionFee > 0 {
			// 减去手续费
			tx.Amount -= transactionFee
		}
		if review {
			w._value.FreezeAmount += tx.Amount
			l.ReviewStatus = wallet.ReviewPending
			l.ReviewRemark = "待审核"
		} else {
			w._value.Balance += tx.Amount
			l.ReviewStatus = wallet.ReviewPass
			l.ReviewTime = int(time.Now().Unix())
		}
		// 保存日志
		l.OuterTxNo = tx.OuterTxNo
		l.TransactionFee = -transactionFee
		l.Balance = w._value.Balance
		l.OuterTxUid = tx.OuterTxUid
		err = w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
		return int(l.Id), err
	}
	return 0, err
}

// ReviewCarryTo 审核入账
func (w *WalletImpl) ReviewCarryTo(transactionId int, pass bool, reason string) error {
	l := w._repo.GetLog(w.GetAggregateRootId(), int64(transactionId))
	if l.ReviewStatus != int(enum.ReviewPending) {
		return wallet.ErrNotSupport
	}
	w._value.FreezeAmount -= int(l.ChangeValue)
	l.UpdateTime = int(time.Now().Unix())
	if pass {
		w._value.Balance += l.ChangeValue
		l.ReviewStatus = int(enum.ReviewPass)
		l.Remark = "系统审核通过"
	} else {
		l.ReviewStatus = int(enum.ReviewReject)
		l.Remark = reason
	}
	err := w.saveWalletLog(l)
	if err == nil {
		_, err = w.Save()
	}
	return err
}

func (w *WalletImpl) Charge(value int, by int, title, outerNo string, remark string, operatorUid int, operatorName string) error {
	needOprUid := by == wallet.CServiceAgentCharge
	err := w.checkValueOpu(value, needOprUid, operatorUid, operatorName)
	if err == nil {
		if value < 0 {
			value = -value
		}
		var kind = by
		// 用户或客服充值、才会计入累计充值记录
		switch by {
		case wallet.CUserCharge, wallet.CServiceAgentCharge:
			kind = wallet.KCharge
			w._value.TotalCharge += value
		case wallet.CSystemCharge:
			kind = wallet.KCarry
		case wallet.CRefundCharge:
			kind = wallet.KPaymentOrderRefund
		default:
			if by < 20 {
				panic("wallet can't charge by internal defined kind")
			}
		}
		w._value.Balance += value
		// 保存日志
		l := w.createWalletLog(kind, value, title, operatorUid, operatorName)
		l.OuterTxNo = outerNo
		l.Remark = remark
		l.Balance = w._value.Balance
		err := w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
}

func (w *WalletImpl) Refund(value int, kind int, title, outerNo string, operatorUid int, operatorName string) error {
	err := w.checkValueOpu(value, false, operatorUid, operatorName)
	if err == nil {
		if value < 0 {
			value = -value
		}
		if !(kind == wallet.KPaymentOrderRefund ||
			kind == wallet.KTransferRefund ||
			kind == wallet.KWithdrawRefund) {
			panic("not support refund kind")
		}
		switch kind {
		// 扣减总支付金额
		case wallet.KPaymentOrderRefund:
			w._value.TotalPay -= value
		}
		w._value.Balance += value
		// 保存日志
		l := w.createWalletLog(kind, value, title, operatorUid, operatorName)
		l.OuterTxNo = outerNo
		l.Balance = w._value.Balance
		err := w.saveWalletLog(l)
		if err == nil {
			_, err = w.Save()
		}
	}
	return err
}

func (w *WalletImpl) Transfer(toWalletId int64, value int, transactionFee int, title, toTitle, remark string) error {
	if value == 0 {
		return wallet.ErrAmountZero
	}
	if value < 0 {
		value = -value
	}
	if transactionFee < 0 {
		transactionFee = -transactionFee
	}
	var tw = w._repo.GetWallet(int(toWalletId))
	err := w.checkWalletState(tw, true)
	if err != nil {
		return err
	}
	// 验证金额
	if w._value.Balance < value+transactionFee {
		return wallet.ErrOutOfAmount
	}
	w._value.Balance -= value + transactionFee
	tradeNo := domain.NewTradeNo(8, int(w._value.UserId))
	l := w.createWalletLog(wallet.KTransferOut, -value, title, 0, "")
	l.TransactionFee = -transactionFee
	l.OuterTxNo = tradeNo
	l.Remark = remark
	l.Balance = w._value.Balance
	err = w.saveWalletLog(l)
	if err == nil {
		if _, err = w.Save(); err == nil {
			err = tw.ReceiveTransfer(int64(w.GetAggregateRootId()), value, tradeNo, toTitle, remark)
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
	l.OuterTxNo = tradeNo
	l.Remark = remark
	l.Balance = w._value.Balance
	err := w.saveWalletLog(l)
	if err == nil {
		_, err = w.Save()
	}
	return err
}

// 申请提现,kind：提现方式,返回info_id,交易号 及错误,value为提现金额,transactionFee为手续费
func (w *WalletImpl) RequestWithdrawal(tx wallet.WithdrawTransaction) (int, string, error) {
	if tx.Amount == 0 {
		return 0, "", wallet.ErrAmountZero
	}
	if tx.Amount < 0 {
		tx.Amount = -tx.Amount
	}
	if tx.Kind != wallet.KWithdrawToBankCard &&
		tx.Kind != wallet.KWithdrawToPayWallet && tx.Kind < 30 {
		return 0, "", wallet.ErrNotSupportWithdrawKind
	}
	// 判断是否暂停提现
	if wallet.WithdrawIsPaused {
		return 0, "", wallet.ErrTakeOutPause
	}
	// 判断最低提现和最高提现
	if tx.Amount < wallet.MinWithdrawAmount {
		return 0, "", wallet.ErrLessThanMinWithdrawAmount
	}
	if tx.Amount > wallet.MaxWithdrawAmount {
		return 0, "", wallet.ErrMoreThanMinWithdrawAmount
	}
	// 余额是否不足
	if w._value.Balance < tx.Amount+tx.TransactionFee {
		return 0, "", wallet.ErrOutOfAmount
	}
	w._value.Balance -= tx.Amount
	l := w.createWalletLog(
		tx.Kind,
		-(tx.Amount - tx.TransactionFee),
		tx.TransactionTitle,
		0,
		"")
	l.TransactionFee = -tx.TransactionFee
	l.OuterTxNo = domain.NewTradeNo(8, int(w._value.UserId))
	l.ReviewStatus = wallet.ReviewPending
	l.ReviewRemark = ""
	l.BankName = tx.BankName
	l.AccountNo = tx.AccountNo
	l.AccountName = tx.AccountName
	l.Balance = w._value.Balance
	err := w.saveWalletLog(l)
	if err == nil {
		_, err = w.Save()
	}
	return l.Id, l.OuterTxNo, err
}

func (w *WalletImpl) ReviewWithdrawal(transactionId int, pass bool, remark string, operatorUid int, operatorName string) error {
	if err := w.checkValueOpu(1, true, operatorUid, operatorName); err != nil {
		return err
	}
	l := w.getLog(int64(transactionId))
	if l == nil {
		return wallet.ErrNoSuchAccountLog
	}
	if l.Kind != wallet.KWithdrawToBankCard && l.Kind != wallet.KWithdrawToPayWallet {
		return wallet.ErrNotSupport
	}
	if l.ReviewStatus != wallet.ReviewPending {
		return wallet.ErrWithdrawState
	}
	l.ReviewTime = int(time.Now().Unix())
	if pass {
		l.ReviewStatus = wallet.ReviewPass
	} else {
		l.ReviewRemark = remark
		l.ReviewStatus = wallet.ReviewReject
		err := w.Refund(-(l.TransactionFee + int(l.ChangeValue)), wallet.KWithdrawRefund, "提现退回",
			l.OuterTxNo, 0, "")
		if err != nil {
			return err
		}
	}
	l.OprUid = operatorUid
	l.OprName = operatorName
	l.UpdateTime = int(time.Now().Unix())
	return w.saveWalletLog(l)
}

func (w *WalletImpl) FinishWithdrawal(transactionId int, outerNo string) error {
	l := w.getLog(int64(transactionId))
	if l == nil {
		return wallet.ErrNoSuchAccountLog
	}
	if l.ReviewStatus != wallet.ReviewPass {
		return wallet.ErrWithdrawState
	}
	l.OuterTxNo = outerNo
	l.ReviewStatus = wallet.ReviewConfirm
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
	if reviewStatus, ok := opt["review_status"]; ok {
		where.WriteString(" AND review_status=")
		where.WriteString(reviewStatus)
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
	if operatorName, ok := opt["op_name"]; ok {
		where.WriteString(" AND op_name='")
		where.WriteString(operatorName)
		where.WriteString("'")
	}
	if operatorUid, ok := opt["op_uid"]; ok {
		where.WriteString(" AND op_uid='")
		where.WriteString(operatorUid)
	}
	return w._repo.PagingWalletLog(int64(w.GetAggregateRootId()), w.NodeId(),
		begin, over, where.String(), sort)
}
