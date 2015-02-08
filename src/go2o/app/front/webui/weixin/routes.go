package weixin

import (
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web"
	_ "github.com/atnet/gof/web/mvc"
	_ "net/http"
)

var (
	routes *web.RouteMap
)

//注册路由
func RegisterRoutes(context app.Context) {
	//	var mc *MainController //= &MainController{Context: context} //入口控制器
	//
	//	routes.Add("^/", func(weixin http.ResponseWriter, r *http.Request) {
	//			mvc.HandleRequest(mc, weixin, r)
	//		})
}
