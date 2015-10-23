/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-01-08 21:35
 * description :
 * history :
 */

package daemon

import (
	"fmt"
	"github.com/jsix/gof"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"time"
)

func orderDaemon(app gof.App) {
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

func confirmNewOrder(app gof.App, dfs []DaemonFunc) {
	if i, _ := appCtx.Storage().GetInt(variable.KvHaveNewCreatedOrder); i == enum.TRUE {
		appCtx.Log().Printf("[ DAEMON][ ORDER][ CONFIRM] - begin invoke confirm handler.")
		if dfs == nil || len(dfs) == 0 {
			confirmOrderQueue(app)
		} else {
			for _, v := range dfs {
				v(app)
			}
		}
		appCtx.Storage().Set(variable.KvHaveNewCreatedOrder, enum.FALSE)
	}
}

func completedOrderObs(app gof.App, dfs []DaemonFunc) {
	if i, _ := appCtx.Storage().GetInt(variable.KvHaveNewCompletedOrder); i == enum.TRUE {
		appCtx.Log().Printf("[ DAEMON][ ORDER][ FINISHED] - begin invoke finish handler.\n")
		for _, v := range dfs {
			v(app)
		}
		appCtx.Storage().Set(variable.KvHaveNewCompletedOrder, enum.FALSE)
	}
}

type orderInfo struct {
	PartnerId int
	OrderNo   string
}

func confirmOrderQueue(app gof.App) {
	var list []*orderInfo = []*orderInfo{}
	appCtx.Db().GetOrm().SelectByQuery(&list, fmt.Sprintf("SELECT partner_id,order_no FROM pt_order WHERE status=%d",
		enum.ORDER_WAIT_CONFIRM))
	for _, v := range list {
		err := dps.ShoppingService.ConfirmOrder(v.PartnerId, v.OrderNo)
		if err != nil {
			appCtx.Log().Printf("[ DAEMON][ ORDER][ ERROR] - %s\n", err.Error())
		}
	}
}
