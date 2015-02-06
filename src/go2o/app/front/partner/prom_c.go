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
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web"
	"go2o/core/domain/interface/promotion"
	"go2o/core/service/dps"
	"html/template"
	"net/http"
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
	e := dps.PromService.GetCoupon(partnerId, id).GetValue()
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
	e := dps.PromService.GetCoupon(partnerId, id).GetValue()
	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {
			(*m)["entity"] = e
		},
		"views/partner/promotion/bind_coupon.html")
}

func (this *promC) BindCoupon_post(w http.ResponseWriter, r *http.Request, partnerId int) {
	var result gof.JsonResult
	r.ParseForm()
	id, err := strconv.Atoi(r.FormValue("id"))
	if err == nil {
		memberIds := strings.TrimSpace(r.FormValue("member_ids"))
		if memberIds == "" {
			result.Message = "请选择会员"
		} else {
			idArr := strings.Split(memberIds, ",")
			err = dps.PromService.BindCoupons(partnerId, id, idArr)
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

	var result gof.JsonResult
	r.ParseForm()
	var e promotion.ValueCoupon
	web.ParseFormToEntity(r.Form, &e)

	const layout string = "2006-01-02 15:04-05"
	bt, _ := time.Parse(layout, r.FormValue("BeginTime"))
	ot, _ := time.Parse(layout, r.FormValue("OverTime"))
	e.BeginTime = bt.Unix()
	e.OverTime = ot.Unix()

	_, err := dps.PromService.SaveCoupon(partnerId, &e)

	if err != nil {
		result = gof.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true, Message: ""}
	}
	w.Write(result.Marshal())
}

func (this *promC) Coupon(w http.ResponseWriter, r *http.Request, partnerId int) {
	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {

		}, "views/partner/promotion/coupon_list.html")
}
