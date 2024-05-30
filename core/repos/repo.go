/**
 * Copyright 2015 @ 56x.net.
 * name : rep
 * author : jarryliu
 * date : 2016-05-24 10:14
 * description :
 * history :
 */
package repos

import (
	"sync"

	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/storage"
)

var (
	mux                 sync.Mutex
	DefaultCacheSeconds int64 = 3600
)

// 处理错误
func handleError(err error) error {
	return domain.HandleError(err, "rep")
	//if err != nil && gof.CurrentApp.Debug() {
	//	gof.CurrentApp.Log().Println("[ GO2O][ Repo][ Error] -", err.Error())
	//}
	//return err
}

// 删除指定前缀的缓存
func PrefixDel(sto storage.Interface, prefix string) error {
	_, err := sto.DeleteWith(prefix)
	if err != nil {
		log.Println("[ Cache][ Clean]: clean by prefix ", prefix, " error:", err)
	}
	return err
}


var (
	DB  db.Connector
	ORM orm.Orm
	STO storage.Interface
)

func Initial(o orm.Orm, sto storage.Interface) {
	ORM = o
	STO = sto
	DB = o.Connector()
	// // 解决依赖
	// r.orderRepo.(*OrderRepoImpl).SetPaymentRepo(r.paymentRepo)
	// // 初始化数据
	// r.memberRepo.GetManager().GetAllBuyerGroups()
}

// 返回orm实例
func GetOrmInstance() orm.Orm {
	return ORM
}

func GetStorageInstance() storage.Interface {
	return STO
}
