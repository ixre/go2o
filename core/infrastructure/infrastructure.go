/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-10 16:23
 * description :
 * history :
 */

package infrastructure

import (
	"github.com/ixre/gof"
)

//todo:....
var DebugMode = false

// get application context
func GetApp() gof.App {
	return gof.CurrentApp
}
