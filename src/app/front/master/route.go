/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package master

import (
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
)

var routes *mvc.Route = mvc.NewRoute(nil)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

func init(){
	registerRoutes()
}
//注册路由
func registerRoutes() {
	//bc := new(baseC)
	//mc := &mainC{} //入口控制器
	lc := &loginC{}
	//routes.Register("shop", new(shopC))             //商家门店控制器

//	routes.Add("/export/getExportData", func(ctx *web.Context) {
//		if b, id := chkLogin(ctx); b {
//			GetExportData(ctx, id)
//		} else {
//			redirect(ctx)
//		}
//	})

	routes.Add("/login",lc.Login)
}
