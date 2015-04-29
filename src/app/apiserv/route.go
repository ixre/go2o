/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package apiserv

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
)

var (
	routes *web.RouteMap = new(web.RouteMap)
)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

func RegisterRoutes(c gof.App) {
	ws := &websocketC{App: c}

	routes.Add("^/ws/", func(ctx *web.Context) {
		//cross ajax request
		ctx.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
		mvc.Handle(ws, ctx, false)
	})

	routes.Add("/", func(ctx *web.Context) {
		ctx.ResponseWriter.Write([]byte("page not found"))
	})
}
