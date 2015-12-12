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
	return ctx.RenderOK("ad/ad_list.html", echox.NewRenderData())
}

// 修改广告
func (this *adC) Edit(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	form := ctx.Request.URL.Query()
	id, _ := strconv.Atoi(form.Get("id"))
	e := dps.AdvertisementService.GetAdvertisement(partnerId, id)

	js, _ := json.Marshal(e)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("ad/ad_edit.html", d)
}

// 保存广告
func (this *adC) Create(ctx *echox.Context) error {
	e := ad.ValueAdvertisement{
		Enabled: 1,
	}
	js, _ := json.Marshal(e)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("ad/ad_edit.html", d)
}

// 删除广告
func (this *adC) Del_post(ctx *echox.Context) error {
	ctx.Request.ParseForm()
	form := ctx.Request.Form
	var result gof.Message
	partnerId := getPartnerId(ctx)
	adId, _ := strconv.Atoi(form.Get("id"))
	err := dps.AdvertisementService.DelAdvertisement(partnerId, adId)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}

	return ctx.JSON(http.StatusOK, result)
}

func (this *adC) SaveAd_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
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
	return ctx.JSON(http.StatusOK, result)
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
	return ctx.RenderOK("ad/ad_data3.html", d)
}

// 创建广告图片
func (this *adC) CreateAdImage(ctx *echox.Context) error {
	adId, _ := strconv.Atoi(ctx.Query("ad_id"))
	e := ad.ValueImage{
		Enabled:         1,
		AdvertisementId: adId,
		LinkUrl:         "http://",
		ImageUrl:        ctx.App.Config().GetString(variable.NoPicPath),
	}

	js, _ := json.Marshal(e)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("ad/ad_image.html", d)
}

// 保存广告
func (this *adC) EditAdImage(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	adId, _ := strconv.Atoi(ctx.Query("ad_id"))
	imgId, _ := strconv.Atoi(ctx.Query("id"))

	e := dps.AdvertisementService.GetValueAdImage(partnerId, adId, imgId)

	js, _ := json.Marshal(e)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("ad/ad_image.html", d)
}

func (this *adC) SaveImage_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
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
	return ctx.JSON(http.StatusOK, result)
}

func (this *adC) Del_image_post(ctx *echox.Context) error {
	ctx.Request.ParseForm()
	form := ctx.Request.Form
	var result gof.Message
	partnerId := getPartnerId(ctx)
	adId, _ := strconv.Atoi(form.Get("ad_id"))
	imgId, _ := strconv.Atoi(form.Get("id"))
	err := dps.AdvertisementService.DelAdImage(partnerId, adId, imgId)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}

	return ctx.JSON(http.StatusOK, result)
}
