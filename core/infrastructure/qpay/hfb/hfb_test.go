package hfb

import (
	"go.etcd.io/etcd/clientv3"
	"go2o/core/infrastructure"
	"go2o/core/infrastructure/qpay"
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
func init(){
	// 默认的ETCD端点
	etcdEndPoints := []string{"http://127.0.0.1:2379"}
	cfg := clientv3.Config{
		Endpoints:   etcdEndPoints,
		DialTimeout: 5 * time.Second,
	}
	s,_:= infrastructure.NewEtcdStorage(cfg)
	h = NewHfb(s)
}

func TestCardBin(t *testing.T){
	bankCardNo := "6227000010990006191"
	r:= h.QueryCardBin(bankCardNo)
	t.Logf("%#v",r)
}
