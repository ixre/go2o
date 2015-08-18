/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-11 21:33
 * description :
 * history :
 */
package ols

import (
	"github.com/jrsix/gof/web"
	"go2o/src/core/infrastructure/alipay"
)

type paymentC struct {
}

func (this *paymentC) Create(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
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

func (this *paymentC) Return(ctx *web.Context) {
	r, _ := ctx.Request, ctx.Response
	alipay.ReturnFunc(r, nil)
	ctx.Response.Write([]byte("支付完成"))
}

func (this *paymentC) Notify_post(ctx *web.Context) {
	r, _ := ctx.Request, ctx.Response
	alipay.NotifyFunc(r, nil)
}
