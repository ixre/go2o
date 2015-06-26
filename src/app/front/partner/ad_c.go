/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/service/dps"
	"html/template"
	"time"
	"strconv"
	"go2o/src/core/domain/interface/content"
)

var _ mvc.Filter = new(adC)

// 广告控制器
type adC struct {
	*baseC
}

//广告列表
func (this *contentC) List(ctx *web.Context) {
	dps.AdvertisementService.GetAdvertisement(this.GetPartnerId(ctx),0)
	ctx.App.Template().Execute(ctx.ResponseWriter, gof.TemplateDataMap{
	}, "views/partner/ad/ad_list.html")
}

// 修改广告
func (this *contentC) Edit(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	form := ctx.Request.URL.Query()
	id, _ := strconv.Atoi(form.Get("id"))
	e := dps.AdvertisementService.GetAdvertisement(partnerId,id)

	js, _ := json.Marshal(e)

	ctx.App.Template().Execute(ctx.ResponseWriter,
		gof.TemplateDataMap{
			"entity": template.JS(js),
		},
		"views/partner/ad/ad_edit.html")
}

// 保存广告
func (this *contentC) Create(ctx *web.Context) {
	e := content.ValuePage{
		Enabled:1,
	}

	js, _ := json.Marshal(e)

	ctx.App.Template().Execute(ctx.ResponseWriter,
		gof.TemplateDataMap{
			"entity": template.JS(js),
		},
		"views/partner/ad/ad_edit.html")
}

func (this *contentC) SaveAd_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r := ctx.Request
	r.ParseForm()

	var result gof.Message

	e := content.ValuePage{}
	web.ParseFormToEntity(r.Form, &e)

	//更新
	e.UpdateTime = time.Now().Unix()
	e.PartnerId = partnerId

	id, err := dps.ContentService.SavePage(partnerId, &e)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
		result.Data = id
	}
	this.ResultOutput(ctx,result)
}
