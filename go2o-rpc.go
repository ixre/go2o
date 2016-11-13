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
	"go2o/core"
	"go2o/core/service/dps"
	"go2o/core/service/thrift"
	"os"
)

func main() {
	var (
		addr  string
		conf  string
		debug bool
		trace bool
	)

	flag.StringVar(&addr, "addr", "localhost:14288", "Address to listen to")
	flag.StringVar(&conf, "conf", "app.conf", "Config file path")
	flag.BoolVar(&debug, "debug", false, "Enable debug")
	flag.BoolVar(&trace, "trace", false, "Enable trace")
	flag.Parse()

	newApp := core.NewMainApp(conf)
	if !newApp.Init(debug, trace) {
		os.Exit(1)
	}
	gof.CurrentApp = newApp
	dps.Init(newApp)

	err := thrift.Listen(addr, false)
	if err != nil {
		log.Println("error running ", addr, " :", err.Error())
	}
}
