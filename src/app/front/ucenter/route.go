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
)

var (
	routes *mvc.Route = mvc.NewRoute(nil)
)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

func registerRoutes() {
	mc := &mainC{}
	bc := &basicC{}
	routes.Register("main", mc)
	routes.Register("basic", bc)
	routes.Register("order", &orderC{})
	routes.Register("account", &accountC{})
	routes.Register("login", &loginC{})
	routes.Add("/logout", mc.Logout)
	routes.Add("/", mc.Index)
}

func init(){
	registerRoutes()
}