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
	"go2o/src/core/domain/interface/merchant/mss"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// 消息发送控制器
type mssC struct {
}

//邮件模板列表
func (this *mssC) Mail_template_list(ctx *echox.Context) error {
	d := ctx.NewData()
	return ctx.RenderOK("mss.mail_tpl_list.html", d)
}

// 修改广告
func (this *mssC) Edit_mail_tpl(ctx *echox.Context) error {
	merchantId := getMerchantId(ctx)
	id, _ := strconv.Atoi(ctx.Query("id"))
	e, _ := dps.MerchantService.GetMailTemplate(merchantId, id)

	js, _ := json.Marshal(e)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	return ctx.Render(http.StatusOK, "mss.edit_mail_tpl.html", d)
}

// 创建邮箱模板
func (this *mssC) Create_mail_tpl(ctx *echox.Context) error {
	e := mss.MailTemplate{
		Enabled: 1,
	}
	js, _ := json.Marshal(e)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	return ctx.Render(http.StatusOK, "mss.edit_mail_tpl.html", d)
}

// 删除邮件模板(POST)
func (this *mssC) Del_mail_tpl(ctx *echox.Context) error {
	req := ctx.HttpRequest()
	if req.Method == "POST" {

		req.ParseForm()
		var result gof.Message
		merchantId := getMerchantId(ctx)
		adId, _ := strconv.Atoi(req.FormValue("id"))
		err := dps.MerchantService.DeleteMailTemplate(merchantId, adId)

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
	merchantId := getMerchantId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		r.ParseForm()

		var result gof.Message

		e := mss.MailTemplate{}
		web.ParseFormToEntity(r.Form, &e)

		//更新
		e.MerchantId = merchantId

		id, err := dps.MerchantService.SaveMailTemplate(merchantId, &e)

		if err != nil {
			result.Message = err.Error()
		} else {
			result.Result = true
			result.Data = id
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

func (this *mssC) getMailTemplateOpts(merchantId int) string {
	return getMailTemplateOpts(merchantId)
}

// 设置
func (this *mssC) Mss_setting(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.mss_setting_post(ctx)
	}
	merchantId := getMerchantId(ctx)
	e := dps.MerchantService.GetKeyMapsByKeyword(merchantId, "mss_")
	js, _ := json.Marshal(e)

	d := ctx.NewData()
	d.Map["mailTplOpt"] = template.HTML(this.getMailTemplateOpts(merchantId))
	d.Map["entity"] = template.JS(js)

	return ctx.Render(http.StatusOK, "mss.setting.html", d)
}

// 保存设置
func (this *mssC) mss_setting_post(ctx *echox.Context) error {
	var result gof.Message
	merchantId := getMerchantId(ctx)
	req := ctx.HttpRequest()
	req.ParseForm()
	var data map[string]string = make(map[string]string, 0)
	for k, v := range req.Form {
		if strings.HasPrefix(k, "mss_") {
			data[k] = v[0]
		}
	}

	err := dps.MerchantService.SaveKeyMaps(merchantId, data)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	return ctx.JSON(http.StatusOK, result)
}
