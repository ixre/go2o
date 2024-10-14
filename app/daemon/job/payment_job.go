package job

import (
	"log"

	"github.com/ixre/go2o/core/infrastructure/locker"
	"github.com/ixre/go2o/core/inject"
)

// 检测已过期的订单并标记
func CheckExpiresPaymentOrder() {
	jobName := "payment.heckExpiresPaymentOrder"

	if !locker.Lock(jobName, 600) {
		return
	}
	defer locker.Unlock(jobName)
	repo := inject.GetPaymentRepo()
	//log.Println("[ job]: start sync wallet log to clickhouse..")
	// job := getJob(jobName)
	// lastId := job.GetValue().LastExecIndex
	size := 200
	lastId := 0
	for {
		list := repo.GetAwaitCloseOrders(lastId, size)
		l := len(list)
		if l > 0 {
			lastId = list[len(list)-1].GetAggregateRootId()
			for _, o := range list {
				//log.Println("取消.", o.GetAggregateRootId())
				if err := o.Cancel(); err != nil {
					log.Println("[ GO2O][ Job]: Cancel expired payment order failed! ", err.Error())
				}
			}
		}
		if l < size {
			break
		}
	}
}
