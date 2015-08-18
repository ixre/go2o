/**
 * Copyright 2014 @ z3q.net.
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
	"github.com/jrsix/gof"
	guitl "github.com/jrsix/gof/util"
	"github.com/jrsix/gof/web"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/infrastructure/payment"
	"go2o/src/core/service/dps"
	"net/http"
	"strconv"
	"strings"
)

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

func (this *PaymentC) getAliPayment(ctx *web.Context) payment.IPayment {
	var p payment.IPayment

	if guitl.IsMobileAgent(ctx.Request.UserAgent()) {
		p = &payment.AliPayWap{
			Partner: "2088021187655650",
			Key:     "3aijnz4020um0c7iq0ayanaqqcxtxk5i",
			Seller:  "1515827759@qq.com",
			PrivateKey: `-----BEGIN ENCRYPTED PRIVATE KEY-----
	MIICoTAbBgkqhkiG9w0BBQMwDgQIu5+aUzZxXVUCAggABIICgMJ/Hc9wnG//S/Rg
	7f4ReeAP0cybL3dO5MUHUnQDa6gGebSPMmB/rAXjtMmTjzyYCjhbJyRkSHKNWLIx
	zZdHMv4T7A/cHK0owRKSNnW7AzY4seCBbnsLteVo33PBAF5u+tO9O5maBm2Rv9xi
	3gtSH2gh5RoGrjF85VsRm0Vn/e/4Q3dB/IN5YZt2W76GJpPpLk9ltgmJYFcJ0c3X
	LAqK08RQGN7TfHptYFydrtBftFI6kj1jn7Qs3h7Uqc1qpnDroWqSio7IWE6cF8Nx
	XD94xuIPLBVlRjGPIZq1PwaIO0cfcAkcD8JZqVMCn99c9x5MW0DFzNjotithZB2v
	ApooLlhYqoLOrPGpUW8aOnaJ15/awMsJYtyvjF4/IkY6Q1xVqwCTKnNq9aMlmKZU
	W+8gnJxpqVRNCUC6fuJhLU2fPD85RWfHoAq8iNxz1nz8KHiVVh3FSwS1RyxV6amH
	ozar9aGZPlh1zT649h51YSpLy/q2pJwfl78a97ArAqXCCltLF/oMDqwcs4BqM9qP
	PUcSt0k6mURvLwBe2ztop4xTFONn5DizAvEmdTO1YHOQlqDXbxSfO9gH7Yj5fmoL
	AdebjiSZfR//1dvePyM8wkk67PdWItxuNGKg7TeZCxfsGkYsq4t38rRNHmSvevV0
	c9XWpbqupJy/g8OsP1Afj4F+9W3wBkhiMFvidIvJcTnkvmxJGz+dJb/feBr10Il+
	+CVucgZdPkMQoREe+FDV3G3K1ZaoGLvbZwcUBsyF0X/l3TIjgjQxuW8j+1NMkstF
	TsARjljf7udXaOCK7Uf6vujC2Zk3UI/39LSJ12WAB/Wgc1TEhBq/e6hnyGXtbuBf
	Ibg2Kl0=
	-----END ENCRYPTED PRIVATE KEY-----`,
		}
	} else {
		p = &payment.AliPay{
			Partner: "2088021187655650",
			Key:     "3aijnz4020um0c7iq0ayanaqqcxtxk5i",
			Seller:  "1515827759@qq.com",
		}
	}
	return p
}
func (this *PaymentC) Create(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
	qs := r.URL.Query()
	partnerId := this.GetPartnerId(ctx)
	orderNo := qs.Get("order_no")
	paymentOpt := qs.Get("pay_opt")

	if len(orderNo) != 0 {
		ctx.Session().Set("current_payment", orderNo)
		ctx.Session().Save()

		if paymentOpt == "alipay" {
			aliPayObj := this.getAliPayment(ctx)
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
	//this.paymentFail(ctx,nil)
	//return
	aliPayObj := this.getAliPayment(ctx)
	result := aliPayObj.Return(ctx.Request)
	partnerId := this.GetPartnerId(ctx)
	if len(result.OrderNo) == 0 {
		result.OrderNo = ctx.Session().Get("current_payment").(string)
	}
	order := dps.ShoppingService.GetOrderByNo(partnerId, result.OrderNo)
	if result.Status == payment.StatusTradeSuccess {
		this.handleOrder(order, "alipay", &result)
		this.paymentSuccess(ctx, order, &result)
		return
	}

	this.paymentFail(ctx, order, &result)
}

func (this *PaymentC) Notify_post(ctx *web.Context) {
	path := ctx.Request.URL.Path
	lastSeg := strings.Split(path[strings.LastIndex(path, "/")+1:], "_")
	paymentOpt := lastSeg[1]
	partnerId, _ := strconv.Atoi(lastSeg[0])
	payment.Debug(" [ Notify] - URL - %s - %d -  %s", ctx.Request.RequestURI, partnerId, paymentOpt)

	if paymentOpt == "alipay" {
		aliPayObj := this.getAliPayment(ctx)
		result := aliPayObj.Notify(ctx.Request)
		order := dps.ShoppingService.GetOrderByNo(partnerId, result.OrderNo)
		if result.Status == payment.StatusTradeSuccess {
			this.handleOrder(order, "alipay", &result)
			payment.Debug("payment ok")
			ctx.Response.Write([]byte("success"))
			return
		}
	}
	ctx.Response.Write([]byte("fail"))
}

func (this *PaymentC) handleOrder(order *shopping.ValueOrder, sp string, result *payment.Result) error {
	if order != nil {
		//		if order.PayFee != result.Fee{
		//			return errors.New(fmt.Sprintf("error order fee %f / %f",order.PayFee,result.Fee))
		//		}
		if order.IsPaid == 1 {
			return errors.New("order has paid")
		}
		return dps.ShoppingService.PayForOrderOnlineTrade(order.PartnerId, order.OrderNo, sp, result.TradeNo)
	}
	return errors.New("no such order")
}

func (this *PaymentC) paymentSuccess(ctx *web.Context,
	order *shopping.ValueOrder, result *payment.Result) {
	p := this.GetPartner(ctx)
	siteConf := this.GetSiteConf(ctx)

	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"partner": p,
		"conf":    siteConf,
		"order":   order,
	},
		"views/shop/ols/{device}/payment_success.html",
		"views/shop/ols/{device}/inc/header.html",
		"views/shop/ols/{device}/inc/footer.html")
}

func (this *PaymentC) paymentFail(ctx *web.Context,
	order *shopping.ValueOrder, result *payment.Result) {
	p := this.GetPartner(ctx)
	siteConf := this.GetSiteConf(ctx)

	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{
		"partner": p,
		"conf":    siteConf,
		"order":   order,
	},
		"views/shop/ols/{device}/payment_fail.html",
		"views/shop/ols/{device}/inc/header.html",
		"views/shop/ols/{device}/inc/footer.html")
}
