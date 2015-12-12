/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

//关于regexp的资料
//http://www.cnblogs.com/golove/archive/2013/08/20/3270918.html

import (
	"encoding/json"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
)

type shopC struct {
}

func (this *shopC) ShopList(ctx *echox.Context) error {
	d := ctx.NewData()
	return ctx.RenderOK("shop/shop_list.html", d)
}

//修改门店信息
func (this *shopC) Create(ctx *echox.Context) error {
	d := ctx.NewData()
	d.Map["entity"] = template.JS("{}")
	return ctx.RenderOK("shop/create.html", d)
}

//修改门店信息
func (this *shopC) Modify(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	id, _ := strconv.Atoi(ctx.Query("id"))
	shop := dps.PartnerService.GetShopValueById(partnerId, id)
	entity, _ := json.Marshal(shop)

	d := ctx.NewData()
	d.Map["entity"] = template.JS(entity)
	return ctx.RenderOK("shop/modify.html", d)
}

//修改门店信息
func (this *shopC) SaveShop_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	var result gof.Message
	r.ParseForm()

	shop := partner.ValueShop{}
	web.ParseFormToEntity(r.Form, &shop)

	id, err := dps.PartnerService.SaveShop(partnerId, &shop)

	if err != nil {
		result = gof.Message{Result: true, Message: err.Error()}
	} else {
		result = gof.Message{Result: true, Message: "", Data: id}
	}
	return ctx.JSON(http.StatusOK, result)
}

func (this *shopC) Del_post(ctx *echox.Context) error {
	var result gof.Message
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	r.ParseForm()
	shopId, err := strconv.Atoi(r.FormValue("id"))
	if err == nil {
		err = dps.PartnerService.DeleteShop(partnerId, shopId)
	}

	if err != nil {
		result.Message = err.Error()
	} else {
		result.Result = true
	}
	return ctx.JSON(http.StatusOK, result)
}
