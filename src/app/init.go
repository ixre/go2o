/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-16 21:38
 * description :
 * history :
 */

package app

import (
	"github.com/jsix/gof"
	"go2o/src/core/service/dps"
)

func Init(app gof.App) {
	dps.Init(app)
	gof.CurrentApp = app
}
