package partner

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/delivery"
	"go2o/src/core/service/dps"
	"html/template"
)

var _ mvc.Filter = new(coverageAreaC)

type coverageAreaC struct {
	*baseC
}

func (this *coverageAreaC) CoverageAreList(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.ResponseWriter, nil, "views/partner/delivery/coverage_area_list.html")
}

func (this *coverageAreaC) Create(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.ResponseWriter, gof.TemplateDataMap{
		"entity": template.JS("{}"),
	}, "views/partner/delivery/create.html")
}

func (this *coverageAreaC) SaveArea_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.Message
	r.ParseForm()

	coverageArea := delivery.CoverageValue{}
	web.ParseFormToEntity(r.Form, &coverageArea)

	id, err := dps.DeliverService.CreateCoverageArea(&coverageArea)

	if err != nil {
		result = gof.Message{Result: true, Message: err.Error()}
	} else {
		result = gof.Message{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}
