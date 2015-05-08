package partner

import (
	"fmt"
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
	ctx.App.Template().Render(ctx.ResponseWriter, "views/partner/delivery/coverage_area_list.html", nil)
}

func (this *coverageAreaC) Create(ctx *web.Context) {
	ctx.App.Template().Render(ctx.ResponseWriter, "views/partner/delivery/create.html", func(m *map[string]interface{}) {
		(*m)["entity"] = template.JS("{}")
	})
}

func (this *coverageAreaC) SaveArea_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.JsonResult
	r.ParseForm()

	coverageArea := delivery.CoverageValue{}
	web.ParseFormToEntity(r.Form, &coverageArea)

	id, err := dps.DeliveyService.CreateCoverageArea(&coverageArea)
	fmt.Println(id, "=====\n ERROR:", err)
	if err != nil {
		result = gof.JsonResult{Result: true, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}
