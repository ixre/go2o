//go:build wireinject

package inject

import (
	"github.com/google/wire"
	impl2 "github.com/ixre/go2o/internal/core/dao/impl"
	"github.com/ixre/go2o/internal/core/query"
	"github.com/ixre/go2o/internal/core/repos"
	spd "github.com/ixre/go2o/internal/core/service/impl"
	"github.com/ixre/go2o/internal/external"
	"github.com/ixre/go2o/pkg/event"
	"github.com/ixre/go2o/pkg/event/handler"
	"github.com/ixre/go2o/pkg/initial/provide"
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
	repos.NewArticleRepo,
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
	repos.NewApprovalRepository,
	repos.NewPageRepo,
	repos.NewArticleCategoryRepo,
	repos.NewInvoiceTenantRepo,
	repos.NewChatRepo,
	repos.NewWorkorderRepo,
	repos.NewRbacRepo,
	repos.NewSysAppRepo,
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
	query.NewWorkQuery,
	query.NewWalletQuery,
	query.NewInvoiceQuery,
	query.NewAdvertisementQuery,
	query.NewStatisticsQuery,
	query.NewPaymentQuery,
	query.NewSystemQuery,
)

var daoProvideSets = wire.NewSet(
	queryProvideSets,
	impl2.NewCommDao,
	impl2.NewPortalDao,
)

var InjectProvideSets = wire.NewSet(
	daoProvideSets,
	spd.NewStatusService,
	spd.NewRegistryService,
	spd.NewMerchantService,
	spd.NewPromotionService,
	spd.NewSystemService,
	spd.NewMemberService,
	spd.NewShopService,
	spd.NewProductService,
	spd.NewItemService,
	spd.NewShoppingService,
	spd.NewCartService,
	spd.NewAfterSalesService,
	spd.NewAdvertisementService,
	spd.NewPaymentService,
	spd.NewQuickPayService,
	spd.NewMessageService,
	spd.NewExpressService,
	spd.NewShipmentService,
	spd.NewContentService,
	spd.NewWalletService,
	spd.NewCodeService,
	spd.NewQueryService,
	spd.NewRbacService,
	spd.NewAppService,
	spd.NewPortalService,
	spd.NewPersonFinanceService,
	spd.NewExecutionService,
	spd.NewCheckService,
	spd.NewInvoiceService,
	spd.NewChatService,
	spd.NewWorkorderService,
	spd.NewApprovalService,
	spd.NewServiceProviderService,
	// 事件
	event.NewEventSource,
	handler.NewEventHandler,
	handler.NewPaymentEventHandler,
	handler.NewMerchantEventHandler,
	handler.NewInvoiceEventHandler,
	// 其他
	external.NewSPConfig,
)

func GetStationQueryService() *query.StationQuery {
	panic(wire.Build(queryProvideSets))
}

// GetMerchantQueryService 商户查询服务
func GetMerchantQueryService() *query.MerchantQuery {
	panic(wire.Build(queryProvideSets))
}

// GetMemberQueryService 会员查询服务
func GetMemberQueryService() *query.MemberQuery {
	panic(wire.Build(queryProvideSets))
}

// GetContentQuery 获取内容查询服务
func GetContentQueryService() *query.ContentQuery {
	panic(wire.Build(queryProvideSets))
}

func GetWorkQueryService() *query.WorkQuery {
	panic(wire.Build(queryProvideSets))
}

// GetWalletQueryService 获取钱包查询服务
func GetWalletQueryService() *query.WalletQuery {
	panic(wire.Build(queryProvideSets))
}

// GetInvoiceQueryService 获取发票查询服务
func GetInvoiceQueryService() *query.InvoiceQuery {
	panic(wire.Build(queryProvideSets))
}

// GetOrderQueryService 获取广告查询服务
func GetAdvertisementQueryService() *query.AdvertisementQuery {
	panic(wire.Build(queryProvideSets))
}

// GetStatisticsQueryService 获取统计查询服务
func GetStatisticsQueryService() *query.StatisticsQuery {
	panic(wire.Build(queryProvideSets))
}

// GetPaymentQueryService 获取支付查询服务
func GetPaymentQueryService() *query.PaymentQuery {
	panic(wire.Build(queryProvideSets))
}

// GetSystemQueryService 获取系统查询服务
func GetSystemQueryService() *query.SystemQuery {
	panic(wire.Build(queryProvideSets))
}
