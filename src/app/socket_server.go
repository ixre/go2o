/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-16 19:03
 * description :
 * history :
 */

package app

import (
	"fmt"
	"github.com/jsix/gof"
	"go2o/src/core"
	"go2o/src/core/infrastructure"
	"go2o/src/core/service"
	"os"
	"strconv"
)

func RunSocket(ctx gof.App, port int, debug, trace bool) {

	if gcx, ok := ctx.(*core.MainApp); ok {
		if !gcx.Loaded {
			gcx.Init(debug, trace)
		}
	} else {
		fmt.Println("app context err")
		os.Exit(1)
		return
	}

	if debug {
		fmt.Println("[Started]:Socket server (with debug) running on port [" +
			strconv.Itoa(port) + "]:")
		infrastructure.DebugMode = true
	} else {
		fmt.Println("[Started]:Socket server running on port [" +
			strconv.Itoa(port) + "]:")
	}
	service.ServerListen("tcp", ":"+strconv.Itoa(port), ctx)
}
