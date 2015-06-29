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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"html/template"
	"strconv"
	"strings"
	"time"
)

var _ mvc.Filter = new(promC)

type promC struct {
	*baseC
}


func (this *promC) getLevelDropDownList(ctx *web.Context) string {
	buf := bytes.NewBufferString("")
	lvs := dps.PartnerService.GetMemberLevels(this.GetPartnerId(ctx))
	for _, v := range lvs {
		if v.Enabled == 1 {
			buf.WriteString(fmt.Sprintf(`<option value="%d">%s</option>`, v.Value, v.Name))
		}
	}
	return buf.String()
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
	e := &promotion.ValuePromotion{
		Enabled: 1,
	}
	e2 := &promotion.ValueCashBack{
		BackType: 1,
	}
	js, _ := json.Marshal(e)
	js2, _ := json.Marshal(e2)

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
	e, e2 := dps.PromService.GetPromotion(id)

	js, _ := json.Marshal(e)
	js2, _ := json.Marshal(e2)

	var goodsInfo string
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
	e := &promotion.ValuePromotion{
		Enabled: 1,
	}
	e2 := &promotion.ValueCoupon{
		AllowEnable: 1,
	}

	js, _ := json.Marshal(e)
	js2, _ := json.Marshal(e2)

	levelDr := this.getLevelDropDownList(ctx)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity":     template.JS(js),
			"entity2":    template.JS(js2),
			"levelDr": template.HTML(levelDr),
		},
		"views/partner/promotion/coupon.html")
}

func (this *promC) Edit_coupon(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	id, _ := strconv.Atoi(form.Get("id"))
	e, e2 := dps.PromService.GetPromotion(id)

	if e.PartnerId != this.GetPartnerId(ctx){
		this.ErrorOutput(ctx,promotion.ErrNoSuchPromotion.Error())
		return
	}

	js, _ := json.Marshal(e)
	js2, _ := json.Marshal(e2)

	levelDr := this.getLevelDropDownList(ctx)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity":     template.JS(js),
			"entity2":    template.JS(js2),
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

/************ NORMAL *******************/
//　绑定优惠券操作页
func (this *promC) BindCoupon(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	id, _ := strconv.Atoi(r.URL.Query().Get("coupon_id"))
	e,_ := dps.PromService.GetPromotion(id)
	if e.PartnerId != this.GetPartnerId(ctx){
		this.ErrorOutput(ctx,promotion.ErrNoSuchPromotion.Error())
		return
	}
	ctx.App.Template().Execute(w,
		gof.TemplateDataMap{
			"entity": e,
		},
		"views/partner/promotion/bind_coupon.html")
}

func (this *promC) BindCoupon_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
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
		result.Result = false
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	w.Write(result.Marshal())
}
