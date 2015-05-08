/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
)

var routes *mvc.Route = mvc.NewRoute(nil)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

//注册路由
func RegisterRoutes() {
	bc := new(baseC)
	mc := &mainC{} //入口控制器
	lc := &loginC{}
<<<<<<< HEAD
	routes.Register("shop", new(shopC))             //商家门店控制器
	routes.Register("goods", new(goodsC))           //商品控制器
	routes.Register("comm", new(commC))             // 通用控制器
	routes.Register("order", new(orderC))           // 订单控制器
	routes.Register("category", new(categoryC))     // 商品分类控制器
	routes.Register("conf", new(configC))           // 商户设置控制器
	routes.Register("prom", new(promC))             // 促销控制器
	routes.Register("delivery", new(coverageAreaC)) // 配送区域控制器
=======
	routes.RegisterController("shop", new(shopC))             //商家门店控制器
	routes.RegisterController("goods", new(goodsC))           //商品控制器
	routes.RegisterController("comm", new(commC))             // 通用控制器
	routes.RegisterController("order", new(orderC))           // 订单控制器
	routes.RegisterController("category", new(categoryC))     // 商品分类控制器
	routes.RegisterController("conf", new(configC))           // 商户设置控制器
	routes.RegisterController("prom", new(promC))             // 促销控制器
	routes.RegisterController("delivery", new(coverageAreaC)) // 配送区域控制器
>>>>>>> 55b2cb6c58ebd6b2d1e8bbbd81858ff12b1b2eee

	routes.Add("/export/getExportData", func(ctx *web.Context) {
		if b, id := chkLogin(ctx); b {
			GetExportData(ctx, id)
		} else {
			redirect(ctx)
		}
	})

	routes.Add("/login", func(ctx *web.Context) {
		mvc.Handle(lc, ctx, true)
	})

	routes.Add("^/[^/]*$", func(ctx *web.Context) {
		if bc.Requesting(ctx) {
			mvc.Handle(mc, ctx, true)
		}
		bc.RequestEnd(ctx)
	})

}
