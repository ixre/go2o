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
	"github.com/atnet/gof/web/ui/tree"
	"go2o/core/domain/interface/sale"
	"go2o/core/infrastructure/format"
	"go2o/core/service/dps"
	"html/template"
	"net/http"
	"strconv"
)

type categoryC struct {
	app.Context
}

//分类树形功能
func (this *categoryC) Category(w http.ResponseWriter, r *http.Request) {

	this.Context.Template().Execute(w, func(m *map[string]interface{}) {
		(*m)["nopicUrl"] = format.GetGoodsImageUrl("")
	}, "views/partner/category/category.html")
}

//分类Json数据
func (c *categoryC) CategoryJson(w http.ResponseWriter, r *http.Request, ptId int) {
	var node *tree.TreeNode = dps.SaleService.GetCategoryTreeNode(ptId)
	json, _ := json.Marshal(node)
	w.Write(json)
}

//分类树形功能
func (this *categoryC) CategorySelect(w http.ResponseWriter, r *http.Request) {
	this.Context.Template().Render(w,
		"views/partner/category/category_select.html",
		nil)
}

//分类Json数据
func (this *categoryC) CreateCategory(w http.ResponseWriter, r *http.Request, ptId int) {

	var node *tree.TreeNode = dps.SaleService.GetCategoryTreeNode(ptId)
	json, _ := json.Marshal(node)

	this.Context.Template().Render(w,
		"views/partner/category/category_create.html",
		func(m *map[string]interface{}) {
			(*m)["treeJson"] = template.JS(json)
		})

}

func (this *categoryC) EditCategory(w http.ResponseWriter, r *http.Request, ptId int) {

	r.ParseForm()
	id, _ := strconv.Atoi(r.Form.Get("id"))
	var category *sale.ValueCategory = dps.SaleService.GetCategory(ptId, id)
	//fmt.Println(category)
	json, _ := json.Marshal(category)

	this.Context.Template().Render(w,
		"views/partner/category/category_edit.html",
		func(m *map[string]interface{}) {
			(*m)["entity"] = template.JS(json)
		})

}

//修改门店信息
func (this *categoryC) SaveCategory_post(w http.ResponseWriter, r *http.Request, ptId int) {
	var result gof.JsonResult
	r.ParseForm()

	e := sale.ValueCategory{}
	web.ParseFormToEntity(r.Form, &e)

	id, err := dps.SaleService.SaveCategory(ptId, &e)
	if err != nil {
		result = gof.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true, Message: "", Data: id}
	}
	w.Write(result.Marshal())
}

func (this *categoryC) DelCategory_post(w http.ResponseWriter, r *http.Request, ptId int) {
	var result gof.JsonResult
	r.ParseForm()
	categoryId, _ := strconv.Atoi(r.Form.Get("id"))

	//删除分类
	err := dps.SaleService.DeleteCategory(ptId, categoryId)

	//id, err := dao.SaveCategory(&entity)
	if err != nil {
		result = gof.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true, Message: ""}
	}
	w.Write(result.Marshal())
}
