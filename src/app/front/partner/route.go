/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/app/session"
	"net/http"
	"net/url"
)

var routes *web.RouteMap = new(web.RouteMap)

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
func RegisterRoutes(context gof.App) {
	mc := &mainC{App: context} //入口控制器
	sc := &shopC{context}      //商家门店控制器
	gc := &goodsC{context}     //食谱控制器
	lc := &loginC{App: context}
	cc := &commC{App: context}
	oc := &orderC{App: context}
	cat_c := &categoryC{App: context}
	conf_c := &configC{App: context}
	prom_c := &promC{App: context}

	routes.Add("^/comm/", func(ctx *web.Context) {
		mvc.Handle(cc, ctx, true)
	})

	routes.Add("^/pt/shop/",
		func(ctx *web.Context) {
			if b, id := chkLogin(ctx.Request); b {
				mvc.Handle(sc, ctx, true, id)
			} else {
				redirect(ctx)
			}
		})

	routes.Add("^/pt/category/",
		func(ctx *web.Context) {
			if b, id := chkLogin(ctx.Request); b {
				mvc.Handle(cat_c, ctx, true, id)
			} else {
				redirect(ctx)
			}
		})

	routes.Add("^/pt/order/",
		func(ctx *web.Context) {
			if b, id := chkLogin(ctx.Request); b {
				mvc.Handle(oc, ctx, true, id)
			} else {
				redirect(ctx)
			}
		})

	routes.Add("^/pt/goods/",
		func(ctx *web.Context) {
			if b, id := chkLogin(ctx.Request); b {
				mvc.Handle(gc, ctx, true, id)
			} else {
				redirect(ctx)
			}
		})

	routes.Add("^/pt/prom/",
		func(ctx *web.Context) {
			if b, id := chkLogin(ctx.Request); b {
				mvc.Handle(prom_c, ctx, true, id)
			} else {
				redirect(ctx)
			}
		})

	routes.Add("^/export/getExportData$", func(ctx *web.Context) {
		if b, id := chkLogin(ctx.Request); b {
			GetExportData(ctx, id)
		} else {
			redirect(ctx)
		}
	})

	routes.Add("^/pt/conf/",
		func(ctx *web.Context) {
			if b, id := chkLogin(ctx.Request); b {
				mvc.Handle(conf_c, ctx, true, id)
			} else {
				redirect(ctx)
			}
		})

	routes.Add("^/login$", func(ctx *web.Context) {
		mvc.Handle(lc, ctx, true)
	})

	routes.Add("^/", func(ctx *web.Context) {
		if b, id := chkLogin(ctx.Request); b {
			mvc.Handle(mc, ctx, true, id)
		} else {
			redirect(ctx)
		}
	})
}
