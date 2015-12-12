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
	"go2o/src/x/echox"
	"net/http"
)

var _ mvc.Filter = new(adC)

// 消息发送控制器
type mssC struct {
}

//邮件模板列表
func (this *mssC) Mail_template_list(ctx *echox.Context)error{
	d := echox.NewRenderData()
	return ctx.Render(http.StatusOK, "mss/mail_tpl_list.html", d)
}

// 修改广告
func (this *mssC) Edit_mail_tpl(ctx *echox.Context)error{
	partnerId := getPartnerId(ctx)
	form := ctx.Request.URL.Query()
	id, _ := strconv.Atoi(form.Get("id"))
	e, _ := dps.PartnerService.GetMailTemplate(partnerId, id)

	js, _ := json.Marshal(e)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.Render(http.StatusOK,"mss/edit_mail_tpl.html",d)
}

// 创建邮箱模板
func (this *mssC) Create_mail_tpl(ctx *echox.Context)error{
	e := mss.MailTemplate{
		Enabled: 1,
	}
	js, _ := json.Marshal(e)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.Render(http.StatusOK,"mss/edit_mail_tpl.html",d)
}

// 删除邮件模板
func (this *mssC) Del_mail_tpl_post(ctx *echox.Context)error{
	ctx.Request.ParseForm()
	form := ctx.Request.Form
	var result gof.Message
	partnerId := getPartnerId(ctx)
	adId, _ := strconv.Atoi(form.Get("id"))
	err := dps.PartnerService.DeleteMailTemplate(partnerId, adId)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}

	return ctx.JSON(http.StatusOK,result)
}

// 保存邮件模板
func (this *mssC) Save_mail_tpl_post(ctx *echox.Context)error{
	partnerId := getPartnerId(ctx)
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
	return ctx.JSON(http.StatusOK,result)
}

func (this *mssC) getMailTemplateOpts(partnerId int) string {
	return getMailTemplateOpts(partnerId)
}

// 设置
func (this *mssC) Mss_setting(ctx *echox.Context)error{
	partnerId := getPartnerId(ctx)
	e := dps.PartnerService.GetKeyMapsByKeyword(partnerId, "mss_")
	js, _ := json.Marshal(e)

	d := echox.NewRenderData()
	d.Map["mailTplOpt"] = template.HTML(this.getMailTemplateOpts(partnerId))
	d.Map["entity"]= template.JS(js)

	return ctx.Render(http.StatusOK,"mss/mss_setting.html",d)
}

// 保存设置
func (this *mssC) Mss_setting_post(ctx *echox.Context)error{
	var result gof.Message
	partnerId := getPartnerId(ctx)
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
	return ctx.JSON(http.StatusOK,result)
}
