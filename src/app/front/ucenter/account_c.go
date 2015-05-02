/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/pager"
	"go2o/src/app/front"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"go2o/src/core/service/goclient"
	"html/template"
	"strconv"
)

type accountC struct {
	gof.App
}

func (this *accountC) IncomeLog(ctx *web.Context,
	m *member.ValueMember, p *partner.ValuePartner, conf *partner.SiteConf) {

	ctx.App.Template().Execute(ctx.ResponseWriter, func(mp *map[string]interface{}) {
		v := *mp
		v["conf"] = conf
		v["record"] = 15
		v["partner"] = p
		v["member"] = m
	}, "views/ucenter/account/income_log.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *accountC) IncomeLog_post(ctx *web.Context,
	m *member.ValueMember) {
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	page, _ := strconv.Atoi(r.FormValue("page"))
	size, _ := strconv.Atoi(r.FormValue("size"))

	n, rows := dps.MemberService.QueryIncomeLog(m.Id, page, size, "", "record_time DESC")

	p := pager.NewUrlPager(pager.TotalPage(n, size), page, pager.JavaScriptPagerGetter)

	pager := &front.Pager{Total: n, Rows: rows, Text: p.PagerString()}

	js, _ := json.Marshal(pager)
	w.Write(js)
}

func (this *accountC) ApplyCash(ctx *web.Context,
	m *member.ValueMember, p *partner.ValuePartner, conf *partner.SiteConf) {

	_, w := ctx.Request, ctx.ResponseWriter
	acc, err := goclient.Member.GetMemberAccount(m.Id, m.LoginToken)
	bank, err := goclient.Member.GetBankInfo(m.Id, m.LoginToken)

	if err != nil {
		w.Write([]byte("error:" + err.Error()))
		return
	}

	js, _ := json.Marshal(bank)
	ctx.App.Template().Execute(w, func(m *map[string]interface{}) {
		v := *m
		v["conf"] = conf
		v["record"] = 15
		v["partner"] = p
		v["member"] = m
		v["account"] = acc
		v["entity"] = template.JS(js)
	}, "views/ucenter/account/apply_cash.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *accountC) ApplyCash_post(ctx *web.Context,
	m *member.ValueMember, p *partner.ValuePartner, conf *partner.SiteConf) {
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.JsonResult
	r.ParseForm()
	e := new(member.BankInfo)
	web.ParseFormToEntity(r.Form, e)
	e.MemberId = m.Id
	err := goclient.Member.SaveBankInfo(m.Id, m.LoginToken, e)

	if err != nil {
		result = gof.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true}
	}
	w.Write(result.Marshal())

}

func (this *accountC) IntegralExchange(ctx *web.Context,
	m *member.ValueMember, p *partner.ValuePartner, conf *partner.SiteConf) {

	acc, _ := goclient.Member.GetMemberAccount(m.Id, m.LoginToken)

	ctx.App.Template().Execute(ctx.ResponseWriter, func(m *map[string]interface{}) {
		v := *m
		v["conf"] = conf
		v["record"] = 15
		v["partner"] = p
		v["member"] = m
		v["account"] = acc
	}, "views/ucenter/account/integral_exchange.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}
