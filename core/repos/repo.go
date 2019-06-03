/**
 * Copyright 2015 @ to2.net.
 * name : rep
 * author : jarryliu
 * date : 2016-05-24 10:14
 * description :
 * history :
 */
package repos

import (
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/storage"
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
	//	gof.CurrentApp.Log().Println("[ Go2o][ Repo][ Error] -", err.Error())
	//}
	//return err
}

// 删除指定前缀的缓存
func PrefixDel(sto storage.Interface, prefix string) error {
	rds := sto.(storage.IRedisStorage)
	_, err := rds.DelWith(prefix)
	if err != nil {
		log.Println("[ Cache][ Clean]: clean by prefix ", prefix, " error:", err)
	}
	return err
}
