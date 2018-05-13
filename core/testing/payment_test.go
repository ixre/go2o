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
	ip := repo.CreatePaymentOrder(&payment.Order{
		TradeNo:     tradeNo,
		TradeType:   "ppi-charge",
		SellerId:    0,
		OrderType:   0,
		OrderId:     0,
		Subject:     "充值",
		BuyerId:     22149,
		PayUid:      1,
		TotalAmount: 1,
		//FinalFee:    1,
		PayFlag:    payment.PBankCard | payment.POutSP,
		OutTradeNo: "",
		SubmitTime: unix,
		PaidTime:   0,
		State:      0,
	})
	if err := ip.Submit(); err != nil {
		t.Error(err)
		t.Failed()
	}
	//ip.TradeNoPrefix("CZ")
	ip.PaymentFinish("alipay", "1234567890")
	t.Log("订单号：", ip.TradeNo())
}
