package app

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ixre/go2o/internal/core"
	"github.com/ixre/go2o/pkg/app/daemon"
	"github.com/ixre/go2o/pkg/event/events"
	"github.com/ixre/go2o/pkg/event/msq"
	"github.com/ixre/go2o/pkg/infrastructure/etcd"
	"github.com/ixre/go2o/pkg/initial"
	"github.com/ixre/go2o/pkg/initial/bootstrap"
	"github.com/ixre/go2o/pkg/inject"

	"github.com/ixre/go2o/pkg/service"
	"github.com/ixre/gof/domain/eventbus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	etcdEndPoints string
	host          string
	port          int
	mqAddr        string
	debug         bool
	trace         bool
	resetCache    bool // 重置缓存
	runDaemon     bool // 运行daemon
	help          bool
	showVersion   bool
	newApp        *bootstrap.AppConfigLoader
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
	flag.StringVar(&etcdEndPoints, "endpoint", getEtcdAddress(), "etcd endpoints")
	flag.StringVar(&mqAddr, "mqs", getNatsAddress(),
		"nats cluster address, like: 192.168.1.1:4222,192.168.1.2:4222")
	flag.BoolVar(&runDaemon, "d", true, "run daemon")
	flag.IntVar(&port, "port", 1427, "gRPC service port")
	flag.BoolVar(&showVersion, "v", false, "print version")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.BoolVar(&trace, "trace", false, "enable trace")
	flag.BoolVar(&resetCache, "reset-cache", false, "force reset cache")
	flag.BoolVar(&help, "help", false, "command usage")
	flag.Parse()
}

func Run(ch chan bool, confFile string, after func(cfg *clientv3.Config, debug bool)) {
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
	newApp = bootstrap.NewApp(confFile, &cfg)
	if debug {
		go AutoInstall()
	}
	if !initial.Init(newApp, debug, trace) {
		os.Exit(1)
	}
	// 监听信号
	initial.WatchSignals(ch)
	// 初始化第三方配置
	inject.GetSPConfig().Configure()
	// 初始化分布式锁
	etcd.InitializeLocker(&cfg)
	if resetCache {
		// 重置缓存
		initial.ResetCache()
		os.Exit(0)
	}
	// 运行RPC服务
	service.ServeRPC(ch, &cfg, port)
	// 注册服务发现
	service.RegisterServiceDiscovery(&cfg, host, port)
	// 初始化producer
	_ = msq.Configure(msq.NATS,
		strings.Split(mqAddr, ","),
		inject.GetRegistryRepo(),
	)
	// 初始化事件
	inject.GetEventSource().Bind()
	// 发布应用初始化事件
	eventbus.Dispatch(&events.AppInitialEvent{})
	InitialModules()
	if runDaemon {
		go daemon.Run(newApp)
	}
	// 启动后运行
	if after != nil {
		after(&cfg, debug)
	}
	// 运行REST API
	//go restapi.Run(ch, newApp, apiPort)
	<-ch
}
