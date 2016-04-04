/**
 * Copyright 2015 @ z3q.net.
 * name : go2o-tcpserve.go
 * author : jarryliu
 * date : 2015-11-23 15:52
 * description :
 * history :
 */
package main

import (
	"flag"
	"fmt"
	"github.com/jsix/gof"
	"go2o/src/app/tcpserve"
	"go2o/src/cache"
	"go2o/src/core"
	"go2o/src/core/service/dps"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		port int
		ch   chan bool = make(chan bool)
		logOutput bool
		conf string
	)

	flag.IntVar(&port, "port", 14197, "")
	flag.StringVar(&conf, "conf", "app.conf", "")
	flag.BoolVar(&logOutput,"l",false,"log output")
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Ltime | log.Ldate | log.Lshortfile)

	gof.CurrentApp = core.NewMainApp(conf)
	dps.Init(gof.CurrentApp)
	cache.Initialize(gof.CurrentApp.Storage())

	ts := tcpserve.NewServe(logOutput)
	ts.RegisterJob(tcpserve.MemberSummaryNotifyJob) //注册会员信息通知
	ts.RegisterJob(tcpserve.AccountNotifyJob) //注册账户通知任务
	go ts.Listen(fmt.Sprintf(":%d", port)) //启动服务

	// 检测退出信号
	go func(mainCh chan bool) {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGKILL)
		for {
			switch <-ch {
			case syscall.SIGKILL, syscall.SIGTERM:
				log.Println("[ Tcp][ Term] - tcp serve has term!")
				close(mainCh)
			}
		}
	}(ch)
	log.Println("[ TCP][ SERVE] - socket is serve on port :", port)
	<-ch
}
