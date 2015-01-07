package glob

import (
	"github.com/garyburd/redigo/redis"
	"github.com/newmin/gof"
	"github.com/newmin/gof/app"
	"github.com/newmin/gof/db"
	"github.com/newmin/gof/log"
	"github.com/newmin/gof/web"
)

var (
	activeContext *AppContext
	_             app.Context = new(AppContext)
)

//application context
// impment of web.Application
type AppContext struct {
	Loaded      bool
	config      *gof.Config
	Redis       *redis.Pool
	dbConnector db.Connector
	debugMode   bool
	template    *web.TemplateWrapper
	logger      log.ILogger
}

func (this *AppContext) Db() db.Connector {
	return this.dbConnector
}

func (this *AppContext) Template() *web.TemplateWrapper {
	return this.template
}

func (this *AppContext) Config() *gof.Config {
	return this.config
}

func (this *AppContext) Source() interface{} {
	return this
}

func (this *AppContext) Debug() bool {
	return this.debugMode
}

func (this *AppContext) Log() log.ILogger {
	if this.logger == nil {
		var flag int = 0
		if this.debugMode {
			flag = log.LOpen | log.LESource | log.LStdFlags
		}
		this.logger = log.NewLogger(nil, " O2O", flag)
	}
	return this.logger
}

func (this *AppContext) Init(debug, trace bool) {
	this.debugMode = debug
	cfg := this.config
	activeContext.Redis = createRedisPool(cfg)
	//todo: check redis connected
	activeContext.template = initTemplate(cfg)
	activeContext.dbConnector = getDb(cfg, this.Log())

	if trace {
		activeContext.Db().GetOrm().SetTrace(this.debugMode)
	}

	this.Loaded = true
}

//create a new context of application
func NewContext() *AppContext {
	if activeContext == nil {
		cfg, err := gof.NewConfig("conf/boot.conf")
		cfg.Set("exp_fee_bit", float64(1.5))

		if err != nil {
			log.Fatalln("[Error]:", err.Error())
		}
		activeContext = &AppContext{
			config: cfg,
		}
	}
	return activeContext
}

//当前上下文对象
func CurrContext() *AppContext {
	if activeContext == nil {
		//activeContext = (app.Context).Source().(*AppContext)
		activeContext = NewContext()
	}
	return activeContext
}
