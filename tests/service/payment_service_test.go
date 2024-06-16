package service

import (
	"context"
	"testing"

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
