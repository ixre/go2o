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
	"time"
)

// 守护进程服务
type DaemonService func(gof.App)
type DaemonFunc func(gof.App)

var (
	appCtx                 *core.MainApp
	services               map[string]DaemonService = map[string]DaemonService{}
	tickerDuration                                  = 5 * time.Second // 间隔5秒执行
	tickerInvokeFunc       []DaemonFunc             = []DaemonFunc{}
	newOrderObserver       []DaemonFunc             = []DaemonFunc{confirmOrderQueue}
	completedOrderObserver []DaemonFunc             = []DaemonFunc{}

//newMemberObserver []DaemonFunc = []DaemonFunc{orderDaemon}
)

func RegisterService(name string, service DaemonService) {
	if _, ok := services[name]; ok {
		panic("service named " + name + " is registed!")
	}
	services[name] = service
}

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

// 添加定时执行任务(默认5秒)
func AddTickerFunc(f DaemonFunc) {
	tickerInvokeFunc = append(tickerInvokeFunc, f)
}

// 获取订单处理函数
func orderDaemonService(app gof.App) {
	AddTickerFunc(func(app gof.App) {
		confirmNewOrder(app, newOrderObserver)

		if completedOrderObserver != nil && len(completedOrderObserver) != 0 {
			completedOrderObs(app, completedOrderObserver)
		}
	})
	orderDaemon(app)
}

// 添加新的订单处理函数
func AddNewOrderFunc(f DaemonFunc) {
	newOrderObserver = append(newOrderObserver, f)
}

// 添加已完成订单处理函数
func AddCompletedOrderFunc(f DaemonFunc) {
	completedOrderObserver = append(completedOrderObserver, f)
}

func RegisterByName(arr []string) {
	for _, v := range arr {
		switch v {
		case "mail":
			RegisterService("mail", func(app gof.App) {
				AddTickerFunc(startMailQueue)
			})
		case "order":
			RegisterService("order", orderDaemonService)
		}
	}
}

func Start() {
	tk := time.NewTicker(tickerDuration)
	defer func() {
		tk.Stop()
	}()

	for name, s := range services {
		log.Println("** [ Go2o][ Daemon][ Booted] - ", name, " daemon running")
		go s(appCtx)
	}

	for {
		select {
		case <-tk.C:
			for _, f := range tickerInvokeFunc {
				f(appCtx)
			}
		}
	}
}

func recoverDaemon() {

}
