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
	"github.com/ixre/go2o/core/event/msq"
	"github.com/ixre/go2o/core/infrastructure/locker"
	"github.com/ixre/go2o/core/initial/provide"
	"github.com/ixre/go2o/core/variable"
)

var startJobs = make([]func(), 0)

func Startup(job func()) {
	startJobs = append(startJobs, job)
}

func Init(a *AppImpl, debug, trace bool) bool {
	provide.Configure(a)
	a._debugMode = debug
	// 初始化并发锁
	locker.Configure(a.Storage())
	// 初始化变量
	variable.Domain = a.Config().GetString(variable.ServerDomain)
	a.Loaded = true
	for _, f := range startJobs {
		go f()
	}
	return true
}

func AppDispose() {
	//GetRedisPool().Close()
	msq.Close()
	//if clickhouse.connInstance != nil{
	//	clickhouse.connInstance.Close()
	//}
}
