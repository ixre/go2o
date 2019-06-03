/**
 * Copyright 2015 @ to2.net.
 * name : defer.go
 * author : jarryliu
 * date : 2016-01-06 17:43
 * description :
 * history :
 */
package daemon

import (
	"log"
	"runtime/debug"
)

func Recover() {
	if r := recover(); r != nil {
		log.Println("[ Daemon][ Recover]-", r, "\n", string(debug.Stack()))
	}
}
