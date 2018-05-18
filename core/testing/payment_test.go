package testing

import (
	"go2o/core/domain/interface/payment"
	"go2o/core/factory"
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
		TradeNo:   tradeNo,
		TradeType: "ppi-charge",
		SellerId:  0,
		OrderType: 0,
		Subject:   "充值",
		BuyerId:   22149,
		//PayUid:      1,
		ItemAmount: 1,
		PaymentFlag: payment.FlagBankCard | payment.FlagOutSp | payment.FlagBalance |
			payment.FlagIntegral | payment.FlagWallet,
		OutTradeNo: "",
		SubmitTime: unix,
		PaidTime:   0,
	})
	if err := ip.Submit(); err != nil {
		t.Error(err)
		t.Failed()
	}
	//err := ip.BalanceDiscount("支付订单")
	//if err != nil{
	//	t.Error(err)
	//}

	//amount, err := ip.IntegralDiscount(100, false)
	//if err != nil {
	//	t.Error(err)
	//}
	//println("---amount=", amount)

	err := ip.PaymentByWallet("")
	if err != nil {
		t.Error(err)
	}

	//ip.TradeNoPrefix("CZ")
	//ip.PaymentFinish("alipay", "1234567890")
	t.Log("订单号：", ip.TradeNo(), "; 订单状态:", ip.State())
}

// 测试支付单交易完成
func TestPaymentOrderTradeFinish(t *testing.T) {
	tradeNo := "IC6180515221155668"
	ic := factory.Repo.GetPaymentRepo().GetPaymentOrder(tradeNo)
	if ic == nil {
		t.Errorf("支付单:%s不存在", tradeNo)
		t.Failed()
	}
	err := ic.TradeFinish()
	if err != nil {
		t.Errorf("支付单%s支付失败：%s", tradeNo, err.Error())
		t.Failed()
	}
}
