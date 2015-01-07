/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-12 16:55
 * description :
 * history :
 */

package partner

import (
	"com/domain/interface/enum"
	"com/ording/cache"
	"com/ording/dproxy"
	"encoding/json"
	"html/template"
	"net/http"
	"github.com/newmin/gof/app"
	"strings"
)

type orderC struct {
	app.Context
}

func (this *orderC) List(w http.ResponseWriter, r *http.Request, partnerId int) {

	shopsJson := cache.GetShopsJson(partnerId)
	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {
			(*m)["shops"] = template.JS(shopsJson)
		}, "views/partner/order/order_list.html")
}

func (this *orderC) Cancel(w http.ResponseWriter, r *http.Request, partnerId int) {
	this.Context.Template().Execute(w, nil, "views/partner/order/cancel.html")

}

func (this *orderC) Cancel_post(w http.ResponseWriter, r *http.Request, partnerId int) {
	r.ParseForm()
	reason := r.FormValue("reason")
	err := dproxy.ShoppingService.CancelOrder(partnerId,
		r.FormValue("order_no"), reason)

	if err == nil {
		w.Write([]byte("{result:true}"))
	} else {
		w.Write([]byte(`{result:false,message:"` + err.Error() + `"}`))
	}
}

func (this *orderC) View(w http.ResponseWriter, r *http.Request, partnerId int) {
	r.ParseForm()
	e := dproxy.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		w.Write([]byte("无效订单"))
		return
	}

	member := dproxy.MemberService.GetMember(e.MemberId)
	if member == nil {
		w.Write([]byte("无效订单"))
		return
	}

	e.ItemsInfo = strings.Replace(e.ItemsInfo, "\n", "<br />", -1)
	if len(e.Note) == 0 {
		e.Note = "无备注"
	}

	js, _ := json.Marshal(e)

	var shopName string
	var payment string
	var orderStateText string
	if e.ShopId == 0 {
		shopName = "未指定"
	} else {
		shopName = dproxy.PartnerService.GetShopValueById(partnerId, e.ShopId).Name
	}
	payment = enum.GetPaymentName(e.PayMethod)
	orderStateText = enum.OrderState(e.Status).String()

	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(js)
			(*m)["member"] = member
			(*m)["shopName"] = shopName
			(*m)["payment"] = payment
			(*m)["state"] = orderStateText

		}, "views/partner/order/order_view.html")
}

func (this *orderC) Setup(w http.ResponseWriter, r *http.Request, partnerId int) {
	r.ParseForm()
	e := dproxy.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		w.Write([]byte("无效订单"))
		return
	}

	if e.ShopId == 0 {
		this.setShop(w, r, partnerId, e)
	} else {
		this.setState(w, r, partnerId, e)
	}
}

func (this *orderC) OrderSetup_post(w http.ResponseWriter, r *http.Request, partnerId int) {
	r.ParseForm()
	err := dproxy.ShoppingService.HandleOrder(partnerId, r.FormValue("order_no"))
	if err != nil {
		w.Write([]byte("{result:false,message:'" + err.Error() + "'}"))
	} else {
		w.Write([]byte("{result:true}"))
	}
}
