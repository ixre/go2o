/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : discovery.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-08 07:02
 * description :
 * history :
 */
package service

import (
	"fmt"
	"net"

	"github.com/ixre/go2o/core/etcd"
	"github.com/ixre/gof/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var service = "Go2oService"
var ttl int64 = 10

// RegisterServiceDiscovery 注册服务发现
func RegisterServiceDiscovery(cfg *clientv3.Config, host string, port int) {
	r, err := etcd.NewRegistry(service, ttl, *cfg)
	if err != nil {
		panic(err)
	}
	ip := host
	if len(ip) == 0 {
		ip = resolveIp()
	}
	_, err = r.Register(ip, port)
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("[ GO2O][ INFO]: service registration discovery successfully. node: %s:%d", ip, port))
}

func resolveIp() string {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalln(err.Error())
	}
	for _, address := range addrList {
		// 检查ip地址判断是否回环地址
		if i, ok := address.(*net.IPNet); ok && !i.IP.IsLoopback() {
			if i.IP.To4() != nil {
				return i.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
