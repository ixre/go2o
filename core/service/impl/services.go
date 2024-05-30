/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:20
 * description :
 * history :
 */

package impl

import (
	"encoding/json"
	"strings"

	"github.com/ixre/go2o/core/dao/impl"
	"github.com/ixre/go2o/core/domain/tmp"
	"github.com/ixre/go2o/core/event"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/repos"
	"github.com/ixre/go2o/core/repos/clickhouse"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
)

var (
	DB db.Connector
)
var (

	// 状态服务
	StatusService *StatusServiceImpl
	// 注册表服务
	RegistryService *registryService
	PromService     *PromotionService
	// 基础服务
	FoundationService *foundationService
	// 会员服务
	MemberService *memberService
	// 商户服务
	MerchantService *merchantService
	// 商店服务
	ShopService proto.ShopServiceServer
	// 产品服务
	ProductService *productService
	// 商品服务
	ItemService *itemService
	// 购物服务
	OrderService *orderServiceImpl
	// 购物车服务
	CartService *cartServiceImpl
	// 售后服务
	AfterSalesService *afterSalesService
	// 支付服务
	PaymentService *paymentService
	// 快捷支付服务
	QuickPayService *quickPayServiceImpl
	// 消息服务
	MessageService *messageService
	// 快递服务
	ExpressService *expressServiceImpl
	// 配送服务
	ShipmentService *shipmentServiceImpl
	// 内容服务
	ContentService *contentService
	// 广告服务
	AdService *advertisementService
	// 钱包服务
	WalletService *walletServiceImpl
	// 个人金融服务
	PersonFinanceService *personFinanceService
	// 门户数据服务
	PortalService *portalService
	// 查询服务
	QueryService *queryService

	// ExecuteService 执行任务服务
	ExecuteService *executionServiceImpl

	CommonDao *impl.CommonDao
	// AppService APP服务
	AppService *appServiceImpl
	// RbacService 权限服务
	RbacService *rbacServiceImpl
	// CodeService 条码服务
	CodeService *codeServiceImpl
)

// 处理错误
func handleError(err error) error {
	return domain.HandleError(err, "service")
	//if err != nil && gof.CurrentApp.Debug() {
	//	gof.CurrentApp.Log().Println("[ GO2O][ Repo][ Error] -", err.Error())
	//}
	//return err
}

func Init(ctx gof.App) {
	Context := ctx
	db := Context.Db()
	sto := Context.Storage()
	o := orm.NewOrm(db.Driver(), db.Raw())
	tmp.SetORM(o)
	orm.CacheProxy(o, sto)
	// 初始化clickhouse
	initializeClickhouse(ctx)
	// 初始化服务
	initService(ctx, db, o, sto)
	// 初始化事件
	event.InitEvent()
	// 初始化数据
	InitData(o)
}

// InitTestService 初始化测试服务
func InitTestService(ctx gof.App, db db.Connector, orm orm.Orm, sto storage.Interface) {
	initService(ctx, db, orm, sto)
	// 初始化clickhouse
	initializeClickhouse(ctx)
	// 初始化事件
	event.InitEvent()
}

// InitializeClickhouse 初始化clickhouse查询连接
func initializeClickhouse(app gof.App) {
	cfg := app.Config()
	server := cfg.GetString("clickhouse_server")
	servers := strings.Split(server, ",")
	database := cfg.GetString("clickhouse_database")
	password := cfg.GetString("clickhouse_password")
	clickhouse.Configure(servers, database, password)
}

// 初始化服务
func initService(ctx gof.App, db db.Connector, orm orm.Orm, sto storage.Interface) {
	OrmMapping(orm)
	//[]string{"127.0.0.1:9000"}
	repos.Initial(orm, sto)
	// registryRepo := inject.GetRegistryRepo()
	// proMRepo := inject.GetProModelRepo()
	// valueRepo := inject.GetValueRepo()
	// mssRepo := inject.GetMssRepo()
	// expressRepo := inject.GetExpressRepo()
	// shipRepo := inject.GetShipmentRepo()
	// memberRepo := inject.GetMemberRepo()
	// productRepo := inject.GetProductRepo()
	// catRepo := inject.GetCategoryRepo()
	// itemRepo := inject.GetItemRepo()
	// tagSaleRepo := inject.GetSaleLabelRepo()
	// promRepo := inject.GetPromotionRepo()

	// shopRepo := inject.GetShopRepo()
	// mchRepo := inject.GetMerchantRepo()
	// cartRepo := inject.GetCartRepo()
	// personFinanceRepo := inject.GetPersonFinanceRepository()
	// deliveryRepo := inject.GetDeliveryRepo()
	// contentRepo := inject.GetContentRepo()
	// adRepo := inject.GetAdRepo()
	// orderRepo := inject.GetOrderRepo()
	// paymentRepo := inject.GetPaymentRepo()
	// asRepo := inject.GetAfterSalesRepo()
	// notifyRepo := inject.GetNotifyRepo()
	// jobRepo := inject.GetJobRepo()

	/** Params **/
	// memberQue := query.NewMemberQuery(orm)
	// mchQuery := query.NewMerchantQuery(ctx)
	// contentQue := query.NewContentQuery(orm)
	// goodsQuery := query.NewItemQuery(orm)
	// shopQuery := query.NewShopQuery(orm, sto)
	// orderQuery := query.NewOrderQuery(orm)
	// afterSalesQuery := query.NewAfterSalesQuery(db)

	/** Service **/
	// StatusService = NewStatusService()
	// RegistryService = NewRegistryService(valueRepo, registryRepo)
	// ProductService = NewProductService(proMRepo, catRepo, productRepo)
	// FoundationService = NewFoundationService(valueRepo, registryRepo, sto, notifyRepo)
	// PromService = NewPromotionService(promRepo)
	// OrderService = NewShoppingService(orderRepo, cartRepo, memberRepo,
	// 	productRepo, itemRepo, mchRepo, shopRepo,
	// 	paymentRepo, shipRepo, expressRepo, orderQuery)
	// CartService = NewCartService(cartRepo, itemRepo, mchRepo, shopRepo)
	// AfterSalesService = NewAfterSalesService(asRepo, afterSalesQuery, orderRepo)
	// MerchantService = NewMerchantService(mchRepo, memberRepo, mchQuery, orderQuery)
	// ShopService = NewShopService(shopRepo, mchRepo, shopRepo, registryRepo, shopQuery)
	// MemberService = NewMemberService(MerchantService, memberRepo, registryRepo, memberQue, orderQuery, valueRepo)
	// ItemService = NewSaleService(sto, catRepo, itemRepo, goodsQuery, tagSaleRepo, proMRepo, mchRepo, valueRepo)
	// PaymentService = NewPaymentService(paymentRepo, orderRepo, memberRepo)
	// QuickPayService = NewQuickPayService(sto, registryRepo)
	// MessageService = NewMessageService(mssRepo)
	// ExpressService = NewExpressService(expressRepo)
	// ShipmentService = NewShipmentService(shipRepo, deliveryRepo, expressRepo)
	// ContentService = NewContentService(contentRepo, contentQue)
	// AdService = NewAdvertisementService(adRepo, sto)
	// PersonFinanceService = NewPersonFinanceService(personFinanceRepo, memberRepo)

	// WalletService = NewWalletService(inject.GetWalletRepo())
	// ExecuteService = NewExecDataService(sto, jobRepo)

	// CommonDao = impl.NewCommDao(orm, sto, adRepo, catRepo)
	// portalDao := impl.NewPortalDao(orm)
	// PortalService = NewPortalService(orm, CommonDao, portalDao)
	// QueryService = NewQueryService(orm, sto, catRepo)
	// AppService = NewAppService(sto, orm)
	// RbacService = NewRbacService(sto, orm, registryRepo)
	// CodeService = NewCodeService(sto, orm)
}

// 服务工具类，实现的服务组合此类,可直接调用其方法
type serviceUtil struct{}

// 返回失败的结果
func (s serviceUtil) failResult(msg string) *proto.Result {
	return s.failCodeResult(1, msg)
}

// 返回错误的结果
func (s serviceUtil) error(err error) *proto.Result {
	if err == nil {
		return s.success(nil)
	}
	return s.failResult(err.Error())
}

// 返回结果
func (s serviceUtil) result(err error) *proto.Result {
	if err == nil {
		return s.success(nil)
	}
	return s.error(err)
}

// 返回自定义编码的结果
func (s serviceUtil) resultWithCode(code int, message string) *proto.Result {
	return &proto.Result{ErrCode: int32(code), ErrMsg: message, Data: map[string]string{}}
}

// 返回失败的结果
func (s serviceUtil) errorCodeResult(code int, err error) *proto.Result {
	return &proto.Result{ErrCode: int32(code), ErrMsg: err.Error(), Data: map[string]string{}}
}

// 返回失败的结果
func (s serviceUtil) failCodeResult(code int, msg string) *proto.Result {
	return &proto.Result{ErrCode: int32(code), ErrMsg: msg, Data: map[string]string{}}
}

// 返回成功的结果
func (s serviceUtil) success(data map[string]string) *proto.Result {
	if data == nil {
		data = map[string]string{}
	}
	return &proto.Result{ErrCode: 0, ErrMsg: "", Data: data}
}

// 将int32数组装换为int数组
func (s serviceUtil) intArray(values []int32) []int {
	arr := make([]int, len(values))
	for i, v := range values {
		arr[i] = int(v)
	}
	return arr
}

// 转换为JSON
func (s serviceUtil) json(data interface{}) string {
	if data == nil {
		return "{}"
	}
	r, err := json.Marshal(data)
	if err != nil {
		return "{\"error\":\"parse error:" + err.Error() + "\"}"
	}
	return string(r)
}

// 分页响应结果
func (s serviceUtil) pagingResult(total int, data interface{}) *proto.SPagingResult {
	switch data.(type) {
	case string:
		return &proto.SPagingResult{
			Count:  int32(total),
			Data:   data.(string),
			Extras: map[string]string{},
		}
	}
	r, _ := json.Marshal(data)
	return &proto.SPagingResult{
		ErrCode: 0,
		ErrMsg:  "",
		Count:   int32(total),
		Data:    string(r),
		Extras:  map[string]string{},
	}
}
