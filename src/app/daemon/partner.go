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
<<<<<<< HEAD
	"github.com/jsix/gof"
	"go2o/src/core/service/dps"
	"log"
=======
	"go2o/src/core/service/dps"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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
<<<<<<< HEAD

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
=======
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
