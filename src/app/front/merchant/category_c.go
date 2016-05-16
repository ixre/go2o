/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package merchant

import (
	"encoding/json"
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/ui/tree"
	"go2o/src/app/cache"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
)

type categoryC struct {
}

//分类树形功能
func (this *categoryC) All_category(ctx *echox.Context) error {
	d := ctx.NewData()
	d.Map["no_pic_url"] = format.GetGoodsImageUrl("")
	return ctx.RenderOK("category.index.html", d)
}

//分类Json数据
func (this *categoryC) CategoryJson(ctx *echox.Context) error {
	merchantId := getMerchantId(ctx)
	var node *tree.TreeNode = dps.SaleService.GetCategoryTreeNode(merchantId)
	return ctx.JSON(http.StatusOK, node)
}

//分类树形功能
func (this *categoryC) CategorySelect(ctx *echox.Context) error {
	d := ctx.NewData()
	return ctx.RenderOK("category.select.html", d)
}

//分类Json数据
func (this *categoryC) CreateCategory(ctx *echox.Context) error {
	merchantId := getMerchantId(ctx)

	cateOpts := cache.GetDropOptionsOfCategory(merchantId)

	e := &sale.ValueCategory{
		Icon: ctx.App.Config().GetString(variable.NoPicPath),
	}
	eJson, _ := json.Marshal(e)

	d := ctx.NewData()
	d.Map = map[string]interface{}{
		"Entity":   template.JS(eJson),
		"CateOpts": template.HTML(cateOpts),
		"Icon":     format.GetGoodsImageUrl(e.Icon),
	}
	return ctx.RenderOK("category.create.html", d)
}

func (this *categoryC) EditCategory(ctx *echox.Context) error {
	merchantId := getMerchantId(ctx)
	r := ctx.HttpRequest()
	r.ParseForm()
	id, _ := strconv.Atoi(r.Form.Get("id"))
	e, _ := dps.SaleService.GetCategory(merchantId, id)
	eJson, _ := json.Marshal(e)

	re := regexp.MustCompile(fmt.Sprintf("<option class=\"opt\\d+\" value=\"%d\">[^>]+>", id))
	originOpts := cache.GetDropOptionsOfCategory(merchantId)
	cateOpts := re.ReplaceAll(originOpts, nil)

	d := ctx.NewData()
	d.Map = map[string]interface{}{
		"Entity":   template.JS(eJson),
		"CateOpts": template.HTML(cateOpts),
		"Icon":     format.GetGoodsImageUrl(e.Icon),
	}
	return ctx.RenderOK("category.edit.html", d)
}

//修改门店信息
func (this *categoryC) SaveCategory(ctx *echox.Context) error {
	merchantId := getMerchantId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		var result gof.Message
		r.ParseForm()

		e := sale.ValueCategory{}
		web.ParseFormToEntity(r.Form, &e)

		id, err := dps.SaleService.SaveCategory(merchantId, &e)
		if err != nil {
			result = gof.Message{Result: false, Message: err.Error()}
		} else {
			result = gof.Message{Result: true, Message: "", Data: id}
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

func (this *categoryC) DelCategory(ctx *echox.Context) error {
	merchantId := getMerchantId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		var result gof.Message
		r.ParseForm()
		categoryId, _ := strconv.Atoi(r.FormValue("id"))

		//删除分类
		err := dps.SaleService.DeleteCategory(merchantId, categoryId)
		if err != nil {
			result = gof.Message{Result: false, Message: err.Error()}
		} else {
			result = gof.Message{Result: true, Message: ""}
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}
