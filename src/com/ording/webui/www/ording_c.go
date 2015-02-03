/**
 * Copyright 2014 @ Ops.
 * name :
 * author : newmin
 * date : 2013-11-05 17:08
 * description :
 * history :
 */
package www

import (
	"bytes"
	"com/domain/interface/enum"
	"com/domain/interface/member"
	"com/domain/interface/partner"
	"com/infrastructure/format"
	"com/ording"
	"com/ording/cache/apicache"
	"com/ording/entity"
	"com/service/goclient"
	"com/share/variable"
	"fmt"
	"html/template"
	"net/http"
	"github.com/atnet/gof/app"
	"strconv"
	"strings"
	"time"
)

type ordingC struct {
	app.Context
}

func (this *ordingC) Index(w http.ResponseWriter, r *http.Request, p *entity.Partner, mm *member.ValueMember) {
	if b, siteConf := GetSiteConf(w, p); b {
		categories := apicache.GetCategories(this.Context, p.Id, p.Secret)
		this.Context.Template().Execute(w, func(m *map[string]interface{}) {
			(*m)["partner"] = p
			(*m)["title"] = "在线订餐-" + p.Name
			(*m)["categories"] = template.HTML(categories)
			(*m)["member"] = mm
			(*m)["conf"] = siteConf
		},
			"views/web/www/ording.html",
			"views/web/www/inc/header.html",
			"views/web/www/inc/footer.html")
	}
}

func (this *ordingC) List(w http.ResponseWriter, r *http.Request, p *entity.Partner) {
	const getNum int = -1 //-1表示全部
	categoryId, err := strconv.Atoi(r.URL.Query().Get("cid"))
	if err != nil {
		w.Write([]byte(`{"error":"yes"}`))
		return
	}
	items, err := goclient.Partner.GetItems(p.Id, p.Secret, categoryId, getNum)
	if err != nil {
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	buf := bytes.NewBufferString("<ul>")

	noPicPath := this.Context.Config().GetString(variable.NoPicPath)

	for _, v := range items {

		if len(v.Image) == 0 {
			v.Image = noPicPath
		}

		buf.WriteString(fmt.Sprintf(`
			<li>
				<div class="it_items">
                        <img src="%s/%s" alt="%s"/>
                        <h3 class="name">%s<em>%s</em></h3>
                        <span class="srice">原价:￥%s</span>
                        <span class="sprice">优惠价:￥%s</span>
                        <a href="javascript:cart.add({'id':'%d','name':'%s','price':%s});" class="add">&nbsp;</a>
                </div>
             </li>
		`, this.Context.Config().GetString(variable.ImageServer),
			v.Image, v.Name, v.Name, v.Note, format.FormatFloat(v.Price),
			format.FormatFloat(v.SalePrice),
			v.Id, v.Name, format.FormatFloat(v.SalePrice)))
	}
	buf.WriteString("</ul>")
	w.Write(buf.Bytes())
}

func (this *ordingC) Order(w http.ResponseWriter, r *http.Request,
	p *entity.Partner, mm *member.ValueMember) {
	//cart=%u91CE%u5C71%u6912%u7092%u8089*5*1*2|2
	if b, siteConf := GetSiteConf(w, p); b {
		ck, err := r.Cookie("cart")
		if err != nil || ck.Value == "" {
			this.OrderEmpty(w, r, p, mm, siteConf)
			return
		}

		cartData := ording.CartCookieFmt(ck.Value)
		cart, err := goclient.Partner.GetShoppingCart(p.Id, p.Secret, cartData)
		if err != nil {
			w.Write([]byte("订单异常，请清空重新下单"))
			return
		}

		cart.Summary = strings.Replace(cart.Summary, "\n", "<br />", -1)

		this.Context.Template().Execute(w, func(m *map[string]interface{}) {
			(*m)["partner"] = p
			(*m)["title"] = "订单确认-" + p.Name
			(*m)["member"] = mm
			(*m)["cart"] = cart
			(*m)["promFee"] = cart.TotalFee - cart.OrderFee
			(*m)["summary"] = template.HTML(cart.Summary)
			(*m)["conf"] = siteConf
		},
			"views/web/www/order_confirm.html",
			"views/web/www/inc/header.html",
			"views/web/www/inc/footer.html")
	}
}

func (this *ordingC) OrderEmpty(w http.ResponseWriter, r *http.Request,
	p *entity.Partner, mm *member.ValueMember, conf *partner.SiteConf) {
	this.Context.Template().Execute(w, func(m *map[string]interface{}) {
		(*m)["partner"] = p
		(*m)["title"] = "订单确认-" + p.Name
		(*m)["member"] = mm
		(*m)["conf"] = conf
	},
		"views/web/www/order_empty.html",
		"views/web/www/inc/header.html",
		"views/web/www/inc/footer.html")
}

func (this *ordingC) Finish(w http.ResponseWriter, r *http.Request,
	p *entity.Partner, mm *member.ValueMember) {
	// 清除购物车
	cookie, _ := r.Cookie("cart")
	if cookie != nil {
		cookie.Expires = time.Now().Add(time.Hour * 24 * -30)
		cookie.Path = "/"
		http.SetCookie(w, cookie)
	}

	if b, siteConf := GetSiteConf(w, p); b {
		orderNo := r.URL.Query().Get("order_no")
		order, err := goclient.Partner.GetOrderByNo(p.Id, p.Secret, orderNo)
		if err != nil {
			this.OrderEmpty(w, r, p, mm, siteConf)
			return
		}

		if b, siteConf := GetSiteConf(w, p); b {
			this.Context.Template().Execute(w, func(m *map[string]interface{}) {
				(*m)["partner"] = p
				(*m)["title"] = "订单成功-" + p.Name
				(*m)["member"] = mm
				(*m)["conf"] = siteConf
				(*m)["order"] = order
			},
				"views/web/www/order_finish.html",
				"views/web/www/inc/header.html",
				"views/web/www/inc/footer.html")
		}
	}
}

func (this *ordingC) SubmitOrder_post(w http.ResponseWriter, r *http.Request,
	p *entity.Partner, mm *member.ValueMember) {
	if p == nil || mm == nil {
		w.Write([]byte(`{"result":false,"tag":"101"}`)) //未登录
		return
	}
	r.ParseForm()
	ck, err := r.Cookie("cart")
	if err != nil || ck.Value == "" {
		w.Write([]byte(`{"result":false,"tag":"102"}`)) //购物车为空
		return
	}
	cart := ording.CartCookieFmt(ck.Value)
	deliverAddrId, _ := strconv.Atoi(r.FormValue("AddrId"))
	couponCode := r.FormValue("CouponCode")
	order_no, err := goclient.Partner.SubmitOrder(p.Id, p.Secret, mm.Id,
		0, enum.PAY_OFFLINE, deliverAddrId, cart, couponCode, r.FormValue("note"))
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"result":false,"tag":"109","message":"%s"}`, err.Error())))
		return
	}

	w.Write([]byte(`{"result":true,"data":"` + order_no + `"}`))
}
