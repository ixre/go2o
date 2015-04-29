/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/app/session"
	"net/http"
	"net/url"
)

var routes *mvc.Route = mvc.NewRoute(nil)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

func chkLogin(r *http.Request) (b bool, partnerId int) {
	//todo:仅仅做了id的检测，没有判断有效性
	i, err := session.GetLSession().GetPartnerIdFromCookie(r)
	return err == nil, i
}

func redirect(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	w.Write([]byte("<script>window.parent.location.href='/login?return_url=" +
		url.QueryEscape(r.URL.String()) + "'</script>"))
}

//注册路由
func RegisterRoutes() {
	mc := &mainC{} //入口控制器
	lc := &loginC{}
	routes.RegisterController("shop",&shopC{})      //商家门店控制器
    routes.RegisterController("goods",&goodsC{})    //商品控制器
	routes.RegisterController("comm",&commC{})
	routes.RegisterController("order",&orderC{})
	routes.RegisterController("category", &categoryC{})
	routes.RegisterController("conf",&configC{})
	routes.RegisterController("prom", &promC{})

	routes.Add("^/export/getExportData$", func(ctx *web.Context) {
		if b, id := chkLogin(ctx.Request); b {
			GetExportData(ctx, id)
		} else {
			redirect(ctx)
		}
	})

	register("shop")
	register("goods")
	register("comm")
	register("order")
	register("category")
	register("conf")
	register("prom")


	routes.Add("^/login$", func(ctx *web.Context) {
		mvc.Handle(lc, ctx, true)
	})

	routes.Add("^/[^/]*$", func(ctx *web.Context) {
		if b, id := chkLogin(ctx.Request); b {
			mvc.Handle(mc, ctx, true, id)
		} else {
			redirect(ctx)
		}
	})

}

func register(name string){
	routes.Add("^/"+name+"/*", func(ctx *web.Context) {
		if b, id := chkLogin(ctx.Request); b {
			mvc.Handle(routes.GetController(name), ctx, true, id)
		} else {
			redirect(ctx)
		}
	})
}
