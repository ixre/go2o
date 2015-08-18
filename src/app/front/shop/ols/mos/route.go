/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package mos

import (
	"github.com/jrsix/gof/web"
	"net/http"
)

var (
	routes *web.RouteMap = new(web.RouteMap)
)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

func handleError(w http.ResponseWriter, err error) {
	w.Write([]byte(`<span style="color:red">` + err.Error() + `</span>`))
}

//注册路由
func registerRoutes() {
	mc := new(mainC)
	routes.Add("/", mc.Index)
}

func init() {
	registerRoutes()
}
