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
	"encoding/json"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"go2o/src/app/cache"
	"go2o/src/core/domain/interface/merchant"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"time"
)

type configC struct {
}

//资料配置
func (this *configC) Profile(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.profile_post(ctx)
	}
	merchantId := getMerchantId(ctx)
	p, _ := dps.MerchantService.GetMerchant(merchantId)
	p.Pwd = ""
	p.ExpiresTime = time.Now().Unix()

	js, _ := json.Marshal(p)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("conf.profile.html", d)
}

func (this *configC) profile_post(ctx *echox.Context) error {
	merchantId := getMerchantId(ctx)
	r := ctx.HttpRequest()
	var result gof.Message
	r.ParseForm()

	e := merchant.MerchantValue{}
	web.ParseFormToEntity(r.Form, &e)

	//更新
	origin, _ := dps.MerchantService.GetMerchant(merchantId)
	e.ExpiresTime = origin.ExpiresTime
	e.JoinTime = origin.JoinTime
	e.LastLoginTime = origin.LastLoginTime
	e.LoginTime = origin.LoginTime
	e.Pwd = origin.Pwd
	e.UpdateTime = time.Now().Unix()
	e.Id = merchantId

	id, err := dps.MerchantService.SaveMerchant(merchantId, &e)
	result.Error(err)
	if err == nil {
		cache.DelMerchantCache(merchantId)
		result.Data = id
	}
	return ctx.JSON(http.StatusOK, result)
}

//站点配置
func (this *configC) SiteConf(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.siteConf_post(ctx)
	}
	merchantId := getMerchantId(ctx)
	conf := dps.MerchantService.GetSiteConf(merchantId)
	js, _ := json.Marshal(conf)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	d.Map["Logo"] = format.GetResUrl(conf.Logo)
	return ctx.RenderOK("conf.site_conf.html", d)
}

func (this *configC) siteConf_post(ctx *echox.Context) error {
	merchantId := getMerchantId(ctx)
	r := ctx.HttpRequest()
	var result gof.Message
	r.ParseForm()

	e := merchant.SiteConf{}
	web.ParseFormToEntity(r.Form, &e)

	//更新
	origin := dps.MerchantService.GetSiteConf(merchantId)
	e.Host = origin.Host
	e.MerchantId = merchantId

	err := dps.MerchantService.SaveSiteConf(merchantId, &e)
	result.Error(err)
	if err == nil {
		cache.DelMerchantCache(merchantId)
	}
	return ctx.JSON(http.StatusOK, result)
}

//销售配置
func (this *configC) SaleConf(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.saleConf_post(ctx)
	}
	merchantId := getMerchantId(ctx)
	conf := dps.MerchantService.GetSaleConf(merchantId)
	js, _ := json.Marshal(conf)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("conf.sale_conf.html", d)
}

func (this *configC) saleConf_post(ctx *echox.Context) error {
	merchantId := getMerchantId(ctx)
	r := ctx.HttpRequest()
	var result gof.Message
	r.ParseForm()

	e := merchant.SaleConf{}
	web.ParseFormToEntity(r.Form, &e)

	e.MerchantId = merchantId

	err := dps.MerchantService.SaveSaleConf(merchantId, &e)
	result.Error(err)
	if err == nil {
		cache.DelMerchantCache(merchantId)
	}
	return ctx.JSON(http.StatusOK, result)
}
