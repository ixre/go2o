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
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"net/http"
	"strings"
)

type UserC struct {
}

func (this *UserC) Login(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.login_post(ctx)
	}
	p := getPartner(ctx)
	r := ctx.Request()
	var tipStyle string
	var returnUrl string = r.URL.Query().Get("return_url")
	if len(returnUrl) == 0 {
		tipStyle = " hidden"
	}

	siteConf := getSiteConf(ctx)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":  p,
		"conf":     siteConf,
		"tipStyle": tipStyle,
	}
	return ctx.RenderOK("login.html", d)

}

func (this *UserC) login_post(ctx *echox.Context) error {
	r := ctx.Request()
	r.ParseForm()
	var result gof.Message
	partnerId := GetPartnerId(r, ctx.Session)
	usr, pwd := r.FormValue("usr"), r.FormValue("pwd")

	pwd = strings.TrimSpace(pwd)

	b, m, err := dps.MemberService.Login(partnerId, usr, pwd)

	if b {
		ctx.Session.Set("member", m)
		ctx.Session.Save()
		result.Result = true
	} else {
		if err != nil {
			result.Message = err.Error()
		} else {
			result.Message = "登陆失败"
		}
	}

	return ctx.JSON(http.StatusOK, result)
}

func (this *UserC) Register(ctx *echox.Context) error {
	p := getPartner(ctx)
	inviCode := ctx.Request().URL.Query().Get("invi_code")
	siteConf := getSiteConf(ctx)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":   p,
		"conf":      siteConf,
		"invi_code": inviCode,
	}
	return ctx.RenderOK("register.html", d)
}

// 验证用户(POST)
func (this *UserC) ValidUsr(ctx *echox.Context) error {
	r := ctx.Request()
	if r.Method == "POST" {
		var msg gof.Message
		r.ParseForm()
		usr := r.FormValue("usr")
		err := dps.MemberService.CheckUsr(usr, 0)
		if err == nil {
			msg.Result = true
		} else {
			msg.Message = err.Error()
		}
		return ctx.JSON(http.StatusOK, msg)
	}
	return nil
}

// 验证推荐人(POST)
func (this *UserC) Valid_invitation(ctx *echox.Context) error {
	r := ctx.Request()
	if r.Method == "POST" {
		r.ParseForm()
		msg := gof.Message{Result: true}
		code := r.FormValue("invi_code")
		if len(code) > 0 {
			memberId := dps.MemberService.GetMemberIdByInvitationCode(code)
			if memberId <= 0 {
				msg.Result = false
				msg.Message = "推荐人无效"
			}
		}
		return ctx.JSON(http.StatusOK, msg)
	}
	return nil
}

// 提交注册信息(POST)
func (this *UserC) PostRegisterInfo(ctx *echox.Context) error {
	r := ctx.Request()
	if r.Method == "POST" {
		r.ParseForm()
		var result gof.Message
		var member member.ValueMember

		web.ParseFormToEntity(r.Form, &member)
		code := r.FormValue("invi_code")

		if i := strings.Index(r.RemoteAddr, ":"); i != -1 {
			member.RegIp = r.RemoteAddr[:i]
		}

		var memberId int
		var partnerId int
		var err error

		partnerId = GetSessionPartnerId(ctx)
		if len(member.Usr) == 0 || len(member.Pwd) == 0 {
			result.Message = "1000:注册信息不完整"
			return ctx.JSON(http.StatusOK, result)
		}

		if err = dps.PartnerService.CheckRegisterMode(partnerId, code); err != nil {
			result.Message = err.Error()
			return ctx.JSON(http.StatusOK, result)
		}

		var invId int
		if len(code) > 0 {
			invId = dps.MemberService.GetMemberIdByInvitationCode(code)
			if invId <= 0 {
				result.Message = "1011：推荐码不正确"
				return ctx.JSON(http.StatusOK, result)
			}
		}

		member.Pwd = domain.MemberSha1Pwd(member.Pwd)
		memberId, err = dps.MemberService.SaveMember(&member)
		if err == nil {
			err = dps.MemberService.SaveRelation(memberId, "", invId, partnerId)
		}

		if err != nil {
			result.Message = err.Error()
		} else {
			result.Result = true
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 跳转到会员中心
// url : /user/jump_m
func (this *UserC) JumpToMCenter(ctx *echox.Context) error {
	m := GetMember(ctx)
	var location string
	if m == nil {
		location = "/user/login?return_url=/user/jump_m"
	} else {
		location = fmt.Sprintf("http://%s%s/partner_connect?device=%s&sessionId=%s&mid=%d&token=%s",
			variable.DOMAIN_PREFIX_MEMBER,
			ctx.App.Config().GetString(variable.ServerDomain),
			util.GetBrownerDevice(ctx.Request()),
			ctx.Session.GetSessionId(),
			m.Id,
			m.DynamicToken,
		)
	}
	return ctx.Redirect(302, location)
}

// 退出
func (this *UserC) Logout(ctx *echox.Context) error {
	ctx.Session.Set("member", nil)
	ctx.Session.Save()
	ctx.Response().Write([]byte(fmt.Sprintf(`<html><head><title>正在退出...</title></head><body>
			3秒后将自动返回到首页... <br />
			<iframe src="http://%s%s/partner_disconnect" width="0" height="0" frameBorder="0"></iframe>
			<script>window.onload=function(){location.replace('/')}</script></body></html>`,
		variable.DOMAIN_PREFIX_MEMBER,
		ctx.App.Config().GetString(variable.ServerDomain),
	)))
	return nil
}
