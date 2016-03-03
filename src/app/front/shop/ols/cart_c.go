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
<<<<<<< HEAD
	"github.com/jsix/gof"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
=======
	"encoding/json"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CartC struct {
<<<<<<< HEAD
}

// 购物车API(POST)
func (this *CartC) CartApiHandle(ctx *echox.Context) error {
	r := ctx.Request()
	if r.Method != "POST" {
		return nil
	}
	r.ParseForm()
	p := getPartner(ctx)
	m := GetMember(ctx)
=======
	*BaseC
}

// 购物车
func (this *CartC) CartApiHandle(ctx *web.Context) {
	if !this.BaseC.Requesting(ctx) {
		ctx.Response.Write([]byte("invalid request"))
		return
	}

	r, _ := ctx.Request, ctx.Response
	p := this.BaseC.GetPartner(ctx)
	m := this.BaseC.GetMember(ctx)
	r.ParseForm()
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	var action = strings.ToLower(r.FormValue("action"))
	var cartKey = r.FormValue("cart.key")
	var memberId int
	if m != nil {
		memberId = m.Id
	}

	switch action {
	case "get":
<<<<<<< HEAD
		return this.cart_GetCart(ctx, p, memberId, cartKey)
	case "add":
		return this.cart_AddItem(ctx, p, memberId, cartKey)
	case "remove":
		return this.cart_RemoveItem(ctx, p, memberId, cartKey)
	}
	return nil
}

func (this *CartC) cart_GetCart(ctx *echox.Context, p *partner.ValuePartner,
	memberId int, cartKey string) error {
=======
		this.cart_GetCart(ctx, p, memberId, cartKey)
	case "add":
		this.cart_AddItem(ctx, p, memberId, cartKey)
	case "remove":
		this.cart_RemoveItem(ctx, p, memberId, cartKey)
	}

}

func (this *CartC) cart_GetCart(ctx *web.Context, p *partner.ValuePartner,
	memberId int, cartKey string) {

	//time.Sleep(time.Second*10)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	cart := dps.ShoppingService.GetShoppingCart(p.Id, memberId, cartKey)

	if cart.Items != nil {
		for _, v := range cart.Items {
			v.GoodsImage = format.GetGoodsImageUrl(v.GoodsImage)
		}
	}

	// 持续保存cookie
<<<<<<< HEAD
	ck, err := ctx.Request().Cookie("_cart")
=======
	ck, err := ctx.Request.Cookie("_cart")
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	if err != nil {
		ck = &http.Cookie{
			Name: "_cart",
			Path: "/",
		}
	}
	ck.Value = cart.CartKey
	ck.Expires = time.Now().Add(time.Hour * 48)
<<<<<<< HEAD
	http.SetCookie(ctx.Response(), ck)
	return ctx.JSON(http.StatusOK, cart)
}

func (this *CartC) cart_AddItem(ctx *echox.Context,
	p *partner.ValuePartner, memberId int, cartKey string) error {
	r := ctx.Request()
=======
	http.SetCookie(ctx.Response, ck)

	d, _ := json.Marshal(cart)
	ctx.Response.Write(d)
}

func (this *CartC) cart_AddItem(ctx *web.Context,
	p *partner.ValuePartner, memberId int, cartKey string) {
	r := ctx.Request
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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
<<<<<<< HEAD
	return ctx.JSON(http.StatusOK, d)
}

func (this *CartC) cart_RemoveItem(ctx *echox.Context,
	p *partner.ValuePartner, memberId int, cartKey string) error {
	var result gof.Message
	r := ctx.Request()
=======
	ctx.Response.JsonOutput(d)
}

func (this *CartC) cart_RemoveItem(ctx *web.Context,
	p *partner.ValuePartner, memberId int, cartKey string) {
	var result gof.Message
	r := ctx.Request
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	goodsId, _ := strconv.Atoi(r.FormValue("id"))
	num, _ := strconv.Atoi(r.FormValue("num"))
	err := dps.ShoppingService.SubCartItem(p.Id, memberId, cartKey, goodsId, num)
	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
<<<<<<< HEAD
	return ctx.JSON(http.StatusOK, result)
}

func (this *CartC) Index(ctx *echox.Context) error {
	return ctx.RenderOK("cart.html", ctx.NewData())
=======
	ctx.Response.JsonOutput(result)
}

func (this *CartC) Index(ctx *web.Context) {
	this.BaseC.ExecuteTemplate(ctx, gof.TemplateDataMap{},
		"views/shop/ols/{device}/cart.html",
		"views/shop/ols/{device}/inc/header.html",
		"views/shop/ols/{device}/inc/footer.html")
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
