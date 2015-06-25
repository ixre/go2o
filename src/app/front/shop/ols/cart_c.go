/**
 * Copyright 2015 @ S1N1 Team.
 * name : cart_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ols

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"strconv"
	"strings"
	"net/http"
	"time"
)

type CartC struct {
	*BaseC
}

// 购物车
func (this *CartC) CartApiHandle(ctx *web.Context) {
	if !this.BaseC.Requesting(ctx) {
		ctx.ResponseWriter.Write([]byte("invalid request"))
		return
	}

	r, _ := ctx.Request, ctx.ResponseWriter
	p := this.BaseC.GetPartner(ctx)
	m := this.BaseC.GetMember(ctx)
	r.ParseForm()
	var action = strings.ToLower(r.FormValue("action"))
	var cartKey = r.FormValue("cart.key")
	var memberId int
	if m != nil {
		memberId = m.Id
	}

	switch action {
	case "get":
		this.cart_GetCart(ctx, p, memberId, cartKey)
	case "add":
		this.cart_AddItem(ctx, p, memberId, cartKey)
	case "remove":
		this.cart_RemoveItem(ctx, p, memberId, cartKey)
	}

}

func (this *CartC) cart_GetCart(ctx *web.Context, p *partner.ValuePartner,
	memberId int, cartKey string) {
	cart := dps.ShoppingService.GetShoppingCart(p.Id, memberId, cartKey)

	if cart.Items != nil {
		for _, v := range cart.Items {
			v.GoodsImage = format.GetGoodsImageUrl(v.GoodsImage)
		}
	}

	// 持续保存cookie
	ck,err := ctx.Request.Cookie("_cart");
	if err != nil{
		ck = &http.Cookie{
			Name:"_cart",
			Path:"/",
		}
	}
	ck.Value = cart.CartKey
	ck.Expires = time.Now().Add(time.Hour * 48)
	http.SetCookie(ctx.ResponseWriter,ck)

	d, _ := json.Marshal(cart)
	ctx.ResponseWriter.Write(d)
}

func (this *CartC) cart_AddItem(ctx *web.Context,
	p *partner.ValuePartner, memberId int, cartKey string) {
	r := ctx.Request
	goodsId, _ := strconv.Atoi(r.FormValue("id"))
	num, _ := strconv.Atoi(r.FormValue("num"))
	item := dps.ShoppingService.AddCartItem(p.Id, memberId, cartKey, goodsId, num)

	var d map[string]interface{} = make(map[string]interface{})
	if item == nil {
		d["error"] = "商品不存在"
	} else {
		item.GoodsImage = format.GetGoodsImageUrl(item.GoodsImage)
		d["item"] = item
	}
	this.BaseC.JsonOutput(ctx, d)
}

func (this *CartC) cart_RemoveItem(ctx *web.Context,
	p *partner.ValuePartner, memberId int, cartKey string) {
	var msg gof.Message
	r := ctx.Request
	goodsId, _ := strconv.Atoi(r.FormValue("id"))
	num, _ := strconv.Atoi(r.FormValue("num"))
	err := dps.ShoppingService.SubCartItem(p.Id, memberId, cartKey, goodsId, num)
	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
	this.BaseC.ResultOutput(ctx, msg)
}

func (this *CartC) Index(ctx *web.Context) {
	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{},
		"views/shop/ols/{device}/cart.html",
		"views/shop/ols/{device}/inc/header.html",
		"views/shop/ols/{device}/inc/footer.html")
}
