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
	"github.com/jsix/gof/web/form"
	"go2o/src/app/util"
	utils "github.com/jsix/gof/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"net/http"
	"net/url"
	"strings"
)

type UserC struct {
}

func (this *UserC) Login(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.login_post(ctx)
	}
	p := getPartner(ctx)
	r := ctx.HttpRequest()
	var tipStyle string
	var returnUrl string = r.URL.Query().Get("return_url")
	if len(returnUrl) == 0 {
		tipStyle = " hidden"
	}

	siteConf := getSiteConf(ctx)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"Partner":  p,
		"Conf":     siteConf,
		"TipStyle": tipStyle,
	}
	return ctx.RenderOK("login.html", d)

}

func (this *UserC) login_post(ctx *echox.Context) error {
	r := ctx.HttpRequest()
	//return ctx.String(http.StatusNotFound,r.FormValue("usr"))
	//r.ParseForm()
	var result gof.Result
	partnerId := GetPartnerId(ctx)
	usr, pwd := r.FormValue("usr"), r.FormValue("pwd")

	pwd = strings.TrimSpace(pwd)
	m, err := dps.MemberService.TryLogin(partnerId, usr, pwd, true)
	if err == nil {
		ctx.Session.Set("member", m)
		err = ctx.Session.Save()
	} else {
		if err != nil {
			result.ErrCode = 1
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 1
			result.ErrMsg = "登陆失败"
		}
	}
	return ctx.JSON(http.StatusOK, result)
}

func (this *UserC) Register(ctx *echox.Context) error {
	p := getPartner(ctx)
	inviCode := ctx.HttpRequest().URL.Query().Get("invi_code")
	siteConf := getSiteConf(ctx)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"Partner":   p,
		"Conf":      siteConf,
		"Invi_code": inviCode,
	}
	return ctx.RenderOK("register.html", d)
}

// 验证用户(POST)
func (this *UserC) ValidUsr(ctx *echox.Context) error {
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		var msg gof.Result
		r.ParseForm()
		usr := r.FormValue("usr")
		err := dps.MemberService.CheckUsr(usr, 0)
		if err == nil {
	msg.ErrCode = 0
		} else {
			msg.ErrMsg = err.Error()
		}
		return ctx.JSON(http.StatusOK, msg)
	}
	return nil
}

// 验证推荐人(POST)
func (this *UserC) Valid_invitation(ctx *echox.Context) error {
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		r.ParseForm()
		msg := gof.Result{}
		code := r.FormValue("invi_code")
		if len(code) > 0 {
			memberId := dps.MemberService.GetMemberIdByInvitationCode(code)
			if memberId <= 0 {
				msg.ErrMsg = "推荐人无效"
				msg.ErrCode = 1
			}
		}
		return ctx.JSON(http.StatusOK, msg)
	}
	return nil
}

// 提交注册信息(POST)
func (this *UserC) PostRegisterInfo(ctx *echox.Context) error {
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		r.ParseForm()
		var result gof.Result
		var member member.ValueMember

		form.ParseEntity(r.Form, &member)
		code := r.FormValue("invi_code")

		if i := strings.Index(r.RemoteAddr, ":"); i != -1 {
			member.RegIp = r.RemoteAddr[:i]
		}

		var memberId int
		var partnerId int
		var err error

		partnerId = GetSessionPartnerId(ctx)
		if len(member.Usr) == 0 || len(member.Pwd) == 0 {
			result.ErrMsg = "1000:注册信息不完整"
			return ctx.JSON(http.StatusOK, result)
		}

		member.Pwd = domain.MemberSha1Pwd(member.Pwd)
		memberId, err = dps.MemberService.RegisterMember(partnerId, &member, "", code)
		if err != nil {
			result.ErrMsg = err.Error()
			result.ErrCode = 1
		} else {
			result.Data = map[string]string{"memberId":utils.Str(memberId)}
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 跳转到会员中心
// url : /user/jump_uc
func (this *UserC) JumpToMCenter(ctx *echox.Context) error {
	returnUrl := ctx.Query("url")
	m := GetMember(ctx)
	var location string
	if m == nil {
		location = "/user/login?return_url=" +
			url.QueryEscape("/user/jump_uc?url="+returnUrl)
	} else {
		location = fmt.Sprintf("http://%s%s/partner_connect?device=%s&sessionId=%s&mid=%d&token=%s&url=%s",
			variable.DOMAIN_PREFIX_MEMBER,
			ctx.App.Config().GetString(variable.ServerDomain),
			util.GetBrownerDevice(ctx.HttpRequest()),
			ctx.Session.GetSessionId(),
			m.Id,
			m.DynamicToken,
			url.QueryEscape(returnUrl),
		)
	}
	return ctx.Redirect(302, location)
}

// 退出
func (this *UserC) Logout(ctx *echox.Context) error {
	ctx.Session.Set("member", nil)
	ctx.Session.Save()
	ctx.HttpResponse().Write([]byte(fmt.Sprintf(`<html><head><title>正在退出...</title></head><body>
			3秒后将自动返回到首页... <br />
			<iframe src="http://%s%s/partner_disconnect" width="0" height="0" frameBorder="0"></iframe>
			<script>window.onload=function(){location.replace('/')}</script></body></html>`,
		variable.DOMAIN_PREFIX_MEMBER,
		ctx.App.Config().GetString(variable.ServerDomain),
	)))
	return nil
}
