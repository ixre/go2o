package app

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ixre/go2o/app/daemon"
	"github.com/ixre/go2o/core"
	"github.com/ixre/go2o/core/etcd"
	"github.com/ixre/go2o/core/event/msq"
	"github.com/ixre/go2o/core/initial"
	"github.com/ixre/go2o/core/service"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	confFile      string
	etcdEndPoints string
	host          string
	port          int
	//apiPort       int
	mqAddr      string
	debug       bool
	trace       bool
	runDaemon   bool // 运行daemon
	help        bool
	showVersion bool
	newApp      *initial.AppImpl
)

func getNatsAddress() string {
	defaultMqAddr := os.Getenv("GO2O_NATS_ADDR")
	if len(defaultMqAddr) == 0 {
		defaultMqAddr = "127.0.0.1:4222"
	}
	return defaultMqAddr
}
func getEtcdAddress() string {
	defaultEtcd := os.Getenv("GO2O_ETCD_ADDR")
	if len(defaultEtcd) == 0 {
		defaultEtcd = "127.0.0.1:2379"
	}
	return defaultEtcd
}

func ParseFlags() {
	flag.StringVar(&confFile, "conf", "app.conf", "")
	flag.StringVar(&etcdEndPoints, "endpoint", getEtcdAddress(), "etcd endpoints")
	flag.StringVar(&mqAddr, "mqs", getNatsAddress(),
		"nats cluster address, like: 192.168.1.1:4222,192.168.1.2:4222")
	flag.BoolVar(&runDaemon, "d", true, "run daemon")
	flag.IntVar(&port, "port", 1427, "gRPC service port")
	//flag.IntVar(&apiPort, "apiport", 1428, "api service port")
	flag.BoolVar(&showVersion, "v", false, "print version")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.BoolVar(&trace, "trace", false, "enable trace")
	flag.BoolVar(&help, "help", false, "command usage")
	flag.Parse()
}

func Run(ch chan bool, after func(*clientv3.Config)) {
	if help {
		flag.Usage()
		os.Exit(0)
		return
	}
	if showVersion {
		fmt.Printf("go2o version v%s \n", core.Version)
		return
	}
	if len(host) == 0 {
		host = os.Getenv("GO2O_SERVER_HOST")
	}
	//confFile = "./app_dev.conf"
	// if runDaemon {
	// 	appFlag = appFlag | FlagDaemon
	// }
	// appFlag = appFlag | FlagRpcServe
	// setting log flags
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Ltime | log.Ldate)

	etcdEndPoints := strings.Split(etcdEndPoints, ",")
	if len(etcdEndPoints) == 0 {
		log.Fatalln("etcd endpoints not specified")
	}

	cfg := clientv3.Config{
		Endpoints:   etcdEndPoints,
		DialTimeout: 5 * time.Second,
	}
	newApp = initial.NewApp(confFile, &cfg)
	if debug {
		go AutoInstall()
	}
	if !initial.Init1(newApp, debug, trace) {
		os.Exit(1)
	}
	go core.SignalNotify(ch, initial.AppDispose)
	//impl.Init(newApp)
	//runGoMicro()
	// 初始化分布式锁
	etcd.InitializeLocker(&cfg)
	// 运行RPC服务
	service.ServeRPC(ch, &cfg, port)
	service.RegisterServiceDiscovery(&cfg, host, port)
	// 初始化producer
	_ = msq.Configure(msq.NATS, strings.Split(mqAddr, ","))
	// initial service client
	//service.ConfigureClient(&cfg, "")
	if runDaemon {
		go daemon.Run(newApp)
	}
	// 启动后运行
	if after != nil {
		after(&cfg)
	}
	// 运行REST API
	//go restapi.Run(ch, newApp, apiPort)
	<-ch
}
