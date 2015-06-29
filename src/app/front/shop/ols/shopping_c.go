/**
 * Copyright 2014 @ ops.
 * name :
 * author : jarryliu
 * date : 2013-11-26 21:09
 * description :
 * history :
 */
package ols

import (
	"encoding/json"
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/dto"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var _ mvc.Filter = new(ShoppingC)

type ShoppingC struct {
	*BaseC
}

func (this *ShoppingC) prepare(ctx *web.Context) bool {
	return this.BaseC.CheckMemberLogin(ctx)
}

// 订单确认
func (this *ShoppingC) Confirm(ctx *web.Context) {

	if !this.prepare(ctx) {
		return
	}

	r := ctx.Request
	p := this.BaseC.GetPartner(ctx)
	m := this.BaseC.GetMember(ctx)

	siteConf := this.BaseC.GetSiteConf(ctx)
	// 获取购物车
	var cartKey string
	ck, err := r.Cookie("_cart")
	if err == nil {
		cartKey = ck.Value
	}
	cart := dps.ShoppingService.GetShoppingCart(p.Id, m.Id, cartKey)
	if cart.Items == nil || len(cart.Items) == 0 {
		this.OrderEmpty(ctx, p, m, siteConf)
		return
	}

	// 配送地址
	var deliverId int
	var paymentOpt int = 1
	var deliverOpt int = 1
	var settle *dto.SettleMeta = dps.ShoppingService.GetCartSettle(p.Id, m.Id, cart.CartKey)
	if settle.Deliver != nil {
		deliverId = settle.Deliver.Id
		ph := settle.Deliver.Phone
		if len(ph) == 11 {
			settle.Deliver.Phone = strings.Replace(ph, ph[3:7], "****", 1)
		}
		if settle.PaymentOpt > 0 {
			paymentOpt = settle.PaymentOpt
		} else {
			paymentOpt = 1
		}

		if settle.DeliverOpt > 0 {
			deliverOpt = settle.DeliverOpt
		} else {
			deliverOpt = 1
		}
	}

	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"partner":     p,
		"title":       "订单确认-" + p.Name,
		"member":      m,
		"cart":        cart,
		"cartDetails": template.HTML(format.CartDetails(cart)),
		"promFee":     cart.TotalFee - cart.OrderFee,
		"summary":     template.HTML(cart.Summary),
		"conf":        siteConf,
		"settle":      settle,
		"deliverId":   deliverId,
		"deliverOpt":  deliverOpt,
		"paymentOpt":  paymentOpt,
	},
		"views/shop/ols/{device}/order_confirm.html",
		"views/shop/ols/{device}/inc/header.html",
		"views/shop/ols/{device}/inc/footer.html")

}

// 订单持久化
func (this *ShoppingC) BuyingPersist_post(ctx *web.Context) {
	if !this.prepare(ctx) {
		return
	}
	r, w := ctx.Request, ctx.Response
	p := this.BaseC.GetPartner(ctx)
	m := this.BaseC.GetMember(ctx)
	var err error
	r.ParseForm()

	var deliverId int
	var paymentOpt int
	var deliverOpt int
	var shopId int

	deliverId, err = strconv.Atoi(r.FormValue("deliver_id"))
	if err != nil {
		goto rsp
	}
	paymentOpt, err = strconv.Atoi(r.FormValue("pay_opt"))
	if err != nil {
		goto rsp
	}
	deliverOpt, err = strconv.Atoi(r.FormValue("deliver_opt"))
	if err != nil {
		goto rsp
	}
	shopId, err = strconv.Atoi(r.FormValue("shop_id"))
	if err != nil {
		goto rsp
	}

	err = dps.ShoppingService.PrepareSettlePersist(p.Id, m.Id, shopId, paymentOpt, deliverOpt, deliverId)

rsp:
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{result:false,"message":"%s"}`, err.Error())))
	} else {
		w.Write([]byte("{result:true}"))
	}
}

// 配送地址管理
func (this *ShoppingC) GetDeliverAddrs(ctx *web.Context) {
	if !this.prepare(ctx) {
		return
	}
	r := ctx.Request
	m := this.GetMember(ctx)
	addrs := dps.MemberService.GetDeliverAddress(m.Id)
	var selId int
	if sel := r.URL.Query().Get("sel"); sel != "" {
		selId, _ = strconv.Atoi(sel)
	}

	js, _ := json.Marshal(addrs)
	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"addrs": template.JS(js),
		"sel":   selId,
	}, "views/shop/ols/{device}/deliver_address.html")
}
func (this *ShoppingC) SaveDeliverAddr_post(ctx *web.Context) {
	if !this.prepare(ctx) {
		return
	}
	m := this.GetMember(ctx)
	r, w := ctx.Request, ctx.Response
	r.ParseForm()
	var e member.DeliverAddress
	web.ParseFormToEntity(r.Form, &e)
	e.MemberId = m.Id
	b, err := dps.MemberService.SaveDeliverAddress(m.Id, &e)
	if err == nil {
		if b > 0 {
			w.Write([]byte(`{"result":true}`))
		} else {
			w.Write([]byte(`{"result":false}`))
		}
	} else {
		w.Write([]byte(fmt.Sprintf(`{"result":false,"message":"%s"}`, err.Error())))
	}
}

// 应用卡券
func (this *ShoppingC) Apply_post(ctx *web.Context) {
	r := ctx.Request
	r.ParseForm()
	if atype := r.URL.Query().Get("type"); atype == "coupon" {
		this.applyCoupon(ctx)
	}
}
func (this *ShoppingC) applyCoupon(ctx *web.Context) {
	if !this.prepare(ctx) {
		return
	}

	p := this.GetPartner(ctx)
	m := this.GetMember(ctx)

	var message string = "购物车还是空的!"
	code := ctx.Request.FormValue("code")
	order, _, err := dps.ShoppingService.BuildOrder(p.Id, m.Id, "", code)
	if err != nil {
		message = err.Error()
	} else {
		d, _ := json.Marshal(order)
		ctx.Response.Write(d)
		return
	}

	ctx.Response.Write([]byte(`{"result":false,"message":"` + message + `"}`))
}

// 提交订单
func (this *ShoppingC) Submit_0_post(ctx *web.Context) {
	if !this.prepare(ctx) {
		return
	}
	r, w := ctx.Request, ctx.Response
	p := this.GetPartner(ctx)
	m := this.GetMember(ctx)

	r.ParseForm()
	if p == nil || m == nil {
		w.Write([]byte(`{"result":false,"tag":"101"}`)) //未登录
		return
	}
	couponCode := r.FormValue("coupon_code")
	order_no, err := dps.ShoppingService.SubmitOrder(p.Id, m.Id, couponCode)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"result":false,"tag":"109","message":"%s"}`, err.Error())))
		return
	}

	// 清空购物车
	this.emptyShoppingCart(ctx)

	w.Write([]byte(`{"result":true,"data":"` + order_no + `"}`))
}

// 清除购物车
func (this *ShoppingC) emptyShoppingCart(ctx *web.Context) {
	cookie, _ := ctx.Request.Cookie("_cart")
	if cookie != nil {
		cookie.Expires = time.Now().Add(time.Hour * 24 * -30)
		cookie.Path = "/"
		http.SetCookie(ctx.Response, cookie)
	}
}

func (this *ShoppingC) OrderEmpty(ctx *web.Context, p *partner.ValuePartner,
	m *member.ValueMember, conf *partner.SiteConf) {
	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"partner": p,
		"title":   "订单确认-" + p.Name,
		"member":  m,
		"conf":    conf,
	},
		"views/shop/ols/{device}/order_empty.html",
		"views/shop/ols/{device}/inc/header.html",
		"views/shop/ols/{device}/inc/footer.html")
}

func (this *ShoppingC) Order_finish(ctx *web.Context) {
	if !this.prepare(ctx) {
		return
	}
	r := ctx.Request

	p := this.GetPartner(ctx)
	m := this.GetMember(ctx)

	siteConf := this.GetSiteConf(ctx)
	var payHtml string // 支付HTML

	orderNo := r.URL.Query().Get("order_no")
	order := dps.ShoppingService.GetOrderByNo(p.Id, orderNo)
	if order != nil {
		if order.PaymentOpt == 2 {
			payHtml = fmt.Sprintf(`<div class="payment_button"><a href="/pay/create?pay_opt=alipay&order_no=%s" target="_blank">%s</a></div>`,
				order.OrderNo, "在线支付")
		}

		this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
			"partner": p,
			"title":   "订单成功-" + p.Name,
			"member":  m,
			"conf":    siteConf,
			"order":   order,
			"payHtml": template.HTML(payHtml),
		},
			"views/shop/ols/{device}/order_finish.html",
			"views/shop/ols/{device}/inc/header.html",
			"views/shop/ols/{device}/inc/footer.html")
	} else {
		this.OrderEmpty(ctx, p, m, siteConf)
	}

}

// 购买中转
func (this *ShoppingC) Index(ctx *web.Context) {
	if !this.prepare(ctx) {
		return
	}
	w := ctx.Response
	w.Header().Add("Location", "/buy/confirm")
	w.WriteHeader(302)
}
