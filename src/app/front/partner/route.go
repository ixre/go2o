/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/mvc"
	"go2o/src/app/util"
)

var routes *mvc.Route = mvc.NewRoute(nil)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

//注册路由
func registerRoutes() {
	//bc := new(baseC)
	mc := &mainC{} //入口控制器
	lc := &loginC{}
	routes.Register("main", new(mainC))
	routes.Register("shop", new(shopC))             //商家门店控制器
	routes.Register("goods", new(goodsC))           //商品控制器
	routes.Register("comm", new(commC))             // 通用控制器
	routes.Register("order", new(orderC))           // 订单控制器
	routes.Register("category", new(categoryC))     // 商品分类控制器
	routes.Register("conf", new(configC))           // 商户设置控制器
	routes.Register("prom", new(promC))             // 促销控制器
	routes.Register("delivery", new(coverageAreaC)) // 配送区域控制器
	routes.Register("member", new(memberC))
	routes.Register("sale", new(saleC))
	routes.Register("content", new(contentC))
	routes.Register("ad", new(adC))
	routes.Register("mss", new(mssC))
	routes.Register("editor", new(editorC))

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

	routes.Add("/upload.cgi", mc.Upload_post)

	// 静态文件处理
	routes.Add("/static/*", util.HttpStaticFileHandler)
	routes.Add("/img/*",util.HttpImageFileHandler)

	// 首页
	routes.Add("/", mc.Index)

}

func init() {
	registerRoutes()
}
