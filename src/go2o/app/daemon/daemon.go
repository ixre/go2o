/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-01-08 21:01
 * description :
 * history :
 */

package daemon

import (
	"github.com/atnet/gof/app"
	"time"
)

var (
	gCTX app.Context
)

func Run(ctx app.Context) {
	gCTX = ctx
	daemon()
}

func daemon() {
	go partnerDaemon()
	go orderDaemon()
}

func recoverDaemon() {

}

func partnerDaemon() {
	defer recoverDaemon()
}

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
