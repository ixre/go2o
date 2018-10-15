/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jsix/gof"
	gfmt "github.com/jsix/gof/util/fmt"
	//"github.com/jsix/gof/web"
	"go2o/src/app/cache"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"github.com/jsix/gof/web/form"
)

type goodsC struct {
}

//货品列表
func (this *goodsC) Item_list(ctx *echox.Context) error {
	cateOpts := cache.GetDropOptionsOfCategory(getPartnerId(ctx))

	d := ctx.NewData()
	d.Map["cate_opts"] = template.HTML(cateOpts)
	d.Map["no_pic_url"] = format.GetGoodsImageUrl("")
	return ctx.RenderOK("goods.item_list.html", d)
}

//货品选择
func (this *goodsC) Goods_select(ctx *echox.Context) error {
	cateOpts := cache.GetDropOptionsOfCategory(getPartnerId(ctx))
	d := ctx.NewData()
	d.Map["cate_opts"] = template.HTML(cateOpts)
	d.Map["no_pic_url"] = format.GetGoodsImageUrl("")
	return ctx.RenderOK("goods.select.html", d)
}

func (this *goodsC) Create(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	shopChks := cache.GetShopCheckboxs(partnerId, "")
	cateOpts := cache.GetDropOptionsOfCategory(partnerId)

	e := &sale.ValueItem{
		Image: ctx.App.Config().GetString(variable.NoPicPath),
	}
	js, _ := json.Marshal(e)

	d := ctx.NewData()
	d.Map = map[string]interface{}{
		"entity":    template.JS(js),
		"shop_chk":  template.HTML(shopChks),
		"cate_opts": template.HTML(cateOpts),
		"Image":     format.GetGoodsImageUrl(e.Image),
	}
	return ctx.RenderOK("goods.create_goods.html", d)
}

func (this *goodsC) Edit(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	var e *sale.ValueItem
	ss := dps.SaleService
	id, _ := strconv.Atoi(ctx.Query("item_id"))
	e = ss.GetValueItem(partnerId, id)
	if e == nil {
		return ctx.StringOK("商品不存在")
	}
	e.Description = ""
	js, _ := json.Marshal(e)
	gs := ss.GetGoodsBySku(partnerId, e.Id, 0) //todo:???
	shopChks := cache.GetShopCheckboxs(partnerId, e.ApplySubs)
	cateOpts := cache.GetDropOptionsOfCategory(partnerId)

	d := ctx.NewData()
	d.Map = map[string]interface{}{
		"entity":    template.JS(js),
		"shop_chk":  template.HTML(shopChks),
		"cate_opts": template.HTML(cateOpts),
		"Image":     format.GetGoodsImageUrl(e.Image),
		"gs":        gs,
	}
	return ctx.RenderOK("goods.update_goods.html", d)
}

// 保存商品描述
func (this *goodsC) Item_info(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	var e *sale.ValueItem
	id, _ := strconv.Atoi(r.URL.Query().Get("item_id"))
	e = dps.SaleService.GetValueItem(partnerId, id)
	if e == nil {
		return ctx.String(http.StatusOK, "商品不存在")
	}

	d := ctx.NewData()
	d.Map = map[string]interface{}{
		"item_id":   e.Id,
		"item_info": template.HTML(e.Description),
	}
	return ctx.RenderOK("goods.item_info.html", d)
}

// 保存货品描述信息(POST)
func (this *goodsC) Save_item_info(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		r.ParseForm()
		id, _ := strconv.Atoi(r.FormValue("ItemId"))
		info := r.FormValue("Info")

		var result gof.Result
		err := dps.SaleService.SaveItemInfo(partnerId, id, info)

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
		}

		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 保存货品信息(POST)
func (this *goodsC) SaveItem(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		ss := dps.SaleService
		var result gof.Result
		r.ParseForm()
		e := sale.ValueItem{}
		form.ParseEntity(r.Form, &e)
		e.State = 1 //todo: 暂时使用
		id, err := ss.SaveItem(partnerId, &e)
		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			gs := ss.GetValueGoodsBySku(partnerId, id, 0) //todo: ??? SKU
			gs.StockNum, _ = strconv.Atoi(r.FormValue("StockNum"))
			gs.SaleNum, _ = strconv.Atoi(r.FormValue("SaleNum"))
			price, _ := strconv.ParseFloat(r.FormValue("SalePrice"), 32)
			gs.SalePrice = float32(price)
			ss.SaveGoods(partnerId, gs)
			result.ErrCode = 0
			var data = make(map[string]string)
			data["id"] = fmt.Sprintf("%d", id)
			result.Data = data
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 删除商品信息(POST)
func (this *goodsC) Del_goods(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	var result gof.Result
	if r.Method == "POST" {
		r.ParseForm()
		id, _ := strconv.Atoi(r.FormValue("id"))
		err := dps.SaleService.DeleteGoods(partnerId, id)

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 删除货品信息(POST)
func (this *goodsC) Del_item(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		var result gof.Result

		r.ParseForm()
		id, _ := strconv.Atoi(r.FormValue("id"))
		err := dps.SaleService.DeleteItem(partnerId, id)

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

// 设置销售标签
func (this *goodsC) SetSaleTag(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	goodsId, _ := strconv.Atoi(ctx.Query("id"))

	var tags []*sale.ValueSaleTag = dps.SaleService.GetAllSaleTags(partnerId)
	tagsHtml := getSaleTagsCheckBoxHtml(tags)

	var chkTags []*sale.ValueSaleTag = dps.SaleService.GetItemSaleTags(partnerId, goodsId)
	strArr := make([]string, len(chkTags))
	for i, v := range chkTags {
		strArr[i] = strconv.Itoa(v.Id)
	}

	tagVal := strings.Join(strArr, ",")

	d := ctx.NewData()
	d.Map = map[string]interface{}{
		"goodsId":  goodsId,
		"tagsHtml": template.HTML(tagsHtml),
		"tagValue": tagVal,
	}
	return ctx.RenderOK("goods.set_sale_tag.html", d)
}

// 保存销售标签(POST)
func (this *goodsC) SaveGoodsSTag(ctx *echox.Context) error {
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		r.ParseForm()
		var result gof.Result
		goodsId, err := strconv.Atoi(r.FormValue("GoodsId"))
		if err == nil {
			tags := strings.Split(r.FormValue("SaleTags"), ",")
			var ids []int = []int{}
			for _, v := range tags {
				if i, err := strconv.Atoi(v); err == nil {
					ids = append(ids, i)
				}
			}

			partnerId := getPartnerId(ctx)
			err = dps.SaleService.SaveItemSaleTags(partnerId, goodsId, ids)
		}

		if err != nil {
			result.ErrMsg = err.Error()
		} else {
			result.ErrCode = 0
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

func (this *goodsC) ItemCtrl(ctx *echox.Context) error {

	itemId, _ := strconv.Atoi(ctx.Query("item_id"))

	d := ctx.NewData()
	d.Map["item_id"] = itemId
	return ctx.RenderOK("goods.item_ctrl.html", d)
}

func (this *goodsC) LvPrice(ctx *echox.Context) error {
	if ctx.Request().Method == "POST" {
		return this.lvPrice_post(ctx)
	}
	partnerId := getPartnerId(ctx)
	//todo: should be goodsId
	itemId, _ := strconv.Atoi(ctx.Query("item_id"))
	goods := dps.SaleService.GetGoodsBySku(partnerId, itemId, 0)
	lvs := dps.PartnerService.GetMemberLevels(partnerId)
	var prices []*sale.MemberPrice = dps.SaleService.GetGoodsLevelPrices(partnerId, goods.GoodsId)

	var buf *bytes.Buffer = bytes.NewBufferString("")

	var fmtFunc = func(level int, levelName string, id int, price float32, enabled int) {
		buf.WriteString(fmt.Sprintf(`
		<tr>
                <td><input type="hidden" field="Id_%d" value="%d"/>
                    %s</td>
                <td align="center"><input type="number" field="Price_%d" value="%s"/></td>
                <td align="center"><input type="checkbox" field="Enabled_%d" %s/></td>
            </tr>
		`, level, id, levelName, level, format.FormatFloat(price), level,
			gfmt.BoolString(enabled == 1, "checked=\"checked\"", "")))
	}

	var b bool
	for _, v := range lvs {
		b = false
		for _, v1 := range prices {
			if v.Value == v1.Level {
				fmtFunc(v.Value, v.Name, v1.Id, v1.Price, v1.Enabled)
				b = true
				break
			}
		}
		if !b {
			fmtFunc(v.Value, v.Name, 0, goods.SalePrice, 0)
		}
	}

	d := ctx.NewData()
	d.Map = map[string]interface{}{
		"goods":   goods,
		"setHtml": template.HTML(buf.String()),
	}
	return ctx.RenderOK("goods.level_price.html", d)
}

func (this *goodsC) lvPrice_post(ctx *echox.Context) error {
	req := ctx.HttpRequest()
	req.ParseForm()
	goodsId, err := strconv.Atoi(req.FormValue("goodsId"))
	if err != nil {
		return ctx.JSON(http.StatusOK, gof.Result{ErrMsg: err.Error()})
	}

	var priceSet []*sale.MemberPrice = []*sale.MemberPrice{}
	var id int
	var price float64
	var lv int
	var enabled int

	for k, _ := range req.Form {
		if strings.HasPrefix(k, "Id_") {
			if lv, err = strconv.Atoi(k[3:]); err == nil {
				id, _ = strconv.Atoi(req.Form.Get(k))
				price, _ = strconv.ParseFloat(req.FormValue(fmt.Sprintf("Price_%d", lv)), 32)
				if req.FormValue(fmt.Sprintf("Enabled_%d", lv)) == "on" {
					enabled = 1
				} else {
					enabled = 0
				}

				priceSet = append(priceSet, &sale.MemberPrice{
					Id:      id,
					Level:   lv,
					GoodsId: goodsId,
					Price:   float32(price),
					Enabled: enabled,
				})
			} else {
				return ctx.JSON(http.StatusOK, gof.Result{ErrMsg: err.Error()})
			}
		}
	}

	partnerId := getPartnerId(ctx)
	err = dps.SaleService.SaveMemberPrices(partnerId, goodsId, priceSet)
	if err != nil {
		return ctx.JSON(http.StatusOK, gof.Result{ErrMsg: err.Error()})
	} else {
		return ctx.JSON(http.StatusOK, gof.Result{ErrCode: 0})
	}
}
