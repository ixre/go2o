/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package api

import (
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
)

var (
	routes *mvc.Route = mvc.NewRoute(nil)
)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

func init() {

	pc := &partnerC{}
	mc := &memberC{}

//	ws := &websocketC{App: c}

//	routes.Add("^/ws/", func(ctx *web.Context) {
//		//cross ajax request
//		ctx.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
//		mvc.Handle(ws, ctx, false)
//	})
//
//	routes.Add("/", func(ctx *web.Context) {
//		ctx.ResponseWriter.Write([]byte("page not found"))
//	})

	routes.Register("partner",pc)
	routes.Register("member",mc)
	routes.Add("/mm_login",mc.login)   // 会员登陆接口
	routes.Add("/",HandleApi)

}
