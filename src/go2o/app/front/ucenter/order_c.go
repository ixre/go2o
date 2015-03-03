/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	"encoding/json"
	"fmt"
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web/pager"
	"go2o/app/front"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/partner"
	"go2o/core/service/dps"
	"net/http"
	"strconv"
)

type orderC struct {
	app.Context
}

func (this *orderC) Complete(w http.ResponseWriter, r *http.Request, memberId int) {
	this.Context.Template().Render(w,
		"views/ucenter/order/complete.html",
		func(m *map[string]interface{}) {

		})
}

func (this *orderC) Orders(w http.ResponseWriter, r *http.Request, m *member.ValueMember,
	p *partner.ValuePartner, conf *partner.SiteConf) {
	this.Context.Template().Execute(w,
		func(mp *map[string]interface{}) {
			v := *mp
			v["partner"] = p
			v["conf"] = conf
			v["partner_host"] = conf.Host
			v["member"] = m
		}, "views/ucenter/order/order_list.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *orderC) Completed(w http.ResponseWriter, r *http.Request, m *member.ValueMember,
	p *partner.ValuePartner, conf *partner.SiteConf) {

	this.Context.Template().Execute(w,
		func(mp *map[string]interface{}) {
			v := *mp
			v["partner"] = p
			v["conf"] = conf
			v["partner_host"] = conf.Host
			v["member"] = m
			v["state"] = enum.ORDER_COMPLETED
		},
		"views/ucenter/order/order_completed.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *orderC) Canceled(w http.ResponseWriter, r *http.Request, m *member.ValueMember,
	p *partner.ValuePartner, conf *partner.SiteConf) {

	this.Context.Template().Execute(w,
		func(mp *map[string]interface{}) {
			v := *mp
			v["partner"] = p
			v["conf"] = conf
			v["partner_host"] = conf.Host
			v["member"] = m
			v["state"] = enum.ORDER_CANCEL
		},
		"views/ucenter/order/order_cancel.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *orderC) Orders_post(w http.ResponseWriter, r *http.Request, m *member.ValueMember) {
	r.ParseForm()
	page, _ := strconv.Atoi(r.FormValue("page"))
	size, _ := strconv.Atoi(r.FormValue("size"))
	state := r.FormValue("state")

	var where string
	if state != "" {
		where = fmt.Sprintf("status IN (%s)", state)
	}

	n, rows := dps.MemberService.QueryPagerOrder(m.Id, page, size, where, "")

	p := pager.NewUrlPager(pager.TotalPage(n, size), page, pager.JavaScriptPagerGetter)

	pager := &front.Pager{Total: n, Rows: rows, Text: p.PagerString()}

	js, _ := json.Marshal(pager)
	w.Write(js)
}
