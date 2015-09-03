/**
 * Copyright 2015 @ z3q.net.
 * name : user_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ols

import (
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/goclient"
	"go2o/src/core/variable"
	"strings"
)

type userC struct {
	*baseC
}

func (this *userC) Login(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	p := this.GetPartner(ctx)
	var tipStyle string
	var returnUrl string = r.URL.Query().Get("return_url")
	if len(returnUrl) == 0 {
		tipStyle = " hidden"
	}

	pa := this.GetPartnerApi(ctx)

	if b, siteConf := GetSiteConf(w, p, pa); b {
		ctx.App.Template().Execute(w, gof.TemplateDataMap{
			"partner":  p,
			"conf":     siteConf,
			"tipStyle": tipStyle,
		},
			"views/shop/ols/{device}/login.html",
			"views/shop/ols/{device}/inc/header.html",
			"views/shop/ols/{device}/inc/footer.html")
	}
}

func (this *userC) Login_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	r.ParseForm()
	usr, pwd := r.Form.Get("usr"), r.Form.Get("pwd")
	result, _ := goclient.Member.Login(usr, pwd)

	if result.Result {
		ctx.Session().Set("member", result.Member)
		ctx.Session().Save()
		w.Write([]byte("{result:true}"))
		return
	}
	w.Write([]byte("{result:false,message:'" + result.Message + "'}"))
}

func (this *userC) Register(ctx *web.Context) {
	_, w := ctx.Request, ctx.Response
	p := this.GetPartner(ctx)
	pa := this.GetPartnerApi(ctx)
	if b, siteConf := GetSiteConf(w, p, pa); b {
		ctx.App.Template().Execute(w, gof.TemplateDataMap{
			"partner": p,
			"conf":    siteConf,
		},
			"views/shop/ols/{device}/register.html",
			"views/shop/ols/{device}/inc/header.html",
			"views/shop/ols/{device}/inc/footer.html")
	}
}

func (this *userC) ValidUsr_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	p := this.GetPartner(ctx)
	pa := this.GetPartnerApi(ctx)
	r.ParseForm()
	usr := r.FormValue("usr")
	b := goclient.Partner.UserIsExist(p.Id, pa.ApiSecret, usr)
	if !b {
		w.Write([]byte(`{"result":true}`))
	} else {
		w.Write([]byte(`{"result":false}`))
	}
}

func (this *userC) PostRegisterInfo_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	p := this.GetPartner(ctx)
	r.ParseForm()
	var member member.ValueMember
	web.ParseFormToEntity(r.Form, &member)
	if i := strings.Index(r.RemoteAddr, ":"); i != -1 {
		member.RegIp = r.RemoteAddr[:i]
	}
	b, err := goclient.Partner.RegisterMember(&member, p.Id, 0, "")
	if b {
		w.Write([]byte(`{"result":true}`))
	} else {
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"result":false,"message":"%s"}`, err.Error())))
		} else {
			w.Write([]byte(`{"result":false}`))
		}

	}
}

//跳转到会员中心
// url : /user/jump_m
func (this *userC) member(ctx *web.Context) {
	m := this.GetMember(ctx)
	var location string
	if m == nil {
		location = "/login?return_url=/member"
	} else {
		location = fmt.Sprintf("http://%s.%s/login/partner_connect?sessionId=%s&mid=%d&token=%s",
			variable.DOMAIN_PREFIX_MEMBER,
			ctx.App.Config().GetString(variable.ServerDomain),
			ctx.Session().GetSessionId(),
			m.Id,
			m.DynamicToken,
		)
	}
	ctx.Response.Write([]byte("<script>window.parent.location.replace('" + location + "')</script>"))
}

//退出
func (this *userC) Logout(ctx *web.Context) {
	ctx.Session().Set("member", nil)
	ctx.Session().Save()
	ctx.Response.Write([]byte(fmt.Sprintf(`<html><head><title>正在退出...</title></head><body>
			3秒后将自动返回到首页... <br />
			<iframe src="http://%s.%s/login/partner_disconnect" width="0" height="0" frameBorder="0"></iframe>
			<script>window.onload=function(){location.replace('/')}</script></body></html>`,
		variable.DOMAIN_PREFIX_MEMBER,
		ctx.App.Config().GetString(variable.ServerDomain),
	)))
}
