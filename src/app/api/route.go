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
	"github.com/jrsix/gof/web"
	"github.com/jrsix/gof/web/mvc"
)

var (
	Routes     web.Route = new(web.RouteMap)
	PathPrefix           = "/go2o_api_v1"
)

//处理请求
func Handle(ctx *web.Context) {
	Routes.Handle(ctx)
}

// 处理请求
func HandleUrlFunc(v interface{}) func(*web.Context) {
	return func(ctx *web.Context) {
		mvc.HandlePath(v, ctx, splitPath(ctx), false)
	}
}

func splitPath(ctx *web.Context) string {
	return ctx.Request.URL.Path[len(PathPrefix):]
}

func init() {
	bc := new(BaseC)
	pc := &partnerC{bc}
	mc := &MemberC{bc}
	gc := &getC{bc}

	Routes.Add("/", ApiTest)
	Routes.Add(PathPrefix+"/mm_login", mc.Login)           // 会员登陆接口
	Routes.Add(PathPrefix+"/mm_register", mc.Register)     // 会员登陆接口
	Routes.Add(PathPrefix+"/member/*", HandleUrlFunc(mc))  // 会员接口
	Routes.Add(PathPrefix+"/partner/*", HandleUrlFunc(pc)) // 会员接口

	// 会员接口
	Routes.Add("/go2o_api_v1/get/*", HandleUrlFunc(gc))
}
