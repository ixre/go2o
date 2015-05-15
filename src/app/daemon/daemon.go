/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-01-08 21:01
 * description :
 * history :
 */

package daemon

import (
	"github.com/atnet/gof"
	"time"
)

var (
	gCTX gof.App
)

func Run(ctx gof.App) {
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
