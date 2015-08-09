/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-01-08 21:35
 * description :
 * history :
 */

package daemon

import (
	"go2o/src/core/service/dps"
	"time"
	"go2o/src/core/variable"
)

func orderDaemon() {

	confirmNewOrder()

	defer recoverDaemon()
	for {
		ids := getPartners()
		for _, v := range ids {
			autoSetOrder(v)
		}
		time.Sleep(time.Minute * CRON_ORDER_SETUP_MINUTE)
	}
}

func autoSetOrder(partnerId int) {
	f := func(err error) {
		appCtx.Log().PrintErr(err)
	}
	dps.ShoppingService.OrderAutoSetup(partnerId, f)
}

func confirmNewOrder(){
	if i, _ := appCtx.Storage().GetInt(variable.KvHaveNewOrder); i > 0 {
		sendQueue()
		appCtx.Storage().Set(variable.KvHaveNewOrder, 0)
	}
	time.Sleep(time.Second * 5)
	confirmNewOrder()
}

type orderInfo struct{
	PartnerId int
	OrderNo string
}
func confirmOrderQueue(){
	var list []*orderInfo = []*orderInfo{}
	appCtx.Db().GetOrm().SelectByQuery(&list,"SELECT order_no,partner_id FROM pt_order WHERE status=")
}
