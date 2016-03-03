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
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
<<<<<<< HEAD
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
=======
	"github.com/jsix/gof/web/mvc"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"html/template"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"strconv"
	"strings"
	"time"
)

<<<<<<< HEAD
type promC struct {
}

func (this *promC) List(ctx *echox.Context) error {
	flag, _ := strconv.Atoi(ctx.Query("flag"))
	d := ctx.NewData()
	d.Map["flag"] = flag
	return ctx.RenderOK(fmt.Sprintf("prom.p%d_list.html", flag), d)
}

// 删除促销(POST)
func (this *promC) Del(ctx *echox.Context) error {
	req := ctx.Request()
	if req.Method == "POST" {
		req.ParseForm()
		var result gof.Message
		partnerId := getPartnerId(ctx)
		promId, _ := strconv.Atoi(req.FormValue("id"))

		err := dps.PromService.DelPromotion(partnerId, promId)

		if err != nil {
			result.Message = err.Error()
		} else {
			result.Result = true
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 创建返现促销
func (this *promC) Create_cb(ctx *echox.Context) error {
=======
var _ mvc.Filter = new(promC)

type promC struct {
	*baseC
}

func (this *promC) List(ctx *web.Context) {
	var flag int
	flag, _ = strconv.Atoi(ctx.Request.URL.Query().Get("flag"))

	ctx.App.Template().Execute(ctx.Response, gof.TemplateDataMap{
		"flag": flag,
	}, fmt.Sprintf("views/partner/promotion/p%d_list.html", flag))
}

// 删除促销
func (this *promC) Del_post(ctx *web.Context) {
	ctx.Request.ParseForm()
	form := ctx.Request.Form
	var result gof.Message
	partnerId := this.GetPartnerId(ctx)
	promId, _ := strconv.Atoi(form.Get("id"))

	err := dps.PromService.DelPromotion(partnerId, promId)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	ctx.Response.JsonOutput(result)
}

// 创建返现促销
func (this *promC) Create_cb(ctx *web.Context) {
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	e := &promotion.ValuePromotion{
		Enabled: 1,
	}
	e2 := &promotion.ValueCashBack{
		BackType: 1,
	}
	js, _ := json.Marshal(e)
	js2, _ := json.Marshal(e2)

<<<<<<< HEAD
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"entity":    template.JS(js),
		"entity2":   template.JS(js2),
		"goods_cls": "hidden",
	}

	return ctx.RenderOK("prom.cash_back.html", d)
}

func (this *promC) Edit_cb(ctx *echox.Context) error {
	id, _ := strconv.Atoi(ctx.Query("id"))
=======
	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity":    template.JS(js),
			"entity2":   template.JS(js2),
			"goods_cls": "hidden",
		},
		"views/partner/promotion/cash_back.html")
}

func (this *promC) Edit_cb(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	id, _ := strconv.Atoi(form.Get("id"))
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	e, e2 := dps.PromService.GetPromotion(id)

	js, _ := json.Marshal(e)
	js2, _ := json.Marshal(e2)

	var goodsInfo string
<<<<<<< HEAD
	goods := dps.SaleService.GetValueGoods(getPartnerId(ctx), e.GoodsId)
	goodsInfo = fmt.Sprintf("%s<span>(销售价：%s)</span>", goods.Name, format.FormatFloat(goods.SalePrice))

	d := ctx.NewData()

	d.Map = gof.TemplateDataMap{
		"entity":     template.JS(js),
		"entity2":    template.JS(js2),
		"goods_info": template.HTML(goodsInfo),
		"goods_cls":  "",
	}

	return ctx.RenderOK("prom.cash_back.html", d)
}

// 保存现金返现(POST)
func (this *promC) Save_cb(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	if r.Method == "POST" {
		r.ParseForm()

		var result gof.Message

		e := promotion.ValuePromotion{}
		web.ParseFormToEntity(r.Form, &e)
		e2 := promotion.ValueCashBack{}
		web.ParseFormToEntity(r.Form, &e2)

		e.PartnerId = partnerId
		e.TypeFlag = promotion.TypeFlagCashBack

		id, err := dps.PromService.SaveCashBackPromotion(partnerId, &e, &e2)

		if err != nil {
			result.Message = err.Error()
		} else {
			result.Result = true
			result.Data = id
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 创建优惠券
func (this *promC) Create_coupon(ctx *echox.Context) error {
=======
	goods := dps.SaleService.GetValueGoods(this.GetPartnerId(ctx), e.GoodsId)
	goodsInfo = fmt.Sprintf("%s<span>(销售价：%s)</span>", goods.Name, format.FormatFloat(goods.SalePrice))

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity":     template.JS(js),
			"entity2":    template.JS(js2),
			"goods_info": template.HTML(goodsInfo),
			"goods_cls":  "",
		},
		"views/partner/promotion/cash_back.html")
}

// 保存现金返现
func (this *promC) Save_cb_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r := ctx.Request
	r.ParseForm()

	var result gof.Message

	e := promotion.ValuePromotion{}
	web.ParseFormToEntity(r.Form, &e)
	e2 := promotion.ValueCashBack{}
	web.ParseFormToEntity(r.Form, &e2)

	e.PartnerId = partnerId
	e.TypeFlag = promotion.TypeFlagCashBack

	id, err := dps.PromService.SaveCashBackPromotion(partnerId, &e, &e2)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
		result.Data = id
	}
	ctx.Response.JsonOutput(result)
}

// 创建优惠券
func (this *promC) Create_coupon(ctx *web.Context) {
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	e := &promotion.ValuePromotion{
		Enabled: 1,
	}
	e2 := &promotion.ValueCoupon{
		BeginTime: time.Now().Unix(),
		OverTime:  time.Now().Add(time.Hour * 24 * 30).Unix(),
		Discount:  100,
	}

	js, _ := json.Marshal(e)
	js2, _ := json.Marshal(e2)

<<<<<<< HEAD
	levelDr := getLevelDropDownList(getPartnerId(ctx))

	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"entity":  template.JS(js),
		"entity2": template.JS(js2),
		"levelDr": template.HTML(levelDr),
	}

	return ctx.RenderOK("prom.coupon.html", d)
}

func (this *promC) Edit_coupon(ctx *echox.Context) error {
	id, _ := strconv.Atoi(ctx.Query("id"))
	e, e2 := dps.PromService.GetPromotion(id)

	if e.PartnerId != getPartnerId(ctx) {
		return ctx.StringOK(promotion.ErrNoSuchPromotion.Error())
=======
	levelDr := getLevelDropDownList(this.GetPartnerId(ctx))

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity":  template.JS(js),
			"entity2": template.JS(js2),
			"levelDr": template.HTML(levelDr),
		},
		"views/partner/promotion/coupon.html")
}

func (this *promC) Edit_coupon(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	id, _ := strconv.Atoi(form.Get("id"))
	e, e2 := dps.PromService.GetPromotion(id)

	if e.PartnerId != this.GetPartnerId(ctx) {
		this.ErrorOutput(ctx, promotion.ErrNoSuchPromotion.Error())
		return
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	}

	js, _ := json.Marshal(e)
	js2, _ := json.Marshal(e2)

<<<<<<< HEAD
	levelDr := getLevelDropDownList(getPartnerId(ctx))

	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"entity":  template.JS(js),
		"entity2": template.JS(js2),
		"levelDr": template.HTML(levelDr),
	}

	return ctx.RenderOK("prom.coupon.html", d)
}

// 保存优惠券(POST)
func (this *promC) Save_coupon(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	if r.Method == "POST" {
		r.ParseForm()

		var result gof.Message

		e := promotion.ValuePromotion{}
		web.ParseFormToEntity(r.Form, &e)
		e2 := promotion.ValueCoupon{}
		web.ParseFormToEntity(r.Form, &e2)

		e.PartnerId = partnerId
		e.TypeFlag = promotion.TypeFlagCoupon

		const layout string = "2006-01-02 15:04:05"
		bt, _ := time.Parse(layout, r.FormValue("BeginTime"))
		ot, _ := time.Parse(layout, r.FormValue("OverTime"))
		e2.BeginTime = bt.Unix()
		e2.OverTime = ot.Unix()

		id, err := dps.PromService.SaveCoupon(partnerId, &e, &e2)

		if err != nil {
			result.Message = err.Error()
		} else {
			result.Result = true
			result.Data = id
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

//　绑定优惠券操作页
func (this *promC) Bind_coupon(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.bind_coupon_post(ctx)
	}
	id, _ := strconv.Atoi(ctx.Query("coupon_id"))
	e, e2 := dps.PromService.GetPromotion(id)
	if e.PartnerId != getPartnerId(ctx) {
		return ctx.StringOK(promotion.ErrNoSuchPromotion.Error())

	}

	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"entity":  e,
		"entity2": e2,
	}

	return ctx.RenderOK("prom.bind_coupon.html", d)
}

func (this *promC) bind_coupon_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
=======
	levelDr := getLevelDropDownList(this.GetPartnerId(ctx))

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity":  template.JS(js),
			"entity2": template.JS(js2),
			"levelDr": template.HTML(levelDr),
		},
		"views/partner/promotion/coupon.html")
}

// 保存优惠券
func (this *promC) Save_coupon_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r := ctx.Request
	r.ParseForm()

	var result gof.Message

	e := promotion.ValuePromotion{}
	web.ParseFormToEntity(r.Form, &e)
	e2 := promotion.ValueCoupon{}
	web.ParseFormToEntity(r.Form, &e2)

	e.PartnerId = partnerId
	e.TypeFlag = promotion.TypeFlagCoupon

	const layout string = "2006-01-02 15:04:05"
	bt, _ := time.Parse(layout, r.FormValue("BeginTime"))
	ot, _ := time.Parse(layout, r.FormValue("OverTime"))
	e2.BeginTime = bt.Unix()
	e2.OverTime = ot.Unix()

	id, err := dps.PromService.SaveCoupon(partnerId, &e, &e2)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
		result.Data = id
	}
	ctx.Response.JsonOutput(result)
}

//　绑定优惠券操作页
func (this *promC) Bind_coupon(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	id, _ := strconv.Atoi(r.URL.Query().Get("coupon_id"))
	e, e2 := dps.PromService.GetPromotion(id)
	if e.PartnerId != this.GetPartnerId(ctx) {
		this.ErrorOutput(ctx, promotion.ErrNoSuchPromotion.Error())
		return
	}
	ctx.App.Template().Execute(w,
		gof.TemplateDataMap{
			"entity":  e,
			"entity2": e2,
		},
		"views/partner/promotion/bind_coupon.html")
}

func (this *promC) Bind_coupon_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	var result gof.Message
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
<<<<<<< HEAD
=======
		result.Result = false
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		result.Message = err.Error()
	} else {
		result.Result = true
	}
<<<<<<< HEAD

	return ctx.JSON(http.StatusOK, result)
=======
	w.Write(result.Marshal())
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
