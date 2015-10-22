/**
 * Copyright 2015 @ z3q.net.
 * name : rest_server.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package app

import (
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"go2o/src/app/api"
	"go2o/src/core/service/goclient"
	"go2o/src/core/variable"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	API_DOMAIN   string
	API_HOST_CHK bool = false // 必须匹配Host
)

func RunRestApi(app gof.App, port int) {
	log.Println("** [ Go2o][ API][ Booted] - Api server running on port " + strconv.Itoa(port))

	//socket client
	time.Sleep(time.Second * 2) //等待启动Socket
	API_DOMAIN = app.Config().GetString(variable.ApiDomain)
	goclient.Configure("tcp", app.Config().GetString(variable.ClientSocketServer), app)

	var in *web.Interceptor = web.NewInterceptor(app, func(ctx *web.Context) {
		host := ctx.Request.URL.Host
		// todo: path compare
		if API_HOST_CHK && host != API_DOMAIN {
			http.Error(ctx.Response, "no such file", http.StatusNotFound)
			return
		}
		api.Handle(ctx)
	})

	//启动服务
	err := http.ListenAndServe(":"+strconv.Itoa(port), in)

	if err != nil {
		app.Log().Fatalln("ListenAndServer ", err)
	}
}
