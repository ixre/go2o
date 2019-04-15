/**
 * Copyright 2015 @ z3q.net.
 * name : testing
 * author : jarryliu
 * date : 2016-06-15 08:31
 * description :
 * history :
 */
package ti

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/storage"
	"go2o/core"
	"go2o/core/factory"
)

var (
	app     *testingApp
	Factory *factory.RepoFactory
)
var (
	REDIS_DB = "1"
)

func GetApp() gof.App {
	if app == nil {
		app = new(testingApp)
		app.Config().Set("conf_path", "../../conf")
		app.Config().Set("redis_host", "127.0.0.1")
		app.Config().Set("redis_db", REDIS_DB)
		app.Config().Set("redis_port", "6379")
		//app.Config().Set("redis_auth", "123456")
		app.Config().Set("redis_auth", "")
		app.Config().Set("db_name", "gcy_v3")
		gof.CurrentApp = app
		app.Init(true, true)
	}
	return app
}

var _ gof.App = new(testingApp)

// application context
// implement of web.Application
type testingApp struct {
	Loaded        bool
	_confFilePath string
	_config       *gof.Config
	_redis        *redis.Pool
	_dbConnector  db.Connector
	_debugMode    bool
	_template     *gof.Template
	_logger       log.ILogger
	_storage      storage.Interface
	_registry     *gof.Registry
}

func newMainApp(confPath string) *testingApp {
	return &testingApp{
		_confFilePath: confPath,
	}
}

func (t *testingApp) Registry() *gof.Registry {
	if t._registry == nil {
		t._registry, _ = gof.NewRegistry("../../conf", ":")
	}
	return t._registry
}

func (t *testingApp) Db() db.Connector {
	if t._dbConnector == nil {
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=Local",
			"root",
			"123456",
			"127.0.0.1",
			"3306",
			t.Config().GetString("db_name"),
			"utf8",
		)
		connector := db.NewConnector("mysql", connStr, t.Log(), false)
		core.OrmMapping(connector)
		t._dbConnector = connector
	}
	return t._dbConnector
}

func (t *testingApp) Storage() storage.Interface {
	if t._storage == nil {
		t._storage = storage.NewRedisStorage(t.Redis())
	}
	return t._storage
}

func (t *testingApp) Config() *gof.Config {
	if t._config == nil {
		if t._confFilePath == "" {
			t._config = gof.NewConfig()
		} else {
			if cfg, err := gof.LoadConfig(t._confFilePath); err == nil {
				t._config = cfg
			} else {
				log.Fatalln(err)
			}
		}
	}
	return t._config
}

func (t *testingApp) Source() interface{} {
	return t
}

func (t *testingApp) Debug() bool {
	return t._debugMode
}

func (t *testingApp) Log() log.ILogger {
	if t._logger == nil {
		var flag int = 0
		if t._debugMode {
			flag = log.LOpen | log.LESource | log.LStdFlags
		}
		t._logger = log.NewLogger(nil, " O2O", flag)
	}
	return t._logger
}

func (t *testingApp) Redis() *redis.Pool {
	if t._redis == nil {
		t._redis = core.CreateRedisPool(t.Config())
	}
	return t._redis
}

func (t *testingApp) Init(debug, trace bool) bool {
	t._debugMode = debug

	if trace {
		t.Db().GetOrm().SetTrace(t._debugMode)
	}
	t.Loaded = true
	return true
}

func init() {
	app := GetApp()
	conn := app.Db()
	sto := app.Storage()
	confPath := app.Config().GetString("conf_path")
	Factory = (&factory.RepoFactory{}).Init(conn, sto, confPath)
}
