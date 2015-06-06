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
	"go2o/src/core/service/dps"
)


type mainC struct {
	*BaseC
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
		m := this.GetMember(ctx)

		if this.HandleIndexGo(ctx) {
			return
		}

		siteConf := this.GetSiteConf(ctx)
		newGoods := dps.SaleService.GetValueGoodsBySaleTag(p.Id,"new-goods",0,12)
		hotSales := dps.SaleService.GetValueGoodsBySaleTag(p.Id,"hot-sales",0,12)

		ctx.App.Template().Execute(w, gof.TemplateDataMap{
			"partner": p,
			"conf":    siteConf,
			"newGoods" :newGoods,
			"hotSales" : hotSales,
			"member":m,
		},
		"views/shop/{device}/index.html",
		"views/shop/{device}/inc/header.html",
		"views/shop/{device}/inc/footer.html")

	}
}
