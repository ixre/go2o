/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:20
 * description :
 * history :
 */

package rsi

import (
	"github.com/ixre/gof"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"go2o/app"
	"go2o/core/dao"
	"go2o/core/factory"
	"go2o/core/infrastructure/domain"
	"go2o/core/query"
	"go2o/core/service/auto_gen/rpc/ttype"
	"go2o/core/service/auto_gen/rpc/wallet_service"
	"go2o/core/variable"
	"strconv"
	"strings"
	"time"
)

var (
	fact        *factory.RepoFactory
	PromService *promotionService
	// 基础服务
	FoundationService *foundationService
	// 会员服务
	MemberService *memberService
	// 商户服务
	MerchantService *merchantService
	// 商店服务
	ShopService *shopServiceImpl
	// 产品服务
	ProductService *productService
	// 商品服务
	ItemService *itemService
	// 购物服务
	ShoppingService *orderServiceImpl
	// 售后服务
	AfterSalesService *afterSalesService
	// 支付服务
	PaymentService *paymentService
	// 消息服务
	MssService *mssService
	// 快递服务
	ExpressService *expressService
	// 配送服务
	ShipmentService *shipmentServiceImpl
	// 内容服务
	ContentService *contentService
	// 广告服务
	AdService *adService
	// 钱包服务
	WalletService wallet_service.WalletService
	// 个人金融服务
	PersonFinanceService *personFinanceService
	// 门户数据服务
	PortalService *portalService

	CommonDao *dao.CommonDao
)

// 处理错误
func handleError(err error) error {
	return domain.HandleError(err, "service")
	//if err != nil && gof.CurrentApp.Debug() {
	//	gof.CurrentApp.Log().Println("[ Go2o][ Repo][ Error] -", err.Error())
	//}
	//return err
}

func Init(ctx gof.App, appFlag int, confDir string) {
	Context := ctx
	db := Context.Db()
	orm := db.GetOrm()
	sto := Context.Storage()

	// 初始化服务
	initService(ctx, db, orm, sto, confDir)
	// RPC
	if appFlag&app.FlagRpcServe == app.FlagRpcServe {
		initRpcServe(ctx)
	}
}

func initService(ctx gof.App, db db.Connector, orm orm.Orm,
	sto storage.Interface, confPath string) {
	fact = (&factory.RepoFactory{}).Init(db, sto, confPath)
	proMRepo := fact.GetProModelRepo()
	valueRepo := fact.GetValueRepo()
	mssRepo := fact.GetMssRepo()
	expressRepo := fact.GetExpressRepo()
	shipRepo := fact.GetShipmentRepo()
	memberRepo := fact.GetMemberRepo()
	productRepo := fact.GetProductRepo()
	catRepo := fact.GetCategoryRepo()
	itemRepo := fact.GetItemRepo()
	tagSaleRepo := fact.GetSaleLabelRepo()
	promRepo := fact.GetPromotionRepo()

	shopRepo := fact.GetShopRepo()
	mchRepo := fact.GetMerchantRepo()
	cartRepo := fact.GetCartRepo()
	personFinanceRepo := fact.GetPersonFinanceRepository()
	deliveryRepo := fact.GetDeliveryRepo()
	contentRepo := fact.GetContentRepo()
	adRepo := fact.GetAdRepo()
	orderRepo := fact.GetOrderRepo()
	paymentRepo := fact.GetPaymentRepo()
	asRepo := fact.GetAfterSalesRepo()

	/** Query **/
	memberQue := query.NewMemberQuery(db)
	mchQuery := query.NewMerchantQuery(ctx)
	contentQue := query.NewContentQuery(db)
	goodsQuery := query.NewItemQuery(db)
	shopQuery := query.NewShopQuery(ctx)
	orderQuery := query.NewOrderQuery(db)
	afterSalesQuery := query.NewAfterSalesQuery(db)

	/** Service **/
	ProductService = NewProService(proMRepo, catRepo, productRepo)
	FoundationService = NewFoundationService(valueRepo)
	PromService = NewPromotionService(promRepo)
	ShoppingService = NewShoppingService(orderRepo, cartRepo, memberRepo,
		productRepo, itemRepo, mchRepo, orderQuery)
	AfterSalesService = NewAfterSalesService(asRepo, afterSalesQuery, orderRepo)
	MerchantService = NewMerchantService(mchRepo, memberRepo, mchQuery, orderQuery)
	ShopService = NewShopService(shopRepo, mchRepo, shopQuery)
	MemberService = NewMemberService(MerchantService, memberRepo, memberQue, orderQuery, valueRepo)
	ItemService = NewSaleService(sto, catRepo, itemRepo, goodsQuery, tagSaleRepo, proMRepo, mchRepo, valueRepo)
	PaymentService = NewPaymentService(paymentRepo, orderRepo, memberRepo)
	MssService = NewMssService(mssRepo)
	ExpressService = NewExpressService(expressRepo)
	ShipmentService = NewShipmentService(shipRepo, deliveryRepo, expressRepo)
	ContentService = NewContentService(contentRepo, contentQue)
	AdService = NewAdvertisementService(adRepo, sto)
	PersonFinanceService = NewPersonFinanceService(personFinanceRepo, memberRepo)

	WalletService = NewWalletService(fact.GetWalletRepo())

	CommonDao = dao.NewCommDao(orm, sto, adRepo, catRepo)
	PortalService = NewPortalService(CommonDao)
}

// RPC服务初始化
func initRpcServe(ctx gof.App) {
	gf := ctx.Config().GetString
	mp := make(map[string]string)
	domain := gf("domain")
	hash := gf("url_hash")
	if hash == "" {
		hash = crypto.Md5([]byte(strconv.Itoa(int(time.Now().Unix()))))[8:14]
	}
	ssl := gf("ssl_enabled")
	prefix := "http://"
	if ssl == "true" || ssl == "1" {
		prefix = "https://"
	}
	mp[variable.DEnabledSSL] = gf("ssl_enabled")
	mp[variable.DStaticServer] = gf("static_server")
	mp[variable.DImageServer] = gf("image_server")
	mp[variable.DUrlHash] = hash
	mp[variable.DRetailPortal] = strings.Join([]string{prefix,
		variable.DOMAIN_PREFIX_PORTAL, domain}, "")
	mp[variable.DWholesalePortal] = strings.Join([]string{prefix,
		variable.DOMAIN_PREFIX_WHOLESALE_PORTAL, domain}, "")
	mp[variable.DUCenter] = strings.Join([]string{prefix,
		variable.DOMAIN_PREFIX_MEMBER, domain}, "")
	mp[variable.DPassport] = strings.Join([]string{prefix,
		variable.DOMAIN_PREFIX_PASSPORT, domain}, "")
	mp[variable.DMerchant] = strings.Join([]string{prefix,
		variable.DOMAIN_PREFIX_MERCHANT, domain}, "")
	mp[variable.DHApi] = strings.Join([]string{prefix,
		variable.DOMAIN_PREFIX_HApi, domain}, "")

	mp[variable.DRetailMobilePortal] = strings.Join([]string{prefix,
		variable.DOMAIN_PREFIX_PORTAL_MOBILE, domain}, "")
	mp[variable.DWholesaleMobilePortal] = strings.Join([]string{prefix,
		variable.DOMAIN_PREFIX_M_WHOLESALE, domain}, "")
	mp[variable.DMobilePassport] = strings.Join([]string{prefix,
		variable.DOMAIN_PREFIX_M_PASSPORT, domain}, "")
	mp[variable.DMobileUCenter] = strings.Join([]string{prefix,
		variable.DOMAIN_PREFIX_M_MEMBER, domain}, "")

	fact.GetValueRepo().SavesRegistry(mp)
}

// 服务工具类，实现的服务组合此类,可直接调用其方法
type serviceUtil struct{}

// 返回失败的结果
func (s serviceUtil) failResult(msg string) *ttype.Result_ {
	return s.failCodeResult(1, msg)
}

// 返回错误的结果
func (s serviceUtil) error(err error) *ttype.Result_ {
	return s.failResult(err.Error())
}

// 返回结果
func (s serviceUtil) result(err error) *ttype.Result_ {
	if err == nil {
		return s.success(nil)
	}
	return s.error(err)
}

// 返回失败的结果
func (s serviceUtil) errorCodeResult(code int, err error) *ttype.Result_ {
	return &ttype.Result_{ErrCode: int32(code), ErrMsg: err.Error(), Data: map[string]string{}}
}

// 返回失败的结果
func (s serviceUtil) failCodeResult(code int, msg string) *ttype.Result_ {
	return &ttype.Result_{ErrCode: int32(code), ErrMsg: msg, Data: map[string]string{}}
}

// 返回成功的结果
func (s serviceUtil) success(data map[string]string) *ttype.Result_ {
	return &ttype.Result_{ErrCode: 0, ErrMsg: "", Data: data}
}

// 将int32数组装换为int数组
func (s serviceUtil) intArray(values []int32) []int {
	arr := make([]int, len(values))
	for i, v := range values {
		arr[i] = int(v)
	}
	return arr
}
