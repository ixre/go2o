/**
 * Copyright 2015 @ z3q.net.
 * name : testing
 * author : jarryliu
 * date : 2016-06-15 08:31
 * description :
 * history :
 */
package include

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jsix/gof"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/log"
	"github.com/jsix/gof/storage"
	"go2o/core"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/repository"
)

var (
	app gof.App
)

func GetApp() gof.App {
	if app == nil {

		app = new(testingApp)
		app.Config().Set("redis_host", "172.16.69.128")
		app.Config().Set("redis_db", "10")
		app.Config().Set("redis_port", "6379")
		app.Config().Set("redis_auth", "123456")
		gof.CurrentApp = app
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
}

func newMainApp(confPath string) *testingApp {
	return &testingApp{
		_confFilePath: confPath,
	}
}

func (t *testingApp) Db() db.Connector {
	if t._dbConnector == nil {
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=Local",
			"root",
			"",
			"172.16.69.128",
			"3306",
			"txmall",
			"utf8",
		)
		connector := db.NewSimpleConnector("mysql", connStr, t.Log(), 5000, false)
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

func (t *testingApp) Template() *gof.Template {
	return t._template
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

var (
	ProMRepo promodel.IProModelRepo
)

func init() {
	app := GetApp()
	db := app.Db()
	orm := db.GetOrm()
	sto := app.Storage()
	ProMRepo = repository.NewProModelRepo(db, orm)

	goodsRepo := repository.NewGoodsRepo(db)
	valRepo := repository.NewValueRepo(db, sto)
	userRepo := repository.NewUserRepo(db)
	notifyRepo := repository.NewNotifyRepo(db)
	mssRepo := repository.NewMssRepo(db, notifyRepo, valRepo)
	expressRepo := repository.NewExpressRepo(db, valRepo)
	shipRepo := repository.NewShipmentRepo(db, expressRepo)
	memberRepo := repository.NewMemberRepo(app.Storage(), db, mssRepo, valRepo)
	itemRepo := repository.NewItemRepo(db)
	tagSaleRepo := repository.NewTagSaleRepo(db)
	promRepo := repository.NewPromotionRepo(db, goodsRepo, memberRepo)
	cateRepo := repository.NewCategoryRepo(db, valRepo, sto)
	saleRepo := repository.NewSaleRepo(db, cateRepo, valRepo, tagSaleRepo,
		itemRepo, expressRepo, goodsRepo, promRepo)
	cartRepo := repository.NewCartRepo(db, memberRepo, goodsRepo)
	shopRepo := repository.NewShopRepo(db, sto)
	mchRepo := repository.NewMerchantRepo(db, sto, shopRepo, userRepo, memberRepo, mssRepo, valRepo)
	//personFinanceRepo := repository.NewPersonFinanceRepository(db, memberRepo)
	deliveryRepo := repository.NewDeliverRepo(db)
	//contentRepo := repository.NewContentRepo(db)
	//adRepo := repository.NewAdvertisementRepo(db)
	spRepo := repository.NewOrderRepo(app.Storage(), db, mchRepo, nil, saleRepo, cartRepo, goodsRepo,
		promRepo, memberRepo, deliveryRepo, expressRepo, shipRepo, valRepo)
	payRepo := repository.NewPaymentRepo(app.Storage(), db, memberRepo, spRepo, valRepo)

	goodsRepo.SetSaleRepo(saleRepo) //fixed
	spRepo.SetPaymentRepo(payRepo)
}
