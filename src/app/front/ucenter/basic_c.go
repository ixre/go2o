/**
 * Copyright 2015 @ z3q.net.
 * name : basic_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ucenter

import (
	"encoding/json"
	"errors"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/mvc"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
	"html/template"
	"strconv"
	"go2o/src/core/dto"
)

var _ mvc.Filter = new(baseC)

type basicC struct {
	*baseC
}

func (this *basicC) Profile(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	mm := this.GetMember(ctx)
	js, _ := json.Marshal(mm)
	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"member":       mm,
		"entity":       template.JS(js),
	}, "views/ucenter/{device}/profile.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}

func (this *basicC) Pwd(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	mm := this.GetMember(ctx)
	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"member":       mm,
	}, "views/ucenter/{device}/pwd.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}

func (this *basicC) Pwd_post(ctx *web.Context) {
	r := ctx.Request
	var result gof.Message
	r.ParseForm()
	m := this.GetMember(ctx)
	var oldPwd, newPwd, rePwd string
	oldPwd = r.FormValue("OldPwd")
	newPwd = r.FormValue("NewPwd")
	rePwd = r.FormValue("RePwd")
	var err error
	if newPwd != rePwd {
		err = errors.New("两次密码输入不一致")
	} else {
		err = dps.MemberService.ModifyPassword(m.Id, oldPwd, newPwd)
	}
	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
		this.ReCacheMember(ctx, m.Id)
	}
	ctx.Response.JsonOutput(result)
}

// 交易密码
func (this *basicC) Trade_pwd(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	mm := this.GetMember(ctx)

	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"member":       mm,
		"notFirstSet":  len(mm.TradePwd) != 0,
	}, "views/ucenter/{device}/trade_pwd.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}
func (this *basicC) Trade_pwd_post(ctx *web.Context) {
	r := ctx.Request
	var result gof.Message
	r.ParseForm()
	m := this.GetMember(ctx)
	var oldPwd, newPwd, rePwd string
	oldPwd = r.FormValue("OldPwd")
	newPwd = r.FormValue("NewPwd")
	rePwd = r.FormValue("RePwd")
	var err error
	if newPwd != rePwd {
		err = errors.New("两次密码输入不一致")
	} else {
		err = dps.MemberService.ModifyTradePassword(m.Id, oldPwd, newPwd)
	}
	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
		this.ReCacheMember(ctx, m.Id)
	}
	ctx.Response.JsonOutput(result)
}

func (this *basicC) Profile_post(ctx *web.Context) {
	mm := this.GetMember(ctx)
	r := ctx.Request
	var result gof.Message
	r.ParseForm()
	m := new(member.ValueMember)
	web.ParseFormToEntity(r.Form, m)
	m.Id = mm.Id
	_, err := dps.MemberService.SaveMember(m)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
		this.ReCacheMember(ctx, m.Id)
	}
	ctx.Response.JsonOutput(result)
}

func (this *basicC) Deliver(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"member":       m,
	}, "views/ucenter/{device}/deliver.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}

func (this *basicC) Deliver_post(ctx *web.Context) {
	m := this.GetMember(ctx)
	add := dps.MemberService.GetDeliverAddress(m.Id)
	js, _ := json.Marshal(add)
	ctx.Response.Write([]byte(`{"rows":` + string(js) + `}`))
}

func (this *basicC) SaveDeliver_post(ctx *web.Context) {
	m := this.GetMember(ctx)
	var result gof.Message
	r := ctx.Request
	r.ParseForm()
	var e member.DeliverAddress
	web.ParseFormToEntity(r.Form, &e)
	e.MemberId = m.Id
	_, err := dps.MemberService.SaveDeliverAddress(m.Id, &e)
	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	ctx.Response.JsonOutput(result)
}

func (this *basicC) DeleteDeliver_post(ctx *web.Context) {
	r := ctx.Request
	var result gof.Message
	m := this.GetMember(ctx)
	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("id"))

	err := dps.MemberService.DeleteDeliverAddress(m.Id, id)
	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	ctx.Response.JsonOutput(result)
}

func (this *basicC) Member_cln_filter(ctx *web.Context){
	r := ctx.Request
	key :=  r.URL.Query().Get("key")


	if len(key) < 3{
		ctx.Response.JsonOutput(gof.Message{
			Message:"length less more",
		})
		return
	}

	var list []*dto.SimpleMember
	partnerId := this.GetPartner(ctx).Id
	list = dps.MemberService.FilterMemberByUsrOrPhone(partnerId,key)
	ctx.Response.JsonOutput(list)
}