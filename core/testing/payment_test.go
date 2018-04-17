package testing

import (
	"go2o/core/domain/interface/payment"
	"go2o/core/infrastructure/domain"
	"go2o/core/testing/ti"
	"testing"
	"time"
)

// 测试充值订单
func TestCreateChargePaymentOrder(t *testing.T) {
	repo := ti.Factory.GetPaymentRepo()
	unix := time.Now().Unix()
	ip := repo.CreatePaymentOrder(&payment.PaymentOrder{
		TradeNo:          domain.NewTradeNo(0),
		TradeType:        "ppi-charge",
		VendorId:         0,
		Type:             0,
		OrderId:          0,
		Subject:          "充值",
		BuyUser:          1,
		PaymentUser:      1,
		TotalAmount:      0.01,
		BalanceDiscount:  0,
		IntegralDiscount: 0,
		SystemDiscount:   0,
		CouponDiscount:   0,
		SubAmount:        0,
		AdjustmentAmount: 0,
		FinalFee:         0.01,
		PayFlag:          payment.SignOnlinePay | payment.SignWalletAccount,
		PaymentSign:      0,
		OuterNo:          "",
		CreateTime:       unix,
		PaidTime:         0,
		State:            0,
	})
	if _, err := ip.Commit(); err != nil {
		t.Error(err)
		t.Failed()
	}
	ip.TradeNoPrefix("CZ")
	t.Log("订单号：", ip.GetTradeNo())
}
