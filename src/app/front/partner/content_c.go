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
	"go2o/src/core/domain/interface/content"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
=======
	"github.com/jsix/gof/web/mvc"
	"go2o/src/core/domain/interface/content"
	"go2o/src/core/service/dps"
	"html/template"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"strconv"
	"time"
)

<<<<<<< HEAD
type contentC struct {
}

//商品列表
func (this *contentC) Page_list(ctx *echox.Context) error {

	d := echox.NewRenderData()
	return ctx.RenderOK("content.page_list.html", d)
}

// 修改页面
func (this *contentC) Page_edit(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	id, _ := strconv.Atoi(ctx.Query("id"))
	e := dps.ContentService.GetPage(partnerId, id)

	js, _ := json.Marshal(e)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("content.page_edit.html", d)
}

// 保存页面
func (this *contentC) Page_create(ctx *echox.Context) error {
=======
var _ mvc.Filter = new(contentC)

type contentC struct {
	*baseC
}

//商品列表
func (this *contentC) Page_list(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.Response, gof.TemplateDataMap{}, "views/partner/content/page_list.html")
}

// 修改页面
func (this *contentC) Page_edit(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	form := ctx.Request.URL.Query()
	id, _ := strconv.Atoi(form.Get("id"))
	e := dps.ContentService.GetPage(partnerId, id)

	js, _ := json.Marshal(e)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity": template.JS(js),
		},
		"views/partner/content/page_edit.html")
}

// 保存页面
func (this *contentC) Page_create(ctx *web.Context) {
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	e := content.ValuePage{
		Enabled: 1,
	}

	js, _ := json.Marshal(e)
<<<<<<< HEAD
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("content.page_edit.html", d)
}

func (this *contentC) SavePage(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	if r.Method == "POST" {

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
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

func (this *contentC) Page_del(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	if r.Method == "POST" {
		r.ParseForm()

		var result gof.Message
		id, _ := strconv.Atoi(r.FormValue("id"))
		err := dps.ContentService.DeletePage(partnerId, id)

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
		"views/partner/content/page_edit.html")
}

func (this *contentC) SavePage_post(ctx *web.Context) {
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
	ctx.Response.JsonOutput(result)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
