/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ols

import (
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"github.com/atnet/gof/util"
	"go2o/src/app/front/shop/ols/mos"
)

var (
	routes *mvc.Route = mvc.NewRoute(nil)
)

//处理请求
func Handle(ctx *web.Context) {
//	mos.Handle(ctx)
//	return
	if util.IsMobileAgent(ctx.Request.UserAgent()){
		mos.Handle(ctx)
	}else {
		routes.Handle(ctx)
	}
}

//注册路由
func registerRoutes() {
	mc := &mainC{}

	sp := &shoppingC{}
	pc := &paymentC{}
	cc := &cartC{}
	uc := &userC{}
	lc := &listC{}

	routes.Register("buy", sp)
	routes.Register("shopping", sp)
	routes.Register("list", lc)
	routes.Register("cart", cc)
	routes.Register("user", uc)

	//处理错误
	routes.DeferFunc(func(ctx *web.Context) {
		if err, ok := recover().(error); ok {
			handleCustomError(ctx.ResponseWriter, ctx, err)
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


func init(){
	registerRoutes()
}