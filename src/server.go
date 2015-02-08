/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-16 21:45
 * description :
 * history :
 */

package main

import (
	"flag"
	"fmt"
	a "github.com/atnet/gof/app"
	"go2o/app"
	"go2o/app/daemon"
	"go2o/core/share/glob"
	"go2o/core/share/variable"
	"os"
	"strconv"
	"strings"
)

func main() {
	var (
		ch         chan bool = make(chan bool)
		httpPort   int
		socketPort int
		mode       string //启动模式: h开启http,s开启socket,a开启所有
		debug      bool
		trace      bool
		help       bool
		ctx        a.Context
	)

	ctx = glob.NewContext()
	httpPort, _ = strconv.Atoi(ctx.Config().GetString(variable.ServerPort))
	socketPort, _ = strconv.Atoi(ctx.Config().GetString(variable.SocketPort))
	flag.IntVar(&httpPort, "port", httpPort, "web server port")
	flag.IntVar(&socketPort, "port2", socketPort, "socket server port")
	flag.StringVar(&mode, "mode", "sh", "boot mode.'h'- boot http service,'s'- boot socket service")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.BoolVar(&trace, "trace", false, "enable trace")
	flag.BoolVar(&help, "help", false, "command usage")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if gcx, ok := ctx.(*glob.AppContext); ok {
		gcx.Init(debug, trace)
	} else {
		fmt.Println("app context err")
		os.Exit(1)
		return
	}
	app.Init(ctx)

	daemon.Run(ctx)

	var booted bool
	if strings.Contains(mode, "s") {
		booted = true
		go app.RunSocket(ctx, socketPort, debug, trace)
	}

	if strings.Contains(mode, "h") {
		booted = true
		go app.RunWeb(ctx, httpPort, debug, trace)
	}

	if booted {
		<-ch
	}
}
