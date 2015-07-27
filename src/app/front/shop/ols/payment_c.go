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
	"errors"
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core/infrastructure/payment"
	"go2o/src/core/service/dps"
	"net/http"
)

var aliPayObj *payment.AliPay = &payment.AliPay{
	Partner: "2088021187655650",
	Key:     "2b0y8466t9mh82s5ajptub2p3wgjzwmh",
	Seller:  "1515827759@qq.com",
}

type PaymentC struct {
	*BaseC
}

func getDomain(r *http.Request) string {
	var proto string
	if r.Proto == "HTTPS" {
		proto = "https://"
	} else {
		proto = "http://"
	}
	return proto + r.Host
}
func (this *PaymentC) Create(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	qs := r.URL.Query()
	partnerId := this.GetPartnerId(ctx)
	orderNo := qs.Get("order_no")
	paymentOpt := qs.Get("pay_opt")

	if len(orderNo) != 0 {
		if paymentOpt == "alipay" {
			domain := getDomain(ctx.Request)
			returnUrl := fmt.Sprintf("%s/pay/return?pay_opt=alipay&partner_id=%d", domain, partnerId)
			notifyUrl := fmt.Sprintf("%s/pay/return?pay_opt=alipay&partner_id=%d", domain, partnerId)
			gateway := aliPayObj.CreateGateway(orderNo, 0.01, "在线支付订单", "订单号："+orderNo, notifyUrl, returnUrl)
			html := "<html><head><meta charset=\"utf-8\"/></head><body>" + gateway + "</body></html>"
			w.Write([]byte(html))
			return
		}
	}

	w.Write([]byte("订单不存在"))
}

func (this *PaymentC) Return(ctx *web.Context) {
	result := aliPayObj.Return(ctx.Request)
	if result.Status == payment.StatusTradeSuccess {
		if err := this.handleOrder(ctx, &result); err == nil {
			this.paymentSuccess(ctx, &result)
			return
		}
	}
	this.paymentFail(ctx, &result)
}

func (this *PaymentC) Notify_post(ctx *web.Context) {
	result := aliPayObj.Notify(ctx.Request)
	if result.Status == payment.StatusTradeSuccess {
		if err := this.handleOrder(ctx, &result); err == nil {
			ctx.Response.Write([]byte("success"))
			return
		}
	}
	ctx.Response.Write([]byte("fail"))
}

func (this *PaymentC) handleOrder(ctx *web.Context, result *payment.Result) error {
	partnerId := this.GetPartnerId(ctx)
	order := dps.ShoppingService.GetOrderByNo(partnerId, result.OrderNo)
	if order != nil {
		//		if order.PayFee != result.Fee{
		//			return errors.New(fmt.Sprintf("error order fee %f / %f",order.PayFee,result.Fee))
		//		}
		if order.IsPaid == 1 {
			return errors.New("order has paid")
		}
		return dps.ShoppingService.PayForOrder(partnerId, order.OrderNo)
	}
	return errors.New("no such order")
}

func (this *PaymentC) paymentSuccess(ctx *web.Context, result *payment.Result) {
	p := this.GetPartner(ctx)
	siteConf := this.GetSiteConf(ctx)

	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"partner": p,
		"conf":    siteConf,
	},
		"views/shop/ols/{device}/payment_success.html",
		"views/shop/ols/{device}/inc/header.html",
		"views/shop/ols/{device}/inc/footer.html")
}

func (this *PaymentC) paymentFail(ctx *web.Context, result *payment.Result) {
	p := this.GetPartner(ctx)
	siteConf := this.GetSiteConf(ctx)

	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"partner": p,
		"conf":    siteConf,
	},
		"views/shop/ols/{device}/payment_fail.html",
		"views/shop/ols/{device}/inc/header.html",
		"views/shop/ols/{device}/inc/footer.html")
}
