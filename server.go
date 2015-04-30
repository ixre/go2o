/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-12-16 21:45
 * description :
 * history :
 */

package main

import (
	"flag"
	"github.com/atnet/gof"
	"go2o/src/app"
	"go2o/src/app/daemon"
	"go2o/src/core"
	"os"
	"runtime"
	"strings"
)

func main() {
	var (
		ch         chan bool = make(chan bool)
		httpPort   int
		socketPort int
		mode       string //启动模式: h开启http,s开启socket,a开启所有
		debug      bool
		trace      bool
		help       bool
		newApp     gof.App
	)

	newApp = &core.MainApp{}
	flag.IntVar(&httpPort, "port", 8091, "web server port")
	flag.IntVar(&socketPort, "port2", 8003, "socket server port")
	flag.StringVar(&mode, "mode", "sh", "boot mode.'h'- boot http service,'s'- boot socket service")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.BoolVar(&trace, "trace", false, "enable trace")
	flag.BoolVar(&help, "help", false, "command usage")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if gcx, ok := newApp.(*core.MainApp); !ok || !gcx.Init(debug, trace) {
		os.Exit(1)
	}

	core.SetGlobalApp(newApp)
	app.Init(newApp)
	daemon.Run(newApp)

	var booted bool
	if strings.Contains(mode, "s") {
		booted = true
		go app.RunSocket(newApp, socketPort, debug, trace)
	}

	if strings.Contains(mode, "h") {
		booted = true
		go app.RunWeb(newApp, httpPort, debug, trace)
	}

	if booted {
		<-ch
	}
}
