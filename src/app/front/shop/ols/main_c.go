/**
 * Copyright 2013 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-02-04 20:13
 * description :
 * history :
 */
package ols

import (
	"github.com/jsix/gof"
<<<<<<< HEAD
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
=======
	"github.com/jsix/gof/web"
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"strings"
)

type MainC struct {
<<<<<<< HEAD
}

// 处理跳转
func (this *MainC) HandleIndexGo(ctx *echox.Context) bool {
	r, w := ctx.Request(), ctx.Response()
=======
	*BaseC
}

// 处理跳转
func (this *MainC) HandleIndexGo(ctx *web.Context) bool {
	r, w := ctx.Request, ctx.Response
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	g := r.URL.Query().Get("go")
	if g == "buy" {
		w.Header().Add("Location", "/list")
		w.WriteHeader(302)
		return true
	}
	return false
}

// 更换访问设备
<<<<<<< HEAD
func (this *MainC) change_device(ctx *echox.Context) error {
	r := ctx.Request()
	util.SetDeviceByUrlQuery(ctx.Response(), r)
	toUrl := r.URL.Query().Get("return_url")
	if len(toUrl) == 0 {
		toUrl = r.Referer()
=======
func (this *MainC) Change_device(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	util.SetDeviceByUrlQuery(ctx, &form)

	toUrl := form.Get("return_url")
	if len(toUrl) == 0 {
		toUrl = ctx.Request.Referer()
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		if len(toUrl) == 0 {
			toUrl = "/"
		}
	}
<<<<<<< HEAD
	return ctx.Redirect(302, toUrl)
}

// Member session connect
func (this *MainC) Msc(ctx *echox.Context) error {
	form := ctx.Request().URL.Query()
	util.SetDeviceByUrlQuery(ctx.Response(), ctx.Request())
	util.MemberHttpSessionConnect(ctx, func(memberId int) {
		v := ctx.Session.Get("member")
=======

	ctx.Response.Header().Add("Location", toUrl)
	ctx.Response.WriteHeader(302)
}

// Member session connect
func (this *MainC) Msc(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	util.SetDeviceByUrlQuery(ctx, &form)

	ok, memberId := util.MemberHttpSessionConnect(ctx, func(memberId int) {
		v := ctx.Session().Get("member")
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		var m *member.ValueMember
		if v != nil {
			m = v.(*member.ValueMember)
			if m.Id != memberId { // 如果会话冲突
				m = nil
			}
		}

		if m == nil {
			m = dps.MemberService.GetMember(memberId)
<<<<<<< HEAD
			ctx.Session.Set("member", m)
			ctx.Session.Save()
		}
	})

=======
			ctx.Session().Set("member", m)
			ctx.Session().Save()
		}
	})

	if ok {
		ctx.Items["client_member_id"] = memberId

	}

>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	rtu := form.Get("return_url")
	if len(rtu) == 0 {
		rtu = "/"
	}
<<<<<<< HEAD
	return ctx.Redirect(302, rtu)
}

// Member session disconnect
func (this *MainC) Msd(ctx *echox.Context) error {
	if util.MemberHttpSessionDisconnect(ctx) {
		ctx.Session.Set("member", nil)
		ctx.Session.Save()
		return ctx.StringOK("disconnect success")
	}
	return ctx.StringOK("disconnect fail")
}

func (this *MainC) T(ctx *echox.Context) error {
	path := ctx.Request().URL.Path
	var i int = strings.LastIndex(path, "/")
	ivCode := path[i+1:]
	ctx.Response().Header().Add("Location", "/user/register?invi_code="+
		ivCode+"&"+ctx.Request().URL.RawQuery)
	ctx.Response().WriteHeader(302)
	return nil
}

func (this *MainC) Index(ctx *echox.Context) error {
	p := getPartner(ctx)
	m := GetMember(ctx)

	if this.HandleIndexGo(ctx) {
		return nil
	}

	siteConf := getSiteConf(ctx)
	newGoods := dps.SaleService.GetValueGoodsBySaleTag(p.Id, "new-goods", 0, 12)
	hotSales := dps.SaleService.GetValueGoodsBySaleTag(p.Id, "hot-sales", 0, 12)

	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":  p,
		"conf":     siteConf,
		"newGoods": newGoods,
		"hotSales": hotSales,
		"member":   m,
	}
	return ctx.RenderOK("index.html", d)
}

func (this *MainC) App(ctx *echox.Context) error {
	p := getPartner(ctx)
	m := GetMember(ctx)
	siteConf := getSiteConf(ctx)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner": p,
		"conf":    siteConf,
		"member":  m,
	}
	return ctx.RenderOK("app.html", d)
=======
	ctx.Response.Header().Add("Location", rtu)
	ctx.Response.WriteHeader(302)
}

// Member session disconnect
func (this *MainC) Msd(ctx *web.Context) {
	if util.MemberHttpSessionDisconnect(ctx) {
		ctx.Session().Set("member", nil)
		ctx.Session().Save()
		ctx.Response.Write([]byte("disconnect success"))
	} else {
		ctx.Response.Write([]byte("disconnect fail"))
	}
}

func (this *MainC) T(ctx *web.Context) {
	path := ctx.Request.URL.Path
	var i int = strings.LastIndex(path, "/")
	ivCode := path[i+1:]
	ctx.Response.Header().Add("Location", "/user/register.htm?invi_code="+ivCode+"&"+ctx.Request.URL.RawQuery)
	ctx.Response.WriteHeader(302)
}

func (this *MainC) Index(ctx *web.Context) {
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

func (this *MainC) App(ctx *web.Context) {
	if this.BaseC.Requesting(ctx) {
		p := this.BaseC.GetPartner(ctx)
		m := this.BaseC.GetMember(ctx)
		siteConf := this.BaseC.GetSiteConf(ctx)

		this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
			"partner": p,
			"conf":    siteConf,
			"member":  m,
		},
			"views/shop/ols/{device}/app.html",
			"views/shop/ols/{device}/inc/header.html",
			"views/shop/ols/{device}/inc/footer.html")
	}
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
