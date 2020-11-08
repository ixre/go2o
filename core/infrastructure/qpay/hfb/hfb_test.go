package hfb

import (
	"encoding/hex"
	"github.com/ixre/gof/crypto"
	"go.etcd.io/etcd/clientv3"
	"go2o/core/infrastructure"
	"go2o/core/infrastructure/qpay"
	"strconv"
	"strings"
	"testing"
	"time"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : hfb_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-05 20:01
 * description :
 * history :
 */

var h qpay.QuickPayProvider

func init() {
	// 默认的ETCD端点
	etcdEndPoints := []string{"http://127.0.0.1:2379"}
	cfg := clientv3.Config{
		Endpoints:   etcdEndPoints,
		DialTimeout: 5 * time.Second,
	}
	s, _ := infrastructure.NewEtcdStorage(cfg)
	h = NewHfb(s)
}

func TestCardBin(t *testing.T) {
	bankCardNo := "6227000010990006191"
	r := h.QueryCardBin(bankCardNo)
	t.Logf("%#v", r)
}

func TestHfbImpl_RequestBankSideAuth(t *testing.T) {
	bankCardNo := "6227000010990006191"
	nonce := strconv.Itoa(int(time.Now().Unix()))
	r, err := h.RequestBankSideAuth(nonce, bankCardNo, "闫雪龙",
		"22011219850823101X", "13810512111")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(r.AuthForm)
}
func TestHfbImpl_DirectPayment(t *testing.T) {
	bankCardNo := "6227000010990006191"
	rsp, err := h.QueryBankAuth(bankCardNo)
	if err == nil {
		t.Logf("%#v", rsp)
	}
	if rsp.Code != 0{
		t.Log("卡片未授权成功")
		t.FailNow()
	}
	orderNo := "BZ"+time.Now().Format("20060102150405")
	ret, err := h.DirectPayment(orderNo, 1, "补差价链接", rsp.BankAuthToken,
		"127.0.0.1", "http://www.go2o-dev.56x.net/qpay/notify_url",
		"")
	if err != nil{
		t.Error(err)
		t.FailNow()
	}
	t.Log("支付成功,订单号：",orderNo,"第三方订单号：",ret.BillNo)
}



func TestHfbImpl_QueryPaymentStatus(t *testing.T) {
	orderNo := "bz20201107224102"
	billTime := "20201107224102"
	ret,err := h.QueryPaymentStatus(orderNo,map[string]string{"agent_bill_time":billTime})
	if err != nil{
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%#v",ret)
}

func TestHfbImpl_QueryBankAuth(t *testing.T) {
	//bankCardNo := "6227000010990006191"
	bankCardNo := "9559981014352796313"

	r, err := h.QueryBankAuth(bankCardNo)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(r)
}

func TestHfbImpl_BatchTransfer(t *testing.T) {
	orderNo := "BZ"+time.Now().Format("20060102150405")
	nonce := strconv.Itoa(int(time.Now().Unix()))
	batch := []*qpay.CardTransferReq{
		{
			OrderNo:         orderNo,
			PersonTransfer:  true,
			TradeFee:        1, //0.01
			BankCardNo:      "6227000010990006191",
			BankCode:        "2",
			BankAccountName: "闫雪龙",
			Subject:         "",
			Province:        "-",
			City:            "-",
			StoreName:       "-",
		},
	}
	r,err := h.BatchTransfer(orderNo,batch,nonce,"http://www.go2o-dev.56x.net/qpay/callback")
	if err != nil{
		t.Error(err)
		t.FailNow()
	}else{
		t.Logf("%#v",r)
	}
}

func TestHfbImpl_Encrypt3DES(t *testing.T) {
	detail_data := "agent_id=2126129&batch_amt=0.01&batch_no=bz20201108080915&batch_num=1&detail_data=bz20201108080915^2^0^6227000010990006191^闫雪龙^0.01^上游结算款^^^&ext_param1=1604794155&key=0e05aac3be0746269f114bd7&notify_url=http://www.go2o-dev.56x.net/qpay/callback&version=3";
	key := []byte("4865534446254C0F8837DFB3")
	bytes ,_ := crypto.EncryptECB3DES([]byte(detail_data),key)
	s := hex.EncodeToString(bytes)
	t.Log(strings.ToUpper(s))
}


