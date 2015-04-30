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
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/app/front"
)

var _ mvc.Filter = new(mainC)

type mainC struct {
	Base *baseC
	*front.WebCgi
}

func (this *mainC) Requesting(ctx *web.Context) bool {
	return this.Base.Requesting(ctx)
}
func (this *mainC) RequestEnd(ctx *web.Context) {
	this.Base.RequestEnd(ctx)
}

//入口
func (this *mainC) Index(ctx *web.Context) {
	ctx.ResponseWriter.Write([]byte("<script>location.replace('/dashboard')</script>"))
}

func (this *mainC) Logout(ctx *web.Context) {
	ctx.Session().Destroy()
	ctx.ResponseWriter.Write([]byte("<script>location.replace('/login')</script>"))
}

//商户首页
func (this *mainC) Dashboard(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	pt, _ := this.Base.GetPartner(ctx)

	var mf gof.TemplateMapFunc = func(m *map[string]interface{}) {
		(*m)["partner"] = pt
		(*m)["loginIp"] = r.Header.Get("USER_ADDRESS")
	}
	ctx.App.Template().Render(w, "views/partner/dashboard.html", mf)
}

//商户汇总页
func (this *mainC) Summary(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	pt, _ := this.Base.GetPartner(ctx)

	ctx.App.Template().Render(w,
		"views/partner/summary.html",
		func(m *map[string]interface{}) {
			(*m)["partner"] = pt
			(*m)["loginIp"] = r.Header.Get("USER_ADDRESS")
		})
}

func (this *mainC) Upload_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	partnerId := this.Base.GetPartnerId(ctx)
	r.ParseMultipartForm(20 * 1024 * 1024 * 1024) //20M
	for f := range r.MultipartForm.File {
		w.Write(this.WebCgi.Upload(f, ctx, fmt.Sprintf("%d/item_pic/", partnerId)))
		break
	}
}
