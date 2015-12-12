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
	"go2o/src/core/domain/interface/content"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var _ mvc.Filter = new(contentC)

type contentC struct {
}

//商品列表
func (this *contentC) Page_list(ctx *echox.Context) error {

	d := echox.NewRenderData()
	return ctx.Render(http.StatusOK, "content/page_list.html", d)
}

// 修改页面
func (this *contentC) Page_edit(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	form := ctx.Request.URL.Query()
	id, _ := strconv.Atoi(form.Get("id"))
	e := dps.ContentService.GetPage(partnerId, id)

	js, _ := json.Marshal(e)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.Render(http.StatusOK, "content/page_edit.html", d)
}

// 保存页面
func (this *contentC) Page_create(ctx *echox.Context) error {
	e := content.ValuePage{
		Enabled: 1,
	}

	js, _ := json.Marshal(e)
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS(js)
	return ctx.Render(http.StatusOK, "content/page_edit.html", d)
}

func (this *contentC) SavePage_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
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
	return ctx.JSON(http.StatusOK, result)
}
