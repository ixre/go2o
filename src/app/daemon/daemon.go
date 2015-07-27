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
)

var (
	gCTX     gof.App
	mailChan chan int
)

func Run(ctx gof.App) {
	gCTX = ctx
	//daemon()
	startMailQueue()
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
