/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package mobi

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"net/http"
)

var (
	routes *web.RouteMap = new(web.RouteMap)
)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

func handleError(w http.ResponseWriter, err error) {
	w.Write([]byte(`<span style="color:red">` + err.Error() + `</span>`))
}

//注册路由
func RegisterRoutes(c gof.App) {
	mc := &mainC{App: c}

	getPartner := func(r *http.Request) (*partner.ValuePartner, error) {
		partnerId := dps.PartnerService.GetPartnerIdByHost(r.Host)
		return dps.PartnerService.GetPartner(partnerId)
	}

	routes.Add("/", func(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		if p, err := getPartner(r); err == nil {
			mvc.Handle(mc, ctx, true, p)
		} else {
			handleError(w, err)
		}
	})
}
