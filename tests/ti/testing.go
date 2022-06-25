/**
 * Copyright 2015 @ 56x.net.
 * name : testing
 * author : jarryliu
 * date : 2016-06-15 08:31
 * description :
 * history :
 */
package ti

import (
	"github.com/ixre/go2o/core"
	"github.com/ixre/go2o/core/repos"
	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"go.etcd.io/etcd/clientv3"
	"os"
	"time"
)

var (
	Factory *repos.RepoFactory
	_conn   db.Connector
	_orm    orm.Orm
)
var (
	REDIS_DB = "1"
)

func GetApp() gof.App {
	return gof.CurrentApp
}

func GetOrm() orm.Orm {
	return _orm
}

func GetConnector() db.Connector {
	return _conn
}

func init() {
	// 默认的ETCD端点
	etcdEndPoints := []string{"http://127.0.0.1:2379"}
	cfg := clientv3.Config{
		Endpoints:   etcdEndPoints,
		DialTimeout: 5 * time.Second,
	}
	confPath := "app.conf"
	for {
		_, err := os.Stat(confPath)
		if err == nil {
			break
		}
		confPath = "../" + confPath
	}
	app := core.NewApp(confPath, &cfg)
	gof.CurrentApp = app
	core.Init(app, false, false)
	_conn = app.Db()
	sto := app.Storage()
	_orm = orm.NewOrm(_conn.Driver(), _conn.Raw())
	impl.InitTestService(app, _conn, _orm, sto)
	Factory = (&repos.RepoFactory{}).Init(_orm, sto)
}
