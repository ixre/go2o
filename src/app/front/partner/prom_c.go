/**
 * Copyright 2014 @ ops.
 * name :
 * author : jarryliu
 * date : 2013-12-04 08:21
 * description :
 * history :
 */
package partner

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/service/dps"
	"html/template"
	"strconv"
	"strings"
	"time"
)

var _ mvc.Filter = new(promC)

type promC struct {
	Base *baseC
}

func (this *promC) Requesting(ctx *web.Context) bool {
	return this.Base.Requesting(ctx)
}
func (this *promC) RequestEnd(ctx *web.Context) {
	this.Base.RequestEnd(ctx)
}

func (this *promC) CreateCoupon(ctx *web.Context) {
	//partnerId := this.Base.GetPartnerId(ctx)
	ctx.App.Template().Execute(ctx.ResponseWriter,
		func(m *map[string]interface{}) {
		},
		"views/partner/promotion/create_coupon.html")
}

func (this *promC) EditCoupon(ctx *web.Context) {
	partnerId := this.Base.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	e := dps.PromService.GetCoupon(partnerId, id).GetValue()
	js, _ := json.Marshal(e)

	ctx.App.Template().Execute(w,
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(js)
		},
		"views/partner/promotion/edit_coupon.html")
}

//　绑定优惠券操作页
func (this *promC) BindCoupon(ctx *web.Context) {
	partnerId := this.Base.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
	id, _ := strconv.Atoi(r.URL.Query().Get("coupon_id"))
	e := dps.PromService.GetCoupon(partnerId, id).GetValue()
	ctx.App.Template().Execute(w,
		func(m *map[string]interface{}) {
			(*m)["entity"] = e
		},
		"views/partner/promotion/bind_coupon.html")
}

func (this *promC) BindCoupon_post(ctx *web.Context) {
	partnerId := this.Base.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
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

func (this *promC) SaveCoupon_post(ctx *web.Context) {
	partnerId := this.Base.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter

	var result gof.JsonResult
	r.ParseForm()
	var e promotion.ValueCoupon
	web.ParseFormToEntity(r.Form, &e)

	const layout string = "2006-01-02 15:04:05"
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

func (this *promC) Coupon(ctx *web.Context) {
	//partnerId := this.Base.GetPartnerId(ctx)
	ctx.App.Template().Execute(ctx.ResponseWriter,
		func(m *map[string]interface{}) {

		}, "views/partner/promotion/coupon_list.html")
}
