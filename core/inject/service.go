//go:build wireinject

package inject

import (
	"github.com/google/wire"
	"github.com/ixre/go2o/core/event"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/go2o/core/sp"
)

// GetSPConfig 获取第三方服务自动配置
func GetSPConfig() *sp.ServiceProviderConfiguration {
	panic(wire.Build(InjectProvideSets))
}

func GetEventSource() *event.EventSource {
	panic(wire.Build(InjectProvideSets))
}

// 状态服务
func GetStatusService() proto.StatusServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 注册表服务
func GetRegistryService() proto.RegistryServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// func GetPromService() impl.PromotionService {
// 	panic(wire.Build(InjectProvideSets))
// }

// 基础服务
func GetSystemService() proto.SystemServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 会员服务
func GetMemberService() proto.MemberServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 商户服务
func GetMerchantService() proto.MerchantServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 商店服务
func GetShopService() proto.ShopServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 产品服务
func GetProductService() proto.ProductServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 商品服务
func GetItemService() proto.ItemServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 购物服务
func GetOrderService() proto.OrderServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 购物车服务
func GetCartService() proto.CartServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 售后服务
func GetAfterSalesService() proto.AfterSalesServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 支付服务
func GetPaymentService() proto.PaymentServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 快捷支付服务
func GetQuickPayService() proto.QuickPayServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 消息服务
func GetMessageService() proto.MessageServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 快递服务
func GetExpressService() proto.ExpressServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 配送服务
func GetShipmentService() proto.ShipmentServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 内容服务
func GetContentService() proto.ContentServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 广告服务
func GetAdService() proto.AdvertisementServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 钱包服务
func GetWalletService() proto.WalletServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 个人金融服务
func GetPersonFinanceService() proto.FinanceServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 门户数据服务
func GetPortalService() proto.PortalServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// 查询服务
func GetQueryService() proto.QueryServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// ExecuteService 执行任务服务
func GetExecuteService() proto.ExecutionServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// func GetCommonDao() impl.CommonDao {
// 	panic(wire.Build(InjectProvideSets))
// }

// AppService APP服务
func GetAppService() proto.AppServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// RbacService 权限服务
func GetRbacService() proto.RbacServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// CodeService 条码服务
func GetCodeService() proto.CodeServiceServer {
	panic(wire.Build(InjectProvideSets))
}

// NewCheckService 校验服务
func GetCheckService() proto.CheckServiceServer {
	panic(wire.Build(InjectProvideSets))
}

func GetInvoiceService() proto.InvoiceServiceServer {
	panic(wire.Build(InjectProvideSets))
}

func GetChatService() proto.ChatServiceServer {
	panic(wire.Build(InjectProvideSets))
}

func GetWorkorderService() proto.WorkorderServiceServer {
	panic(wire.Build(InjectProvideSets))
}