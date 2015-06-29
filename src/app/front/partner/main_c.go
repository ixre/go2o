/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/app/front"
)

var _ mvc.Filter = new(mainC)

type mainC struct {
	*baseC
	*front.WebCgi
}

//入口
func (this *mainC) Index(ctx *web.Context) {
	if this.baseC.Requesting(ctx) {
		ctx.Response.Write([]byte("<script>location.replace('/main/dashboard')</script>"))
	}
	this.baseC.RequestEnd(ctx)
}

func (this *mainC) Logout(ctx *web.Context) {
	ctx.Session().Destroy()
	ctx.Response.Write([]byte("<script>location.replace('/login')</script>"))
}

//商户首页
func (this *mainC) Dashboard(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	pt, _ := this.GetPartner(ctx)

	var mf gof.TemplateDataMap = gof.TemplateDataMap{
		"partner": pt,
		"loginIp": r.Header.Get("USER_ADDRESS"),
	}
	ctx.App.Template().Execute(w, mf, "views/partner/dashboard.html")
}

//商户汇总页
func (this *mainC) Summary(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	pt, _ := this.GetPartner(ctx)

	ctx.App.Template().Execute(w,
		gof.TemplateDataMap{
			"partner": pt,
			"loginIp": r.Header.Get("USER_ADDRESS"),
		},
		"views/partner/summary.html")
}

func (this *mainC) Upload_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	partnerId := this.GetPartnerId(ctx)
	r.ParseMultipartForm(20 * 1024 * 1024 * 1024) //20M
	for f := range r.MultipartForm.File {
		w.Write(this.WebCgi.Upload(f, ctx, fmt.Sprintf("%d/item_pic/", partnerId)))
		break
	}
}
