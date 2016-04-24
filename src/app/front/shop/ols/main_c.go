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
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"strings"
)

type MainC struct {
}

// 处理跳转
func (this *MainC) HandleIndexGo(ctx *echox.Context) bool {
	r, w := ctx.Request(), ctx.Response()
	g := r.URL.Query().Get("go")
	if g == "buy" {
		w.Header().Add("Location", "/list")
		w.WriteHeader(302)
		return true
	}
	return false
}

// 更换访问设备
func (this *MainC) change_device(ctx *echox.Context) error {
	r := ctx.Request()
	util.SetDeviceByUrlQuery(ctx.Response(), r)
	toUrl := r.URL.Query().Get("return_url")
	if len(toUrl) == 0 {
		toUrl = r.Referer()
		if len(toUrl) == 0 {
			toUrl = "/"
		}
	}
	return ctx.Redirect(302, toUrl)
}

// Member session connect
func (this *MainC) Msc(ctx *echox.Context) error {
	form := ctx.Request().URL.Query()
	util.SetDeviceByUrlQuery(ctx.Response(), ctx.Request())
	util.MemberHttpSessionConnect(ctx, func(memberId int) {
		v := ctx.Session.Get("member")
		var m *member.ValueMember
		if v != nil {
			m = v.(*member.ValueMember)
			if m.Id != memberId { // 如果会话冲突
				m = nil
			}
		}

		if m == nil {
			m = dps.MemberService.GetMember(memberId)
			ctx.Session.Set("member", m)
			ctx.Session.Save()
		}
	})

	rtu := form.Get("return_url")
	if len(rtu) == 0 {
		rtu = "/"
	}
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

func (this *MainC) T(c *echox.Context) error {
	path := c.Request().URL.Path
	var i int = strings.LastIndex(path, "/")
	ivCode := path[i+1:]
	device := c.Query("device")
	if len(device) > 0 {
		util.SetBrownerDevice(c.Response(), c.Request(), device) //设置访问设备
	}
	c.Response().Header().Add("Location", "/user/register?invi_code="+
		ivCode+"&"+c.Request().URL.RawQuery)
	c.Response().WriteHeader(302)
	return nil
}

func (this *MainC) Index(ctx *echox.Context) error {
	p := getPartner(ctx)
	m := GetMember(ctx)

	if this.HandleIndexGo(ctx) {
		return nil
	}

	siteConf := getSiteConf(ctx)
	newGoods := dps.SaleService.GetValueGoodsBySaleTag(p.Id, "new-goods", "", 0, 12)
	hotSales := dps.SaleService.GetValueGoodsBySaleTag(p.Id, "hot-sales", "", 0, 12)

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

func (this *MainC) MallEntry(ctx *echox.Context) error {
	p := getPartner(ctx)
	m := GetMember(ctx)

	if this.HandleIndexGo(ctx) {
		return nil
	}

	siteConf := getSiteConf(ctx)
	newGoods := dps.SaleService.GetValueGoodsBySaleTag(p.Id, "new-goods", "", 0, 12)
	hotSales := dps.SaleService.GetValueGoodsBySaleTag(p.Id, "hot-sales", "", 0, 12)

	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":  p,
		"conf":     siteConf,
		"newGoods": newGoods,
		"hotSales": hotSales,
		"member":   m,
	}
	return ctx.RenderOK("mall_entry.html", d)
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
}
