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
	"com/domain/interface/member"
	"com/ording"
	"com/ording/entity"
	"com/service/goclient"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"github.com/newmin/gof/app"
	"github.com/newmin/gof/web"
)

type shoppingC struct {
	app.Context
}

func (this *shoppingC) GetDeliverAddrs(w http.ResponseWriter, r *http.Request,
	p *entity.Partner, m *member.ValueMember) {
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
	r *http.Request, p *entity.Partner, m *member.ValueMember) {
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
	p *entity.Partner, m *member.ValueMember) {

	r.ParseForm()
	var message string = "购物车还是空的!"
	code := r.FormValue("code")
	ck, err := r.Cookie("cart")
	if err == nil {
		cartData := ording.CartCookieFmt(ck.Value)
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
