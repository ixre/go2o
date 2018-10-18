/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-16 21:45
 * description :
 * history :
 */

package main

import (
	"flag"
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/storage"
	"github.com/jsix/gof/web/session"
	"go2o/src/app"
	"go2o/src/app/cache"
	"go2o/src/app/daemon"
	"go2o/src/app/restapi"
	"go2o/src/core"
	"go2o/src/core/service/dps"
	"log"
	"os"
	"runtime"
	"go2o/src/fix"
)

func main() {
	var (
		ch        chan bool = make(chan bool)
		confFile  string
		httpPort  int
		restPort  int
		debug     bool
		trace     bool
		runDaemon bool // 运行daemon
		help      bool
		newApp    *core.MainApp
	)

	flag.IntVar(&httpPort, "port", 14190, "web server port")
	flag.IntVar(&restPort, "restport", 14191, "rest api port")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.BoolVar(&trace, "trace", false, "enable trace")
	flag.BoolVar(&help, "help", false, "command usage")
	flag.StringVar(&confFile, "conf", "app.conf", "")
	flag.BoolVar(&runDaemon, "d", false, "run daemon")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Ltime | log.Ldate | log.Lshortfile)

	runtime.GOMAXPROCS(runtime.NumCPU())
	newApp = core.NewMainApp(confFile)
	if !newApp.Init(debug, trace) {
		os.Exit(1)
	}
	fix.CustomFix()
	go fix.SignalNotify(ch)

	if v := newApp.Config().GetInt("server_port"); v != 0 {
		httpPort = v
	}
	if v := newApp.Config().GetInt("api_service_port"); v != 0 {
		restPort = v
	}
	gof.CurrentApp = newApp
	dps.Init(newApp)
	cache.Initialize(storage.NewRedisStorage(newApp.Redis()))
	core.RegisterTypes()
	session.Initialize(newApp.Storage(), "",false)
	
	if runDaemon {
		go daemon.Run(newApp)
	}

	go app.Run(ch, newApp, fmt.Sprintf(":%d", httpPort)) //运行HTTP

	go restapi.Run(newApp, restPort) // 运行REST API

	<-ch

	os.Exit(1) // 退出
}
