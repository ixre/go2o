package ucenter

import (
	"com/domain/interface/member"
	"com/domain/interface/partner"
	"com/ording/entity"
	"com/service/goclient"
	"net/http"
	"net/url"
	"github.com/newmin/gof/app"
	"github.com/newmin/gof/web"
	"github.com/newmin/gof/web/mvc"
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

func redirect(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<script>window.parent.location.href='/login?return_url=" +
		url.QueryEscape(r.URL.String()) + "'</script>"))
}

func RegisterRoutes(c app.Context) {
	mc := &mainC{Context: c}
	oc := &orderC{Context: c}
	ac := &accountC{Context: c}
	lc := &loginC{Context: c}

	routes.Add("^/order/", func(w http.ResponseWriter, r *http.Request) {
		if m, p, host := chkLogin(r); m != nil {
			mvc.HandleRequest(oc, w, r, true, m, p, host)
		} else {
			redirect(w, r)
		}
	})

	routes.Add("^/account/", func(w http.ResponseWriter, r *http.Request) {
		if m, p, host := chkLogin(r); m != nil {
			mvc.HandleRequest(ac, w, r, true, m, p, host)
		} else {
			redirect(w, r)
		}
	})

	routes.Add("^/login", func(w http.ResponseWriter, r *http.Request) {
		mvc.HandleRequest(lc, w, r, true)
	})

	routes.Add("/", func(w http.ResponseWriter, r *http.Request) {
		if m, p, host := chkLogin(r); m != nil {
			mvc.HandleRequest(mc, w, r, true, m, p, host)
		} else {
			redirect(w, r)
		}
	})
}

func chkLogin(r *http.Request) (m *member.ValueMember, p *entity.Partner, conf *partner.SiteConf) {
	cookie, err := r.Cookie("ms_token")
	if err != nil {
		return nil, nil, nil
	}
	arr := strings.Split(cookie.Value, "$")
	id, _ := strconv.Atoi(arr[0])
	token := arr[1]

	m, err = goclient.Member.GetMember(id, token)

	if err != nil {
		return nil, nil, nil
	}
	m.LoginToken = token

	p, err = goclient.Member.GetBindPartner(id, token)
	if err != nil {
		return nil, nil, nil
	}

	siteConf, err := goclient.Partner.GetSiteConf(p.Id, p.Secret)
	//host, err = goclient.Partner.GetHost(p.Id, p.Secret)
	return m, p, siteConf
}
