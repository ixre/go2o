/**
 * Copyright 2014 @ ops.
 * name :
 * author : newmin
 * date : 2013-12-04 08:21
 * description :
 * history :
 */
package partner

import (
	"com/domain/interface/promotion"
	"com/ording/dproxy"
	"encoding/json"
	"html/template"
	"net/http"
	"ops/cf"
	"ops/cf/app"
	"ops/cf/web"
	"strconv"
	"strings"
	"time"
)

type promC struct {
	app.Context
}

func (this *promC) CreateCoupon(w http.ResponseWriter, r *http.Request, partnerId int) {
	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {
		},
		"views/partner/promotion/create_coupon.html")
}

func (this *promC) EditCoupon(w http.ResponseWriter, r *http.Request, partnerId int) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	e := dproxy.PromService.GetCoupon(partnerId, id).GetValue()
	js, _ := json.Marshal(e)

	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(js)
		},
		"views/partner/promotion/edit_coupon.html")
}

//　绑定优惠券操作页
func (this *promC) BindCoupon(w http.ResponseWriter, r *http.Request, partnerId int) {
	id, _ := strconv.Atoi(r.URL.Query().Get("coupon_id"))
	e := dproxy.PromService.GetCoupon(partnerId, id).GetValue()
	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {
			(*m)["entity"] = e
		},
		"views/partner/promotion/bind_coupon.html")
}

func (this *promC) BindCoupon_post(w http.ResponseWriter, r *http.Request, partnerId int) {
	var result cf.JsonResult
	r.ParseForm()
	id, err := strconv.Atoi(r.FormValue("id"))
	if err == nil {
		memberIds := strings.TrimSpace(r.FormValue("member_ids"))
		if memberIds == "" {
			result.Message = "请选择会员"
		} else {
			idArr := strings.Split(memberIds, ",")
			err = dproxy.PromService.BindCoupons(partnerId, id, idArr)
		}
	}
	if err != nil {
		result.Result = false
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	w.Write(result.Marshal())
}

func (this *promC) SaveCoupon_post(w http.ResponseWriter, r *http.Request, partnerId int) {

	var result cf.JsonResult
	r.ParseForm()
	var e promotion.ValueCoupon
	web.ParseFormToEntity(r.Form, &e)
	dt := time.Now()

	if e.Id > 0 {
		o := dproxy.PromService.GetCoupon(partnerId, e.Id).GetValue()
		e.PtId = partnerId
		e.Amount = o.Amount
		e.CreateTime = o.CreateTime
	} else {
		e.PtId = partnerId
		e.CreateTime = dt
	}
	e.UpdateTime = dt
	_, err := dproxy.PromService.SaveCoupon(partnerId, &e)

	if err != nil {
		result = cf.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = cf.JsonResult{Result: true, Message: ""}
	}
	w.Write(result.Marshal())
}

func (this *promC) Coupon(w http.ResponseWriter, r *http.Request, partnerId int) {
	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {

		}, "views/partner/promotion/coupon_list.html")
}
