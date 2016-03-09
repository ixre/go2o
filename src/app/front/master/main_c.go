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
<<<<<<< HEAD
	"go2o/src/app/front"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/x/echox"
	"net/http"
)

type mainC struct {
=======
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/mvc"
	"go2o/src/app/front"
	"go2o/src/core/service/dps"

	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/infrastructure/domain"
)

var _ mvc.Filter = new(mainC)

type mainC struct {
	*baseC
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	*front.WebCgi
}

//入口
<<<<<<< HEAD
func (this *mainC) Index(ctx *echox.Context) error {
	return ctx.StringOK("<script>location.replace('/dashboard')</script>")
}

func (this *mainC) Logout(ctx *echox.Context) error {
	ctx.Session.Destroy()
	ctx.Response().Write([]byte("<script>location.replace('/login')</script>"))
	return nil
}

//商户首页
func (this *mainC) Dashboard(ctx *echox.Context) error {
	d := ctx.NewData()
	d.Map["loginIp"] = ctx.Request().Header.Get("USER_ADDRESS")
	return ctx.RenderOK("dashboard.html", d)
}

//商户汇总页
func (this *mainC) Summary(ctx *echox.Context) error {
	d := ctx.NewData()
	d.Map["loginIp"] = ctx.Request().Header.Get("USER_ADDRESS")
	return ctx.RenderOK("summary.html", d)
}

// 导出数据
func (this *mainC) exportData(ctx *echox.Context) error {
	ctx.Response().Header().Set("Content-Type", "application/json")
	ctx.Response().Write(GetExportData(ctx.Request()))
	return nil
}

func (this *mainC) Upload_post(ctx *echox.Context) error {
	r, w := ctx.Request(), ctx.Response()
=======
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	r.ParseMultipartForm(20 * 1024 * 1024 * 1024) //20M
	for f := range r.MultipartForm.File {
		w.Write(this.WebCgi.Upload(f, ctx, fmt.Sprintf("master/item_pic/")))
		break
	}
<<<<<<< HEAD
	return nil
}

//登陆
func (this *mainC) Login(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.login_post(ctx)
	}
	return ctx.RenderOK("login.html", ctx.NewData())
}
func (this *mainC) login_post(ctx *echox.Context) error {
	r := ctx.Request()
=======
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	var msg gof.Message
	r.ParseForm()
	usr, pwd := r.Form.Get("uid"), r.Form.Get("pwd")

	if domain.Md5Pwd(pwd, usr) == ctx.App.Config().GetString("webmaster_valid_md5") {
<<<<<<< HEAD
		ctx.Session.Set("is_master", 1)
		if err := ctx.Session.Save(); err != nil {
=======
		ctx.Session().Set("master_id", 1)
		if err := ctx.Session().Save(); err != nil {
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
			msg.Message = err.Error()
		} else {
			msg.Result = true
		}
	} else {
		msg.Message = "用户或密码不正确！"
	}
<<<<<<< HEAD
	return ctx.JSON(http.StatusOK, msg)
=======
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
