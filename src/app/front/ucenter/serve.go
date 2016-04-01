/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	mw "github.com/labstack/echo/middleware"
	"go2o/src/app/util"
	"go2o/src/x/echox"
	"net/http"
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
	gc := &getC{}
	riseC := &personFinanceRiseC{}

	s.Static("/static/", "./public/static/") //静态资源
	s.Getx("/", mc.Index)
	s.Getx("/logout", mc.Logout)
	s.Anyx("/login", lc.Index)
	s.Getx("/change_device", mc.Change_device)
	s.Getx("/msc", mc.Msc)
	s.Getx("/msd", mc.Msd)
	s.Getx("/partner_connect", lc.Partner_connect)
	s.Getx("/partner_disconnect", lc.Partner_disconnect)
	s.Danyx("/basic/:action", bc)
	s.Danyx("/order/:action", oc)
	s.Danyx("/account/:action", ac)
	s.Getx("/get/qr/:code/:size", gc.GetQR)
	s.Danyx("/finance/rise/:action", riseC)
}

var (
	waitInit   bool = true
	pcServe    *echox.Echo
	mobiServe  *echox.Echo
	embedServe *echox.Echo
)

func getServe(path string) *echox.Echo {
	s := echox.New()
	s.SetTemplateRender(path)
	s.Use(mw.Recover())
	s.Use(memberLogonCheck)
	registerRoutes(s)
	return s
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if waitInit {
		pcServe = getServe("public/views/ucenter/pc")
		mobiServe = getServe("public/views/ucenter/mobi")
		embedServe = getServe("public/views/ucenter/app_embed")
		waitInit = false
	}
	switch util.GetBrownerDevice(r) {
	default:
	case util.DevicePC:
		pcServe.ServeHTTP(w, r)
	case util.DeviceTouchPad, util.DeviceMobile:
		mobiServe.ServeHTTP(w, r)
	case util.DeviceAppEmbed:
		embedServe.ServeHTTP(w, r)
	}
}
