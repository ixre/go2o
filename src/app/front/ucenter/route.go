/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/app/front"
)

var (
	routes *mvc.Route = mvc.NewRoute(nil)
	moRoutes *mvc.Route = mvc.NewRoute(nil)
)

//处理请求
func Handle(ctx *web.Context) {
	switch front.GetBrownerDevice(ctx) {
		default:
		case front.DevicePC:
		ctx.Items["device_view_dir"] = "pc"
		routes.Handle(ctx)
		case front.DeviceTouchPad, front.DeviceMobile:
		ctx.Items["device_view_dir"] = "touchpad"
		moRoutes.Handle(ctx)
	}
}

func registerRoutes() {
	mc := &mainC{}
	bc := &basicC{}
	oc :=  &orderC{}
	ac :=  &accountC{}
	lc := &loginC{}
	routes.Register("main", mc)
	routes.Register("basic", bc)
	routes.Register("order",oc)
	routes.Register("account",ac)
	routes.Register("login", lc)
	routes.Add("/logout", mc.Logout)
	routes.Add("/device", mc.changeDevice)
	routes.Add("/", mc.Index)

	// 注册触屏版路由
	moRoutes.Register("main", mc)
	moRoutes.Register("basic", bc)
	moRoutes.Register("order",oc)
	moRoutes.Register("account",ac)
	moRoutes.Register("login", lc)
	moRoutes.Add("/logout", mc.Logout)
	moRoutes.Add("/device", mc.changeDevice)
	moRoutes.Add("/", mc.Index)
}

func init() {
	registerRoutes()
}
