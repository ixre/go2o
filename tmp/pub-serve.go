/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-16 21:45
 * description :
 * history :
 */

package main

import (
	"flag"
	"fmt"
	"go2o/app/pub"
	"go2o/core"
	"go2o/fix"
	"log"
	"os"
	"runtime"
)

func main() {
	var (
		ch       chan bool = make(chan bool)
		confFile string
		httpPort int
		help     bool
		newApp   *core.AppImpl
	)

	flag.IntVar(&httpPort, "port", 14280, "web server port")
	flag.BoolVar(&help, "help", false, "command usage")
	flag.StringVar(&confFile, "conf", "app.conf", "")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Ltime | log.Ldate | log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())
	newApp = core.NewApp(confFile)
	go fix.SignalNotify(ch)
	go pub.Listen(ch, newApp, fmt.Sprintf(":%d", httpPort))

	<-ch
}
