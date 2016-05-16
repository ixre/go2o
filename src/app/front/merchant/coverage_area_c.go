package merchant

import (
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"go2o/src/core/domain/interface/delivery"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
)

type coverageAreaC struct {
}

func (this *coverageAreaC) CoverageAreList(ctx *echox.Context) error {
	d := ctx.NewData()
	return ctx.RenderOK("delivery.coverage_area_list.html", d)
}

func (this *coverageAreaC) Create(ctx *echox.Context) error {
	d := ctx.NewData()
	d.Map["entity"] = template.JS("{}")
	return ctx.RenderOK("delivery.create_area.html", d)
}

// 保存配送区域(POST)
func (this *coverageAreaC) SaveArea(ctx *echox.Context) error {
	r := ctx.HttpRequest()
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
}
