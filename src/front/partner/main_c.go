/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/mvc"
	"go2o/src/front"
	"github.com/labstack/echo"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
)

var _ mvc.Filter = new(mainC)

type mainC struct {
	*baseC
	*front.WebCgi
}

//入口
func (this *mainC) Index(ctx *echo.Context)(err error) {

	_,err = ctx.Response().Write([]byte("<script>location.replace('/main/dashboard')</script>"))

	//todo:??
//	if this.baseC.Requesting(ctx) {
//		ctx.Response.Write([]byte("<script>location.replace('/main/dashboard')</script>"))
//	}
//	this.baseC.RequestEnd(ctx)
	return err
}

func (this *mainC) Logout(ctx *PartnerContext)error {
//	ctx.Session().Destroy()
//	ctx.Response.Write([]byte("<script>location.replace('/login')</script>"))
	return nil
}

//商户首页
func (this *mainC) Dashboard(ctx *PartnerContext)error {
	pt, _ := dps.PartnerService.GetPartner(ctx.PartnerId)

	dm := echox.NewRendData()
	dm.Data = gof.TemplateDataMap{
		"partner": pt,
		"loginIp": ctx.Echo.Request().Header.Get("USER_ADDRESS"),
	}
	return ctx.Echo.Render(200,"dashboard.html",dm)
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
