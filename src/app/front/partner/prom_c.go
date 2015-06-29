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

func (this *promC) List(ctx *web.Context) {
	var flag int
	flag, _ = strconv.Atoi(ctx.Request.URL.Query().Get("flag"))

	ctx.App.Template().Execute(ctx.Response, gof.TemplateDataMap{
		"flag": flag,
	}, fmt.Sprintf("views/partner/promotion/p%d_list.html", flag))
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

// 删除现金返现
func (this *promC) Del_cb_post(ctx *web.Context) {
	ctx.Request.ParseForm()
	form := ctx.Request.Form
	var result gof.Message
	partnerId := this.GetPartnerId(ctx)
	adId, _ := strconv.Atoi(form.Get("ad_id"))
	imgId, _ := strconv.Atoi(form.Get("id"))
	err := dps.AdvertisementService.DelAdImage(partnerId, adId, imgId)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	ctx.Response.JsonOutput(result)
}

func (this *promC) CreateCoupon(ctx *web.Context) {

	levelDr := this.getLevelDropDownList(ctx)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity":  "{}",
			"levelDr": template.HTML(levelDr),
		},
		"views/partner/promotion/create_coupon.html")
}

func (this *promC) EditCoupon(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	e := dps.PromService.GetCoupon(partnerId, id).GetValue()
	js, _ := json.Marshal(e)

	levelDr := this.getLevelDropDownList(ctx)

	ctx.App.Template().Execute(w,
		gof.TemplateDataMap{
			"entity":  template.JS(js),
			"levelDr": template.HTML(levelDr),
		},
		"views/partner/promotion/edit_coupon.html")
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

//　绑定优惠券操作页
func (this *promC) BindCoupon(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	id, _ := strconv.Atoi(r.URL.Query().Get("coupon_id"))
	e := dps.PromService.GetCoupon(partnerId, id).GetValue()
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

func (this *promC) SaveCoupon_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response

	var result gof.Message
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
		result = gof.Message{Result: false, Message: err.Error()}
	} else {
		result = gof.Message{Result: true, Message: ""}
	}
	w.Write(result.Marshal())
}

func (this *promC) Coupon(ctx *web.Context) {
	//partnerId := this.GetPartnerId(ctx)
	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{}, "views/partner/promotion/coupon_list.html")
}
