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
)

var (
	routes web.Route =new(web.RouteMap)
)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

func init() {
	bc := new(baseC)
	//pc := &partnerC{}
	mc := &memberC{baseC:bc}

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


	routes.Add("/",HandleApi)
	routes.Add("/go2o_api_v1/mm_login",mc.login)    // 会员登陆接口
	routes.Add("/go2o_api_v1/mm_register",mc.register)    // 会员登陆接口
	routes.Add("^/go2o_api_v1/member/",mc.handle)	// 会员接口
}
