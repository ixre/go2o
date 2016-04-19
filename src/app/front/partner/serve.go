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
	"github.com/jsix/gof/web/session"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"go2o/src/x/echox"
	"net/url"
	"strings"
)

func GetServe() *echox.Echo {

	s := echox.New()
	s.SetTemplateRender("public/views/partner")
	s.Use(mw.Recover())
	s.Use(partnerLogonCheck) // 判断商户登陆状态

	s.Static("/static/", "./public/static/") //静态资源

	mc := new(mainC) //入口控制器
	s.Get("/", mc.Index)
	s.Anyx("/login", mc.Login)
	s.Postx("/export/getExportData", mc.exportData) //数据导出
	s.Postx("/upload.cgi", mc.Upload_post)          //上传文件
	s.Aanyx("/main/:action", mc)
	s.Aanyx("/shop/:action", new(shopC))             //商家门店控制器
	s.Aanyx("/goods/:action", new(goodsC))           //商品控制器
	s.Aanyx("/order/:action", new(orderC))           // 订单控制器
	s.Aanyx("/category/:action", new(categoryC))     // 商品分类控制器
	s.Aanyx("/conf/:action", new(configC))           // 商户设置控制器
	s.Aanyx("/prom/:action", new(promC))             // 促销控制器
	s.Aanyx("/delivery/:action", new(coverageAreaC)) // 配送区域控制器
	s.Aanyx("/member/:action", new(memberC))
	s.Aanyx("/sale/:action", new(saleC))
	s.Aanyx("/content/:action", new(contentC))
	s.Aanyx("/ad/:action", new(adC))
	s.Aanyx("/mss/:action", new(mssC))
	s.Aanyx("/editor/:action", new(editorC))
	s.Aanyx("/finance/:action", new(financeC))
	s.Aanyx("/person_finance/:action", new(personFinanceC))
	return s
}

func partnerLogonCheck(h echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		path := ctx.Request().URL.Path
		if path == "/login" || strings.HasPrefix(path, "/static/") {
			return h(ctx)
		}
		session := session.Default(ctx.Response(), ctx.Request())
		if id := session.Get("partner_id"); id != nil {
			ctx.Set("partner_id", id.(int))
			return h(ctx)
		}
		ctx.Response().Header().Set("Location", "/login?return_url="+
			url.QueryEscape(ctx.Request().URL.String()))
		ctx.Response().WriteHeader(302)
		return nil
	}
}
