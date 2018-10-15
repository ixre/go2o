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
	//"github.com/jsix/gof/web"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/jsix/gof/web/form"
)

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
	req := ctx.HttpRequest()
	if req.Method == "POST" {
		req.ParseForm()
		var result gof.Result
		partnerId := getPartnerId(ctx)
		promId, _ := strconv.Atoi(req.FormValue("id"))

		err := dps.PromService.DelPromotion(partnerId, promId)

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 创建返现促销
func (this *promC) Create_cb(ctx *echox.Context) error {
	e := &promotion.ValuePromotion{
		Enabled: 1,
	}
	e2 := &promotion.ValueCashBack{
		BackType: 1,
	}
	js, _ := json.Marshal(e)
	js2, _ := json.Marshal(e2)

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
	e, e2 := dps.PromService.GetPromotion(id)

	js, _ := json.Marshal(e)
	js2, _ := json.Marshal(e2)

	var goodsInfo string
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
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		r.ParseForm()

		var result gof.Result

		e := promotion.ValuePromotion{}
		form.ParseEntity(r.Form, &e)
		e2 := promotion.ValueCashBack{}
		form.ParseEntity(r.Form, &e2)

		e.PartnerId = partnerId
		e.TypeFlag = promotion.TypeFlagCashBack

		id, err := dps.PromService.SaveCashBackPromotion(partnerId, &e, &e2)

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
			var data = make(map[string]string)
			data["id"] = fmt.Sprintf("%d", id)
			result.Data = data
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 创建优惠券
func (this *promC) Create_coupon(ctx *echox.Context) error {
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
	}

	js, _ := json.Marshal(e)
	js2, _ := json.Marshal(e2)

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
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		r.ParseForm()

		var result gof.Result

		e := promotion.ValuePromotion{}
		form.ParseEntity(r.Form, &e)
		e2 := promotion.ValueCoupon{}
		form.ParseEntity(r.Form, &e2)

		e.PartnerId = partnerId
		e.TypeFlag = promotion.TypeFlagCoupon

		const layout string = "2006-01-02 15:04:05"
		bt, _ := time.Parse(layout, r.FormValue("BeginTime"))
		ot, _ := time.Parse(layout, r.FormValue("OverTime"))
		e2.BeginTime = bt.Unix()
		e2.OverTime = ot.Unix()

		id, err := dps.PromService.SaveCoupon(partnerId, &e, &e2)

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
			var data = make(map[string]string)
			data["id"] = fmt.Sprintf("%d", id)
			result.Data = data
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
	r := ctx.HttpRequest()
	var result gof.Result
	r.ParseForm()
	id, err := strconv.Atoi(r.FormValue("id"))
	if err == nil {
		memberIds := strings.TrimSpace(r.FormValue("member_ids"))
		if memberIds == "" {
			result.ErrMsg = "请选择会员"
		} else {
			idArr := strings.Split(memberIds, ",")
			err = dps.PromService.BindCoupons(partnerId, id, idArr)
		}
	}
	if err != nil {
		result.ErrMsg = err.Error()
	} else {
		result.ErrCode = 0
	}

	return ctx.JSON(http.StatusOK, result)
}
