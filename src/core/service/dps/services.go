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
	"go2o/src/core/query"
	"go2o/src/core/repository"
)

var (
	PromService          *promotionService
	ShoppingService      *shoppingService
	MemberService        *memberService

	// 商户服务
	MerchantService      *merchantService

	// 销售服务
	SaleService          *saleService

	DeliverService       *deliveryService

	// 内容服务
	ContentService       *contentService

	// 广告服务
	AdService *adService

	// 个人金融服务
	PersonFinanceService *personFinanceService
)

func Init(ctx gof.App) {
	Context := ctx
	db := Context.Db()

	/** Repository **/
	userRep := repository.NewUserRep(db)
	memberRep := repository.NewMemberRep(db)
	goodsRep := repository.NewGoodsRep(db)
	tagSaleRep := repository.NewTagSaleRep(db)
	promRep := repository.NewPromotionRep(db, goodsRep, memberRep)
	saleRep := repository.NewSaleRep(db, tagSaleRep, goodsRep, promRep)
	mssRep := repository.NewMssRep(db)
	partnerRep := repository.NewMerchantRep(db, userRep, memberRep, mssRep)
	memberRep.SetMerchantRep(partnerRep)
	personFinanceRep := repository.NewPersonFinanceRepository(db, memberRep)

	deliveryRep := repository.NewDeliverRep(db)
	contentRep := repository.NewContentRep(db)
	adRep := repository.NewAdvertisementRep(db)
	spRep := repository.NewShoppingRep(db, partnerRep, saleRep, goodsRep, promRep, memberRep, deliveryRep)

	/** Query **/
	memberQue := query.NewMemberQuery(db)
	partnerQue := query.NewMerchantQuery(ctx)
	contentQue := query.NewContentQuery(db)
	goodsQuery := query.NewGoodsQuery(db)

	/** Service **/
	PromService = NewPromotionService(promRep)
	ShoppingService = NewShoppingService(spRep)
	MerchantService = NewMerchantService(partnerRep, saleRep, partnerQue)
	MemberService = NewMemberService(MerchantService, memberRep, memberQue)
	SaleService = NewSaleService(saleRep, goodsRep, goodsQuery)
	DeliverService = NewDeliveryService(deliveryRep)
	ContentService = NewContentService(contentRep, contentQue)
	AdService = NewAdvertisementService(adRep)
	PersonFinanceService = NewPersonFinanceService(personFinanceRep, memberRep)
}
