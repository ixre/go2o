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
	"github.com/garyburd/redigo/redis"
	"github.com/jsix/gof"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/log"
	"github.com/jsix/gof/storage"
	"go2o/core"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/repository"
)

var (
	app *testingApp
)

func GetApp() gof.App {
	if app == nil {
		app = new(testingApp)
		app.Config().Set("redis_host", "172.16.69.128")
		app.Config().Set("redis_db", "10")
		app.Config().Set("redis_port", "6379")
		app.Config().Set("redis_auth", "123456")
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
	ProMRepo       promodel.IProModelRepo
	AfterSalesRepo afterSales.IAfterSalesRepo
	OrderRepo      order.IOrderRepo
	ExpressRepo    express.IExpressRepo
	ValueRepo      valueobject.IValueRepo
	ItemRepo       item.IGoodsItemRepo
	ProductRepo    product.IProductRepo
	CatRepo        product.ICategoryRepo
	MemberRepo     member.IMemberRepo
	MchRepo        merchant.IMerchantRepo
	CartRepo       cart.ICartRepo
	ShipmentRepo   shipment.IShipmentRepo
)

func init() {
	app := GetApp()
	db := app.Db()
	orm := db.GetOrm()
	sto := app.Storage()
	proMRepo := repository.NewProModelRepo(db, orm)
	valueRepo := repository.NewValueRepo(db, sto)
	userRepo := repository.NewUserRepo(db)
	notifyRepo := repository.NewNotifyRepo(db)
	mssRepo := repository.NewMssRepo(db, notifyRepo, valueRepo)
	expressRepo := repository.NewExpressRepo(db, valueRepo)
	shipRepo := repository.NewShipmentRepo(db, expressRepo)
	memberRepo := repository.NewMemberRepo(sto, db, mssRepo, valueRepo)
	productRepo := repository.NewProductRepo(db, proMRepo, valueRepo)
	itemWsRepo := repository.NewItemWholesaleRepo(db)
	itemRepo := repository.NewGoodsItemRepo(db, productRepo,
		proMRepo, itemWsRepo, expressRepo, valueRepo)
	//tagSaleRepo := repository.NewTagSaleRepo(db, valRepo)
	promRepo := repository.NewPromotionRepo(db, itemRepo, memberRepo)
	catRepo := repository.NewCategoryRepo(db, valueRepo, sto)
	//afterSalesRepo := repository.NewAfterSalesRepo(db)
	cartRepo := repository.NewCartRepo(db, memberRepo, itemRepo)
	shopRepo := repository.NewShopRepo(db, sto)
	wholesaleRepo := repository.NewWholesaleRepo(db)
	mchRepo := repository.NewMerchantRepo(db, sto, wholesaleRepo, shopRepo, userRepo, memberRepo, mssRepo, valueRepo)
	//personFinanceRepo := repository.NewPersonFinanceRepository(db, memberRepo)
	deliveryRepo := repository.NewDeliverRepo(db)
	//contentRepo := repository.NewContentRepo(db)
	//adRepo := repository.NewAdvertisementRepo(db, sto)
	orderRepo := repository.NewOrderRepo(sto, db, mchRepo, nil, productRepo, cartRepo, itemRepo,
		promRepo, memberRepo, deliveryRepo, expressRepo, shipRepo, valueRepo)
	paymentRepo := repository.NewPaymentRepo(sto, db, memberRepo, orderRepo, valueRepo)
	asRepo := repository.NewAfterSalesRepo(db, orderRepo, memberRepo, paymentRepo)

	orderRepo.SetPaymentRepo(paymentRepo)

	ProMRepo = proMRepo
	AfterSalesRepo = asRepo
	OrderRepo = orderRepo
	ExpressRepo = expressRepo
	ValueRepo = valueRepo
	ItemRepo = itemRepo
	ProductRepo = productRepo
	CatRepo = catRepo
	MemberRepo = memberRepo
	MchRepo = mchRepo
	CartRepo = cartRepo
	ShipmentRepo = shipRepo
}
