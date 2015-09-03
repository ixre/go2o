/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"encoding/json"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/mvc"
	"go2o/src/core/domain/interface/ad"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"html/template"
	"strconv"
)

var _ mvc.Filter = new(adC)

// 广告控制器
type adC struct {
	*baseC
}

//广告列表
func (this *adC) List(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.Response, gof.TemplateDataMap{}, "views/partner/ad/ad_list.html")
}

// 修改广告
func (this *adC) Edit(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	form := ctx.Request.URL.Query()
	id, _ := strconv.Atoi(form.Get("id"))
	e := dps.AdvertisementService.GetAdvertisement(partnerId, id)

	js, _ := json.Marshal(e)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity": template.JS(js),
		},
		"views/partner/ad/ad_edit.html")
}

// 保存广告
func (this *adC) Create(ctx *web.Context) {
	e := ad.ValueAdvertisement{
		Enabled: 1,
	}

	js, _ := json.Marshal(e)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity": template.JS(js),
		},
		"views/partner/ad/ad_edit.html")
}

// 删除广告
func (this *adC) Del_post(ctx *web.Context) {
	ctx.Request.ParseForm()
	form := ctx.Request.Form
	var result gof.Message
	partnerId := this.GetPartnerId(ctx)
	adId, _ := strconv.Atoi(form.Get("id"))
	err := dps.AdvertisementService.DelAdvertisement(partnerId, adId)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}

	ctx.Response.JsonOutput(result)
}

func (this *adC) SaveAd_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r := ctx.Request
	r.ParseForm()

	var result gof.Message

	e := ad.ValueAdvertisement{}
	web.ParseFormToEntity(r.Form, &e)

	//更新
	e.PartnerId = partnerId

	id, err := dps.AdvertisementService.SaveAdvertisement(partnerId, &e)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
		result.Data = id
	}
	ctx.Response.JsonOutput(result)
}

func (this *adC) Ad_data1(ctx *web.Context) {
	ctx.Response.Write([]byte(`<span style="color:red">暂只支持轮播广告</span>`))
}

func (this *adC) Ad_data2(ctx *web.Context) {
	ctx.Response.Write([]byte(`<span style="color:red">暂只支持轮播广告</span>`))
}

//轮播广告
func (this *adC) Ad_data3(ctx *web.Context) {
	dps.AdvertisementService.GetAdvertisement(this.GetPartnerId(ctx), 0)
	ctx.App.Template().Execute(ctx.Response, gof.TemplateDataMap{
		"adId": ctx.Request.URL.Query().Get("id"),
	}, "views/partner/ad/ad_data3.html")
}

// 创建广告图片
func (this *adC) CreateAdImage(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	adId, _ := strconv.Atoi(form.Get("ad_id"))
	e := ad.ValueImage{
		Enabled:         1,
		AdvertisementId: adId,
		LinkUrl:         "http://",
		ImageUrl:        ctx.App.Config().GetString(variable.NoPicPath),
	}

	js, _ := json.Marshal(e)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity": template.JS(js),
		},
		"views/partner/ad/ad_image.html")
}

// 保存广告
func (this *adC) EditAdImage(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	partnerId := this.GetPartnerId(ctx)
	adId, _ := strconv.Atoi(form.Get("ad_id"))
	imgId, _ := strconv.Atoi(form.Get("id"))

	e := dps.AdvertisementService.GetValueAdImage(partnerId, adId, imgId)

	js, _ := json.Marshal(e)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity": template.JS(js),
		},
		"views/partner/ad/ad_image.html")
}

func (this *adC) SaveImage_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r := ctx.Request
	r.ParseForm()

	var result gof.Message

	e := ad.ValueImage{}
	web.ParseFormToEntity(r.Form, &e)

	id, err := dps.AdvertisementService.SaveImage(partnerId, e.AdvertisementId, &e)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
		result.Data = id
	}
	ctx.Response.JsonOutput(result)
}

func (this *adC) Del_image_post(ctx *web.Context) {
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
