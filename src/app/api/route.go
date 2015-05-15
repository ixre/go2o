/**
 * Copyright 2014 @ S1N1 Team.
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
	Routes web.Route = new(web.RouteMap)
)

//处理请求
func Handle(ctx *web.Context) {
	Routes.Handle(ctx)
}

func init() {
	bc := new(BaseC)
	//pc := &partnerC{}
	mc := &MemberC{BaseC: bc}

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

	Routes.Add("/", ApiTest)
	Routes.Add("/go2o_api_v1/mm_login", mc.login)       // 会员登陆接口
	Routes.Add("/go2o_api_v1/mm_register", mc.register) // 会员登陆接口
	Routes.Add("^/go2o_api_v1/member/", mc.handle)      // 会员接口
}
