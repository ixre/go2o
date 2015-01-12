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
	"com/ording/dproxy"
	"github.com/atnet/gof/app"
)

func Init(ctx app.Context) {
	dproxy.Init(ctx)
}
