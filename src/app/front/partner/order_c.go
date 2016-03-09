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
<<<<<<< HEAD
=======
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/mvc"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"go2o/src/cache"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
<<<<<<< HEAD
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
	d := echox.NewRenderData()
	d.Map["shops"] = template.JS(shopsJson)
	return ctx.Render(http.StatusOK, "order.list.html", d)
}

func (this *orderC) WaitPaymentList(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	shopsJson := cache.GetShopsJson(partnerId)

	d := echox.NewRenderData()
	d.Map["shops"] = template.JS(shopsJson)
	return ctx.Render(http.StatusOK, "order.waitpay_list.html", d)
}

func (this *orderC) Cancel(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.cancel_post(ctx)
	}
	d := echox.NewRenderData()
	return ctx.Render(http.StatusOK, "order.cancel.html", d)
}

func (this *orderC) cancel_post(ctx *echox.Context) error {
	result := gof.Message{}
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
=======
	"html/template"
	"strings"
)

var _ mvc.Filter = new(orderC)

type orderC struct {
	*baseC
}

func (this *orderC) List(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	shopsJson := cache.GetShopsJson(partnerId)
	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"shops": template.JS(shopsJson),
		}, "views/partner/order/order_list.html")
}

func (this *orderC) WaitPaymentList(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	shopsJson := cache.GetShopsJson(partnerId)
	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"shops": template.JS(shopsJson),
		}, "views/partner/order/order_waitpay_list.html")
}

func (this *orderC) Cancel(ctx *web.Context) {
	//partnerId := this.GetPartnerId(ctx)
	ctx.App.Template().Execute(ctx.Response, nil, "views/partner/order/cancel.html")

}

func (this *orderC) Cancel_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	r.ParseForm()
	reason := r.FormValue("reason")
	err := dps.ShoppingService.CancelOrder(partnerId,
		r.FormValue("order_no"), reason)

	if err == nil {
<<<<<<< HEAD
		result.Result = true
	} else {
		result.Message = err.Error()
	}
	return ctx.JSON(http.StatusOK, result)
}

func (this *orderC) View(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	r.ParseForm()
	e := dps.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		return ctx.String(http.StatusOK, "无效订单")
=======
		w.Write([]byte("{result:true}"))
	} else {
		w.Write([]byte(`{result:false,message:"` + err.Error() + `"}`))
	}
}

func (this *orderC) View(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	r.ParseForm()
	e := dps.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		w.Write([]byte("无效订单"))
		return
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	}

	member := dps.MemberService.GetMember(e.MemberId)
	if member == nil {
<<<<<<< HEAD
		return ctx.StringOK("无效订单")
=======
		w.Write([]byte("无效订单"))
		return
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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

<<<<<<< HEAD
	d := echox.NewRenderData()
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
	r := ctx.Request()
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
=======
	ctx.App.Template().Execute(w,
		gof.TemplateDataMap{
			"entity":   template.JS(js),
			"member":   member,
			"shopName": shopName,
			"payment":  payment,
			"state":    orderStateText,
		}, "views/partner/order/order_view.html")
}

func (this *orderC) Setup(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	r.ParseForm()
	e := dps.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		w.Write([]byte("无效订单"))
		return
	}

	if e.ShopId == 0 {
		this.setShop(ctx, partnerId, e)
	} else {
		this.setState(ctx, partnerId, e)
	}
}

// 锁定，防止重复下单，返回false,表示正在处理订单
func (this *orderC) lockOrder(ctx *web.Context) bool {
	s := ctx.Session()
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	v := s.Get("pt_order_lock")
	if v != nil {
		return false
	}
	s.Set("pt_order_lock", "1")
	s.Save()
	return true
}
<<<<<<< HEAD
func (this *orderC) releaseOrder(ctx *echox.Context) {
	ctx.Session.Remove("pt_order_lock")
	ctx.Session.Save()
}

// 订单流程(POST)
func (this *orderC) OrderSetup(ctx *echox.Context) error {
	if !this.lockOrder(ctx) {
		return ctx.String(http.StatusOK, "请勿频繁操作")
	}
	var msg gof.Message
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	if r.Method == "POST" {
		r.ParseForm()
		err := dps.ShoppingService.HandleOrder(partnerId, r.FormValue("order_no"))
		this.releaseOrder(ctx)
		if err != nil {
			msg.Message = err.Error()
		} else {
			msg.Result = true
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
	r := ctx.Request()
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

	d := echox.NewRenderData()
	d.Map["shopName"] = shopName
	d.Map["order"] = *e
	return ctx.Render(http.StatusOK, "order.payment.html", d)

}

func (this *orderC) payment_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
=======
func (this *orderC) releaseOrder(ctx *web.Context) {
	ctx.Session().Remove("pt_order_lock")
	ctx.Session().Save()
}

func (this *orderC) OrderSetup_post(ctx *web.Context) {
	if !this.lockOrder(ctx) {
		return
	}

	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	r.ParseForm()
	err := dps.ShoppingService.HandleOrder(partnerId, r.FormValue("order_no"))

	this.releaseOrder(ctx)
	if err != nil {
		w.Write([]byte("{result:false,message:'" + err.Error() + "'}"))
	} else {
		w.Write([]byte("{result:true}"))
	}
}

func (this *orderC) Payment(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	r.ParseForm()
	e := dps.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		w.Write([]byte("无效订单"))
	} else if e.IsPaid == 1 {
		w.Write([]byte("订单已付款"))
	} else {
		var shopName string
		if e.ShopId == 0 {
			shopName = "未指定"
		} else {
			shopName = dps.PartnerService.GetShopValueById(partnerId, e.ShopId).Name
		}

		ctx.App.Template().Execute(w, gof.TemplateDataMap{
			"shopName": shopName,
			"order":    *e,
		}, "views/partner/order/payment.html")
	}
}

func (this *orderC) Payment_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	r.ParseForm()
	orderNo := r.FormValue("orderNo")

	order := dps.ShoppingService.GetOrderByNo(partnerId, orderNo)

	err := dps.MemberService.Charge(partnerId, order.MemberId, member.TypeBalanceSystemCharge, "系统充值", "", order.PayFee)
	if err == nil {
		err = dps.ShoppingService.PayForOrderByManager(partnerId, orderNo)
	}

	if err != nil {
<<<<<<< HEAD
		return ctx.String(http.StatusOK, "{result:false,message:'"+err.Error()+"'}")
	} else {
		return ctx.String(http.StatusOK, "{result:true,message:'付款成功'}")
=======
		w.Write([]byte("{result:false,message:'" + err.Error() + "'}"))
	} else {
		w.Write([]byte("{result:true,message:'付款成功'}"))
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	}
}
