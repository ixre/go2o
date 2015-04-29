/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/goclient"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	routes *web.RouteMap = new(web.RouteMap)
)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

func redirect(ctx *web.Context) {
	ctx.ResponseWriter.Write([]byte("<script>window.parent.location.href='/login?return_url=" +
		url.QueryEscape(ctx.Request.URL.String()) + "'</script>"))
}

func RegisterRoutes(c gof.App) {
	mc := &mainC{App: c}
	oc := &orderC{App: c}
	ac := &accountC{App: c}
	lc := &loginC{App: c}

	routes.Add("^/order/", func(ctx *web.Context) {
		if m, p, host := chkLogin(ctx.Request); m != nil {
			mvc.Handle(oc, ctx, true, m, p, host)
		} else {
			redirect(ctx)
		}
	})

	routes.Add("^/account/", func(ctx *web.Context) {
		if m, p, host := chkLogin(ctx.Request); m != nil {
			mvc.Handle(ac, ctx, true, m, p, host)
		} else {
			redirect(ctx)
		}
	})

	routes.Add("^/login", func(ctx *web.Context) {
		mvc.Handle(lc, ctx, true)
	})

	routes.Add("/", func(ctx *web.Context) {
		if m, p, host := chkLogin(ctx.Request); m != nil {
			mvc.Handle(mc, ctx, true, m, p, host)
		} else {
			redirect(ctx)
		}
	})
}

func chkLogin(r *http.Request) (m *member.ValueMember, p *partner.ValuePartner, conf *partner.SiteConf) {
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
