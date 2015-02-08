package ucenter

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/pager"
	"go2o/app/front"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/partner"
	"go2o/core/service/dps"
	"go2o/core/service/goclient"
	"html/template"
	"net/http"
	"strconv"
)

type accountC struct {
	app.Context
}

func (this *accountC) IncomeLog(w http.ResponseWriter, r *http.Request,
	m *member.ValueMember, p *partner.ValuePartner, conf *partner.SiteConf) {

	this.Context.Template().Execute(w, func(mp *map[string]interface{}) {
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

func (this *accountC) IncomeLog_post(w http.ResponseWriter, r *http.Request,
	m *member.ValueMember) {

	r.ParseForm()
	page, _ := strconv.Atoi(r.FormValue("page"))
	size, _ := strconv.Atoi(r.FormValue("size"))

	n, rows := dps.MemberService.QueryIncomeLog(m.Id, page, size, "", "record_time DESC")

	p := pager.NewUrlPager(pager.TotalPage(n, size), page, pager.JavaScriptPagerGetter)

	pager := &front.Pager{Total: n, Rows: rows, Text: p.PagerString()}

	js, _ := json.Marshal(pager)
	w.Write(js)
}

func (this *accountC) ApplyCash(w http.ResponseWriter, r *http.Request,
	m *member.ValueMember, p *partner.ValuePartner, conf *partner.SiteConf) {
	acc, err := goclient.Member.GetMemberAccount(m.Id, m.LoginToken)
	bank, err := goclient.Member.GetBankInfo(m.Id, m.LoginToken)

	if err != nil {
		w.Write([]byte("error:" + err.Error()))
		return
	}

	js, _ := json.Marshal(bank)
	this.Context.Template().Execute(w, func(m *map[string]interface{}) {
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

func (this *accountC) ApplyCash_post(w http.ResponseWriter, r *http.Request,
	m *member.ValueMember, p *partner.ValuePartner, conf *partner.SiteConf) {
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

func (this *accountC) IntegralExchange(w http.ResponseWriter, r *http.Request,
	m *member.ValueMember, p *partner.ValuePartner, conf *partner.SiteConf) {

	acc, _ := goclient.Member.GetMemberAccount(m.Id, m.LoginToken)

	this.Context.Template().Execute(w, func(m *map[string]interface{}) {
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
