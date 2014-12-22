package apiserv

import (
	"net/http"
	"ops/cf/app"
	"ops/cf/web"
	"ops/cf/web/mvc"
)

var (
	routes *web.RouteMap = new(web.RouteMap)
)

//处理请求
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	routes.HandleRequest(w, r)
}

//
//func redirect(weixin http.ResponseWriter, r *http.Request) {
//	weixin.Write([]byte("<script>window.parent.location.href='/login?return_url=" +
//		url.QueryEscape(r.URL.String()) + "'</script>"))
//}

func RegistRoutes(c app.Context) {
	ws := &websocketC{Context: c}

	routes.Add("^/ws/", func(w http.ResponseWriter, r *http.Request) {
		//cross ajax request
		w.Header().Add("Access-Control-Allow-Origin", "*")
		mvc.HandleRequest(ws, w, r, false)
	})

	routes.Add("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("page not found"))
	})
}
