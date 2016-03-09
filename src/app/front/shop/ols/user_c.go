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
<<<<<<< HEAD
=======
	"github.com/jsix/gof/web/mvc"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
<<<<<<< HEAD
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
=======
	"strings"
)

var _ mvc.Filter = new(UserC)

type UserC struct {
	*BaseC
}

func (this *UserC) Login(ctx *web.Context) {
	p := this.BaseC.GetPartner(ctx)
	r := ctx.Request
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	var tipStyle string
	var returnUrl string = r.URL.Query().Get("return_url")
	if len(returnUrl) == 0 {
		tipStyle = " hidden"
	}

<<<<<<< HEAD
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
	//return ctx.String(http.StatusNotFound,r.FormValue("usr"))
	//r.ParseForm()
	var result gof.Message
	partnerId := GetPartnerId(r, ctx.Session)
	usr, pwd := r.FormValue("usr"), r.FormValue("pwd")

	pwd = strings.TrimSpace(pwd)
	m, err := dps.MemberService.TryLogin(partnerId, usr, pwd, true)
	if err == nil {
		result.Result = true
		ctx.Session.Set("member", m)
		err = ctx.Session.Save()
=======
	siteConf := this.BaseC.GetSiteConf(ctx)
	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"partner":  p,
		"conf":     siteConf,
		"tipStyle": tipStyle,
	},
		"views/shop/ols/{device}/login.html",
		"views/shop/ols/{device}/inc/header.html",
		"views/shop/ols/{device}/inc/footer.html")

}

func (this *UserC) Login_post(ctx *web.Context) {
	r := ctx.Request
	r.ParseForm()
	var result gof.Message
	partnerId := this.BaseC.GetPartnerId(ctx)
	usr, pwd := r.Form.Get("usr"), r.Form.Get("pwd")

	pwd = strings.TrimSpace(pwd)

	b, m, err := dps.MemberService.Login(partnerId, usr, pwd)

	if b {
		ctx.Session().Set("member", m)
		ctx.Session().Save()
		result.Result = true
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	} else {
		if err != nil {
			result.Message = err.Error()
		} else {
			result.Message = "登陆失败"
		}
	}
<<<<<<< HEAD
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

		member.Pwd = domain.MemberSha1Pwd(member.Pwd)
		memberId, err = dps.MemberService.RegisterMember(partnerId, &member, "", code)
		if err != nil {
			result.Message = err.Error()
		} else {
			result.Result = true
			result.Data = memberId
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
=======
	ctx.Response.JsonOutput(result)
}

func (this *UserC) Register(ctx *web.Context) {
	p := this.BaseC.GetPartner(ctx)
	inviCode := ctx.Request.URL.Query().Get("invi_code")

	siteConf := this.BaseC.GetSiteConf(ctx)
	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"partner":   p,
		"conf":      siteConf,
		"invi_code": inviCode,
	},
		"views/shop/ols/{device}/register.html",
		"views/shop/ols/{device}/inc/header.html",
		"views/shop/ols/{device}/inc/footer.html")

}

func (this *UserC) ValidUsr_post(ctx *web.Context) {
	r := ctx.Request
	var msg gof.Message
	r.ParseForm()
	usr := r.FormValue("usr")
	err := dps.MemberService.CheckUsr(usr, 0)
	if err == nil {
		msg.Result = true
	} else {
		msg.Message = err.Error()
	}
	ctx.Response.JsonOutput(msg)
}

func (this *UserC) Valid_invitation_post(ctx *web.Context) {
	var result gof.Message = gof.Message{Result: true}
	ctx.Request.ParseForm()
	code := ctx.Request.FormValue("invi_code")
	if len(code) > 0 {
		memberId := dps.MemberService.GetMemberIdByInvitationCode(code)
		if memberId <= 0 {
			result.Result = false
			result.Message = "推荐人无效"
		}
	}
	ctx.Response.JsonOutput(result)
}

func (this *UserC) PostRegisterInfo_post(ctx *web.Context) {
	ctx.Request.ParseForm()
	var result gof.Message
	var member member.ValueMember

	web.ParseFormToEntity(ctx.Request.Form, &member)
	code := ctx.Request.FormValue("invi_code")

	if i := strings.Index(ctx.Request.RemoteAddr, ":"); i != -1 {
		member.RegIp = ctx.Request.RemoteAddr[:i]
	}

	var memberId int
	var partnerId int
	var err error

	partnerId = this.GetPartnerId(ctx)
	if len(member.Usr) == 0 || len(member.Pwd) == 0 {
		result.Message = "1000:注册信息不完整"
		ctx.Response.JsonOutput(result)
		return
	}

	if err = dps.PartnerService.CheckRegisterMode(partnerId, code); err != nil {
		result.Message = err.Error()
		ctx.Response.JsonOutput(result)
		return
	}

	var invId int
	if len(code) > 0 {
		invId = dps.MemberService.GetMemberIdByInvitationCode(code)
		if invId <= 0 {
			result.Message = "1011：推荐码不正确"
			ctx.Response.JsonOutput(result)
			return
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
	ctx.Response.JsonOutput(result)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}

// 跳转到会员中心
// url : /user/jump_m
<<<<<<< HEAD
func (this *UserC) JumpToMCenter(ctx *echox.Context) error {
	m := GetMember(ctx)
=======
func (this *UserC) JumpToMCenter(ctx *web.Context) {
	w := ctx.Response
	m := this.BaseC.GetMember(ctx)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	var location string
	if m == nil {
		location = "/user/login?return_url=/user/jump_m"
	} else {
<<<<<<< HEAD
		location = fmt.Sprintf("http://%s%s/partner_connect?device=%s&sessionId=%s&mid=%d&token=%s",
			variable.DOMAIN_PREFIX_MEMBER,
			ctx.App.Config().GetString(variable.ServerDomain),
			util.GetBrownerDevice(ctx.Request()),
			ctx.Session.GetSessionId(),
=======
		location = fmt.Sprintf("http://%s.%s/login/partner_connect?device=%s&sessionId=%s&mid=%d&token=%s",
			variable.DOMAIN_PREFIX_MEMBER,
			ctx.App.Config().GetString(variable.ServerDomain),
			util.GetBrownerDevice(ctx),
			ctx.Session().GetSessionId(),
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
			m.Id,
			m.DynamicToken,
		)
	}
<<<<<<< HEAD
	return ctx.Redirect(302, location)
}

// 退出
func (this *UserC) Logout(ctx *echox.Context) error {
	ctx.Session.Set("member", nil)
	ctx.Session.Save()
	ctx.Response().Write([]byte(fmt.Sprintf(`<html><head><title>正在退出...</title></head><body>
			3秒后将自动返回到首页... <br />
			<iframe src="http://%s%s/partner_disconnect" width="0" height="0" frameBorder="0"></iframe>
=======
	w.Header().Add("Location", location)
	w.WriteHeader(302)
}

// 退出
func (this *UserC) Logout(ctx *web.Context) {
	ctx.Session().Set("member", nil)
	ctx.Session().Save()
	ctx.Response.Write([]byte(fmt.Sprintf(`<html><head><title>正在退出...</title></head><body>
			3秒后将自动返回到首页... <br />
			<iframe src="http://%s.%s/login/partner_disconnect" width="0" height="0" frameBorder="0"></iframe>
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
			<script>window.onload=function(){location.replace('/')}</script></body></html>`,
		variable.DOMAIN_PREFIX_MEMBER,
		ctx.App.Config().GetString(variable.ServerDomain),
	)))
<<<<<<< HEAD
	return nil
=======
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
