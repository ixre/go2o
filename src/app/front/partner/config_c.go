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
	"encoding/json"
	"github.com/jsix/gof"
	//"github.com/jsix/gof/web"
	"go2o/src/app/cache"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"time"
	"github.com/jsix/gof/web/form"
	"fmt"
)

type configC struct {
}

//资料配置
func (this *configC) Profile(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.profile_post(ctx)
	}
	partnerId := getPartnerId(ctx)
	p, _ := dps.PartnerService.GetPartner(partnerId)
	p.Pwd = ""
	p.ExpiresTime = time.Now().Unix()

	js, _ := json.Marshal(p)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("conf.profile.html", d)
}

func (this *configC) profile_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	var result gof.Result
	r.ParseForm()

	e := partner.ValuePartner{}
	form.ParseEntity(r.Form, &e)

	//更新
	origin, _ := dps.PartnerService.GetPartner(partnerId)
	e.ExpiresTime = origin.ExpiresTime
	e.JoinTime = origin.JoinTime
	e.LastLoginTime = origin.LastLoginTime
	e.LoginTime = origin.LoginTime
	e.Pwd = origin.Pwd
	e.UpdateTime = time.Now().Unix()
	e.Id = partnerId

	id, err := dps.PartnerService.SavePartner(partnerId, &e)
	result.Error(err)
	if err == nil {
		cache.DelPartnerCache(partnerId)
		var data = make(map[string]string)
		data["id"] = fmt.Sprintf("%d", id)
		result.Data = data
	}
	return ctx.JSON(http.StatusOK, result)
}

//站点配置
func (this *configC) SiteConf(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.siteConf_post(ctx)
	}
	partnerId := getPartnerId(ctx)
	conf := dps.PartnerService.GetSiteConf(partnerId)
	js, _ := json.Marshal(conf)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	d.Map["Logo"] = format.GetResUrl(conf.Logo)
	return ctx.RenderOK("conf.site_conf.html", d)
}

func (this *configC) siteConf_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	var result gof.Result
	r.ParseForm()

	e := partner.SiteConf{}
	form.ParseEntity(r.Form, &e)

	//更新
	origin := dps.PartnerService.GetSiteConf(partnerId)
	e.Host = origin.Host
	e.PartnerId = partnerId

	err := dps.PartnerService.SaveSiteConf(partnerId, &e)
	result.Error(err)
	if err == nil {
		cache.DelPartnerCache(partnerId)
	}
	return ctx.JSON(http.StatusOK, result)
}

//销售配置
func (this *configC) SaleConf(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.saleConf_post(ctx)
	}
	partnerId := getPartnerId(ctx)
	conf := dps.PartnerService.GetSaleConf(partnerId)
	js, _ := json.Marshal(conf)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("conf.sale_conf.html", d)
}

func (this *configC) saleConf_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	var result gof.Result
	r.ParseForm()

	e := partner.SaleConf{}
	form.ParseEntity(r.Form, &e)

	e.PartnerId = partnerId

	err := dps.PartnerService.SaveSaleConf(partnerId, &e)
	result.Error(err)
	if err == nil {
		cache.DelPartnerCache(partnerId)
	}
	return ctx.JSON(http.StatusOK, result)
}
