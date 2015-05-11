/**
 * Copyright 2015 @ S1N1 Team.
 * name : user_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package www

import (
	"fmt"
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/goclient"
	"go2o/src/core/variable"
	"strings"
)

type userC struct {
	*baseC
}

func (this *userC) Login(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	p := this.GetPartner(ctx)
	var tipStyle string
	var returnUrl string = r.URL.Query().Get("return_url")
	if len(returnUrl) == 0 {
		tipStyle = " hidden"
	}

	pa := this.GetPartnerApi(ctx)

	if b, siteConf := GetSiteConf(w, p, pa); b {
		ctx.App.Template().Execute(w, func(m *map[string]interface{}) {
			mv := *m
			mv["partner"] = p
			mv["title"] = "会员登录－" + siteConf.SubTitle
			mv["conf"] = siteConf
			mv["tipStyle"] = tipStyle
		},
			"views/web/www/login.html",
			"views/web/www/inc/header.html",
			"views/web/www/inc/footer.html")
	}
}

func (this *userC) Login_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
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
	_, w := ctx.Request, ctx.ResponseWriter
	p := this.GetPartner(ctx)
	pa := this.GetPartnerApi(ctx)
	if b, siteConf := GetSiteConf(w, p, pa); b {
		ctx.App.Template().Execute(w, func(m *map[string]interface{}) {
			(*m)["partner"] = p
			(*m)["title"] = "会员注册－" + siteConf.SubTitle
			(*m)["conf"] = siteConf
		},
			"views/web/www/register.html",
			"views/web/www/inc/header.html",
			"views/web/www/inc/footer.html")
	}
}

func (this *userC) ValidUsr_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
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

func (this *userC) PostRegistInfo_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
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
// url : /user/g2m
func (this *userC) member(ctx *web.Context) {
	m := this.GetMember(ctx)
	var location string
	if m == nil {
		location = "/login?return_url=/member"
	} else {
		location = fmt.Sprintf("http://%s.%s/login/partner_connect?sessionId=%s&mid=%d&token=%s",
			variable.DOMAIN_MEMBER_PREFIX,
			ctx.App.Config().GetString(variable.ServerDomain),
			ctx.Session().GetSessionId(),
			m.Id,
			m.DynamicToken,
		)
	}
	ctx.ResponseWriter.Write([]byte("<script>window.parent.location.replace('" + location + "')</script>"))
}

//退出
func (this *userC) Logout(ctx *web.Context) {
	ctx.Session().Set("member", nil)
	ctx.Session().Save()
	ctx.ResponseWriter.Write([]byte(fmt.Sprintf(`<html><head><title>正在退出...</title></head><body>
			3秒后将自动返回到首页... <br />
			<iframe src="http://%s.%s/login/partner_disconnect" width="0" height="0" frameBorder="0"></iframe>
			<script>window.onload=function(){location.replace('/')}</script></body></html>`,
		variable.DOMAIN_MEMBER_PREFIX,
		ctx.App.Config().GetString(variable.ServerDomain),
	)))
}
