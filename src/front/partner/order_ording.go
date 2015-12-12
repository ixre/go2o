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
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/service/dps"
	"html/template"
	"strconv"
)

func (this *orderC) setShop(ctx *web.Context,
	partnerId int, order *shopping.ValueOrder) {
	shopDr := cache.GetShopDropList(partnerId, -1)

	isNoShop := len(shopDr) == 0

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"shopDr":  template.HTML(shopDr),
			"noShop":  isNoShop,
			"orderNo": order.OrderNo,
		}, "views/partner/order/order_setup_setshop.html")
}

func (this *orderC) SetShop_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
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

	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"button":   template.HTML(button),
			"descript": template.HTML(descript),
			"order_no": order.OrderNo,
		}, "views/partner/order/order_setup_setstate.html")
}
