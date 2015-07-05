/**
 * Copyright 2014 @ S1N1 Team.
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
	"go2o/src/core/service/dps"
	"go2o/src/core/service/goclient"
	"html/template"
	"strconv"
)

type accountC struct {
	*baseC
}

func (this *accountC) Income_log(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"conf":    conf,
		"partner": p,
		"member":  m,
	}, "views/ucenter/{device}/account/income_log.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}

func (this *accountC) Income_log_post(ctx *web.Context) {
	m := this.GetMember(ctx)
	r := ctx.Request
	r.ParseForm()
	page, _ := strconv.Atoi(r.FormValue("page"))
	size, _ := strconv.Atoi(r.FormValue("size"))

	n, rows := dps.MemberService.QueryIncomeLog(m.Id, page, size, "", "record_time DESC")
	p := pager.NewUrlPager(pager.TotalPage(n, size), page, pager.GetterJavaScriptPager)
	pager := &front.Pager{Total: n, Rows: rows, Text: p.PagerString()}
	ctx.Response.JsonOutput(pager)
}

func (this *accountC) Apply_cash(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
	acc := dps.MemberService.GetAccount(m.Id)
	bank := dps.MemberService.GetBank(m.Id)

	js, _ := json.Marshal(bank)
	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"conf":    conf,
		"record":  15,
		"partner": p,
		"member":  m,
		"account": acc,
		"entity":  template.JS(js),
	}, "views/ucenter/{device}/account/apply_cash.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}

func (this *accountC) Apply_cash_post(ctx *web.Context) {
	m := this.GetMember(ctx)
	r, w := ctx.Request, ctx.Response
	var result gof.Message
	r.ParseForm()
	e := new(member.BankInfo)
	web.ParseFormToEntity(r.Form, e)
	e.MemberId = m.Id
	err := goclient.Member.SaveBankInfo(m.Id, m.DynamicToken, e)

	if err != nil {
		result = gof.Message{Result: false, Message: err.Error()}
	} else {
		result = gof.Message{Result: true}
	}
	w.Write(result.Marshal())

}

func (this *accountC) Integral_exchange(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
	acc, _ := goclient.Member.GetMemberAccount(m.Id, m.DynamicToken)

	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"conf":    conf,
		"record":  15,
		"partner": p,
		"member":  m,
		"account": acc,
	}, "views/ucenter/{device}/account/integral_exchange.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}
