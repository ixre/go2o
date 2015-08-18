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
	"github.com/jrsix/gof"
	"github.com/jrsix/gof/db"
	"github.com/jrsix/gof/log"
	"github.com/jrsix/gof/storage"
	"github.com/garyburd/redigo/redis"
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
	_storage      gof.Storage
}

func NewMainApp(confPath string) *MainApp {
	return &MainApp{
		_confFilePath: confPath,
	}
}

func (this *MainApp) Db() db.Connector {
	if this._dbConnector == nil {
		this._dbConnector = getDb(this.Config(), this._debugMode, this.Log())
	}
	return this._dbConnector
}

func (this *MainApp) Storage() gof.Storage {
	if this._storage == nil {
		this._storage = storage.NewRedisStorage(this.Redis())
	}
	return this._storage
}

func (this *MainApp) Template() *gof.Template {
	if this._template == nil {
		this._template = initTemplate(this.Config())
	}
	return this._template
}

func (this *MainApp) Config() *gof.Config {
	if this._config == nil {
		if cfg, err := gof.LoadConfig(this._confFilePath); err == nil {
			this._config = cfg
			cfg.Set("exp_fee_bit", float64(1))
		} else {
			log.Fatalln(err)
		}
	}
	return this._config
}

func (this *MainApp) Source() interface{} {
	return this
}

func (this *MainApp) Debug() bool {
	return this._debugMode
}

func (this *MainApp) Log() log.ILogger {
	if this._logger == nil {
		var flag int = 0
		if this._debugMode {
			flag = log.LOpen | log.LESource | log.LStdFlags
		}
		this._logger = log.NewLogger(nil, " O2O", flag)
	}
	return this._logger
}

func (this *MainApp) Redis() *redis.Pool {
	if this._redis == nil {
		this._redis = createRedisPool(this.Config())
	}
	return this._redis
}

func (this *MainApp) Init(debug, trace bool) bool {
	this._debugMode = debug

	if trace {
		this.Db().GetOrm().SetTrace(this._debugMode)
	}
	this.Loaded = true
	return true
}
