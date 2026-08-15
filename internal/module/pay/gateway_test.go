package pay

import (
	"github.com/ixre/gof/storage"
	"strconv"
	"testing"
	"time"
)

// 测试提交支付请求到网关
func TestGateway_Submit(t *testing.T) {
	pool := storage.NewRedisPool("127.0.0.1",
		6379, 10, "123456", 10000, 200)
	st := storage.NewRedisStorage(pool)
	gw := NewGateway(st)
	userId := int64(1)
	token := gw.CreatePostToken(userId)
	tradeNo := strconv.Itoa(int(time.Now().UnixNano()))
	data := map[string]string{
		"token":         token,
		"trade_no":      tradeNo,
		"subject":       "线下支付",
		"amount":        "100",
		"procedure_fee": "0",
		"vendor_id":     "1",
		"shop_id":       "16",
		"notify_url":    "http://m.ts.com/trade/epay_notify",
	}
	err := gw.Submit(userId, data)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = gw.CheckAndPayment(userId, tradeNo, "189405")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
