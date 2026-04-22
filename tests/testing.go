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

	"github.com/ixre/go2o/pkg/event/msq"
	"github.com/ixre/go2o/pkg/infrastructure/domain"
	"github.com/ixre/go2o/pkg/infrastructure/etcd"
	"github.com/ixre/go2o/pkg/initial"
	"github.com/ixre/go2o/pkg/initial/bootstrap"
	"github.com/ixre/go2o/pkg/initial/provide"
	"github.com/ixre/go2o/pkg/inject"
	"github.com/ixre/go2o/pkg/interface/domain/registry"
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
	confPath := "local.conf"
	// 默认的ETCD端点
	natsEndPoints := []string{"localhost:4222"}
	etcdEndPoints := []string{"localhost:2379"}

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
	app := bootstrap.NewApp(confPath, &cfg)
	initial.Init(app, false, false)
	// 初始化第三方配置
	inject.GetSPConfig().Configure()
	// 初始化nats
	msq.Configure(msq.NATS, natsEndPoints, inject.GetRegistryRepo())
	// 初始化分布式锁
	etcd.InitializeLocker(&cfg)
	// 初始化事件
	inject.GetEventSource().Bind()
	// 初始化私钥
	initPrivateKey()
}

func initPrivateKey() {
	repo := inject.GetRegistryRepo()
	key, _ := repo.GetValue(registry.SysPrivateKey)
	domain.ConfigureHmacPrivateKey(key)
}
