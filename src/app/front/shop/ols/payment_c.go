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
	"github.com/jsix/gof"
	guitl "github.com/jsix/gof/util"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/infrastructure/payment"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"net/http"
	"strconv"
	"strings"
)

type PaymentC struct {
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

func (this *PaymentC) getAliPayment(ctx *echox.Context) payment.IPayment {
	var p payment.IPayment

	if guitl.IsMobileAgent(ctx.Request().UserAgent()) {
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
func (this *PaymentC) Create(ctx *echox.Context) error {
	r, w := ctx.Request(), ctx.Response()
	qs := r.URL.Query()
	partnerId := GetSessionPartnerId(ctx)
	orderNo := qs.Get("order_no")
	paymentOpt := qs.Get("pay_opt")

	var order *shopping.ValueOrder
	if len(orderNo) > 0 {
		//	dps.ShoppingService.PayForOrderOnlineTrade(partnerId,orderNo,"alipay","")
		//	ctx.Response.Header().Add("Location", fmt.Sprintf("/buy/payment?order_no=%s", orderNo))
		//	ctx.Response.WriteHeader(302)
		//	return
		order = dps.ShoppingService.GetOrderByNo(partnerId, orderNo)
	}

	if order != nil {
		if order.IsPaid == enum.TRUE {
			w.Header().Add("Location", fmt.Sprintf("/buy/payment?order_no=%s", order.OrderNo))
			w.WriteHeader(302)
			return nil
		}
		ctx.Session.Set("current_payment", orderNo)
		ctx.Session.Save()

		//order.PayFee = 0.01

		if paymentOpt == "alipay" || paymentOpt == strconv.Itoa(enum.PaymentOnlinePay) {
			aliPayObj := this.getAliPayment(ctx)
			domain := getDomain(r)
			returnUrl := fmt.Sprintf("%s/pay/return_alipay", domain)
			notifyUrl := fmt.Sprintf("%s/pay/notify/%d_alipay", domain, partnerId)
			if len(order.Subject) == 0 {
				order.Subject = "在线支付订单"
			}
			gateway := aliPayObj.CreateGateway(orderNo, order.PayFee, order.Subject, "订单号："+orderNo, notifyUrl, returnUrl)
			html := "<html><head><meta charset=\"utf-8\"/></head><body>" + gateway + "</body></html>"
			w.Write([]byte(html))

			payment.Debug(" [ Submit] - %s - %s", orderNo, notifyUrl)
			return nil
		}
	}

	w.Write([]byte("订单不存在"))
	return nil
}

func (this *PaymentC) Return_alipay(ctx *echox.Context) error {
	//this.paymentFail(ctx,nil)
	//return
	aliPayObj := this.getAliPayment(ctx)
	result := aliPayObj.Return(ctx.Request())
	partnerId := GetSessionPartnerId(ctx)
	if len(result.OrderNo) == 0 {
		result.OrderNo = ctx.Session.Get("current_payment").(string)
	}
	order := dps.ShoppingService.GetOrderByNo(partnerId, result.OrderNo)
	if result.Status == payment.StatusTradeSuccess {
		this.handleOrder(order, "alipay", &result)
		return this.paymentSuccess(ctx, order, &result)
	}

	return this.paymentFail(ctx, order, &result)
}

func (this *PaymentC) Notify_post(ctx *echox.Context) error {
	r := ctx.Request()
	path := r.URL.Path
	lastSeg := strings.Split(path[strings.LastIndex(path, "/")+1:], "_")
	paymentOpt := lastSeg[1]
	partnerId, _ := strconv.Atoi(lastSeg[0])
	payment.Debug(" [ Notify] - URL - %s - %d -  %s", r.RequestURI, partnerId, paymentOpt)

	if paymentOpt == "alipay" {
		aliPayObj := this.getAliPayment(ctx)
		result := aliPayObj.Notify(ctx.Request())
		order := dps.ShoppingService.GetOrderByNo(partnerId, result.OrderNo)
		if result.Status == payment.StatusTradeSuccess {
			this.handleOrder(order, "alipay", &result)
			payment.Debug("payment ok")
			return ctx.StringOK("success")
		}
	}
	return ctx.StringOK("fail")
}

func (this *PaymentC) handleOrder(order *shopping.ValueOrder, sp string, result *payment.Result) error {
	if order != nil {
		//		if order.PayFee != result.Fee{
		//			return errors.New(fmt.Sprintf("error order fee %f / %f",order.PayFee,result.Fee))
		//		}
		if order.IsPaid == 1 {
			return errors.New("order has paid")
		}
		return dps.ShoppingService.PayForOrderOnlineTrade(order.PartnerId,
			order.OrderNo, sp, result.TradeNo)
	}
	return errors.New("no such order")
}

func (this *PaymentC) paymentSuccess(ctx *echox.Context,
	order *shopping.ValueOrder, result *payment.Result) error {
	p := getPartner(ctx)
	siteConf := getSiteConf(ctx)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner": p,
		"conf":    siteConf,
		"order":   order,
	}
	return ctx.RenderOK("payment_success.html", d)
}

func (this *PaymentC) paymentFail(ctx *echox.Context,
	order *shopping.ValueOrder, result *payment.Result) error {
	p := getPartner(ctx)
	siteConf := getSiteConf(ctx)
	d := ctx.NewData()
	d.Map = gof.TemplateDataMap{
		"partner": p,
		"conf":    siteConf,
		"order":   order,
	}
	return ctx.RenderOK("payment_fail.html", d)
}
