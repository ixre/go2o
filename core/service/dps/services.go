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
	PromService     *promotionService
	ShoppingService *shoppingService

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

	// 消息服务
	MssService *mssService

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

	valRep := repository.NewValueRep(db)
	userRep := repository.NewUserRep(db)
	mssRep := repository.NewMssRep(db)
	memberRep := repository.NewMemberRep(db, mssRep, valRep)
	itemRep := repository.NewItemRep(db)
	goodsRep := repository.NewGoodsRep(db)
	tagSaleRep := repository.NewTagSaleRep(db)
	promRep := repository.NewPromotionRep(db, goodsRep, memberRep)
	cateRep := repository.NewCategoryRep(db, valRep)
	saleRep := repository.NewSaleRep(db, cateRep, valRep, tagSaleRep,
		itemRep, goodsRep, promRep)
	shopRep := repository.NewShopRep(db)
	mchRep := repository.NewMerchantRep(db, shopRep, userRep, mssRep, valRep)
	memberRep.SetMerchantRep(mchRep)
	personFinanceRep := repository.NewPersonFinanceRepository(db, memberRep)

	deliveryRep := repository.NewDeliverRep(db)
	contentRep := repository.NewContentRep(db)
	adRep := repository.NewAdvertisementRep(db)
	spRep := repository.NewShoppingRep(db, mchRep, saleRep, goodsRep,
		promRep, memberRep, deliveryRep, valRep)

	/** Query **/
	memberQue := query.NewMemberQuery(db)
	partnerQue := query.NewMerchantQuery(ctx)
	contentQue := query.NewContentQuery(db)
	goodsQuery := query.NewGoodsQuery(db)
	shopQuery := query.NewShopQuery(ctx)

	/** Service **/
	BaseService = NewPlatformService(valRep)
	PromService = NewPromotionService(promRep)
	ShoppingService = NewShoppingService(spRep, saleRep, itemRep, goodsRep, mchRep)
	MerchantService = NewMerchantService(mchRep, saleRep, partnerQue)
	ShopService = NewShopService(shopRep, mchRep, shopQuery)
	MemberService = NewMemberService(MerchantService, memberRep, memberQue)
	SaleService = NewSaleService(saleRep, cateRep, goodsRep, goodsQuery)
	MssService = NewMssService(mssRep)
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
