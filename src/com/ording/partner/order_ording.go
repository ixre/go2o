package partner

import (
	"com/domain/interface/enum"
	"com/domain/interface/shopping"
	"com/ording/cache"
	"com/ording/dproxy"
	"html/template"
	"net/http"
	"strconv"
)

func (this *orderC) setShop(w http.ResponseWriter, r *http.Request,
	partnerId int, order *shopping.ValueOrder) {
	shopDr := cache.GetShopDropList(partnerId, -1)

	isNoShop := len(shopDr) == 0

	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {
			(*m)["shopDr"] = template.HTML(shopDr)
			(*m)["noShop"] = isNoShop
			(*m)["orderNo"] = order.OrderNo

		}, "views/partner/order/order_setup_setshop.html")
}

func (this *orderC) SetShop_post(w http.ResponseWriter, r *http.Request, partnerId int) {
	r.ParseForm()

	shopId, err := strconv.Atoi(r.FormValue("shopId"))
	if err == nil {
		err = dproxy.ShoppingService.SetDeliverShop(partnerId,
			r.FormValue("order_no"), shopId)
	}
	if err != nil {
		w.Write([]byte("{result:false,message:'" + err.Error() + "'}"))
	} else {
		w.Write([]byte("{result:true}"))
	}
}

func (this *orderC) setState(w http.ResponseWriter, r *http.Request,
	partnerId int, order *shopping.ValueOrder) {

	var descript string
	var button string

	switch order.Status {
	case enum.ORDER_COMPLETED:
		descript = `<span style="color:green">订单已经完成！</span>`
	case enum.ORDER_CANCEL:
		descript = `<span style="color:red">订单已经作废！</span>`
	case enum.ORDER_CREATED:
		descript = "确认订单无误后，点击按钮进行下一步.."
		button = `<input type="button" id="btn2" value="确认订单"/>`
	case enum.ORDER_CONFIRMED:
		button = `<input type="button" id="btn2" value="开始处理订单"/>`
	case enum.ORDER_PROCESSING:
		button = `<input type="button" id="btn2" value="开始配送"/>`
	case enum.ORDER_SENDING:
		if order.IsPayed == 1 {
			descript = `<span style="color:red">订单尚未付款,如果已经付款，请手动付款！</span>`
		} else {
			descript = "如果已收货，点击按钮完成订单"
			button = `<input type="button" id="btn2" value="完成订单"/>`
		}
	}

	this.Context.Template().Execute(w,
		func(m *map[string]interface{}) {
			(*m)["button"] = template.HTML(button)
			(*m)["descript"] = template.HTML(descript)
			(*m)["order_no"] = order.OrderNo
		}, "views/partner/order/order_setup_setstate.html")
}
