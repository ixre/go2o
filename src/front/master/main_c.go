/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package master

import (
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/mvc"
	"go2o/src/front"
	"go2o/src/core/service/dps"

	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/infrastructure/domain"
)

var _ mvc.Filter = new(mainC)

type mainC struct {
	*baseC
	*front.WebCgi
}

//入口
func (this *mainC) Index(ctx *web.Context) {
	ctx.Response.Write([]byte("<script>location.replace('/dashboard')</script>"))
}

func (this *mainC) Logout(ctx *web.Context) {
	ctx.Session().Destroy()
	ctx.Response.Write([]byte("<script>location.replace('/login')</script>"))
}

//商户首页
func (this *mainC) Dashboard(ctx *web.Context) {
	if this.Requesting(ctx) {
		r, w := ctx.Request, ctx.Response
		d := gof.TemplateDataMap{
			"loginIp": r.Header.Get("USER_ADDRESS"),
		}
		ctx.App.Template().Execute(w, d, "views/master/dashboard.html")
	}
}

//商户汇总页
func (this *mainC) Summary(ctx *web.Context) {
	if this.Requesting(ctx) {
		r, w := ctx.Request, ctx.Response
		d := gof.TemplateDataMap{
			"loginIp": r.Header.Get("USER_ADDRESS"),
		}
		ctx.App.Template().Execute(w, d, "views/master/summary.html")
	}
}

// 导出数据
func (this *mainC) exportData(ctx *web.Context) {
	if this.Requesting(ctx) {
		GetExportData(ctx)
	}
}

func (this *mainC) Upload_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	r.ParseMultipartForm(20 * 1024 * 1024 * 1024) //20M
	for f := range r.MultipartForm.File {
		w.Write(this.WebCgi.Upload(f, ctx, fmt.Sprintf("master/item_pic/")))
		break
	}
}

//登陆
func (this *mainC) Login(ctx *web.Context) {
	if ctx.Request.Method == "POST" {
		this.Login_post(ctx)
	} else {
		ctx.App.Template().Execute(ctx.Response, nil, "views/master/login.html")
	}
}
func (this *mainC) Login_post(ctx *web.Context) {
	r := ctx.Request
	var msg gof.Message
	r.ParseForm()
	usr, pwd := r.Form.Get("uid"), r.Form.Get("pwd")

	if domain.Md5Pwd(pwd, usr) == ctx.App.Config().GetString("webmaster_valid_md5") {
		ctx.Session().Set("master_id", 1)
		if err := ctx.Session().Save(); err != nil {
			msg.Message = err.Error()
		} else {
			msg.Result = true
		}
	} else {
		msg.Message = "用户或密码不正确！"
	}
	ctx.Response.Write(msg.Marshal())
}

//验证登陆
func (pb *mainC) ValidLogin(usr string, pwd string) (*partner.ValuePartner, bool, string) {
	var message string
	var result bool
	var pt *partner.ValuePartner
	var err error

	id := dps.PartnerService.Verify(usr, pwd)

	if id == -1 {
		result = false
		message = "用户或密码不正确！"
	} else {
		pt, err = dps.PartnerService.GetPartner(id)
		if err != nil {
			message = err.Error()
			result = false
		} else {
			result = true
		}
	}
	return pt, result, message
}
