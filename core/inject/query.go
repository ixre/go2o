//go:build wireinject

package inject

import (
	"github.com/google/wire"
	impl2 "github.com/ixre/go2o/core/dao/impl"
	"github.com/ixre/go2o/core/event"
	"github.com/ixre/go2o/core/event/handler"
	"github.com/ixre/go2o/core/initial/provide"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/core/repos"
	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/sp"
)

var provideSets = wire.NewSet(
	provide.GetOrm,
	provide.GetGOrm,
	provide.GetOrmInstance,
	provide.GetStorageInstance,
	provide.GetApp,
	provide.GetDb,
	repos.NewSystemRepo,
	repos.NewRegistryRepo,
	repos.NewProModelRepo,
	repos.NewValueRepo,
	repos.NewUserRepo,
	repos.NewWalletRepo,
	repos.NewNotifyRepo,
	repos.NewMssRepo,
	repos.NewExpressRepo,
	repos.NewShipmentRepo,
	repos.NewMemberRepo,
	repos.NewProductRepo,
	repos.NewItemWholesaleRepo,
	repos.NewCategoryRepo,
	repos.NewShopRepo,
	repos.NewGoodsItemRepo,
	repos.NewAfterSalesRepo,
	repos.NewCartRepo,
	repos.NewContentRepo,
	repos.NewMerchantRepo,
	repos.NewOrderRepo,
	repos.NewPaymentRepo,
	repos.NewPromotionRepo,
	repos.NewStationRepo,
	repos.NewTagSaleRepo,
	repos.NewWholesaleRepo,
	repos.NewPersonFinanceRepository,
	repos.NewDeliverRepo,
	repos.NewAdvertisementRepo,
	repos.NewJobRepository,
	repos.NewStaffRepo,
)
var queryProvideSets = wire.NewSet(
	provideSets,
	query.NewStationQuery,
	query.NewMerchantQuery,
	query.NewOrderQuery,
	query.NewMemberQuery,
	query.NewShopQuery,
	query.NewItemQuery,
	query.NewAfterSalesQuery,
	query.NewContentQuery,
)

var daoProvideSets = wire.NewSet(
	queryProvideSets,
	impl2.NewCommDao,
	impl2.NewPortalDao,
)

var InjectProvideSets = wire.NewSet(
	daoProvideSets,
	impl.NewStatusService,
	impl.NewRegistryService,
	impl.NewMerchantService,
	impl.NewPromotionService,
	impl.NewFoundationService,
	impl.NewMemberService,
	impl.NewShopService,
	impl.NewProductService,
	impl.NewItemService,
	impl.NewShoppingService,
	impl.NewCartService,
	impl.NewAfterSalesService,
	impl.NewAdvertisementService,
	impl.NewPaymentService,
	impl.NewQuickPayService,
	impl.NewMessageService,
	impl.NewExpressService,
	impl.NewShipmentService,
	impl.NewContentService,
	impl.NewWalletService,
	impl.NewCodeService,
	impl.NewQueryService,
	impl.NewRbacService,
	impl.NewAppService,
	impl.NewPortalService,
	impl.NewPersonFinanceService,
	impl.NewExecutionService,
	impl.NewCheckService,
	// 事件
	event.NewEventSource,
	handler.NewEventHandler,
	// 其他
	sp.NewSPConfig,
)

func GetStationQueryService() *query.StationQuery {
	panic(wire.Build(queryProvideSets))
}

func GetMerchantQueryService() *query.MerchantQuery {
	panic(wire.Build(queryProvideSets))
}

func GetSPConfig() *sp.ServiceProviderConfiguration {
	panic(wire.Build(InjectProvideSets))
}
