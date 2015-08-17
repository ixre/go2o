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
	"fmt"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"time"
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

func confirmNewOrder() {
	for {
		if i, _ := appCtx.Storage().GetInt(variable.KvHaveNewOrder); i == enum.FALSE {
			appCtx.Log().Printf("[ DAEMON][ ORDER][ CONFIRM] - begin confirm")
			confirmOrderQueue()
			appCtx.Storage().Set(variable.KvHaveNewOrder, enum.TRUE)
		}
		time.Sleep(time.Second * 5)
	}
}

type orderInfo struct {
	PartnerId int
	OrderNo   string
}

func confirmOrderQueue() {
	var list []*orderInfo = []*orderInfo{}
	appCtx.Db().GetOrm().SelectByQuery(&list, fmt.Sprintf("SELECT partner_id,order_no FROM pt_order WHERE status=%d",
		enum.ORDER_CREATED))
	for _, v := range list {
		err := dps.ShoppingService.ConfirmOrder(v.PartnerId, v.OrderNo)
		if err != nil {
			appCtx.Log().Printf("[ DAEMON][ ORDER][ ERROR] - %s\n", err.Error())
		}
	}
}
