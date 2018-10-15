package partner

import (
	"github.com/jsix/gof"
	//"github.com/jsix/gof/web"
	"go2o/src/core/domain/interface/delivery"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"github.com/jsix/gof/web/form"
	"fmt"
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
		var result gof.Result
		r.ParseForm()

		coverageArea := delivery.CoverageValue{}
		form.ParseEntity(r.Form, &coverageArea)

		id, err := dps.DeliverService.CreateCoverageArea(&coverageArea)

		if err != nil {
			result = gof.Result{ErrMsg: err.Error()}
		} else {
			var data = make(map[string]string)
			data["id"] = fmt.Sprintf("%d", id)
			result = gof.Result{ErrCode: 0, ErrMsg: "", Data: data}
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}
