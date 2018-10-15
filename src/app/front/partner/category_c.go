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
	//"github.com/jsix/gof/web"
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
	"github.com/jsix/gof/web/form"
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
	partnerId := getPartnerId(ctx)
	var node *tree.TreeNode = dps.SaleService.GetCategoryTreeNode(partnerId)
	return ctx.JSON(http.StatusOK, node)
}

//分类树形功能
func (this *categoryC) CategorySelect(ctx *echox.Context) error {
	d := ctx.NewData()
	return ctx.RenderOK("category.select.html", d)
}

//分类Json数据
func (this *categoryC) CreateCategory(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)

	cateOpts := cache.GetDropOptionsOfCategory(partnerId)

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
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	r.ParseForm()
	id, _ := strconv.Atoi(r.Form.Get("id"))
	e, _ := dps.SaleService.GetCategory(partnerId, id)
	eJson, _ := json.Marshal(e)

	re := regexp.MustCompile(fmt.Sprintf("<option class=\"opt\\d+\" value=\"%d\">[^>]+>", id))
	originOpts := cache.GetDropOptionsOfCategory(partnerId)
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
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		var result gof.Result
		r.ParseForm()

		e := sale.ValueCategory{}
		form.ParseEntity(r.Form, &e)

		id, err := dps.SaleService.SaveCategory(partnerId, &e)
		if err != nil {
			result = gof.Result{ErrCode: 1, ErrMsg: err.Error()}
		} else {
			var data = make(map[string]string)
			data["id"] = fmt.Sprintf("%d", id)
			result = gof.Result{ErrCode: 0, ErrMsg: "", Data: data}
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}

func (this *categoryC) DelCategory(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	r := ctx.HttpRequest()
	if r.Method == "POST" {
		var result gof.Result
		r.ParseForm()
		categoryId, _ := strconv.Atoi(r.FormValue("id"))

		//删除分类
		err := dps.SaleService.DeleteCategory(partnerId, categoryId)
		if err != nil {
			result = gof.Result{ErrCode: 1, ErrMsg: err.Error()}
		} else {
			result = gof.Result{ErrCode: 0, ErrMsg: ""}
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return nil
}
