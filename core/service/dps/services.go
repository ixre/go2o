/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:20
 * description :
 * history :
 */

package dps

import (
	"github.com/jsix/gof"
	"go2o/core/infrastructure/domain"
	"go2o/core/query"
	"go2o/core/repository"
)

var (
	PromService *promotionService

	// 基础服务
	BaseService *platformService

	// 会员服务
	MemberService *memberService

	// 商户服务
	MerchantService *merchantService

	// 商店服务
	ShopService *shopService

	// 销售服务
	SaleService *saleService

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
	DeliverService *deliveryService

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

	/** Repository **/

	goodsRep := repository.NewGoodsRep(db)
	valRep := repository.NewValueRep(db)
	userRep := repository.NewUserRep(db)
	notifyRep := repository.NewNotifyRep(db)
	mssRep := repository.NewMssRep(db, notifyRep, valRep)
	expressRep := repository.NewExpressRep(db, valRep)
	shipRep := repository.NewShipmentRep(db, expressRep)
	memberRep := repository.NewMemberRep(db, mssRep, valRep)
	itemRep := repository.NewItemRep(db)
	tagSaleRep := repository.NewTagSaleRep(db)
	promRep := repository.NewPromotionRep(db, goodsRep, memberRep)
	cateRep := repository.NewCategoryRep(db, valRep)
	saleRep := repository.NewSaleRep(db, cateRep, valRep, tagSaleRep,
		itemRep, goodsRep, promRep)
	//afterSalesRep := repository.NewAfterSalesRep(db)
	cartRep := repository.NewCartRep(db, memberRep, goodsRep)
	shopRep := repository.NewShopRep(db)
	mchRep := repository.NewMerchantRep(db, shopRep, userRep, mssRep, valRep)
	personFinanceRep := repository.NewPersonFinanceRepository(db, memberRep)
	deliveryRep := repository.NewDeliverRep(db)
	contentRep := repository.NewContentRep(db)
	adRep := repository.NewAdvertisementRep(db)
	spRep := repository.NewOrderRep(db, mchRep, nil, saleRep, cartRep, goodsRep,
		promRep, memberRep, deliveryRep, expressRep, shipRep, valRep)
	payRep := repository.NewPaymentRep(db, memberRep, spRep, valRep)
	asRep := repository.NewAfterSalesRep(db, spRep, memberRep)

	goodsRep.SetSaleRep(saleRep) //fixed
	memberRep.SetMerchantRep(mchRep)
	spRep.SetPaymentRep(payRep)

	/** Query **/
	memberQue := query.NewMemberQuery(db)
	mchQuery := query.NewMerchantQuery(ctx)
	contentQue := query.NewContentQuery(db)
	goodsQuery := query.NewGoodsQuery(db)
	shopQuery := query.NewShopQuery(ctx)
	orderQuery := query.NewOrderQuery(db)
	afterSalesQuery := query.NewAfterSalesQuery(db)

	/** Service **/
	BaseService = NewPlatformService(valRep)
	PromService = NewPromotionService(promRep)
	ShoppingService = NewShoppingService(spRep, saleRep, cartRep,
		itemRep, goodsRep, mchRep, orderQuery)
	AfterSalesService = NewAfterSalesService(asRep, afterSalesQuery, spRep)
	MerchantService = NewMerchantService(mchRep, saleRep, mchQuery, orderQuery)
	ShopService = NewShopService(shopRep, mchRep, shopQuery)
	MemberService = NewMemberService(MerchantService, memberRep, memberQue, orderQuery)
	SaleService = NewSaleService(saleRep, cateRep, goodsRep, goodsQuery)
	PaymentService = NewPaymentService(payRep, spRep)
	MssService = NewMssService(mssRep)
	ExpressService = NewExpressService(expressRep)
	DeliverService = NewDeliveryService(deliveryRep)
	ContentService = NewContentService(contentRep, contentQue)
	AdService = NewAdvertisementService(adRep)
	PersonFinanceService = NewPersonFinanceService(personFinanceRep, memberRep)

	//m := memberRep.GetMember(1)
	//d := m.ProfileManager().GetDeliverAddress()[0]
	//v := d.GetValue()
	//v.Province = 440000
	//v.City = 440600
	//v.District = 440605
	//d.SetValue(&v)
	//d.Save()
}
