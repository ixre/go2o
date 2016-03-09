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
	"github.com/jsix/gof/web"
<<<<<<< HEAD
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
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
=======
	"github.com/jsix/gof/web/mvc"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"html/template"
	"strconv"
)

var _ mvc.Filter = new(saleC)

type saleC struct {
	*baseC
}

func (this *saleC) TagList(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.Response, nil, "views/partner/sale/sale_tag_list.html")
}

//修改门店信息
func (this *saleC) Edit_stag(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	entity := dps.SaleService.GetSaleTag(partnerId, id)
	bys, _ := json.Marshal(entity)

<<<<<<< HEAD
	d := ctx.NewData()
	d.Map["entity"] = template.JS(bys)
	return ctx.RenderOK("sale.tag_form.html", d)
}

func (this *saleC) Create_stag(ctx *echox.Context) error {
=======
	ctx.App.Template().Execute(w,
		gof.TemplateDataMap{
			"entity": template.JS(bys),
		},
		"views/partner/sale/sale_tag_form.html")
}

func (this *saleC) Create_stag(ctx *web.Context) {
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	entity := sale.ValueSaleTag{
		GoodsImage: ctx.App.Config().GetString(variable.NoPicPath),
	}
	bys, _ := json.Marshal(entity)

<<<<<<< HEAD
	d := ctx.NewData()
	d.Map["entity"] = template.JS(bys)
	return ctx.RenderOK("sale.tag_form.html", d)
}

// 保存销售标签(POST)
func (this *saleC) Save_stag(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	if r.Method == "POST" {
		var result gof.Message
		r.ParseForm()

		e := sale.ValueSaleTag{}
		web.ParseFormToEntity(r.Form, &e)
		e.PartnerId = getPartnerId(ctx)

		id, err := dps.SaleService.SaveSaleTag(partnerId, &e)

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
	r := ctx.Request()
	var result gof.Message
	if r.Method == "POST" {
		r.ParseForm()
		partnerId := getPartnerId(ctx)
		id, err := strconv.Atoi(r.FormValue("id"))
		if err == nil {
			err = dps.SaleService.DeleteSaleTag(partnerId, id)
		}

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
			"entity": template.JS(bys),
		},
		"views/partner/sale/sale_tag_form.html")
}

func (this *saleC) Save_stag_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r := ctx.Request
	var result gof.Message
	r.ParseForm()

	e := sale.ValueSaleTag{}
	web.ParseFormToEntity(r.Form, &e)
	e.PartnerId = this.GetPartnerId(ctx)

	id, err := dps.SaleService.SaveSaleTag(partnerId, &e)

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
		result.Data = id
	}
	ctx.Response.JsonOutput(result)
}

func (this *saleC) Del_stag_post(ctx *web.Context) {
	r := ctx.Request
	var result gof.Message
	r.ParseForm()
	partnerId := this.GetPartnerId(ctx)
	id, err := strconv.Atoi(r.FormValue("id"))
	if err == nil {
		err = dps.SaleService.DeleteSaleTag(partnerId, id)
	}

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	ctx.Response.JsonOutput(result)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
