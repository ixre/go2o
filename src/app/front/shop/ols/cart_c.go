/**
 * Copyright 2015 @ z3q.net.
 * name : cart_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ols

import (
	"github.com/jsix/gof"
	"go2o/src/core/domain/interface/merchant"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CartC struct {
}

// 购物车API(POST)
func (this *CartC) CartApiHandle(ctx *echox.Context) error {
	r := ctx.HttpRequest()
	if r.Method != "POST" {
		return nil
	}
	r.ParseForm()
	p := getPartner(ctx)
	m := GetMember(ctx)
	var action = strings.ToLower(r.FormValue("action"))
	var cartKey = r.FormValue("cart.key")
	var memberId int
	if m != nil {
		memberId = m.Id
	}

	switch action {
	case "get":
		return this.cart_GetCart(ctx, p, memberId, cartKey)
	case "add":
		return this.cart_AddItem(ctx, p, memberId, cartKey)
	case "remove":
		return this.cart_RemoveItem(ctx, p, memberId, cartKey)
	}
	return nil
}

func (this *CartC) cart_GetCart(ctx *echox.Context, p *merchant.MerchantValue,
	memberId int, cartKey string) error {
	cart := dps.ShoppingService.GetShoppingCart(p.Id, memberId, cartKey)

	if cart.Items != nil {
		for _, v := range cart.Items {
			v.GoodsImage = format.GetGoodsImageUrl(v.GoodsImage)
		}
	}

	// 持续保存cookie
	ck, err := ctx.HttpRequest().Cookie("_cart")
	if err != nil {
		ck = &http.Cookie{
			Name: "_cart",
			Path: "/",
		}
	}
	ck.Value = cart.CartKey
	ck.Expires = time.Now().Add(time.Hour * 48)
	http.SetCookie(ctx.HttpResponse(), ck)
	return ctx.JSON(http.StatusOK, cart)
}

func (this *CartC) cart_AddItem(ctx *echox.Context,
	p *merchant.MerchantValue, memberId int, cartKey string) error {
	r := ctx.HttpRequest()
	goodsId, _ := strconv.Atoi(r.FormValue("id"))
	num, _ := strconv.Atoi(r.FormValue("num"))
	item, err := dps.ShoppingService.AddCartItem(p.Id, memberId, cartKey, goodsId, num)

	var d map[string]interface{} = make(map[string]interface{})
	if err != nil {
		d["error"] = err.Error()
	} else {
		item.GoodsImage = format.GetGoodsImageUrl(item.GoodsImage)
		d["item"] = item
	}
	return ctx.JSON(http.StatusOK, d)
}

func (this *CartC) cart_RemoveItem(ctx *echox.Context,
	p *merchant.MerchantValue, memberId int, cartKey string) error {
	var result gof.Message
	r := ctx.HttpRequest()
	goodsId, _ := strconv.Atoi(r.FormValue("id"))
	num, _ := strconv.Atoi(r.FormValue("num"))
	err := dps.ShoppingService.SubCartItem(p.Id, memberId, cartKey, goodsId, num)
	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	return ctx.JSON(http.StatusOK, result)
}

func (this *CartC) Index(ctx *echox.Context) error {
	return ctx.RenderOK("cart.html", ctx.NewData())
}
