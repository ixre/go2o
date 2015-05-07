package partner

import (
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"html/template"
    "github.com/atnet/gof"
    "go2o/src/core/domain/interface/delivery"
    "go2o/src/core/service/dps"
    "fmt"
)

var _ mvc.Filter = new(converageAreaC)

type converageAreaC struct {
	Base *baseC
}

func (this *converageAreaC) Requesting(ctx *web.Context) bool {
	return this.Base.Requesting(ctx)
}
func (this *converageAreaC) RequestEnd(ctx *web.Context) {
	this.Base.RequestEnd(ctx)
}

func (this *converageAreaC) CoverageAreList(ctx *web.Context) {
	ctx.App.Template().Render(ctx.ResponseWriter, "views/partner/delivery/converage_area_list.html", nil)
}

func (this *converageAreaC) Create(ctx *web.Context) {
	ctx.App.Template().Render(ctx.ResponseWriter, "views/partner/delivery/create.html", func(m *map[string]interface{}) {
		(*m)["entity"] = template.JS("{}")
	})
}

func (this *converageAreaC) SaveArea_post(ctx *web.Context) {
    r, w := ctx.Request, ctx.ResponseWriter
    var result gof.JsonResult
    r.ParseForm()

    converageArea := delivery.CoverageValue{}
    web.ParseFormToEntity(r.Form, &converageArea)

    id, err := dps.DeliveyService.CreateConverageArea(&converageArea)
    fmt.Println(id,"=====\n ERROR:",err)
    if err != nil {
        result = gof.JsonResult{Result: true, Message: err.Error()}
    } else {
        result = gof.JsonResult{Result: true, Message: "", Data: id}
    }
    w.Write(result.Marshal())
}