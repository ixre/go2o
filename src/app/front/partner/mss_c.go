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
	"go2o/src/core/domain/interface/partner/mss"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"github.com/jsix/gof/web/form"
	"fmt"
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
	partnerId := getPartnerId(ctx)
	id, _ := strconv.Atoi(ctx.Query("id"))
	e, _ := dps.PartnerService.GetMailTemplate(partnerId, id)

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
		var result gof.Result
		partnerId := getPartnerId(ctx)
		adId, _ := strconv.Atoi(req.FormValue("id"))
		err := dps.PartnerService.DeleteMailTemplate(partnerId, adId)

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
		}

		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 保存邮件模板(POST)
func (this *mssC) Save_mail_tpl(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		r.ParseForm()

		var result gof.Result

		e := mss.MailTemplate{}
		form.ParseEntity(r.Form, &e)

		//更新
		e.PartnerId = partnerId

		id, err := dps.PartnerService.SaveMailTemplate(partnerId, &e)

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
			var data = make(map[string]string)
			data["id"] = fmt.Sprintf("%d", id)
			result.Data = data
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

func (this *mssC) getMailTemplateOpts(partnerId int) string {
	return getMailTemplateOpts(partnerId)
}

// 设置
func (this *mssC) Mss_setting(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.mss_setting_post(ctx)
	}
	partnerId := getPartnerId(ctx)
	e := dps.PartnerService.GetKeyMapsByKeyword(partnerId, "mss_")
	js, _ := json.Marshal(e)

	d := ctx.NewData()
	d.Map["mailTplOpt"] = template.HTML(this.getMailTemplateOpts(partnerId))
	d.Map["entity"] = template.JS(js)

	return ctx.Render(http.StatusOK, "mss.setting.html", d)
}

// 保存设置
func (this *mssC) mss_setting_post(ctx *echox.Context) error {
	var result gof.Result
	partnerId := getPartnerId(ctx)
	req := ctx.HttpRequest()
	req.ParseForm()
	var data map[string]string = make(map[string]string, 0)
	for k, v := range req.Form {
		if strings.HasPrefix(k, "mss_") {
			data[k] = v[0]
		}
	}

	err := dps.PartnerService.SaveKeyMaps(partnerId, data)

	if err != nil {
		result.ErrMsg = err.Error()
	} else {
		result.ErrCode = 0
	}
	return ctx.JSON(http.StatusOK, result)
}
