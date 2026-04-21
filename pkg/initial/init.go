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
	"github.com/ixre/go2o/pkg/event/msq"
	"github.com/ixre/go2o/pkg/infrastructure/logger"
	"github.com/ixre/go2o/pkg/initial/bootstrap"
	"github.com/ixre/go2o/pkg/initial/provide"
	"github.com/ixre/go2o/pkg/variable"
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
