/**
 * Copyright 2015 @ S1N1 Team.
 * name : core.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package core

import (
	"github.com/atnet/gof"
)

// 全局APP上下文
var GlobalApp gof.App

func SetGlobalApp(a gof.App) {
	if a != nil {
		GlobalApp = a
	}
}
