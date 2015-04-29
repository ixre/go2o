/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package www

import (
	"errors"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"go2o/src/core/service/goclient"
	"net/http"
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

func handleError(w http.ResponseWriter, err error) {
	w.Write([]byte(`<html><head><title></title></head>` +
		`<body><span style="color:red">` + err.Error() + `</span></body></html>`))
}

//注册路由
func RegisterRoutes(c gof.App) {
	mc := &mainC{App: c}
	sp := &shoppingC{App: c}
	pc := &paymentC{App: c}

	var httpFunc web.HttpContextFunc = func(ctx *web.Context) {
		if err, ok := recover().(error); ok {
			handleCustomError(ctx.ResponseWriter, c, err)
		}
	}
	routes.Defer(httpFunc)

	getPartner := func(r *http.Request) (*partner.ValuePartner, error, *member.ValueMember) {
		var m *member.ValueMember
		cookie, err := r.Cookie("ms_token")
		if err == nil {
			if len(cookie.Value) == 0 {
				err = errors.New("empty cookie")
			}

			arr := strings.Split(cookie.Value, "$")
			id, _ := strconv.Atoi(arr[0])
			token := arr[1]

			m, err = goclient.Member.GetMember(id, token)
			if err == nil {
				m.LoginToken = token
			}
		}

		partnerId := dps.PartnerService.GetPartnerIdByHost(r.Host)
		p, err := dps.PartnerService.GetPartner(partnerId)
		return p, err, m
	}

	// 购物车
	routes.Add("^/cart_api$", func(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		if p, err, m := getPartner(r); err == nil {
			sp.CartApi(ctx, p, m)
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/cart/*", func(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		if p, err, _ := getPartner(r); err == nil {
			sp.Cart(ctx, p)
		} else {
			handleError(w, err)
		}
	})

	// 购买跳转
	routes.Add("^/buy/*$", func(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		if p, err, m := getPartner(r); err == nil {
			sp.BuyRedirect(ctx, p, m)
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/buy/confirm$", func(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		if p, err, m := getPartner(r); err == nil {
			sp.OrderConfirm(ctx, p, m)
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/buy/apply/coupon$", func(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		if p, err, m := getPartner(r); err == nil {
			if r.Method == "POST" {
				sp.ApplyCoupon_post(ctx, p, m)
			}
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/buy/order/persist", func(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		if p, err, m := getPartner(r); err == nil {
			if r.Method == "POST" {
				sp.OrderPersist_post(ctx, p, m)
			}
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/buy/order/submit_order$", func(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		if p, err, m := getPartner(r); err == nil {
			if r.Method == "POST" {
				sp.SubmitOrder_post(ctx, p, m)
			}
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/buy/order/finish$", func(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		if p, err, m := getPartner(r); err == nil {
			sp.OrderFinish(ctx, p, m)
		} else {
			handleError(w, err)
		}
	})

	routes.Add("/(shopping|buy)/", func(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		if p, err, m := getPartner(r); err == nil {
			mvc.Handle(sp, ctx, true, p, m)
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/pay/", func(ctx *web.Context) {
		mvc.Handle(pc, ctx, true)
	})

	routes.Add("/", func(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		if p, err, m := getPartner(r); err == nil {
			mvc.Handle(mc, ctx, true, p, m)
		} else {
			handleError(w, err)
		}
	})
}
