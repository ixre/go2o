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
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/pager"
	"go2o/src/app/front"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/service/dps"
	"strconv"
)

type orderC struct {
	*baseC
}

func (this *orderC) Complete(ctx *web.Context) {
	this.ExecuteTemplate(ctx, nil,
		"views/ucenter/{device}/order/complete.html")
}

// 所有订单
func (this *orderC) All(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
	this.ExecuteTemplate(ctx,
		gof.TemplateDataMap{
			"partner":      p,
			"conf":         conf,
			"partner_host": conf.Host,
			"member":       m,
		}, "views/ucenter/{device}/order/order_list.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}

func (this *orderC) All_post(ctx *web.Context) {
	ctx.Request.ParseForm()
	this.responseList(ctx, "")
}

func (this *orderC) responseList(ctx *web.Context, where string) {
	m := this.GetMember(ctx)
	r, w := ctx.Request, ctx.Response
	page, _ := strconv.Atoi(r.FormValue("page"))
	size, _ := strconv.Atoi(r.FormValue("size"))

	n, rows := dps.MemberService.QueryPagerOrder(m.Id, page, size, where, "")

	p := pager.NewUrlPager(pager.TotalPage(n, size), page, pager.GetterJavaScriptPager)

	p.RecordCount = n
	pager := &front.Pager{Total: n, Rows: rows, Text: p.PagerString()}

	js, _ := json.Marshal(pager)
	w.Write(js)
}

func (this *orderC) Wait_payment(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
	this.ExecuteTemplate(ctx,
		gof.TemplateDataMap{
			"partner":      p,
			"conf":         conf,
			"partner_host": conf.Host,
			"member":       m,
			"state":        enum.ORDER_WAIT_CONFIRM,
		},
		"views/ucenter/{device}/order/order_wait_payment.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}

func (this *orderC) Wait_payment_post(ctx *web.Context) {
	ctx.Request.ParseForm()
	this.responseList(ctx, fmt.Sprintf("is_paid=0 AND status <> %d", enum.ORDER_CANCEL))
}

func (this *orderC) Wait_delivery(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
	this.ExecuteTemplate(ctx,
		gof.TemplateDataMap{
			"partner":      p,
			"conf":         conf,
			"partner_host": conf.Host,
			"member":       m,
		},
		"views/ucenter/{device}/order/order_wait_delivery.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}

func (this *orderC) Wait_delivery_post(ctx *web.Context) {
	ctx.Request.ParseForm()
	var where string = fmt.Sprintf("is_paid=1 AND deliver_time < create_time AND status =  %d", enum.ORDER_WAIT_DELIVERY)
	this.responseList(ctx, where)
}

func (this *orderC) Completed(ctx *web.Context) {

	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
	this.ExecuteTemplate(ctx,
		gof.TemplateDataMap{
			"partner":      p,
			"conf":         conf,
			"partner_host": conf.Host,
			"member":       m,
			"state":        enum.ORDER_COMPLETED,
		},
		"views/ucenter/{device}/order/order_completed.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}

func (this *orderC) Canceled(ctx *web.Context) {

	p := this.GetPartner(ctx)
	conf := this.GetSiteConf(p.Id)
	m := this.GetMember(ctx)
	this.ExecuteTemplate(ctx,
		gof.TemplateDataMap{
			"partner":      p,
			"conf":         conf,
			"partner_host": conf.Host,
			"member":       m,
			"state":        enum.ORDER_CANCEL,
		},
		"views/ucenter/{device}/order/order_cancel.html",
		"views/ucenter/{device}/inc/header.html",
		"views/ucenter/{device}/inc/menu.html",
		"views/ucenter/{device}/inc/footer.html")
}
