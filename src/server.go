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
	"app"
	"com/share/glob"
	"com/share/variable"
	"flag"
	"fmt"
	a "ops/cf/app"
	"os"
	"strconv"
	"strings"
)

func main() {
	var (
		ch         chan int = make(chan int)
		httpPort   int
		socketPort int
		mode       string //启动模式: h开启http,s开启socket,a开启所有
		debug      bool
		help       bool
		ctx        a.Context
	)

	ctx = glob.NewContext()
	httpPort, _ = strconv.Atoi(ctx.Config().Get(variable.ServerPort))
	socketPort, _ = strconv.Atoi(ctx.Config().Get(variable.SocketPort))
	flag.IntVar(&httpPort, "port", httpPort, "web server port")
	flag.IntVar(&socketPort, "port2", socketPort, "socket server port")
	flag.StringVar(&mode, "mode", "sh", "boot mode.'h'- boot http service,'s'- boot socket service")
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.BoolVar(&help, "help", false, "command usage")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if gcx, ok := ctx.(*glob.AppContext); ok {
		gcx.Init(debug)
	} else {
		fmt.Println("app context err")
		os.Exit(1)
		return
	}
	app.Init(ctx)

	var booted bool
	if strings.Contains(mode, "s") {
		booted = true
		go app.RunSocket(ctx, socketPort, debug)
	}

	if strings.Contains(mode, "h") {
		booted = true
		go app.RunWeb(ctx, httpPort, debug)
	}

	if booted {
		<-ch
	}
}
