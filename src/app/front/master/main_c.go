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
	"errors"
	"fmt"
	"github.com/jsix/gof"
	"go2o/src/app/front"
	"go2o/src/core/infrastructure/domain"
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
	ctx.HttpResponse().Write([]byte("<script>location.replace('/login')</script>"))
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
	ctx.HttpResponse().Header().Set("Content-Type", "application/json")
	ctx.HttpResponse().Write(GetExportData(ctx.HttpRequest()))
	return nil
}

func (this *mainC) Upload_post(ctx *echox.Context) error {
	r, w := ctx.HttpRequest(), ctx.HttpResponse()
	r.ParseMultipartForm(64 * 1024 * 1024 * 1024) //64M
	for f := range r.MultipartForm.File {
		w.Write(this.WebCgi.Upload(f, ctx, fmt.Sprintf("master/item_pic/")))
		break
	}
	return nil
}

//登陆
func (this *mainC) Login(ctx *echox.Context) error {
	if ctx.HttpRequest().Method == "POST" {
		return this.login_post(ctx)
	}
	return ctx.RenderOK("login.html", ctx.NewData())
}
func (this *mainC) login_post(ctx *echox.Context) error {
	var msg gof.Message
	ctx.HttpRequest().ParseForm()
	usr, pwd := ctx.HttpRequest().FormValue("uid"), ctx.HttpRequest().FormValue("pwd")
	var err error
	if domain.Md5Pwd(pwd, usr) == ctx.App.Config().GetString("webmaster_valid_md5") {
		ctx.Session.Set("is_master", 1)
		err = ctx.Session.Save()
	} else {
		err = errors.New("用户或密码不正确！")
	}
	return ctx.JSON(http.StatusOK, msg.Error(err))
}
