/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:55
 * description :
 * history :
 */

package partner

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/service/dps"
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
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	entity := dps.SaleService.GetSaleTag(partnerId, id)
	bys, _ := json.Marshal(entity)

	ctx.App.Template().Execute(w,
		gof.TemplateDataMap{
			"entity": template.JS(bys),
		},
		"views/partner/sale/edit_sale_tag.html")
}

func (this *saleC) Create_stag(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"entity": "{}",
		},
		"views/partner/sale/create_sale_tag.html")
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
}
