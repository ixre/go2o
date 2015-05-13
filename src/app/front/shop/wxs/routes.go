/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package weixin

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	_ "github.com/atnet/gof/web/mvc"
	_ "net/http"
)

var (
	routes *web.RouteMap
)

//注册路由
func RegisterRoutes(context gof.App) {
	//	var mc *MainController //= &MainController{Context: context} //入口控制器
	//
	//	routes.Add("^/", func(wxs http.ResponseWriter, r *http.Request) {
	//			mvc.HandleRequest(mc, wxs, r)
	//		})
}
