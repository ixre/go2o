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
	"github.com/atnet/gof"
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
	partnerRep := repository.NewPartnerRep(db, userRep, memberRep)
	promRep := repository.NewPromotionRep(db, memberRep)
	tagSaleRep := repository.NewTagSaleRep(db)
	saleRep := repository.NewSaleRep(db, tagSaleRep)
	deliveryRep := repository.NewDeliverRep(db)
	contentRep := repository.NewContentRep(db)
	adRep := repository.NewAdvertisementRep(db)
	spRep := repository.NewShoppingRep(db, partnerRep, saleRep, promRep, memberRep, deliveryRep)

	/** Query **/
	memberQue := query.NewMemberQuery(db)
	partnerQue := query.NewPartnerQuery(ctx)
	contentQue := query.NewContentQuery(db)

	/** Service **/
	PromService = NewPromotionService(promRep)
	ShoppingService = NewShoppingService(spRep)
	MemberService = NewMemberService(memberRep, memberQue)
	PartnerService = NewPartnerService(partnerRep, partnerQue)
	SaleService = NewSaleService(saleRep)
	DeliverService = NewDeliveryService(deliveryRep)
	ContentService = NewContentService(contentRep, contentQue)
	AdvertisementService = NewAdvertisementService(adRep)
}
