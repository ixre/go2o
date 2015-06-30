/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-10 16:23
 * description :
 * history :
 */

package infrastructure

import (
	"github.com/atnet/gof"
)

//todo:....
var DebugMode bool = false

// get application context
func GetApp() gof.App {
	return gof.CurrentApp
}
