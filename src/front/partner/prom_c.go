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
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"html/template"
	"strconv"
	"strings"
	"time"
	"go2o/src/x/echox"
	"net/http"
)


type promC struct {
}

func (this *promC) List(ctx *echox.Context) error {
	flag, _ := strconv.Atoi(ctx.Query("flag"))
	d := ctx.NewData()
	d.Map["flag"] = flag
	return ctx.RenderOK(fmt.Sprintf("promotion/p%d_list.html", flag),d)
}

// 删除促销
func (this *promC) Del_post(ctx *echox.Context) error {
	ctx.Request.ParseForm()
	form := ctx.Request.Form
	var result gof.Message
	partnerId := getPartnerId(ctx)
	promId, _ := strconv.Atoi(form.Get("id"))

	err := dps.PromService.DelPromotion(partnerId, promId)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	return ctx.JSON(http.StatusOK, result)
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

	return ctx.RenderOK("promotion/cash_back.html",d)
}

func (this *promC) Edit_cb(ctx *echox.Context) error {
	form := ctx.Request.URL.Query()
	id, _ := strconv.Atoi(form.Get("id"))
	e, e2 := dps.PromService.GetPromotion(id)

	js, _ := json.Marshal(e)
	js2, _ := json.Marshal(e2)

	var goodsInfo string
	goods := dps.SaleService.GetValueGoods(getPartnerId(ctx), e.GoodsId)
	goodsInfo = fmt.Sprintf("%s<span>(销售价：%s)</span>", goods.Name, format.FormatFloat(goods.SalePrice))

	d:= ctx.NewData()

	d.Map = gof.TemplateDataMap{
			"entity":     template.JS(js),
			"entity2":    template.JS(js2),
			"goods_info": template.HTML(goodsInfo),
			"goods_cls":  "",
		}

	return ctx.RenderOK("promotion/cash_back.html",d)
}

// 保存现金返现
func (this *promC) Save_cb_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
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
	return ctx.JSON(http.StatusOK, result)
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
	d.Map =  gof.TemplateDataMap{
			"entity":  template.JS(js),
			"entity2": template.JS(js2),
			"levelDr": template.HTML(levelDr),
		}

	return ctx.RenderOK("promotion/coupon.html",d)
}

func (this *promC) Edit_coupon(ctx *echox.Context) error {
	form := ctx.Request.URL.Query()
	id, _ := strconv.Atoi(form.Get("id"))
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

	return ctx.RenderOK("promotion/coupon.html",d)
}

// 保存优惠券
func (this *promC) Save_coupon_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
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
	return ctx.JSON(http.StatusOK, result)
}

//　绑定优惠券操作页
func (this *promC) Bind_coupon(ctx *echox.Context) error {
	r, w := ctx.Request, ctx.Response
	id, _ := strconv.Atoi(r.URL.Query().Get("coupon_id"))
	e, e2 := dps.PromService.GetPromotion(id)
	if e.PartnerId != getPartnerId(ctx) {
		return ctx.StringOK(promotion.ErrNoSuchPromotion.Error())

	}

	d := ctx.NewData()
	d.Map =  gof.TemplateDataMap{
			"entity":  e,
			"entity2": e2,
		}

	return ctx.RenderOK("promotion/bind_coupon.html",d)
}

func (this *promC) Bind_coupon_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
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
		result.Message = err.Error()
	} else {
		result.Result = true
	}

	return ctx.JSON(http.StatusOK,result)
}
