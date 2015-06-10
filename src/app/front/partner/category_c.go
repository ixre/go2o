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
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"github.com/atnet/gof/web/ui/tree"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"html/template"
	"regexp"
	"strconv"
)

var _ mvc.Filter = new(categoryC)

type categoryC struct {
	*baseC
}

//分类树形功能
func (this *categoryC) Category(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.ResponseWriter, gof.TemplateDataMap{
		"no_pic_url": format.GetGoodsImageUrl(""),
	}, "views/partner/category/category.html")
}

//分类Json数据
func (this *categoryC) CategoryJson(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	var node *tree.TreeNode = dps.SaleService.GetCategoryTreeNode(partnerId)
	json, _ := json.Marshal(node)
	ctx.ResponseWriter.Write(json)
}

//分类树形功能
func (this *categoryC) CategorySelect(ctx *web.Context) {
	ctx.App.Template().Execute(ctx.ResponseWriter, nil,
		"views/partner/category/category_select.html")
}

//分类Json数据
func (this *categoryC) CreateCategory(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)

	var node *tree.TreeNode = dps.SaleService.GetCategoryTreeNode(partnerId)
	json, _ := json.Marshal(node)

	ctx.App.Template().Execute(ctx.ResponseWriter,
		gof.TemplateDataMap{
			"treeJson": template.JS(json),
		},
		"views/partner/category/category_create.html")
}

func (this *categoryC) EditCategory(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
	r.ParseForm()
	id, _ := strconv.Atoi(r.Form.Get("id"))
	var category *sale.ValueCategory = dps.SaleService.GetCategory(partnerId, id)
	json, _ := json.Marshal(category)

	re := regexp.MustCompile(fmt.Sprintf("<option class=\"opt\\d+\" value=\"%d\">[^>]+>", id))
	originOpts := cache.GetDropOptionsOfCategory(partnerId)
	cateOpts := re.ReplaceAll(originOpts,nil)


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
	r, w := ctx.Request, ctx.ResponseWriter
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
	r, w := ctx.Request, ctx.ResponseWriter
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
}
