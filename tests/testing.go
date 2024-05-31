/**
 * Copyright 2015 @ 56x.net.
 * name : testing
 * author : jarryliu
 * date : 2016-06-15 08:31
 * description :
 * history :
 */
package tests

import (
	"os"
	"time"

	"github.com/ixre/go2o/core/etcd"
	"github.com/ixre/go2o/core/event/msq"
	"github.com/ixre/go2o/core/initial"
	"github.com/ixre/go2o/core/initial/provide"
	"github.com/ixre/go2o/core/repos"
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func GetApp() gof.App {
	return provide.GetApp()
}

func GetOrm() orm.Orm {
	return provide.GetOrmInstance()
}

func GetConnector() db.Connector {
	return provide.GetDb()
}

func init() {
	// 默认的ETCD端点
	//etcdEndPoints := []string{"http://127.0.0.1:2379"}
	natsEndPoints := []string{"http://go2o.dev:4222"}
	etcdEndPoints := []string{"http://go2o.dev:2379"}
	cfg := clientv3.Config{
		Endpoints:   etcdEndPoints,
		DialTimeout: 5 * time.Second,
	}
	msq.Configure(msq.NATS, natsEndPoints)
	confPath := "app-test.conf"
	for {
		_, err := os.Stat(confPath)
		if err == nil {
			break
		}
		confPath = "../" + confPath
	}
	app := initial.NewApp(confPath, &cfg)
	initial.Init(app, false, false)

	// 初始化分布式锁
	etcd.InitializeLocker(&cfg)
	repos.OrmMapping(provide.GetOrmInstance())
}
