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
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/partner"
	"strconv"
	"strings"
	"go2o/src/core/service/dps"
	"github.com/atnet/gof"
)

type cartC struct {
	*baseC
}

// 购物车
func (this *cartC) cartApi(ctx *web.Context) {
	if !this.Requesting(ctx) {
		ctx.ResponseWriter.Write([]byte("invalid request"))
		return
	}

	r, _ := ctx.Request, ctx.ResponseWriter
	p := this.GetPartner(ctx)
	m := this.GetMember(ctx)
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

func (this *cartC) cart_GetCart(ctx *web.Context,p *partner.ValuePartner,
	memberId int, cartKey string) {
	cart := dps.ShoppingService.GetShoppingCart(p.Id, memberId, cartKey)

	// 如果已经购买，則创建新的购物车
	if cart.IsBought == 1 {
		cart = dps.ShoppingService.GetShoppingCart(p.Id, memberId, "")
	}

	d, _ := json.Marshal(cart)
	ctx.ResponseWriter.Write(d)
}

func (this *cartC) cart_AddItem(ctx *web.Context,
	p *partner.ValuePartner, memberId int, cartKey string) {
	r := ctx.Request
	goodsId, _ := strconv.Atoi(r.FormValue("id"))
	num, _ := strconv.Atoi(r.FormValue("num"))
	item := dps.ShoppingService.AddCartItem(p.Id, memberId, cartKey, goodsId, num)

	var d map[string]interface{}
	if item == nil{
		d["error"] = "商品不存在"
	}else {
		d["item"] = item
	}
	this.JsonOutput(ctx,d)
}

func (this *cartC) cart_RemoveItem(ctx *web.Context,
	p *partner.ValuePartner, memberId int, cartKey string) {
	var msg gof.Message
	r:= ctx.Request
	goodsId, _ := strconv.Atoi(r.FormValue("id"))
	num, _ := strconv.Atoi(r.FormValue("num"))
	err := dps.ShoppingService.SubCartItem(p.Id, memberId, cartKey, goodsId, num)
	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
	this.ResultOutput(ctx,msg)
}

func (this *cartC) cart(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	//todo: 需页面
	if r.URL.Query().Get("edit") == "1" {
		w.Header().Add("Location", "/list")
	} else {
		w.Header().Add("Location", "/buy/confirm")
	}
	w.WriteHeader(302)
}
