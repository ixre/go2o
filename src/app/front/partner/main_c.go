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
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/app/front"
	"go2o/src/app/session"
)

type mainC struct {
	gof.App
	*front.WebCgi
}

//入口
func (this *mainC) Index(ctx *web.Context) {
	ctx.ResponseWriter.Write([]byte("<script>location.replace('/dashboard')</script>"))
}

func (this *mainC) Logout(ctx *web.Context) {
	session.GetLSession().PartnerLogout(ctx)
	ctx.ResponseWriter.Write([]byte("<script>location.replace('/login')</script>"))
}

//商户首页
func (this *mainC) Dashboard(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	pt, err := session.GetLSession().GetCurrentSessionFromCookie(r)
	if err != nil {
		ctx.ResponseWriter.Write([]byte("<script>location.replace('/login')</script>"))
		return
	}

	var mf gof.TemplateMapFunc = func(m *map[string]interface{}) {
		(*m)["partner"] = pt
		(*m)["loginIp"] = r.Header.Get("USER_ADDRESS")
	}
	this.App.Template().Render(w, "views/partner/dashboard.html", mf)
}

//商户汇总页
func (this *mainC) Summary(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	pt, err := session.GetLSession().GetCurrentSessionFromCookie(r)
	if err != nil {
		return
	}
	this.App.Template().Render(w,
		"views/partner/summary.html",
		func(m *map[string]interface{}) {
			(*m)["partner"] = pt
			(*m)["loginIp"] = r.Header.Get("USER_ADDRESS")
		})
}

func (this *mainC) Upload_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	ptid, _ := session.GetLSession().GetPartnerIdFromCookie(r)
	r.ParseMultipartForm(20 * 1024 * 1024 * 1024) //20M
	for f, _ := range r.MultipartForm.File {
		w.Write(this.WebCgi.Upload(f, ctx, fmt.Sprintf("%d/item_pic/", ptid)))
		break
	}
}
