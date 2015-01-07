package weixin

import (
	_ "net/http"
	"github.com/newmin/gof/app"
	"github.com/newmin/gof/web"
	_ "github.com/newmin/gof/web/mvc"
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
