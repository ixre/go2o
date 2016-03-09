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
<<<<<<< HEAD
	"go2o/src/app/tcpserve"
	"go2o/src/cache"
	"go2o/src/core"
	"go2o/src/core/service/dps"
=======
	"go2o/src/cache"
	"go2o/src/core"
	"go2o/src/core/service/dps"
	"go2o/src/tcpserve"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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

<<<<<<< HEAD
	flag.IntVar(&port, "port", 14197, "")
	flag.StringVar(&conf, "conf", "app.conf", "")
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Ltime | log.Ldate | log.Lshortfile)

	gof.CurrentApp = core.NewMainApp(conf)
	dps.Init(gof.CurrentApp)
	cache.Initialize(gof.CurrentApp.Storage())
=======
	flag.IntVar(&port, "port", 1005, "")
	flag.StringVar(&conf, "conf", "app.conf", "")
	flag.Parse()

	gof.CurrentApp = core.NewMainApp(conf)
	dps.Init(gof.CurrentApp)
	cache.Initialize(gof.CurrentApp.Storage())
	tcpserve.DebugOn = true
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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

<<<<<<< HEAD
	log.Println("[ TCP][ SERVE] - socket is serve on port :", port)
=======
	log.Println("[ TCP][ SERVE] - socket is served ... ")
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d

	<-ch
}
