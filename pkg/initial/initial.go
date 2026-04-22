/**
 * Copyright 2015 @ 56x.net.
 * name : types.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package initial

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ixre/go2o/internal/impl/repo"
	"github.com/ixre/go2o/pkg/constants"
	"github.com/ixre/go2o/pkg/event/msq"
	"github.com/ixre/go2o/pkg/infra/logger"
	"github.com/ixre/go2o/pkg/initial/bootstrap"
	"github.com/ixre/go2o/pkg/initial/provide"
)

var startJobs = make([]func(), 0)

var _appInstance *bootstrap.AppConfigLoader

func Startup(job func()) {
	startJobs = append(startJobs, job)
}

func Init(a *bootstrap.AppConfigLoader, debug, trace bool) bool {
	_appInstance = a
	provide.Configure(a)
	//a._debugMode = debug
	// 初始化变量
	constants.Domain = a.Config().GetString(constants.ConfigServerDomain)
	a.Loaded = true
	for _, f := range startJobs {
		go f()
	}

	repo.OrmMapping(provide.GetOrmInstance())

	return true
}

// ResetCache 重置缓存
func ResetCache() {
	sto := provide.GetStorageInstance()
	prefixs := []string{
		"go2o",
		"member",
		"merchant",
	}
	var total int
	for _, v := range prefixs {
		i, err := sto.DeleteWith(v)
		if err != nil {
			logger.Error("reset cache error, %s", err.Error())
			return
		}
		total += i
	}
	logger.Info("reset cache complete! total %s keys", total)
}

func appDispose() {
	//GetRedisPool().Close()
	msq.Close()
	//if clickhouse.connInstance != nil{
	//	clickhouse.connInstance.Close()
	//}
}

func WatchSignals(c chan bool) {
	go watchSignals(c, appDispose)
}

// watchSignals 监听进程信号,并执行操作。比如退出时应释放资源
func watchSignals(c chan bool, fn func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM)
	for {
		switch <-ch {
		case syscall.SIGHUP, syscall.SIGTERM: // 退出时
			log.Println("[ OS][ TERM] - program has exit !")
			fn()
			close(c)
		}
	}
}
