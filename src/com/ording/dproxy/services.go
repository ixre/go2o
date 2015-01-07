/**
 * Copyright 2014 @ Ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-03 23:20
 * description :
 * history :
 */

package dproxy

import (
	"com/repository"
	"github.com/newmin/gof/app"
)

var (
	Context         app.Context
	PromService     *promotionService
	ShoppingService *shoppingService
	MemberService   *memberService
	PartnerService  *partnerService
	SaleService     *saleService
)

func Init(ctx app.Context) {
	Context := ctx
	db := Context.Db()
	partnerRep := repository.NewPartnerRep(db)
	memberRep := &repository.MemberRep{Connector: db}
	promRep := &repository.PromotionRep{Connector: db, MemberRep: memberRep}
	saleRep := repository.NewSaleRep(db)
	spRep := repository.NewShoppingRep(db, partnerRep, saleRep, promRep, memberRep)

	PromService = &promotionService{promRep: promRep}
	ShoppingService = &shoppingService{spRep: spRep}
	MemberService = &memberService{memberRep: memberRep}
	PartnerService = &partnerService{partnerRep: partnerRep}
	SaleService = &saleService{saleRep: saleRep}
}
