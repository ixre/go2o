/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:55
 * description :
 * history :
 */

package merchant

import (
	"encoding/json"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
)

type saleC struct {
}

func (this *saleC) TagList(ctx *echox.Context) error {
	d := ctx.NewData()
	return ctx.RenderOK("sale.tag_list.html", d)
}

//修改门店信息
func (this *saleC) Edit_stag(ctx *echox.Context) error {
	merchantId := getMerchantId(ctx)
	r := ctx.HttpRequest()
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	entity := dps.SaleService.GetSaleTag(merchantId, id)
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
	merchantId := getMerchantId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		var result gof.Message
		r.ParseForm()

		e := sale.ValueSaleTag{}
		web.ParseFormToEntity(r.Form, &e)
		e.MerchantId = getMerchantId(ctx)

		id, err := dps.SaleService.SaveSaleTag(merchantId, &e)

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

// 删除销售标签(POST)
func (this *saleC) Del_stag(ctx *echox.Context) error {
	r := ctx.HttpRequest()
	var result gof.Message
	if r.Method == "POST" {
		r.ParseForm()
		merchantId := getMerchantId(ctx)
		id, err := strconv.Atoi(r.FormValue("id"))
		if err == nil {
			err = dps.SaleService.DeleteSaleTag(merchantId, id)
		}

		if err != nil {
			result.Message = err.Error()
		} else {
			result.Result = true
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}
