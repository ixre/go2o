/**
 * Copyright 2013 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-02-04 20:13
 * description :
 * history :
 */
package ols

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"html/template"
)

//todo: filter valid partner is nil
type mainC struct {
	*baseC
}

// 处理跳转
func (this *mainC) HandleIndexGo(ctx *web.Context) bool {
	r, w := ctx.Request, ctx.ResponseWriter
	g := r.URL.Query().Get("go")
	if g == "buy" {
		w.Header().Add("Location", "/list")
		w.WriteHeader(302)
		return true
	}
	return false
}

func (this *mainC) Index(ctx *web.Context) {
	if this.Requesting(ctx) {
		_, w := ctx.Request, ctx.ResponseWriter
		p := this.GetPartner(ctx)

		if this.HandleIndexGo(ctx) {
			return
		}



		siteConf := this.GetSiteConf(ctx)

			shops := GetShops(ctx.App, p.Id)
			if shops == nil {
				shops = []byte("{}")
			}
			ctx.App.Template().Execute(w, gof.TemplateDataMap{
				"partner": p,
				"conf":    siteConf,
				"shops":   template.HTML(shops),
			},
				"views/shop/ols/index.html",
				"views/shop/ols/inc/header.html",
				"views/shop/ols/inc/footer.html")

	}
}
