/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:55
 * description :
 * history :
 */

package partner

import (
	"encoding/json"
	"github.com/jsix/gof"
	//"github.com/jsix/gof/web"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
	"github.com/jsix/gof/web/form"
	"fmt"
)

type saleC struct {
}

func (this *saleC) TagList(ctx *echox.Context) error {
	d := ctx.NewData()
	return ctx.RenderOK("sale.tag_list.html", d)
}

//修改门店信息
func (this *saleC) Edit_stag(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	entity := dps.SaleService.GetSaleTag(partnerId, id)
	bys, _ := json.Marshal(entity)

	d := ctx.NewData()
	d.Map["entity"] = template.JS(bys)
	return ctx.RenderOK("sale.tag_form.html", d)
}

func (this *saleC) Create_stag(ctx *echox.Context) error {
	entity := sale.ValueSaleTag{
		GoodsImage: ctx.App.Config().GetString(variable.NoPicPath),
	}
	bys, _ := json.Marshal(entity)

	d := ctx.NewData()
	d.Map["entity"] = template.JS(bys)
	return ctx.RenderOK("sale.tag_form.html", d)
}

// 保存销售标签(POST)
func (this *saleC) Save_stag(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		var result gof.Result
		r.ParseForm()

		e := sale.ValueSaleTag{}
		form.ParseEntity(r.Form, &e)
		e.PartnerId = getPartnerId(ctx)

		id, err := dps.SaleService.SaveSaleTag(partnerId, &e)

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

// 删除销售标签(POST)
func (this *saleC) Del_stag(ctx *echox.Context) error {
	r := ctx.HttpRequest()
	var result gof.Result
	if r.Method == "POST" {
		r.ParseForm()
		partnerId := getPartnerId(ctx)
		id, err := strconv.Atoi(r.FormValue("id"))
		if err == nil {
			err = dps.SaleService.DeleteSaleTag(partnerId, id)
		}

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}
