/**
 * Copyright 2014 @ ops.
 * name :
 * author : newmin
 * date : 2013-11-26 21:09
 * description :
 * history :
 */
package www

import (
	"encoding/json"
	"fmt"
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/partner"
	"go2o/core/infrastructure/domain"
	"go2o/core/service/goclient"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type shoppingC struct {
	app.Context
}

func (this *shoppingC) GetDeliverAddrs(w http.ResponseWriter, r *http.Request,
	p *partner.ValuePartner, m *member.ValueMember) {
	addrs, err := goclient.Member.GetDeliverAddrs(m.Id, m.LoginToken)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	js, _ := json.Marshal(addrs)
	this.Context.Template().Execute(w, func(md *map[string]interface{}) {
		(*md)["addrs"] = template.JS(js)
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

func (this *shoppingC) ApplyCoupon_post(w http.ResponseWriter, r *http.Request,
	p *partner.ValuePartner, m *member.ValueMember) {

	r.ParseForm()
	var message string = "购物车还是空的!"
	code := r.FormValue("code")
	ck, err := r.Cookie("cart")
	if err == nil {
		cartData := domain.CartCookieFmt(ck.Value)
		if len(cartData) != 0 {
			json, err := goclient.Partner.BuildOrder(p.Id,
				p.Secret, m.Id, cartData, code)
			if err != nil {
				message = err.Error()
			} else {
				w.Write([]byte(json))
				return
			}
		}
	}

	w.Write([]byte(`{"result":false,"message":"` + message + `"}`))

}

func (this *shoppingC) Order(w http.ResponseWriter, r *http.Request,
	p *partner.ValuePartner, mm *member.ValueMember) {
	//cart=%u91CE%u5C71%u6912%u7092%u8089*5*1*2|2
	if b, siteConf := GetSiteConf(w, p); b {
		ck, err := r.Cookie("cart")
		if err != nil || ck.Value == "" {
			this.Context.Log().PrintErr(err)
			this.OrderEmpty(w, r, p, mm, siteConf)
			return
		}

		//cartData := domain.CartCookieFmt(ck.Value)

		//todo:

		//cart, err := goclient.Partner.GetShoppingCart(p.Id, p.Secret, cartData)
		//if err != nil {
		//	w.Write([]byte("订单异常，请清空重新下单"))
		//	return
		//}

		//		cart.Summary = strings.Replace(cart.Summary, "\n", "<br />", -1)
		//
		//		this.Context.Template().Execute(w, func(m *map[string]interface{}) {
		//			(*m)["partner"] = p
		//			(*m)["title"] = "订单确认-" + p.Name
		//			(*m)["member"] = mm
		//			(*m)["cart"] = cart
		//			(*m)["promFee"] = cart.TotalFee - cart.OrderFee
		//			(*m)["summary"] = template.HTML(cart.Summary)
		//			(*m)["conf"] = siteConf
		//		},
		//			"views/web/www/order_confirm.html",
		//			"views/web/www/inc/header.html",
		//			"views/web/www/inc/footer.html")
	}
}

func (this *shoppingC) OrderEmpty(w http.ResponseWriter, r *http.Request,
	p *partner.ValuePartner, mm *member.ValueMember, conf *partner.SiteConf) {
	this.Context.Template().Execute(w, func(m *map[string]interface{}) {
		(*m)["partner"] = p
		(*m)["title"] = "订单确认-" + p.Name
		(*m)["member"] = mm
		(*m)["conf"] = conf
	},
		"views/web/www/order_empty.html",
		"views/web/www/inc/header.html",
		"views/web/www/inc/footer.html")
}

func (this *shoppingC) Finish(w http.ResponseWriter, r *http.Request,
	p *partner.ValuePartner, mm *member.ValueMember) {
	// 清除购物车
	cookie, _ := r.Cookie("cart")
	if cookie != nil {
		cookie.Expires = time.Now().Add(time.Hour * 24 * -30)
		cookie.Path = "/"
		http.SetCookie(w, cookie)
	}

	if b, siteConf := GetSiteConf(w, p); b {
		orderNo := r.URL.Query().Get("order_no")
		order, err := goclient.Partner.GetOrderByNo(p.Id, p.Secret, orderNo)
		if err != nil {
			this.Context.Log().PrintErr(err)
			this.OrderEmpty(w, r, p, mm, siteConf)
			return
		}

		if b, siteConf := GetSiteConf(w, p); b {
			this.Context.Template().Execute(w, func(m *map[string]interface{}) {
				(*m)["partner"] = p
				(*m)["title"] = "订单成功-" + p.Name
				(*m)["member"] = mm
				(*m)["conf"] = siteConf
				(*m)["order"] = order
			},
				"views/web/www/order_finish.html",
				"views/web/www/inc/header.html",
				"views/web/www/inc/footer.html")
		}
	}
}

func (this *shoppingC) SubmitOrder_post(w http.ResponseWriter, r *http.Request,
	p *partner.ValuePartner, mm *member.ValueMember) {
	if p == nil || mm == nil {
		w.Write([]byte(`{"result":false,"tag":"101"}`)) //未登录
		return
	}
	r.ParseForm()
	ck, err := r.Cookie("cart")
	if err != nil || ck.Value == "" {
		w.Write([]byte(`{"result":false,"tag":"102"}`)) //购物车为空
		return
	}
	cart := domain.CartCookieFmt(ck.Value)
	deliverAddrId, _ := strconv.Atoi(r.FormValue("AddrId"))
	couponCode := r.FormValue("CouponCode")
	order_no, err := goclient.Partner.SubmitOrder(p.Id, p.Secret, mm.Id,
		0, enum.PAY_OFFLINE, deliverAddrId, cart, couponCode, r.FormValue("note"))
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"result":false,"tag":"109","message":"%s"}`, err.Error())))
		return
	}

	w.Write([]byte(`{"result":true,"data":"` + order_no + `"}`))
}

// 购物车
func (this *shoppingC) CartApi(w http.ResponseWriter, r *http.Request,
	p *partner.ValuePartner, m *member.ValueMember) {
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
		this.Cart_AddItem(w, r, p, memberId, cartKey)
	case "remove":
		this.Cart_RemoveItem(w, r, p, memberId, cartKey)
	}
}

func (this *shoppingC) Cart_GetCart(w http.ResponseWriter,
	p *partner.ValuePartner, memberId int, cartKey string) {
	cart := goclient.Partner.GetShoppingCart(p.Id, memberId, cartKey)
	d, _ := json.Marshal(cart)
	w.Write(d)
}

func (this *shoppingC) Cart_AddItem(w http.ResponseWriter, r *http.Request,
	p *partner.ValuePartner, memberId int, cartKey string) {
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

func (this *shoppingC) Cart_RemoveItem(w http.ResponseWriter, r *http.Request,
	p *partner.ValuePartner, memberId int, cartKey string) {
	goodsId, _ := strconv.Atoi(r.FormValue("id"))
	num, _ := strconv.Atoi(r.FormValue("num"))
	err := goclient.Partner.SubCartItem(p.Id, memberId, cartKey, goodsId, num)
	if err != nil{
		w.Write([]byte(`{error:'`+ err.Error()+`'}`))
	}else{
		w.Write([]byte("{}"))
	}
}



func (this *shoppingC) Cart(w http.ResponseWriter, r *http.Request,
	p *partner.ValuePartner) {
	//
	w.Write([]byte("购物车"))
}

// 购买中转
func (this *shoppingC) BuyRedirect(w http.ResponseWriter, r *http.Request,
	p *partner.ValuePartner, mm *member.ValueMember) {
	if mm == nil {
		RedirectLoginPage(w, r.RequestURI)
	} else {
		w.Header().Add("Location", "/buy/ship")
		w.WriteHeader(302)
	}
}
