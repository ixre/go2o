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
	"github.com/jsix/gof/web/session"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"go2o/src/app/util"
	"go2o/src/x/echox"
	"net/url"
	"strings"
)

var routes *mvc.Route = mvc.NewRoute(nil)

//处理请求
func Handle(ctx *echox.Context) error {
	routes.Handle(ctx)
}

//注册路由
func registerRoutes() {
	//bc := new(baseC)
	mc := &mainC{} //入口控制器
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

	routes.Add("/export/getExportData", func(ctx *echox.Context) error {
		if b, id := chkLogin(ctx); b {
			GetExportData(ctx, id)
		} else {
			redirect(ctx)
		}
	})

	routes.Add("/upload.cgi", mc.Upload_post)

	// 静态文件处理
	routes.Add("/static/*", util.HttpStaticFileHandler)
	routes.Add("/img/*", util.HttpImageFileHandler)

	// 首页
	//routes.Add("/", mc.Index)

}

func GetServe() *echox.Echo {
	mc := &mainC{} //入口控制器
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

	s := echox.New()
	r := echox.NewGoTemplateForEcho("public/views/partner")
	s.SetRenderer(r)
	s.Use(mw.Recover())
	s.Use(partnerLogonCheck) // 判断商户登陆状态
	s.Static("/static/", "./public/static/")
	s.Get("/", mc.Index)
	s.Anyx("/login", mc.Login)
	s.Danyx("/main/:action", mc)
	return s
}

func partnerLogonCheck(ctx *echo.Context) error {
	path := ctx.Request().URL.Path
	if path == "/login" || strings.HasPrefix(path, "/static/") {
		return nil
	}
	session := session.Default(ctx.Response(), ctx.Request())
	id := session.Get("partner_id")
	if id != nil {
		ctx.Set("partner_id", id.(int))
		return nil
	}
	ctx.Response().Header().Set("Location", "/login?return_url="+
		url.QueryEscape(ctx.Request().URL.String()))
	ctx.Response().WriteHeader(302)
	ctx.Done()
	return nil
}
