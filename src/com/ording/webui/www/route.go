package www

import (
	"com/domain/interface/member"
	"com/ording/dao"
	"com/ording/entity"
	"com/service/goclient"
	"net/http"
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"strconv"
	"strings"
)

var (
	routes *web.RouteMap = new(web.RouteMap)
)

//处理请求
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	routes.HandleRequest(w, r)
}

func handleError(w http.ResponseWriter, err error) {
	w.Write([]byte(`<html><head><title></title></head>` +
		`<body><span style="color:red">` + err.Error() + `</span></body></html>`))
}

//注册路由
func RegisterRoutes(c app.Context) {
	mc := &mainC{Context: c}
	dc := &ordingC{Context: c}
	sp := &shoppingC{Context: c}

	getPartner := func(r *http.Request) (*entity.Partner, error, *member.ValueMember) {
		var m *member.ValueMember
		cookie, err := r.Cookie("ms_token")
		if err == nil {
			arr := strings.Split(cookie.Value, "$")
			id, _ := strconv.Atoi(arr[0])
			token := arr[1]

			m, err = goclient.Member.GetMember(id, token)
			if err == nil {
				m.LoginToken = token
			}
		}

		p, err := dao.Partner().GetPartnerByHost(r.Host)

		return p, err, m
	}

	routes.Add("/ding/", func(w http.ResponseWriter, r *http.Request) {
		if p, err, m := getPartner(r); err == nil {
			mvc.HandleRequest(dc, w, r, true, p, m)
		} else {
			handleError(w, err)
		}
	})

	routes.Add("/shopping/", func(w http.ResponseWriter, r *http.Request) {
		if p, err, m := getPartner(r); err == nil {
			mvc.HandleRequest(sp, w, r, true, p, m)
		} else {
			handleError(w, err)
		}
	})

	routes.Add("/", func(w http.ResponseWriter, r *http.Request) {
		if p, err, m := getPartner(r); err == nil {
			mvc.HandleRequest(mc, w, r, true, p, m)
		} else {
			handleError(w, err)
		}
	})
}
