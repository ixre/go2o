/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:20
 * description :
 * history :
 */

package dps

import (
	"github.com/jrsix/gof"
	"go2o/src/core/query"
	"go2o/src/core/repository"
)

var (
	Context              gof.App
	PromService          *promotionService
	ShoppingService      *shoppingService
	MemberService        *memberService
	PartnerService       *partnerService
	SaleService          *saleService
	DeliverService       *deliveryService
	ContentService       *contentService
	AdvertisementService *advertisementService
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
	partnerRep := repository.NewPartnerRep(db, userRep, memberRep, mssRep)
	memberRep.SetPartnerRep(partnerRep)

	deliveryRep := repository.NewDeliverRep(db)
	contentRep := repository.NewContentRep(db)
	adRep := repository.NewAdvertisementRep(db)
	spRep := repository.NewShoppingRep(db, partnerRep, saleRep, goodsRep, promRep, memberRep, deliveryRep)

	/** Query **/
	memberQue := query.NewMemberQuery(db)
	partnerQue := query.NewPartnerQuery(ctx)
	contentQue := query.NewContentQuery(db)

	/** Service **/
	PromService = NewPromotionService(promRep)
	ShoppingService = NewShoppingService(spRep)
	MemberService = NewMemberService(memberRep, memberQue)
	PartnerService = NewPartnerService(partnerRep, saleRep, adRep, partnerQue)
	SaleService = NewSaleService(saleRep, goodsRep)
	DeliverService = NewDeliveryService(deliveryRep)
	ContentService = NewContentService(contentRep, contentQue)
	AdvertisementService = NewAdvertisementService(adRep)
}
