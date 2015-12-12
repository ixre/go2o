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

func (this *saleC) TagList(ctx *echox.Context) error {
	ctx.App.Template().Execute(ctx.Response, nil, "views/partner/sale/sale_tag_list.html")
}

//修改门店信息
func (this *saleC) Edit_stag(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	entity := dps.SaleService.GetSaleTag(partnerId, id)
	bys, _ := json.Marshal(entity)

	ctx.App.Template().Execute(w,
		gof.TemplateDataMap{
			"entity": template.JS(bys),
		},
		"views/partner/sale/sale_tag_form.html")
}

func (this *saleC) Create_stag(ctx *echox.Context) error {
	entity := sale.ValueSaleTag{
		GoodsImage: ctx.App.Config().GetString(variable.NoPicPath),
	}
	bys, _ := json.Marshal(entity)

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity": template.JS(bys),
		},
		"views/partner/sale/sale_tag_form.html")
}

func (this *saleC) Save_stag_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request
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

func (this *saleC) Del_stag_post(ctx *echox.Context) error {
	r := ctx.Request
	var result gof.Message
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
