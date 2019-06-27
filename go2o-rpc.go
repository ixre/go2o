package main

/**
 * Copyright 2015 @ at3.net.
 * name : go2o-rpc.go
 * author : jarryliu
 * date : 2016-11-12 19:29
 * description :
 * history :
 */

import (
	"flag"
	"github.com/ixre/gof"
	"github.com/ixre/gof/log"
	"go2o/app"
	"go2o/core"
	"go2o/core/service/rsi"
	rs "go2o/core/service/thrift/service"
	"os"
)

func main() {
	var (
		addr  string
		conf  string
		debug bool
		trace bool
	)

	flag.StringVar(&addr, "addr", ":1427", "Address to listen to")
	flag.StringVar(&conf, "conf", "app.conf", "Config file path")
	flag.BoolVar(&debug, "debug", false, "Enable debug")
	flag.BoolVar(&trace, "trace", false, "Enable trace")
	flag.Parse()

	newApp := core.NewApp(conf)
	if !core.Init(newApp, debug, trace) {
		os.Exit(1)
	}
	gof.CurrentApp = newApp
	rsi.Init(newApp, app.FlagRpcServe)
	//app.Configure(hook.HookUp, newApp, app.FlagRpcServe)

	err := rs.ListenAndServe(addr, false)
	if err != nil {
		log.Println("error running ", addr, " :", err.Error())
	}
}
