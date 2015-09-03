/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ols

import (
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/mvc"
	"go2o/src/app/front/shop/ols/mos"
	"go2o/src/app/util"
)

var (
	routes *mvc.Route = mvc.NewRoute(nil)
)

//处理请求
func Handle(ctx *web.Context) {
	switch util.GetBrownerDevice(ctx) {
	default:
	case util.DevicePC:
		ctx.Items["device_view_dir"] = "pc"
		routes.Handle(ctx)
	case util.DeviceTouchPad, util.DeviceMobile:
		ctx.Items["device_view_dir"] = "touchpad"
		mos.Handle(ctx)
	case util.DeviceAppEmbed:
		ctx.Items["device_view_dir"] = "app_embed"
		routes.Handle(ctx)
	}
}

//注册路由
func registerRoutes() {
	mc := &MainC{}

	sp := &ShoppingC{}
	pc := &PaymentC{}
	cc := &CartC{}
	uc := &UserC{}
	lc := &ListC{}
	ctc := &ContentC{}

	routes.Register("main", mc)
	routes.Register("buy", sp)
	routes.Register("shopping", sp)
	routes.Register("list", lc)
	routes.Register("cart", cc)
	routes.Register("user", uc)
	routes.Register("content", ctc)
	routes.Register("pay", pc)

	//处理错误
	routes.DeferFunc(func(ctx *web.Context) {
		if err, ok := recover().(error); ok {
			HandleCustomError(ctx.Response, ctx, err)
		}
	})

	// 购物车接口
	routes.Add("/cart_api_v1", cc.CartApiHandle)

	// 支付异步提醒
	routes.Add("/pay/notify/*", pc.Notify_post)

	// 首页
	routes.Add("/", mc.Index)
	routes.Add("/st/*", lc.SaleTagGoodsList)
	routes.Add("/user/jump_m", uc.JumpToMCenter)
	routes.Add("^/c-[0-9-]+.htm", lc.List_Index)
	routes.Add("^/item-[0-9-]+.htm", lc.GoodsDetails)
}

func init() {
	registerRoutes()
}
