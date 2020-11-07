package hfb

import (
	"go.etcd.io/etcd/clientv3"
	"go2o/core/infrastructure"
	"go2o/core/infrastructure/qpay"
	"strconv"
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

func TestHfbImpl_QueryBankAuthByNonceId(t *testing.T) {
	bankCardNo := "6227000010990006191"
	nonce := strconv.Itoa(int(time.Now().Unix()))
	r, err := h.RequestBankSideAuth(nonce, bankCardNo, "闫雪龙",
		"22011219850823101X", "13810512111")
	if err == nil {
		//t.Log(r.AuthForm)
		var rsp *qpay.BankAuthQueryResponse
		rsp, err = h.QueryBankAuthByNonceId(r.NonceId)
		if err == nil {
			t.Logf("%#v", rsp)
		}
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
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
