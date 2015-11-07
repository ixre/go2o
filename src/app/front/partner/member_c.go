/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:55
 * description :
 * history :
 */

package partner

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jsix/gof"
	gfmt "github.com/jsix/gof/util/fmt"
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/mvc"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/valueobject"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"html/template"
	"strconv"
	"time"
)

var _ mvc.Filter = new(memberC)

type memberC struct {
	*baseC
}

func (this *memberC) LevelList(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.Response, gof.TemplateDataMap{},
		"views/partner/member/level_list.html")
}

//修改门店信息
func (this *memberC) EditMLevel(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	entity := dps.PartnerService.GetMemberLevelById(partnerId, id)
	bys, _ := json.Marshal(entity)

	ctx.App.Template().Execute(w,
		gof.TemplateDataMap{
			"entity": template.JS(bys),
		},
		"views/partner/member/edit_level.html")
}

func (this *memberC) CreateMLevel(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity": "{}",
		},
		"views/partner/member/create_level.html")
}

func (this *memberC) SaveMLevel_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r := ctx.Request
	var result gof.Message
	r.ParseForm()

	e := valueobject.MemberLevel{}
	web.ParseFormToEntity(r.Form, &e)
	e.PartnerId = this.GetPartnerId(ctx)

	id, err := dps.PartnerService.SaveMemberLevel(partnerId, &e)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
		result.Data = id
	}
	ctx.Response.JsonOutput(result)
}

func (this *memberC) DelMLevel(ctx *web.Context) {
	r := ctx.Request
	var result gof.Message
	r.ParseForm()
	partnerId := this.GetPartnerId(ctx)
	id, err := strconv.Atoi(r.FormValue("id"))
	if err == nil {
		err = dps.PartnerService.DelMemberLevel(partnerId, id)
	}

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	ctx.Response.JsonOutput(result)
}

// 会员列表
func (this *memberC) List(ctx *web.Context) {
	levelDr := getLevelDropDownList(this.GetPartnerId(ctx))
	ctx.App.Template().Execute(ctx.Response, gof.TemplateDataMap{
		"levelDr": template.HTML(levelDr),
	}, "views/partner/member/member_list.html")
}

// 锁定会员
func (this *memberC) Lock_member_post(ctx *web.Context) {
	ctx.Request.ParseForm()
	id, _ := strconv.Atoi(ctx.Request.FormValue("id"))
	partnerId := this.GetPartnerId(ctx)
	var result gof.Message
	if _, err := dps.MemberService.LockMember(partnerId, id); err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	ctx.Response.JsonOutput(result)
}

func (this *memberC) Member_details(ctx *web.Context) {
	memberId, _ := strconv.Atoi(ctx.Request.URL.Query().Get("member_id"))

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"memberId": memberId,
		}, "views/partner/member/member_details.html")
}

// 会员基本信息
func (this *memberC) Member_basic(ctx *web.Context) {
	memberId, _ := strconv.Atoi(ctx.Request.URL.Query().Get("member_id"))
	m := dps.MemberService.GetMember(memberId)
	if m == nil {
		ctx.Response.Write([]byte("no such member"))
	} else {

		lv := dps.PartnerService.GetLevel(this.GetPartnerId(ctx), m.Level)

		ctx.App.Template().Execute(ctx.Response,
			gof.TemplateDataMap{
				"m":  m,
				"lv": lv,
				"sexName": gfmt.BoolString(m.Sex == 1, "先生",
					gfmt.BoolString(m.Sex == 2, "女士", "-")),
				"lastLoginTime": format.HanUnixDateTime(m.LastLoginTime),
				"regTime":       format.HanUnixDateTime(m.RegTime),
			}, "views/partner/member/basic_info.html")
	}
}

// 会员账户信息
func (this *memberC) Member_account(ctx *web.Context) {
	memberId, _ := strconv.Atoi(ctx.Request.URL.Query().Get("member_id"))
	acc := dps.MemberService.GetAccount(memberId)
	if acc != nil {

		ctx.App.Template().Execute(ctx.Response,
			gof.TemplateDataMap{
				"acc": acc,
				"balanceAccountAlias": variable.AliasBalanceAccount,
				"presentAccountAlias": variable.AliasPresentAccount,
				"flowAccountAlias":    variable.AliasFlowAccount,
				"growAccountAlias":    variable.AliasGrowAccount,
				"integralAlias":       variable.AliasIntegral,
				"updateTime":          format.HanUnixDateTime(acc.UpdateTime),
			}, "views/partner/member/account_info.html")
	}
}

func (this *memberC) Reset_pwd_post(ctx *web.Context) {
	var result gof.Message
	ctx.Request.ParseForm()
	memberId, _ := strconv.Atoi(ctx.Request.FormValue("member_id"))
	rl := dps.MemberService.GetRelation(memberId)
	partnerId := this.GetPartnerId(ctx)
	if rl == nil || rl.RegisterPartnerId != partnerId {
		result.Message = "无权进行当前操作"
	} else {
		newPwd := dps.MemberService.ResetPassword(memberId)
		result.Result = true
		result.Message = fmt.Sprintf("重置成功,新密码为: %s", newPwd)
	}
	ctx.Response.JsonOutput(result)
}

// 客服充值
func (this *memberC) Charge(ctx *web.Context) {
	memberId, _ := strconv.Atoi(ctx.Request.URL.Query().Get("member_id"))
	mem := dps.MemberService.GetMemberSummary(memberId)
	if mem == nil {
		ctx.Response.Write([]byte("no such member"))
	} else {
		ctx.App.Template().Execute(ctx.Response,
			gof.TemplateDataMap{
				"m": mem,
			}, "views/partner/member/charge.html")
	}
}

func (this *memberC) Charge_post(ctx *web.Context) {
	var msg gof.Message
	var err error
	ctx.Request.ParseForm()
	partnerId := this.GetPartnerId(ctx)
	memberId, _ := strconv.Atoi(ctx.Request.FormValue("MemberId"))
	amount, _ := strconv.ParseFloat(ctx.Request.FormValue("Amount"), 32)
	if amount < 0 {
		msg.Message = "error amount"
	} else {
		rel := dps.MemberService.GetRelation(memberId)

		if rel.RegisterPartnerId != this.GetPartnerId(ctx) {
			err = partner.ErrPartnerNotMatch
		} else {
			title := fmt.Sprintf("客服充值%f", amount)
			err = dps.MemberService.Charge(partnerId, memberId, member.TypeBalanceServiceCharge, title, "", float32(amount))
		}
		if err != nil {
			msg.Message = err.Error()
		} else {
			msg.Result = true
		}
	}
	ctx.Response.JsonOutput(msg)
}

// 提现列表
func (this *memberC) ApplyRequestList(ctx *web.Context) {
	levelDr := getLevelDropDownList(this.GetPartnerId(ctx))
	ctx.App.Template().Execute(ctx.Response, gof.TemplateDataMap{
		"levelDr": template.HTML(levelDr),
		"kind":    member.KindBalanceApplyCash,
	}, "views/partner/member/apply_request_list.html")
}

// 审核提现请求
func (this *memberC) Pass_apply_req_post(ctx *web.Context) {
	var msg gof.Message
	ctx.Request.ParseForm()
	partnerId := this.GetPartnerId(ctx)
	passed := ctx.Request.FormValue("pass") == "1"
	memberId, _ := strconv.Atoi(ctx.Request.FormValue("member_id"))
	id, _ := strconv.Atoi(ctx.Request.FormValue("id"))

	err := dps.MemberService.ConfirmApplyCash(partnerId, memberId, id, passed, "")

	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
	ctx.Response.JsonOutput(msg)
}

// 退回提现请求
func (this *memberC) Back_apply_req(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	memberId, _ := strconv.Atoi(form.Get("member_id"))
	id, _ := strconv.Atoi(form.Get("id"))

	info := dps.MemberService.GetBalanceInfoById(memberId, id)

	if info != nil {
		ctx.App.Template().Execute(ctx.Response, gof.TemplateDataMap{
			"info":      info,
			"applyTime": time.Unix(info.CreateTime, 0).Format("2006-01-02 15:04:05"),
		}, "views/partner/member/back_apply_req.html")
	}
}

func (this *memberC) Back_apply_req_post(ctx *web.Context) {
	var msg gof.Message
	ctx.Request.ParseForm()
	form := ctx.Request.Form
	partnerId := this.GetPartnerId(ctx)
	memberId, _ := strconv.Atoi(form.Get("MemberId"))
	id, _ := strconv.Atoi(form.Get("Id"))

	err := dps.MemberService.ConfirmApplyCash(partnerId, memberId, id, false, "")
	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
	ctx.Response.JsonOutput(msg)
}

// 提现打款
func (this *memberC) Handle_apply_req(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	memberId, _ := strconv.Atoi(form.Get("member_id"))
	id, _ := strconv.Atoi(form.Get("id"))

	info := dps.MemberService.GetBalanceInfoById(memberId, id)

	if info != nil {
		bank := dps.MemberService.GetBank(memberId)
		ctx.App.Template().Execute(ctx.Response, gof.TemplateDataMap{
			"info":      info,
			"bank":      bank,
			"applyTime": time.Unix(info.CreateTime, 0).Format("2006-01-02 15:04:05"),
		}, "views/partner/member/handle_apply_req.html")
	}
}

func (this *memberC) Handle_apply_req_post(ctx *web.Context) {
	var msg gof.Message
	var err error
	ctx.Request.ParseForm()
	form := ctx.Request.Form
	partnerId := this.GetPartnerId(ctx)
	memberId, _ := strconv.Atoi(form.Get("MemberId"))
	id, _ := strconv.Atoi(form.Get("Id"))
	agree := form.Get("Agree") == "on"
	tradeNo := form.Get("TradeNo")

	if !agree {
		err = errors.New("请同意已知晓并打款选项")
	} else {
		err = dps.MemberService.FinishApplyCash(partnerId, memberId, id, tradeNo)
	}
	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
	ctx.Response.JsonOutput(msg)
}
