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
	"encoding/json"
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/ui/tree"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
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
	d := echox.NewRenderData()
	d.Map["no_pic_url"] = format.GetGoodsImageUrl("")
	return ctx.RenderOK("category.index.html", d)
}

//分类Json数据
func (this *categoryC) CategoryJson(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	var node *tree.TreeNode = dps.SaleService.GetCategoryTreeNode(partnerId)
	return ctx.JSON(http.StatusOK, node)
}

//分类树形功能
func (this *categoryC) CategorySelect(ctx *echox.Context) error {
	d := echox.NewRenderData()
	return ctx.RenderOK("category.select.html", d)
}

//分类Json数据
func (this *categoryC) CreateCategory(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)

	var node *tree.TreeNode = dps.SaleService.GetCategoryTreeNode(partnerId)
	json, _ := json.Marshal(node)

	d := echox.NewRenderData()
	d.Map["treeJson"] = template.JS(json)
	return ctx.RenderOK("category.create.html", d)
}

func (this *categoryC) EditCategory(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	r.ParseForm()
	id, _ := strconv.Atoi(r.Form.Get("id"))
	var category *sale.ValueCategory = dps.SaleService.GetCategory(partnerId, id)
	json, _ := json.Marshal(category)

	re := regexp.MustCompile(fmt.Sprintf("<option class=\"opt\\d+\" value=\"%d\">[^>]+>", id))
	originOpts := cache.GetDropOptionsOfCategory(partnerId)
	cateOpts := re.ReplaceAll(originOpts, nil)

	d := echox.NewRenderData()
	d.Map = map[string]interface{}{
		"entity":    template.JS(json),
		"cate_opts": template.HTML(cateOpts),
	}
	return ctx.RenderOK("category.edit.html", d)
}

//修改门店信息
func (this *categoryC) SaveCategory_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	var result gof.Message
	r.ParseForm()

	e := sale.ValueCategory{}
	web.ParseFormToEntity(r.Form, &e)

	id, err := dps.SaleService.SaveCategory(partnerId, &e)
	if err != nil {
		result = gof.Message{Result: false, Message: err.Error()}
	} else {
		result = gof.Message{Result: true, Message: "", Data: id}
	}
	return ctx.JSON(http.StatusOK, result)
}

func (this *categoryC) DelCategory_post(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	var result gof.Message
	r.ParseForm()
	categoryId, _ := strconv.Atoi(r.FormValue("id"))

	//删除分类
	err := dps.SaleService.DeleteCategory(partnerId, categoryId)
	if err != nil {
		result = gof.Message{Result: false, Message: err.Error()}
	} else {
		result = gof.Message{Result: true, Message: ""}
	}
	return ctx.JSON(http.StatusOK, result)
}
