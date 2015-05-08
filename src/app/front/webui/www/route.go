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
)

var (
	routes *mvc.Route = mvc.NewRoute(nil)
)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

//注册路由
func RegisterRoutes(c gof.App) {
	mc := &mainC{}
	sp := &shoppingC{}
	pc := &paymentC{}
	cc := &cartC{}
	uc := &userC{}

	routes.SingletonRegister("buy", sp)
	routes.SingletonRegister("shopping", sp)
	routes.Register("list", func()mvc.Controller{return &listC{}})
	routes.SingletonRegister("cart", cc)
	routes.SingletonRegister("user", uc)

	//处理错误
	routes.DeferFunc(func(ctx *web.Context) {
		if err, ok := recover().(error); ok {
			handleCustomError(ctx.ResponseWriter, c, err)
		}
	})

	// 购物车接口
	routes.Add("/cart_api_v1", cc.cartApi)
	// 支付
	routes.Add("^/pay/create", pc.Create)
	// 首页
	routes.Add("/", mc.Index)
	routes.Add("/user/g2m", uc.member)
}
