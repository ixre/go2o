/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"encoding/json"
	"github.com/jrsix/gof"
	"github.com/jrsix/gof/web"
	"github.com/jrsix/gof/web/mvc"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"html/template"
	"time"
)

var _ mvc.Filter = new(configC)

type configC struct {
	*baseC
}

//资料配置
func (this *configC) Profile(ctx *web.Context) {

	partnerId := this.GetPartnerId(ctx)
	p, _ := dps.PartnerService.GetPartner(partnerId)
	p.Pwd = ""
	p.ExpiresTime = time.Now().Unix()

	js, _ := json.Marshal(p)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity": template.JS(js),
		},
		"views/partner/conf/profile.html")
}

func (this *configC) Profile_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
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
	w.Write(result.Marshal())
}

//站点配置
func (this *configC) SiteConf(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	conf := dps.PartnerService.GetSiteConf(partnerId)
	js, _ := json.Marshal(conf)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity": template.JS(js),
		},
		"views/partner/conf/site_conf.html")
}

func (this *configC) SiteConf_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
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
	w.Write(result.Marshal())
}

//销售配置
func (this *configC) SaleConf(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	conf := dps.PartnerService.GetSaleConf(partnerId)
	js, _ := json.Marshal(conf)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity": template.JS(js),
		},
		"views/partner/conf/sale_conf.html")
}

func (this *configC) SaleConf_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
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
	w.Write(result.Marshal())
}
