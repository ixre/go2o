/**
 * Copyright 2015 @ z3q.net.
 * name : rep
 * author : jarryliu
 * date : 2016-05-24 10:14
 * description :
 * history :
 */
package repository

import (
	"go2o/core/infrastructure/domain"
	"sync"
)

var (
	mux                 sync.Mutex
	DefaultCacheSeconds int64 = 3600
)

// 处理错误
func handleError(err error) error {
	return domain.HandleError(err, "rep")
	//if err != nil && gof.CurrentApp.Debug() {
	//	gof.CurrentApp.Log().Println("[ Go2o][ Rep][ Error] -", err.Error())
	//}
	//return err
}
