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
	"github.com/jsix/gof/util"
	gfmt "github.com/jsix/gof/util/fmt"
	"github.com/jsix/gof/web/form"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/valueobject"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type memberC struct {
}

func (this *memberC) LevelList(ctx *echox.Context) error {
	d := ctx.NewData()
	return ctx.RenderOK("member.level_list.html", d)
}

//修改门店信息
func (this *memberC) EditMLevel(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	id, _ := strconv.Atoi(ctx.Query("id"))
	entity := dps.PartnerService.GetMemberLevelById(partnerId, id)
	js, _ := json.Marshal(entity)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("member.edit_level.html", d)
}

func (this *memberC) CreateMLevel(ctx *echox.Context) error {
	d := ctx.NewData()
	d.Map["entity"] = template.JS("{}")
	return ctx.RenderOK("member.create_level.html", d)
}

// 保存会员等级(POST)
func (this *memberC) SaveMLevel(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		var result gof.Result
		r.ParseForm()

		e := valueobject.MemberLevel{}
		form.ParseEntity(r.Form, &e)
		e.PartnerId = getPartnerId(ctx)

		id, err := dps.PartnerService.SaveMemberLevel(partnerId, &e)

		if err != nil {
			result.ErrMsg = err.Error()
			result.ErrCode = 1
		} else {
			result.Data = map[string]string{"id":util.Str(id)}
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

func (this *memberC) DelMLevel(ctx *echox.Context) error {
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		var result gof.Result
		r.ParseForm()
		partnerId := getPartnerId(ctx)
		id, err := strconv.Atoi(r.FormValue("id"))
		if err == nil {
			err = dps.PartnerService.DelMemberLevel(partnerId, id)
		}

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 会员列表
func (this *memberC) List(ctx *echox.Context) error {
	levelDr := getLevelDropDownList(getPartnerId(ctx))
	d := ctx.NewData()
	d.Map["levelDr"] = template.HTML(levelDr)
	return ctx.RenderOK("member.list.html", d)
}

// 锁定会员
func (this *memberC) Lock_member(ctx *echox.Context) error {
	req := ctx.HttpRequest()
	if req.Method == "POST" {
		req.ParseForm()
		id, _ := strconv.Atoi(req.FormValue("id"))
		partnerId := getPartnerId(ctx)
		var result gof.Result
		if _, err := dps.MemberService.LockMember(partnerId, id); err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

func (this *memberC) Member_details(ctx *echox.Context) error {
	memberId, _ := strconv.Atoi(ctx.Query("member_id"))

	d := ctx.NewData()
	d.Map["memberId"] = memberId
	return ctx.RenderOK("member.details.html", d)
}

// 会员基本信息
func (this *memberC) Member_basic(ctx *echox.Context) error {
	memberId, _ := strconv.Atoi(ctx.Query("member_id"))
	m := dps.MemberService.GetMember(memberId)
	if m == nil {
		return ctx.String(http.StatusOK, "no such member")
	}
	lv := dps.PartnerService.GetLevel(getPartnerId(ctx), m.Level)
	rl := dps.MemberService.GetRelation(m.Id)
	var invName string = "无"
	if rl.RefereesId > 0 {
		rlm := dps.MemberService.GetMember(rl.RefereesId)
		invName = rlm.Name + "(" + rlm.Usr + ")"
	}
	d := ctx.NewData()
	d.Map = map[string]interface{}{
		"m":       m,
		"lv":      lv,
		"invId":   rl.RefereesId,
		"invName": invName,
		"sexName": gfmt.BoolString(m.Sex == 1, "先生",
			gfmt.BoolString(m.Sex == 2, "女士", "-")),
		"lastLoginTime": format.HanUnixDateTime(m.LastLoginTime),
		"regTime":       format.HanUnixDateTime(m.RegTime),
		"Avatar":        format.GetResUrl(m.Avatar),
	}

	return ctx.RenderOK("member.basic_info.html", d)
}

// 会员账户信息
func (this *memberC) Member_account(ctx *echox.Context) error {
	memberId, _ := strconv.Atoi(ctx.Query("member_id"))
	acc := dps.MemberService.GetAccount(memberId)
	if acc == nil {
		return ctx.String(http.StatusOK, "no such account")
	}

	d := ctx.NewData()
	d.Map = map[string]interface{}{
		"acc": acc,
		"balanceAccountAlias": variable.AliasBalanceAccount,
		"presentAccountAlias": variable.AliasPresentAccount,
		"flowAccountAlias":    variable.AliasFlowAccount,
		"growAccountAlias":    variable.AliasGrowAccount,
		"integralAlias":       variable.AliasIntegral,
		"updateTime":          format.HanUnixDateTime(acc.UpdateTime),
	}
	return ctx.Render(http.StatusOK, "member.account_info.html", d)

}

// 会员收款银行信息
func (this *memberC) Member_bankinfo(ctx *echox.Context) error {
	memberId, _ := strconv.Atoi(ctx.Query("member_id"))
	e := dps.MemberService.GetBank(memberId)
	if e != nil && len(e.Account) > 0 && len(e.AccountName) > 0 &&
		len(e.Name) > 0 && len(e.Network) > 0 {
		d := ctx.NewData()
		d.Map["bank"] = e
		return ctx.RenderOK("member.bank_info.html", d)
	}
	return ctx.String(http.StatusOK, "<span class=\"red\">尚未完善</span>")
}

func (this *memberC) Unlock_bankinfo(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		memberId, _ := strconv.Atoi(ctx.Query("member_id"))
		msg := new(gof.Result)
		err := dps.MemberService.UnlockBankInfo(memberId)
		return ctx.JSON(http.StatusOK, msg.Error(err))
	}
	return nil
}

// 重置密码(POST)
func (this *memberC) Reset_pwd(ctx *echox.Context) error {
	var result gof.Result
	req := ctx.HttpRequest()
	if req.Method == "POST" {
		req.ParseForm()
		memberId, _ := strconv.Atoi(req.FormValue("member_id"))
		rl := dps.MemberService.GetRelation(memberId)
		partnerId := getPartnerId(ctx)
		if rl == nil || rl.RegisterPartnerId != partnerId {
			result.ErrMsg = "无权进行当前操作"
		} else {
			newPwd := dps.MemberService.ResetPassword(memberId)
			result.ErrCode = 0
			result.ErrMsg = fmt.Sprintf("重置成功,新密码为: %s", newPwd)
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 客服充值
func (this *memberC) Charge(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.charge_post(ctx)
	}
	memberId, _ := strconv.Atoi(ctx.Query("member_id"))
	mem := dps.MemberService.GetMemberSummary(memberId)
	if mem == nil {
		return ctx.String(http.StatusOK, "no such member")
	}
	d := ctx.NewData()
	d.Map["m"] = mem
	return ctx.RenderOK("member.charge.html", d)
}
func (this *memberC) charge_post(ctx *echox.Context) error {
	var msg gof.Result
	var err error
	req := ctx.HttpRequest()
	req.ParseForm()
	partnerId := getPartnerId(ctx)
	memberId, _ := strconv.Atoi(req.FormValue("MemberId"))
	amount, _ := strconv.ParseFloat(req.FormValue("Amount"), 32)
	if amount < 0 {
		msg.ErrMsg = "error amount"
	} else {
		rel := dps.MemberService.GetRelation(memberId)

		if rel.RegisterPartnerId != getPartnerId(ctx) {
			err = partner.ErrPartnerNotMatch
		} else {
			title := fmt.Sprintf("[KF]客服充值%s", format.FormatFloat(float32(amount)))
			err = dps.MemberService.Charge(partnerId, memberId,
				member.TypeBalanceServiceCharge, title, "-", float32(amount))
		}
		if err != nil {
			msg.ErrMsg = err.Error()
		} else {
	msg.ErrCode = 0
		}
	}
	return ctx.JSON(http.StatusOK, msg)
}

// 提现列表
func (this *memberC) ApplyRequestList(ctx *echox.Context) error {
	levelDr := getLevelDropDownList(getPartnerId(ctx))

	d := ctx.NewData()
	d.Map["levelDr"] = template.HTML(levelDr)
	d.Map["kind"] = member.KindBalanceApplyCash
	return ctx.RenderOK("member.apply_request_list.html", d)
}

// 审核提现请求
func (this *memberC) Pass_apply_req(ctx *echox.Context) error {
	var msg gof.Result
	req := ctx.HttpRequest()
	if req.Method == "POST" {
		req.ParseForm()
		partnerId := getPartnerId(ctx)
		passed := req.FormValue("pass") == "1"
		memberId, _ := strconv.Atoi(req.FormValue("member_id"))
		id, _ := strconv.Atoi(req.FormValue("id"))

		err := dps.MemberService.ConfirmApplyCash(partnerId, memberId, id, passed, "")

		if err != nil {
			msg.ErrMsg = err.Error()
		} else {
	msg.ErrCode = 0
		}
		return ctx.JSON(http.StatusOK, msg)
	}
	return nil
}

// 退回提现请求
func (this *memberC) Back_apply_req(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.back_apply_req_post(ctx)
	}
	memberId, _ := strconv.Atoi(ctx.Query("member_id"))
	id, _ := strconv.Atoi(ctx.Query("id"))

	info := dps.MemberService.GetBalanceInfoById(memberId, id)

	if info == nil {
		return ctx.String(http.StatusOK, "no such request")
	}

	d := ctx.NewData()
	d.Map["info"] = info
	d.Map["applyTime"] = time.Unix(info.CreateTime, 0).Format("2006-01-02 15:04:05")
	return ctx.RenderOK("member.back_apply_req.html", d)
}

func (this *memberC) back_apply_req_post(ctx *echox.Context) error {
	var msg gof.Result
	req := ctx.HttpRequest()
	req.ParseForm()
	partnerId := getPartnerId(ctx)
	memberId, _ := strconv.Atoi(req.FormValue("MemberId"))
	id, _ := strconv.Atoi(req.FormValue("Id"))

	err := dps.MemberService.ConfirmApplyCash(partnerId, memberId, id, false, "")
	if err != nil {
		msg.ErrMsg = err.Error()
	} else {
msg.ErrCode = 0
	}
	return ctx.JSON(http.StatusOK, msg)
}

// 提现打款
func (this *memberC) Handle_apply_req(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.handle_apply_req_post(ctx)
	}
	memberId, _ := strconv.Atoi(ctx.Query("member_id"))
	id, _ := strconv.Atoi(ctx.Query("id"))

	info := dps.MemberService.GetBalanceInfoById(memberId, id)

	if info == nil {
		return ctx.String(http.StatusOK, "no such info")
	}

	d := ctx.NewData()
	bank := dps.MemberService.GetBank(memberId)
	if info.Amount < 0 {
		info.Amount = -info.Amount
	}
	d.Map = map[string]interface{}{
		"info":      info,
		"bank":      bank,
		"applyTime": time.Unix(info.CreateTime, 0).Format("2006-01-02 15:04:05"),
	}
	return ctx.RenderOK("member.handle_apply_req.html", d)
}

func (this *memberC) handle_apply_req_post(ctx *echox.Context) error {
	var msg gof.Result
	var err error
	req := ctx.HttpRequest()
	req.ParseForm()
	partnerId := getPartnerId(ctx)
	memberId, _ := strconv.Atoi(req.FormValue("MemberId"))
	id, _ := strconv.Atoi(req.FormValue("Id"))
	agree := req.FormValue("Agree") == "on"
	tradeNo := req.FormValue("TradeNo")

	if !agree {
		err = errors.New("请同意已知晓并打款选项")
	} else {
		err = dps.MemberService.FinishApplyCash(partnerId, memberId, id, tradeNo)
	}
	if err != nil {
		msg.ErrMsg = err.Error()
	} else {
msg.ErrCode = 0
	}
	return ctx.JSON(http.StatusOK, msg)
}

// 团队排名列表
func (this *memberC) Team_rank(ctx *echox.Context) error {
	levelDr := getLevelDropDownList(getPartnerId(ctx))
	d := ctx.NewData()
	d.Map["levelDr"] = template.HTML(levelDr)
	return ctx.RenderOK("member.team_rank.html", d)
}

// 邀请关系
func (this *memberC) Invi_relation(c *echox.Context) error {
	partnerId := getPartnerId(c)
	memberId, _ := strconv.Atoi(c.Query("member_id"))
	ms := dps.MemberService
	idArr := []int{memberId}
	rl := ms.GetRelation(memberId)
	for rl != nil && rl.RefereesId > 0 {
		idArr = append(idArr, rl.RefereesId)
		rl = ms.GetRelation(rl.RefereesId)
	}
	//for i, j := 0, len(idArr)-1; i < j; i, j = i+1, j-1 {
	//	//反序
	//	idArr[i], idArr[j] = idArr[j], idArr[i]
	//}
	list, _ := json.Marshal(ms.GetMemberList(partnerId, idArr))
	d := c.NewData()
	d.Map["listJson"] = template.JS(list)
	return c.RenderOK("member.invi_relation.html", d)
}
