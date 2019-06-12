/**
 * Copyright 2014 @ to2.net.
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
	"github.com/ixre/gof"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/web"
	"go2o/app"
	"go2o/app/cache"
	"go2o/app/daemon"
	"go2o/app/restapi"
	"go2o/core"
	"go2o/core/msq"
	"go2o/core/service/rsi"
	rs "go2o/core/service/thrift/service"
	"log"
	"os"
	"runtime"
	"strings"
)

var _ = `

  ###   ###   ###   ###
 #     #  ##    #  #  ##
#     #    #    # #    #
#  #  #   #   ##  #   #
#  #  #   #  #    #   #
 ###   ###   ###   ###


Go2o is Google Go language binding domain-driven design (DDD) O2O open source implementation. Support Online Store
, Offline stores; multi-channel (businesses), multi-store, merchandise, snapshots, orders, sales, payment, distribution and other functions.

Project by a management center (including platform management center, business background, store background), online store (PC shop,
Handheld shops, micro-channel), the member center, open API in four parts.

Go2o using domain-driven design for business depth abstract, theoretical support in most sectors O2O scenarios.
Through open API, you can seamlessly integrate into legacy systems.


Email: jarrysix#gmail.com

`

func main() {
	var (
		ch        = make(chan bool)
		confFile  string
		httpPort  int
		restPort  int
		debug     bool
		trace     bool
		runDaemon bool // 运行daemon
		runRpc    bool // 运行Rpc服务
		help      bool
		showVer   bool
		newApp    *core.AppImpl
		appFlag   = app.FlagWebApp
	)

	flag.IntVar(&httpPort, "port", 14190, "web server port")
	flag.IntVar(&restPort, "restport", 1428, "rest api port")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.BoolVar(&trace, "trace", false, "enable trace")
	flag.BoolVar(&help, "help", false, "command usage")
	flag.StringVar(&confFile, "conf", "app.conf", "")
	flag.BoolVar(&runDaemon, "d", false, "run daemon")
	flag.BoolVar(&runRpc, "r", false, "run rpc service")
	flag.BoolVar(&showVer, "v", false, "print version")
	flag.Parse()

	if runDaemon {
		appFlag = appFlag | app.FlagDaemon
	}
	if runRpc {
		appFlag = appFlag | app.FlagRpcServe
	}
	if help {
		flag.Usage()
		return
	}
	if showVer {
		fmt.Println(fmt.Sprintf("go2o version v%s", core.Version))
		return
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Ltime | log.Ldate | log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())
	newApp = core.NewApp(confFile)
	if debug {
		go app.AutoInstall()
	}
	if !core.Init(newApp, debug, trace) {
		os.Exit(1)
	}
	go core.SignalNotify(ch)
	gof.CurrentApp = newApp
	cache.Initialize(storage.NewRedisStorage(newApp.Redis()))
	web.Initialize(web.Options{
		Storage:    newApp.Storage(),
		XSRFCookie: true,
	})
	app.FsInit(debug)
	rsi.Init(newApp, appFlag)
	//app.Configure(hook.HookUp, newApp, appFlag)

	// 初始化producer
	kafkaAddress := strings.Split(newApp.Config().GetString("kafka_address"), ",")
	msq.Configure(msq.KAFKA, kafkaAddress)
	if runRpc {
		go rs.ListenAndServe(":1427", false)
	}
	if runDaemon {
		go daemon.Run(newApp)
	}
	go restapi.Run(newApp, restPort) // 运行REST API
	<-ch
}
