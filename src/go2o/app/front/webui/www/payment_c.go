/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-11 21:33
 * description :
 * history :
 */
package www

import (
	"github.com/atnet/gof/app"
	"go2o/core/infrastructure/alipay"
	"net/http"
)

type paymentC struct {
	app.Context
}

func (this *paymentC) Create(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	orderNo := qs.Get("order_no")
	paymentOpt := qs.Get("pay_opt")

	if len(orderNo) != 0 {
		if paymentOpt == "alipay" {
			html := "<html><body>" + alipay.CreatePaymentGateWay(orderNo, 0.01, "订单"+orderNo, "", "", "") + "</body></html>"
			w.Write([]byte(html))
			return
		}
	}

	w.Write([]byte("订单不存在"))
}

func (this *paymentC) Return(w http.ResponseWriter, r *http.Request) {
	alipay.ReturnFunc(r, nil)
}

func (this *paymentC) Notify_post(w http.ResponseWriter, r *http.Request) {
	alipay.NotifyFunc(r, nil)
}
