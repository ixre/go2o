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
	"github.com/jsix/gof"
	"go2o/core/infrastructure/domain"
	"go2o/core/query"
	"go2o/core/repository"
)

var (
	PromService *promotionService
	// 基础服务
	FoundationService *foundationService
	// 会员服务
	MemberService *memberService
	// 商户服务
	MerchantService *merchantService
	// 商店服务
	ShopService *shopService
	// 销售服务
	SaleService *saleService
	// 产品服务
	ProductService *productService
	// 购物服务
	ShoppingService *shoppingService
	// 售后服务
	AfterSalesService *afterSalesService
	// 支付服务
	PaymentService *paymentService
	// 消息服务
	MssService *mssService
	// 快递服务
	ExpressService *expressService
	// 配送服务
	ShipmentService *shipmentService
	// 内容服务
	ContentService *contentService
	// 广告服务
	AdService *adService

	// 个人金融服务
	PersonFinanceService *personFinanceService
)

// 处理错误
func handleError(err error) error {
	return domain.HandleError(err, "service")
	//if err != nil && gof.CurrentApp.Debug() {
	//	gof.CurrentApp.Log().Println("[ Go2o][ Rep][ Error] -", err.Error())
	//}
	//return err
}

func Init(ctx gof.App) {
	Context := ctx
	db := Context.Db()
	orm := db.GetOrm()
	sto := Context.Storage()

	/** Repository **/
	proMRepo := repository.NewProModelRepo(db, orm)
	goodsRepo := repository.NewGoodsRepo(db)
	valRepo := repository.NewValueRepo(db, sto)
	userRepo := repository.NewUserRepo(db)
	notifyRepo := repository.NewNotifyRepo(db)
	mssRepo := repository.NewMssRepo(db, notifyRepo, valRepo)
	expressRepo := repository.NewExpressRepo(db, valRepo)
	shipRepo := repository.NewShipmentRepo(db, expressRepo)
	memberRepo := repository.NewMemberRepo(sto, db, mssRepo, valRepo)
	itemRepo := repository.NewItemRepo(db)
	tagSaleRepo := repository.NewTagSaleRepo(db)
	promRepo := repository.NewPromotionRepo(db, goodsRepo, memberRepo)
	cateRepo := repository.NewCategoryRepo(db, valRepo, sto)
	saleRepo := repository.NewSaleRepo(db, cateRepo, valRepo, tagSaleRepo,
		itemRepo, expressRepo, goodsRepo, promRepo)
	//afterSalesRepo := repository.NewAfterSalesRepo(db)
	cartRepo := repository.NewCartRepo(db, memberRepo, goodsRepo)
	shopRepo := repository.NewShopRepo(db, sto)
	mchRepo := repository.NewMerchantRepo(db, sto, shopRepo, userRepo, memberRepo, mssRepo, valRepo)
	personFinanceRepo := repository.NewPersonFinanceRepository(db, memberRepo)
	deliveryRepo := repository.NewDeliverRepo(db)
	contentRepo := repository.NewContentRepo(db)
	adRepo := repository.NewAdvertisementRepo(db, sto)
	spRepo := repository.NewOrderRepo(sto, db, mchRepo, nil, saleRepo, cartRepo, goodsRepo,
		promRepo, memberRepo, deliveryRepo, expressRepo, shipRepo, valRepo)
	paymentRepo := repository.NewPaymentRepo(sto, db, memberRepo, spRepo, valRepo)
	asRepo := repository.NewAfterSalesRepo(db, spRepo, memberRepo, paymentRepo)

	goodsRepo.SetSaleRepo(saleRepo) //fixed
	spRepo.SetPaymentRepo(paymentRepo)

	/** Query **/
	memberQue := query.NewMemberQuery(db)
	mchQuery := query.NewMerchantQuery(ctx)
	contentQue := query.NewContentQuery(db)
	goodsQuery := query.NewGoodsQuery(db)
	shopQuery := query.NewShopQuery(ctx)
	orderQuery := query.NewOrderQuery(db)
	afterSalesQuery := query.NewAfterSalesQuery(db)

	/** Service **/
	ProductService = NewProService(proMRepo)
	FoundationService = NewFoundationService(valRepo)
	PromService = NewPromotionService(promRepo)
	ShoppingService = NewShoppingService(spRepo, saleRepo, cartRepo,
		itemRepo, goodsRepo, mchRepo, orderQuery)
	AfterSalesService = NewAfterSalesService(asRepo, afterSalesQuery, spRepo)
	MerchantService = NewMerchantService(mchRepo, saleRepo, mchQuery, orderQuery)
	ShopService = NewShopService(shopRepo, mchRepo, shopQuery)
	MemberService = NewMemberService(MerchantService, memberRepo, memberQue, orderQuery, valRepo)
	SaleService = NewSaleService(saleRepo, cateRepo, goodsRepo, goodsQuery)
	PaymentService = NewPaymentService(paymentRepo, spRepo)
	MssService = NewMssService(mssRepo)
	ExpressService = NewExpressService(expressRepo)
	ShipmentService = NewShipmentService(shipRepo, deliveryRepo)
	ContentService = NewContentService(contentRepo, contentQue)
	AdService = NewAdvertisementService(adRepo, sto)
	PersonFinanceService = NewPersonFinanceService(personFinanceRepo, memberRepo)

	//m := memberRepo.GetMember(1)
	//d := m.ProfileManager().GetDeliverAddress()[0]
	//v := d.GetValue()
	//v.Province = 440000
	//v.City = 440600
	//v.District = 440605
	//d.SetValue(&v)
	//d.Save()
}
