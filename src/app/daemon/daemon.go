/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-01-08 21:01
 * description :
 * history :
 */

package daemon

import (
	"flag"
<<<<<<< HEAD
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jsix/gof"
	"go2o/src/core"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner/mss"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
=======
	"github.com/jsix/gof"
	"go2o/src/core"
	"go2o/src/core/service/dps"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"log"
	"strings"
	"time"
)

<<<<<<< HEAD
// 守护进程执行的函数
type Func func(gof.App)

// 守护进程服务
type Service interface {
	// 服务名称
	Name() string
	// 设置APP上下文
	SetApp(gof.App)
	// 启动服务
	Start()
	// 处理订单,需根据订单不同的状态,作不同的业务
	// 返回布尔值,如果返回false,则不继续执行
	OrderObs(*shopping.ValueOrder) bool
	// 监视会员修改,@create:是否为新注册会员
	// 返回布尔值,如果返回false,则不继续执行
	MemberObs(m *member.ValueMember, create bool) bool
	// 处理邮件队列
	// 返回布尔值,如果返回false,则不继续执行
	HandleMailQueue([]*mss.MailTask) bool
}

var (
	appCtx           *core.MainApp
	services         []Service      = make([]Service, 0)
	serviceNames     map[string]int = make(map[string]int)
	tickerDuration   time.Duration  = 5 * time.Second // 间隔5秒执行
	tickerInvokeFunc []Func         = []Func{}
)

// 注册服务
func RegisterService(s Service) {
	if s == nil {
		panic("service is nil")
	}
	name := s.Name()
	if _, ok := serviceNames[name]; ok {
		panic("service named " + name + " is registed!")
	}
	serviceNames[name] = len(services)
	services = append(services, s)
}

// 添加定时执行任务(默认5秒)
func AddTickerFunc(f Func) {
	tickerInvokeFunc = append(tickerInvokeFunc, f)
}

// 启动守护进程
func Start() {
	for i, s := range services { //运行自定义服务
		log.Println("** [ Go2o][ Daemon] - (", i, ")", s.Name(), "daemon running")
		s.SetApp(appCtx)
		go s.Start()
	}
	tk := time.NewTicker(tickerDuration)
	defer func() {
		tk.Stop()
	}()
	for { //执行定时任务
		select {
		case <-tk.C:
			for _, f := range tickerInvokeFunc {
				f(appCtx)
			}
		}
	}
}

func recoverDaemon() {

}

type defaultService struct {
	app     gof.App
	sOrder  bool
	sMember bool
	sMail   bool
}

// 注册系统服务
func (this *defaultService) register() {
	if len(services) == 0 {
		RegisterService(this)
	} else {
		services = append([]Service{this}, services...)
	}
}

// 服务名称
func (this *defaultService) Name() string {
	return "sys"
}

// 设置APP上下文
func (this *defaultService) SetApp(a gof.App) {
	this.app = a
}

// 启动服务
func (this *defaultService) Start() {
	//AddTickerFunc(orderDaemon) //订单自动进行流程
	AddTickerFunc(detectOrderExpires) //检查订单过期
	go superviseMemberUpdate(services)
	go superviseOrder(services)
	go startMailQueue(services)
}

// 处理订单,需根据订单不同的状态,作不同的业务
// 返回布尔值,如果返回false,则不继续执行
func (this *defaultService) OrderObs(o *shopping.ValueOrder) bool {
	defer Recover()
	conn := core.GetRedisConn()
	defer conn.Close()
	if this.sOrder {
		if o.Status == enum.ORDER_WAIT_CONFIRM { //确认订单
			dps.ShoppingService.ConfirmOrder(o.PartnerId, o.OrderNo)
		}
	}
	return true
}

//设置过期时间
func (this *defaultService) setOrderExpires(conn redis.Conn, o *shopping.ValueOrder) {
	if o.Status == enum.ORDER_WAIT_PAYMENT { //订单刚创建时,设置过期时间
		ss := dps.PartnerService.GetSaleConf(o.PartnerId)
		t := int64(ss.OrderTimeOutMinute) * 60
		unix := o.CreateTime + t
		conn.Do("SET", this.getExpiresKey(o), unix)
	} else if o.IsPaid == 1 { //删除过期时间
		conn.Do("DEL", this.getExpiresKey(o))
	}
}

func (this *defaultService) getExpiresKey(o *shopping.ValueOrder) string {
	return fmt.Sprintf("%s%d_%s", variable.KvOrderExpiresTime, o.PartnerId, o.OrderNo)
}

// 监视会员修改,@create:是否为新注册会员
// 返回布尔值,如果返回false,则不继续执行
func (this *defaultService) MemberObs(m *member.ValueMember, create bool) bool {
	defer Recover()
	if this.sMember {
		//todo: 执行会员逻辑
	}
	return true
}

// 处理邮件队列
// 返回布尔值,如果返回false,则不继续执行
func (this *defaultService) HandleMailQueue(list []*mss.MailTask) bool {
	defer Recover()
	if !this.sMail {
		handleMailQueue(list)
	}
	return true
=======
// 守护进程服务
type DaemonService func(gof.App)
type DaemonFunc func(gof.App)

var (
	appCtx                 *core.MainApp
	services               map[string]DaemonService = map[string]DaemonService{}
	tickerDuration                                  = 5 * time.Second // 间隔5秒执行
	tickerInvokeFunc       []DaemonFunc             = []DaemonFunc{}
	newOrderObserver       []DaemonFunc             = []DaemonFunc{confirmOrderQueue}
	completedOrderObserver []DaemonFunc             = []DaemonFunc{}

//newMemberObserver []DaemonFunc = []DaemonFunc{orderDaemon}
)

func RegisterService(name string, service DaemonService) {
	if _, ok := services[name]; ok {
		panic("service named " + name + " is registed!")
	}
	services[name] = service
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}

// 运行
func Run(ctx gof.App) {
	if ctx != nil {
		appCtx = ctx.(*core.MainApp)
	} else {
<<<<<<< HEAD
		appCtx = core.NewMainApp("app.conf")
	}

	sMail := appCtx.Config().GetString(variable.SystemMailQueueOff) != "1" //是否关闭系统邮件队列
	//sMail := cnf.GetString(variable.)

	s := &defaultService{
		sMember: true,
		sOrder:  true,
		sMail:   sMail,
	}
	s.register()
=======
		appCtx = getAppCtx("app.conf")
	}
	RegisterByName([]string{"mail", "order"})
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	Start()
}

// 自定义参数运行
func FlagRun() {
	var conf string
	var debug bool
	var trace bool
	var service string
	var serviceArr []string = []string{"mail", "order"}
	var ch chan bool = make(chan bool)
	flag.StringVar(&conf, "conf", "app.conf", "")
	flag.BoolVar(&debug, "debug", true, "")
	flag.BoolVar(&trace, "trace", true, "")
	flag.StringVar(&service, "service", strings.Join(serviceArr, ","), "")

	flag.Parse()

<<<<<<< HEAD
	appCtx = core.NewMainApp(conf)
=======
	appCtx = getAppCtx(conf)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	appCtx.Init(debug, trace)
	gof.CurrentApp = appCtx

	dps.Init(appCtx)

<<<<<<< HEAD
	//todo:???
	//	if service != "all" {
	//		serviceArr = strings.Split(service, ",")
	//	}
	// RegisterByName(serviceArr)

	s := &defaultService{
		sMember: true,
		sOrder:  true,
		sMail:   true,
	}
	s.register()
=======
	if service != "all" {
		serviceArr = strings.Split(service, ",")
	}

	RegisterByName(serviceArr)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	Start()

	<-ch
}
<<<<<<< HEAD
=======

func getAppCtx(conf string) *core.MainApp {
	return core.NewMainApp(conf)
}

// 添加定时执行任务(默认5秒)
func AddTickerFunc(f DaemonFunc) {
	tickerInvokeFunc = append(tickerInvokeFunc, f)
}

// 获取订单处理函数
func orderDaemonService(app gof.App) {
	AddTickerFunc(func(app gof.App) {
		confirmNewOrder(app, newOrderObserver)

		if completedOrderObserver != nil && len(completedOrderObserver) != 0 {
			completedOrderObs(app, completedOrderObserver)
		}
	})
	orderDaemon(app)
}

// 添加新的订单处理函数
func AddNewOrderFunc(f DaemonFunc) {
	newOrderObserver = append(newOrderObserver, f)
}

// 添加已完成订单处理函数
func AddCompletedOrderFunc(f DaemonFunc) {
	completedOrderObserver = append(completedOrderObserver, f)
}

func RegisterByName(arr []string) {
	for _, v := range arr {
		switch v {
		case "mail":
			RegisterService("mail", func(app gof.App) {
				AddTickerFunc(startMailQueue)
			})
		case "order":
			RegisterService("order", orderDaemonService)
		}
	}
}

func Start() {
	tk := time.NewTicker(tickerDuration)
	defer func() {
		tk.Stop()
	}()

	for name, s := range services {
		log.Println("** [ Go2o][ Daemon][ Booted] - ", name, " daemon running")
		go s(appCtx)
	}

	for {
		select {
		case <-tk.C:
			for _, f := range tickerInvokeFunc {
				f(appCtx)
			}
		}
	}
}

func recoverDaemon() {

}
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
