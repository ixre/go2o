/**
 * Copyright 2015 @ S1N1 Team.
 * name : default_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ols

import (
	"encoding/json"
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"net/url"
	"strings"
)

type BaseC struct {
	*BaseC
}

func (this *BaseC) Requesting(ctx *web.Context) bool {
	// 商户不存在
	partnerId := getPartnerId(ctx)
	if partnerId <= 0 {
		renderError(ctx, true, "No such partner.")
		return false
	}

	ctx.Items["partner_id"] = partnerId // 缓存PartnerId

	// 校验商户并缓存
	pt := cache.GetValuePartnerCache(partnerId)
	ctx.Items["partner_ins"] = pt

	// 判断线上商店开通情况
	var conf = cache.GetPartnerSiteConf(partnerId)
	if conf == nil {
		renderError(ctx, true, "线上商店未开通")
		return false
	}

	if conf.State == enum.PARTNER_SITE_CLOSED {
		if strings.TrimSpace(conf.StateHtml) == "" {
			conf.StateHtml = "网站暂停访问，请联系商家：" + pt.Tel
		}
		renderError(ctx, true, conf.StateHtml)
		return false
	}

	ctx.Items["partner_siteconf"] = conf

	return true
}

func (this *BaseC) RequestEnd(*web.Context) {
}

// 获取商户编号
func (this *BaseC) GetPartnerId(ctx *web.Context) int {
	return ctx.GetItem("partner_id").(int)
}
func (this *BaseC) GetPartner(ctx *web.Context) *partner.ValuePartner {
	return ctx.GetItem("partner_ins").(*partner.ValuePartner)
}

func (this *BaseC) GetSiteConf(ctx *web.Context) *partner.SiteConf {
	return ctx.GetItem("partner_siteconf").(*partner.SiteConf)
}

// 获取商户API信息
func (this *BaseC) GetPartnerApi(ctx *web.Context) *partner.ApiInfo {
	return dps.PartnerService.GetApiInfo(getPartnerId(ctx))
}

// 获取会员
func (this *BaseC) GetMember(ctx *web.Context) *member.ValueMember {
	memberIdObj := ctx.Session().Get("member")
	if memberIdObj != nil {
		if o, ok := memberIdObj.(*member.ValueMember); ok {
			return o
		}
	}
	return nil
}

// 检查会员是否登陆
func (this *BaseC) CheckMemberLogin(ctx *web.Context) bool {
	if ctx.Session().Get("member") == nil {
		ctx.ResponseWriter.Header().Add("Location", "/user/login?return_url="+
			url.QueryEscape(ctx.Request.RequestURI))
		ctx.ResponseWriter.WriteHeader(302)
		return false
	}
	return true
}

func renderError(ctx *web.Context, simpleError bool, message string) {
	if simpleError {
		const errTpl string = "<html><body><h1 style='color:red'>%s</h1></body></html>"
		ctx.ResponseWriter.Write([]byte(fmt.Sprintf(errTpl, message)))
	} else {
		//todo: 用模板显示错误
	}
}
func getPartnerId(ctx *web.Context) int {
	//return 104
	currHost := ctx.Request.Host
	host := ctx.Session().Get("webui_host")
	pid := ctx.Session().Get("webui_pid")
	if host == nil || pid == nil || host != currHost {
		partnerId := dps.PartnerService.GetPartnerIdByHost(currHost)
		if partnerId != -1 {
			ctx.Session().Set("webui_host", currHost)
			ctx.Session().Set("webui_pid", partnerId)
			ctx.Session().Save()
		}
		return partnerId
	}
	return pid.(int)
}

// 输出Json
func (this *BaseC) JsonOutput(ctx *web.Context, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		this.ErrorOutput(ctx, err.Error())
	} else {
		ctx.ResponseWriter.Write(b)
	}
}

// 输出错误信息
func (this *BaseC) ErrorOutput(ctx *web.Context, err string) {
	ctx.ResponseWriter.Write([]byte("{error:\"" + err + "\"}"))
}

// 输出错误信息
func (this *BaseC) ResultOutput(ctx *web.Context, result gof.Message) {
	ctx.ResponseWriter.Write([]byte(fmt.Sprintf(
		"{result:%v,code:%d,message:\"%s\"}", result.Result, result.Code, result.Message)))
}

// 执行模板
func (this *BaseC) ExecuteTemplate(ctx *web.Context, d gof.TemplateDataMap, files ...string) {
	newFiles := make([]string, len(files))
	for i, v := range files {
		newFiles[i] = strings.Replace(v, "{device}", ctx.Items["device_view_dir"].(string), -1)
	}
	ctx.App.Template().Execute(ctx.ResponseWriter, d, newFiles...)
}
