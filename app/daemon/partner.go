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
	"go2o/core/service/dps"
	"log"
)

var (
	merchantIds []int
)

func getMerchants() []int {
	if merchantIds == nil {
		merchantIds = dps.MerchantService.GetMerchantsId()
	}
	return merchantIds
}

/***** OLD CODE *****/
// todo: 等待重构

func orderDaemon(app gof.App) {
	defer recoverDaemon()
	ids := getMerchants()
	for _, v := range ids {
		log.Println("--", v)
		autoSetOrder(v)
	}
}

func autoSetOrder(merchantId int) {
	f := func(err error) {
		appCtx.Log().Error(err)
	}
	dps.ShoppingService.OrderAutoSetup(merchantId, f)
}
