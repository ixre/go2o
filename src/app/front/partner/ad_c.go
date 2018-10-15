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
	//"github.com/jsix/gof/web"
	"go2o/src/core/domain/interface/ad"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
	"github.com/jsix/gof/web/form"
	"fmt"
)

// 广告控制器
type adC struct {
}

//广告列表
func (this *adC) List(ctx *echox.Context) error {
	return ctx.RenderOK("ad.list.html", ctx.NewData())
}

// 修改广告
func (this *adC) Edit(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	id, _ := strconv.Atoi(ctx.Query("id"))
	e := dps.AdvertisementService.GetAdvertisement(partnerId, id)

	js, _ := json.Marshal(e)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("ad.edit.html", d)
}

// 保存广告
func (this *adC) Create(ctx *echox.Context) error {
	e := ad.ValueAdvertisement{
		Enabled: 1,
	}
	js, _ := json.Marshal(e)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("ad.edit.html", d)
}

// 删除广告(POST)
func (this *adC) Delete_ad(ctx *echox.Context) error {
	req := ctx.HttpRequest()
	if req.Method == "POST" {
		req.ParseForm()
		var result gof.Result
		partnerId := getPartnerId(ctx)
		adId, _ := strconv.Atoi(req.FormValue("id"))
		err := dps.AdvertisementService.DelAdvertisement(partnerId, adId)

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
		}

		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 保存广告(POST)
func (this *adC) SaveAd(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	req := ctx.HttpRequest()
	if req.Method == "POST" {
		req.ParseForm()

		var result gof.Result

		e := ad.ValueAdvertisement{}
		form.ParseEntity(req.Form, &e)

		//更新
		e.PartnerId = partnerId

		id, err := dps.AdvertisementService.SaveAdvertisement(partnerId, &e)

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

func (this *adC) Ad_data1(ctx *echox.Context) error {
	return ctx.String(http.StatusOK, `<span style="color:red">暂只支持轮播广告</span>`)
}

func (this *adC) Ad_data2(ctx *echox.Context) error {
	return ctx.String(http.StatusOK, `<span style="color:red">暂只支持轮播广告</span>`)
}

//轮播广告
func (this *adC) Ad_data3(ctx *echox.Context) error {
	d := ctx.NewData()
	d.Map["adId"] = ctx.Query("id")
	return ctx.RenderOK("ad.data3.html", d)
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
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("ad.image.html", d)
}

// 保存广告
func (this *adC) EditAdImage(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	adId, _ := strconv.Atoi(ctx.Query("ad_id"))
	imgId, _ := strconv.Atoi(ctx.Query("id"))

	e := dps.AdvertisementService.GetValueAdImage(partnerId, adId, imgId)

	js, _ := json.Marshal(e)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("ad.image.html", d)
}

// 保存图片(POST)
func (this *adC) SaveImage(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		r.ParseForm()

		var result gof.Result

		e := ad.ValueImage{}
		form.ParseEntity(r.Form, &e)

		id, err := dps.AdvertisementService.SaveImage(partnerId, e.AdvertisementId, &e)

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

// 删除广告图片(POST)
func (this *adC) Delete_image(ctx *echox.Context) error {
	req := ctx.HttpRequest()
	if req.Method == "POST" {
		req.ParseForm()
		var result gof.Result
		partnerId := getPartnerId(ctx)
		adId, _ := strconv.Atoi(req.FormValue("ad_id"))
		imgId, _ := strconv.Atoi(req.FormValue("id"))
		err := dps.AdvertisementService.DelAdImage(partnerId, adId, imgId)

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
		}

		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}
