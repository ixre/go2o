/**
 * Copyright 2013 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-02-04 20:13
 * description :
 * history :
 */
package www

import (
	"github.com/atnet/gof/web"
	"go2o/src/app/cache/apicache"
	"html/template"
)

//todo: fiter valid partner is nil
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
		if b, siteConf := GetSiteConf(w, p); b {
			shops := apicache.GetShops(ctx.App, p.Id, p.Secret)
			if shops == nil {
				shops = []byte("{}")
			}
			ctx.App.Template().Execute(w, func(m *map[string]interface{}) {
				(*m)["partner"] = p
				(*m)["conf"] = siteConf
				(*m)["title"] = siteConf.IndexTitle
				(*m)["shops"] = template.HTML(shops)
			},
				"views/web/www/index.html",
				"views/web/www/inc/header.html",
				"views/web/www/inc/footer.html")
		}
	}
}
