/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-01-08 21:35
 * description :
 * history :
 */

package daemon

import (
	"go2o/src/core/service/dps"
	"time"
)

func orderDaemon() {
	defer recoverDaemon()
	for {
		ids := getPartners()
		for _, v := range ids {
			autoSetOrder(v)
		}
		time.Sleep(time.Minute * CRON_ORDER_SETUP_MINUTE)
	}
}

func autoSetOrder(partnerId int) {
	f := func(err error) {
		appCtx.Log().PrintErr(err)
	}
	dps.ShoppingService.OrderAutoSetup(partnerId, f)
}
