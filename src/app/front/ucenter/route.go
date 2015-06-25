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
	"net/http"
	"go2o/src/app/util"
)

var (
	routes   *mvc.Route = mvc.NewRoute(nil)
	moRoutes *mvc.Route = mvc.NewRoute(nil)
)

func GetRouter()*mvc.Route{
	return moRoutes
}

//处理请求
func Handle(ctx *web.Context) {
	switch util.GetBrownerDevice(ctx) {
	default:
	case util.DevicePC:
		ctx.Items["device_view_dir"] = "pc"
		routes.Handle(ctx)
	case util.DeviceTouchPad, util.DeviceMobile:
		ctx.Items["device_view_dir"] = "touchpad"
		moRoutes.Handle(ctx)
	}
}

func registerRoutes() {
	mc := &mainC{}
	bc := &basicC{}
	oc := &orderC{}
	ac := &accountC{}
	lc := &loginC{}

	// 静态文件处理
	sf := func(ctx *web.Context) {
		http.ServeFile(ctx.ResponseWriter, ctx.Request, "."+ctx.Request.URL.Path)
	}

	routes.Register("main", mc)
	routes.Register("basic", bc)
	routes.Register("order", oc)
	routes.Register("account", ac)
	routes.Register("login", lc)
	routes.Add("/logout", mc.Logout)
	routes.Add("/", mc.Index)
	routes.Add("^/static/",sf)

	// 注册触屏版路由
	moRoutes.Register("main", mc)
	moRoutes.Register("basic", bc)
	moRoutes.Register("order", oc)
	moRoutes.Register("account", ac)
	moRoutes.Register("login", lc)
	moRoutes.Add("/logout", mc.Logout)
	moRoutes.Add("/", mc.Index)

	// 为了使用IconFont
	moRoutes.Add("^/static/",sf)
}

func init() {
	registerRoutes()
}
