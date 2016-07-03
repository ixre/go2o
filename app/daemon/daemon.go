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
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jsix/gof"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"github.com/robfig/cron"
	"go2o/core"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/service/dps"
	"go2o/core/variable"
	"log"
	"strings"
	"time"
)

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

	// 处理订单,需根据订单不同的状态,作不同的业务,返回布尔值,如果返回false,则不继续执行
	OrderObs(*order.ValueOrder) bool

	// 监视会员修改,@create:是否为新注册会员,返回布尔值,如果返回false,则不继续执行
	MemberObs(m *member.Member, create bool) bool

	// 通知支付单完成队列,返回布尔值,如果返回false,则不继续执行
	PaymentOrderObs(order *payment.PaymentOrderBean) bool

	// 处理邮件队列,返回布尔值,如果返回false,则不继续执行
	HandleMailQueue([]*mss.MailTask) bool
}

var (
	appCtx           *core.MainApp
	_db              db.Connector
	_orm             orm.Orm
	services         []Service      = make([]Service, 0)
	serviceNames     map[string]int = make(map[string]int)
	tickerDuration   time.Duration  = 20 * time.Second // 间隔20秒执行
	tickerInvokeFunc []Func         = []Func{}
	cronTab          *cron.Cron     = cron.New()
	ticker           *time.Ticker   = time.NewTicker(tickerDuration)
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
	defer func() {
		cronTab.Stop()
		ticker.Stop()
	}()

	for i, s := range services { //运行自定义服务
		log.Println("** [ Go2o][ Daemon] - (", i, ")", s.Name(), "daemon running")
		s.SetApp(appCtx)
		go s.Start()
	}

	startCronTab()
	startTicker() //阻塞

}

func startTicker() {
	for { //执行定时任务
		select {
		case <-ticker.C:
			for _, f := range tickerInvokeFunc {
				go f(appCtx)
			}
		}
	}
}

func startCronTab() {
	//cron
	cronTab.AddFunc("0 0 1 * * *", personFinanceSettle) //个人金融结算,每天2点更新数据
	//cronTab.AddFunc("1 * * * * *", func() { log.Println("grouting -", runtime.NumGoroutine(), runtime.NumCPU()) })
	cronTab.Start()

	go func() {
		personFinanceSettle()
	}()
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
func (this *defaultService) init() {
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
	go supervisePaymentOrderFinish(services)
	go startMailQueue(services)
}

// 处理订单,需根据订单不同的状态,作不同的业务
// 返回布尔值,如果返回false,则不继续执行
func (this *defaultService) OrderObs(o *order.ValueOrder) bool {
	defer Recover()
	conn := core.GetRedisConn()
	defer conn.Close()
	if this.sOrder {
		if o.Status == enum.ORDER_WAIT_CONFIRM { //确认订单
			dps.ShoppingService.ConfirmOrder(o.OrderNo)
		}
	}
	return true
}

// 监视会员修改,@create:是否为新注册会员
// 返回布尔值,如果返回false,则不继续执行
func (this *defaultService) MemberObs(m *member.Member, create bool) bool {
	defer Recover()
	if this.sMember {
		//todo: 执行会员逻辑
	}
	return true
}

// 通知支付单完成队列,返回布尔值,如果返回false,则不继续执行
func (this *defaultService) PaymentOrderObs(order *payment.PaymentOrderBean) bool {
	this.app.Log().Println("---支付单", order.TradeNo, "支付完成")
	return true
}

//设置订单过期时间
func (this *defaultService) setOrderExpires(conn redis.Conn, o *order.ValueOrder) {
	if o.Status == enum.ORDER_WAIT_PAYMENT { //订单刚创建时,设置过期时间
		ss := dps.MerchantService.GetSaleConf(o.VendorId)
		t := int64(ss.OrderTimeOutMinute) * 60
		unix := o.CreateTime + t
		conn.Do("SET", this.getExpiresKey(o), unix)
	} else if o.IsPaid == 1 { //删除过期时间
		conn.Do("DEL", this.getExpiresKey(o))
	}
}

func (this *defaultService) getExpiresKey(o *order.ValueOrder) string {
	return fmt.Sprintf("%s%s", variable.KvOrderExpiresTime, o.OrderNo)
}

// 处理邮件队列
// 返回布尔值,如果返回false,则不继续执行
func (this *defaultService) HandleMailQueue(list []*mss.MailTask) bool {
	defer Recover()
	if !this.sMail {
		handleMailQueue(list)
	}
	return true
}

// 运行
func Run(ctx gof.App) {
	if ctx != nil {
		appCtx = ctx.(*core.MainApp)
	} else {
		appCtx = core.NewMainApp("app.conf")
	}
	_db = appCtx.Db()
	_orm = _db.GetOrm()
	sMail := appCtx.Config().GetString(variable.SystemMailQueueOff) != "1" //是否关闭系统邮件队列
	//sMail := cnf.GetString(variable.)

	s := &defaultService{
		sMember: true,
		sOrder:  true,
		sMail:   sMail,
	}
	s.init()
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

	appCtx = core.NewMainApp(conf)
	appCtx.Init(debug, trace)
	gof.CurrentApp = appCtx

	_db = appCtx.Db()
	_orm = _db.GetOrm()

	dps.Init(appCtx)

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
	s.init()
	Start()

	<-ch
}
