/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
<<<<<<< HEAD
=======
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"go2o/src/cache"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/service/dps"
<<<<<<< HEAD
	"go2o/src/x/echox"
=======
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"html/template"
	"strconv"
)

<<<<<<< HEAD
func (this *orderC) setShop(ctx *echox.Context,
	partnerId int, order *shopping.ValueOrder) error {
=======
func (this *orderC) setShop(ctx *web.Context,
	partnerId int, order *shopping.ValueOrder) {
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	shopDr := cache.GetShopDropList(partnerId, -1)

	isNoShop := len(shopDr) == 0

<<<<<<< HEAD
	d := echox.NewRenderData()
	d.Map = map[string]interface{}{
		"shopDr":  template.HTML(shopDr),
		"noShop":  isNoShop,
		"orderNo": order.OrderNo,
	}
	return ctx.RenderOK("order.setup_setshop.html", d)
}

// 设置门店(POST)
func (this *orderC) SetShop(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	if r.Method == "POST" {
		r.ParseForm()

		shopId, err := strconv.Atoi(r.FormValue("shopId"))
		if err == nil {
			orderNo := r.FormValue("order_no")
			err = dps.ShoppingService.SetDeliverShop(partnerId,
				orderNo, shopId)
			dps.ShoppingService.ConfirmOrder(partnerId, orderNo)
		}
		if err != nil {
			return ctx.StringOK("{result:false,message:'" + err.Error() + "'}")
		}
		return ctx.StringOK("{result:true}")
	}
	return nil
}

func (this *orderC) setState(ctx *echox.Context,
	partnerId int, order *shopping.ValueOrder) error {
=======
	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"shopDr":  template.HTML(shopDr),
			"noShop":  isNoShop,
			"orderNo": order.OrderNo,
		}, "views/partner/order/order_setup_setshop.html")
}

func (this *orderC) SetShop_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	r.ParseForm()

	shopId, err := strconv.Atoi(r.FormValue("shopId"))
	if err == nil {
		err = dps.ShoppingService.SetDeliverShop(partnerId,
			r.FormValue("order_no"), shopId)
	}
	if err != nil {
		w.Write([]byte("{result:false,message:'" + err.Error() + "'}"))
	} else {
		w.Write([]byte("{result:true}"))
	}
}

func (this *orderC) setState(ctx *web.Context,
	partnerId int, order *shopping.ValueOrder) {
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d

	var descript string
	var button string

	switch order.Status {
	case enum.ORDER_COMPLETED:
		descript = `<span style="color:green">订单已经完成！</span>`
	case enum.ORDER_CANCEL:
		descript = `<span style="color:red">订单已经作废！</span>`
	case enum.ORDER_WAIT_CONFIRM:
		descript = "确认订单无误后，点击按钮进行下一步.."
		button = `<span class="ui-button w80 small-button">
                <span class="button-inner">
                    <span class="button-txt">确认订单</span>
                    <input id="btn2"/>
                </span>
            </span>`
	case enum.ORDER_WAIT_DELIVERY:
		button = `<span class="ui-button w80 small-button">
				<span class="button-inner">
				<span class="button-txt">开始配送</span>
				<input id="btn2"/>
				</span>
				</span>`
	case enum.ORDER_WAIT_RECEIVE:
		button = `<span class="ui-button w80 small-button">
					<span class="button-inner">
					<span class="button-txt">确认收货</span>
					<input id="btn2"/>
					</span>
					</span>`
	case enum.ORDER_RECEIVED:
		if order.IsPaid == 0 {
			descript = `<span style="color:red">订单尚未付款,如果已经付款，请人工手动付款！</span>`
		} else {
			descript = "如果已收货，点击按钮完成订单"
			button = `<span class="ui-button w80 small-button">
                <span class="button-inner">
                    <span class="button-txt">完成订单</span>
                    <input id="btn2"/>
                </span>
            </span>`
		}
	}

<<<<<<< HEAD
	d := echox.NewRenderData()
	d.Map = map[string]interface{}{

		"button":   template.HTML(button),
		"descript": template.HTML(descript),
		"order_no": order.OrderNo,
	}
	return ctx.RenderOK("order.setup_setstate.html", d)
=======
	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"button":   template.HTML(button),
			"descript": template.HTML(descript),
			"order_no": order.OrderNo,
		}, "views/partner/order/order_setup_setstate.html")
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
