package domain

import (
	"testing"
	"time"

	"github.com/ixre/go2o/tests/ti"
)

func TestPaymentSuccess(t *testing.T) {
	orderNo := "1220607000313450"
	repo := ti.Factory.GetPaymentRepo()
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
