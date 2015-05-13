/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package master

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
	ctx.ResponseWriter.Write([]byte("<script>location.replace('/dashboard')</script>"))
}

func (this *mainC) Logout(ctx *web.Context) {
	ctx.Session().Destroy()
	ctx.ResponseWriter.Write([]byte("<script>location.replace('/login')</script>"))
}

//商户首页
func (this *mainC) Dashboard(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter

	var mf gof.TemplateMapFunc = func(m *map[string]interface{}) {
		(*m)["loginIp"] = r.Header.Get("USER_ADDRESS")
	}
	ctx.App.Template().Render(w, "views/master/dashboard.html", mf)
}

//商户汇总页
func (this *mainC) Summary(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter

	ctx.App.Template().Render(w,
		"views/master/summary.html",
		func(m *map[string]interface{}) {
			(*m)["loginIp"] = r.Header.Get("USER_ADDRESS")
		})
}

func (this *mainC) Upload_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseMultipartForm(20 * 1024 * 1024 * 1024) //20M
	for f := range r.MultipartForm.File {
		w.Write(this.WebCgi.Upload(f, ctx, fmt.Sprintf("master/item_pic/")))
		break
	}
}
