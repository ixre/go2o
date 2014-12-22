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
	"ops/cf/app"
)

var (
	Context       app.Context
	PromService   *promotionService
	SpService     *shoppingService
	MemberService *memberService
	PartnerService *partnerService
)

func Init(ctx app.Context) {
	Context := ctx
	db := Context.Db()
	partnerRep := &repository.PartnerRep{Connector: db}
	memberRep := &repository.MemberRep{Connector: db}
	promRep := &repository.PromotionRep{Connector: db, MemberRep: memberRep}
	saleRep := &repository.SaleRep{Connector: db}
	spRep := &repository.ShoppingRep{Connector: db, PartnerRep: partnerRep,
		SaleRep: saleRep, PromRep: promRep, MemberRep: memberRep}

	PromService = &promotionService{promRep: promRep}
	SpService = &shoppingService{spRep: spRep}
	MemberService = &memberService{memberRep: memberRep}
	PartnerService = &partnerService{partnerRep:partnerRep}
}
