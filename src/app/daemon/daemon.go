/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-01-08 21:01
 * description :
 * history :
 */

package daemon

import (
	"flag"
	"github.com/jsix/gof"
	"go2o/src/core"
	"go2o/src/core/service/dps"
	"log"
	"strings"
)

// 守护进程服务
type DaemonService func()

var (
	services map[string]DaemonService = map[string]DaemonService{}
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
	RegisterByName([]string{"mail", "order"})
	Start()
}

// 自定义参数运行
func FlagRun() {
	var conf string
	var debug bool
	var trace bool
	var service string
	var serviceArr []string = []string{"mail", "order"}
	var ch chan bool = make(chan bool)
	flag.StringVar(&conf, "conf", "app.conf", "")
	flag.BoolVar(&debug, "debug", true, "")
	flag.BoolVar(&trace, "trace", true, "")
	flag.StringVar(&service, "service", strings.Join(serviceArr, ","), "")

	flag.Parse()

	appCtx = getAppCtx(conf)
	appCtx.Init(debug, trace)
	gof.CurrentApp = appCtx

	dps.Init(appCtx)

	if service != "all" {
		serviceArr = strings.Split(service, ",")
	}

	RegisterByName(serviceArr)
	Start()

	<-ch
}

func getAppCtx(conf string) *core.MainApp {
	return core.NewMainApp(conf)
}

func RegisterService(name string, service DaemonService) {
	if _, ok := services[name]; ok {
		panic("service named " + name + " is registed!")
	}
	services[name] = service
}

func RegisterByName(arr []string) {
	for _, v := range arr {
		switch v {
		case "mail":
			RegisterService("mail", startMailQueue)
		case "order":
			RegisterService("order", orderDaemon)
		}
	}
}

func Start() {
	//log.Println("[ Go2o][ Daemon][ Booted] - Daemon service is running.")
	for name, s := range services {
		log.Println("[ Go2o][ Daemon][ Booted] - ", name, " daemon running")
		go s()
	}
}

func recoverDaemon() {

}
