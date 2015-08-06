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
	"flag"
	"fmt"
	"github.com/atnet/gof"
	"go2o/src/core"
)

var (
	appCtx *core.MainApp
)

// 运行
func Run(ctx gof.App) {
	if ctx != nil {
		appCtx = ctx.(*core.MainApp)
	} else {
		appCtx = getAppCtx("app.conf")
	}

	fmt.Println("[ Go2o][ Daemon][ Booted] - Daemon service is running.")
	//daemon()
	startMailQueue()
}

// 自定义参数运行
func FlagRun() {
	var conf string
	flag.StringVar(&conf, "conf", "app.conf", "")
	Run(getAppCtx(conf))
}

func getAppCtx(conf string) *core.MainApp {
	return core.NewMainApp(conf)
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
