package etcd

import (
	"go.etcd.io/etcd/clientv3"
	"strconv"
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

func TestRegisterService(t *testing.T) {
	service := "Go2oService"
	var ttl int64 = 3
	cfg := clientv3.Config{
		Endpoints:   []string{"http://localhost:2379/"},
		DialTimeout: 5 * time.Second,
	}
	leaseIDList := make([]int64,0)
	for i:= 0;i<3;i++ {
		r, _ := NewRegistry(service, ttl, cfg)
		id,_ := r.Register("127.0.0."+strconv.Itoa(i)+":1428")
		leaseIDList = append(leaseIDList,id)
	}
	time.Sleep(15 * time.Second)
	for _,v := range leaseIDList{
		r, _ := NewRegistry(service, ttl, cfg)
		r.Revoke(v)
	}
}
