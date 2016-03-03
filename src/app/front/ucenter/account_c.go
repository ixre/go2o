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
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
<<<<<<< HEAD
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
=======
	"go2o/src/core/service/goclient"
	"go2o/src/core/variable"
	"html/template"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"strconv"
	"strings"
)

<<<<<<< HEAD
const minAmount float64 = 0.01
=======
const minAmount float64 = 50
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d

var (
	bonusKindStr string
)

type accountC struct {
<<<<<<< HEAD
}

func (this *accountC) Income_log(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.income_log_post(ctx)
	}
	m := getMember(ctx)
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"conf":    conf,
		"partner": p,
		"member":  m,
	}
	return ctx.RenderOK("account.income_log.html", d)
}

func (this *accountC) income_log_post(ctx *echox.Context) error {
	m := getMember(ctx)
	r := ctx.Request()
=======
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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
<<<<<<< HEAD
	return ctx.JSON(http.StatusOK, pager)
}

// 提现申请
func (this *accountC) Apply_cash(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.apply_cash_post(ctx)
	}
	m := getMember(ctx)
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
=======
	ctx.Response.JsonOutput(pager)
}

// 提现申请
func (this *accountC) Apply_cash(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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

<<<<<<< HEAD
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
=======
	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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
<<<<<<< HEAD
	}
	return ctx.RenderOK("account.apply_cash.html", d)
}

func (this *accountC) apply_cash_post(ctx *echox.Context) error {
	var msg gof.Message
	var err error
	r := ctx.Request()
	r.ParseForm()
	partnerId := getPartner(ctx).Id
	strAmount := strings.TrimSpace(r.FormValue("Amount"))
	if len(strAmount) == 0 || strAmount == "NaN" {
		msg.Message = "提现金额错误"
		return ctx.JSON(http.StatusOK, msg)
	}

	amount, err := strconv.ParseFloat(r.FormValue("Amount"), 32)
	if err != nil {
		msg.Message = err.Error()
		return ctx.JSON(http.StatusOK, msg)
	}

	tradePwd := r.FormValue("TradePwd")
	memberId := getMember(ctx).Id
=======
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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
<<<<<<< HEAD
		m := getMember(ctx)
		_, _, err = dps.MemberService.SubmitApplyPresentBalance(partnerId, m.Id,
=======
		m := this.GetMember(ctx)
		err = dps.MemberService.SubmitApplyPresentBalance(partnerId, m.Id,
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
			member.TypeApplyCashToBank, float32(amount), saleConf.ApplyCsn)
	}

toErr:
	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
<<<<<<< HEAD

	return ctx.JSON(http.StatusOK, msg)
}

// 转换活动金到提现账户
func (this *accountC) Convert_f2p(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.convert_f2p_post(ctx)
	}
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	saleConf := dps.PartnerService.GetSaleConf(p.Id)
	m := getMember(ctx)
=======
	ctx.Response.JsonOutput(msg)
}

// 转换活动金到提现账户
func (this *accountC) Convert_f2p(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	saleConf := dps.PartnerService.GetSaleConf(p.Id)
	m := this.GetMember(ctx)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	acc := dps.MemberService.GetAccount(m.Id)

	var commissionStr string
	if saleConf.FlowConvertCsn == 0 {
		commissionStr = "不收取手续费"
	} else {
		commissionStr = fmt.Sprintf("收取<i>%s%s</i>手续费",
			format.FormatFloat(saleConf.FlowConvertCsn*100), "%")
	}

<<<<<<< HEAD
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
=======
	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		"conf":              conf,
		"partner":           p,
		"member":            m,
		"account":           acc,
		"commissionStr":     template.HTML(commissionStr),
		"flowAlias":         variable.AliasFlowAccount,
		"flowConvertSlogan": variable.FlowConvertSlogan,
		"cns":               saleConf.FlowConvertCsn,
		"notSetTradePwd":    len(m.TradePwd) == 0,
<<<<<<< HEAD
	}
	return ctx.RenderOK("account.convert_f2p.html", d)
}

func (this *accountC) convert_f2p_post(ctx *echox.Context) error {
	var msg gof.Message
	var err error
	r := ctx.Request()
	r.ParseForm()
	pt := getPartner(ctx)
	amount, _ := strconv.ParseFloat(r.FormValue("Amount"), 32)
	tradePwd := r.FormValue("TradePwd")
	saleConf := dps.PartnerService.GetSaleConf(pt.Id)

	m := getMember(ctx)

	if _, err = dps.MemberService.VerifyTradePwd(m.Id, tradePwd); err == nil {
		err = dps.MemberService.TransferFlow(m.Id, member.KindBalancePresent,
			float32(amount), saleConf.FlowConvertCsn, domain.NewTradeNo(pt.Id),
=======
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
			fmt.Sprintf("%s转换", variable.AliasFlowAccount),
			fmt.Sprintf("%s转换%s", variable.AliasFlowAccount, variable.AliasPresentAccount))
	}

	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
<<<<<<< HEAD
	return ctx.JSON(http.StatusOK, msg)
}

// 转换活动金到提现账户
func (this *accountC) Transfer_f2m(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.transfer_f2m_post(ctx)
	}
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	saleConf := dps.PartnerService.GetSaleConf(p.Id)
	m := getMember(ctx)
=======
	ctx.Response.JsonOutput(msg)
}

// 转换活动金到提现账户
func (this *accountC) Transfer_f2m(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	saleConf := dps.PartnerService.GetSaleConf(p.Id)
	m := this.GetMember(ctx)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	acc := dps.MemberService.GetAccount(m.Id)

	var commissionStr string
	if saleConf.FlowConvertCsn == 0 {
		commissionStr = "不收取手续费"
	} else {
		commissionStr = fmt.Sprintf("收取<i>%s%s</i>手续费",
			format.FormatFloat(saleConf.FlowConvertCsn*100), "%")
	}

<<<<<<< HEAD
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
=======
	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		"conf":           conf,
		"partner":        p,
		"member":         m,
		"account":        acc,
		"commissionStr":  template.HTML(commissionStr),
		"flowAlias":      variable.AliasFlowAccount,
		"cns":            saleConf.TransCsn,
		"notSetTradePwd": len(m.TradePwd) == 0,
<<<<<<< HEAD
	}

	return ctx.RenderOK("account.transfer_f2m.html", d)
}

func (this *accountC) transfer_f2m_post(ctx *echox.Context) error {
	var msg gof.Message
	var err error
	r := ctx.Request()
	r.ParseForm()
	p := getPartner(ctx)
	toMemberId, _ := strconv.Atoi(r.FormValue("ToId"))
	amount, _ := strconv.ParseFloat(r.FormValue("Amount"), 32)
	tradePwd := r.FormValue("TradePwd")
	saleConf := dps.PartnerService.GetSaleConf(p.Id)
	m := getMember(ctx)

	if toMemberId == m.Id {
		err = errors.New("无法转账到自己账号")
	} else {
		if _, err = dps.MemberService.VerifyTradePwd(m.Id, tradePwd); err == nil {
			err = dps.MemberService.TransferFlowTo(m.Id, toMemberId, member.KindBalanceFlow,
				float32(amount), saleConf.TransCsn, domain.NewTradeNo(p.Id),
				variable.AliasFlowAccount+"转账", "转入"+variable.AliasFlowAccount)
		}
	}
=======
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

>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
<<<<<<< HEAD
	return ctx.JSON(http.StatusOK, msg)
}

// 转账成功提示页面
func (this *accountC) Transfer_success(ctx *echox.Context) error {
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)

	src := ctx.Query("src")
=======
	ctx.Response.JsonOutput(msg)
}

// 转账成功提示页面
func (this *accountC) Transfer_success(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)

	src := ctx.Request.URL.Query().Get("src")
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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

<<<<<<< HEAD
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
=======
	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		"conf":     conf,
		"partner":  p,
		"title":    title,
		"subTitle": subTitle,
		"btnText":  btnText,
<<<<<<< HEAD
		"referer":  ctx.Request().Referer(),
	}
	return ctx.RenderOK("account.transfer_success.html", d)
}

func (this *accountC) Bank_info(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.bank_info_post(ctx)
	}
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	m := getMember(ctx)
	bank := dps.MemberService.GetBank(m.Id)

	js, _ := json.Marshal(bank)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"conf":    conf,
		"partner": p,
		"entity":  template.JS(js),
	}
	return ctx.RenderOK("account.bank_info.html", d)
}

func (this *accountC) bank_info_post(ctx *echox.Context) error {
	m := getMember(ctx)
	r := ctx.Request()
=======
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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
<<<<<<< HEAD
	return ctx.JSON(http.StatusOK, msg)
}

func (this *accountC) Integral_exchange(ctx *echox.Context) error {
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	m := getMember(ctx)
	acc := dps.MemberService.GetAccount(m.Id)

	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
=======
	ctx.Response.JsonOutput(msg)
}

func (this *accountC) Integral_exchange(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
	acc, _ := goclient.Member.GetMemberAccount(m.Id, m.DynamicToken)

	this.ExecuteTemplate(ctx, gof.TemplateDataMap{
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		"conf":    conf,
		"record":  15,
		"partner": p,
		"member":  m,
		"account": acc,
<<<<<<< HEAD
	}
	return ctx.RenderOK("account.integral_exchange.html", d)
=======
	}, "views/ucenter/{device}/account/integral_exchange.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
