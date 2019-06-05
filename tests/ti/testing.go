/**
 * Copyright 2015 @ to2.net.
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
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/storage"
	"go2o/core"
	"go2o/core/msq"
	"go2o/core/repos"
	"time"
)

var (
	app     *testingApp
	Factory *repos.RepoFactory
)
var (
	REDIS_DB = "1"
)

func GetApp() gof.App {
	return gof.CurrentApp
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
		t._registry, _ = gof.NewRegistry("../conf", ":")
	}
	return t._registry
}

func (t *testingApp) Db() db.Connector {
	if t._dbConnector == nil {
		connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			t._config.GetString(core.DbUsr),
			t._config.GetString(core.DbPwd),
			t._config.GetString(core.DbServer),
			t._config.GetString(core.DbPort),
			t._config.GetString(core.DbName))
		conn := db.NewConnector("postgresql", connStr, t.Log(), t._debugMode)
		conn.SetMaxIdleConns(10000)
		conn.SetMaxIdleConns(5000)
		conn.SetConnMaxLifetime(time.Second * 10)
		t._dbConnector = conn
		orm.CacheProxy(t._dbConnector.GetOrm(), t.Storage())
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
	app := core.NewApp("../app_dev.conf")
	gof.CurrentApp = app
	core.Init(app, false, false)
	conn := app.Db()
	sto := app.Storage()
	Factory = (&repos.RepoFactory{}).Init(conn, sto, "../conf")
}

// 初始化producer
func InitMsq() {
	msq.Configure(msq.KAFKA, []string{"127.0.0.1:9092"})
}
