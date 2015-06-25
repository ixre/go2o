/**
 * Copyright 2013 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-02-04 20:13
 * description :
 * history :
 */
package ols

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core/service/dps"
	"go2o/src/app/util"
)

type mainC struct {
	*BaseC
}

// 处理跳转
func (this *mainC) HandleIndexGo(ctx *web.Context) bool {
	r, w := ctx.Request, ctx.ResponseWriter
	g := r.URL.Query().Get("go")
	if g == "buy" {
		w.Header().Add("Location", "/list")
		w.WriteHeader(302)
		return true
	}
	return false
}

// 更换访问设备
func (this *mainC) Change_device(ctx *web.Context) {
	deviceType := ctx.Request.URL.Query().Get("device_type")
	util.SetBrownerDevice(ctx, deviceType)
	urlRef := ctx.Request.Referer()
	if len(urlRef) == 0 {
		urlRef = "/"
	}
	ctx.ResponseWriter.Header().Add("Location", urlRef)
	ctx.ResponseWriter.WriteHeader(302)
}

// Member Session Connect
func (this *mainC) Msc(ctx *web.Context){
	form := ctx.Request.URL.Query()
	util.SetDeviceByUrlQuery(ctx,&form)

	ok, memberId := util.MemberHttpSessionConnect(ctx)
	if ok {
		ctx.Items["client_member_id"] = memberId
		rtu := form.Get("return_url")
		if len(rtu) == 0{
			rtu = "/"
		}
		ctx.ResponseWriter.Header().Add("Location",rtu)
		ctx.ResponseWriter.WriteHeader(302)
	}else {
		ctx.ResponseWriter.Write([]byte("not authorized!"))
	}
}




func (this *mainC) Index(ctx *web.Context) {
	if this.BaseC.Requesting(ctx) {
		p := this.BaseC.GetPartner(ctx)
		m := this.BaseC.GetMember(ctx)

		if this.HandleIndexGo(ctx) {
			return
		}

		siteConf := this.BaseC.GetSiteConf(ctx)
		newGoods := dps.SaleService.GetValueGoodsBySaleTag(p.Id, "new-goods", 0, 12)
		hotSales := dps.SaleService.GetValueGoodsBySaleTag(p.Id, "hot-sales", 0, 12)

		this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
			"partner":  p,
			"conf":     siteConf,
			"newGoods": newGoods,
			"hotSales": hotSales,
			"member":   m,
		},
			"views/shop/ols/{device}/index.html",
			"views/shop/ols/{device}/inc/header.html",
			"views/shop/ols/{device}/inc/footer.html")
	}
}
