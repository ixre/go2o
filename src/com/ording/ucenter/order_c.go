package ucenter

import (
	"com/domain/interface/enum"
	"com/domain/interface/member"
	"com/domain/interface/partner"
	"com/ording"
	"com/ording/dao"
	"com/ording/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/newmin/gof/app"
	"github.com/newmin/gof/web/pager"
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
	p *entity.Partner, conf *partner.SiteConf) {
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
	p *entity.Partner, conf *partner.SiteConf) {

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
	p *entity.Partner, conf *partner.SiteConf) {

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
	state, err := strconv.Atoi(r.FormValue("state"))
	if err != nil {
		state = enum.ORDER_CREATED
	}

	n, rows := dao.Order().GetMemberPagerOrder(m.Id, page, size,
		fmt.Sprintf("status=%d", state), "")

	p := pager.NewUrlPager(pager.TotalPage(n, size), page, pager.JavaScriptPagerGetter)

	pager := &ording.Pager{Total: n, Rows: rows, Text: p.PagerString()}

	js, _ := json.Marshal(pager)
	w.Write(js)
}
