/**
 * Copyright 2014 @ Ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-16 21:38
 * description :
 * history :
 */

package app

import (
	"github.com/atnet/gof/app"
	"go2o/core/service/dps"
)

func Init(ctx app.Context) {
	dps.Init(ctx)
}
