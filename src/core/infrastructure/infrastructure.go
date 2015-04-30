/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-12-10 16:23
 * description :
 * history :
 */

package infrastructure

import (
	"github.com/atnet/gof"
	"go2o/src/core"
)

//todo:....
var DebugMode bool = false

// get application context
func GetApp() gof.App {
	return core.GlobalApp
}
