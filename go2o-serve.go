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
	"github.com/ixre/gof/web"
	"go.etcd.io/etcd/clientv3"
	"go2o/app"
	"go2o/app/restapi"
	"go2o/core"
	"go2o/core/msq"
	"go2o/core/service"
	"go2o/core/service/impl"
	"log"
	"os"
	"strings"
	"time"
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
		ch            = make(chan bool)
		confFile      string
		etcdEndPoints gof.ArrayFlags
		port          int
		apiPort       int
		mqAddr        string
		debug         bool
		trace         bool
		runDaemon     bool // 运行daemon
		help          bool
		showVer       bool
		newApp        *core.AppImpl
		appFlag       = app.FlagWebApp
	)

	defaultMqAddr := os.Getenv("GO2O_NATS_ADDR")
	if len(defaultMqAddr) == 0 {
		defaultMqAddr = "127.0.0.1:4222"
	}
	flag.IntVar(&port, "port", 1427, "thrift service port")
	flag.IntVar(&apiPort, "apiport", 1428, "api service port")
	flag.Var(&etcdEndPoints, "endpoint", "")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.BoolVar(&trace, "trace", false, "enable trace")
	flag.BoolVar(&help, "help", false, "command usage")
	flag.StringVar(&confFile, "conf", "app.conf", "")
	flag.StringVar(&mqAddr, "mqs", defaultMqAddr,
		"mq cluster address, like: 192.168.1.1:4222,192.168.1.2:4222")
	flag.BoolVar(&runDaemon, "d", false, "run daemon")
	flag.BoolVar(&showVer, "v", false, "print version")
	flag.Parse()
	//confFile = "./app_dev.conf"
	if runDaemon {
		appFlag = appFlag | app.FlagDaemon
	}
	appFlag = appFlag | app.FlagRpcServe
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

	// 默认的ETCD端点
	if len(etcdEndPoints) == 0 {
		etcdEndPoints = strings.Split(os.Getenv("GO2O_ETCD_ADDR"),",")
		if len(etcdEndPoints) == 0  || etcdEndPoints[0] ==""{
			etcdEndPoints = []string{"http://127.0.0.1:2379"}
		}
	}
	cfg := clientv3.Config{
		Endpoints:   etcdEndPoints,
		DialTimeout: 5 * time.Second,
	}

	newApp = core.NewApp(confFile, &cfg)
	if debug {
		go app.AutoInstall()
	}
	if !core.Init(newApp, debug, trace) {
		os.Exit(1)
	}
	go core.SignalNotify(ch, core.AppDispose)
	gof.CurrentApp = newApp
	web.Initialize(web.Options{
		Storage:    newApp.Storage(),
		XSRFCookie: true,
	})
	impl.Init(newApp)
	//runGoMicro()
	// 初始化producer
	_ = msq.Configure(msq.NATS, strings.Split(mqAddr, ","))
	// 运行RPC服务
	go service.ServeRPC(ch, &cfg, port)
	service.ConfigureClient(cfg) // initial service client
	if runDaemon {
		//todo: daemon需重构
		//go daemon.Run(newApp)
	}
	// 运行REST API
	go restapi.Run(ch, newApp, apiPort)
	<-ch
}

/*
// todo: v3 还是测试版本
func runGoMicro() {
	r := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	grpc.NewServer(
		server.Name("Greeter"),
		server.Registry(NewRegisterV3(r)))
	s := service.New(
		service.Name("Greeter"),
		service.Address(":1081"),
		)
	//service := micro.NewService(
	//	micro.Name("Greeter"),
	//	//micro.Address(":1081"),
	//	micro.Registry(r),
	//	)
	//service.Init()
	s.Handle(new(grpc.TestServiceImpl))
	//proto.RegisterGreeterServiceHandler(service,new(grpc.TestServiceImpl))
	service.Run()
}

func NewRegisterV3(r registry.Registry) registry2.Registry {
	return &RegisterV3{r}
}
*/
