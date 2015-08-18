/**
 * Copyright 2015 @ z3q.net.
 * name : default_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ucenter

import (
	"github.com/jrsix/gof"
	"github.com/jrsix/gof/web"
	"github.com/jrsix/gof/web/mvc"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"html/template"
	"net/url"
	"strings"
)

var _ mvc.Filter = new(baseC)

type baseC struct {
}

func (this *baseC) Requesting(ctx *web.Context) bool {
	//验证是否登陆
	s := ctx.Session().Get("member")
	if s != nil {
		if m := s.(*member.ValueMember); m != nil {
			ctx.Items["member"] = m
			return true
		}
	}
	ctx.Response.Write([]byte("<script>window.parent.location.href='/login?return_url=" +
		url.QueryEscape(ctx.Request.URL.String()) + "'</script>"))
	return false
}

func (this *baseC) RequestEnd(ctx *web.Context) {
}

// 获取商户
func (this *baseC) GetPartner(ctx *web.Context) *partner.ValuePartner {
	val := ctx.Session().Get("member:rel_partner")
	if val != nil {
		return cache.GetValuePartnerCache(val.(int))
	} else {
		m := this.GetMember(ctx)
		if m != nil {
			rel := dps.MemberService.GetRelation(m.Id)
			ctx.Session().Set("member:rel_partner", rel.RegisterPartnerId)
			ctx.Session().Save()
			return cache.GetValuePartnerCache(rel.RegisterPartnerId)
		}
	}
	return nil
}

// 获取会员
func (this *baseC) GetMember(ctx *web.Context) *member.ValueMember {
	return ctx.Items["member"].(*member.ValueMember)
}

// 重新缓存会员
func (this *baseC) ReCacheMember(ctx *web.Context, memberId int) {
	v := dps.MemberService.GetMember(memberId)
	ctx.Session().Set("member", v)
	ctx.Session().Save()
}

// 获取商户的站点设置
func (this *baseC) GetSiteConf(partnerId int) *partner.SiteConf {
	return dps.PartnerService.GetSiteConf(partnerId)
}

// 输出错误信息
func (this *baseC) errorOutput(ctx *web.Context, err string) {
	ctx.Response.Write([]byte("{error:\"" + err + "\"}"))
}

func executeTemplate(ctx *web.Context, funcMap template.FuncMap, dataMap gof.TemplateDataMap, files ...string) {
	newFiles := make([]string, len(files))
	for i, v := range files {
		newFiles[i] = strings.Replace(v, "{device}", ctx.Items["device_view_dir"].(string), -1)
	}
	if funcMap == nil {
		ctx.App.Template().Execute(ctx.Response, dataMap, newFiles...)
	} else {
		ctx.App.Template().ExecuteWithFunc(ctx.Response, funcMap, dataMap, newFiles...)
	}
}

// 执行模板
func (this *baseC) ExecuteTemplate(ctx *web.Context, d gof.TemplateDataMap, files ...string) {
	executeTemplate(ctx, nil, d, files...)
}

func (this *baseC) ExecuteTemplateWithFunc(ctx *web.Context, funcMap template.FuncMap,
	dataMap gof.TemplateDataMap, files ...string) {
	executeTemplate(ctx, funcMap, dataMap, files...)
}
