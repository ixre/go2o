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
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/partner"
	"go2o/core/service/dps"
	"go2o/core/service/goclient"
	"net/http"
	"strconv"
	"strings"
)

var (
	routes *web.RouteMap = new(web.RouteMap)
)

//处理请求
func Handle(w http.ResponseWriter, r *http.Request) {
	routes.Handle(w, r)
}

func handleError(w http.ResponseWriter, err error) {
	w.Write([]byte(`<html><head><title></title></head>` +
		`<body><span style="color:red">` + err.Error() + `</span></body></html>`))
}

//注册路由
func RegisterRoutes(c app.Context) {
	mc := &mainC{Context: c}
	sp := &shoppingC{Context: c}
	pc := &paymentC{Context: c}

	var httpFunc http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		if err, ok := recover().(error); ok {
			handleCustomError(w, c, err)
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
	routes.Add("^/cart_api$", func(w http.ResponseWriter, r *http.Request) {
		if p, err, m := getPartner(r); err == nil {
			sp.CartApi(w, r, p, m)
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/cart/*", func(w http.ResponseWriter, r *http.Request) {
		if p, err, _ := getPartner(r); err == nil {
			sp.Cart(w, r, p)
		} else {
			handleError(w, err)
		}
	})

	// 购买跳转
	routes.Add("^/buy/*$", func(w http.ResponseWriter, r *http.Request) {
		if p, err, m := getPartner(r); err == nil {
			sp.BuyRedirect(w, r, p, m)
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/buy/confirm$", func(w http.ResponseWriter, r *http.Request) {
		if p, err, m := getPartner(r); err == nil {
			sp.OrderConfirm(w, r, p, m)
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/buy/apply/coupon$", func(w http.ResponseWriter, r *http.Request) {
		if p, err, m := getPartner(r); err == nil {
			if r.Method == "POST" {
				sp.ApplyCoupon_post(w, r, p, m)
			}
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/buy/order/persist", func(w http.ResponseWriter, r *http.Request) {
		if p, err, m := getPartner(r); err == nil {
			if r.Method == "POST" {
				sp.OrderPersist_post(w, r, p, m)
			}
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/buy/order/submit_order$", func(w http.ResponseWriter, r *http.Request) {
		if p, err, m := getPartner(r); err == nil {
			if r.Method == "POST" {
				sp.SubmitOrder_post(w, r, p, m)
			}
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/buy/order/finish$", func(w http.ResponseWriter, r *http.Request) {
		if p, err, m := getPartner(r); err == nil {
			sp.OrderFinish(w, r, p, m)
		} else {
			handleError(w, err)
		}
	})

	routes.Add("/(shopping|buy)/", func(w http.ResponseWriter, r *http.Request) {
		if p, err, m := getPartner(r); err == nil {
			mvc.Handle(sp, w, r, true, p, m)
		} else {
			handleError(w, err)
		}
	})

	routes.Add("^/pay/", func(w http.ResponseWriter, r *http.Request) {
		mvc.Handle(pc, w, r, true)
	})

	routes.Add("/", func(w http.ResponseWriter, r *http.Request) {
		if p, err, m := getPartner(r); err == nil {
			mvc.Handle(mc, w, r, true, p, m)
		} else {
			handleError(w, err)
		}
	})
}
