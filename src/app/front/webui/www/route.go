/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package www

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"net/http"
)

var (
	routes *mvc.Route = mvc.NewRoute(nil)
)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

func handleError(w http.ResponseWriter, err error) {
	w.Write([]byte(`<html><head><title></title></head>` +
		`<body><span style="color:red">` + err.Error() + `</span></body></html>`))
}

//注册路由
func RegisterRoutes(c gof.App) {
	mc := &mainC{}
	sp := &shoppingC{}
	pc := &paymentC{}

	routes.RegisterController("buy",sp)
	routes.RegisterController("shopping",sp)
	//处理错误
	routes.DeferFunc(func(ctx *web.Context) {
		if err, ok := recover().(error); ok {
			handleCustomError(ctx.ResponseWriter, c, err)
		}
	})

	// 购物车接口
	routes.Add("^/cart_api_v1$",sp.cartApi)
	routes.Add("^/cart/*",sp.cart)

	routes.Add("^/pay/", func(ctx *web.Context) {
		mvc.Handle(pc, ctx, true)
	})

	// add route for main controller
	routes.Add("^/[^/]*$", func(ctx *web.Context) {
		mvc.Handle(mc, ctx, true)
	})
}
