/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package master

import (
	"github.com/jsix/gof/web/session"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"go2o/src/x/echox"
	"net/url"
)

//注册路由
func registerRoutes(s *echox.Echo) {
	mc := &mainC{} //入口控制器
	s.Getx("/", mc.Index)
	s.Getx("/dashboard", mc.Dashboard)
	s.Anyx("/login", mc.Login)
	s.Getx("/logout", mc.Logout)
	s.Postx("/upload.cgi", mc.Upload_post)
	s.Danyx("/partner/:action", new(partnerC))
	s.Postx("/export/getExportData", mc.exportData)

}

func GetServe() *echox.Echo {
	s := echox.New()
	s.SetRenderer(echox.NewGoTemplateForEcho("public/views/master"))
	s.Use(mw.Recover())
	s.Use(masterLogonCheck) // 判断商户登陆状态
	registerRoutes(s)
	return s
}

func masterLogonCheck(ctx *echo.Context) error {
	path := ctx.Request().URL.Path
	if path == "/login" {
		return nil
	}
	session := session.Default(ctx.Response(), ctx.Request())
	if id := session.Get("is_master"); id == nil {
		ctx.Response().Header().Set("Location", "/login?return_url="+
			url.QueryEscape(ctx.Request().URL.String()))
		ctx.Response().WriteHeader(302)
		ctx.Done()
	}
	return nil
}
