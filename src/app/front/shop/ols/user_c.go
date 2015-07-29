/**
 * Copyright 2015 @ S1N1 Team.
 * name : user_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ols

import (
	"errors"
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"strings"
)

var _ mvc.Filter = new(UserC)

type UserC struct {
	*BaseC
}

func (this *UserC) Login(ctx *web.Context) {
	p := this.BaseC.GetPartner(ctx)
	r := ctx.Request
	var tipStyle string
	var returnUrl string = r.URL.Query().Get("return_url")
	if len(returnUrl) == 0 {
		tipStyle = " hidden"
	}

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
	b, m, err := dps.MemberService.Login(partnerId, usr, pwd)

	if b {
		ctx.Session().Set("member", m)
		ctx.Session().Save()
		result.Result = true
	} else {
		if err != nil {
			result.Message = err.Error()
		} else {
			result.Message = "登陆失败"
		}
	}
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
	ctx.Request.ParseForm()
	memberId := dps.MemberService.GetMemberIdByInvitationCode(ctx.Request.FormValue("invi_code"))
	var result gof.Message
	result.Result = memberId != 0

	if !result.Result {
		result.Message = "推荐人无效"
	}

	ctx.Response.JsonOutput(result)
}

func (this *UserC) PostRegisterInfo_post(ctx *web.Context) {
	ctx.Request.ParseForm()
	var result gof.Message
	var member member.ValueMember
	web.ParseFormToEntity(ctx.Request.Form, &member)
	if i := strings.Index(ctx.Request.RemoteAddr, ":"); i != -1 {
		member.RegIp = ctx.Request.RemoteAddr[:i]
	}

	var memberId int
	var err error

	if len(member.Usr) == 0 || len(member.Pwd) == 0 {
		err = errors.New("注册信息不完整")
	} else {
		member.Pwd = domain.Md5MemberPwd(member.Pwd)
		memberId, err = dps.MemberService.SaveMember(&member)
		if err == nil {
			invId := dps.MemberService.GetMemberIdByInvitationCode(ctx.Request.FormValue("invi_code"))
			err = dps.MemberService.SaveRelation(memberId, "", invId,
				this.BaseC.GetPartnerId(ctx))
		}
	}

	if err != nil {
		result.Message = "注册失败," + err.Error() + "!"
	} else {
		result.Result = true
	}
	ctx.Response.JsonOutput(result)
}

// 跳转到会员中心
// url : /user/jump_m
func (this *UserC) JumpToMCenter(ctx *web.Context) {
	w := ctx.Response
	m := this.BaseC.GetMember(ctx)
	var location string
	if m == nil {
		location = "/user/login?return_url=/user/jump_m"
	} else {
		location = fmt.Sprintf("http://%s.%s/login/partner_connect?device=%s&sessionId=%s&mid=%d&token=%s",
			variable.DOMAIN_PREFIX_MEMBER,
			ctx.App.Config().GetString(variable.ServerDomain),
			util.GetBrownerDevice(ctx),
			ctx.Session().GetSessionId(),
			m.Id,
			m.DynamicToken,
		)
	}
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
			<script>window.onload=function(){location.replace('/')}</script></body></html>`,
		variable.DOMAIN_PREFIX_MEMBER,
		ctx.App.Config().GetString(variable.ServerDomain),
	)))
}
