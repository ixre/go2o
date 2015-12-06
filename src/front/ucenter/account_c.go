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
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/core/service/goclient"
	"go2o/src/core/variable"
	"go2o/src/front"
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

// 提现申请
func (this *accountC) Apply_cash(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
	acc := dps.MemberService.GetAccount(m.Id)
	saleConf := dps.PartnerService.GetSaleConf(p.Id)

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

	var commissionStr string
	if saleConf.ApplyCsn == 0 {
		commissionStr = "不收取手续费"
	} else {
		commissionStr = fmt.Sprintf("收取<i>%s%s</i>手续费",
			format.FormatFloat(saleConf.ApplyCsn*100), "%")
	}

	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"conf":           conf,
		"partner":        p,
		"member":         m,
		"minAmount":      format.FormatFloat(float32(minAmount)),
		"maxApplyAmount": maxApplyAmount,
		"account":        acc,
		"latestInfo":     template.HTML(latestInfo),
		"commissionStr":  template.HTML(commissionStr),
		"presentAlias":   variable.AliasFlowAccount,
		"cns":            saleConf.ApplyCsn,
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
	memberId := this.GetMember(ctx).Id
	saleConf := dps.PartnerService.GetSaleConf(partnerId)
	bank := dps.MemberService.GetBank(memberId)

	if bank == nil || len(bank.Account) == 0 || len(bank.AccountName) == 0 ||
		len(bank.Network) == 0 {
		err = errors.New("请先设置收款银行信息")
		goto toErr
	}

	if _, err = dps.MemberService.VerifyTradePwd(memberId, tradePwd); err != nil {
		goto toErr
	}

	if amount < minAmount {
		err = errors.New(fmt.Sprintf("必须达到最低提现金额:%s元",
			format.FormatFloat(float32(minAmount))))
	} else {
		m := this.GetMember(ctx)
		err = dps.MemberService.SubmitApplyPresentBalance(partnerId, m.Id,
			member.TypeApplyCashToBank, float32(amount), saleConf.ApplyCsn)
	}

toErr:
	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
	ctx.Response.JsonOutput(msg)
}

// 转换活动金到提现账户
func (this *accountC) Convert_f2p(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	saleConf := dps.PartnerService.GetSaleConf(p.Id)
	m := this.GetMember(ctx)
	acc := dps.MemberService.GetAccount(m.Id)

	var commissionStr string
	if saleConf.FlowConvertCsn == 0 {
		commissionStr = "不收取手续费"
	} else {
		commissionStr = fmt.Sprintf("收取<i>%s%s</i>手续费",
			format.FormatFloat(saleConf.FlowConvertCsn*100), "%")
	}

	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"conf":              conf,
		"partner":           p,
		"member":            m,
		"account":           acc,
		"commissionStr":     template.HTML(commissionStr),
		"flowAlias":         variable.AliasFlowAccount,
		"flowConvertSlogan": variable.FlowConvertSlogan,
		"cns":               saleConf.FlowConvertCsn,
		"notSetTradePwd":    len(m.TradePwd) == 0,
	}, "views/ucenter/{device}/account/convert_f2p.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}

func (this *accountC) Convert_f2p_post(ctx *web.Context) {
	var msg gof.Message
	var err error
	ctx.Request.ParseForm()
	partnerId := this.GetPartner(ctx).Id
	amount, _ := strconv.ParseFloat(ctx.Request.FormValue("Amount"), 32)
	tradePwd := ctx.Request.FormValue("TradePwd")
	saleConf := dps.PartnerService.GetSaleConf(partnerId)

	memberId := this.GetMember(ctx).Id

	if _, err = dps.MemberService.VerifyTradePwd(memberId, tradePwd); err == nil {
		err = dps.MemberService.TransferFlow(memberId, member.KindBalancePresent,
			float32(amount), saleConf.FlowConvertCsn, domain.NewTradeNo(partnerId),
			fmt.Sprintf("%s转换", variable.AliasFlowAccount),
			fmt.Sprintf("%s转换%s", variable.AliasFlowAccount, variable.AliasPresentAccount))
	}

	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
	ctx.Response.JsonOutput(msg)
}

// 转换活动金到提现账户
func (this *accountC) Transfer_f2m(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	saleConf := dps.PartnerService.GetSaleConf(p.Id)
	m := this.GetMember(ctx)
	acc := dps.MemberService.GetAccount(m.Id)

	var commissionStr string
	if saleConf.FlowConvertCsn == 0 {
		commissionStr = "不收取手续费"
	} else {
		commissionStr = fmt.Sprintf("收取<i>%s%s</i>手续费",
			format.FormatFloat(saleConf.FlowConvertCsn*100), "%")
	}

	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"conf":           conf,
		"partner":        p,
		"member":         m,
		"account":        acc,
		"commissionStr":  template.HTML(commissionStr),
		"flowAlias":      variable.AliasFlowAccount,
		"cns":            saleConf.TransCsn,
		"notSetTradePwd": len(m.TradePwd) == 0,
	}, "views/ucenter/{device}/account/transfer_f2m.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}

func (this *accountC) Transfer_f2m_post(ctx *web.Context) {
	var msg gof.Message
	var err error
	ctx.Request.ParseForm()
	form := ctx.Request.Form
	partnerId := this.GetPartner(ctx).Id
	toMemberId, _ := strconv.Atoi(form.Get("ToId"))
	amount, _ := strconv.ParseFloat(form.Get("Amount"), 32)
	tradePwd := form.Get("TradePwd")
	saleConf := dps.PartnerService.GetSaleConf(partnerId)
	memberId := this.GetMember(ctx).Id

	if toMemberId == memberId {
		err = errors.New("无法转账到自己账号")
	} else {
		if _, err = dps.MemberService.VerifyTradePwd(memberId, tradePwd); err == nil {
			err = dps.MemberService.TransferFlowTo(memberId, toMemberId, member.KindBalanceFlow,
				float32(amount), saleConf.TransCsn, domain.NewTradeNo(partnerId),
				variable.AliasFlowAccount+"转账", "转入"+variable.AliasFlowAccount)
		}
	}

	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
	ctx.Response.JsonOutput(msg)
}

// 转账成功提示页面
func (this *accountC) Transfer_success(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)

	src := ctx.Request.URL.Query().Get("src")
	var title, subTitle, btnText string

	switch src {
	case "trans_f2m":
		title = "转账成功"
		subTitle = "转账成功！"
		btnText = "继续转账"
	case "convert_f2p":
		title = "转换成功"
		subTitle = variable.AliasFlowAccount + "转换成功！"
		btnText = "继续转换"
	case "apply_p2b":
		title = "申请成功"
		subTitle = "申请成功，客服将在1-3个工作日完成审核。"
		btnText = "继续提现"
	}

	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"conf":     conf,
		"partner":  p,
		"title":    title,
		"subTitle": subTitle,
		"btnText":  btnText,
		"referer":  ctx.Request.Referer(),
	}, "views/ucenter/{device}/account/transfer_success.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
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
