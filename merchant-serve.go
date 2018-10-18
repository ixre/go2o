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
    "github.com/jsix/gof"
    "github.com/jsix/gof/storage"
    "go2o/src/app/cache"
    "go2o/src/core"
    "go2o/src/core/service/dps"
    "log"
    "os"
    "runtime"
    "go2o/src/fix"
    "go2o/src/app/front/partner"
    "github.com/jsix/gof/web/session"
)

func main() {
    var (
        ch        chan bool = make(chan bool)
        confFile  string
        httpPort  int
        debug     bool
        trace     bool
        help      bool
        newApp    *core.MainApp
    )

    flag.IntVar(&httpPort, "port", 14281, "web server port")
    flag.BoolVar(&debug, "debug", false, "enable debug")
    flag.BoolVar(&trace, "trace", false, "enable trace")
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
    newApp = core.NewMainApp(confFile)
    if !newApp.Init(debug, trace) {
        os.Exit(1)
    }
    fix.CustomFix()
    go fix.SignalNotify(ch)

    gof.CurrentApp = newApp
    dps.Init(newApp)
    core.RegisterTypes()
    cache.Initialize(storage.NewRedisStorage(newApp.Redis()))
    session.Initialize(newApp.Storage(), "",false)
    go partner.Listen(ch, newApp, fmt.Sprintf(":%d", httpPort))

    <-ch
}