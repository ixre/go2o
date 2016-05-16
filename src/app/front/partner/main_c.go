/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package merchant

import (
	"fmt"
	"github.com/jsix/gof"
	"go2o/src/app/front"
	"go2o/src/core/domain/interface/merchant"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"gopkg.in/labstack/echo.v1"
	"net/http"
	"strings"
)

type mainC struct {
	*front.WebCgi
}

//入口
func (this *mainC) Index(ctx *echo.Context) (err error) {

	_, err = ctx.Response().Write([]byte("<script>location.replace('/main/dashboard')</script>"))

	//todo:??
	//	if this.baseC.Requesting(ctx) {
	//		ctx.Response.Write([]byte("<script>location.replace('/main/dashboard')</script>"))
	//	}
	//	this.baseC.RequestEnd(ctx)
	return err
}

//登陆
func (this *mainC) Login(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.Login_post(ctx)
	}
	d := ctx.NewData()
	return ctx.RenderOK("login.html", d)
}

func (this *mainC) Login_post(ctx *echox.Context) error {
	r, w := ctx.HttpRequest(), ctx.HttpResponse()
	r.ParseForm()
	usr, pwd := r.Form.Get("uid"), r.Form.Get("pwd")

	pwd = strings.TrimSpace(pwd)
	pt, result, message := this.ValidLogin(usr, pwd)

	if result {
		ctx.Session.Set("merchant_id", pt.Id)
		if err := ctx.Session.Save(); err != nil {
			result = false
			message = err.Error()
		}
	}

	if result {
		w.Write([]byte("{result:true}"))
	} else {
		w.Write([]byte("{result:false,message:'" + message + "'}"))
	}
	return nil
}

//验证登陆
func (pb *mainC) ValidLogin(usr string, pwd string) (*merchant.MerchantValue, bool, string) {
	var message string
	var result bool
	var pt *merchant.MerchantValue
	var err error

	id := dps.PartnerService.Verify(usr, pwd)

	if id == -1 {
		result = false
		message = "用户或密码不正确！"
	} else {
		pt, err = dps.PartnerService.GetMerchant(id)
		if err != nil {
			message = err.Error()
			result = false
		} else {
			result = true
		}
	}
	return pt, result, message
}

func (this *mainC) Logout(ctx *echox.Context) error {
	ctx.Session.Destroy()
	ctx.HttpResponse().Write([]byte("<script>location.replace('/login')</script>"))
	return nil
}

//商户首页
func (this *mainC) Dashboard(ctx *echox.Context) error {
	pt, _ := dps.PartnerService.GetMerchant(getMerchantId(ctx))

	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":   pt,
		"loginIp":   ctx.Request().Header.Get("USER_ADDRESS"),
		"AliasGrow": variable.AliasGrowAccount,
	}
	return ctx.Render(200, "dashboard.html", d)
}

//商户汇总页
func (this *mainC) Summary(ctx *echox.Context) error {
	r := ctx.HttpRequest()
	pt, _ := dps.PartnerService.GetMerchant(getMerchantId(ctx))
	d := ctx.NewData()
	d.Map["partner"] = pt
	d.Map["loginIp"] = r.Header.Get("USER_ADDRESS")

	return ctx.Render(http.StatusOK, "summary.html", d)
}

// 导出数据
func (this *mainC) exportData(ctx *echox.Context) error {
	ctx.HttpResponse().Header().Set("Content-Type", "application/json")
	ctx.HttpResponse().Write(GetExportData(ctx.HttpRequest(), getMerchantId(ctx)))
	return nil
}

func (this *mainC) Upload_post(ctx *echox.Context) error {
	req := ctx.HttpRequest()
	partnerId := getMerchantId(ctx)
	req.ParseMultipartForm(20 * 1024 * 1024 * 1024) //20M
	for f := range req.MultipartForm.File {
		ctx.HttpResponse().Write(this.WebCgi.Upload(f, ctx, fmt.Sprintf("%d/item_pic/", partnerId)))
	}
	return nil
}

func (this *mainC) GeoLocation(ctx *echox.Context) error {
	this.WebCgi.GeoLocation(ctx)
	return nil
}

//地区Json
//func (this *mainC) ChinaJson(ctx *echox.Context)error{
//	var node *tree.TreeNode = dao.Common().GetChinaTree()
//	json, _ := json.Marshal(node)
//	w.Write(json)
//}
