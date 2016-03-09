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
<<<<<<< HEAD
=======
	"github.com/jsix/gof/web/mvc"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"github.com/jsix/gof/web/ui/tree"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
<<<<<<< HEAD
	"go2o/src/x/echox"
	"html/template"
	"log"
	"net/http"
=======
	"html/template"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"regexp"
	"strconv"
)

<<<<<<< HEAD
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
=======
var _ mvc.Filter = new(categoryC)

type categoryC struct {
	*baseC
}

//分类树形功能
func (this *categoryC) All_category(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.Response, gof.TemplateDataMap{
		"no_pic_url": format.GetGoodsImageUrl(""),
	}, "views/partner/category/category.html")
}

//分类Json数据
func (this *categoryC) CategoryJson(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	var node *tree.TreeNode = dps.SaleService.GetCategoryTreeNode(partnerId)
	json, _ := json.Marshal(node)
	ctx.Response.Write(json)
}

//分类树形功能
func (this *categoryC) CategorySelect(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.Response, nil,
		"views/partner/category/category_select.html")
}

//分类Json数据
func (this *categoryC) CreateCategory(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d

	var node *tree.TreeNode = dps.SaleService.GetCategoryTreeNode(partnerId)
	json, _ := json.Marshal(node)

<<<<<<< HEAD
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
	log.Println("---", category)
=======
	ctx.App.Template().Execute(ctx.Response,
		gof.TemplateDataMap{
			"treeJson": template.JS(json),
		},
		"views/partner/category/category_create.html")
}

func (this *categoryC) EditCategory(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	r.ParseForm()
	id, _ := strconv.Atoi(r.Form.Get("id"))
	var category *sale.ValueCategory = dps.SaleService.GetCategory(partnerId, id)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	json, _ := json.Marshal(category)

	re := regexp.MustCompile(fmt.Sprintf("<option class=\"opt\\d+\" value=\"%d\">[^>]+>", id))
	originOpts := cache.GetDropOptionsOfCategory(partnerId)
	cateOpts := re.ReplaceAll(originOpts, nil)

<<<<<<< HEAD
	d := echox.NewRenderData()
	d.Map = map[string]interface{}{
		"entity":    template.JS(json),
		"cate_opts": template.HTML(cateOpts),
	}
	return ctx.RenderOK("category.edit.html", d)
}

//修改门店信息
func (this *categoryC) SaveCategory(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	if r.Method == "POST" {
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
	return nil
}

func (this *categoryC) DelCategory(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.Request()
	if r.Method == "POST" {
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
	return nil
=======
	ctx.App.Template().Execute(w,
		gof.TemplateDataMap{
			"entity":    template.JS(json),
			"cate_opts": template.HTML(cateOpts),
		},
		"views/partner/category/category_edit.html")
}

//修改门店信息
func (this *categoryC) SaveCategory_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
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
	w.Write(result.Marshal())
}

func (this *categoryC) DelCategory_post(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.Response
	var result gof.Message
	r.ParseForm()
	categoryId, _ := strconv.Atoi(r.Form.Get("id"))

	//删除分类
	err := dps.SaleService.DeleteCategory(partnerId, categoryId)
	if err != nil {
		result = gof.Message{Result: false, Message: err.Error()}
	} else {
		result = gof.Message{Result: true, Message: ""}
	}
	w.Write(result.Marshal())
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
