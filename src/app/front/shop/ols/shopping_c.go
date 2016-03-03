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
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
<<<<<<< HEAD
=======
	"github.com/jsix/gof/web/mvc"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/dto"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
<<<<<<< HEAD
	"go2o/src/x/echox"
=======
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

<<<<<<< HEAD
type ShoppingC struct {
}

func (this *ShoppingC) prepare(ctx *echox.Context) bool {
	return CheckMemberLogin(ctx)
}

// 订单确认
func (this *ShoppingC) Confirm(ctx *echox.Context) error {
	if !this.prepare(ctx) {
		return nil
	}
	r := ctx.Request()
	p := getPartner(ctx)
	m := GetMember(ctx)
	siteConf := getSiteConf(ctx)
=======
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d

	// 获取购物车
	var cartKey string
	ck, err := r.Cookie("_cart")
	if err == nil {
		cartKey = ck.Value
	}
	cart := dps.ShoppingService.GetShoppingCart(p.Id, m.Id, cartKey)
	if cart.Items == nil || len(cart.Items) == 0 {
<<<<<<< HEAD
		return this.OrderEmpty(ctx, p, m, siteConf)
=======
		this.OrderEmpty(ctx, p, m, siteConf)
		return
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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

	acc := dps.MemberService.GetAccount(m.Id)

<<<<<<< HEAD
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
=======
	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		"partner":     p,
		"member":      m,
		"cart":        cart,
		"cartDetails": template.HTML(format.CartDetails(cart)),
		"promFee":     cart.TotalFee - cart.OrderFee,
		"balance":     acc.Balance,
		"summary":     template.HTML(cart.Summary),
		"conf":        siteConf,
		"settle":      settle,
		"deliverId":   deliverId,
		"deliverOpt":  deliverOpt,
		"paymentOpt":  paymentOpt,
<<<<<<< HEAD
	}
	return ctx.RenderOK("order_confirm.html", d)
}

// 订单持久化(POST)
func (this *ShoppingC) BuyingPersist(ctx *echox.Context) error {
	if !this.prepare(ctx) {
		return nil
	}
	r := ctx.Request()
	if r.Method == "POST" {
		p := getPartner(ctx)
		m := GetMember(ctx)
		msg := gof.Message{}
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
			msg.Message = err.Error()
		} else {
			msg.Result = true
		}
		return ctx.JSON(http.StatusOK, msg)
	}
	return nil
}

// 配送地址管理
func (this *ShoppingC) GetDeliverAddress(ctx *echox.Context) error {
	if !this.prepare(ctx) {
		return nil
	}
	r := ctx.Request()
	m := GetMember(ctx)
=======
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
func (this *ShoppingC) GetDeliverAddress(ctx *web.Context) {
	if !this.prepare(ctx) {
		return
	}
	r := ctx.Request
	m := this.BaseC.GetMember(ctx)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	address := dps.MemberService.GetDeliverAddress(m.Id)
	var selId int
	if sel := r.URL.Query().Get("sel"); sel != "" {
		selId, _ = strconv.Atoi(sel)
	}

	js, _ := json.Marshal(address)
<<<<<<< HEAD
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"addrs": template.JS(js),
		"sel":   selId,
	}
	return ctx.RenderOK("deliver_address.html", d)
}

// 保存配送地址(POST)
func (this *ShoppingC) SaveDeliverAddress(ctx *echox.Context) error {
	r := ctx.Request()
	if this.prepare(ctx) && r.Method == "POST" {
		msg := gof.Message{Result: true}
		m := GetMember(ctx)
		r.ParseForm()
		var e member.DeliverAddress
		web.ParseFormToEntity(r.Form, &e)
		e.MemberId = m.Id
		_, err := dps.MemberService.SaveDeliverAddress(m.Id, &e)
		if err != nil {
			msg.Result = false
			msg.Message = err.Error()
		}
		return ctx.JSON(http.StatusOK, msg)
	}
	return nil
}

// 应用卡券(POST)
func (this *ShoppingC) Apply(ctx *echox.Context) error {
	r := ctx.Request()
	if this.prepare(ctx) && r.Method == "POST" {
		r.ParseForm()
		applyType := r.URL.Query().Get("type")
		if applyType == "coupon" {
			this.applyCoupon(ctx)
		}
	}
	return nil
}

func (this *ShoppingC) applyCoupon(ctx *echox.Context) error {
	msg := gof.Message{}
	p := getPartner(ctx)
	m := GetMember(ctx)
	r := ctx.Request()

	code := r.FormValue("code")
	subject := r.FormValue("subject") // not necessary
	data, err := dps.ShoppingService.BuildOrder(p.Id, subject, m.Id, "", code)
	if err != nil {
		msg.Message = err.Error()
	} else {
		return ctx.JSON(http.StatusOK, data)
	}
	return ctx.JSON(http.StatusOK, msg)
}

// 提交订单
func (this *ShoppingC) Submit_0(ctx *echox.Context) error {
	r := ctx.Request()
	if this.prepare(ctx) && r.Method == "POST" {
		p := getPartner(ctx)
		m := GetMember(ctx)

		r.ParseForm()
		if p == nil || m == nil {
			return ctx.StringOK(`{"result":false,"tag":"101"}`) //未登录
		}
		couponCode := r.FormValue("coupon_code")
		subject := r.FormValue("subject") // not necessary

		//this.releaseOrder(ctx)
		// 锁定订单提交
		if !this.lockOrder(ctx) {
			//fmt.Println("--- IS LOCKED")
			return ctx.JSON(http.StatusOK, gof.Message{Message: "请勿重复提交订单"})
		}

		// 是否余额支付
		var useBalanceDiscount bool = r.FormValue("balance_discount") == "1"

		// 提交订单
		order_no, err := dps.ShoppingService.SubmitOrder(p.Id, m.Id,
			subject, couponCode, useBalanceDiscount)

		// 释放订单占用
		this.releaseOrder(ctx)

		if err != nil {
			return ctx.StringOK(fmt.Sprintf(`{"result":false,"tag":"109","message":"%s"}`,
				err.Error()))
		}

		// 清空购物车
		this.emptyShoppingCart(ctx)

		ctx.StringOK(`{"result":true,"data":"` + order_no + `"}`)
		this.releaseOrder(ctx)
	}
	return nil
}

// 锁定，防止重复下单，返回false,表示正在处理订单
func (this *ShoppingC) lockOrder(ctx *echox.Context) bool {
	s := ctx.Session
=======
	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"addrs": template.JS(js),
		"sel":   selId,
	}, "views/shop/ols/{device}/deliver_address.html")
}
func (this *ShoppingC) SaveDeliverAddress_post(ctx *web.Context) {
	if !this.prepare(ctx) {
		return
	}
	m := this.BaseC.GetMember(ctx)
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
	if applyType := r.URL.Query().Get("type"); applyType == "coupon" {
		this.applyCoupon(ctx)
	}
}

func (this *ShoppingC) applyCoupon(ctx *web.Context) {
	var result gof.Message
	if !this.prepare(ctx) {
		result.Message = "请先登陆"
	} else {
		p := this.BaseC.GetPartner(ctx)
		m := this.BaseC.GetMember(ctx)

		code := ctx.Request.FormValue("code")
		data, err := dps.ShoppingService.BuildOrder(p.Id, m.Id, "", code)
		if err != nil {
			result.Message = err.Error()
		} else {
			ctx.Response.JsonOutput(data)
			return
		}
	}

	ctx.Response.JsonOutput(result)
}

// 提交订单
func (this *ShoppingC) Submit_0_post(ctx *web.Context) {
	if !this.prepare(ctx) {
		return
	}

	r, w := ctx.Request, ctx.Response
	p := this.BaseC.GetPartner(ctx)
	m := this.BaseC.GetMember(ctx)

	r.ParseForm()
	if p == nil || m == nil {
		w.Write([]byte(`{"result":false,"tag":"101"}`)) //未登录
		return
	}
	couponCode := r.FormValue("coupon_code")

	//this.releaseOrder(ctx)
	// 锁定订单提交
	if !this.lockOrder(ctx) {
		//fmt.Println("--- IS LOCKED")
		ctx.Response.JsonOutput(gof.Message{Message: "请勿重复提交订单"})
		return
	}

	// 是否余额支付
	var useBalanceDiscount bool = ctx.Request.FormValue("balance_discount") == "1"

	// 提交订单
	order_no, err := dps.ShoppingService.SubmitOrder(p.Id, m.Id, couponCode, useBalanceDiscount)

	// 释放订单占用
	this.releaseOrder(ctx)

	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"result":false,"tag":"109","message":"%s"}`, err.Error())))
		return
	}

	// 清空购物车
	this.emptyShoppingCart(ctx)

	w.Write([]byte(`{"result":true,"data":"` + order_no + `"}`))
	this.releaseOrder(ctx)
}

// 锁定，防止重复下单，返回false,表示正在处理订单
func (this *ShoppingC) lockOrder(ctx *web.Context) bool {
	s := ctx.Session()
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	v := s.Get("shopping_lock")

	//fmt.Println(v)
	if v != nil {
		return false
	}
	s.Set("shopping_lock", "1")
	s.Save()
	return true
}
<<<<<<< HEAD
func (this *ShoppingC) releaseOrder(ctx *echox.Context) {
	s := ctx.Session
	s.Remove("shopping_lock")
	s.Save()
=======
func (this *ShoppingC) releaseOrder(ctx *web.Context) {
	s := ctx.Session()
	s.Remove("shopping_lock")
	s.Save()

>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	//fmt.Println("REMOVED")
}

// 清除购物车
<<<<<<< HEAD
func (this *ShoppingC) emptyShoppingCart(ctx *echox.Context) {
	cookie, _ := ctx.Request().Cookie("_cart")
	if cookie != nil {
		cookie.Expires = time.Now().Add(time.Hour * 24 * -30)
		cookie.Path = "/"
		http.SetCookie(ctx.Response(), cookie)
	}
}

func (this *ShoppingC) OrderEmpty(ctx *echox.Context, p *partner.ValuePartner,
	m *member.ValueMember, conf *partner.SiteConf) error {
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner": p,
		"member":  m,
		"conf":    conf,
	}
	return ctx.RenderOK("order_empty.html", d)
}

func (this *ShoppingC) Payment(ctx *echox.Context) error {
	if !this.prepare(ctx) {
		return nil
	}
	r := ctx.Request()

	p := getPartner(ctx)
	m := GetMember(ctx)
	siteConf := getSiteConf(ctx)
=======
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
		"member":  m,
		"conf":    conf,
	},
		"views/shop/ols/{device}/order_empty.html",
		"views/shop/ols/{device}/inc/header.html",
		"views/shop/ols/{device}/inc/footer.html")
}

func (this *ShoppingC) Payment(ctx *web.Context) {
	if !this.prepare(ctx) {
		return
	}
	r := ctx.Request

	p := this.BaseC.GetPartner(ctx)
	m := this.BaseC.GetMember(ctx)

	siteConf := this.BaseC.GetSiteConf(ctx)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	var payHtml string // 支付HTML
	var payOpt string
	var payHelp string

	orderNo := r.URL.Query().Get("order_no")
	order := dps.ShoppingService.GetOrderByNo(p.Id, orderNo)

	// 已经支付成功
	if order.IsPaid == 1 {
<<<<<<< HEAD
		return this.orderFinish(ctx, p, m, siteConf, order)
=======
		this.orderFinish(ctx, p, m, siteConf, order)
		return
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	}

	if order != nil {
		if order.PaymentOpt == enum.PaymentOnlinePay {
			payHtml = fmt.Sprintf(`<div class="btn_payment"><a class="btn btn-m"
				href="/pay/create?pay_opt=alipay&order_no=%s" target="_blank">%s</a></div>`,
				order.OrderNo, "立即支付")
		}
		payOpt = enum.GetPaymentName(order.PaymentOpt)
		payHelp = enum.GetPaymentHelpContent(order.PaymentOpt)

<<<<<<< HEAD
		d := ctx.NewData()
		d.Map = gof.TemplateDataMap{
=======
		this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
			"partner":      p,
			"member":       m,
			"conf":         siteConf,
			"order":        order,
			"payment_opt":  payOpt,
			"payment_html": template.HTML(payHtml),
			"payment_help": template.HTML(payHelp),
<<<<<<< HEAD
		}
		return ctx.RenderOK("order_payment.html", d)
	}
	return this.OrderEmpty(ctx, p, m, siteConf)
}

func (this *ShoppingC) Order_finish(ctx *echox.Context) error {
	p := getPartner(ctx)
	m := GetMember(ctx)
	siteConf := getSiteConf(ctx)
	orderNo := ctx.Query("order_no")
	order := dps.ShoppingService.GetOrderByNo(p.Id, orderNo)
	return this.orderFinish(ctx, p, m, siteConf, order)
}

func (this *ShoppingC) orderFinish(ctx *echox.Context, p *partner.ValuePartner,
	m *member.ValueMember, siteConf *partner.SiteConf, order *shopping.ValueOrder) error {
	if !this.prepare(ctx) {
		return nil
	}
	var payHtml string // 支付HTML
	var payHelp string
=======
		},
			"views/shop/ols/{device}/order_payment.html",
			"views/shop/ols/{device}/inc/header.html",
			"views/shop/ols/{device}/inc/footer.html")
	} else {
		this.OrderEmpty(ctx, p, m, siteConf)
	}
}

func (this *ShoppingC) Order_finish(ctx *web.Context) {
	p := this.BaseC.GetPartner(ctx)
	m := this.BaseC.GetMember(ctx)
	siteConf := this.BaseC.GetSiteConf(ctx)

	orderNo := ctx.Request.URL.Query().Get("order_no")
	order := dps.ShoppingService.GetOrderByNo(p.Id, orderNo)

	this.orderFinish(ctx, p, m, siteConf, order)
}

func (this *ShoppingC) orderFinish(ctx *web.Context, p *partner.ValuePartner,
	m *member.ValueMember, siteConf *partner.SiteConf, order *shopping.ValueOrder) {
	if !this.prepare(ctx) {
		return
	}

	var payHtml string // 支付HTML
	var payHelp string

>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	if order != nil {
		if order.IsPaid == 0 {
			if order.PaymentOpt == enum.PaymentOnlinePay {
				payHtml = fmt.Sprintf(`<div class="btn_payment"><a class="btn btn-m"
					href="/pay/create?pay_opt=alipay&order_no=%s" target="_blank">%s</a></div>`,
					order.OrderNo, "继续支付")
			}
		}
<<<<<<< HEAD
		d := ctx.NewData()
		d.Map = gof.TemplateDataMap{
=======

		this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
			"partner":      p,
			"member":       m,
			"conf":         siteConf,
			"order":        order,
			"payment_html": template.HTML(payHtml),
			"payment_help": template.HTML(payHelp),
<<<<<<< HEAD
		}
		return ctx.RenderOK("order_finish.html", d)
	}
	return this.OrderEmpty(ctx, p, m, siteConf)
}

// 购买中转
func (this *ShoppingC) Index(ctx *echox.Context) error {
	ctx.Response().Header().Add("Location", "/buy/confirm")
	ctx.Response().WriteHeader(302)
	return nil
=======
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
