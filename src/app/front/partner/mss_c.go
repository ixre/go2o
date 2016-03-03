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
<<<<<<< HEAD
	"go2o/src/core/domain/interface/partner/mss"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
=======
	"github.com/jsix/gof/web/mvc"
	"go2o/src/core/domain/interface/partner/mss"
	"go2o/src/core/service/dps"
	"html/template"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"strconv"
	"strings"
)

<<<<<<< HEAD
// 消息发送控制器
type mssC struct {
}

//邮件模板列表
func (this *mssC) Mail_template_list(ctx *echox.Context) error {
	d := echox.NewRenderData()
	return ctx.RenderOK("mss.mail_tpl_list.html", d)
}

// 修改广告
func (this *mssC) Edit_mail_tpl(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	id, _ := strconv.Atoi(ctx.Query("id"))
	e, _ := dps.PartnerService.GetMailTemplate(partnerId, id)

	js, _ := json.Marshal(e)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.Render(http.StatusOK, "mss.edit_mail_tpl.html", d)
}

// 创建邮箱模板
func (this *mssC) Create_mail_tpl(ctx *echox.Context) error {
=======
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	e := mss.MailTemplate{
		Enabled: 1,
	}
	js, _ := json.Marshal(e)
<<<<<<< HEAD
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.Render(http.StatusOK, "mss.edit_mail_tpl.html", d)
}

// 删除邮件模板(POST)
func (this *mssC) Del_mail_tpl(ctx *echox.Context) error {
	req := ctx.Request()
	if req.Method == "POST" {

		req.ParseForm()
		var result gof.Message
		partnerId := getPartnerId(ctx)
		adId, _ := strconv.Atoi(req.FormValue("id"))
		err := dps.PartnerService.DeleteMailTemplate(partnerId, adId)

		if err != nil {
			result.Message = err.Error()
		} else {
			result.Result = true
		}

		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 保存邮件模板(POST)
func (this *mssC) Save_mail_tpl(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	if r.Method == "POST" {
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
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
=======

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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}

func (this *mssC) getMailTemplateOpts(partnerId int) string {
	return getMailTemplateOpts(partnerId)
}

// 设置
<<<<<<< HEAD
func (this *mssC) Mss_setting(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.mss_setting_post(ctx)
	}
	partnerId := getPartnerId(ctx)
	e := dps.PartnerService.GetKeyMapsByKeyword(partnerId, "mss_")
	js, _ := json.Marshal(e)

	d := echox.NewRenderData()
	d.Map["mailTplOpt"] = template.HTML(this.getMailTemplateOpts(partnerId))
	d.Map["entity"] = template.JS(js)

	return ctx.Render(http.StatusOK, "mss.setting.html", d)
}

// 保存设置
func (this *mssC) mss_setting_post(ctx *echox.Context) error {
	var result gof.Message
	partnerId := getPartnerId(ctx)
	req := ctx.Request()
	req.ParseForm()
	var data map[string]string = make(map[string]string, 0)
	for k, v := range req.Form {
=======
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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
<<<<<<< HEAD
	return ctx.JSON(http.StatusOK, result)
=======
	ctx.Response.JsonOutput(result)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
