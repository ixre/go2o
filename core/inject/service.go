//go:build wireinject

package inject

import (
	"github.com/google/wire"
	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
)

var serviceProvideSets = wire.NewSet(provideSets,
	impl.NewStatusService,
	impl.NewRegistryService,
	// NewRegistryService,
	// NewPromotionService,
	// NewFoundationService,
	// NewMemberService,
	// NewMerchantService,
	// NewShopService,
	// NewProductService,
	// NewItemService,
	// NewOrderService,
	// NewCartService,
	// NewAfterSalesService,
)

// 状态服务
func GetStatusService() proto.StatusServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// 注册表服务
func GetRegistryService() proto.RegistryServiceServer {
	panic(wire.Build(serviceProvideSets))
}

// func GetPromService() promotionService {

// }

// // 基础服务
// func GetFoundationService() foundationService {

// }

// // 会员服务
// func GetMemberService() memberService {

// }

// // 商户服务
// func GetMerchantService() merchantService {

// }

// // 商店服务
// func GetShopService() proto.ShopServiceServer {

// }

// // 产品服务
// func GetProductService() productService {

// }

// // 商品服务
// func GetItemService() itemService {

// }

// // 购物服务
// func GetOrderService() orderServiceImpl {

// }

// // 购物车服务
// func GetCartService() cartServiceImpl {

// }

// // 售后服务
// func GetAfterSalesService() afterSalesService {

// }

// // 支付服务
// func GetPaymentService() paymentService {

// }

// // 快捷支付服务
// func GetQuickPayService() quickPayServiceImpl {

// }

// // 消息服务
// func GetMessageService() messageService {

// }

// // 快递服务
// func GetExpressService() expressServiceImpl {

// }

// // 配送服务
// func GetShipmentService() shipmentServiceImpl {

// }

// // 内容服务
// func GetContentService() contentService {

// }

// // 广告服务
// func GetAdService() advertisementService {

// }

// // 钱包服务
// func GetWalletService() walletServiceImpl {

// }

// // 个人金融服务
// func GetPersonFinanceService() personFinanceService {

// }

// // 门户数据服务
// func GetPortalService() portalService {

// }

// // 查询服务
// func GetQueryService() queryService {

// }

// // ExecuteService 执行任务服务
// func GetExecuteService() executionServiceImpl {

// }

// func GetCommonDao() impl.CommonDao {

// }

// // AppService APP服务
// func GetAppService() appServiceImpl {

// }

// // RbacService 权限服务
// func GetRbacService() rbacServiceImpl {

// }

// // CodeService 条码服务
// func GetCodeService() codeServiceImpl {

// }
