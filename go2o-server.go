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
	"github.com/jsix/gof"
	"github.com/jsix/gof/storage"
	"go2o/src/app"
	"go2o/src/app/daemon"
	"go2o/src/cache"
	"go2o/src/core"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

func main() {
	var (
		ch         chan bool = make(chan bool)
		confFile   string
		httpPort   int
		socketPort int
		restPort   int
		mode       string //启动模式: h开启http,s开启socket,a开启所有
		debug      bool
		trace      bool
		runDaemon  bool // 运行daemon
		help       bool
		newApp     *core.MainApp
	)

	flag.IntVar(&socketPort, "port2", 1001, "socket server port")
	flag.IntVar(&httpPort, "port", 1002, "web server port")
	flag.IntVar(&restPort, "port3", 1003, "rest api port")
	flag.StringVar(&mode, "mode", "shr", "boot mode.'h'- boot http service,'s'- boot socket service")
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

	runtime.GOMAXPROCS(runtime.NumCPU())

	newApp = core.NewMainApp(confFile)
	if !newApp.Init(debug, trace) {
		os.Exit(1)
	}

	go handleSignal(ch)

	if v := newApp.Config().GetInt("server_port"); v != 0 {
		httpPort = v
	}
	if v := newApp.Config().GetInt("socket_port"); v != 0 {
		socketPort = v
	}
	if v := newApp.Config().GetInt("api_service_port"); v != 0 {
		restPort = v
	}

	gof.CurrentApp = newApp
	app.Init(newApp)
	cache.Initialize(storage.NewRedisStorage(newApp.Redis()))
	core.RegisterTypes()

	var booted bool

	if runDaemon {
		go daemon.Run(newApp)
	}

	if strings.Contains(mode, "s") {
		booted = true
		go app.RunSocket(newApp, socketPort, debug, trace)
	}

	if strings.Contains(mode, "h") {
		booted = true
		go app.RunWeb(newApp, httpPort, debug, trace)
	}

	if strings.Contains(mode, "r") {
		booted = true
		go app.RunRestApi(newApp, restPort)
	}

	if booted {
		<-ch
	}
}

func handleSignal(srcCh chan bool) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGHUP:
		//log.Println("[ OS][ TERM] - go2o sighup ...")
		case syscall.SIGTERM: // 退出时
			log.Println("[ OS][ TERM] - go2o server has exit !")
			close(srcCh)
		}
	}
}
