package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/typeconv"
)

func TestGetPaymentOrder(t *testing.T) {
	ret, err := inject.GetPaymentService().GetPaymentOrder(
		context.TODO(),
		&proto.PaymentOrderRequest{
			TradeNo:    "1230227000283186",
			AllowBreak: false,
		})
	if err != nil {
		t.Error(err)
	}
	t.Log(typeconv.MustJson(ret.TradeData))
}

// 测试查询可用于分账的订单
func TestQueryDivideOrders(t *testing.T) {
	memberId := 848
	ret, err := inject.GetPaymentService().QueryDivideOrders(
		context.TODO(),
		&proto.DivideOrdersRequest{
			MemberId:  int64(memberId),
			OrderType: payment.TypeRecharge,
		})
	if err != nil {
		t.Error(err)
	}
	for _, v := range ret.Orders {
		t.Logf("orderNo: %s, amount:%.2f,data:%s", v.TradeNo, float64(v.Amount/100), typeconv.MustJson(v))
	}
}

// 测试退款
func TestReturnRecharge(t *testing.T) {
	ps := inject.GetPaymentService()
	pr, _ := ps.RequestRefund(context.TODO(), &proto.PaymentRefundRequest{
		TradeNo:      "2240909848196227",
		RefundAmount: 850,
		Reason:       "测试退款",
	})
	if pr.Code > 0 {
		t.Error(pr.Message)
		t.FailNow()
	}
}

// 测试退款全部可退金额
func TestRefundAvailPaymentOrder(t *testing.T) {
	ps := inject.GetPaymentService()
	ret, _ := ps.RequestRefundAvail(context.TODO(), &proto.PaymentRefundAvailRequest{
		TradeNo: "2240909848196227",
		Remark:  "测试退款",
	})
	if ret.Code > 0 {
		t.Error(ret.Message)
		t.FailNow()
	}
	t.Logf("ret: %v, 退款金额:%.2f", ret, float32(ret.Amount)/100)
}
