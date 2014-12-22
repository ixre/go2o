package mobi

import (
	"com/ording/dao"
	"com/ording/entity"
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

func handleError(w http.ResponseWriter, err error) {
	w.Write([]byte(`<span style="color:red">` + err.Error() + `</span>`))
}

//注册路由
func RegistRoutes(c app.Context) {
	mc := &mainC{Context: c}

	getPartner := func(r *http.Request) (*entity.Partner, error) {
		return dao.Partner().GetPartnerByHost(r.Host)
	}

	routes.Add("/", func(w http.ResponseWriter, r *http.Request) {
		if p, err := getPartner(r); err == nil {
			mvc.HandleRequest(mc, w, r, true, p)
		} else {
			handleError(w, err)
		}
	})
}
