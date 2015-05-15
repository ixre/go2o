/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:55
 * description :
 * history :
 */

package partner

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/service/dps"
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
	ctx.App.Template().Execute(ctx.ResponseWriter,
		gof.TemplateDataMap{
			"shops": template.JS(shopsJson),
		}, "views/partner/order/order_list.html")
}

func (this *orderC) Cancel(ctx *web.Context) {
	//partnerId := this.GetPartnerId(ctx)
	ctx.App.Template().Execute(ctx.ResponseWriter, nil, "views/partner/order/cancel.html")

}

func (this *orderC) Cancel_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	reason := r.FormValue("reason")
	err := dps.ShoppingService.CancelOrder(partnerId,
		r.FormValue("order_no"), reason)

	if err == nil {
		w.Write([]byte("{result:true}"))
	} else {
		w.Write([]byte(`{result:false,message:"` + err.Error() + `"}`))
	}
}

func (this *orderC) View(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	e := dps.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		w.Write([]byte("无效订单"))
		return
	}

	member := dps.MemberService.GetMember(e.MemberId)
	if member == nil {
		w.Write([]byte("无效订单"))
		return
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
	r, w := ctx.Request, ctx.ResponseWriter
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

func (this *orderC) OrderSetup_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	err := dps.ShoppingService.HandleOrder(partnerId, r.FormValue("order_no"))
	if err != nil {
		w.Write([]byte("{result:false,message:'" + err.Error() + "'}"))
	} else {
		w.Write([]byte("{result:true}"))
	}
}

func (this *orderC) Payment(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
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
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	orderNo := r.FormValue("orderNo")

	err := dps.ShoppingService.PayForOrder(partnerId, orderNo)
	if err != nil {
		w.Write([]byte("{result:false,message:'" + err.Error() + "'}"))
	} else {
		w.Write([]byte("{result:true,message:'付款成功'}"))
	}
}
