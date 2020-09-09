package service

import (
	"go.etcd.io/etcd/clientv3"
	"go2o/core/etcd"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : discovery.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-08 07:02
 * description :
 * history :
 */

var service = "Go2oService"
var ttl int64 = 10

// 注册服务发现
func initRegistry(cfg *clientv3.Config, port int) {
	r, err := etcd.NewRegistry(service, ttl, *cfg)
	if err != nil {
		panic(err)
	}
	_, err = r.Register(port)
	if err != nil {
		panic(err)
	}
}
