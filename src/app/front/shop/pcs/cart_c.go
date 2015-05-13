/**
 * Copyright 2015 @ S1N1 Team.
 * name : cart_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package pcshop

import (
	"encoding/json"
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/goclient"
	"strconv"
	"strings"
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

func (this *cartC) cart_GetCart(ctx *web.Context,
	p *partner.ValuePartner, memberId int, cartKey string) {
	cart := goclient.Partner.GetShoppingCart(p.Id, memberId, cartKey)

	// 如果已经购买，則创建新的购物车
	if cart.IsBought == 1 {
		cart = goclient.Partner.GetShoppingCart(p.Id, memberId, "")
	}

	d, _ := json.Marshal(cart)
	ctx.ResponseWriter.Write(d)
}

func (this *cartC) cart_AddItem(ctx *web.Context,
	p *partner.ValuePartner, memberId int, cartKey string) {
	r, w := ctx.Request, ctx.ResponseWriter
	goodsId, _ := strconv.Atoi(r.FormValue("id"))
	num, _ := strconv.Atoi(r.FormValue("num"))
	item, err := goclient.Partner.AddCartItem(p.Id, memberId, cartKey, goodsId, num)

	var result = make(map[string]interface{}, 2)
	if err != nil {
		result["message"] = err.Error()
	} else {
		result["message"] = ""
		result["item"] = item
	}
	d, _ := json.Marshal(result)
	w.Write(d)
}

func (this *cartC) cart_RemoveItem(ctx *web.Context,
	p *partner.ValuePartner, memberId int, cartKey string) {
	r, w := ctx.Request, ctx.ResponseWriter
	goodsId, _ := strconv.Atoi(r.FormValue("id"))
	num, _ := strconv.Atoi(r.FormValue("num"))
	err := goclient.Partner.SubCartItem(p.Id, memberId, cartKey, goodsId, num)
	if err != nil {
		w.Write([]byte(`{error:'` + err.Error() + `'}`))
	} else {
		w.Write([]byte("{}"))
	}
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
