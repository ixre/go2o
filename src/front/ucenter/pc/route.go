/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package pc

import (
	mw "github.com/labstack/echo/middleware"
	"go2o/src/front/ucenter"
	"go2o/src/x/echox"
)

//处理请求
//func Handle(ctx *web.Context) {
//	switch util.GetBrownerDevice(ctx) {
//	default:
//	case util.DevicePC:
//		ctx.Items["device_view_dir"] = "pc"
//		routes.Handle(ctx)
//	case util.DeviceTouchPad, util.DeviceMobile:
//		ctx.Items["device_view_dir"] = "touchpad"
//		moRoutes.Handle(ctx)
//	case util.DeviceAppEmbed:
//		ctx.Items["device_view_dir"] = "app_embed"
//		routes.Handle(ctx)
//	}
//}

func registerRoutes(s *echox.Echo) {
	mc := &mainC{}
	bc := &basicC{}
	oc := &orderC{}
	ac := &accountC{}
	lc := &loginC{}

	s.Static("/static/", "./public/static/") //静态资源
	s.Getx("/", mc.Index)
	s.Getx("/logout", mc.Logout)
	s.Anyx("/login", lc.Index)
	//	s.Danyx("/main/:action",mc)
	s.Danyx("/basic/:action", bc)
	s.Danyx("/order/:action", oc)
	s.Danyx("/account/:action", ac)
}

func GetServe() *echox.Echo {
	s := echox.New()
	s.SetTemplateRender("public/views/ucenter/pc")
	s.Use(mw.Recover())
	s.Use(ucenter.MemberLogonCheck)
	registerRoutes(s)
	return s
}
