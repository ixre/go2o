/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:55
 * description :
 * history :
 */

package partner

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/valueobject"
	"go2o/src/core/service/dps"
	"html/template"
	"strconv"
	"go2o/src/core/domain/interface/member"
	"fmt"
)

var _ mvc.Filter = new(memberC)

type memberC struct {
	*baseC
}

func (this *memberC) LevelList(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.Response, nil, "views/partner/member/level_list.html")
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
	//partnerId := this.GetPartnerId(ctx)
	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{}, "views/partner/member/member_list.html")
}

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

func (this *memberC) Charge(ctx *web.Context) {
	memberId, _ := strconv.Atoi(ctx.Request.URL.Query().Get("member_id"))
	mem := dps.MemberService.GetMemberSummary(memberId)
	if mem == nil{
		ctx.Response.Write([]byte("no such member"))
	}else {
		ctx.App.Template().Execute(ctx.Response,
			gof.TemplateDataMap{
				"m":mem,
			}, "views/partner/member/charge.html")
	}
}

func (this *memberC) Charge_post(ctx *web.Context) {
	var msg gof.Message
	ctx.Request.ParseForm()
	memberId, _ := strconv.Atoi(ctx.Request.FormValue("MemberId"))
	amount, _ := strconv.ParseFloat(ctx.Request.FormValue("Amount"), 32)
	if amount < 0 {
		msg.Message="error amount"
	}else {
		rel := dps.MemberService.GetRelation(memberId)

		if rel.RegisterPartnerId != this.GetPartnerId(ctx){
			msg.Message = "can not operate"
		}

		title := fmt.Sprintf("客服充值",amount)
		err := dps.MemberService.Charge(memberId,member.TypeBalanceServiceCharge,title,"", float32(amount))

		if err != nil {
			msg.Message = err.Error()
		} else {
			msg.Result = true
		}
	}
	ctx.Response.JsonOutput(msg)
}

