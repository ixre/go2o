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
	"fmt"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/pager"
	"go2o/src/app/front"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/service/dps"
	"strconv"
)

type orderC struct {
	*baseC
}

func (this *orderC) Complete(ctx *web.Context) {
	ctx.App.Template().Render(ctx.ResponseWriter,
		"views/ucenter/order/complete.html",
		func(m *map[string]interface{}) {

		})
}

func (this *orderC) Orders(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf, _ := this.GetSiteConf(p.Id, p.Secret)
	m := this.GetMember(ctx)
	ctx.App.Template().Execute(ctx.ResponseWriter,
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

func (this *orderC) Completed(ctx *web.Context) {

	p := this.GetPartner(ctx)
	conf, _ := this.GetSiteConf(p.Id, p.Secret)
	m := this.GetMember(ctx)
	ctx.App.Template().Execute(ctx.ResponseWriter,
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

func (this *orderC) Canceled(ctx *web.Context) {

	p := this.GetPartner(ctx)
	conf, _ := this.GetSiteConf(p.Id, p.Secret)
	m := this.GetMember(ctx)
	ctx.App.Template().Execute(ctx.ResponseWriter,
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

func (this *orderC) Orders_post(ctx *web.Context) {
	m := this.GetMember(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
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
