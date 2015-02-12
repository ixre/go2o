/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package mobi

import (
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/core/domain/interface/partner"
	"go2o/core/service/dps"
	"net/http"
)

var (
	routes *web.RouteMap = new(web.RouteMap)
)

//处理请求
func Handle(w http.ResponseWriter, r *http.Request) {
	routes.Handle(w, r)
}

func handleError(w http.ResponseWriter, err error) {
	w.Write([]byte(`<span style="color:red">` + err.Error() + `</span>`))
}

//注册路由
func RegisterRoutes(c app.Context) {
	mc := &mainC{Context: c}

	getPartner := func(r *http.Request) (*partner.ValuePartner, error) {
		partnerId := dps.PartnerService.GetPartnerIdByHost(r.Host)
		return dps.PartnerService.GetPartner(partnerId)
	}

	routes.Add("/", func(w http.ResponseWriter, r *http.Request) {
		if p, err := getPartner(r); err == nil {
			mvc.Handle(mc, w, r, true, p)
		} else {
			handleError(w, err)
		}
	})
}
