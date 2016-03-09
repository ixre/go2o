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
<<<<<<< HEAD
	"go2o/src/core/domain/interface/ad"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
)

// 广告控制器
type adC struct {
}

//广告列表
func (this *adC) List(ctx *echox.Context) error {
	return ctx.RenderOK("ad.list.html", echox.NewRenderData())
}

// 修改广告
func (this *adC) Edit(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	id, _ := strconv.Atoi(ctx.Query("id"))
	e := dps.AdvertisementService.GetAdvertisement(partnerId, id)

	js, _ := json.Marshal(e)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("ad.edit.html", d)
}

// 保存广告
func (this *adC) Create(ctx *echox.Context) error {
	e := ad.ValueAdvertisement{
		Enabled: 1,
	}
	js, _ := json.Marshal(e)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("ad.edit.html", d)
}

// 删除广告(POST)
func (this *adC) Delete_ad(ctx *echox.Context) error {
	req := ctx.Request()
	if req.Method == "POST" {
		req.ParseForm()
		var result gof.Message
		partnerId := getPartnerId(ctx)
		adId, _ := strconv.Atoi(req.FormValue("id"))
		err := dps.AdvertisementService.DelAdvertisement(partnerId, adId)

		if err != nil {
			result.Message = err.Error()
		} else {
			result.Result = true
		}

		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 保存广告(POST)
func (this *adC) SaveAd(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	req := ctx.Request()
	if req.Method == "POST" {
		req.ParseForm()

		var result gof.Message

		e := ad.ValueAdvertisement{}
		web.ParseFormToEntity(req.Form, &e)

		//更新
		e.PartnerId = partnerId

		id, err := dps.AdvertisementService.SaveAdvertisement(partnerId, &e)

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

func (this *adC) Ad_data1(ctx *echox.Context) error {
	return ctx.String(http.StatusOK, `<span style="color:red">暂只支持轮播广告</span>`)
}

func (this *adC) Ad_data2(ctx *echox.Context) error {
	return ctx.String(http.StatusOK, `<span style="color:red">暂只支持轮播广告</span>`)
}

//轮播广告
func (this *adC) Ad_data3(ctx *echox.Context) error {
	d := echox.NewRenderData()
	d.Map["adId"] = ctx.Query("id")
	return ctx.RenderOK("ad.data3.html", d)
}

// 创建广告图片
func (this *adC) CreateAdImage(ctx *echox.Context) error {
	adId, _ := strconv.Atoi(ctx.Query("ad_id"))
=======
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	e := ad.ValueImage{
		Enabled:         1,
		AdvertisementId: adId,
		LinkUrl:         "http://",
		ImageUrl:        ctx.App.Config().GetString(variable.NoPicPath),
	}

	js, _ := json.Marshal(e)
<<<<<<< HEAD
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("ad.image.html", d)
}

// 保存广告
func (this *adC) EditAdImage(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	adId, _ := strconv.Atoi(ctx.Query("ad_id"))
	imgId, _ := strconv.Atoi(ctx.Query("id"))
=======

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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d

	e := dps.AdvertisementService.GetValueAdImage(partnerId, adId, imgId)

	js, _ := json.Marshal(e)
<<<<<<< HEAD
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("ad.image.html", d)
}

// 保存图片(POST)
func (this *adC) SaveImage(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	if r.Method == "POST" {
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
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 删除广告图片(POST)
func (this *adC) Delete_image(ctx *echox.Context) error {
	req := ctx.Request()
	if req.Method == "POST" {
		req.ParseForm()
		var result gof.Message
		partnerId := getPartnerId(ctx)
		adId, _ := strconv.Atoi(req.FormValue("ad_id"))
		imgId, _ := strconv.Atoi(req.FormValue("id"))
		err := dps.AdvertisementService.DelAdImage(partnerId, adId, imgId)

		if err != nil {
			result.Message = err.Error()
		} else {
			result.Result = true
		}

		return ctx.JSON(http.StatusOK, result)
	}
	return nil
=======

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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
