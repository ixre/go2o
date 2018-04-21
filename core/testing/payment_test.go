package testing

import (
	"go2o/core/domain/interface/payment"
	"go2o/core/infrastructure/domain"
	"go2o/core/testing/ti"
	"testing"
	"time"
)

func TestCreateTradeNo(t *testing.T) {
	for i := 0; i < 10; i++ {
		println(domain.NewTradeNo(0, i))
	}
}

// 测试充值订单
func TestCreateChargePaymentOrder(t *testing.T) {
	repo := ti.Factory.GetPaymentRepo()
	unix := time.Now().Unix()
	tradeNo := domain.NewTradeNo(0, 0)
	ip := repo.CreatePaymentOrder(&payment.PaymentOrder{
		TradeNo:          tradeNo,
		TradeType:        "ppi-charge",
		VendorId:         0,
		Type:             0,
		OrderId:          0,
		Subject:          "充值",
		BuyUser:          22149,
		PaymentUser:      1,
		TotalAmount:      0.01,
		BalanceDiscount:  0,
		IntegralDiscount: 0,
		SystemDiscount:   0,
		CouponDiscount:   0,
		SubAmount:        0,
		AdjustmentAmount: 0,
		FinalFee:         1.01,
		PayFlag:          0,
		PaymentSign:      payment.SignOnlinePay | payment.SignWalletAccount,
		OuterNo:          "",
		CreateTime:       unix,
		PaidTime:         0,
		State:            0,
	})
	if err := ip.Commit(); err != nil {
		t.Error(err)
		t.Failed()
	}
	//ip.TradeNoPrefix("CZ")
	ip.PaymentFinish("alipay", "1234567890")
	t.Log("订单号：", ip.GetTradeNo())
}
