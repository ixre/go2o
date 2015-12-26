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
		conf string
	)

	flag.IntVar(&port, "port", 14197, "")
	flag.StringVar(&conf, "conf", "app.conf", "")
	flag.Parse()

	gof.CurrentApp = core.NewMainApp(conf)
	dps.Init(gof.CurrentApp)
	cache.Initialize(gof.CurrentApp.Storage())
	tcpserve.DebugOn = true
	go tcpserve.ListenTcp(fmt.Sprintf(":%d", port))
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
