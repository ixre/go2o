package partner

import (
	"com/ording/dao"
	"com/ording/entity"
	"encoding/json"
	"html/template"
	"net/http"
	"ops/cf"
	"ops/cf/app"
	"ops/cf/web"
	"time"
	"com/domain/interface/partner"
	"com/ording/dproxy"
)

type configC struct {
	app.Context
}

//资料配置
func (this *configC) Profile(w http.ResponseWriter, r *http.Request, partnerId int) {
	p := dao.Partner().GetPartnerById(partnerId)
	p.Pwd = ""
	p.Secret = ""
	p.Expires = time.Now()

	js, _ := json.Marshal(p)

	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(js)
		},
		"views/partner/conf/profile.html")
}

func (this *configC) Profile_post(w http.ResponseWriter, r *http.Request, partnerId int) {
	var result cf.JsonResult
	r.ParseForm()

	e := entity.Partner{}
	web.ParseFormToEntity(r.Form, &e)

	//更新
	origin := dao.Partner().GetPartnerById(partnerId)
	e.Expires = origin.Expires
	e.Secret = origin.Secret
	e.JoinTime = origin.JoinTime
	e.LastLoginTime = origin.LastLoginTime
	e.LoginTime = origin.LoginTime
	e.Pwd = origin.Pwd
	e.UpdateTime = time.Now()
	e.Id = partnerId

	id, err := dao.Partner().SavePartner(&e)

	if err != nil {
		result = cf.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = cf.JsonResult{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}

//站点配置
func (this *configC) SiteConf(w http.ResponseWriter, r *http.Request, partnerId int) {
	conf := dproxy.PartnerService.GetSiteConf(partnerId)
	js, _ := json.Marshal(conf)

	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(js)
		},
		"views/partner/conf/site_conf.html")
}

func (this *configC) SiteConf_post(w http.ResponseWriter, r *http.Request, partnerId int) {
	var result cf.JsonResult
	r.ParseForm()

	e := partner.SiteConf{}
	web.ParseFormToEntity(r.Form, &e)

	//更新
	origin :=dproxy.PartnerService.GetSiteConf(partnerId)
	e.Host = origin.Host
	e.PtId = partnerId

	err := dproxy.PartnerService.SaveSiteConf(partnerId,&e)

	if err != nil {
		result = cf.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = cf.JsonResult{Result: true, Message: ""}
	}
	w.Write(result.Marshal())
}

//销售配置
func (this *configC) SaleConf(w http.ResponseWriter, r *http.Request, partnerId int) {
	conf := dproxy.PartnerService.GetSaleConf(partnerId)
	js, _ := json.Marshal(conf)

	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(js)
		},
		"views/partner/conf/sale_conf.html")
}

func (this *configC) SaleConf_post(w http.ResponseWriter, r *http.Request, partnerId int) {
	var result cf.JsonResult
	r.ParseForm()

	e := partner.SaleConf{}
	web.ParseFormToEntity(r.Form, &e)

	e.PtId = partnerId

	err := dproxy.PartnerService.SaveSaleConf(partnerId,&e)

	if err != nil {
		result = cf.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = cf.JsonResult{Result: true, Message: ""}
	}
	w.Write(result.Marshal())
}
