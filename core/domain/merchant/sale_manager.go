/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: sale_manager.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2017-01-27 12:23:22
 * description: 商户交易管理
 * history:
 */
package merchant

import (
	"errors"
	"fmt"
	"time"

	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/invoice"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	rbac "github.com/ixre/go2o/core/domain/interface/rabc"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
	"github.com/ixre/gof/domain/eventbus"
)

var _ merchant.IMerchantTransactionManager = new(TransactionManagerImpl)

type TransactionManagerImpl struct {
	mch           *merchantImpl
	mchId         int
	_mchRepo      merchant.IMerchantRepo
	_invoiceRepo  invoice.IInvoiceRepo
	_billRepo     fw.Repository[merchant.MerchantBill]
	currentBill   *merchant.MerchantBill
	_rbacRepo     rbac.IRbacRepository
	_registryRepo registry.IRegistryRepo
}

// AdjustBillAmount implements merchant.IMerchantTransactionManager.
func (s *TransactionManagerImpl) AdjustBillAmount(amountType merchant.BillAmountType, amount int, txFee int) error {
	bill := s.GetCurrentBill()
	switch amountType {
	case merchant.BillAmountTypeShop:
		bill.ShopOrderCount += types.Ternary(amount > 0, 1, -1)
		bill.ShopTotalAmount += amount
	case merchant.BillAmountTypeStore:
		bill.StoreOrderCount += types.Ternary(amount > 0, 1, -1)
		bill.StoreTotalAmount += amount
	case merchant.BillAmountTypeOther:
		bill.OtherOrderCount += types.Ternary(amount > 0, 1, -1)
		bill.OtherTotalAmount += amount
	}
	bill.TotalTxFee += txFee
	bill.UpdateTime = int(time.Now().Unix())
	_, err := s._billRepo.Save(bill)
	return err
}

// GenerateBill implements merchant.IMerchantTransactionManager.
func (s *TransactionManagerImpl) GenerateBill() error {
	now := time.Now().Unix()
	bill := s.GetCurrentBill()
	if now < int64(bill.EndTime) {
		return errors.New("账单尚未到截止时间")
	}
	if bill.OtherTotalAmount+bill.ShopTotalAmount+bill.StoreTotalAmount < bill.TotalTxFee {
		return errors.New("账单异常:总交易费大于总交易额")
	}
	bill.Status = int(merchant.BillStatusGenerated)
	bill.BuildTime = int(now)
	bill.UpdateTime = int(now)
	_, err := s._billRepo.Save(bill)
	return err
}

// GetBillByTime implements merchant.IMerchantTransactionManager.
func (s *TransactionManagerImpl) GetBillByTime(billTime int) *merchant.MerchantBill {
	month, _ := fw.GetMonthStartEndUnix(int64(billTime))
	bill := s._billRepo.FindBy("mch_id = ? AND bill_time = ?", s.mchId, month)
	return bill
}

// GetCurrentBill implements merchant.IMerchantTransactionManager.
func (s *TransactionManagerImpl) GetCurrentBill() *merchant.MerchantBill {
	unix := time.Now().Unix()
	month, end := fw.GetMonthStartEndUnix(unix)
	if s.currentBill != nil && s.currentBill.BillTime == int(month) {
		// 如果当前账单存在,则直接返回
		return s.currentBill
	}
	bill := s._billRepo.FindBy("mch_id = ? AND bill_time = ?", s.mchId, month)
	if bill == nil {
		bill = &merchant.MerchantBill{
			Id:               0,
			MchId:            s.mchId,
			BillTime:         int(month),
			BillMonth:        time.Unix(month, 0).Format("2006-01"),
			StartTime:        int(month),
			EndTime:          int(end),
			ShopOrderCount:   0,
			StoreOrderCount:  0,
			ShopTotalAmount:  0,
			StoreTotalAmount: 0,
			OtherOrderCount:  0,
			OtherTotalAmount: 0,
			Status:           int(merchant.BillStatusPending),
			ReviewerId:       0,
			ReviewerName:     "",
			ReviewTime:       0,
			CreateTime:       int(unix),
			BuildTime:        0,
			UpdateTime:       int(unix),
		}
		s._billRepo.Save(bill)
		s.currentBill = bill
	}
	return bill
}

// ReviewBill implements merchant.IMerchantTransactionManager.
func (s *TransactionManagerImpl) ReviewBill(billId int, reviewerId int) error {
	bill := s._billRepo.FindBy("mch_id = ? AND id = ?", s.mchId, billId)
	if bill == nil {
		return errors.New("账单不存在")
	}
	iu := s._rbacRepo.GetRbacAggregateRoot().GetUser(reviewerId)
	if iu == nil {
		return errors.New("复核人不存在")
	}
	if bill.Status != int(merchant.BillStatusGenerated) {
		return errors.New("账单尚未生成或已结算")
	}
	bill.ReviewerId = reviewerId
	bill.ReviewerName = iu.GetValue().Nickname
	bill.ReviewTime = int(time.Now().Unix())
	bill.Status = int(merchant.BillStatusReviewed)
	bill.UpdateTime = int(time.Now().Unix())
	s._billRepo.Save(bill)
	return nil
}

// SettleBill implements merchant.IMerchantTransactionManager.
func (s *TransactionManagerImpl) SettleBill(billId int) error {
	bill := s._billRepo.FindBy("mch_id = ? AND id = ?", s.mchId, billId)
	if bill == nil {
		return errors.New("账单不存在")
	}
	if bill.Status != int(merchant.BillStatusReviewed) {
		return errors.New("账单尚未复核或已结算")
	}
	conf := s.mch.ConfManager().GetSettleConf()
	if len(conf.SubMchNo) == 0 {
		return errors.New("商户尚未在支付平台入网,无法完成结算")
	}
	bill.Status = int(merchant.BillStatusSettled)
	bill.UpdateTime = int(time.Now().Unix())
	_, err := s._billRepo.Save(bill)
	if err == nil {
		// 生成结算备注
		platformName, _ := s._registryRepo.GetValue(registry.PlatformName)
		remark := fmt.Sprintf("%s结算账单%d", platformName, bill.BillMonth)
		// 发送结算事件
		go eventbus.Publish(&merchant.MerchantBillSettleEvent{
			MchId:         s.mchId,
			MchName:       s.mch.GetValue().MchName,
			SubMerchantNo: conf.SubMchNo,
			SettleRemark:  remark,
			Bill:          bill,
		})
	}
	return err
}

func newSaleManagerImpl(id int, m *merchantImpl, mchRepo merchant.IMerchantRepo,
	invoiceRepo invoice.IInvoiceRepo,
	rbacRepo rbac.IRbacRepository,
	registryRepo registry.IRegistryRepo) merchant.IMerchantTransactionManager {
	return &TransactionManagerImpl{
		mchId:         id,
		mch:           m,
		_mchRepo:      mchRepo,
		_billRepo:     m._repo.BillRepo(),
		_invoiceRepo:  invoiceRepo,
		_rbacRepo:     rbacRepo,
		_registryRepo: registryRepo,
	}
}

// 计算交易手续费
func (s *TransactionManagerImpl) MathTransactionFee(tradeType int, amount int) (int, error) {
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
