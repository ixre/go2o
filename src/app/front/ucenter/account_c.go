/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/pager"
	"go2o/src/app/front"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/core/service/goclient"
	"html/template"
	"strconv"
	"strings"
)

const minAmount float64 = 50

var (
	bonusKindStr string
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

	if len(bonusKindStr) == 0 {
		bonusKindStr = strings.Join([]string{
			strconv.Itoa(member.KindBalanceTransfer),
			strconv.Itoa(member.KindBalancePresent),
			strconv.Itoa(member.KindBalanceFlow),
		}, ",")
	}

	n, rows := dps.MemberService.QueryIncomeLog(m.Id, page, size,
		" AND kind IN ("+bonusKindStr+")", "create_time DESC")
	p := pager.NewUrlPager(pager.TotalPage(n, size), page, pager.GetterJavaScriptPager)

	p.RecordCount = n
	pager := &front.Pager{Total: n, Rows: rows, Text: p.PagerString()}
	ctx.Response.JsonOutput(pager)
}

func (this *accountC) Apply_cash(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
	acc := dps.MemberService.GetAccount(m.Id)

	var latestInfo string = dps.MemberService.GetLatestApplyCashText(m.Id)
	if len(latestInfo) != 0 {
		latestInfo = "<div class=\"info\">" + latestInfo + "</div>"
	}

	var maxApplyAmount int
	if acc.PresentBalance < float32(minAmount) {
		maxApplyAmount = 0
	} else {
		maxApplyAmount = int(acc.PresentBalance)
	}

	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"conf":           conf,
		"partner":        p,
		"member":         m,
		"minAmount":      format.FormatFloat(float32(minAmount)),
		"maxApplyAmount": maxApplyAmount,
		"account":        acc,
		"latestInfo":     template.HTML(latestInfo),
		"notSetTradePwd": len(m.TradePwd) == 0,
	}, "views/ucenter/{device}/account/apply_cash.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}

func (this *accountC) Apply_cash_post(ctx *web.Context) {
	var msg gof.Message
	var err error
	ctx.Request.ParseForm()
	partnerId := this.GetPartner(ctx).Id
	amount, _ := strconv.ParseFloat(ctx.Request.FormValue("Amount"), 32)
	tradePwd := ctx.Request.FormValue("TradePwd")
	if amount < minAmount {
		err = errors.New(fmt.Sprintf("必须达到最低提现金额:%s元",
			format.FormatFloat(float32(minAmount))))
	} else {
		m := this.GetMember(ctx)
		err = dps.MemberService.SubmitApplyCash(partnerId, m.Id,
			tradePwd, member.TypeApplyCashToBank, float32(amount))
	}

	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
	ctx.Response.JsonOutput(msg)
}

func (this *accountC) Bank_info(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
	bank := dps.MemberService.GetBank(m.Id)

	js, _ := json.Marshal(bank)
	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"conf":    conf,
		"partner": p,
		"entity":  template.JS(js),
	}, "views/ucenter/{device}/account/bank_info.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}

func (this *accountC) Bank_info_post(ctx *web.Context) {
	m := this.GetMember(ctx)
	r := ctx.Request
	var msg gof.Message
	r.ParseForm()
	e := new(member.BankInfo)
	web.ParseFormToEntity(r.Form, e)
	e.MemberId = m.Id
	err := dps.MemberService.SaveBankInfo(e)

	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
	ctx.Response.JsonOutput(msg)
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
