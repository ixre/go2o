/**
 * Copyright 2014 @ ops.
 * name :
 * author : jarryliu
 * date : 2013-11-26 21:09
 * description :
 * history :
 */
package www

import (
	"encoding/json"
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/dto"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/goclient"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type shoppingC struct {
	gof.App
}

// 订单确认
func (this *shoppingC) OrderConfirm(ctx *web.Context,
	p *partner.ValuePartner, m *member.ValueMember) {
	r, w := ctx.Request, ctx.ResponseWriter
	if m == nil {
		RedirectLoginPage(w, r.RequestURI)
		return
	}

	if b, siteConf := GetSiteConf(w, p); b {
		// 获取购物车
		var cartKey string
		ck, err := r.Cookie("_cart")
		if err == nil {
			cartKey = ck.Value
		}
		cart := goclient.Partner.GetShoppingCart(p.Id, m.Id, cartKey)
		if cart.Items == nil || len(cart.Items) == 0 {
			this.OrderEmpty(ctx, p, m, siteConf)
			return
		}

		// 配送地址
		var deliverId int
		var paymentOpt int = 1
		var deliverOpt int = 1
		var settle *dto.SettleMeta = goclient.Partner.GetCartSettle(p.Id, m.Id, cart.CartKey)
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

		this.App.Template().Execute(w, func(mp *map[string]interface{}) {
			(*mp)["partner"] = p
			(*mp)["title"] = "订单确认-" + p.Name
			(*mp)["member"] = m
			(*mp)["cart"] = cart
			(*mp)["cartDetails"] = template.HTML(format.CartDetails(cart))
			(*mp)["promFee"] = cart.TotalFee - cart.OrderFee
			(*mp)["summary"] = template.HTML(cart.Summary)
			(*mp)["conf"] = siteConf
			(*mp)["settle"] = settle
			(*mp)["deliverId"] = deliverId
			(*mp)["deliverOpt"] = deliverOpt
			(*mp)["paymentOpt"] = paymentOpt
		},
			"views/web/www/order_confirm.html",
			"views/web/www/inc/header.html",
			"views/web/www/inc/footer.html")
	}
}

func (this *shoppingC) GetDeliverAddrs(ctx *web.Context,
	p *partner.ValuePartner, m *member.ValueMember) {
	r, w := ctx.Request, ctx.ResponseWriter
	addrs, err := goclient.Member.GetDeliverAddrs(m.Id, m.LoginToken)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	var selId int
	if sel := r.URL.Query().Get("sel"); sel != "" {
		selId, _ = strconv.Atoi(sel)
	}

	js, _ := json.Marshal(addrs)
	this.App.Template().Execute(w, func(md *map[string]interface{}) {
		(*md)["addrs"] = template.JS(js)
		(*md)["sel"] = selId
	}, "views/web/www/profile/deliver_address.html")
}

func (this *shoppingC) SaveDeliverAddr_post(w http.ResponseWriter,
	r *http.Request, p *partner.ValuePartner, m *member.ValueMember) {
	r.ParseForm()
	var e member.DeliverAddress
	web.ParseFormToEntity(r.Form, &e)
	e.MemberId = m.Id
	b, err := goclient.Member.SaveDeliverAddr(m.Id, m.LoginToken, &e)
	if err == nil {
		if b {
			w.Write([]byte(`{"result":true}`))
		} else {
			w.Write([]byte(`{"result":false}`))
		}
	} else {
		w.Write([]byte(fmt.Sprintf(`{"result":false,"message":"%s"}`, err.Error())))
	}
}

func (this *shoppingC) ApplyCoupon_post(ctx *web.Context,
	p *partner.ValuePartner, m *member.ValueMember) {
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	var message string = "购物车还是空的!"
	code := r.FormValue("code")
	json, err := goclient.Partner.BuildOrder(p.Id, p.Secret, m.Id, code)
	if err != nil {
		message = err.Error()
	} else {
		w.Write([]byte(json))
		return
	}

	w.Write([]byte(`{"result":false,"message":"` + message + `"}`))

}

func (this *shoppingC) OrderEmpty(ctx *web.Context,
	p *partner.ValuePartner, mm *member.ValueMember, conf *partner.SiteConf) {
	this.App.Template().Execute(ctx.ResponseWriter, func(m *map[string]interface{}) {
		(*m)["partner"] = p
		(*m)["title"] = "订单确认-" + p.Name
		(*m)["member"] = mm
		(*m)["conf"] = conf
	},
		"views/web/www/order_empty.html",
		"views/web/www/inc/header.html",
		"views/web/www/inc/footer.html")
}

func (this *shoppingC) OrderFinish(ctx *web.Context,
	p *partner.ValuePartner, mm *member.ValueMember) {
	r, w := ctx.Request, ctx.ResponseWriter
	// 清除购物车
	//	cookie, _ := r.Cookie("cart")
	//	if cookie != nil {
	//		cookie.Expires = time.Now().Add(time.Hour * 24 * -30)
	//		cookie.Path = "/"
	//		http.SetCookie(w, cookie)
	//	}

	if b, siteConf := GetSiteConf(w, p); b {
		var payHtml string // 支付HTML

		orderNo := r.URL.Query().Get("order_no")
		order, err := goclient.Partner.GetOrderByNo(p.Id, p.Secret, orderNo)
		if err != nil {
			this.App.Log().PrintErr(err)
			this.OrderEmpty(ctx, p, mm, siteConf)
			return
		}

		if order.PaymentOpt == 2 {
			payHtml = fmt.Sprintf(`<div class="payment_button"><a href="/pay/create?pay_opt=alipay&order_no=%s" target="_blank">%s</a></div>`,
				order.OrderNo, "在线支付")
		}

		if b, siteConf := GetSiteConf(w, p); b {
			this.App.Template().Execute(w, func(m *map[string]interface{}) {
				(*m)["partner"] = p
				(*m)["title"] = "订单成功-" + p.Name
				(*m)["member"] = mm
				(*m)["conf"] = siteConf
				(*m)["order"] = order
				(*m)["payHtml"] = template.HTML(payHtml)
			},
				"views/web/www/order_finish.html",
				"views/web/www/inc/header.html",
				"views/web/www/inc/footer.html")
		}
	}
}

func (this *shoppingC) OrderPersist_post(ctx *web.Context,
	p *partner.ValuePartner, m *member.ValueMember) {
	r, w := ctx.Request, ctx.ResponseWriter
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

	err = goclient.Partner.OrderPersist(p.Id, m.Id, shopId, paymentOpt, deliverOpt, deliverId)

rsp:
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
	} else {
		w.Write([]byte("{}"))
	}
}

func (this *shoppingC) SubmitOrder_post(ctx *web.Context,
	p *partner.ValuePartner, mm *member.ValueMember) {
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	if p == nil || mm == nil {
		w.Write([]byte(`{"result":false,"tag":"101"}`)) //未登录
		return
	}
	couponCode := r.FormValue("coupon_code")
	order_no, err := goclient.Partner.SubmitOrder(p.Id, p.Secret, mm.Id, couponCode)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"result":false,"tag":"109","message":"%s"}`, err.Error())))
		return
	}

	w.Write([]byte(`{"result":true,"data":"` + order_no + `"}`))
}

// 购物车
func (this *shoppingC) CartApi(ctx *web.Context,
	p *partner.ValuePartner, m *member.ValueMember) {
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	var action = strings.ToLower(r.FormValue("action"))
	var cartKey = r.FormValue("cart.key")
	var memberId int
	if m != nil {
		memberId = m.Id
	}

	switch action {
	case "get":
		this.Cart_GetCart(w, p, memberId, cartKey)
	case "add":
		this.Cart_AddItem(ctx, p, memberId, cartKey)
	case "remove":
		this.Cart_RemoveItem(ctx, p, memberId, cartKey)
	}
}

func (this *shoppingC) Cart_GetCart(w http.ResponseWriter,
	p *partner.ValuePartner, memberId int, cartKey string) {
	cart := goclient.Partner.GetShoppingCart(p.Id, memberId, cartKey)
	d, _ := json.Marshal(cart)
	w.Write(d)
}

func (this *shoppingC) Cart_AddItem(ctx *web.Context,
	p *partner.ValuePartner, memberId int, cartKey string) {
	r, w := ctx.Request, ctx.ResponseWriter
	goodsId, _ := strconv.Atoi(r.FormValue("id"))
	num, _ := strconv.Atoi(r.FormValue("num"))
	item, err := goclient.Partner.AddCartItem(p.Id, memberId, cartKey, goodsId, num)

	var result = make(map[string]interface{}, 2)
	if err != nil {
		result["message"] = err.Error()
	} else {
		result["message"] = ""
		result["item"] = item
	}
	d, _ := json.Marshal(result)
	w.Write(d)
}

func (this *shoppingC) Cart_RemoveItem(ctx *web.Context,
	p *partner.ValuePartner, memberId int, cartKey string) {
	r, w := ctx.Request, ctx.ResponseWriter
	goodsId, _ := strconv.Atoi(r.FormValue("id"))
	num, _ := strconv.Atoi(r.FormValue("num"))
	err := goclient.Partner.SubCartItem(p.Id, memberId, cartKey, goodsId, num)
	if err != nil {
		w.Write([]byte(`{error:'` + err.Error() + `'}`))
	} else {
		w.Write([]byte("{}"))
	}
}

func (this *shoppingC) Cart(ctx *web.Context,
	p *partner.ValuePartner) {
	r, w := ctx.Request, ctx.ResponseWriter
	//todo: 需页面
	if r.URL.Query().Get("edit") == "1" {
		w.Header().Add("Location", "/list")
	} else {
		w.Header().Add("Location", "/buy/confirm")
	}
	w.WriteHeader(302)
}

// 购买中转
func (this *shoppingC) BuyRedirect(ctx *web.Context,
	p *partner.ValuePartner, mm *member.ValueMember) {
	r, w := ctx.Request, ctx.ResponseWriter
	if mm == nil {
		RedirectLoginPage(w, r.RequestURI)
	} else {
		w.Header().Add("Location", "/buy/confirm")
		w.WriteHeader(302)
	}
}
