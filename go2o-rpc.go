/**
 * Copyright 2015 @ at3.net.
 * name : go2o-rpc.go
 * author : jarryliu
 * date : 2016-11-12 19:29
 * description :
 * history :
 */
package main

import (
	"flag"
	"github.com/jsix/gof"
	"github.com/jsix/gof/log"
	"go2o-web/src/hook"
	"go2o/app"
	"go2o/core"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift"
	"os"
)

func main() {
	var (
		addr    string
		conf    string
		confDir string
		debug   bool
		trace   bool
	)

	flag.StringVar(&addr, "addr", ":14280", "Address to listen to")
	flag.StringVar(&conf, "conf", "app.conf", "Config file path")
	flag.StringVar(&confDir, "conf-dir", "./conf", "config file directory")
	flag.BoolVar(&debug, "debug", false, "Enable debug")
	flag.BoolVar(&trace, "trace", false, "Enable trace")
	flag.Parse()

	newApp := core.NewApp(conf)
	if !core.Init(newApp, debug, trace) {
		os.Exit(1)
	}
	gof.CurrentApp = newApp
	rsi.Init(newApp, app.FlagRpcServe, confDir)
	app.Configure(hook.HookUp, newApp, app.FlagRpcServe)

	err := thrift.ListenAndServe(addr, false)
	if err != nil {
		log.Println("error running ", addr, " :", err.Error())
	}
}
