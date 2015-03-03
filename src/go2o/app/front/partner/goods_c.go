/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"encoding/json"
	"github.com/atnet/gof"
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web"
	"go2o/app/cache"
	"go2o/core/domain/interface/sale"
	"go2o/core/service/dps"
	"html/template"
	"net/http"
	"strconv"
)

type goodsC struct {
	Context app.Context
}

//食物列表
func (this *goodsC) List(w http.ResponseWriter, r *http.Request) {
	/*
		'''
		菜单列表
		'''
		req=web.input(cid=-1,returnUri='')
		_dataurl='index?m=food&act=foods&ajax=1&cid=%s'%(req.category_id)

		return render.foods(dataurl=_dataurl)
	*/
	r.ParseForm()
	//cid,_:= strconv.Atoi(r.FormValue("cid"))
	this.Context.Template().Render(w,
		"views/partner/goods/list.html",
		nil)
}

func (this *goodsC) Create(w http.ResponseWriter, r *http.Request, ptId int) {
	shopChks := cache.GetShopCheckboxs(ptId, "")
	cateOpts := cache.GetDropOptionsOfCategory(ptId)

	this.Context.Template().Render(w,
		"views/partner/goods/create_goods.html",
		func(m *map[string]interface{}) {
			(*m)["shop_chk"] = template.HTML(shopChks)
			(*m)["cate_opts"] = template.HTML(cateOpts)
		})
}

func (this *goodsC) Edit(w http.ResponseWriter, r *http.Request, partnerId int) {
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

	this.Context.Template().Render(w,
		"views/partner/goods/update_goods.html",
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(js)
			(*m)["shop_chk"] = template.HTML(shopChks)
			(*m)["cate_opts"] = template.HTML(cateOpts)
		})
}

func (this *goodsC) SaveItem_post(w http.ResponseWriter, r *http.Request, partnerId int) {
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

func (this *goodsC) Del_post(w http.ResponseWriter, r *http.Request, partnerId int) {
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
