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
	"go2o/src/core/domain/interface/enum"
	"fmt"
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
	for {
		if i, _ := appCtx.Storage().GetInt(variable.KvHaveNewOrder); i > 0 {
			confirmOrderQueue()
			appCtx.Storage().Set(variable.KvHaveNewOrder, 0)
		}
		time.Sleep(time.Second * 5
	}
}

type orderInfo struct{
	PartnerId int
	OrderNo string
}
func confirmOrderQueue(){
	var list []*orderInfo = []*orderInfo{}
	appCtx.Db().GetOrm().SelectByQuery(&list,fmt.Sprintf("SELECT partner_id,order_no FROM pt_order WHERE status=%d",
		enum.ORDER_CREATED))
	for _,v := range list{
		err := dps.ShoppingService.ConfirmOrder(v.PartnerId,v.OrderNo)
		if err != nil{
			appCtx.Log().Printf("[ DAEMON][ ERROR] - %s\n",err.Error())
		}
	}
}
