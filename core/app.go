/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package core

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jsix/gof"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"github.com/jsix/gof/log"
	"github.com/jsix/gof/storage"
	"go2o/core/variable"
)

var _ gof.App = new(MainApp)

// application context
// implement of web.Application
type MainApp struct {
	Loaded        bool
	_confFilePath string
	_config       *gof.Config
	_redis        *redis.Pool
	_dbConnector  db.Connector
	_debugMode    bool
	_template     *gof.Template
	_logger       log.ILogger
	_storage      storage.Interface
}

func NewMainApp(confPath string) *MainApp {
	return &MainApp{
		_confFilePath: confPath,
	}
}

func (a *MainApp) Db() db.Connector {
	if a._dbConnector == nil {
		a._dbConnector = getDb(a.Config(), a._debugMode, a.Log())
		orm.CacheProxy(a._dbConnector.GetOrm(), a.Storage())
	}
	return a._dbConnector
}

func (a *MainApp) Storage() storage.Interface {
	if a._storage == nil {
		a._storage = storage.NewRedisStorage(a.Redis())
	}
	return a._storage
}

func (a *MainApp) Template() *gof.Template {
	if a._template == nil {
		a._template = initTemplate(a.Config())
	}
	return a._template
}

func (a *MainApp) Config() *gof.Config {
	if a._config == nil {
		if cfg, err := gof.LoadConfig(a._confFilePath); err == nil {
			a._config = cfg
			variable.Domain = a._config.GetString(variable.ServerDomain)
			cfg.Set("exp_fee_bit", float64(1))
		} else {
			log.Fatalln(err)
		}
	}
	return a._config
}

func (a *MainApp) Source() interface{} {
	return a
}

func (a *MainApp) Debug() bool {
	return a._debugMode
}

func (a *MainApp) Log() log.ILogger {
	if a._logger == nil {
		var flag int = 0
		if a._debugMode {
			flag = log.LOpen | log.LESource | log.LStdFlags
		}
		a._logger = log.NewLogger(nil, " O2O", flag)
	}
	return a._logger
}

func (a *MainApp) Redis() *redis.Pool {
	if a._redis == nil {
		a._redis = CreateRedisPool(a.Config())
	}
	return a._redis
}

func (a *MainApp) Init(debug, trace bool) bool {
	a._debugMode = debug

	if trace {
		a.Db().GetOrm().SetTrace(a._debugMode)
	}
	a.Loaded = true
	return true
}
