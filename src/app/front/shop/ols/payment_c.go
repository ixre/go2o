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
	"strconv"
	"strings"
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
			returnUrl := fmt.Sprintf("%s/pay/return_alipay", domain)
			notifyUrl := fmt.Sprintf("%s/pay/notify/%d_alipay", domain, partnerId)
			gateway := aliPayObj.CreateGateway(orderNo, 0.01, "在线支付订单", "订单号："+orderNo, notifyUrl, returnUrl)
			html := "<html><head><meta charset=\"utf-8\"/></head><body>" + gateway + "</body></html>"
			w.Write([]byte(html))

			payment.Debug(" [ Submit] - %s - %s", orderNo, notifyUrl)

			return
		}
	}

	w.Write([]byte("订单不存在"))
}

func (this *PaymentC) Return_alipay(ctx *web.Context) {
	result := aliPayObj.Return(ctx.Request)
	if result.Status == payment.StatusTradeSuccess {
		if err := this.handleOrder(this.GetPartnerId(ctx), "alipay", &result); err == nil {
			this.paymentSuccess(ctx, &result)
			return
		}
	}
	this.paymentFail(ctx, &result)
}

func (this *PaymentC) Notify_post(ctx *web.Context) {
	path := ctx.Request.URL.Path
	lastSeg := strings.Split(path[strings.LastIndex(path, "/")+1:], "_")
	paymentOpt := lastSeg[1]
	partnerId, _ := strconv.Atoi(lastSeg[0])

	payment.Debug(" [ Notify] - URL - %s - %d -  %s", ctx.Request.RequestURI, partnerId, paymentOpt)

	if paymentOpt == "alipay" {
		result := aliPayObj.Notify(ctx.Request)
		if result.Status == payment.StatusTradeSuccess {
			err := this.handleOrder(partnerId, "alipay", &result)
			if err == nil {
				payment.Debug("payment ok")
				ctx.Response.Write([]byte("success"))
				return
			}
			payment.Debug(" payment fail, %s", err.Error())
		}
	}
	ctx.Response.Write([]byte("fail"))
}

func (this *PaymentC) handleOrder(partnerId int, sp string, result *payment.Result) error {
	order := dps.ShoppingService.GetOrderByNo(partnerId, result.OrderNo)
	if order != nil {
		//		if order.PayFee != result.Fee{
		//			return errors.New(fmt.Sprintf("error order fee %f / %f",order.PayFee,result.Fee))
		//		}
		if order.IsPaid == 1 {
			return errors.New("order has paid")
		}
		return dps.ShoppingService.PayForOrderOnlineTrade(partnerId, order.OrderNo, sp, result.TradeNo)
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
