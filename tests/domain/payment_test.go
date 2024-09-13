/**
 * Copyright (C) 2007-2024 fze.NET,All rights reserved.
 *
 * name : payment_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2015-09-08 10:02
 * description : 支付测试
 * history :
 */
package domain

import (
	"math"
	"testing"
	"time"

	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/inject"
)

// 测试支付完成
func TestPaymentSuccess(t *testing.T) {
	orderNo := "2240907702773876"
	repo := inject.GetPaymentRepo()
	ip := repo.GetPaymentOrder(orderNo)
	if ip == nil {
		t.Error("no such order")
		t.FailNow()
	}
	err := ip.PaymentFinish("test", "123456789")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	time.Sleep(time.Second * 2)
}

// 测试钱包支付单钱包抵扣
func TestWalletDeductPaymentOrder(t *testing.T) {
	orderNo := "1230306000808485"
	repo := inject.GetPaymentRepo()
	ip := repo.GetPaymentOrder(orderNo)
	if ip == nil {
		t.Error("no such order")
		t.FailNow()
	}
	err := ip.WalletDeduct("抵扣")
	if err != nil {
		t.Log("----", err.Error())
		t.FailNow()
	}
}

func TestCreateTradeNo(t *testing.T) {
	for i := 0; i < 10; i++ {
		println(domain.NewTradeNo(0, i))
	}
}

// 测试充值订单
func TestCreateChargePaymentOrder(t *testing.T) {
	repo := inject.GetPaymentRepo()
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
		TotalAmount: 1,
		PayFlag: domain.MathPaymentMethodFlag([]int{
			payment.MBankCard, payment.MPaySP, payment.MBalance,
			payment.MIntegral, payment.MWallet}),
		OutTradeNo: "",
		SubmitTime: int(unix),
		PaidTime:   0,
	})
	if err := ip.Submit(); err != nil {
		t.Error(err)
		t.Failed()
	}
	//err := ip.BalanceDeduct("支付订单")
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
	ic := inject.GetPaymentRepo().GetPaymentOrder(tradeNo)
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

// 测试拆分支付单均摊抵扣金额
func TestBreakPaymentOrderAVGDeductAmount(t *testing.T) {
	finalAmount := 10039
	deductAmount := 38
	finalAmount1 := 3021
	finalAmount2 := 7018

	avgAmount1 := int(math.Round(float64(deductAmount) * (float64(finalAmount1) / float64(finalAmount))))
	avgAmount2 := int(math.Round(float64(deductAmount) * (float64(finalAmount2) / float64(finalAmount))))
	t.Log(avgAmount1, avgAmount2)

	if avgAmount1+avgAmount2 != deductAmount {
		t.FailNow()
	}
}

func TestCancelPaymentOrder(t *testing.T) {
	orderNo := "1230326007400338"
	p := inject.GetPaymentRepo().GetPaymentOrder(orderNo)
	if err := p.Cancel(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
