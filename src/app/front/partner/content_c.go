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
	"go2o/src/core/domain/interface/content"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
	"time"
	"github.com/jsix/gof/web/form"
	"fmt"
)

type contentC struct {
}

//商品列表
func (this *contentC) Page_list(ctx *echox.Context) error {

	d := ctx.NewData()
	return ctx.RenderOK("content.page_list.html", d)
}

// 修改页面
func (this *contentC) Page_edit(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	id, _ := strconv.Atoi(ctx.Query("id"))
	e := dps.ContentService.GetPage(partnerId, id)

	js, _ := json.Marshal(e)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("content.page_edit.html", d)
}

// 保存页面
func (this *contentC) Page_create(ctx *echox.Context) error {
	e := content.ValuePage{
		Enabled: 1,
	}

	js, _ := json.Marshal(e)
	d := ctx.NewData()
	d.Map["entity"] = template.JS(js)
	return ctx.RenderOK("content.page_edit.html", d)
}

func (this *contentC) SavePage(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {

		r.ParseForm()

		var result gof.Result

		e := content.ValuePage{}
		form.ParseEntity(r.Form, &e)

		//更新
		e.UpdateTime = time.Now().Unix()
		e.PartnerId = partnerId

		id, err := dps.ContentService.SavePage(partnerId, &e)

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

func (this *contentC) Page_del(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		r.ParseForm()

		var result gof.Result
		id, _ := strconv.Atoi(r.FormValue("id"))
		err := dps.ContentService.DeletePage(partnerId, id)

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}
