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
	"go2o/src/core/domain/interface/partner/mss"
	"go2o/src/core/service/dps"
	"html/template"
	"strconv"
	"strings"
)

var _ mvc.Filter = new(adC)

// 消息发送控制器
type mssC struct {
	*baseC
}

//邮件模板列表
func (this *mssC) Mail_template_list(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.Response, gof.TemplateDataMap{}, "views/partner/mss/mail_tpl_list.html")
}

// 修改广告
func (this *mssC) Edit_mail_tpl(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	form := ctx.Request.URL.Query()
	id, _ := strconv.Atoi(form.Get("id"))
	e, _ := dps.PartnerService.GetMailTemplate(partnerId, id)

	js, _ := json.Marshal(e)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity": template.JS(js),
		},
		"views/partner/mss/edit_mail_tpl.html")
}

// 创建邮箱模板
func (this *mssC) Create_mail_tpl(ctx *web.Context) {
	e := mss.MailTemplate{
		Enabled: 1,
	}
	js, _ := json.Marshal(e)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity": template.JS(js),
		},
		"views/partner/mss/edit_mail_tpl.html")
}

// 删除邮件模板
func (this *mssC) Del_mail_tpl_post(ctx *web.Context) {
	ctx.Request.ParseForm()
	form := ctx.Request.Form
	var result gof.Message
	partnerId := this.GetPartnerId(ctx)
	adId, _ := strconv.Atoi(form.Get("id"))
	err := dps.PartnerService.DeleteMailTemplate(partnerId, adId)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}

	ctx.Response.JsonOutput(result)
}

// 保存邮件模板
func (this *mssC) Save_mail_tpl_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r := ctx.Request
	r.ParseForm()

	var result gof.Message

	e := mss.MailTemplate{}
	web.ParseFormToEntity(r.Form, &e)

	//更新
	e.PartnerId = partnerId

	id, err := dps.PartnerService.SaveMailTemplate(partnerId, &e)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
		result.Data = id
	}
	ctx.Response.JsonOutput(result)
}

func (this *mssC) getMailTemplateOpts(partnerId int) string {
	return getMailTemplateOpts(partnerId)
}

// 设置
func (this *mssC) Mss_setting(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	e := dps.PartnerService.GetKeyMapsByKeyword(partnerId, "mss_")
	js, _ := json.Marshal(e)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"mailTplOpt": template.HTML(this.getMailTemplateOpts(partnerId)),
			"entity":     template.JS(js),
		},
		"views/partner/mss/mss_setting.html")
}

// 保存设置
func (this *mssC) Mss_setting_post(ctx *web.Context) {
	var result gof.Message
	partnerId := this.GetPartnerId(ctx)
	ctx.Request.ParseForm()
	var data map[string]string = make(map[string]string, 0)
	for k, v := range ctx.Request.Form {
		if strings.HasPrefix(k, "mss_") {
			data[k] = v[0]
		}
	}

	err := dps.PartnerService.SaveKeyMaps(partnerId, data)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	ctx.Response.JsonOutput(result)
}
