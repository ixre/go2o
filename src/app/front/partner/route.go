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
	routes.RegisterController("shop", new(shopC))             //商家门店控制器
	routes.RegisterController("goods", new(goodsC))           //商品控制器
	routes.RegisterController("comm", new(commC))             // 通用控制器
	routes.RegisterController("order", new(orderC))           // 订单控制器
	routes.RegisterController("category", new(categoryC))     // 商品分类控制器
	routes.RegisterController("conf", new(configC))           // 商户设置控制器
	routes.RegisterController("prom", new(promC))             // 促销控制器
	routes.RegisterController("delivery", new(coverageAreaC)) // 配送区域控制器

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
