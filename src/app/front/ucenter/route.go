/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
)

var (
	routes *mvc.Route = mvc.NewRoute(nil)
)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}


func RegisterRoutes(c gof.App) {
	mc := &mainC{}
    routes.RegisterController("order",&orderC{})
    routes.RegisterController("account",&accountC{})
	routes.RegisterController("login",&loginC{})

	routes.Add("/logout",mc.Logout)
    routes.Add("/",mc.Index)
}

