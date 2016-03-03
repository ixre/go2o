package partner

import (
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
<<<<<<< HEAD
	"go2o/src/core/domain/interface/delivery"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
)

type coverageAreaC struct {
}

func (this *coverageAreaC) CoverageAreList(ctx *echox.Context) error {
	d := echox.NewRenderData()
	return ctx.RenderOK("delivery.coverage_area_list.html", d)
}

func (this *coverageAreaC) Create(ctx *echox.Context) error {
	d := echox.NewRenderData()
	d.Map["entity"] = template.JS("{}")
	return ctx.RenderOK("delivery.create_area.html", d)
}

// 保存配送区域(POST)
func (this *coverageAreaC) SaveArea(ctx *echox.Context) error {
	r := ctx.Request()
	if r.Method == "POST" {
		var result gof.Message
		r.ParseForm()

		coverageArea := delivery.CoverageValue{}
		web.ParseFormToEntity(r.Form, &coverageArea)

		id, err := dps.DeliverService.CreateCoverageArea(&coverageArea)

		if err != nil {
			result = gof.Message{Message: err.Error()}
		} else {
			result = gof.Message{Result: true, Message: "", Data: id}
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
=======
	"github.com/jsix/gof/web/mvc"
	"go2o/src/core/domain/interface/delivery"
	"go2o/src/core/service/dps"
	"html/template"
)

var _ mvc.Filter = new(coverageAreaC)

type coverageAreaC struct {
	*baseC
}

func (this *coverageAreaC) CoverageAreList(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.Response, nil, "views/partner/delivery/coverage_area_list.html")
}

func (this *coverageAreaC) Create(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.Response, gof.TemplateDataMap{
		"entity": template.JS("{}"),
	}, "views/partner/delivery/create.html")
}

func (this *coverageAreaC) SaveArea_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
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
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
