package ucenter

import (
	"com/domain/interface/enum"
	"com/domain/interface/member"
	"com/ording"
	"com/ording/dao"
	"com/ording/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"ops/cf/app"
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
	p *entity.Partner, conf *entity.SiteConf) {
	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {
			(*m)["partner"] = p
			(*m)["conf"] = conf
			(*m)["partner_host"] = conf.Host
			(*m)["member"] = m
		}, "views/ucenter/order/orders.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *orderC) Orders_post(w http.ResponseWriter, r *http.Request, m *member.ValueMember) {
	r.ParseForm()
	page, _ := strconv.Atoi(r.FormValue("page"))
	size, _ := strconv.Atoi(r.FormValue("size"))
	status, err := strconv.Atoi(r.FormValue("status"))
	if err != nil {
		status = enum.ORDER_CREATED
	}

	n, rows := dao.Order().GetMemberPagerOrder(m.Id, page, size,
		fmt.Sprintf("status=%d", status), "")

	pager := &ording.Pager{Total: n, Rows: rows}

	js, _ := json.Marshal(pager)
	w.Write(js)
}
