/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"html/template"
	"time"
)

type configC struct {
}

//资料配置
func (this *configC) Profile(ctx *web.Context, partnerId int) {
	p, _ := dps.PartnerService.GetPartner(partnerId)
	p.Pwd = ""
	p.Secret = ""
	p.ExpiresTime = time.Now().Unix()

	js, _ := json.Marshal(p)

	ctx.App.Template().Execute(ctx.ResponseWriter,
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(js)
		},
		"views/partner/conf/profile.html")
}

func (this *configC) Profile_post(ctx *web.Context, partnerId int) {
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.JsonResult
	r.ParseForm()

	e := partner.ValuePartner{}
	web.ParseFormToEntity(r.Form, &e)

	//更新
	origin, _ := dps.PartnerService.GetPartner(partnerId)
	e.ExpiresTime = origin.ExpiresTime
	e.Secret = origin.Secret
	e.JoinTime = origin.JoinTime
	e.LastLoginTime = origin.LastLoginTime
	e.LoginTime = origin.LoginTime
	e.Pwd = origin.Pwd
	e.UpdateTime = time.Now().Unix()
	e.Id = partnerId

	id, err := dps.PartnerService.SavePartner(partnerId, &e)

	if err != nil {
		result = gof.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}

//站点配置
func (this *configC) SiteConf(ctx *web.Context, partnerId int) {
	conf := dps.PartnerService.GetSiteConf(partnerId)
	js, _ := json.Marshal(conf)

	ctx.App.Template().Execute(ctx.ResponseWriter,
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(js)
		},
		"views/partner/conf/site_conf.html")
}

func (this *configC) SiteConf_post(ctx *web.Context, partnerId int) {
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.JsonResult
	r.ParseForm()

	e := partner.SiteConf{}
	web.ParseFormToEntity(r.Form, &e)

	//更新
	origin := dps.PartnerService.GetSiteConf(partnerId)
	e.Host = origin.Host
	e.PartnerId = partnerId

	err := dps.PartnerService.SaveSiteConf(partnerId, &e)

	if err != nil {
		result = gof.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true, Message: ""}
	}
	w.Write(result.Marshal())
}

//销售配置
func (this *configC) SaleConf(ctx *web.Context, partnerId int) {
	conf := dps.PartnerService.GetSaleConf(partnerId)
	js, _ := json.Marshal(conf)

	ctx.App.Template().Execute(ctx.ResponseWriter,
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(js)
		},
		"views/partner/conf/sale_conf.html")
}

func (this *configC) SaleConf_post(ctx *web.Context, partnerId int) {
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.JsonResult
	r.ParseForm()

	e := partner.SaleConf{}
	web.ParseFormToEntity(r.Form, &e)

	e.PartnerId = partnerId

	err := dps.PartnerService.SaveSaleConf(partnerId, &e)

	if err != nil {
		result = gof.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true, Message: ""}
	}
	w.Write(result.Marshal())
}
