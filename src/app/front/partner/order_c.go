/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:55
 * description :
 * history :
 */

package partner

import (
	"encoding/json"
	"github.com/jsix/gof"
	"go2o/src/app/cache"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strings"
)

type orderC struct {
}

func (this *orderC) List(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	shopsJson := cache.GetShopsJson(partnerId)
	d := ctx.NewData()
	d.Map["shops"] = template.JS(shopsJson)
	return ctx.Render(http.StatusOK, "order.list.html", d)
}

func (this *orderC) WaitPaymentList(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	shopsJson := cache.GetShopsJson(partnerId)

	d := ctx.NewData()
	d.Map["shops"] = template.JS(shopsJson)
	return ctx.Render(http.StatusOK, "order.waitpay_list.html", d)
}

func (this *orderC) Cancel(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.cancel_post(ctx)
	}
	d := ctx.NewData()
	return ctx.Render(http.StatusOK, "order.cancel.html", d)
}

func (this *orderC) cancel_post(ctx *echox.Context) error {
	result := gof.Result{}
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	r.ParseForm()
	reason := r.FormValue("reason")
	err := dps.ShoppingService.CancelOrder(partnerId,
		r.FormValue("order_no"), reason)

	if err == nil {
		result.ErrCode = 0
	} else {
		result.ErrMsg = err.Error()
	}
	return ctx.JSON(http.StatusOK, result)
}

func (this *orderC) View(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	r.ParseForm()
	e := dps.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		return ctx.String(http.StatusOK, "无效订单")
	}

	member := dps.MemberService.GetMember(e.MemberId)
	if member == nil {
		return ctx.StringOK("无效订单")
	}

	e.ItemsInfo = strings.Replace(e.ItemsInfo, "\n", "<br />", -1)
	if len(e.Note) == 0 {
		e.Note = "无备注"
	}

	js, _ := json.Marshal(e)

	var shopName string
	var payment string
	var orderStateText string
	if e.ShopId == 0 {
		shopName = "未指定"
	} else {
		shopName = dps.PartnerService.GetShopValueById(partnerId, e.ShopId).Name
	}
	payment = enum.GetPaymentName(e.PaymentOpt)
	orderStateText = enum.OrderState(e.Status).String()

	d := ctx.NewData()
	d.Map = map[string]interface{}{
		"entity":   template.JS(js),
		"member":   member,
		"shopName": shopName,
		"payment":  payment,
		"state":    orderStateText,
	}
	return ctx.Render(http.StatusOK, "order.view.html", d)
}

func (this *orderC) Setup(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	e := dps.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		return ctx.String(http.StatusOK, "无效订单")
	}

	if e.ShopId == 0 {
		return this.setShop(ctx, partnerId, e)
	}
	return this.setState(ctx, partnerId, e)
}

// 锁定，防止重复下单，返回false,表示正在处理订单
func (this *orderC) lockOrder(ctx *echox.Context) bool {
	s := ctx.Session
	v := s.Get("pt_order_lock")
	if v != nil {
		return false
	}
	s.Set("pt_order_lock", "1")
	s.Save()
	return true
}
func (this *orderC) releaseOrder(ctx *echox.Context) {
	ctx.Session.Remove("pt_order_lock")
	ctx.Session.Save()
}

// 订单流程(POST)
func (this *orderC) OrderSetup(ctx *echox.Context) error {
	if !this.lockOrder(ctx) {
		return ctx.String(http.StatusOK, "请勿频繁操作")
	}
	var msg gof.Result
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		r.ParseForm()
		err := dps.ShoppingService.HandleOrder(partnerId, r.FormValue("order_no"))
		this.releaseOrder(ctx)
		if err != nil {
			msg.ErrMsg = err.Error()
		} else {
	msg.ErrCode = 0
		}
		return ctx.JSON(http.StatusOK, msg)
	}
	return nil
}

func (this *orderC) Payment(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.payment_post(ctx)
	}
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	r.ParseForm()
	e := dps.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		return ctx.String(http.StatusOK, "无效订单")
	} else if e.IsPaid == 1 {
		return ctx.String(http.StatusOK, "订单已付款")
	}

	var shopName string
	if e.ShopId == 0 {
		shopName = "未指定"
	} else {
		shopName = dps.PartnerService.GetShopValueById(partnerId, e.ShopId).Name
	}

	d := ctx.NewData()
	d.Map["shopName"] = shopName
	d.Map["order"] = *e
	return ctx.Render(http.StatusOK, "order.payment.html", d)

}

func (this *orderC) payment_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	r.ParseForm()
	orderNo := r.FormValue("orderNo")

	order := dps.ShoppingService.GetOrderByNo(partnerId, orderNo)

	err := dps.MemberService.Charge(partnerId, order.MemberId, member.TypeBalanceSystemCharge, "系统充值(订单付款)", "", order.PayFee)
	if err == nil {
		err = dps.ShoppingService.PayForOrderByManager(partnerId, orderNo)
	}

	if err != nil {
		return ctx.String(http.StatusOK, "{result:false,message:'"+err.Error()+"'}")
	} else {
		return ctx.String(http.StatusOK, "{result:true,message:'付款成功'}")
	}
}
