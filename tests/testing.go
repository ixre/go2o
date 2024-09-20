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
	"github.com/ixre/go2o/core/inject"
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
	confPath := "app-test.conf"
	// 默认的ETCD端点
	natsEndPoints := []string{"http://go2o.dev:4222"}
	etcdEndPoints := []string{"http://go2o.dev:2379"}

	// 加载配置文件
	for {
		_, err := os.Stat(confPath)
		if err == nil {
			break
		}
		confPath = "../" + confPath
	}
	cfgFile, err := gof.LoadConfig(confPath)
	if err != nil {
		panic(err)
	}
	if v := cfgFile.GetString("go2o_nats_addr"); len(v) > 0 {
		natsEndPoints = []string{v}
	}
	if v := cfgFile.GetString("go2o_etcd_addr"); len(v) > 0 {
		etcdEndPoints = []string{v}
	}
	// 连接nats和etcd
	cfg := clientv3.Config{
		Endpoints:   etcdEndPoints,
		DialTimeout: 5 * time.Second,
	}
	app := initial.NewApp(confPath, &cfg)
	initial.Init(app, false, false)
	repos.OrmMapping(provide.GetOrmInstance())
	// 初始化第三方配置
	inject.GetSPConfig().Configure()
	// 初始化nats
	msq.Configure(msq.NATS, natsEndPoints, inject.GetRegistryRepo())
	// 初始化分布式锁
	etcd.InitializeLocker(&cfg)
	// 初始化事件
	inject.GetEventSource().Bind()
}
