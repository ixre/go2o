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
	"strconv"
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

// CreateBill 创建账单
func (s *transactionManagerImpl) CreateBill(billType int, unixtime int) merchant.IBillDomain {
	unix := time.Now().Unix()
	bill := &merchant.MerchantBill{}
	bill.BillType = billType
	bill.MchId = s.mchId
	bill.CreateTime = int(unix)
	bill.UpdateTime = int(unix)
	bill.Status = int(merchant.BillStatusPending)
	bill.BillMonth = time.Unix(int64(unixtime), 0).Format("2006-01")
	if billType == merchant.BillTypeDaily {
		// 日度账单
		startTime, endTime := util.GetStartEndUnix(time.Unix(int64(unixtime), 0))
		bill.StartTime = int(startTime)
		bill.EndTime = int(endTime)
		bill.BillTime = int(startTime)
	} else if billType == merchant.BillTypeMonthly {
		// 月度账单
		startTime, endTime := fw.GetMonthStartEndUnix(int64(unixtime))
		bill.StartTime = int(startTime)
		bill.EndTime = int(endTime)
		bill.BillTime = int(startTime)
	} else {
		panic("not support bill type")
	}
	return newBillDomainImpl(bill, s.mch, s._mchRepo, s._rbacRepo, s._registryRepo, s)
}

// GetDailyBill implements merchant.IMerchantTransactionManager.
func (s *transactionManagerImpl) GetDailyBill(billTime int) merchant.IBillDomain {
	startTime, _ := util.GetStartEndUnix(time.Unix(int64(billTime), 0))
	bill := s._billRepo.FindBy("bill_type = ? AND mch_id = ? AND bill_time = ? ",
		merchant.BillTypeDaily,
		s.mchId,
		startTime)
	if bill != nil {
		return newBillDomainImpl(bill,
			s.mch,
			s._mchRepo,
			s._rbacRepo,
			s._registryRepo,
			s,
		)
	}
	return nil
}

// GetBill 获取指定账单
func (s *transactionManagerImpl) GetBill(billId int) merchant.IBillDomain {
	v := s._billRepo.Get(billId)
	if v != nil && v.MchId == s.mchId {
		return newBillDomainImpl(v,
			s.mch,
			s._mchRepo,
			s._rbacRepo,
			s._registryRepo,
			s,
		)
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
	startTime, _ := util.GetStartEndUnix(time.Now())
	if s.currentBill != nil && s.currentBill.BillTime == int(startTime) {
		// 如果当前账单存在,则直接返回
		return newBillDomainImpl(
			s.currentBill,
			s.mch,
			s._mchRepo,
			s._rbacRepo,
			s._registryRepo,
			s,
		)
	}
	bill := s._billRepo.FindBy("mch_id = ? AND bill_type=? AND bill_time = ?", s.mchId, merchant.BillTypeDaily, startTime)
	if bill == nil {
		bill = s.CreateBill(merchant.BillTypeDaily, int(startTime)).Value()
		s._billRepo.Save(bill)
	}
	s.currentBill = bill
	return newBillDomainImpl(s.currentBill, s.mch, s._mchRepo, s._rbacRepo, s._registryRepo, s)
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

var _ merchant.IBillDomain = new(billDomainImpl)

// 账单领域
type billDomainImpl struct {
	_mch          merchant.IMerchantAggregateRoot
	_value        *merchant.MerchantBill
	_repo         merchant.IMerchantRepo
	_rbacRepo     rbac.IRbacRepository
	_registryRepo registry.IRegistryRepo
	_manager      merchant.IMerchantTransactionManager
}

func newBillDomainImpl(value *merchant.MerchantBill,
	mch merchant.IMerchantAggregateRoot,
	mchRepo merchant.IMerchantRepo,
	rbacRepo rbac.IRbacRepository,
	registryRepo registry.IRegistryRepo,
	manager merchant.IMerchantTransactionManager) merchant.IBillDomain {
	return &billDomainImpl{
		_mch:          mch,
		_value:        value,
		_repo:         mchRepo,
		_rbacRepo:     rbacRepo,
		_registryRepo: registryRepo,
		_manager:      manager,
	}
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
	b._value.FinalAmount = b._value.TxAmount - b._value.TxFee - b._value.RefundAmount
	b._value.UpdateTime = int(time.Now().Unix())
	return b.save()
}

// GenerateBill implements merchant.IBillDomain.
func (b *billDomainImpl) Generate() error {
	period := b._manager.GetSettlementPeriod()
	if b._value.BillType == merchant.BillTypeMonthly {
		// 月度账单
		if period == merchant.BillTypeDaily {
			return errors.New("结算周期为日结,无需生成月度账单")
		}
		return b.generateMonthlyBill()
	}
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
	if b._value.BillType == merchant.BillTypeDaily {
		// 日结账单,不需商户确认
		b._value.Status = int(merchant.BillStatusWaitReview)
		if period == merchant.BillTypeDaily {
			b._value.SettleStatus = merchant.BillWaitSettlement
			// 结算结果
			b._value.SettleResult = merchant.BillResultPending
		} else {
			// 结算周期为：月结, 不需结算
			b._value.SettleStatus = merchant.BillNoSettlement
			// 结算结果
			b._value.SettleResult = merchant.BillResultNone
		}
	} else {
		// 月结账单
		if period == merchant.BillTypeMonthly {
			// 结算周期为: 需要商户确认
			b._value.Status = int(merchant.BillStatusWaitConfirm)
			b._value.SettleStatus = merchant.BillWaitSettlement
			// 结算结果
			b._value.SettleResult = merchant.BillResultPending
		} else {
			// 结算周期为：日结, 不需要商户确认
			b._value.Status = int(merchant.BillStatusWaitReview)
			b._value.SettleStatus = merchant.BillNoSettlement
			// 结算结果
			b._value.SettleResult = merchant.BillResultNone
		}
	}
	if b._value.TxAmount == 0 && b._value.TxFee == 0 {
		// 未产生金额的账单直接更改为待复核，无需商户确认
		b._value.Status = int(merchant.BillStatusReviewed)
		b._value.SettleStatus = merchant.BillNoSettlement
		b._value.SettleResult = merchant.BillResultNone
	}
	b._value.UpdateTime = int(now)
	b._value.BillRemark = ""
	b._value.BuildTime = int(now)
	return b.save()
}

// GenerateMonthlyBill 生成月度账单
func (s *billDomainImpl) generateMonthlyBill() error {
	unix := time.Now().Unix()
	mchId := s._value.MchId
	startTime, endTime := fw.GetMonthStartEndUnix(int64(s._value.BillTime))
	bill := s._repo.BillRepo().FindBy("mch_id = ? AND bill_type = ? AND bill_time = ?",
		mchId,
		merchant.BillTypeMonthly,
		startTime)
	if endTime > unix {
		return errors.New("未到出账时间")
	}
	year, month := time.Unix(int64(s._value.BillTime), 0).Year(), time.Unix(int64(s._value.BillTime), 0).Month()
	fmt.Printf("生成月度账单, 商户ID:%d, 账单日期:%d-%d\n", mchId, year, month)
	txData := s._manager.(*transactionManagerImpl).queryBillAmount(int(startTime), int(endTime))
	if bill == nil {
		// 如果账单不存在,则使用当前创建的账单
		bill = s._value
	}
	if bill.Status == int(merchant.BillStatusReviewed) {
		return errors.New("账单已经复核,无法重新生成")
	}

	// 更新账单数据
	bill.TxCount = txData.TransactionCount
	bill.TxAmount = txData.TxAmount
	bill.TxFee = txData.TxFee
	bill.RefundAmount = txData.RefundAmount
	bill.BuildTime = int(unix)
	// 更新状态
	if bill.TxAmount > 0 || bill.TxFee > 0 {
		// 月结账单,需要商户确认
		bill.Status = int(merchant.BillStatusWaitConfirm)
		bill.SettleStatus = merchant.BillWaitSettlement
		bill.SettleResult = merchant.BillResultPending
	} else {
		// 无交易金额，则不结算
		bill.Status = int(merchant.BillStatusReviewed)
		bill.SettleStatus = merchant.BillNoSettlement
		bill.SettleResult = merchant.BillResultNone
	}
	_, err := s._repo.BillRepo().Save(bill)
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
	return b.save()
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
		b._value.SettleResult = merchant.BillResultSuccess
	}
	b._value.UpdateTime = int(time.Now().Unix())
	return b.save()
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
	// 扣除商户余额
	err := b._mch.Account().Consume("账单结算-"+b._value.BillMonth,
		b._value.FinalAmount, strconv.Itoa(b._value.Id), "系统扣除")
	if err != nil {
		return errors.New("结算扣款失败:" + err.Error())
	}
	b._value.SettleStatus = merchant.BillSettlemented // 结算成功
	b._value.UpdateTime = int(time.Now().Unix())
	err = b.save()
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
func (b *billDomainImpl) UpdateSettleInfo(spCode string, settleTxNo string, settleResult int, message string) error {
	if b._value.SettleStatus != merchant.BillSettlemented {
		return errors.New("账单尚未结算")
	}
	if b._value.SettleStatus == merchant.BillNoSettlement {
		return errors.New("账单无需结算")
	}
	// if b._value.SettleResult == merchant.BillResultSuccess {
	// 	return errors.New("结算结果为成功,无需更新")
	// }
	b._value.SettleSpCode = spCode
	b._value.SettleTxNo = settleTxNo
	b._value.SettleRemark = message
	b._value.SettleResult = settleResult
	b._value.UpdateTime = int(time.Now().Unix())
	return b.save()
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
	return b.save()
}

func (b *billDomainImpl) save() error {
	b.updateAmount()
	_, err := b._repo.BillRepo().Save(b._value)
	return err
}

func (b *billDomainImpl) updateAmount() {
	b._value.FinalAmount = b._value.TxAmount - b._value.TxFee - b._value.RefundAmount
}
