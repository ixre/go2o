/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-01-08 21:35
 * description :
 * history :
 */

package daemon

import (
	"go2o/core/service/dps"
)

func autoSetOrder(partnerId int) {
	f := func(err error) {
		gCTX.Log().PrintErr(err)
	}
	dps.ShoppingService.OrderAutoSetup(partnerId, f)
}
