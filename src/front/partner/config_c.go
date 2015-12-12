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
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/mvc"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"time"
)

var _ mvc.Filter = new(configC)

type configC struct {
}

//资料配置
func (this *configC) Profile(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	p, _ := dps.PartnerService.GetPartner(partnerId)
	p.Pwd = ""
	p.ExpiresTime = time.Now().Unix()

	js, _ := json.Marshal(p)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.Render(http.StatusOK, "conf/profile.html", d)
}

func (this *configC) Profile_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	var result gof.Message
	r.ParseForm()

	e := partner.ValuePartner{}
	web.ParseFormToEntity(r.Form, &e)

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

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
		result.Data = id
	}
	return ctx.JSON(http.StatusOK, result)
}

//站点配置
func (this *configC) SiteConf(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	conf := dps.PartnerService.GetSiteConf(partnerId)
	js, _ := json.Marshal(conf)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.Render(http.StatusOK, "conf/site_conf.html", d)
}

func (this *configC) SiteConf_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	var result gof.Message
	r.ParseForm()

	e := partner.SiteConf{}
	web.ParseFormToEntity(r.Form, &e)

	//更新
	origin := dps.PartnerService.GetSiteConf(partnerId)
	e.Host = origin.Host
	e.PartnerId = partnerId

	err := dps.PartnerService.SaveSiteConf(partnerId, &e)

	if err != nil {
		result = gof.Message{Result: false, Message: err.Error()}
	} else {
		result = gof.Message{Result: true, Message: ""}
	}
	return ctx.JSON(http.StatusOK, result)
}

//销售配置
func (this *configC) SaleConf(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	conf := dps.PartnerService.GetSaleConf(partnerId)
	js, _ := json.Marshal(conf)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.Render(http.StatusOK, "conf/sale_conf.html", d)
}

func (this *configC) SaleConf_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	var result gof.Message
	r.ParseForm()

	e := partner.SaleConf{}
	web.ParseFormToEntity(r.Form, &e)

	e.PartnerId = partnerId

	err := dps.PartnerService.SaveSaleConf(partnerId, &e)

	if err != nil {
		result = gof.Message{Result: false, Message: err.Error()}
	} else {
		result = gof.Message{Result: true, Message: ""}
	}
	return ctx.JSON(http.StatusOK, result)
}
