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
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web/pager"
	"go2o/src/app/front"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"net/http"
	"strconv"
)

type orderC struct {
}

func (this *orderC) Complete(ctx *echox.Context) error {
	return ctx.RenderOK("order.complete.html", ctx.NewData())
}

// 所有订单
func (this *orderC) All(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.all_post(ctx)
	}
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	m := getMember(ctx)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"member":       m,
	}
	return ctx.RenderOK("order.list.html", d)
}

func (this *orderC) all_post(ctx *echox.Context) error {
	return this.responseList(ctx, "")
}

func (this *orderC) responseList(ctx *echox.Context, where string) error {
	m := getMember(ctx)
	r := ctx.HttpRequest()
	r.ParseForm()
	page, _ := strconv.Atoi(r.FormValue("page"))
	size, _ := strconv.Atoi(r.FormValue("size"))

	n, rows := dps.MemberService.QueryPagerOrder(m.Id, page, size, where, "")

	p := pager.NewUrlPager(pager.MathPages(n, size), page, "")

	p.Total = n
	pager := &front.Pager{Total: n, Rows: rows, Text: p.PagerString()}

	return ctx.JSON(http.StatusOK, pager)
}

func (this *orderC) Wait_payment(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.wait_payment_post(ctx)
	}
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	m := getMember(ctx)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"member":       m,
		"state":        enum.ORDER_WAIT_CONFIRM,
	}
	return ctx.RenderOK("order.wait_payment.html", d)
}

func (this *orderC) wait_payment_post(ctx *echox.Context) error {
	return this.responseList(ctx, fmt.Sprintf("is_paid=0 AND status <> %d", enum.ORDER_CANCEL))
}

func (this *orderC) Wait_delivery(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.wait_payment_post(ctx)
	}
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	m := getMember(ctx)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"member":       m,
	}
	return ctx.RenderOK("order.wait_delivery.html", d)
}

func (this *orderC) wait_delivery_post(ctx *echox.Context) error {
	var where string = fmt.Sprintf("is_paid=1 AND deliver_time < create_time AND status =  %d", enum.ORDER_WAIT_DELIVERY)
	return this.responseList(ctx, where)
}

func (this *orderC) Completed(ctx *echox.Context) error {

	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	m := getMember(ctx)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"member":       m,
		"state":        enum.ORDER_COMPLETED,
	}
	return ctx.RenderOK("order.completed.html", d)
}

func (this *orderC) Canceled(ctx *echox.Context) error {
	p := getPartner(ctx)
	conf := getSiteConf(p.Id)
	m := getMember(ctx)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner":      p,
		"conf":         conf,
		"partner_host": conf.Host,
		"member":       m,
		"state":        enum.ORDER_CANCEL,
	}
	return ctx.RenderOK("order.cancel.html", d)
}
