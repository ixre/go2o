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
	"go2o/core/repository"
)

var (
	app gof.App
)

func GetApp() gof.App {
	if app == nil {

		app = new(testingApp)
		app.Config().Set("redis_host", "rds.flm-dev.redvp.com")
		app.Config().Set("redis_db", "2")
		app.Config().Set("redis_port", "6379")
		app.Config().Set("redis_auth", "")
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
			"123456",
			"dbs.flm-dev.redvp.com",
			"3306",
			"flm",
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
				cfg.Set("exp_fee_bit", float64(1))
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
	db := app.Db()
	sto := app.Storage()
	goodsRep := repository.NewGoodsRep(db)
	valRep := repository.NewValueRep(db, sto)
	userRep := repository.NewUserRep(db)
	notifyRep := repository.NewNotifyRep(db)
	mssRep := repository.NewMssRep(db, notifyRep, valRep)
	expressRep := repository.NewExpressRep(db, valRep)
	shipRep := repository.NewShipmentRep(db, expressRep)
	memberRep := repository.NewMemberRep(app.Storage(), db, mssRep, valRep)
	itemRep := repository.NewItemRep(db)
	tagSaleRep := repository.NewTagSaleRep(db)
	promRep := repository.NewPromotionRep(db, goodsRep, memberRep)
	cateRep := repository.NewCategoryRep(db, valRep, sto)
	saleRep := repository.NewSaleRep(db, cateRep, valRep, tagSaleRep,
		itemRep, expressRep, goodsRep, promRep)
	cartRep := repository.NewCartRep(db, memberRep, goodsRep)
	shopRep := repository.NewShopRep(db, sto)
	mchRep := repository.NewMerchantRep(db, sto, shopRep, userRep, memberRep, mssRep, valRep)
	//personFinanceRep := repository.NewPersonFinanceRepository(db, memberRep)
	deliveryRep := repository.NewDeliverRep(db)
	//contentRep := repository.NewContentRep(db)
	//adRep := repository.NewAdvertisementRep(db)
	spRep := repository.NewOrderRep(app.Storage(), db, mchRep, nil, saleRep, cartRep, goodsRep,
		promRep, memberRep, deliveryRep, expressRep, shipRep, valRep)
	payRep := repository.NewPaymentRep(app.Storage(), db, memberRep, spRep, valRep)

	goodsRep.SetSaleRep(saleRep) //fixed
	spRep.SetPaymentRep(payRep)
}
