/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package apiserv

import (
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"net/http"
)

var (
	routes *web.RouteMap = new(web.RouteMap)
)

//处理请求
func Handle(w http.ResponseWriter, r *http.Request) {
	routes.Handle(w, r)
}

//
//func redirect(weixin http.ResponseWriter, r *http.Request) {
//	weixin.Write([]byte("<script>window.parent.location.href='/login?return_url=" +
//		url.QueryEscape(r.URL.String()) + "'</script>"))
//}

func RegisterRoutes(c app.Context) {
	ws := &websocketC{Context: c}

	routes.Add("^/ws/", func(w http.ResponseWriter, r *http.Request) {
		//cross ajax request
		w.Header().Add("Access-Control-Allow-Origin", "*")
		mvc.Handle(ws, w, r, false)
	})

	routes.Add("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("page not found"))
	})
}
