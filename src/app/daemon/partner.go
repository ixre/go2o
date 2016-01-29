/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-01-08 21:35
 * description :
 * history :
 */

package daemon

import (
	"github.com/jsix/gof"
	"go2o/src/core/service/dps"
	"log"
)

var (
	partnerIds []int
)

func getPartners() []int {
	if partnerIds == nil {
		partnerIds = dps.PartnerService.GetPartnersId()
	}
	return partnerIds
}

/***** OLD CODE *****/
// todo: 等待重构

func orderDaemon(app gof.App) {
	defer recoverDaemon()
	ids := getPartners()
	for _, v := range ids {
		log.Println("--", v)
		autoSetOrder(v)
	}
}

func autoSetOrder(partnerId int) {
	f := func(err error) {
		appCtx.Log().PrintErr(err)
	}
	dps.ShoppingService.OrderAutoSetup(partnerId, f)
}
