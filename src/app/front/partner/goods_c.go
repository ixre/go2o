/**
 * Copyright 2014 @ S1N1 Team.
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
	"github.com/atnet/gof/web/mvc"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"html/template"
	"strconv"
	"strings"
)

var _ mvc.Filter = new(goodsC)

type goodsC struct {
	*baseC
}

//食物列表
func (this *goodsC) List(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()

	cateOpts := cache.GetDropOptionsOfCategory(this.GetPartnerId(ctx))
	ctx.App.Template().Execute(w, gof.TemplateDataMap{
		"cate_opts":  template.HTML(cateOpts),
		"no_pic_url": format.GetGoodsImageUrl(""),
	}, "views/partner/goods/goods_list.html")
}

func (this *goodsC) Create(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	shopChks := cache.GetShopCheckboxs(partnerId, "")
	cateOpts := cache.GetDropOptionsOfCategory(partnerId)

	ctx.App.Template().Execute(ctx.ResponseWriter, gof.TemplateDataMap{
		"shop_chk":  template.HTML(shopChks),
		"cate_opts": template.HTML(cateOpts),
	},
		"views/partner/goods/create_goods.html")
}

func (this *goodsC) Edit(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
	var e *sale.ValueItem
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	e = dps.SaleService.GetValueItem(partnerId, id)
	if e == nil {
		w.Write([]byte("商品不存在"))
		return
	}
	js, _ := json.Marshal(e)

	shopChks := cache.GetShopCheckboxs(partnerId, e.ApplySubs)
	cateOpts := cache.GetDropOptionsOfCategory(partnerId)

	ctx.App.Template().Execute(w,
		gof.TemplateDataMap{
			"entity":    template.JS(js),
			"shop_chk":  template.HTML(shopChks),
			"cate_opts": template.HTML(cateOpts),
		},
		"views/partner/goods/update_goods.html")
}

func (this *goodsC) SaveItem_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.Message
	r.ParseForm()

	e := sale.ValueItem{}
	web.ParseFormToEntity(r.Form, &e)

	id, err := dps.SaleService.SaveItem(partnerId, &e)

	if err != nil {
		result = gof.Message{Result: true, Message: err.Error()}
	} else {
		result = gof.Message{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}

func (this *goodsC) Del_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.Message

	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("id"))
	err := dps.SaleService.DeleteGoods(partnerId, id)

	if err != nil {
		result = gof.Message{Result: true, Message: err.Error()}
	} else {
		result = gof.Message{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}

func (this *goodsC) SetSaleTag(ctx *web.Context) {
	r := ctx.Request
	r.ParseForm()
	partnerId := this.GetPartnerId(ctx)
	goodsId, _ := strconv.Atoi(r.URL.Query().Get("id"))

	var tags []*sale.ValueSaleTag = dps.SaleService.GetAllSaleTags(partnerId)
	tagsHtml := getSaleTagsCheckBoxHtml(tags)

	var chkTags []*sale.ValueSaleTag = dps.SaleService.GetItemSaleTags(partnerId, goodsId)
	strArr := make([]string, len(chkTags))
	for i, v := range chkTags {
		strArr[i] = strconv.Itoa(v.Id)
	}

	tagVal := strings.Join(strArr, ",")

	ctx.App.Template().Execute(ctx.ResponseWriter, gof.TemplateDataMap{
		"goodsId":  goodsId,
		"tagsHtml": template.HTML(tagsHtml),
		"tagValue": tagVal,
	}, "views/partner/goods/set_sale_tag.html")
}

func (this *goodsC) SaveGoodsSTag_post(ctx *web.Context) {
	r := ctx.Request
	var msg gof.Message
	goodsId, err := strconv.Atoi(r.FormValue("GoodsId"))
	if err == nil {
		tags := strings.Split(r.FormValue("SaleTags"), ",")
		var ids []int = []int{}
		for _, v := range tags {
			if i, err := strconv.Atoi(v); err == nil {
				ids = append(ids, i)
			}
		}

		partnerId := this.GetPartnerId(ctx)
		err = dps.SaleService.SaveItemSaleTags(partnerId, goodsId, ids)
	}

	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
	this.ResultOutput(ctx, msg)
}
