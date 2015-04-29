/**
 * Copyright 2014 @ ops Inc.
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
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"html/template"
	"strconv"
)

//==========================  后台逻辑 =============================//
type shopC struct {
	gof.App
}

func (this *shopC) ShopList(ctx *web.Context) {
	ctx.App.Template().Render(ctx.ResponseWriter, "views/partner/shop/shop_list.html", nil)
}

//修改门店信息
func (this *shopC) Create(ctx *web.Context) {
	ctx.App.Template().Render(ctx.ResponseWriter,
		"views/partner/shop/create.html",
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS("{}")
		})
}

//修改门店信息
func (this *shopC) Modify(ctx *web.Context, partnerId int) {
	r, w := ctx.Request, ctx.ResponseWriter
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	shop := dps.PartnerService.GetShopValueById(partnerId, id)
	entity, _ := json.Marshal(shop)

	ctx.App.Template().Render(w,
		"views/partner/shop/modify.html",
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(entity)
		})
}

//修改门店信息
func (this *shopC) SaveShop_post(ctx *web.Context, partnerId int) {
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.JsonResult
	r.ParseForm()

	shop := partner.ValueShop{}
	web.ParseFormToEntity(r.Form, &shop)

	id, err := dps.PartnerService.SaveShop(partnerId, &shop)

	if err != nil {
		result = gof.JsonResult{Result: true, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}

func (this *shopC) Del_post(ctx *web.Context, partnerId int) {
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	shopId, err := strconv.Atoi(r.FormValue("id"))
	if err == nil {
		err = dps.PartnerService.DeleteShop(partnerId, shopId)
	}

	if err != nil {
		w.Write([]byte("{result:false,message:'" + err.Error() + "'}"))
	} else {
		w.Write([]byte("{result:true}"))
	}
}
