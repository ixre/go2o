/**
 * Copyright 2014 @ Ops Inc.
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
	Context         gof.App
	PromService     *promotionService
	ShoppingService *shoppingService
	MemberService   *memberService
	PartnerService  *partnerService
	SaleService     *saleService
	DeliveyService  *deliveryService
)

func Init(ctx gof.App) {
	Context := ctx
	db := Context.Db()

	/** Repository **/
	userRep := repository.NewUserRep(db)
	partnerRep := repository.NewPartnerRep(db, userRep)
	memberRep := repository.NewMemberRep(db)
	promRep := repository.NewPromotionRep(db, memberRep)
	saleRep := repository.NewSaleRep(db)
	deliveryRep := repository.NewDeliverRep(db)
	spRep := repository.NewShoppingRep(db, partnerRep, saleRep, promRep, memberRep, deliveryRep)

	/** Query **/
	mq := query.NewMemberQuery(db)
	pq := query.NewPartnerQuery(db)

	/** Service **/
	PromService = NewPromotionService(promRep)
	ShoppingService = NewShoppingService(spRep)
	MemberService = NewMemberService(memberRep, mq)
	PartnerService = NewPartnerService(partnerRep, pq)
	SaleService = NewSaleService(saleRep)
	DeliveyService = NewDeliveryService(deliveryRep)
}
