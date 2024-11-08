/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: transaction_manager.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2017-01-27 12:23:22
 * description: 商户交易管理
 * history:
 */
package merchant

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/invoice"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	rbac "github.com/ixre/go2o/core/domain/interface/rabc"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
	"github.com/ixre/go2o/core/infrastructure/util"
	"github.com/ixre/gof/domain/eventbus"
	"github.com/ixre/gof/typeconv"
)

var _ merchant.IMerchantTransactionManager = new(transactionManagerImpl)

// 商户交易管理
type transactionManagerImpl struct {
	mch           *merchantImpl
	mchId         int
	_mchRepo      merchant.IMerchantRepo
	_invoiceRepo  invoice.IInvoiceRepo
	_billRepo     fw.Repository[merchant.MerchantBill]
	currentBill   *merchant.MerchantBill
	_rbacRepo     rbac.IRbacRepository
	_registryRepo registry.IRegistryRepo
	_walletRepo   wallet.IWalletRepo
	_walletId     int
}

func newTransactionManagerImpl(id int, m *merchantImpl, mchRepo merchant.IMerchantRepo,
	invoiceRepo invoice.IInvoiceRepo,
	walletRepo wallet.IWalletRepo,
	rbacRepo rbac.IRbacRepository,
	registryRepo registry.IRegistryRepo) merchant.IMerchantTransactionManager {
	return &transactionManagerImpl{
		mchId:         id,
		mch:           m,
		_mchRepo:      mchRepo,
		_billRepo:     m._repo.BillRepo(),
		_invoiceRepo:  invoiceRepo,
		_rbacRepo:     rbacRepo,
		_registryRepo: registryRepo,
		_walletRepo:   walletRepo,
	}
}

// GetSettlementPeriod 获取结算周期
func (s *transactionManagerImpl) GetSettlementPeriod() int {
	period, _ := s._registryRepo.GetValue(registry.MerchantSettlementPeriod)
	i := typeconv.Int(period)
	if i == 0 {
		// 默认月结
		return 2
	}
	return i
}

// getWalletId 获取商户钱包ID
func (s *transactionManagerImpl) getWalletId() int {
	if s._walletId == 0 {
		iw := s._walletRepo.GetWalletByUserId(int64(s.mchId), wallet.TMerchant)
		if iw != nil {
			s._walletId = iw.GetAggregateRootId()
		}
	}
	return s._walletId
}

// GetBillByTime implements merchant.IMerchantTransactionManager.
func (s *transactionManagerImpl) GetBillByTime(billTime int) merchant.IBillDomain {
	month, _ := fw.GetMonthStartEndUnix(int64(billTime))
	bill := s._billRepo.FindBy("mch_id = ? AND bill_time = ?", s.mchId, month)
	if bill != nil {
		return newBillDomainImpl(bill, s._mchRepo, s._rbacRepo, s._registryRepo, s)
	}
	return nil
}

// GetBill 获取指定账单
func (s *transactionManagerImpl) GetBill(billId int) merchant.IBillDomain {
	v := s._billRepo.Get(billId)
	if v != nil && v.MchId == s.mchId {
		return newBillDomainImpl(v, s._mchRepo, s._rbacRepo, s._registryRepo, s)
	}
	return nil
}

// 计算交易手续费
func (s *transactionManagerImpl) MathTransactionFee(tradeType int, amount int) (int, error) {
	cm := s.mch.ConfManager()
	conf := cm.GetTradeConf(tradeType)
	if conf == nil {
		//todo: 应使用系统默认的比例进行手续费
		return 0, nil
	}
	// 免费
	if conf.Flag&merchant.TFlagFree == merchant.TFlagFree {
		return 0, nil
	}
	switch conf.AmountBasis {
	case enum.AmountBasisNotSet: // 免费
		return 0, nil
	case enum.AmountBasisByAmount: // 按订单单数，收取金额
		return conf.TransactionFee, nil
	case enum.AmountBasisByPercent: // 按订单金额，收取百分比
		return int(float64(amount*conf.TradeRate) / enum.RATE_PERCENT), nil
	default:
		panic("not support amount basis")
	}
}

// GetCurrentDailyBill implements merchant.IMerchantTransactionManager.
func (s *transactionManagerImpl) GetCurrentDailyBill() merchant.IBillDomain {
	unix := time.Now().Unix()
	startTime, endTime := util.GetStartEndUnix(time.Now())
	if s.currentBill != nil && s.currentBill.BillTime == int(startTime) {
		// 如果当前账单存在,则直接返回
		return newBillDomainImpl(s.currentBill, s._mchRepo, s._rbacRepo, s._registryRepo, s)
	}
	bill := s._billRepo.FindBy("mch_id = ? AND bill_type=? AND bill_time = ?", s.mchId, merchant.BillTypeDaily, startTime)
	if bill == nil {
		bill = &merchant.MerchantBill{
			Id:           0,
			MchId:        s.mchId,
			BillType:     merchant.BillTypeDaily,
			BillTime:     int(startTime),
			BillMonth:    time.Unix(startTime, 0).Format("2006-01"),
			StartTime:    int(startTime),
			EndTime:      int(endTime),
			TxCount:      0,
			TxAmount:     0,
			TxFee:        0,
			RefundAmount: 0,
			Status:       int(merchant.BillStatusPending),
			ReviewerId:   0,
			ReviewerName: "",
			ReviewRemark: "",
			ReviewTime:   0,
			BillRemark:   "",
			UserRemark:   "",
			SettleSpCode: "",
			SettleTxNo:   "",
			CreateTime:   int(unix),
			BuildTime:    0,
			UpdateTime:   int(unix),
		}
		s._billRepo.Save(bill)
	}
	s.currentBill = bill
	return newBillDomainImpl(s.currentBill, s._mchRepo, s._rbacRepo, s._registryRepo, s)
}

// TransactionAmounts 交易统计数据
type transactionData struct {
	// Amount 交易金额
	Amount int
	// TransactionFee 交易手续费
	TransactionFee int
	// TransactionCount 交易次数
	TransactionCount int
}

// queryBillAmount 查询账单金额,返回交易金额，手续费，退款金额
func (s *transactionManagerImpl) queryBillAmount(beginTime, endTime int) (ret struct {
	TxAmount         int
	TxFee            int
	RefundAmount     int
	TransactionCount int
}) {
	walletId := s.getWalletId()
	rep := s._walletRepo.LogRepo().Raw()
	// 查询金额函数
	f := func(kind int) *transactionData {
		total := transactionData{}
		rep.Model(wallet.WalletLog{}).Select("count(*) as transaction_count, sum(change_value) as amount, sum(transaction_fee) as transaction_fee").
			Where("wallet_id = ? AND kind = ? AND create_time BETWEEN ? AND ?",
				walletId, kind, beginTime, endTime).
			Find(&total)
		return &total
	}
	total := f(wallet.KCarry)
	refund := f(wallet.KRefund)

	ret.TxAmount = total.Amount - refund.Amount
	ret.TxFee = int(math.Abs(float64(total.TransactionFee + refund.TransactionFee)))
	ret.RefundAmount = refund.Amount
	ret.TransactionCount = total.TransactionCount
	return ret
}

// GenerateMonthlyBill 生成月度账单
func (s *transactionManagerImpl) GenerateMonthlyBill(year, month int) error {
	period := s.GetSettlementPeriod()
	unix := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).Unix()
	startTime, endTime := fw.GetMonthStartEndUnix(unix)
	bill := s._billRepo.FindBy("mch_id = ? AND bill_type = ? AND bill_time = ?",
		s.mchId,
		merchant.BillTypeMonthly,
		startTime)
	txData := s.queryBillAmount(int(startTime), int(endTime))
	if bill == nil {
		bill = &merchant.MerchantBill{
			Id:           0,
			MchId:        s.mchId,
			BillType:     merchant.BillTypeMonthly,
			BillTime:     int(startTime),
			BillMonth:    time.Unix(startTime, 0).Format("2006-01"),
			StartTime:    int(startTime),
			EndTime:      int(endTime),
			Status:       int(merchant.BillStatusWaitConfirm), // 待商户确认
			ReviewerId:   0,
			ReviewerName: "",
			ReviewRemark: "",
			ReviewTime:   0,
			BillRemark:   "",
			UserRemark:   "",
			SettleStatus: merchant.BillWaitSettlement, // 待结算
			SettleSpCode: "",
			SettleTxNo:   "",
			SettleRemark: "",
			CreateTime:   int(unix),
			BuildTime:    int(unix),
			UpdateTime:   int(unix),
		}
	} else {
		if bill.Status == int(merchant.BillStatusReviewed) {
			return errors.New("账单已经复核,无法重新生成")
		}
	}
	// 更新账单数据
	bill.TxCount = txData.TransactionCount
	bill.TxAmount = txData.TxAmount
	bill.TxFee = txData.TxFee
	bill.RefundAmount = txData.RefundAmount
	// 更新状态
	if period == merchant.BillTypeDaily {
		// 日结账单,不需商户确认
		bill.Status = int(merchant.BillStatusReviewed)
		// 月结账单不需结算
		bill.SettleStatus = merchant.BillNoSettlement
	} else {
		// 月结账单,需要商户确认
		bill.Status = int(merchant.BillStatusWaitConfirm)
		// 待结算
		bill.SettleStatus = merchant.BillWaitSettlement
	}
	_, err := s._billRepo.Save(bill)
	return err
}

var _ merchant.IBillDomain = new(billDomainImpl)

// 账单领域
type billDomainImpl struct {
	_value        *merchant.MerchantBill
	_repo         merchant.IMerchantRepo
	_rbacRepo     rbac.IRbacRepository
	_registryRepo registry.IRegistryRepo
	_manager      merchant.IMerchantTransactionManager
}

func newBillDomainImpl(value *merchant.MerchantBill,
	mchRepo merchant.IMerchantRepo,
	rbacRepo rbac.IRbacRepository,
	registryRepo registry.IRegistryRepo,
	manager merchant.IMerchantTransactionManager) merchant.IBillDomain {
	return &billDomainImpl{_value: value, _repo: mchRepo, _rbacRepo: rbacRepo, _registryRepo: registryRepo, _manager: manager}
}

// Value 获取值
func (b *billDomainImpl) Value() *merchant.MerchantBill {
	return types.DeepClone(b._value)
}

func (b *billDomainImpl) UpdateAmount() error {
	if b._value.Status == int(merchant.BillStatusReviewed) {
		return errors.New("账单已复核,无法更新账单金额")
	}
	txData := b._manager.(*transactionManagerImpl).queryBillAmount(b._value.StartTime, b._value.EndTime)
	b._value.TxCount = txData.TransactionCount
	b._value.TxAmount = txData.TxAmount
	b._value.TxFee = txData.TxFee
	b._value.RefundAmount = txData.RefundAmount
	b._value.UpdateTime = int(time.Now().Unix())
	_, err := b._repo.BillRepo().Save(b._value)
	return err
}

// GenerateBill implements merchant.IBillDomain.
func (b *billDomainImpl) Generate() error {
	now := time.Now().Unix()
	if now < int64(b._value.EndTime) {
		return errors.New("账单尚未到截止时间")
	}
	if b._value.TxAmount-b._value.RefundAmount < b._value.TxFee {
		b._value.BillRemark = "账单异常:总交易费大于总交易额"
		_, err := b._repo.BillRepo().Save(b._value)
		return err
	}
	// 更新账单金额
	err := b.UpdateAmount()
	if err != nil {
		return err
	}
	period := b._manager.GetSettlementPeriod()
	if b._value.BillType == merchant.BillTypeDaily {
		// 日结账单,不需商户确认
		b._value.Status = int(merchant.BillStatusWaitReview)
		if period == merchant.BillTypeDaily {
			b._value.SettleStatus = merchant.BillWaitSettlement
		} else {
			// 结算周期为：月结, 不需结算
			b._value.SettleStatus = merchant.BillNoSettlement
		}
	} else {
		// 月结账单
		if period == merchant.BillTypeMonthly {
			// 结算周期为: 需要商户确认
			b._value.Status = int(merchant.BillStatusWaitConfirm)
			b._value.SettleStatus = merchant.BillWaitSettlement
		} else {
			// 结算周期为：日结, 不需要商户确认
			b._value.Status = int(merchant.BillStatusWaitReview)
			b._value.SettleStatus = merchant.BillNoSettlement
		}
	}
	if b._value.TxAmount == 0 && b._value.TxFee == 0 {
		// 未产生金额的账单直接更改为待复核，无需商户确认
		b._value.Status = int(merchant.BillStatusWaitReview)
	}
	b._value.UpdateTime = int(now)
	b._value.BillRemark = ""
	_, err = b._repo.BillRepo().Save(b._value)
	return err
}

// GetDomainId implements merchant.IBillDomain.
func (b *billDomainImpl) GetDomainId() int {
	return b._value.Id
}

// Confirm 确认账单
func (b *billDomainImpl) Confirm() error {
	if b._value.BillType == merchant.BillTypeDaily {
		return errors.New("日结账单无需确认")
	}
	if b._value.Status != int(merchant.BillStatusWaitConfirm) {
		return errors.New("账单尚未生成或已结算")
	}
	b._value.Status = int(merchant.BillStatusWaitReview)
	b._value.UpdateTime = int(time.Now().Unix())
	_, err := b._repo.BillRepo().Save(b._value)
	return err
}

// ReviewBill implements merchant.IBillDomain.
func (b *billDomainImpl) Review(reviewerId int, remark string) error {
	iu := b._rbacRepo.GetRbacAggregateRoot().GetUser(reviewerId)
	if iu == nil {
		return errors.New("复核人不存在")
	}
	if b._value.Status == int(merchant.BillStatusReviewed) {
		return errors.New("账单已经复核")
	}
	if b._value.Status != int(merchant.BillStatusWaitReview) {
		return errors.New("账单非待复核状态")
	}
	b._value.ReviewerId = reviewerId
	b._value.ReviewerName = iu.GetValue().Nickname
	b._value.ReviewTime = int(time.Now().Unix())
	b._value.ReviewRemark = remark
	b._value.Status = int(merchant.BillStatusReviewed)
	if b._value.TxAmount == 0 && b._value.TxFee == 0 {
		// 如果账单金额为零，则自动结算完成
		b._value.SettleStatus = merchant.BillSettlemented
	}
	b._value.UpdateTime = int(time.Now().Unix())
	_, err := b._repo.BillRepo().Save(b._value)
	return err
}

// SettleBill implements merchant.IBillDomain.
func (b *billDomainImpl) Settle() error {
	if b._value.Status != int(merchant.BillStatusReviewed) {
		return errors.New("账单尚未复核或已结算")
	}
	if b._value.SettleStatus == merchant.BillSettlemented {
		return errors.New("账单已经结算")
	}
	if b._value.SettleStatus == merchant.BillNoSettlement {
		return errors.New("账单无需结算")
	}
	mch := b._repo.GetMerchant(b._value.MchId)
	if mch == nil {
		return errors.New("商户不存在")
	}
	conf := mch.ConfManager().GetSettleConf()
	if len(conf.SubMchNo) == 0 {
		return errors.New("商户尚未在支付平台入网,无法完成结算")
	}
	b._value.SettleStatus = merchant.BillSettlemented // 结算成功
	b._value.UpdateTime = int(time.Now().Unix())
	_, err := b._repo.BillRepo().Save(b._value)
	if err == nil {
		mchName := mch.GetValue().MchName
		// 生成结算备注
		platformName, _ := b._registryRepo.GetValue(registry.PlatformName)
		remark := fmt.Sprintf("%s结算账单:%s-%s", platformName, mchName, b._value.BillMonth)
		// 发送结算事件
		go eventbus.Dispatch(&merchant.MerchantBillSettleEvent{
			MchId:         b._value.MchId,
			MchName:       mchName,
			SubMerchantNo: conf.SubMchNo,
			SettleRemark:  remark,
			Bill:          b._value,
		})
	}
	return err
}

// UpdateSettleInfo 更新结算信息
func (b *billDomainImpl) UpdateSettleInfo(spCode string, settleTxNo string, message string) error {
	if b._value.SettleStatus != merchant.BillSettlemented {
		return errors.New("账单尚未结算")
	}
	b._value.SettleSpCode = spCode
	b._value.SettleTxNo = settleTxNo
	b._value.SettleRemark = message
	b._value.UpdateTime = int(time.Now().Unix())
	_, err := b._repo.BillRepo().Save(b._value)
	return err
}

// UpdateBillAmount 更新账单金额
func (b *billDomainImpl) UpdateBillAmount(amount int, txFee int) error {
	if amount == 0 {
		return errors.New("金额不能为零")
	}
	b._value.TxAmount += amount
	b._value.TxFee += txFee
	b._value.TxCount += 1
	if amount < 0 {
		// 退款金额
		b._value.RefundAmount += -amount
	}
	b._value.UpdateTime = int(time.Now().Unix())
	_, err := b._repo.BillRepo().Save(b._value)
	return err
}
