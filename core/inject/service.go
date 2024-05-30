//go:build wireinject

package inject

import (
	"github.com/google/wire"
	"github.com/ixre/go2o/core/event"
	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
)

var serviceProvideSets = wire.NewSet(provideSets,
	impl.NewStatusService,
	impl.NewRegistryService,
	impl.NewPromotionService,
	impl.NewMerchantService,
	impl.NewRegistryService,
	impl.NewPromotionService,
	impl.NewFoundationService,
	impl.NewMemberService,
	impl.NewMerchantService,
	impl.NewShopService,
	impl.NewProductService,
	//impl.NewItemService,
	//impl.NewOrderService,
	impl.NewCartService,
	impl.NewAfterSalesService,
	// 事件
	event.NewEventSource,
	event.NewEventHandler
)

// 状态服务
func GetStatusService() proto.StatusServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 注册表服务
func GetRegistryService() proto.RegistryServiceServer {
	panic(wire.Build(serviceProvideSets))
}

func GetPromService() impl.PromotionService {
	panic(wire.Build(serviceProvideSets))
}

// 基础服务
func GetFoundationService() proto.FoundationServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 会员服务
func GetMemberService() proto.MemberServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 商户服务
func GetMerchantService() proto.MerchantServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 商店服务
func GetShopService() proto.ShopServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 产品服务
func GetProductService() proto.ProductServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 商品服务
func GetItemService() proto.ItemServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 购物服务
func GetOrderService() proto.OrderServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 购物车服务
func GetCartService() proto.CartServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 售后服务
func GetAfterSalesService() proto.AfterSalesServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 支付服务
func GetPaymentService() proto.PaymentServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 快捷支付服务
func GetQuickPayService() proto.QuickPayServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 消息服务
func GetMessageService() proto.MessageServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 快递服务
func GetExpressService() proto.ExpressServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 配送服务
func GetShipmentService() proto.ShipmentServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 内容服务
func GetContentService() proto.ContentServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 广告服务
func GetAdService() proto.AdvertisementServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 钱包服务
func GetWalletService() proto.WalletServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// // 个人金融服务
// func GetPersonFinanceService() personFinanceService {
// 	panic(wire.Build(serviceProvideSets))
// }

// // 门户数据服务
// func GetPortalService() portalService {
// 	panic(wire.Build(serviceProvideSets))
// }

// // 查询服务
// func GetQueryService() proto.QueryServiceServer {
// 	panic(wire.Build(serviceProvideSets))
// }

// // ExecuteService 执行任务服务
// func GetExecuteService() executionServiceImpl {
// 	panic(wire.Build(serviceProvideSets))
// }

// func GetCommonDao() impl.CommonDao {
// 	panic(wire.Build(serviceProvideSets))
// }

// // AppService APP服务
// func GetAppService() appServiceImpl {
// 	panic(wire.Build(serviceProvideSets))
// }

// // RbacService 权限服务
// func GetRbacService() rbacServiceImpl {
// 	panic(wire.Build(serviceProvideSets))
// }

// // CodeService 条码服务
// func GetCodeService() codeServiceImpl {
// 	panic(wire.Build(serviceProvideSets))
// }
