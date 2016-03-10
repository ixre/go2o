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
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/dto"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
	"go2o/src/core/variable"
	"bytes"
	"fmt"
	fmt2 "github.com/jsix/gof/util/fmt"
)

type basicC struct {
}

func getFieldHtml(alias,note,field string,show bool)string{
	const tpl string = `<div class=\"fl %s\">
			<div class="label">%s：</div>
		<div class="in">
		<input class="ui-box ui-validate" tipin="tip_%s" field="%s"/>
		<span class="fv" id="tip_%s">%s</span>
		</div>
		</div>`
	return fmt.Sprintf(tpl,fmt2.BoolString(show,""," hidden"),
		alias,field,field,field,note)
}
func getExtFormFields()*bytes.Buffer{
	buf := bytes.NewBufferString(nil)
	buf.WriteString(getFieldHtml(
		variable.AliasMemberExt1,
		variable.AliasMemberExt1Note,
		variable.AliasMemberExt1Show,"Ext1"))
	buf.WriteString(getFieldHtml(
		variable.AliasMemberExt2,
		variable.AliasMemberExt2Note,
		variable.AliasMemberExt2Show,"Ext2"))
	buf.WriteString(getFieldHtml(
		variable.AliasMemberExt3,
		variable.AliasMemberExt3Note,
		variable.AliasMemberExt3Show,"Ext3"))
	buf.WriteString(getFieldHtml(
		variable.AliasMemberExt4,
		variable.AliasMemberExt4Note,
		variable.AliasMemberExt4Show,"Ext4"))
	buf.WriteString(getFieldHtml(
		variable.AliasMemberExt5,
		variable.AliasMemberExt5Note,
		variable.AliasMemberExt5Show,"Ext5"))
	buf.WriteString(getFieldHtml(
		variable.AliasMemberExt6,
		variable.AliasMemberExt6Note,
		variable.AliasMemberExt6Show,"Ext6"))
}

func (this *basicC) Profile(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.profile_post(ctx)
	}
	mm := getMember(ctx)
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	js, _ := json.Marshal(mm)
	d := ctx.NewData()


	d.Map = gof.TemplateDataMap{
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"member":       mm,
		"entity":       template.JS(js),
		"aliasIM":      variable.AliasMemberIM,
		"extFields":template.HTML(getExtFormFields().String()),
	}
	return ctx.RenderOK("profile.html", d)
}

func (this *basicC) profile_post(ctx *echox.Context) error {
	mm := getMember(ctx)
	r := ctx.Request()
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
		reCacheMember(ctx, m.Id)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (this *basicC) Pwd(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.pwd_post(ctx)
	}
	mm := getMember(ctx)
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"member":       mm,
	}
	return ctx.RenderOK("pwd.html", d)
}

func (this *basicC) pwd_post(ctx *echox.Context) error {
	r := ctx.Request()
	var result gof.Message
	r.ParseForm()
	m := getMember(ctx)
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
		reCacheMember(ctx, m.Id)
	}
	return ctx.JSON(http.StatusOK, result)
}

// 交易密码
func (this *basicC) Trade_pwd(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.trade_pwd_post(ctx)
	}
	m := getMember(ctx)
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"member":       m,
		"notFirstSet":  len(m.TradePwd) != 0,
	}
	return ctx.RenderOK("trade_pwd.html", d)
}
func (this *basicC) trade_pwd_post(ctx *echox.Context) error {
	r := ctx.Request()
	var result gof.Message
	r.ParseForm()
	m := getMember(ctx)
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
		reCacheMember(ctx, m.Id)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (this *basicC) Deliver(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.deliver_post(ctx)
	}
	m := getMember(ctx)
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"member":       m,
	}
	return ctx.RenderOK("delivery.html", d)
}

func (this *basicC) deliver_post(ctx *echox.Context) error {
	m := getMember(ctx)
	add := dps.MemberService.GetDeliverAddress(m.Id)
	js, _ := json.Marshal(add)
	ctx.Response().Write([]byte(`{"rows":` + string(js) + `}`))
	return nil
}

func (this *basicC) SaveDeliver(ctx *echox.Context) error {
	r := ctx.Request()
	if r.Method == "POST" {
		m := getMember(ctx)
		var result gof.Message
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
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

func (this *basicC) DeleteDeliver(ctx *echox.Context) error {
	r := ctx.Request()
	if r.Method == "POST" {
		var result gof.Message
		m := getMember(ctx)
		r.ParseForm()
		id, _ := strconv.Atoi(r.FormValue("id"))

		err := dps.MemberService.DeleteDeliverAddress(m.Id, id)
		if err != nil {
			result.Message = err.Error()
		} else {
			result.Result = true
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

func (this *basicC) Member_cln_filter(ctx *echox.Context) error {
	key := ctx.Query("key")
	if len(key) < 3 {
		return ctx.JSON(http.StatusOK, gof.Message{
			Message: "length less more",
		})
	}
	var list []*dto.SimpleMember
	partnerId := getPartner(ctx).Id
	list = dps.MemberService.FilterMemberByUsrOrPhone(partnerId, key)
	return ctx.JSON(http.StatusOK, list)
}
