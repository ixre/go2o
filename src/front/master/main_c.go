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
	"go2o/src/core/infrastructure/domain"
	"go2o/src/front"
	"go2o/src/x/echox"
	"net/http"
)

type mainC struct {
	*front.WebCgi
}

//入口
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
	r.ParseMultipartForm(20 * 1024 * 1024 * 1024) //20M
	for f := range r.MultipartForm.File {
		w.Write(this.WebCgi.Upload(f, ctx, fmt.Sprintf("master/item_pic/")))
		break
	}
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
	var msg gof.Message
	r.ParseForm()
	usr, pwd := r.Form.Get("uid"), r.Form.Get("pwd")

	if domain.Md5Pwd(pwd, usr) == ctx.App.Config().GetString("webmaster_valid_md5") {
		ctx.Session.Set("is_master", 1)
		if err := ctx.Session.Save(); err != nil {
			msg.Message = err.Error()
		} else {
			msg.Result = true
		}
	} else {
		msg.Message = "用户或密码不正确！"
	}
	return ctx.JSON(http.StatusOK, msg)
}
