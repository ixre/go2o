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
type Service func(gof.App)

// 守护进程执行的函数
type Func func(gof.App)

var (
	appCtx                 *core.MainApp
	services               map[string]Service = map[string]Service{}
	tickerDuration                            = 5 * time.Second // 间隔5秒执行
	tickerInvokeFunc       []Func             = []Func{}
	orderCreatedObserver   []Func             = []Func{confirmOrderQueue} //订单提交通知
	orderCompletedObserver []Func             = []Func{}                  //订单完成通知

	//newMemberObserver []DaemonFunc = []DaemonFunc{orderDaemon}
)

// 注册服务
func RegisterService(name string, service Service) {
	if service == nil {
		panic("service is nil")
	}
	if _, ok := services[name]; ok {
		panic("service named " + name + " is registed!")
	}
	services[name] = service
}

// 添加定时执行任务(默认5秒)
func AddTickerFunc(f Func) {
	tickerInvokeFunc = append(tickerInvokeFunc, f)
}

// 添加新的订单处理函数
func AppendOrderCreatedObserver(f Func) {
	orderCreatedObserver = append(orderCreatedObserver, f)
}

// 添加已完成订单处理函数
func AppendOrderCompletedObserver(f Func) {
	orderCompletedObserver = append(orderCompletedObserver, f)
}

// 启动守护进程
func Start() {
	loadNesTasks()
	for name, s := range services { //运行自定义服务
		log.Println("** [ Go2o][ Daemon][ Booted] - ", name, " daemon running")
		go s(appCtx)
	}

	tk := time.NewTicker(tickerDuration)
	defer func() {
		tk.Stop()
	}()
	for { //执行定时任务
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

func loadNesTasks() {
	AddTickerFunc(orderDaemon)
}

// 获取订单处理函数
func orderService(app gof.App) {
	AddTickerFunc(func(app gof.App) {
		confirmNewOrder(app, orderCreatedObserver)     //确认新订单
		completedOrderObs(app, orderCompletedObserver) //通知完成的订单
	})
}

func RegisterByName(arr []string) {
	for _, v := range arr {
		switch v {
		case "mail":
			RegisterService("mail", func(app gof.App) {
				AddTickerFunc(startMailQueue)
			})
		case "order":
			RegisterService("order", orderService)
		}
	}
}

// 运行
func Run(ctx gof.App) {
	if ctx != nil {
		appCtx = ctx.(*core.MainApp)
	} else {
		appCtx = core.NewMainApp("app.conf")
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

	appCtx = core.NewMainApp(conf)
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
