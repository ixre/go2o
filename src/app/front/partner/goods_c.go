/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/app/cache"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/service/dps"
	"html/template"
	"strconv"
)

type goodsC struct {
	gof.App
}

//食物列表
func (this *goodsC) List(ctx *web.Context) {
	/*
		'''
		菜单列表
		'''
		req=web.input(cid=-1,returnUri='')
		_dataurl='index?m=food&act=foods&ajax=1&cid=%s'%(req.category_id)

		return render.foods(dataurl=_dataurl)
	*/
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	//cid,_:= strconv.Atoi(r.FormValue("cid"))
	ctx.App.Template().Render(w,
		"views/partner/goods/list.html",
		nil)
}

func (this *goodsC) Create(ctx *web.Context, ptId int) {
	shopChks := cache.GetShopCheckboxs(ptId, "")
	cateOpts := cache.GetDropOptionsOfCategory(ptId)

	ctx.App.Template().Render(ctx.ResponseWriter,
		"views/partner/goods/create_goods.html",
		func(m *map[string]interface{}) {
			(*m)["shop_chk"] = template.HTML(shopChks)
			(*m)["cate_opts"] = template.HTML(cateOpts)
		})
}

func (this *goodsC) Edit(ctx *web.Context, partnerId int) {
	r, w := ctx.Request, ctx.ResponseWriter
	var e *sale.ValueGoods
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	e = dps.SaleService.GetValueGoods(partnerId, id)
	if e == nil {
		w.Write([]byte("商品不存在"))
		return
	}
	js, _ := json.Marshal(e)

	shopChks := cache.GetShopCheckboxs(partnerId, e.ApplySubs)
	cateOpts := cache.GetDropOptionsOfCategory(partnerId)

	ctx.App.Template().Render(w,
		"views/partner/goods/update_goods.html",
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(js)
			(*m)["shop_chk"] = template.HTML(shopChks)
			(*m)["cate_opts"] = template.HTML(cateOpts)
		})
}

func (this *goodsC) SaveItem_post(ctx *web.Context, partnerId int) {
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.JsonResult
	r.ParseForm()

	e := sale.ValueGoods{}
	web.ParseFormToEntity(r.Form, &e)

	id, err := dps.SaleService.SaveGoods(partnerId, &e)

	if err != nil {
		result = gof.JsonResult{Result: true, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}

func (this *goodsC) Del_post(ctx *web.Context, partnerId int) {
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.JsonResult

	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("id"))
	err := dps.SaleService.DeleteGoods(partnerId, id)

	if err != nil {
		result = gof.JsonResult{Result: true, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}
