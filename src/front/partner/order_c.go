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
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/mvc"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
	"html/template"
	"strings"
	"go2o/src/x/echox"
	"net/http"
)

type orderC struct {
}

func (this *orderC) List(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	shopsJson := cache.GetShopsJson(partnerId)
	d := echox.NewRenderData()
	d.Map["shops"] = template.JS(shopsJson)
	return ctx.Render(http.StatusOK,"order/order_list.html",d)
}

func (this *orderC) WaitPaymentList(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	shopsJson := cache.GetShopsJson(partnerId)

	d := echox.NewRenderData()
	d.Map["shops"] = template.JS(shopsJson)
	return ctx.Render(http.StatusOK,"order/order_waitpay_list.html",d)
}

func (this *orderC) Cancel(ctx *echox.Context) error {
	d := echox.NewRenderData()
	return ctx.Render(http.StatusOK,"order/cancel.html",d)
}

func (this *orderC) Cancel_post(ctx *echox.Context) error {
	result := gof.Message{}
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	r.ParseForm()
	reason := r.FormValue("reason")
	err := dps.ShoppingService.CancelOrder(partnerId,
		r.FormValue("order_no"), reason)

	if err == nil {
		result.Result =true
	} else {
		result.Message = err.Error()
	}
	return ctx.JSON(http.StatusOK,result)
}

func (this *orderC) View(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	r.ParseForm()
	e := dps.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		return ctx.String(http.StatusOK,"无效订单")
	}

	member := dps.MemberService.GetMember(e.MemberId)
	if member == nil {
		return ctx.String(http.StatusOK,"无效订单")
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

	d := echox.NewRenderData()
	d.Map = map[string]interface{}{
		"entity":   template.JS(js),
		"member":   member,
		"shopName": shopName,
		"payment":  payment,
		"state":    orderStateText,
	}
	return ctx.Render(http.StatusOK,"order/order_view.html",d)
}

func (this *orderC) Setup(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	r.ParseForm()
	e := dps.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		return ctx.String(http.StatusOK,"无效订单")
	}

	if e.ShopId == 0 {
		return this.setShop(ctx, partnerId, e)
	}
		return this.setState(ctx, partnerId, e)
}

// 锁定，防止重复下单，返回false,表示正在处理订单
func (this *orderC) lockOrder(ctx *echox.Context) bool {
	s := ctx.Session()
	v := s.Get("pt_order_lock")
	if v != nil {
		return false
	}
	s.Set("pt_order_lock", "1")
	s.Save()
	return true
}
func (this *orderC) releaseOrder(ctx *echox.Context) {
	ctx.Session().Remove("pt_order_lock")
	ctx.Session().Save()
}

func (this *orderC) OrderSetup_post(ctx *echox.Context) error {
	if !this.lockOrder(ctx) {
		return ctx.String(http.StatusOK,"请勿频繁操作")
	}

	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	r.ParseForm()
	err := dps.ShoppingService.HandleOrder(partnerId, r.FormValue("order_no"))

	this.releaseOrder(ctx)
	if err != nil {
		return ctx.String(http.StatusOK,"{result:false,message:'" + err.Error() + "'}"))
	}
		return ctx.String(http.StatusOK,"{result:true}")

}

func (this *orderC) Payment(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	r.ParseForm()
	e := dps.ShoppingService.GetOrderByNo(partnerId, r.FormValue("order_no"))
	if e == nil {
		return ctx.String(http.StatusOK,"无效订单")
	} else if e.IsPaid == 1 {
		return ctx.String(http.StatusOK,"订单已付款")
	}

		var shopName string
		if e.ShopId == 0 {
			shopName = "未指定"
		} else {
			shopName = dps.PartnerService.GetShopValueById(partnerId, e.ShopId).Name
		}

	d:= echox.NewRenderData()
	d.Map["shopName"] = shopName
	d.Map["order"] = *e
	return ctx.Render(http.StatusOK,"order/payment.html",d)

}

func (this *orderC) Payment_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	r.ParseForm()
	orderNo := r.FormValue("orderNo")

	order := dps.ShoppingService.GetOrderByNo(partnerId, orderNo)

	err := dps.MemberService.Charge(partnerId, order.MemberId, member.TypeBalanceSystemCharge, "系统充值", "", order.PayFee)
	if err == nil {
		err = dps.ShoppingService.PayForOrderByManager(partnerId, orderNo)
	}

	if err != nil {
		return ctx.String(http.StatusOK,"{result:false,message:'" + err.Error() + "'}")
	} else {
		return ctx.String(http.StatusOK,"{result:true,message:'付款成功'}")
	}
}
