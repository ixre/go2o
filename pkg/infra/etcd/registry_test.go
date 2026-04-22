package etcd

import (
	"go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : registry_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-02 16:44
 * description :
 * history :
 */
var service = "Go2oService"
var ttl int64 = 3
var cfg = clientv3.Config{
	Endpoints:   []string{"http://localhost:2379/"},
	DialTimeout: 5 * time.Second,
	Username:    "",
	Password:    "",
}

func TestRegisterService(t *testing.T) {
	arr := make([]Registry, 0)
	r, _ := NewRegistry(service, ttl, cfg)
	for i := 0; i < 3; i++ {
		_, _ = r.Register("", 10+i)
		arr = append(arr, r)
	}
	time.Sleep(15 * time.Second)
	r.Stop()
}

func TestSelector(t *testing.T) {
	selector, _ := NewSelector(service, cfg, AlgRoundRobin)
	for {
		next, err := selector.Next()
		if err != nil {
			t.Logf("select node error %s", err.Error())
		} else {
			t.Log("selected:" + next.Addr)
		}
		time.Sleep(time.Second * 3)
	}
}
