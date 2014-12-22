package weixin

import (
	_ "net/http"
	"ops/cf/app"
	"ops/cf/web"
	_ "ops/cf/web/mvc"
)

var (
	routes *web.RouteMap
)

//注册路由
func RegistRoutes(context app.Context) {
	//	var mc *MainController //= &MainController{Context: context} //入口控制器
	//
	//	routes.Add("^/", func(weixin http.ResponseWriter, r *http.Request) {
	//			mvc.HandleRequest(mc, weixin, r)
	//		})
}
