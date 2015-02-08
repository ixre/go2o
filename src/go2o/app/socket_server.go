/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-16 19:03
 * description :
 * history :
 */

package app

import (
	"fmt"
	"github.com/atnet/gof/app"
	"go2o/core/infrastructure"
	"go2o/core/service"
	"go2o/core/share/glob"
	"os"
	_ "runtime/debug"
	"strconv"
)

func RunSocket(ctx app.Context, port int, debug, trace bool) {

	if gcx, ok := ctx.(*glob.AppContext); ok {
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
