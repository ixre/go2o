/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package master

import (
	"github.com/jrsix/gof/web"
	"github.com/jrsix/gof/web/mvc"
)

var routes *mvc.Route = mvc.NewRoute(nil)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

func init() {
	registerRoutes()
}

//注册路由
func registerRoutes() {
	//bc := new(baseC)
	mc := &mainC{} //入口控制器
	routes.Register("partner", &partnerC{})

	//routes.Register("shop", new(shopC))             //商家门店控制器

	routes.Add("/export/getExportData", mc.exportData)

	routes.Add("/dashboard", mc.Dashboard)
	routes.Add("/login", mc.Login)
	routes.Add("/logout", mc.Logout)
	routes.Add("/upload.cgi", mc.Upload_post)

}
